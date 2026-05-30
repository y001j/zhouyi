package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// 访问码一年有效
const accessCodeTTL = 365 * 24 * time.Hour

// AccessCode 访问码（支持多次使用：MaxUses=0 视为 1，等价于一次性码）
type AccessCode struct {
	Code      string     `json:"code"`
	CreatedAt time.Time  `json:"createdAt"`
	ExpiresAt time.Time  `json:"expiresAt"`
	MaxUses   int        `json:"maxUses,omitempty"`   // 总可用次数；0 表示历史一次性码
	UsedCount int        `json:"usedCount,omitempty"` // 已使用次数
	UsedAt    *time.Time `json:"usedAt,omitempty"`    // 最后一次使用时间；MaxUses 用尽时与"已使用"等价
	UsedBy    string     `json:"usedBy,omitempty"`    // 最后一次使用 IP
	Note      string     `json:"note,omitempty"`      // 管理员备注

	// reserved: 该码已被某请求预留，请求处理完成（commit/release）前其他请求不可用。
	// 不持久化（每次启动重置，避免崩溃后码被永久卡住）。
	reserved bool `json:"-"`
}

// effectiveMaxUses 兼容历史数据：MaxUses=0 视为 1（一次性码）
func (c *AccessCode) effectiveMaxUses() int {
	if c.MaxUses <= 0 {
		return 1
	}
	return c.MaxUses
}

// remainingUses 剩余可用次数
func (c *AccessCode) remainingUses() int {
	r := c.effectiveMaxUses() - c.UsedCount
	if r < 0 {
		return 0
	}
	return r
}

// exhausted 是否已用尽
func (c *AccessCode) exhausted() bool {
	return c.remainingUses() == 0
}

// authStore 管理访问码与管理员 session
type authStore struct {
	mu       sync.Mutex
	path     string
	codes    map[string]*AccessCode
	sessions map[string]time.Time // token -> 过期时间
	password string
}

const sessionTTL = 12 * time.Hour

var store *authStore

// readPasswordFromConfig 从 ./config.json 或 ./admin.conf 读取密码
// config.json 格式：{"adminPassword": "..."}
// admin.conf 格式：纯文本（一行），整行作为密码
func readPasswordFromConfig() string {
	// 优先 JSON 配置
	if b, err := os.ReadFile("./config.json"); err == nil {
		var cfg struct {
			AdminPassword string `json:"adminPassword"`
		}
		if json.Unmarshal(b, &cfg) == nil && cfg.AdminPassword != "" {
			return cfg.AdminPassword
		}
	}
	// 其次纯文本
	if b, err := os.ReadFile("./admin.conf"); err == nil {
		return strings.TrimSpace(string(b))
	}
	return ""
}

func initAuth(filePath string) error {
	pw := os.Getenv("ADMIN_PASSWORD")
	if pw == "" {
		pw = readPasswordFromConfig()
	}
	if pw == "" {
		return errors.New("管理员密码未设置：请通过环境变量 ADMIN_PASSWORD 或配置文件 ./config.json（{\"adminPassword\":\"...\"}）提供")
	}
	store = &authStore{
		path:     filePath,
		codes:    map[string]*AccessCode{},
		sessions: map[string]time.Time{},
		password: pw,
	}
	return store.load()
}

func (s *authStore) load() error {
	b, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var list []*AccessCode
	if err := json.Unmarshal(b, &list); err != nil {
		return err
	}
	for _, c := range list {
		s.codes[c.Code] = c
	}
	return nil
}

