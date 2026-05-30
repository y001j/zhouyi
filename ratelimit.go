package main

import (
	"net/http"
	"sync"
	"time"
)

// 极简滑动窗口限频：每个 (bucket, ip) 维护近 window 时间内的失败时间戳数组。
// 用于阻断暴力破解（管理员密码）与暴力枚举（访问码），仅对"失败请求"计数。

type rateBucket struct {
	window time.Duration
	max    int
}

var (
	rlBuckets = map[string]rateBucket{
		"login": {window: 10 * time.Minute, max: 10}, // 每 IP 10 分钟内最多 10 次登录失败
		"code":  {window: 1 * time.Minute, max: 30},  // 每 IP 1 分钟内最多 30 次访问码失败
	}
	rlMu      sync.Mutex
	rlEvents  = map[string][]time.Time{} // key: bucket + "|" + ip
	rlSweepAt time.Time
)

// rateLimitAllow 判断当前请求是否允许：
// - 返回 true 表示放行（仍需上层执行业务）
// - 返回 false 表示已超频，应直接 429
func rateLimitAllow(bucket, ip string) bool {
	b, ok := rlBuckets[bucket]
	if !ok {
		return true
	}
	now := time.Now()
	key := bucket + "|" + ip
	rlMu.Lock()
	defer rlMu.Unlock()
	maybeSweepLocked(now)
	events := rlEvents[key]
	cutoff := now.Add(-b.window)
	// 丢弃窗口外旧事件
	kept := events[:0]
	for _, t := range events {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	if len(kept) >= b.max {
		rlEvents[key] = kept
		return false
	}
	rlEvents[key] = kept
	return true
}

// rateLimitRecordFailure 记录一次失败（仅失败才计入限额）。
func rateLimitRecordFailure(bucket, ip string) {
	now := time.Now()
	key := bucket + "|" + ip
	rlMu.Lock()
	defer rlMu.Unlock()
	rlEvents[key] = append(rlEvents[key], now)
}

// maybeSweepLocked 周期性清理过期 key，避免内存无限增长。caller 持锁。
func maybeSweepLocked(now time.Time) {
	if now.Sub(rlSweepAt) < 5*time.Minute {
		return
	}
	rlSweepAt = now
	for key, events := range rlEvents {
		// 找到 key 对应的 bucket
		i := -1
		for j, c := range key {
			if c == '|' {
				i = j
				break
			}
		}
		if i <= 0 {
			delete(rlEvents, key)
			continue
		}
		b, ok := rlBuckets[key[:i]]
		if !ok {
			delete(rlEvents, key)
			continue
		}
		cutoff := now.Add(-b.window)
		kept := events[:0]
		for _, t := range events {
			if t.After(cutoff) {
				kept = append(kept, t)
			}
		}
		if len(kept) == 0 {
			delete(rlEvents, key)
		} else {
			rlEvents[key] = kept
		}
	}
}

// writeRateLimited 统一返回 429。
func writeRateLimited(w http.ResponseWriter, msg string) {
	if msg == "" {
		msg = "请求过于频繁，请稍后再试"
	}
	w.Header().Set("Retry-After", "60")
	writeError(w, http.StatusTooManyRequests, msg)
}
