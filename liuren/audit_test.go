package liuren

import "testing"

// TestDunGan_JiaZiXun 甲子旬：子→甲 … 酉→癸；戌亥旬空
func TestDunGan_JiaZiXun(t *testing.T) {
	// 甲子 jiaziIndex=0
	expect := []struct {
		z  Zhi
		g  Gan
		ok bool
	}{
		{Zi, Jia, true}, {Chou, Yi, true}, {Yin, Bing, true},
		{Mao, Ding, true}, {Chen, Wu1, true}, {Si, Ji, true},
		{Wu, Geng, true}, {Wei, Xin, true}, {Shen, Ren, true}, {You, Gui, true},
		{Xu, -1, false}, {Hai, -1, false}, // 旬空
	}
	for _, e := range expect {
		got := DunGan(e.z, 0)
		if e.ok && got != e.g {
			t.Errorf("甲子旬 %s 应遁 %s, 实际 %s", e.z, e.g, got)
		}
		if !e.ok && got != -1 {
			t.Errorf("甲子旬 %s 应空亡无遁干", e.z)
		}
	}
}

// TestDunGan_JiaXuXun 甲戌旬（jiaziIndex=10）：戌→甲、亥→乙 …；申酉旬空
func TestDunGan_JiaXuXun(t *testing.T) {
	if DunGan(Xu, 10) != Jia {
		t.Errorf("甲戌旬 戌 应遁甲, 实际 %s", DunGan(Xu, 10))
	}
	if DunGan(Hai, 10) != Yi {
		t.Errorf("甲戌旬 亥 应遁乙, 实际 %s", DunGan(Hai, 10))
	}
	if DunGan(Shen, 10) != -1 {
		t.Error("甲戌旬 申 应旬空")
	}
}

// TestMaoXing_ZhongMo_ShunXu 昴星阳日：中传=支上神、末传=干上神（《六壬大全》"刚日先辰后日"）
func TestMaoXing_ZhongMo_ShunXu(t *testing.T) {
	// 构造无克无遥场景。甲寅日不可（是八专）。
	// 甲申日、月将子、占时辰 → 天盘子压辰，offset=(Zi-Chen+12)%12=8
	ctx := &Context{
		Gan: Jia, DayZhi: Shen, ZhanShi: Chen, YueJiang: Zi,
		JiaziIndex: 20, ZhouYe: IsDayByZhi(Chen),
	}
	pan := DivineWithContext(ctx)
	if pan.SanChuan.Method == "昴星法" {
		ganUp := pan.TianPan[GanJiGong[ctx.Gan]]
		zhiUp := pan.TianPan[ctx.DayZhi]
		if pan.SanChuan.Zhong.Zhi != zhiUp {
			t.Errorf("昴星阳日 中传应=支上神 %s, 实际 %s", zhiUp, pan.SanChuan.Zhong.Zhi)
		}
		if pan.SanChuan.Mo.Zhi != ganUp {
			t.Errorf("昴星阳日 末传应=干上神 %s, 实际 %s", ganUp, pan.SanChuan.Mo.Zhi)
		}
	}
}

// TestMaoXing_Direct_Yang 阳日（刚日）昴星：中=支上、末=干上
//
// 据《六壬大全》卷一入式法："刚日先辰而后日"——辰=日支、日=日干。
func TestMaoXing_Direct_Yang(t *testing.T) {
	ctx := &Context{
		Gan: Jia, DayZhi: Shen, ZhanShi: Chen, YueJiang: Zi,
		JiaziIndex: 20,
	}
	tp := BuildTianPan(ctx.YueJiang, ctx.ZhanShi)
	ganUp := tp[GanJiGong[ctx.Gan]]
	zhiUp := tp[ctx.DayZhi]
	sc := maoXingMethod(ctx, tp, [4]Ke{})
	if sc.Zhong.Zhi != zhiUp {
		t.Errorf("阳日昴星 中传应=支上神 %s, 实际 %s", zhiUp, sc.Zhong.Zhi)
	}
	if sc.Mo.Zhi != ganUp {
		t.Errorf("阳日昴星 末传应=干上神 %s, 实际 %s", ganUp, sc.Mo.Zhi)
	}
}

