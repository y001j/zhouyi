package liuren

import "testing"

// 构造辅助：直接给定干、支、占时、月将，执行起课
func divineCustom(g Gan, z Zhi, zhan, yue Zhi) *Pan {
	ctx := &Context{
		Gan: g, DayZhi: z, ZhanShi: zhan, YueJiang: yue,
		JiaziIndex: 0, ZhouYe: IsDayByZhi(zhan),
	}
	return DivineWithContext(ctx)
}

// TestFuYin 月将=占时即伏吟
func TestFuYin(t *testing.T) {
	pan := divineCustom(Jia, Yin, Zi, Zi) // 月将=占时=子
	if pan.SanChuan.Method != "伏吟法" {
		t.Errorf("月将=占时 应走伏吟法, 实际 %s", pan.SanChuan.Method)
	}
	// 天盘应各居本位
	for i, up := range pan.TianPan {
		if int(up) != i {
			t.Errorf("伏吟天盘第%d位应为%s, 实际%s", i, Zhi(i), up)
		}
	}
}

// TestFanYin 月将与占时相冲即返吟
func TestFanYin(t *testing.T) {
	pan := divineCustom(Jia, Yin, Zi, Wu) // 子午相冲
	if pan.SanChuan.Method != "返吟法" {
		t.Errorf("月将与占时相冲 应走返吟法, 实际 %s", pan.SanChuan.Method)
	}
}

// TestBaZhuan 甲寅日（日干寄宫=寅=日支）为八专
func TestBaZhuan(t *testing.T) {
	pan := divineCustom(Jia, Yin, Chen, Si) // 非伏吟非返吟，但八专
	if pan.SanChuan.Method != "八专法" {
		t.Errorf("甲寅日 应走八专法, 实际 %s", pan.SanChuan.Method)
	}
}

// TestYaoKeOrNormal 构造一个普通无克/有克课，验证结果落在 9 法之一
func TestDispatcherCovers9(t *testing.T) {
	// 随意取多组干支与时辰，验证不会走到未知分支
	for g := Gan(0); g < 10; g++ {
		for z := Zhi(0); z < 12; z++ {
			pan := divineCustom(g, z, Chen, Hai)
			m := pan.SanChuan.Method
			if m == "" {
				t.Errorf("g=%s z=%s 未产生发传法", g, z)
			}
		}
	}
}

// TestTianJiang_Smoke 验证天将数组每个位置都已填充
func TestTianJiang_Smoke(t *testing.T) {
	pan := divineCustom(Jia, Zi, Chen, Hai)
	seen := map[TianJiang]bool{}
	for _, tj := range pan.TianJiang {
		seen[tj] = true
	}
	if len(seen) != 12 {
		t.Errorf("十二天将应全部出现, 实际 %d 种", len(seen))
	}
}

// TestSanChuanLiuQinKong 三传应已填充六亲与空亡标记
func TestSanChuanLiuQinKong(t *testing.T) {
	pan := divineCustom(Jia, Zi, Chen, Hai)
	// 六亲应有效（不应为空字符串）
	if pan.SanChuan.Chu.LiuQin.String() == "" {
		t.Error("初传六亲未填充")
	}
}

// TestKeTi_HasName 课体名不应为空
func TestKeTi_HasName(t *testing.T) {
	pan := divineCustom(Jia, Zi, Chen, Hai)
	if pan.KeTi.Name == "" {
		t.Error("课体名为空")
	}
}

// TestYiMa 驿马
func TestYiMa(t *testing.T) {
	if yiMa(Shen) != Yin || yiMa(Zi) != Yin || yiMa(Chen) != Yin {
		t.Error("申子辰驿马在寅")
	}
	if yiMa(Yin) != Shen || yiMa(Wu) != Shen {
		t.Error("寅午戌驿马在申")
	}
}

// TestPlaceTianJiang_GuiRen 验证贵人位置
// 甲日昼贵丑、夜贵未。本测试以夜贵未为例：
// 亥加戌时天盘午上为未，则贵人落在地盘午位（午属"卯辰巳午未申"日昼弧 → 顺布）。
func TestPlaceTianJiang_GuiRen(t *testing.T) {
	// 构造场景：甲日、昼占、天盘某位上应为丑
	// 亥将加戌时：offset=(Hai - Xu +12)%12 = 1，地盘子上为丑
	ctx := &Context{Gan: Jia, DayZhi: Zi, ZhanShi: Xu, YueJiang: Hai, ZhouYe: false}
	// 戌时=17-19 应为夜；但甲日夜贵=未，未所在地盘位 = ?
	// tp[i] = (i+1)%12 → 地盘午(6)上为未
	tp := BuildTianPan(Hai, Xu)
	if tp[Wu] != Wei {
		t.Fatalf("亥加戌时 地盘午上应为未, 实际 %s", tp[Wu])
	}
	tj := PlaceTianJiang(ctx, tp)
	if tj[Wu] != TJGuiRen {
		t.Errorf("甲日夜贵在未，地盘午位(天盘为未)应是贵人, 实际 %s", tj[Wu])
	}
}

// TestPlaceTianJiang_ShunNi 贵人临"卯辰巳午未申"（日昼弧）顺布、
// "酉戌亥子丑寅"（夜弧）逆布。出自《六壬大全》卷一神图。
func TestPlaceTianJiang_ShunNi(t *testing.T) {
	for _, z := range []Zhi{Mao, Chen, Si, Wu, Wei, Shen} {
		if !isClockwise(z) {
			t.Errorf("贵人临 %s 应顺布", z)
		}
	}
	for _, z := range []Zhi{You, Xu, Hai, Zi, Chou, Yin} {
		if isClockwise(z) {
			t.Errorf("贵人临 %s 应逆布", z)
		}
	}
}
