package liuren

import "testing"

// 辅助：构造完整盘面（走正式起课流程）
func scenePan(g Gan, z Zhi, zhan, yue Zhi) *Pan {
	ctx := &Context{
		Gan: g, DayZhi: z, ZhanShi: zhan, YueJiang: yue,
		JiaziIndex: 0, ZhouYe: IsDayByZhi(zhan),
	}
	return DivineWithContext(ctx)
}

// TestBiFaCatalog_Has100 目录应有 100 条且序号从 1 到 100
func TestBiFaCatalog_Has100(t *testing.T) {
	cat := BiFaCatalog()
	if len(cat) != 100 {
		t.Fatalf("毕法赋目录应 100 条, 实际 %d", len(cat))
	}
	seen := map[int]bool{}
	for _, e := range cat {
		seen[e.Number] = true
		if e.Title == "" {
			t.Errorf("第 %d 条标题空缺", e.Number)
		}
	}
	for i := 1; i <= 100; i++ {
		if !seen[i] {
			t.Errorf("第 %d 条缺失", i)
		}
	}
}

// TestBiFa_LuLinShen 旺禄临身 —— 乙卯日干上卯
func TestBiFa_LuLinShen(t *testing.T) {
	// 乙寄辰，要天盘辰位上为卯 → offset=(3-4+12)%12=11；即月将卯加占时辰
	pan := scenePan(Yi, Mao, Chen, Mao)
	entries := MatchBiFa(pan)
	has7 := false
	for _, e := range entries {
		if e.Number == 7 {
			has7 = true
		}
	}
	if !has7 {
		t.Error("乙卯日干上卯 应命中毕法第 7 条（旺禄临身）")
	}
}

// TestBiFa_KuiDuTianMen 戌加亥 —— 第 51 条
func TestBiFa_KuiDuTianMen(t *testing.T) {
	// 天盘亥位上为戌 → offset=(Xu - Hai + 12)%12 = 11；即月将戌加占时亥
	pan := scenePan(Ren, Wu, Hai, Xu)
	if pan.TianPan[Hai] != Xu {
		t.Fatalf("setup 错：天盘亥上应为戌，实际 %s", pan.TianPan[Hai])
	}
	entries := MatchBiFa(pan)
	for _, e := range entries {
		if e.Number == 51 {
			return
		}
	}
	t.Error("戌加亥且发用戌 应命中第 51 条（魁度天门）")
}

// TestBiFa_GangSaiGuiHu 辰加寅 —— 第 52 条
func TestBiFa_GangSaiGuiHu(t *testing.T) {
	// 天盘寅上为辰 → offset=(Chen-Yin+12)%12=2；即月将辰加占时寅
	pan := scenePan(Ji, Chou, Yin, Chen)
	if pan.TianPan[Yin] != Chen {
		t.Fatalf("setup 错：天盘寅上应为辰, 实际 %s", pan.TianPan[Yin])
	}
	entries := MatchBiFa(pan)
	for _, e := range entries {
		if e.Number == 52 {
			return
		}
	}
	t.Error("辰加寅 应命中第 52 条（罡塞鬼户）")
}

// TestBiFa_CatalogLookup 目录查询工具
func TestBiFa_CatalogLookup(t *testing.T) {
	e := lookupCatalog(51)
	if e.Title != "魁度天门关隔定" {
		t.Errorf("第 51 条标题应为 '魁度天门关隔定', 实际 %q", e.Title)
	}
	e100 := lookupCatalog(100)
	if e100.Number != 100 {
		t.Error("第 100 条应存在")
	}
}

// TestBiFa_PhaseChain 确认 MatchBiFa 对正常盘不会错误命中过多
func TestBiFa_PhaseChain(t *testing.T) {
	pan := scenePan(Jia, Zi, Chen, Hai)
	entries := MatchBiFa(pan)
	// 至少返回一个（盘面肯定有一些特征）；不要一口气命中 50+
	if len(entries) > 40 {
		t.Errorf("单次匹配 %d 条过多", len(entries))
	}
}