// caller 必须持有 mu
func (s *authStore) saveLocked() error {
	list := make([]*AccessCode, 0, len(s.codes))
	for _, c := range s.codes {
		list = append(list, c)
	}
	b, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func randomToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// ---------- 管理员 session ----------

func (s *authStore) login(password string) (string, error) {
	if password != s.password {
		return "", errors.New("管理员密码错误")
	}
	tok, err := randomToken(24)
	if err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cleanupSessionsLocked()
	s.sessions[tok] = time.Now().Add(sessionTTL)
	return tok, nil
}

func (s *authStore) cleanupSessionsLocked() {
	now := time.Now()
	for t, exp := range s.sessions {
		if exp.Before(now) {
			delete(s.sessions, t)
		}
	}
}

func (s *authStore) checkSession(tok string) bool {
	if tok == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	exp, ok := s.sessions[tok]
	if !ok {
		return false
	}
	if exp.Before(time.Now()) {
		delete(s.sessions, tok)
		return false
	}
	return true
}

// ---------- 访问码 ----------

// 生成 n 个新访问码，每个码可用 maxUses 次
func (s *authStore) generateCodes(n int, maxUses int, note string) ([]*AccessCode, error) {
	if n <= 0 {
		n = 1
	}
	if n > 50 {
		n = 50
	}
	if maxUses <= 0 {
		maxUses = 1
	}
	if maxUses > 9999 {
		maxUses = 9999
	}
	out := make([]*AccessCode, 0, n)
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < n; i++ {
		var code string
		for {
			t, err := randomToken(4) // 8 hex chars
			if err != nil {
				return nil, err
			}
			code = strings.ToUpper(t)
			if _, dup := s.codes[code]; !dup {
				break
			}
		}
		c := &AccessCode{
			Code:      code,
			CreatedAt: now,
			ExpiresAt: now.Add(accessCodeTTL),
			MaxUses:   maxUses,
			Note:      note,
		}
		s.codes[code] = c
		out = append(out, c)
	}
	if err := s.saveLocked(); err != nil {
		return nil, err
	}
	return out, nil
}

// 列出全部访问码（按创建时间倒序由调用方处理）
func (s *authStore) listCodes() []*AccessCode {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*AccessCode, 0, len(s.codes))
	for _, c := range s.codes {
		out = append(out, c)
	}
	return out
}

// deleteCodesByScope 按范围删除访问码：scope 取值 "used" / "expired"。
// 已被 reserve 的码（正在被请求消耗）不会被删除，避免与并发请求竞争。
// 返回实际删除数量。
func (s *authStore) deleteCodesByScope(scope string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	victims := make([]string, 0)
	for code, c := range s.codes {
		if c.reserved {
			continue
		}
		switch scope {
		case "used":
			if c.exhausted() {
				victims = append(victims, code)
			}
		case "expired":
			// 已过期：未用尽 且 已过期；用尽的归"已使用"
			if !c.exhausted() && now.After(c.ExpiresAt) {
				victims = append(victims, code)
			}
		default:
			return 0, errors.New("不支持的删除范围")
		}
	}
	if len(victims) == 0 {
		return 0, nil
	}
	// 备份被删除项以便落盘失败时回滚
	backup := make(map[string]*AccessCode, len(victims))
	for _, code := range victims {
		backup[code] = s.codes[code]
		delete(s.codes, code)
	}
	if err := s.saveLocked(); err != nil {
		for code, c := range backup {
			s.codes[code] = c
		}
		return 0, err
	}
	return len(victims), nil
}

// reserveCode 原子地校验并预留访问码（不写盘）。
// 返回 normalized code + nil 表示预留成功；调用方必须随后调用 commitCode 或 releaseCode。
// 同一时刻只允许一个请求预留一个码（互斥）。
func (s *authStore) reserveCode(code string) (string, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return "", errors.New("缺少访问码")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.codes[code]
	if !ok {
		return "", errors.New("访问码无效")
	}
	if c.exhausted() {
		return "", errors.New("访问码已被使用")
	}
	if c.reserved {
		return "", errors.New("访问码正在使用中")
	}
	if time.Now().After(c.ExpiresAt) {
		return "", errors.New("访问码已过期")
	}
	c.reserved = true
	return code, nil
}

// commitCode 落盘消费已预留的码（次数 +1）。仅在请求成功（2xx）时调用。
func (s *authStore) commitCode(code, ip string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.codes[code]
	if !ok {
		return errors.New("访问码不存在")
	}
	prevUsedAt := c.UsedAt
	prevUsedBy := c.UsedBy
	prevCount := c.UsedCount
	now := time.Now()
	c.UsedCount++
	c.UsedAt = &now
	c.UsedBy = ip
	c.reserved = false
	if err := s.saveLocked(); err != nil {
		// 落盘失败：回滚到预留前状态，让用户可再试
		c.UsedCount = prevCount
		c.UsedAt = prevUsedAt
		c.UsedBy = prevUsedBy
		c.reserved = false
		return err
	}
	return nil
}

// releaseCode 释放已预留但未消费的码（业务失败、handler 异常等）。
func (s *authStore) releaseCode(code string) {
	if code == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok := s.codes[code]; ok {
		c.reserved = false
	}
}

// ---------- HTTP 处理 ----------

func clientIP(r *http.Request) string {
	if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
		if i := strings.IndexByte(xf, ','); i >= 0 {
			return strings.TrimSpace(xf[:i])
		}
		return strings.TrimSpace(xf)
	}
	// r.RemoteAddr 形如 "1.2.3.4:54321" 或 "[::1]:54321"，剥掉端口
	addr := r.RemoteAddr
	if host, _, err := net.SplitHostPort(addr); err == nil {
		return host
	}
	return addr
}

const adminCookieName = "zy_admin"

func handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "请使用 POST")
		return
	}
	ip := clientIP(r)
	if !rateLimitAllow("login", ip) {
		writeRateLimited(w, "登录尝试过多，请 10 分钟后再试")
		return
	}
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rateLimitRecordFailure("login", ip)
		writeError(w, http.StatusBadRequest, "请求体解析失败")
		return
	}
	tok, err := store.login(req.Password)
	if err != nil {
		rateLimitRecordFailure("login", ip)
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     adminCookieName,
		Value:    tok,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessionTTL.Seconds()),
	})
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func handleAdminLogout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(adminCookieName); err == nil {
		store.mu.Lock()
		delete(store.sessions, c.Value)
		store.mu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{
		Name: adminCookieName, Value: "", Path: "/", MaxAge: -1,
	})
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func adminAuthed(r *http.Request) bool {
	c, err := r.Cookie(adminCookieName)
	if err != nil {
		return false
	}
	return store.checkSession(c.Value)
}

func handleAdminCodes(w http.ResponseWriter, r *http.Request) {
	if !adminAuthed(r) {
		writeError(w, http.StatusUnauthorized, "请先登录管理员")
		return
	}
	switch r.Method {
	case http.MethodGet:
		list := store.listCodes()
		// 序列化时附加状态与剩余次数
		type item struct {
			*AccessCode
			Status         string `json:"status"`
			EffectiveMax   int    `json:"effectiveMaxUses"`
			RemainingUses  int    `json:"remainingUses"`
		}
		now := time.Now()
		out := make([]item, 0, len(list))
		for _, c := range list {
			st := "未使用"
			if c.exhausted() {
				st = "已使用"
			} else if now.After(c.ExpiresAt) {
				st = "已过期"
			} else if c.UsedCount > 0 {
				st = "使用中" // 已部分消耗、尚未用尽且未过期
			}
			out = append(out, item{
				AccessCode:    c,
				Status:        st,
				EffectiveMax:  c.effectiveMaxUses(),
				RemainingUses: c.remainingUses(),
			})
		}
		writeJSON(w, http.StatusOK, out)
	case http.MethodPost:
		var req struct {
			Count   int    `json:"count"`
			MaxUses int    `json:"maxUses"`
			Note    string `json:"note"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Count <= 0 {
			req.Count = 1
		}
		if req.MaxUses <= 0 {
			req.MaxUses = 1
		}
		codes, err := store.generateCodes(req.Count, req.MaxUses, req.Note)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, codes)
	case http.MethodDelete:
		var req struct {
			Scope string `json:"scope"` // "used" | "expired"
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Scope == "" {
			req.Scope = r.URL.Query().Get("scope")
		}
		n, err := store.deleteCodesByScope(req.Scope)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": n})
	default:
		writeError(w, http.StatusMethodNotAllowed, "只支持 GET / POST / DELETE")
	}
}

// statusRecorder 包装 ResponseWriter 以拦截状态码。
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (sr *statusRecorder) WriteHeader(code int) {
	if !sr.wroteHeader {
		sr.status = code
		sr.wroteHeader = true
	}
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	if !sr.wroteHeader {
		sr.status = http.StatusOK
		sr.wroteHeader = true
	}
	return sr.ResponseWriter.Write(b)
}

// requireCode 起卦中间件（两步消费）：
// 1) 管理员 cookie 直通
// 2) 否则先 reserve 访问码（原子互斥），handler 处理完后：
//    - 状态码 2xx → commit（永久消费）
//    - 否则 → release（释放给用户重试）
func requireCode(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if adminAuthed(r) {
			next(w, r)
			return
		}
		ip := clientIP(r)
		if !rateLimitAllow("code", ip) {
			writeRateLimited(w, "访问码尝试过多，请稍后再试")
			return
		}
		code := r.Header.Get("X-Access-Code")
		if code == "" {
			code = r.URL.Query().Get("code")
		}
		normalized, err := store.reserveCode(code)
		if err != nil {
			rateLimitRecordFailure("code", ip)
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		sr := &statusRecorder{ResponseWriter: w}
		// 用 panic recover 兜底：handler panic 时也要释放码
		committed := false
		defer func() {
			if rec := recover(); rec != nil {
				store.releaseCode(normalized)
				// 让 panic 继续，由 net/http 默认处理
				panic(rec)
			}
			if committed {
				return
			}
			if sr.status >= 200 && sr.status < 300 {
				if err := store.commitCode(normalized, ip); err != nil {
					// commit 失败仅记录，码已回滚为可用
					return
				}
				committed = true
			} else {
				store.releaseCode(normalized)
			}
		}()
		next(sr, r)
	}
}
