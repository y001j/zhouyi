package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"zhouyi/liuren"
	"zhouyi/qimen"
)

//go:embed web/*
var webFS embed.FS

// ===== 请求/响应类型 =====

type divineRequest struct {
	Method       string `json:"method"`       // coin | yarrow | number
	Upper        int    `json:"upper"`        // number 法：上卦数
	Lower        int    `json:"lower"`        // number 法：下卦数
	Changing     int    `json:"changing"`     // number 法：变爻位置 0-6
	Question     string  `json:"question"`     // 所问之事
	QuestionType string  `json:"questionType"` // career/wealth/...
	Longitude    float64 `json:"longitude"`    // 真太阳时校正用经度（东经为正），0/越界视为未提供
}

type lineView struct {
	Pos      int    `json:"pos"`      // 1-6
	Value    int    `json:"value"`    // 6/7/8/9
	IsYang   bool   `json:"isYang"`
	IsChange bool   `json:"isChange"`
	Symbol   string `json:"symbol"`
	TypeName string `json:"typeName"` // 老阴/少阳/少阴/老阳
	LineName string `json:"lineName"` // 初九 九二 ...
	Text     string `json:"text"`     // 本卦该爻爻辞
}

type trigramView struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Nature string `json:"nature"`
}

type hexagramView struct {
	Number   int         `json:"number"`
	Name     string      `json:"name"`
	Symbol   string      `json:"symbol"`
	Upper    trigramView `json:"upper"`
	Lower    trigramView `json:"lower"`
	Judgment string      `json:"judgment"`
	Image    string      `json:"image"`
	Lines    []struct {
		Position int    `json:"position"`
		Type     string `json:"type"`
		Text     string `json:"text"`
	} `json:"lines"`
}

type derivedView struct {
	Mutual   *hexagramView `json:"mutual,omitempty"`
	Opposite *hexagramView `json:"opposite,omitempty"`
	Reverse  *hexagramView `json:"reverse,omitempty"`
}

type timingView struct {
	SolarTime      string        `json:"solarTime"`
	LunarDesc      string        `json:"lunarDesc"`
	GanZhiSummary  string        `json:"ganZhiSummary"`
	JieQiName      string        `json:"jieQiName"`
	JieQiDay       string        `json:"jieQiDay"`
	NextJieQiName  string        `json:"nextJieQiName"`
	NextJieQiDate  string        `json:"nextJieQiDate"`
	MonthlyHex     *hexagramView `json:"monthlyHex,omitempty"`
	MonthlyHexNote string        `json:"monthlyHexNote"`
}

type divineResponse struct {
	Method        string        `json:"method"`
	MethodLabel   string        `json:"methodLabel"`
	Time          string        `json:"time"`
	Question      string        `json:"question"`
	QuestionType  string        `json:"questionType"`
	QuestionLabel string        `json:"questionLabel"`
	Lines         []lineView    `json:"lines"`
	ChangingPos   []int         `json:"changingPos"`
	MainHex       *hexagramView `json:"mainHex"`
	ChangeHex     *hexagramView `json:"changeHex,omitempty"`
	Derived       derivedView   `json:"derived"`
	Timing        *timingView   `json:"timing,omitempty"`
	Guide         string        `json:"guide"`         // 解卦指引文字
	Prompt        string        `json:"prompt"`        // AI 提示词
	Interpret     string        `json:"interpret"`     // 终端风格的完整解卦文字
}

// ===== 类型转换辅助 =====

func toTrigramView(name string) trigramView {
	t, ok := Trigrams[name]
	if !ok {
		return trigramView{Name: name}
	}
	return trigramView{Name: t.Name, Symbol: t.Symbol, Nature: t.Nature}
}

func toHexagramView(h *Hexagram) *hexagramView {
	if h == nil {
		return nil
	}
	v := &hexagramView{
		Number:   h.Number,
		Name:     h.Name,
		Symbol:   h.Symbol,
		Upper:    toTrigramView(h.Upper),
		Lower:    toTrigramView(h.Lower),
		Judgment: h.Judgment,
		Image:    h.Image,
	}
	for _, ln := range h.Lines {
		v.Lines = append(v.Lines, struct {
			Position int    `json:"position"`
			Type     string `json:"type"`
			Text     string `json:"text"`
		}{ln.Position, ln.Type, ln.Text})
	}
	return v
}

