package liuren

import "testing"

// TestYing_Classify 五种应机
func TestYing_Classify(t *testing.T) {
	// 甲木日
	// 水生木 → 救应
	if ClassifyYing(Jia, Zi) != YingSheng {
		t.Error("甲日 年命上神为子(水) 应为救应")
	}
	// 金克木 → 损应
	if ClassifyYing(Jia, Shen) != YingKe {
		t.Error("甲日 年命上神为申(金) 应为损应")
	}
	// 木生火 → 脱应
	if ClassifyYing(Jia, Si) != YingTuo {
		t.Error("甲日 年命上神为巳(火) 应为脱应")
	}
	// 木克土 → 制应
	if ClassifyYing(Jia, Chen) != YingZhi {
		t.Error("甲日 年命上神为辰(土) 应为制应")
	}
	// 木木 → 比应
	if ClassifyYing(Jia, Yin) != YingBiHe {
		t.Error("甲日 年命上神为寅(木) 应为比应")
	}
}

// TestNianMing_IncludesYing 起课后 NianMing 应含 Ying 字段
func TestNianMing_IncludesYing(t *testing.T) {
	bm := Yin
	ctx := &Context{
		Gan: Jia, DayZhi: Zi, ZhanShi: Chen, YueJiang: Hai,
		JiaziIndex: 0, BenMing: &bm, BirthYear: 1990, Gender: "男",
	}
	pan := DivineWithContext(ctx)
	if pan.NianMing == nil || pan.NianMing.BenMing == nil {
		t.Fatal("NianMing.BenMing 应存在")
	}
	// 有 Ying 枚举值（至少能 toString）
	if pan.NianMing.BenMing.Ying.String() == "" {
		t.Error("Ying 枚举应有名字")
	}
}

// TestBiFa_10_XiuMu 朽木：初传卯且卯空
func TestBiFa_10_XiuMu(t *testing.T) {
	// 甲午旬 jiaziIndex=30 空辰巳——卯不空；
	// 甲辰旬 jiaziIndex=40 空寅卯——卯空。
	pan := &Pan{
		SanChuan: SanChuan{
			Chu: ChuanEntry{Zhi: Mao},
		},
		Ctx: &Context{Gan: Jia, DayZhi: Zi, JiaziIndex: 40},
	}
	entries := MatchBiFa(pan)
	found := false
	for _, e := range entries {
		if e.Number == 10 {
			found = true
		}
	}
	if !found {
		t.Error("甲辰旬 初传卯 应命中第 10 条（朽木难雕）")
	}
}

// TestBiFa_19_TaiCai 胎财：日干胎神=妻财
func TestBiFa_19_TaiCai(t *testing.T) {
	// 甲木日，胎在酉(金)；甲木克土为财，酉金非土——不成立
	// 庚金日，胎在卯(木)；庚金克木为财，卯木为财 → 应命中
	pan := &Pan{
		SanChuan: SanChuan{Chu: ChuanEntry{Zhi: Zi}},
		Ctx:      &Context{Gan: Geng, DayZhi: Zi, JiaziIndex: 0},
	}
	entries := MatchBiFa(pan)
	found := false
	for _, e := range entries {
		if e.Number == 19 {
			found = true
		}
	}
	if !found {
		t.Error("庚日胎卯为木=庚之财，应命中第 19 条")
	}
}

// TestBiFa_42_SanQi 三奇：三传全为甲戊庚或乙丙丁（按遁干）
func TestBiFa_42_SanQi(t *testing.T) {
	// 甲子旬：子=甲、辰=戊、午=庚 → 三传子辰午即甲戊庚三奇
	pan := &Pan{
		SanChuan: SanChuan{
			Chu:   ChuanEntry{Zhi: Zi},  // 甲
			Zhong: ChuanEntry{Zhi: Chen}, // 戊
			Mo:    ChuanEntry{Zhi: Wu},   // 庚
		},
		Ctx: &Context{Gan: Jia, DayZhi: Zi, JiaziIndex: 0},
	}
	entries := MatchBiFa(pan)
	found := false
	for _, e := range entries {
		if e.Number == 42 {
			found = true
		}
	}
	if !found {
		t.Error("三传甲戊庚 应命中第 42 条（三奇）")
	}
}

// TestBiFa_88_GanZhiChengMu 干支乘墓
func TestBiFa_88_GanZhiChengMu(t *testing.T) {
	// 甲木日，墓未；日支子，支墓(水)为辰
	// 构造天盘：甲寄寅，寅位上神=未；子位上神=辰
	// → offset 必须让 tp[寅]=未、tp[子]=辰，但这两个条件不独立；让 tp[i]=(i+5)%12 → tp[寅]=未(7)=错
	// 直接构造 pan 手填 tp：
	var tp [12]Zhi
	for i := range tp {
		tp[i] = Zhi(i)
	}
	tp[Yin] = Wei // 寅位上神=未
	tp[Zi] = Chen // 子位上神=辰
	pan := &Pan{
		TianPan:   tp,
		TianJiang: [12]TianJiang{},
		Ctx:       &Context{Gan: Jia, DayZhi: Zi, JiaziIndex: 0},
		SanChuan:  SanChuan{Chu: ChuanEntry{Zhi: Wei}, Zhong: ChuanEntry{Zhi: Chen}, Mo: ChuanEntry{Zhi: Chen}},
	}
	entries := MatchBiFa(pan)
	found := false
	for _, e := range entries {
		if e.Number == 88 {
			found = true
		}
	}
	if !found {
		t.Error("甲日干上未、日支子上辰 应命中第 88 条（干支乘墓）")
	}
}

// TestBiFa_MatchCount_Reasonable 新增条目后命中数不应暴涨
func TestBiFa_MatchCount_Reasonable(t *testing.T) {
	pan := scenePan(Jia, Zi, Chen, Hai)
	entries := MatchBiFa(pan)
	if len(entries) > 50 {
		t.Errorf("普通盘面 命中 %d 条过多，规则可能过宽", len(entries))
	}
}