// TestMaoXing_Direct_Yin 阴日（柔日）昴星：中=干上、末=支上
//
// 据《六壬大全》卷一入式法："柔日先日而后辰"。
func TestMaoXing_Direct_Yin(t *testing.T) {
	ctx := &Context{
		Gan: Yi, DayZhi: You, ZhanShi: Chen, YueJiang: Zi,
		JiaziIndex: 21,
	}
	tp := BuildTianPan(ctx.YueJiang, ctx.ZhanShi)
	ganUp := tp[GanJiGong[ctx.Gan]]
	zhiUp := tp[ctx.DayZhi]
	sc := maoXingMethod(ctx, tp, [4]Ke{})
	if sc.Zhong.Zhi != ganUp {
		t.Errorf("阴日昴星 中传应=干上神 %s, 实际 %s", ganUp, sc.Zhong.Zhi)
	}
	if sc.Mo.Zhi != zhiUp {
		t.Errorf("阴日昴星 末传应=支上神 %s, 实际 %s", zhiUp, sc.Mo.Zhi)
	}
}

// TestYaoKe_Naming_HaoshiTanshe 验证蒿矢/弹射命名遵循经典：
//
//	蒿矢=神（上神）遥克日干（外来克我）
//	弹射=日干遥克神（我克外）
//
// 见《六壬大全》卷一入式法、《大六壬指南》卷一心印赋。
// 历史 BUG：原代码命名颠倒，本测试锁定经典派系。
func TestYaoKe_Naming_HaoshiTanshe(t *testing.T) {
	// 直接构造遥克场景。手工合成 ke：上神克日（蒿矢）
	// 甲日（木），上神 庚（金）→ 金克木 → 蒿矢
	tp := BuildTianPan(Zi, Zi) // 任意天盘，对 yaoKeMethod 的语义不影响
	ctx := &Context{Gan: Jia, DayZhi: Yin, ZhanShi: Zi, YueJiang: Zi}
	// 第一课：上=申（金）下=寅（木实际为甲寄宫）
	ke := [4]Ke{
		{Index: 1, Upper: Shen, Lower: Yin}, // 申金克甲木 → 蒿矢
		{Index: 2, Upper: Yin, Lower: Yin},  // 比和
		{Index: 3, Upper: Yin, Lower: Yin},  // 比和
		{Index: 4, Upper: Yin, Lower: Yin},  // 比和
	}
	sc, ok := yaoKeMethod(ctx, tp, ke)
	if !ok {
		t.Fatal("应当识别为遥克")
	}
	if sc.Note != "蒿矢课" {
		t.Errorf("上神(申)克日干(甲) 应判为蒿矢课, 实际 %s", sc.Note)
	}
	// 反例：日干克上神 → 弹射
	ke2 := [4]Ke{
		{Index: 1, Upper: Wei, Lower: Yin}, // 甲木克未土 → 弹射
		{Index: 2, Upper: Yin, Lower: Yin},
		{Index: 3, Upper: Yin, Lower: Yin},
		{Index: 4, Upper: Yin, Lower: Yin},
	}
	sc2, ok2 := yaoKeMethod(ctx, tp, ke2)
	if !ok2 {
		t.Fatal("应当识别为遥克")
	}
	if sc2.Note != "弹射课" {
		t.Errorf("日干(甲)克上神(未) 应判为弹射课, 实际 %s", sc2.Note)
	}
}

// TestKeTi_HuShiZhuanPeng 验证 keti.go 课名查表：虎视转蓬课能命中
//
// 历史 BUG：keti.go 写为"虎视转篷"（草字头篷）但 santraan_yao.go 输出"虎视转蓬"，
// 导致查表永远 miss、断辞输出"——"。本测试锁定字形一致。
func TestKeTi_HuShiZhuanPeng(t *testing.T) {
	if _, ok := ketiSummaryTable["虎视转蓬课"]; !ok {
		t.Error("虎视转蓬课 应在 ketiSummaryTable 中（曾因错别字\"篷\"miss）")
	}
}