func toLineViews(r DivinationResult) []lineView {
	typeNames := map[int]string{6: "老阴（变）", 7: "少阳", 8: "少阴", 9: "老阳（变）"}
	out := make([]lineView, 6)
	for i := 0; i < 6; i++ {
		li := lineInfo(r.Lines[i])
		yy := "九"
		if !li.IsYang {
			yy = "六"
		}
		posName := []string{"初", "二", "三", "四", "五", "上"}[i]
		text := ""
		if r.MainHex != nil {
			text = r.MainHex.Lines[i].Text
		}
		out[i] = lineView{
			Pos:      i + 1,
			Value:    r.Lines[i],
			IsYang:   li.IsYang,
			IsChange: li.IsChange,
			Symbol:   li.Symbol,
			TypeName: typeNames[r.Lines[i]],
			LineName: posName + yy,
			Text:     text,
		}
	}
	return out
}

func toTimingView(ti *TimingInfo) *timingView {
	if ti == nil {
		return nil
	}
	return &timingView{
		SolarTime:      ti.SolarTime.Format("2006-01-02 15:04"),
		LunarDesc:      ti.LunarDesc,
		GanZhiSummary:  ti.GanZhiSummary,
		JieQiName:      ti.JieQiName,
		JieQiDay:       ti.JieQiDay,
		NextJieQiName:  ti.NextJieQiName,
		NextJieQiDate:  ti.NextJieQiDate,
		MonthlyHex:     toHexagramView(ti.MonthlyHex),
		MonthlyHexNote: ti.MonthlyHexNote,
	}
}

func methodLabel(m DivinationMethod) string {
	return map[DivinationMethod]string{
		CoinMethod:   "铜钱法（金钱卦）",
		YarrowMethod: "蓍草法（大衍揲蓍）",
		NumberMethod: "数字起卦法",
	}[m]
}

// ===== HTTP 处理 =====

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func handleDivine(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "请使用 POST")
		return
	}
	var req divineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体解析失败: "+err.Error())
		return
	}

	var result DivinationResult
	switch strings.ToLower(strings.TrimSpace(req.Method)) {
	case "", "coin":
		result = DivineByCoins()
	case "yarrow":
		result = DivineByYarrow()
	case "number":
		result = DivineByNumber(req.Upper, req.Lower, req.Changing)
	default:
		writeError(w, http.StatusBadRequest, "未知的起卦方式："+req.Method)
		return
	}

	qt := ParseQuestionType(req.QuestionType)
	result.QuestionType = qt
	if result.MainHex == nil {
		writeError(w, http.StatusInternalServerError, "起卦失败")
		return
	}

	// 真太阳时校正：把 result.Time（钟表时）映射到本地真太阳平时
	result.Time = ApplyTrueSolarTime(result.Time, ResolveLongitude(req.Longitude))

	derived := DeriveHexagrams(result.Lines)
	timing := CaptureTiming(result.Time)

	resp := divineResponse{
		Method:        req.Method,
		MethodLabel:   methodLabel(result.Method),
		Time:          result.Time.Format("2006-01-02 15:04:05"),
		Question:      req.Question,
		QuestionType:  string(qt),
		QuestionLabel: QuestionTypeLabel(qt),
		Lines:         toLineViews(result),
		ChangingPos:   result.ChangingPos,
		MainHex:       toHexagramView(result.MainHex),
		ChangeHex:     toHexagramView(result.ChangeHex),
		Derived: derivedView{
			Mutual:   toHexagramView(derived.Mutual),
			Opposite: toHexagramView(derived.Opposite),
			Reverse:  toHexagramView(derived.Reverse),
		},
		Timing:    toTimingView(timing),
		Guide:     interpretationGuide(result),
		Prompt:    GenerateAIPrompt(result, req.Question),
		Interpret: InterpretResult(result),
	}
	writeJSON(w, http.StatusOK, resp)
}

