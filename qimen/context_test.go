package qimen

import (
	"testing"
	"time"
)

// TestBuildContextSmoke 最简冒烟测试：确保 BuildContext 能跑通并返回合理字段。
func TestBuildContextSmoke(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tests := []struct {
		name string
		when time.Time
	}{
		{"2024-01-01 00:00", time.Date(2024, 1, 1, 0, 0, 0, 0, loc)},
		{"2024-06-21 12:00", time.Date(2024, 6, 21, 12, 0, 0, 0, loc)},
		{"2023-11-15 15:30", time.Date(2023, 11, 15, 15, 30, 0, 0, loc)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, err := BuildContext(tc.when)
			if err != nil {
				t.Fatalf("BuildContext err: %v", err)
			}
			// 基本字段不为空
			if ctx.DayGZ == "" || ctx.HourGZ == "" {
				t.Errorf("empty four-pillar: day=%q hour=%q", ctx.DayGZ, ctx.HourGZ)
			}
			if ctx.Dun != "阳遁" && ctx.Dun != "阴遁" {
				t.Errorf("unexpected Dun: %q", ctx.Dun)
			}
			if ctx.Yuan != "上元" && ctx.Yuan != "中元" && ctx.Yuan != "下元" {
				t.Errorf("unexpected Yuan: %q", ctx.Yuan)
			}
			if ctx.Ju < 1 || ctx.Ju > 9 {
				t.Errorf("Ju out of range: %d", ctx.Ju)
			}
			if _, ok := XunshouToDungan[ctx.Xunshou]; !ok {
				t.Errorf("bad Xunshou: %q", ctx.Xunshou)
			}
			t.Logf("%s → %s", tc.name, ctx.Summary())
		})
	}
}