// TestBaZhuan_Yin_FromZhiYin 八专阴日：初传从支上阴神（四课上神）逆数三位
//
// 据《六壬粹言》卷一 p.29："柔日从支阴遁逆数三神"。
// 支阴 = 支上阴神 = 四课上神 = tianpan[tianpan[DayZhi]]。
func TestBaZhuan_Yin_FromZhiYin(t *testing.T) {
	// 丁未日：丁寄未，与日支未同位，是八专。占阴日。
	// 设定一个无克场景：选 月将 寅、占时 子。
	ctx := &Context{
		Gan: Ding, DayZhi: Wei, ZhanShi: Zi, YueJiang: Yin,
		JiaziIndex: 43, // 丁未在甲辰旬第 4 位
	}
	tp := BuildTianPan(ctx.YueJiang, ctx.ZhanShi)
	sc := baZhuanMethod(ctx, tp, [4]Ke{})
	if sc.Method != "八专法" {
		t.Fatalf("丁未日应走八专法, 实际 %s", sc.Method)
	}
	// 期望初传 = 支阴神 - 2
	zhiYin := tp[tp[ctx.DayZhi]]
	expected := Zhi((int(zhiYin) - 2 + 12) % 12)
	if sc.Chu.Zhi != expected {
		t.Errorf("丁未日八专 初传应=%s（支阴%s逆数3位）, 实际 %s",
			expected, zhiYin, sc.Chu.Zhi)
	}
	// 中末传应=干上神
	ganUp := tp[GanJiGong[ctx.Gan]]
	if sc.Zhong.Zhi != ganUp || sc.Mo.Zhi != ganUp {
		t.Errorf("丁未日八专 中末应=干上神 %s, 实际 中%s 末%s",
			ganUp, sc.Zhong.Zhi, sc.Mo.Zhi)
	}
}

// TestBaZhuan_Yang_FromGanShang 八专阳日：初传从干上神顺数三位
func TestBaZhuan_Yang_FromGanShang(t *testing.T) {
	// 甲寅日：甲寄寅，与日支寅同位，是八专。占阳日。
	ctx := &Context{
		Gan: Jia, DayZhi: Yin, ZhanShi: Zi, YueJiang: Yin,
		JiaziIndex: 50, // 甲寅在甲寅旬第 0 位
	}
	tp := BuildTianPan(ctx.YueJiang, ctx.ZhanShi)
	sc := baZhuanMethod(ctx, tp, [4]Ke{})
	if sc.Method != "八专法" {
		t.Fatalf("甲寅日应走八专法, 实际 %s", sc.Method)
	}
	ganUp := tp[GanJiGong[ctx.Gan]]
	expected := Zhi((int(ganUp) + 2) % 12)
	if sc.Chu.Zhi != expected {
		t.Errorf("甲寅日八专 初传应=%s（干上%s顺数3位）, 实际 %s",
			expected, ganUp, sc.Chu.Zhi)
	}
	if sc.Zhong.Zhi != ganUp || sc.Mo.Zhi != ganUp {
		t.Errorf("甲寅日八专 中末应=干上神 %s, 实际 中%s 末%s",
			ganUp, sc.Zhong.Zhi, sc.Mo.Zhi)
	}
}

// TestBieZe_Yang_GanLiuhe 别责阳日初传：取日干六合的阴干寄宫上之天盘神
//
// 据《六壬大全》"刚日干合上头神"、《指南》"甲己庚乙丙辛丁壬戊癸六合也"。
// 例：甲日 → 与之合者己 → 己寄未 → 初传 = 天盘 未 位上之神。
func TestBieZe_Yang_GanLiuhe(t *testing.T) {
	cases := []struct {
		gan      Gan
		expectAt Zhi // 预期初传应取自 tp[expectAt]
	}{
		{Jia, Wei},   // 甲↔己，己寄未
		{Bing, Xu},   // 丙↔辛，辛寄戌
		{Wu1, Chou},  // 戊↔癸，癸寄丑
		{Geng, Chen}, // 庚↔乙，乙寄辰
		{Ren, Wei},   // 壬↔丁，丁寄未
	}
	for _, c := range cases {
		ctx := &Context{
			Gan: c.gan, DayZhi: Shen, ZhanShi: Chen, YueJiang: Zi,
			JiaziIndex: 0,
		}
		tp := BuildTianPan(ctx.YueJiang, ctx.ZhanShi)
		sc := bieZeMethod(ctx, tp, [4]Ke{})
		expected := tp[c.expectAt]
		if sc.Chu.Zhi != expected {
			t.Errorf("%s日别责 初传应=tp[%s]=%s, 实际 %s",
				c.gan, c.expectAt, expected, sc.Chu.Zhi)
		}
		// 中末传应=干上神
		ganUp := tp[GanJiGong[c.gan]]
		if sc.Zhong.Zhi != ganUp || sc.Mo.Zhi != ganUp {
			t.Errorf("%s日别责 中末应=干上神 %s, 实际 中%s 末%s",
				c.gan, ganUp, sc.Zhong.Zhi, sc.Mo.Zhi)
		}
	}
}