func handleHexagramList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "请使用 GET")
		return
	}
	type item struct {
		Number int    `json:"number"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
		Upper  string `json:"upper"`
		Lower  string `json:"lower"`
	}
	list := make([]item, 0, len(Hexagrams))
	for _, h := range Hexagrams {
		list = append(list, item{h.Number, h.Name, h.Symbol, h.Upper, h.Lower})
	}
	writeJSON(w, http.StatusOK, list)
}

func handleHexagramDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "请使用 GET")
		return
	}
	nStr := r.URL.Query().Get("n")
	n, err := strconv.Atoi(nStr)
	if err != nil || n < 1 || n > 64 {
		writeError(w, http.StatusBadRequest, "请提供 1-64 之间的卦号 n")
		return
	}
	writeJSON(w, http.StatusOK, toHexagramView(&Hexagrams[n-1]))
}

// ===== 大六壬 API =====

type liurenRequest struct {
	Question     string  `json:"question"`
	QuestionType string  `json:"questionType"`
	TimeISO      string  `json:"time"`      // 可选：ISO 8601；留空取服务器当前时间
	BenMing      string  `json:"benMing"`   // 可选：本命（生肖或地支）
	BirthYear    int     `json:"birthYear"` // 可选：出生公历年
	Gender       string  `json:"gender"`    // 可选：男/女
	Longitude    float64 `json:"longitude"` // 真太阳时校正用经度
}

type liurenKeView struct {
	Index    int    `json:"index"`
	Upper    string `json:"upper"`
	Lower    string `json:"lower"`
	Relation string `json:"relation"`
}

type liurenChuanView struct {
	Name      string `json:"name"`
	Zhi       string `json:"zhi"`
	TianJiang string `json:"tianJiang"`
	LiuQin    string `json:"liuQin"`
	IsKong    bool   `json:"isKong"`
}

type shenShaView struct {
	Name string `json:"name"`
	Zhi  string `json:"zhi"`
	Desc string `json:"desc"`
}

type ketiView struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
}

type nianMingView struct {
	Zhi       string `json:"zhi"`
	Upper     string `json:"upper"`
	TianJiang string `json:"tianJiang"`
	LiuQin    string `json:"liuQin"`
	IsKong    bool   `json:"isKong"`
	Ying      string `json:"ying"`
	YingDesc  string `json:"yingDesc"`
}

type biFaView struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Note   string `json:"note"`
}

type tabooView struct {
	TianJiang string `json:"tianJiang"`
	Zhi       string `json:"zhi"`
	Note      string `json:"note"`
}

type liurenPanView struct {
	Time          string            `json:"time"`
	DayGan        string            `json:"dayGan"`
	DayZhi        string            `json:"dayZhi"`
	JiaziIndex    int               `json:"jiaziIndex"`
	ZhanShi       string            `json:"zhanShi"`
	YueJiang      string            `json:"yueJiang"`
	YueJiangAlt   string            `json:"yueJiangAlt"`
	QiName        string            `json:"qiName"`
	ZhouYe        string            `json:"zhouYe"`
	XunKong       [2]string         `json:"xunKong"`
	DiPan         [12]string        `json:"diPan"`
	TianPan       [12]string        `json:"tianPan"`
	TianJiang     [12]string        `json:"tianJiang"`
	TianJiangSrt  [12]string        `json:"tianJiangShort"`
	SiKe          [4]liurenKeView   `json:"siKe"`
	SanChuan      [3]liurenChuanView `json:"sanChuan"`
	Method        string            `json:"method"`
	KeTiName      string            `json:"keTiName"`
	KeTiSummary   string            `json:"keTiSummary"`
	Tags          []ketiView        `json:"tags"`
	ShenSha       []shenShaView     `json:"shenSha"`
	BenMing       *nianMingView     `json:"benMing,omitempty"`
	XingNian      *nianMingView     `json:"xingNian,omitempty"`
	BMXNRel       string            `json:"bmxnRel,omitempty"`
	BMXNDesc      string            `json:"bmxnDesc,omitempty"`
	BiFa          []biFaView        `json:"biFa"`
	BiFaCatalog   []biFaView        `json:"biFaCatalog"`
	Taboos        []tabooView       `json:"taboos"`
	Guide         string            `json:"guide"`
	Interpret     string            `json:"interpret"`
	Prompt        string            `json:"prompt"`
	QuestionLabel string            `json:"questionLabel"`
}

func handleLiuRen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "请使用 POST")
		return
	}
	var req liurenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体解析失败: "+err.Error())
		return
	}
	t := time.Now()
	if req.TimeISO != "" {
		if parsed, err := time.Parse(time.RFC3339, req.TimeISO); err == nil {
			t = parsed
		}
	}
	t = ApplyTrueSolarTime(t, ResolveLongitude(req.Longitude))
	ctx, err := liuren.BuildContext(t)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "起课失败: "+err.Error())
		return
	}
	if req.BenMing != "" {
		if z := parseShengXiao(req.BenMing); z >= 0 {
			ctx.BenMing = &z
		}
	}
	if req.BirthYear > 0 {
		ctx.BirthYear = req.BirthYear
		ctx.Gender = req.Gender
	}
	qt := ParseQuestionType(req.QuestionType)
	ctx.QuestionType = string(qt)
	pan := liuren.DivineWithContext(ctx)
	focus := FocusGuideLiuRen(qt)
	leishen := LeiShenDirective(pan, qt)
	prompt := liuren.GenerateAIPrompt(pan, req.Question, focus, leishen)

	v := liurenPanView{
		Time:          pan.Ctx.Time.Format("2006-01-02 15:04:05"),
		DayGan:        pan.Ctx.Gan.String(),
		DayZhi:        pan.Ctx.DayZhi.String(),
		JiaziIndex:    pan.Ctx.JiaziIndex + 1,
		ZhanShi:       pan.Ctx.ZhanShi.String(),
		YueJiang:      pan.Ctx.YueJiang.String(),
		YueJiangAlt:   liuren.ZhiBieMing[pan.Ctx.YueJiang],
		QiName:        pan.Ctx.QiName,
		ZhouYe:        map[bool]string{true: "昼占", false: "夜占"}[pan.Ctx.ZhouYe],
		Method:        pan.SanChuan.Method,
		KeTiName:      pan.KeTi.Name,
		KeTiSummary:   pan.KeTi.Summary,
		Guide:         liuren.InterpretGuide(pan),
		Interpret:     liuren.Render(pan),
		Prompt:        prompt,
		QuestionLabel: QuestionTypeLabel(qt),
	}
	kong := pan.Ctx.XunKongPair()
	v.XunKong = [2]string{kong[0].String(), kong[1].String()}
	for i := 0; i < 12; i++ {
		v.DiPan[i] = liuren.Zhi(i).String()
		v.TianPan[i] = pan.TianPan[i].String()
		v.TianJiang[i] = pan.TianJiang[i].String()
		v.TianJiangSrt[i] = pan.TianJiang[i].Short()
	}
	for i, ke := range pan.SiKe {
		v.SiKe[i] = liurenKeView{Index: ke.Index, Upper: ke.Upper.String(), Lower: ke.Lower.String(), Relation: ke.Relation}
	}
	chuans := [3]liuren.ChuanEntry{pan.SanChuan.Chu, pan.SanChuan.Zhong, pan.SanChuan.Mo}
	for i, ce := range chuans {
		v.SanChuan[i] = liurenChuanView{
			Name: ce.Name, Zhi: ce.Zhi.String(),
			TianJiang: ce.TianJiang.String(), LiuQin: ce.LiuQin.String(),
			IsKong: ce.IsKong,
		}
	}
	// 附加课格
	for _, tg := range pan.Tags {
		v.Tags = append(v.Tags, ketiView{Name: tg.Name, Summary: tg.Summary})
	}
	// 神煞
	for _, ss := range pan.ShenSha {
		if ss.Zhi < 0 {
			continue
		}
		v.ShenSha = append(v.ShenSha, shenShaView{Name: ss.Name, Zhi: ss.Zhi.String(), Desc: ss.Desc})
	}
	// 年命/行年
	if pan.NianMing != nil {
		if pan.NianMing.BenMing != nil {
			bm := pan.NianMing.BenMing
			v.BenMing = &nianMingView{Zhi: bm.Zhi.String(), Upper: bm.Upper.String(),
				TianJiang: bm.TianJiang.String(), LiuQin: bm.LiuQin.String(), IsKong: bm.IsKong,
				Ying: bm.Ying.String(), YingDesc: bm.Ying.Desc()}
		}
		if pan.NianMing.XingNian != nil {
			xn := pan.NianMing.XingNian
			v.XingNian = &nianMingView{Zhi: xn.Zhi.String(), Upper: xn.Upper.String(),
				TianJiang: xn.TianJiang.String(), LiuQin: xn.LiuQin.String(), IsKong: xn.IsKong,
				Ying: xn.Ying.String(), YingDesc: xn.Ying.Desc()}
		}
		v.BMXNRel = pan.NianMing.BMXNRel
		v.BMXNDesc = pan.NianMing.BMXNDesc
	}
	// 毕法赋（命中）
	for _, bf := range pan.BiFa {
		v.BiFa = append(v.BiFa, biFaView{Number: bf.Number, Title: bf.Title, Text: bf.Text, Note: bf.Note})
	}
	// 毕法赋全文 100 条
	for _, bf := range liuren.BiFaCatalog() {
		v.BiFaCatalog = append(v.BiFaCatalog, biFaView{Number: bf.Number, Title: bf.Title, Text: bf.Text, Note: bf.Note})
	}
	// 天将乘临禁忌
	for _, tb := range pan.Taboos {
		v.Taboos = append(v.Taboos, tabooView{
			TianJiang: tb.TianJiang.String(),
			Zhi:       tb.DiZhi.String(),
			Note:      tb.Note,
		})
	}
	writeJSON(w, http.StatusOK, v)
}

// ===== 互参 API =====

type huCanRequest struct {
	Question     string  `json:"question"`
	QuestionType string  `json:"questionType"`
	TimeISO      string  `json:"time"`
	Longitude    float64 `json:"longitude"`
}

type huCanResponse struct {
	Time        string         `json:"time"`
	Question    string         `json:"question"`
	MainHex     *hexagramView  `json:"mainHex"`
	ChangeHex   *hexagramView  `json:"changeHex,omitempty"`
	ChangingPos []int          `json:"changingPos"`
	LiuRen      *liurenPanView `json:"liuren"` // 完整六壬盘视图
	Qimen       *qimenPanView  `json:"qimen"`  // 完整奇门盘视图（阶段 3 新增）
	Prompt      string         `json:"prompt"`
}

func handleHuCan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "请使用 POST")
		return
	}
	var req huCanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体解析失败: "+err.Error())
		return
	}
	t := time.Now()
	if req.TimeISO != "" {
		if parsed, err := time.Parse(time.RFC3339, req.TimeISO); err == nil {
			t = parsed
		}
	}
	t = ApplyTrueSolarTime(t, ResolveLongitude(req.Longitude))
	qt := ParseQuestionType(req.QuestionType)
	res, err := HuCanDivine(t, req.Question, qt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "互参起占失败: "+err.Error())
		return
	}
	resp := huCanResponse{
		Time:        t.Format("2006-01-02 15:04:05"),
		Question:    req.Question,
		MainHex:     toHexagramView(res.Zhouyi.MainHex),
		ChangeHex:   toHexagramView(res.Zhouyi.ChangeHex),
		ChangingPos: res.Zhouyi.ChangingPos,
		LiuRen:      buildLiuRenView(res.LiuRenPan),
		Qimen:       buildQimenView(res.QimenPan),
		Prompt:      HuCanPrompt(res),
	}
	writeJSON(w, http.StatusOK, resp)
}

// buildLiuRenView 从 liuren.Pan 构造响应视图（互参 API 用；单侧 API 另有更全的版本）
func buildLiuRenView(pan *liuren.Pan) *liurenPanView {
	if pan == nil {
		return nil
	}
	v := &liurenPanView{
		Time:        pan.Ctx.Time.Format("2006-01-02 15:04:05"),
		DayGan:      pan.Ctx.Gan.String(),
		DayZhi:      pan.Ctx.DayZhi.String(),
		JiaziIndex:  pan.Ctx.JiaziIndex + 1,
		ZhanShi:     pan.Ctx.ZhanShi.String(),
		YueJiang:    pan.Ctx.YueJiang.String(),
		YueJiangAlt: liuren.ZhiBieMing[pan.Ctx.YueJiang],
		QiName:      pan.Ctx.QiName,
		ZhouYe:      map[bool]string{true: "昼占", false: "夜占"}[pan.Ctx.ZhouYe],
		Method:      pan.SanChuan.Method,
		KeTiName:    pan.KeTi.Name,
		KeTiSummary: pan.KeTi.Summary,
	}
	kong := pan.Ctx.XunKongPair()
	v.XunKong = [2]string{kong[0].String(), kong[1].String()}
	for i := 0; i < 12; i++ {
		v.DiPan[i] = liuren.Zhi(i).String()
		v.TianPan[i] = pan.TianPan[i].String()
		v.TianJiang[i] = pan.TianJiang[i].String()
		v.TianJiangSrt[i] = pan.TianJiang[i].Short()
	}
	for i, ke := range pan.SiKe {
		v.SiKe[i] = liurenKeView{Index: ke.Index, Upper: ke.Upper.String(), Lower: ke.Lower.String(), Relation: ke.Relation}
	}
	chuans := [3]liuren.ChuanEntry{pan.SanChuan.Chu, pan.SanChuan.Zhong, pan.SanChuan.Mo}
	for i, ce := range chuans {
		v.SanChuan[i] = liurenChuanView{
			Name: ce.Name, Zhi: ce.Zhi.String(),
			TianJiang: ce.TianJiang.String(), LiuQin: ce.LiuQin.String(),
			IsKong: ce.IsKong,
		}
	}
	for _, tag := range pan.Tags {
		v.Tags = append(v.Tags, ketiView{Name: tag.Name, Summary: tag.Summary})
	}
	return v
}

// buildQimenView 从 qimen.Pan 构造响应视图
func buildQimenView(pan *qimen.Pan) *qimenPanView {
	if pan == nil {
		return nil
	}
	v := &qimenPanView{
		Time:         pan.Ctx.Time.Format("2006-01-02 15:04:05"),
		YearGZ:       pan.Ctx.YearGZ,
		MonthGZ:      pan.Ctx.MonthGZ,
		DayGZ:        pan.Ctx.DayGZ,
		HourGZ:       pan.Ctx.HourGZ,
		JieQi:        pan.Ctx.JieQi,
		Dun:          pan.Ctx.Dun,
		Yuan:         pan.Ctx.Yuan,
		Ju:           pan.Ctx.Ju,
		Xunshou:      pan.Ctx.Xunshou,
		Dungan:       pan.Ctx.Dungan,
		XunKong:      pan.Ctx.XunKong,
		YiMaZhi:      pan.YiMaZhi,
		ZhiFuStar:    pan.ZhiFuStar,
		ZhiFuPalace:  pan.ZhiFuPalace,
		ZhiShiGate:   pan.ZhiShiGate,
		ZhiShiPalace: pan.ZhiShiPalace,
		WangXiangXS:  pan.WangXiangXS,
		Render:       pan.Render(),
	}
	for i, c := range pan.Cells {
		v.Cells[i] = qimenCellView{
			PalaceFei: c.PalaceFei, PalaceName: c.PalaceName,
			EarthStem: c.EarthStem, HeavenStem: c.HeavenStem,
			Star: c.Star, Door: c.Door, God: c.God,
			Branches:      c.Branches,
			PalaceWuXing:  c.PalaceWuXing,
			HeavenWuXing:  c.HeavenWuXing,
			EarthWuXing:   c.EarthWuXing,
			StarWangShuai: c.StarWangShuai,
			DoorPalaceRel: c.DoorPalaceRel,
			IsDoorPo:      c.IsDoorPo, IsDoorSheng: c.IsDoorSheng,
			IsVoid: c.IsVoid, IsYima: c.IsYima,
			IsTianStemMu: c.IsTianStemMu, IsEarthStemMu: c.IsEarthStemMu,
			IsJiXing: c.IsJiXing,
		}
	}
	for _, h := range qimen.DetectPatterns(pan) {
		palName := ""
		if h.PalaceFei >= 0 && h.PalaceFei < 9 {
			palName = pan.Cells[h.PalaceFei].PalaceName
		}
		v.Patterns = append(v.Patterns, qimenPatternView{
			Name: h.Name, Category: h.Category,
			PalaceName: palName, Classic: h.Classic,
			Summary: h.Summary, AuspiceScore: h.AuspiceScore,
		})
	}
	return v
}

// parseShengXiao 主包版：字符串→大六壬地支
func parseShengXiao(s string) liuren.Zhi {
	s = strings.TrimSpace(s)
	m := map[string]liuren.Zhi{
		"鼠": liuren.Zi, "牛": liuren.Chou, "虎": liuren.Yin, "兔": liuren.Mao,
		"龙": liuren.Chen, "蛇": liuren.Si, "马": liuren.Wu, "羊": liuren.Wei,
		"猴": liuren.Shen, "鸡": liuren.You, "狗": liuren.Xu, "猪": liuren.Hai,
		"子": liuren.Zi, "丑": liuren.Chou, "寅": liuren.Yin, "卯": liuren.Mao,
		"辰": liuren.Chen, "巳": liuren.Si, "午": liuren.Wu, "未": liuren.Wei,
		"申": liuren.Shen, "酉": liuren.You, "戌": liuren.Xu, "亥": liuren.Hai,
	}
	if z, ok := m[s]; ok {
		return z
	}
	return -1
}

// ===== 奇门遁甲 API =====

type qimenRequest struct {
	Question     string  `json:"question"`
	QuestionType string  `json:"questionType"`
	TimeISO      string  `json:"timeISO,omitempty"`
	Longitude    float64 `json:"longitude"`
}

type qimenCellView struct {
	PalaceFei     int      `json:"palaceFei"`
	PalaceName    string   `json:"palaceName"`
	EarthStem     string   `json:"earthStem"`
	HeavenStem    string   `json:"heavenStem"`
	Star          string   `json:"star"`
	Door          string   `json:"door"`
	God           string   `json:"god"`
	Branches      []string `json:"branches"`
	PalaceWuXing  string   `json:"palaceWuXing"`
	HeavenWuXing  string   `json:"heavenWuXing"`
	EarthWuXing   string   `json:"earthWuXing"`
	StarWangShuai string   `json:"starWangShuai"`
	DoorPalaceRel string   `json:"doorPalaceRel"`
	IsDoorPo      bool     `json:"isDoorPo"`
	IsDoorSheng   bool     `json:"isDoorSheng"`
	IsVoid        bool     `json:"isVoid"`
	IsYima        bool     `json:"isYima"`
	IsTianStemMu  bool     `json:"isTianStemMu"`
	IsEarthStemMu bool     `json:"isEarthStemMu"`
	IsJiXing      bool     `json:"isJiXing"`
}

type qimenPatternView struct {
	Name         string `json:"name"`
	Category     string `json:"category"`
	PalaceName   string `json:"palaceName"`
	Classic      string `json:"classic"`
	Summary      string `json:"summary"`
	AuspiceScore int    `json:"auspiceScore"`
}

type qimenPanView struct {
	Time         string          `json:"time"`
	YearGZ       string          `json:"yearGZ"`
	MonthGZ      string          `json:"monthGZ"`
	DayGZ        string          `json:"dayGZ"`
	HourGZ       string          `json:"hourGZ"`
	JieQi        string          `json:"jieqi"`
	Dun          string          `json:"dun"`
	Yuan         string          `json:"yuan"`
	Ju           int             `json:"ju"`
	Xunshou      string          `json:"xunshou"`
	Dungan       string          `json:"dungan"`
	XunKong      [2]string       `json:"xunKong"`
	YiMaZhi      string          `json:"yiMaZhi"`
	ZhiFuStar    string          `json:"zhiFuStar"`
	ZhiFuPalace  string          `json:"zhiFuPalace"`
	ZhiShiGate   string          `json:"zhiShiGate"`
	ZhiShiPalace string          `json:"zhiShiPalace"`
	Cells        [9]qimenCellView `json:"cells"`
	WangXiangXS  [5]string       `json:"wangXiangXS"`
	Patterns     []qimenPatternView `json:"patterns"` // 命中格局
	Render       string          `json:"render"` // ASCII 盘面
	Prompt       string          `json:"prompt"` // AI 提示词
}

func handleQimen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "请使用 POST")
		return
	}
	var req qimenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "请求体解析失败: "+err.Error())
		return
	}
	t := time.Now()
	if req.TimeISO != "" {
		if parsed, err := time.Parse(time.RFC3339, req.TimeISO); err == nil {
			t = parsed
		}
	}
	t = ApplyTrueSolarTime(t, ResolveLongitude(req.Longitude))
	pan, err := qimen.BuildPan(t)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "起局失败: "+err.Error())
		return
	}
	qt := ParseQuestionType(req.QuestionType)
	// 用奇门自己的 FocusGuide（奇门化侧重语言，围绕星/门/神/干/宫断）
	focus := qimen.FocusGuide(string(qt))
	prompt := qimen.GenerateAIPrompt(pan, req.Question, focus, string(qt))

	v := qimenPanView{
		Time:         pan.Ctx.Time.Format("2006-01-02 15:04:05"),
		YearGZ:       pan.Ctx.YearGZ,
		MonthGZ:      pan.Ctx.MonthGZ,
		DayGZ:        pan.Ctx.DayGZ,
		HourGZ:       pan.Ctx.HourGZ,
		JieQi:        pan.Ctx.JieQi,
		Dun:          pan.Ctx.Dun,
		Yuan:         pan.Ctx.Yuan,
		Ju:           pan.Ctx.Ju,
		Xunshou:      pan.Ctx.Xunshou,
		Dungan:       pan.Ctx.Dungan,
		XunKong:      pan.Ctx.XunKong,
		YiMaZhi:      pan.YiMaZhi,
		ZhiFuStar:    pan.ZhiFuStar,
		ZhiFuPalace:  pan.ZhiFuPalace,
		ZhiShiGate:   pan.ZhiShiGate,
		ZhiShiPalace: pan.ZhiShiPalace,
		WangXiangXS:  pan.WangXiangXS,
		Render:       pan.Render(),
		Prompt:       prompt,
	}
	for i, c := range pan.Cells {
		v.Cells[i] = qimenCellView{
			PalaceFei:     c.PalaceFei,
			PalaceName:    c.PalaceName,
			EarthStem:     c.EarthStem,
			HeavenStem:    c.HeavenStem,
			Star:          c.Star,
			Door:          c.Door,
			God:           c.God,
			Branches:      c.Branches,
			PalaceWuXing:  c.PalaceWuXing,
			HeavenWuXing:  c.HeavenWuXing,
			EarthWuXing:   c.EarthWuXing,
			StarWangShuai: c.StarWangShuai,
			DoorPalaceRel: c.DoorPalaceRel,
			IsDoorPo:      c.IsDoorPo,
			IsDoorSheng:   c.IsDoorSheng,
			IsVoid:        c.IsVoid,
			IsYima:        c.IsYima,
			IsTianStemMu:  c.IsTianStemMu,
			IsEarthStemMu: c.IsEarthStemMu,
			IsJiXing:      c.IsJiXing,
		}
	}
	// 命中格局
	for _, h := range qimen.DetectPatterns(pan) {
		palName := ""
		if h.PalaceFei >= 0 && h.PalaceFei < 9 {
			palName = pan.Cells[h.PalaceFei].PalaceName
		}
		v.Patterns = append(v.Patterns, qimenPatternView{
			Name:         h.Name,
			Category:     h.Category,
			PalaceName:   palName,
			Classic:      h.Classic,
			Summary:      h.Summary,
			AuspiceScore: h.AuspiceScore,
		})
	}
	writeJSON(w, http.StatusOK, v)
}

func handleQuestionTypes(w http.ResponseWriter, r *http.Request) {
	type item struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}
	types := []QuestionType{QTCareer, QTWealth, QTRelation, QTHealth, QTDecision, QTTiming, QTOther}
	out := make([]item, 0, len(types))
	for _, t := range types {
		out = append(out, item{string(t), QuestionTypeLabel(t)})
	}
	writeJSON(w, http.StatusOK, out)
}

// runServer 启动 HTTP 服务
func runServer(addr string) {
	if err := initAuth("./codes.json"); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/divine", requireCode(handleDivine))
	mux.HandleFunc("/api/hexagrams", handleHexagramList)
	mux.HandleFunc("/api/hexagram", handleHexagramDetail)
	mux.HandleFunc("/api/question-types", handleQuestionTypes)
	mux.HandleFunc("/api/liuren/divine", requireCode(handleLiuRen))
	mux.HandleFunc("/api/huican/divine", requireCode(handleHuCan))
	mux.HandleFunc("/api/qimen/divine", requireCode(handleQimen))

	mux.HandleFunc("/api/admin/login", handleAdminLogin)
	mux.HandleFunc("/api/admin/logout", handleAdminLogout)
	mux.HandleFunc("/api/admin/codes", handleAdminCodes)

	// 静态文件：embed 的 web 目录。/ 默认提供 index.html（即扉页"卜筮明心"）。
	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", http.FileServer(http.FS(sub)))

	srv := &http.Server{
		Addr:         addr,
		Handler:      logMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	fmt.Printf("  周易服务已启动：http://%s\n", addr)
	fmt.Println("  Ctrl+C 退出")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
