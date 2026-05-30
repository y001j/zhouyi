package qimen

import (
	"testing"
	"time"
)

// TestBuildPan_Smoke 构建端到端集成测试。
// 确保 3 个典型时刻能跑通并产出合理的盘面。
func TestBuildPan_Smoke(t *testing.T) {
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
			pan, err := BuildPan(tc.when)
			if err != nil {
				t.Fatalf("BuildPan err: %v", err)
			}

			// 结构完整性：8 个非中5宫都应该有星/门/神
			for i := 0; i < 9; i++ {
				c := pan.Cells[i]
				if c.PalaceName == "" {
					t.Errorf("cell[%d] 无宫位名", i)
				}
				if c.EarthStem == "" {
					t.Errorf("cell[%d] 无地盘干", i)
				}
				if i == 4 {
					// 中5宫：天禽恒在、无门无神
					if c.Star != "天禽" {
						t.Errorf("中5宫 star=%q, want 天禽", c.Star)
					}
					if c.Door != "" || c.God != "" {
						t.Errorf("中5宫应无门无神: door=%q god=%q", c.Door, c.God)
					}
				} else {
					if c.Star == "" {
						t.Errorf("cell[%d] 无九星", i)
					}
					if c.Door == "" {
						t.Errorf("cell[%d] 无八门", i)
					}
					if c.God == "" {
						t.Errorf("cell[%d] 无八神", i)
					}
				}
			}

			if pan.ZhiFuStar == "" || pan.ZhiFuPalace == "" {
				t.Errorf("值符缺失：star=%q palace=%q", pan.ZhiFuStar, pan.ZhiFuPalace)
			}
			if pan.ZhiShiGate == "" || pan.ZhiShiPalace == "" {
				t.Errorf("值使缺失：gate=%q palace=%q", pan.ZhiShiGate, pan.ZhiShiPalace)
			}

			t.Log("\n" + pan.Render())
		})
	}
}
