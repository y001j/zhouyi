package liuren

import (
	"testing"
	"time"
)

// TestTianPan_HaiJiaChen 亥将加辰时：天盘亥应压在地盘辰上
func TestTianPan_HaiJiaChen(t *testing.T) {
	tp := BuildTianPan(Hai, Chen)
	if tp[Chen] != Hai {
		t.Fatalf("亥加辰: 地盘辰上应为亥, 实际 %s", tp[Chen])
	}
	// 既然亥压辰，则子压巳、丑压午……以此类推
	expected := []struct {
		di, tian Zhi
	}{
		{Chen, Hai}, {Si, Zi}, {Wu, Chou}, {Wei, Yin},
		{Shen, Mao}, {You, Chen}, {Xu, Si}, {Hai, Wu},
		{Zi, Wei}, {Chou, Shen}, {Yin, You}, {Mao, Xu},
	}
	for _, e := range expected {
		if tp[e.di] != e.tian {
			t.Errorf("地盘 %s 上应为天盘 %s, 实际 %s", e.di, e.tian, tp[e.di])
		}
	}
}

// TestGanJiGong 验证寄宫表
func TestGanJiGong(t *testing.T) {
	cases := []struct {
		g  Gan
		ji Zhi
	}{
		{Jia, Yin}, {Yi, Chen}, {Bing, Si}, {Ding, Wei},
		{Wu1, Si}, {Ji, Wei}, {Geng, Shen}, {Xin, Xu},
		{Ren, Hai}, {Gui, Chou},
	}
	for _, c := range cases {
		if GanJiGong[c.g] != c.ji {
			t.Errorf("%s 寄 %s, 实际 %s", c.g, c.ji, GanJiGong[c.g])
		}
	}
}

// TestSiKe_XuShen 书例：戊子日，辰将加申时
// 月将=辰，占时=申，日=戊子。
// 戊寄巳。
//   offset = (4 - 8 + 12) % 12 = 8
//   天盘(i) = (i+8)%12
//   地盘巳(5)上天盘 = (5+8)%12 = 1 = 丑 → 一课 丑戊
//   一课上神=丑，丑(1)上天盘 = (1+8)%12 = 9 = 酉 → 二课 酉丑
//   三课：日支=子(0)上天盘 = 8 = 申 → 三课 申子
//   三课上神=申(8)上天盘 = (8+8)%12 = 4 = 辰 → 四课 辰申
func TestSiKe_WuZi_ChenJiaShen(t *testing.T) {
	// 不用真实时间，直接构造上下文
	ctx := &Context{
		Gan: Wu1, DayZhi: Zi, ZhanShi: Shen, YueJiang: Chen,
	}
	tp := BuildTianPan(ctx.YueJiang, ctx.ZhanShi)
	ke := BuildSiKe(ctx, tp)

	cases := []struct{ up, lo Zhi }{
		{Chou, Si}, // 一课 丑巳（戊寄巳）
		{You, Chou},
		{Shen, Zi},
		{Chen, Shen},
	}
	for i, c := range cases {
		if ke[i].Upper != c.up || ke[i].Lower != c.lo {
			t.Errorf("第%d课 应为 %s%s, 实际 %s%s", i+1, c.up, c.lo, ke[i].Upper, ke[i].Lower)
		}
	}
}

// TestXunKong_JiaZi 甲子旬空戌亥
func TestXunKong_JiaZi(t *testing.T) {
	// 甲子日：jiaziIndex = 0
	if !IsXunKong(Xu, 0) || !IsXunKong(Hai, 0) {
		t.Error("甲子旬应空戌亥")
	}
	if IsXunKong(Zi, 0) || IsXunKong(You, 0) {
		t.Error("甲子旬不应空子/酉")
	}
}

// TestResolveZhanShi 占时支
func TestResolveZhanShi(t *testing.T) {
	cases := []struct {
		hour int
		z    Zhi
	}{
		{0, Zi}, {1, Chou}, {2, Chou}, {3, Yin}, {11, Wu}, {12, Wu},
		{13, Wei}, {23, Zi},
	}
	for _, c := range cases {
		tm := time.Date(2026, 4, 23, c.hour, 0, 0, 0, time.Local)
		got := ResolveZhanShi(tm)
		if got != c.z {
			t.Errorf("%02d 时应 %s, 实际 %s", c.hour, c.z, got)
		}
	}
}

// TestBuildContext_Smoke 冒烟测试：任意时刻能成功构造上下文
func TestBuildContext_Smoke(t *testing.T) {
	tm := time.Date(2026, 4, 23, 10, 30, 0, 0, time.Local)
	ctx, err := BuildContext(tm)
	if err != nil {
		t.Fatal(err)
	}
	if ctx.Gan < 0 || ctx.DayZhi < 0 {
		t.Error("日干支未解析")
	}
	if _, ok := QiToYueJiang[ctx.QiName]; !ok {
		t.Errorf("月将中气名不在表中: %q", ctx.QiName)
	}
}

// TestRelation 生克关系
func TestRelation(t *testing.T) {
	// 子(水) 与 巳(火)：水克火 → "上克下"（上子下巳）
	if RelationOfZhi(Zi, Si) != "上克下" {
		t.Errorf("子(水)上 巳(火)下 应为上克下, 实际 %s", RelationOfZhi(Zi, Si))
	}
	// 巳(火) 与 子(水)：下克上 → "下贼上"
	if RelationOfZhi(Si, Zi) != "下贼上" {
		t.Error("巳(火)上 子(水)下 应为下贼上")
	}
	// 子(水) 与 亥(水)：比和
	if RelationOfZhi(Zi, Hai) != "比和" {
		t.Error("子亥同水应比和")
	}
}
