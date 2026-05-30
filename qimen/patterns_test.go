package qimen

import (
	"strings"
	"testing"
	"time"
)

// 本测试构造"真实时刻"的盘，扫描命中的格局并做合理性检查。
// 由于格局命中依赖具体干支组合，我们主要做：
//   (1) 多个不同时刻下 DetectPatterns 不报错；
//   (2) 至少命中了一些吉/凶格（不要求具体是哪个——不同时刻结果不同）；
//   (3) 每一条 hit 的字段完整。
func TestDetectPatterns_Smoke(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	times := []time.Time{
		time.Date(2024, 1, 1, 12, 0, 0, 0, loc),
		time.Date(2024, 6, 21, 12, 0, 0, 0, loc),
		time.Date(2023, 11, 15, 15, 30, 0, 0, loc),
		time.Date(2026, 4, 24, 2, 30, 0, 0, loc),
		time.Date(2026, 4, 24, 14, 30, 0, 0, loc),
	}
	totalHits := 0
	for _, tt := range times {
		pan, err := BuildPan(tt)
		if err != nil {
			t.Fatalf("BuildPan(%v): %v", tt, err)
		}
		hits := DetectPatterns(pan)
		t.Logf("%s → 命中 %d 个格局", tt.Format("2006-01-02 15:04"), len(hits))
		for _, h := range hits {
			if h.Name == "" {
				t.Errorf("空 name: %+v", h)
			}
			if h.Category == "" {
				t.Errorf("空 category: %+v", h)
			}
			if h.Classic == "" {
				t.Errorf("空 classic: %+v", h)
			}
			t.Logf("  · %s[%s] %s：%s", h.Name, h.Category, aucStr(h.AuspiceScore), h.Summary)
		}
		totalHits += len(hits)
	}
	// 至少要有一些格局命中——不然说明识别规则太严苛
	if totalHits < 3 {
		t.Errorf("5 个时刻共命中 %d 个格局，疑似识别过严", totalHits)
	}
}

func aucStr(s int) string {
	switch {
	case s >= 2:
		return "大吉"
	case s == 1:
		return "吉"
	case s == -1:
		return "凶"
	case s <= -2:
		return "大凶"
	}
	return "平"
}

// TestDetectSanQiDeShi_Rule 验证三奇得使规则判定
func TestDetectSanQiDeShi_Rule(t *testing.T) {
	// 构造一个盘：强制让丙加戊（青龙返首 = 飞鸟跌穴） + 乙加己（三奇得使）
	// 这个通过真实时刻做不到，直接做单元逻辑：看规则函数在预设盘上是否识别
	pan := fakePan(
		[9]string{"戊", "己", "庚", "辛", "壬", "癸", "丁", "丙", "乙"}, // 地盘：阳遁一局
		[9]string{"戊", "乙", "庚", "辛", "壬", "癸", "丁", "丙", "乙"}, // 天盘：刻意做一个"乙加己"（位置1 的天盘=乙、地盘=己）
	)
	hits := detectSanQiDeShi(pan)
	foundYi := false
	for _, h := range hits {
		if h.Name == "三奇得使" && strings.Contains(h.Classic, "乙") {
			foundYi = true
		}
	}
	if !foundYi {
		t.Errorf("未识别到乙奇得使，hits=%v", hits)
	}
}

// TestLeiShenDirective_Smoke 生成类神直指文本，验证包含关键字段
func TestLeiShenDirective_Smoke(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	pan, _ := BuildPan(time.Date(2024, 6, 21, 12, 0, 0, 0, loc))

	for _, qt := range []string{"career", "wealth", "relation", "health", "decision", "timing"} {
		text := LeiShenDirective(pan, qt)
		if text == "" {
			t.Errorf("%s 类神直指为空", qt)
			continue
		}
		if !strings.Contains(text, "- **") {
			t.Errorf("%s 类神格式异常: %s", qt, text[:min(len(text), 100)])
		}
		t.Logf("--- %s ---\n%s", qt, text)
	}
}

// TestExtractQimenSignals_Smoke 验证信号提取跑通
func TestExtractQimenSignals_Smoke(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	pan, _ := BuildPan(time.Date(2024, 6, 21, 12, 0, 0, 0, loc))
	sig := ExtractQimenSignals(pan)
	mustHave := []string{
		"日干",
		"值符",
		"值使",
		"三奇分布",
		"格局走势",
	}
	for _, m := range mustHave {
		if !strings.Contains(sig, m) {
			t.Errorf("signals 缺失 %q", m)
		}
	}
	t.Logf("signals:\n%s", sig)
}

// ============ 辅助 ============

// fakePan 构造一个用于规则验证的假盘（只填地盘/天盘干；其他字段不准但够规则测试用）
func fakePan(earth, heaven [9]string) *Pan {
	var cells [9]Cell
	for i := 0; i < 9; i++ {
		cells[i] = Cell{
			PalaceFei:  i,
			PalaceName: Palaces[i],
			EarthStem:  earth[i],
			HeavenStem: heaven[i],
		}
	}
	return &Pan{
		Ctx: &Context{
			Dun:     "阳遁",
			Ju:      1,
			Dungan:  "戊",
			Xunshou: "甲子",
			DayGan:  "甲",
			HourGan: "甲",
		},
		Cells: cells,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