// TestBieZe_Yin_UpperNotLower 别责阴日初传应为天盘神而非地支本身
func TestBieZe_Yin_UpperNotLower(t *testing.T) {
	// 乙日，日支取一个能让 sanHeNext(支) 在天盘上映射到另一支的情况
	// 乙亥日，DaySupport=亥，sanHeNext(亥)=卯（亥卯未三合）→ 初传=天盘 卯 位上的神
	ctx := &Context{
		Gan: Yi, DayZhi: Hai, ZhanShi: Zi, YueJiang: Chen,
		JiaziIndex: 11,
	}
	tp := BuildTianPan(ctx.YueJiang, ctx.ZhanShi)
	sc := bieZeMethod(ctx, tp, [4]Ke{})
	// 初传应为 tp[卯]，不是 Mao 本身
	expected := tp[sanHeNext(ctx.DayZhi)]
	if sc.Chu.Zhi != expected {
		t.Errorf("别责阴日 初传应为天盘神 %s（tp[%s]），实际 %s",
			expected, sanHeNext(ctx.DayZhi), sc.Chu.Zhi)
	}
}

// TestSheHai_Direction 涉害"行来本家"顺行：smoke + 方向断言
//
// 《大全》卷一 p473「涉害行来本家止」是顺行（与天盘旋转反向）。
// 构造月将子加时午（offset=6 = 全天盘地盘相冲），任取一个上神，
// 顺行步数=逆行步数=6，无方向偏差；
// 再换一个偏移量 offset=1（亥加戌），手算顺行经过的位置作为基线断言。
func TestSheHai_Direction(t *testing.T) {
	// 1. smoke
	tp := BuildTianPan(Chen, Zi)
	k := Ke{Upper: Yin, Lower: Shen, Relation: "下贼上"}
	n, _ := sheHaiCount(tp, k)
	if n < 0 || n > 11 {
		t.Errorf("涉害计数应在 [0,11] 范围, 实际 %d", n)
	}

	// 2. 方向断言：月将亥加占时戌（offset=1），上神戌临地盘酉。
	//    本家=戌，从酉顺行 1 步即到戌——途中不经任何位（步数=1，"本位不计"）。
	//    此时 sheHaiCount 应返回 0；若错用逆行则会经过 申-未-午-巳-辰-卯-寅-丑-子-亥 共 10 位有受克可能。
	tp2 := BuildTianPan(Hai, Xu)
	k2 := Ke{Upper: Xu, Lower: You}
	n2, _ := sheHaiCount(tp2, k2)
	if n2 != 0 {
		t.Errorf("涉害顺行场景（offset=1, 戌临酉）应返回 0 步内无涉害, 实际 %d；可能仍是逆行", n2)
	}
}

// TestDingShen 金日/水日逢丁的识别
func TestDingShen(t *testing.T) {
	// 甲子旬 卯=丁
	if !IsDingShen(Mao, 0) {
		t.Error("甲子旬 卯 应为丁神")
	}
	if IsDingShen(Chen, 0) {
		t.Error("甲子旬 辰 不是丁神")
	}
}

