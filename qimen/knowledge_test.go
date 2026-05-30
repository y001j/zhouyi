package qimen

import (
	"testing"
	"time"
)

// TestDoorPalaceTable_Completeness 验证 72 条表完整（8 门 × 9 宫）且无遗漏
func TestDoorPalaceTable_Completeness(t *testing.T) {
	if len(DoorPalaceTable) != 72 {
		t.Fatalf("DoorPalaceTable 条数 = %d，应为 72", len(DoorPalaceTable))
	}
	seen := map[[2]int]bool{}
	for _, e := range DoorPalaceTable {
		key := [2]int{e.DoorIdx, e.PalaceIdx}
		if seen[key] {
			t.Errorf("重复条目：门%d 宫%d", e.DoorIdx, e.PalaceIdx)
		}
		seen[key] = true

		if e.DoorIdx < 0 || e.DoorIdx > 7 {
			t.Errorf("DoorIdx 越界: %d", e.DoorIdx)
		}
		if e.PalaceIdx < 1 || e.PalaceIdx > 9 {
			t.Errorf("PalaceIdx 越界: %d", e.PalaceIdx)
		}
		if e.Summary == "" {
			t.Errorf("空 Summary: %+v", e)
		}
	}
	// 覆盖性：8 × 9 = 72，seen 应正好 72 条
	if len(seen) != 72 {
		t.Errorf("唯一条目数 %d ≠ 72（有遗漏）", len(seen))
	}
}

// TestDetectDoorPalaceHits_Smoke 构建一张真实盘面，看自动命中多少条
func TestDetectDoorPalaceHits_Smoke(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	pan, err := BuildPan(time.Date(2024, 6, 21, 12, 0, 0, 0, loc))
	if err != nil {
		t.Fatalf("BuildPan: %v", err)
	}
	hits := DetectDoorPalaceHits(pan)
	// 一张盘面有 8 个非中5宫，每格恰好一扇门 → 应该能稳定命中 8 条
	if len(hits) != 8 {
		t.Errorf("命中 %d 条，预期 8 条", len(hits))
	}
	for _, h := range hits {
		t.Logf("  [%s] %s", AuspiceLabelShort(h.Auspice), h.Summary)
	}
}
