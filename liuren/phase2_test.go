package liuren

import "testing"

// TestComputeXingNian 男1岁起寅顺；女1岁起申逆
func TestComputeXingNian(t *testing.T) {
	if ComputeXingNian("男", 1) != Yin {
		t.Error("男1岁行年应为寅")
	}
	if ComputeXingNian("男", 13) != Yin {
		t.Error("男13岁行年应回到寅")
	}
	if ComputeXingNian("女", 1) != Shen {
		t.Error("女1岁行年应为申")
	}
	if ComputeXingNian("女", 3) != Wu {
		t.Error("女3岁行年应为午（申-2 逆行）")
	}
}

// TestShenSha_YiMa 申子辰驿马在寅
func TestShenSha_YiMa(t *testing.T) {
	if yiMa(Shen) != Yin || yiMa(Zi) != Yin || yiMa(Chen) != Yin {
		t.Error("申子辰驿马应在寅")
	}
}

// TestShenSha_TaoHua 申子辰桃花在酉；寅午戌桃花在卯
func TestShenSha_TaoHua(t *testing.T) {
	if taoHua(Shen) != You || taoHua(Zi) != You || taoHua(Chen) != You {
		t.Error("申子辰桃花应在酉")
	}
	if taoHua(Yin) != Mao || taoHua(Wu) != Mao || taoHua(Xu) != Mao {
		t.Error("寅午戌桃花应在卯")
	}
}

// TestShenSha_HuaGai 华盖（墓库位）
func TestShenSha_HuaGai(t *testing.T) {
	if huaGai(Shen) != Chen || huaGai(Zi) != Chen || huaGai(Chen) != Chen {
		t.Error("申子辰华盖在辰")
	}
}

// TestShenSha_JiangXing 将星（帝旺位）
func TestShenSha_JiangXing(t *testing.T) {
	if jiangXing(Shen) != Zi || jiangXing(Zi) != Zi {
		t.Error("申子辰将星在子")
	}
	if jiangXing(Yin) != Wu || jiangXing(Wu) != Wu {
		t.Error("寅午戌将星在午")
	}
}

// TestKeTiTags_SanYang 三传皆阳 → 三阳课
func TestKeTiTags_SanYang(t *testing.T) {
	pan := &Pan{
		SanChuan: SanChuan{
			Chu:   ChuanEntry{Zhi: Zi, TianJiang: TJGuiRen},
			Zhong: ChuanEntry{Zhi: Yin, TianJiang: TJQingLong},
			Mo:    ChuanEntry{Zhi: Chen, TianJiang: TJLiuHe},
		},
		TianPan:   BuildTianPan(Zi, Zi),
		TianJiang: [12]TianJiang{},
		Ctx:       &Context{Gan: Jia, DayZhi: Zi},
	}
	tags := KeTiTags(pan)
	found := false
	for _, tg := range tags {
		if tg.Name == "三阳课" {
			found = true
		}
	}
	if !found {
		t.Error("三传皆阳(子寅辰) 应识别三阳课")
	}
}

// TestKeTiTags_SanHeju 三传三合
func TestKeTiTags_SanHeju(t *testing.T) {
	pan := &Pan{
		SanChuan: SanChuan{
			Chu:   ChuanEntry{Zhi: Shen},
			Zhong: ChuanEntry{Zhi: Zi},
			Mo:    ChuanEntry{Zhi: Chen},
		},
		TianPan:   BuildTianPan(Zi, Zi),
		TianJiang: [12]TianJiang{},
		Ctx:       &Context{Gan: Jia, DayZhi: Zi},
	}
	tags := KeTiTags(pan)
	found := false
	for _, tg := range tags {
		if tg.Name == "三合课" {
			found = true
		}
	}
	if !found {
		t.Error("申子辰三传 应识别三合课")
	}
}

// TestBiFa_ZhongMoKong 中末空亡 → 毕法第 82 条「不行传者考初时」
func TestBiFa_ZhongMoKong(t *testing.T) {
	pan := &Pan{
		SanChuan: SanChuan{
			Chu:   ChuanEntry{Zhi: Zi, TianJiang: TJGuiRen},
			Zhong: ChuanEntry{Zhi: Chou, TianJiang: TJTengShe, IsKong: true},
			Mo:    ChuanEntry{Zhi: Yin, TianJiang: TJZhuQue, IsKong: true},
		},
		TianPan:   BuildTianPan(Zi, Zi),
		TianJiang: [12]TianJiang{},
		Ctx:       &Context{Gan: Jia, DayZhi: Zi, JiaziIndex: 0},
	}
	entries := MatchBiFa(pan)
	matched := false
	for _, e := range entries {
		if e.Number == 82 {
			matched = true
		}
	}
	if !matched {
		t.Error("中末空亡 应命中毕法第 82 条")
	}
}

// TestDivineWithContext_NianMingFlow 带本命、出生年起课，盘面应含 NianMing
func TestDivineWithContext_NianMingFlow(t *testing.T) {
	bm := Yin
	ctx := &Context{
		Gan: Jia, DayZhi: Zi, ZhanShi: Chen, YueJiang: Hai,
		JiaziIndex: 0, BenMing: &bm, BirthYear: 1990, Gender: "男",
	}
	pan := DivineWithContext(ctx)
	if pan.NianMing == nil || pan.NianMing.BenMing == nil {
		t.Error("指定本命时 应填充 NianMing.BenMing")
	}
	if pan.NianMing.BenMing.Zhi != Yin {
		t.Error("本命应为寅")
	}
}