// TestBaZhuan_XinYou 辛酉日：寄宫戌、日支酉，干支同属金行紧邻，仍归八专。
//
// 这是六个八专日中唯一寄宫≠日支的特例（《六壬大全》卷一）。
// 历史 BUG：原触发条件 GanJiGong[Gan]==DayZhi 漏掉辛酉日。
func TestBaZhuan_XinYou(t *testing.T) {
	if !isBaZhuanDay(Xin, You) {
		t.Fatal("辛酉日应识别为八专日")
	}
	// 选无克场景：辛酉日 月将寅、占时子（任意非伏吟非返吟即可）
	ctx := &Context{
		Gan: Xin, DayZhi: You, ZhanShi: Zi, YueJiang: Yin,
		JiaziIndex: 57, // 辛酉在甲寅旬第 7 位
	}
	pan := DivineWithContext(ctx)
	if pan.SanChuan.Method != "八专法" {
		t.Errorf("辛酉日应走八专法, 实际 %s", pan.SanChuan.Method)
	}
}

// TestBaZhuan_AllSixDays 八专日完整白名单：6 个、且仅这 6 个
func TestBaZhuan_AllSixDays(t *testing.T) {
	want := map[Gan]Zhi{
		Jia: Yin, Ding: Wei, Ji: Wei,
		Geng: Shen, Xin: You, Gui: Chou,
	}
	for g, z := range want {
		if !isBaZhuanDay(g, z) {
			t.Errorf("%s%s 日应为八专", g, z)
		}
	}
	// 反例：寄宫=日支但不在白名单的（实际不存在，遍历所有干支验证白名单不漏不多）
	hits := 0
	for g := Jia; g <= Gui; g++ {
		for z := Zi; z <= Hai; z++ {
			if isBaZhuanDay(g, z) {
				hits++
			}
		}
	}
	if hits != 6 {
		t.Errorf("八专日应恰好 6 个, 实际 %d", hits)
	}
}

// TestTianJiang_ShunBu_GuiRenAtMao 贵人临卯（日昼弧）应顺布
//
// 历史 BUG：原 isClockwise 划分将卯归"亥子丑寅卯辰顺布"，
// 与主流《六壬大全》"卯辰巳午未申顺、酉戌亥子丑寅逆"相反。
// 此测试锁定主流派系，防止回退。
func TestTianJiang_ShunBu_GuiRenAtMao(t *testing.T) {
	if !isClockwise(Mao) {
		t.Error("贵人临卯（日昼弧）应顺布")
	}
	// 顺布：贵→腾(下一位顺时针)
	// 构造甲日昼贵丑加于地盘卯：取 月将卯、占时丑 → offset=2，地盘丑(1)上为卯(3)？
	// 直接验证函数即可
}

// TestTianJiang_NiBu_GuiRenAtYou 贵人临酉（夜弧）应逆布
func TestTianJiang_NiBu_GuiRenAtYou(t *testing.T) {
	if isClockwise(You) {
		t.Error("贵人临酉（夜弧）应逆布")
	}
}

// TestTianJiang_BoundaryArc 顺逆弧分界线 全 12 位枚举
func TestTianJiang_BoundaryArc(t *testing.T) {
	shun := map[Zhi]bool{Mao: true, Chen: true, Si: true, Wu: true, Wei: true, Shen: true}
	for z := Zi; z <= Hai; z++ {
		want := shun[z]
		if isClockwise(z) != want {
			t.Errorf("贵人临 %s 顺逆判定错误：want=%v got=%v", z, want, isClockwise(z))
		}
	}
}

// TestGuiRen_Consistency 贵人口诀 5 组验证
func TestGuiRen_Consistency(t *testing.T) {
	cases := []struct {
		g        Gan
		day, nig Zhi
	}{
		{Jia, Chou, Wei}, {Wu1, Chou, Wei}, {Geng, Chou, Wei},
		{Yi, Zi, Shen}, {Ji, Zi, Shen},
		{Bing, Hai, You}, {Ding, Hai, You},
		{Ren, Si, Mao}, {Gui, Si, Mao},
		{Xin, Wu, Yin},
	}
	for _, c := range cases {
		if GuiRenByGan[c.g][0] != c.day || GuiRenByGan[c.g][1] != c.nig {
			t.Errorf("%s 贵人应昼%s夜%s, 实际 昼%s夜%s",
				c.g, c.day, c.nig, GuiRenByGan[c.g][0], GuiRenByGan[c.g][1])
		}
	}
}
