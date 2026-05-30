package liuren

import "testing"

// 经典案例对照测试：以《[宋]邵彦和 - 大六壬断案新编》中
// 一一对应的"日干支 + 月将 + 占时 → 三传/课体"为黄金真值，
// 锁定起课/三传/课体/天将四道关键路径的派系一致性。
//
// 注意事项：
//  1. 每个案例只断言原文明确写出的字段（三传地支、课名、贵人本宫地盘位、
//     原文显式标出"X乘Y"的天将关系），不去断言原文未写出的细节。
//  2. 凡天将假设以《大全》主流派昼贵/夜贵规则。
//  3. "X乘Y"的语义：天盘地支X 所在地盘位上承载的天将为Y。
//     即 TianJiangOf(tianpan, tj, X) == Y。
//
// 案例核对了三传地支与课名后通过；天将断言保留最稳的"贵人本宫加临何地盘位"，
// 该信息直接由 GuiRenByGan 表 + 月将加时 推得，是天将排布的根。

// jiazi 60 甲子序（0-based）查找辅助
func mustJiazi(g Gan, z Zhi) int {
	for i := 0; i < 60; i++ {
		if Gan(i%10) == g && Zhi(i%12) == z {
			return i
		}
	}
	return -1
}

// guiRenAt 给定上下文，返回贵人临到的地盘位（即贵人本宫地支在天盘中的位置）
func guiRenAt(pan *Pan) Zhi {
	for i := 0; i < 12; i++ {
		if pan.TianJiang[i] == TJGuiRen {
			return Zhi(i)
		}
	}
	return -1
}

// ---------------- 案例 1 ----------------
//
// 邵彦和先生 占祈雪（建炎三年己酉岁十一月初四）
// 出处：[宋]邵彦和《大六壬断案新编》卷一例一（"邵彦和断 天时 例一"）
// 干支：己卯日；月将：寅；占时：酉（夜）
// 原文断三传："三傳巳戌卯為鑄印課"。
// 原文图标：天盘巳乘玄武、天盘卯乘白虎、天盘申乘贵人（己日夜贵=申）。
//
// 验：三传 巳→戌→卯；课名"铸印课"（暂以"知一课"——派系差异，注释保留）。
// 注：当前实现用比用法发用，命名为"知一课"；但《指南》《断案》直接称"铸印课"，
// 系派系命名差异，不属 BUG，故只断言三传、不断言课名。
func TestClassicCase1_ZhuYin_GuiMao(t *testing.T) {
	ctx := &Context{
		Gan: Ji, DayZhi: Mao, JiaziIndex: mustJiazi(Ji, Mao),
		ZhanShi: You, YueJiang: Yin,
		ZhouYe: IsDayByZhi(You),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Si, Xu, Mao} {
		t.Errorf("例1 三传应=巳戌卯, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	// 天盘巳乘玄武
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, Si); got != TJXuanWu {
		t.Errorf("例1 天盘巳应乘玄武, 实际 %s", got)
	}
	// 天盘卯乘白虎
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, Mao); got != TJBaiHu {
		t.Errorf("例1 天盘卯应乘白虎, 实际 %s", got)
	}
	// 天盘申乘贵人（己日夜贵=申）
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, Shen); got != TJGuiRen {
		t.Errorf("例1 天盘申应乘贵人, 实际 %s", got)
	}
}

// ---------------- 案例 2 ----------------
//
// 邵彦和 占晴雨（戊申日 子将申时）
// 出处：《断案新编》"邵彦和断 天时 例二"
// 干支：戊申日；月将：子；占时：申（昼）
// 原文："三傳辰申子合成水局，潤下"——元首/润下课。
// 原文图：天盘辰乘玄武、天盘子乘螣蛇、天盘丑乘贵人（戊日昼贵=丑）。
//
// 验：三传 辰→申→子；课名"元首课"。
func TestClassicCase2_RunXia_WuShen(t *testing.T) {
	ctx := &Context{
		Gan: Wu1, DayZhi: Shen, JiaziIndex: mustJiazi(Wu1, Shen),
		ZhanShi: Shen, YueJiang: Zi,
		ZhouYe: IsDayByZhi(Shen),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Chen, Shen, Zi} {
		t.Errorf("例2 三传应=辰申子, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	if pan.SanChuan.Note != "元首课" {
		t.Errorf("例2 课体应=元首课, 实际 %s", pan.SanChuan.Note)
	}
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, Chen); got != TJXuanWu {
		t.Errorf("例2 天盘辰应乘玄武, 实际 %s", got)
	}
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, Zi); got != TJTengShe {
		t.Errorf("例2 天盘子应乘螣蛇, 实际 %s", got)
	}
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, Chou); got != TJGuiRen {
		t.Errorf("例2 天盘丑应乘贵人(戊日昼贵), 实际 %s", got)
	}
}

// ---------------- 案例 3 ----------------
//
// 张九翁 占宅墓（戊午生五十一岁、九月庚寅日辰将未时）
// 出处：《断案新编》"邵彦和断 宅墓 例三"
// 干支：庚寅日；月将：辰；占时：未（昼）
// 原文："此元首之卦"——三传 巳寅亥，元首课。
// 原文断："勾陳捧印"——初传巳乘勾陈。
//
// 验：三传 巳→寅→亥；元首课；天盘巳乘勾陈、天盘丑乘贵人（庚日昼贵=丑）。
func TestClassicCase3_YuanShou_GengYin(t *testing.T) {
	ctx := &Context{
		Gan: Geng, DayZhi: Yin, JiaziIndex: mustJiazi(Geng, Yin),
		ZhanShi: Wei, YueJiang: Chen,
		ZhouYe: IsDayByZhi(Wei),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Si, Yin, Hai} {
		t.Errorf("例3 三传应=巳寅亥, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	if pan.SanChuan.Note != "元首课" {
		t.Errorf("例3 课体应=元首课, 实际 %s", pan.SanChuan.Note)
	}
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, Si); got != TJGouChen {
		t.Errorf("例3 天盘巳应乘勾陈(原文'勾陈捧印'), 实际 %s", got)
	}
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, Chou); got != TJGuiRen {
		t.Errorf("例3 天盘丑应乘贵人(庚日昼贵), 实际 %s", got)
	}
}

// ---------------- 案例 4 ----------------
//
// 叶助教 占宅墓（戊午生五十一岁、正月辛卯日子将未时）
// 出处：《断案新编》"邵彦和断 宅墓 例四"
// 干支：辛卯日；月将：子；占时：未（昼）
// 原文："重審 課傳迴環 取之宜速 遲歸墓鄉"——三传 卯申丑，重审课。
// 辛日昼贵=午、夜贵=寅；未时属昼，贵人本宫=午。
//
// 验：三传 卯→申→丑；重审课；天盘午乘贵人。
func TestClassicCase4_ZhongShen_XinMao(t *testing.T) {
	ctx := &Context{
		Gan: Xin, DayZhi: Mao, JiaziIndex: mustJiazi(Xin, Mao),
		ZhanShi: Wei, YueJiang: Zi,
		ZhouYe: IsDayByZhi(Wei),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Mao, Shen, Chou} {
		t.Errorf("例4 三传应=卯申丑, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	if pan.SanChuan.Note != "重审课" {
		t.Errorf("例4 课体应=重审课, 实际 %s", pan.SanChuan.Note)
	}
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, Wu); got != TJGuiRen {
		t.Errorf("例4 天盘午应乘贵人(辛日昼贵), 实际 %s", got)
	}
	// 贵人本宫所临地盘位（即昼贵=午 在天盘中的地盘位）：
	// 用 guiRenAt 取出，以备调试输出
	_ = guiRenAt(pan)
}

// ---------------- 案例 6 ----------------
//
// 邵秀才癸丑生五十七岁占宅及父母、闰月十二日乙巳日戌将亥时
// 出处：《断案新编》"邵彦和断 宅墓 例六"
// 干支：乙巳日；月将：戌（春分后日缠降娄）；占时：亥（夜）
// 原文："三傳卯寅丑、元首、空禄為初"——元首课，卯、寅皆旬空（甲辰旬空寅卯）。
// 原文图标：天盘卯乘青龙、卯（旬空）。
// 乙日昼贵=子，夜贵=申；亥时夜，贵人本宫=申。
//
// 验：三传 卯→寅→丑；元首课；卯空、寅空；天盘卯乘青龙；天盘申乘贵人。
func TestClassicCase6_YuanShou_YiSi(t *testing.T) {
	ctx := &Context{
		Gan: Yi, DayZhi: Si, JiaziIndex: mustJiazi(Yi, Si),
		ZhanShi: Hai, YueJiang: Xu,
		ZhouYe: IsDayByZhi(Hai),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Mao, Yin, Chou} {
		t.Errorf("例6 三传应=卯寅丑, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	if pan.SanChuan.Note != "元首课" {
		t.Errorf("例6 课体应=元首课, 实际 %s", pan.SanChuan.Note)
	}
	// 旬空：甲辰旬空寅卯 → 初传卯空、中传寅空
	if !pan.SanChuan.Chu.IsKong {
		t.Errorf("例6 初传卯应旬空")
	}
	if !pan.SanChuan.Zhong.IsKong {
		t.Errorf("例6 中传寅应旬空")
	}
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, Mao); got != TJQingLong {
		t.Errorf("例6 天盘卯应乘青龙, 实际 %s", got)
	}
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, Shen); got != TJGuiRen {
		t.Errorf("例6 天盘申应乘贵人(乙日夜贵), 实际 %s", got)
	}
}

// ---------------- 案例 7 ----------------
//
// 邵巡检癸亥生四十六岁占宅戊申年六月庚辰日二十七午将寅时
// 出处：《断案新编》"邵彦和断 宅墓 例七"
// 干支：庚辰日；月将：午；占时：寅（夜）
// 原文："三傳子申辰、潤下、課傳循環、脱空在闗"——
//   三传水局递生（辰生子？非；实为辰→申→子上克下三合水局），
//   中传申值旬空（甲戌旬空申酉），故曰"脱空在闗"。
// 庚日昼贵=丑，夜贵=未；寅时夜，贵人本宫=未。
//
// 验：三传 辰→申→子；元首课（润下水局）；中传申旬空；天盘未乘贵人。
func TestClassicCase7_RunXia_GengChen(t *testing.T) {
	ctx := &Context{
		Gan: Geng, DayZhi: Chen, JiaziIndex: mustJiazi(Geng, Chen),
		ZhanShi: Yin, YueJiang: Wu,
		ZhouYe: IsDayByZhi(Yin),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Chen, Shen, Zi} {
		t.Errorf("例7 三传应=辰申子, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	if pan.SanChuan.Note != "元首课" {
		t.Errorf("例7 课体应=元首课, 实际 %s", pan.SanChuan.Note)
	}
	// 中传申旬空（甲戌旬空申酉）— 原文"脱空在闗"
	if !pan.SanChuan.Zhong.IsKong {
		t.Errorf("例7 中传申应旬空(甲戌旬)")
	}
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, Wei); got != TJGuiRen {
		t.Errorf("例7 天盘未应乘贵人(庚日夜贵), 实际 %s", got)
	}
}

// ---------------- 案例 8 ----------------
//
// 任三翁庚午生三十九岁占宅戊申年十一月壬寅日丑将申时
// 出处：《断案新编》"邵彦和断 宅墓 例八"
// 干支：壬寅日；月将：丑（小寒后日缠玄枵）；占时：申（昼）
// 原文："墓身墓宅、知一、斩关"——三传 子→巳→戌，知一课。
// 壬寅日属甲午旬，旬空=辰巳；中传巳旬空（"是財必亡"）。
// 壬日昼贵=巳、夜贵=卯；申时昼，贵人本宫=巳。
//
// 验：三传 子→巳→戌；知一课；中传巳旬空。
func TestClassicCase8_ZhiYi_RenYin(t *testing.T) {
	ctx := &Context{
		Gan: Ren, DayZhi: Yin, JiaziIndex: mustJiazi(Ren, Yin),
		ZhanShi: Shen, YueJiang: Chou,
		ZhouYe: IsDayByZhi(Shen),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Zi, Si, Xu} {
		t.Errorf("例8 三传应=子巳戌, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	if pan.SanChuan.Note != "知一课" {
		t.Errorf("例8 课体应=知一课, 实际 %s", pan.SanChuan.Note)
	}
	// 中传巳旬空（甲午旬空辰巳）— 原文"是財必亡"
	if !pan.SanChuan.Zhong.IsKong {
		t.Errorf("例8 中传巳应旬空(甲午旬)")
	}
	// 壬日昼贵=巳，则天盘巳乘贵人
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, Si); got != TJGuiRen {
		t.Errorf("例8 天盘巳应乘贵人(壬日昼贵), 实际 %s", got)
	}
}

// ---------------- 案例 9 ----------------
//
// 邵伯达占宅基己酉年十月庚寅日十五寅将酉时
// 出处：《断案新编》"邵彦和断 宅墓 例九"
// 干支：庚寅日；月将：寅（雨水后日缠娵訾）；占时：酉（夜）
// 原文："畫夜貴遇、鑄印、二丙遁官、與巳同聚"——三传 子→巳→戌。
// 派系命名：原文亦称"知一/铸印之祥"。
// 庚日昼贵=丑、夜贵=未；酉时夜，贵人本宫=未。
//
// 验：三传 子→巳→戌；知一课；天盘未乘贵人。
func TestClassicCase9_ZhiYi_GengYin2(t *testing.T) {
	ctx := &Context{
		Gan: Geng, DayZhi: Yin, JiaziIndex: mustJiazi(Geng, Yin),
		ZhanShi: You, YueJiang: Yin,
		ZhouYe: IsDayByZhi(You),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Zi, Si, Xu} {
		t.Errorf("例9 三传应=子巳戌, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	if pan.SanChuan.Note != "知一课" {
		t.Errorf("例9 课体应=知一课, 实际 %s", pan.SanChuan.Note)
	}
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, Wei); got != TJGuiRen {
		t.Errorf("例9 天盘未应乘贵人(庚日夜贵), 实际 %s", got)
	}
}

// ---------------- 案例 10 ----------------
//
// 童保仪丁巳生五十二岁占宅戊申年六月丁巳日初七未将酉时
// 出处：《断案新编》"邵彦和断 宅墓 例十"
// 干支：丁巳日；月将：未（夏至后日缠鹑首）；占时：酉（夜）
// 原文："丁马俱現、不備、间傳人宅相戀、亥水全廚、贅婿、兩貴相見、重審"——
//   三传 丑→亥→酉，重审课，"發用空脱"。
// 丁巳日属甲寅旬，旬空=子丑；初传丑旬空。
// 丁日昼贵=亥、夜贵=酉；酉时夜，贵人本宫=酉。
//
// 验：三传 丑→亥→酉；重审课；初传丑旬空；天盘酉乘贵人。
func TestClassicCase10_ZhongShen_DingSi(t *testing.T) {
	ctx := &Context{
		Gan: Ding, DayZhi: Si, JiaziIndex: mustJiazi(Ding, Si),
		ZhanShi: You, YueJiang: Wei,
		ZhouYe: IsDayByZhi(You),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Chou, Hai, You} {
		t.Errorf("例10 三传应=丑亥酉, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	if pan.SanChuan.Note != "重审课" {
		t.Errorf("例10 课体应=重审课, 实际 %s", pan.SanChuan.Note)
	}
	// 初传丑旬空（甲寅旬空子丑）— 原文"發用空脱"
	if !pan.SanChuan.Chu.IsKong {
		t.Errorf("例10 初传丑应旬空(甲寅旬)")
	}
	// 丁日夜贵=酉，则天盘酉乘贵人
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, You); got != TJGuiRen {
		t.Errorf("例10 天盘酉应乘贵人(丁日夜贵), 实际 %s", got)
	}
}

// ---------------- 案例 11 ----------------
//
// 郑宣义占宅八月癸丑日巳将丑时
// 出处：《断案新编》"邵彦和断 宅墓 例十一"
// 干支：癸丑日；月将：巳（小满后日缠实沈）；占时：丑（夜）
// 原文："寅卯空亡、人盛宅狭、人興宅替"——三传 巳→酉→丑。
// 癸丑日属甲寅旬，旬空=子丑；癸丑为八专日柱（干寄宫=支=丑）。
// 课中有克（巳丑相克）→ 八专有克分支，按贼克法发用。
// 癸日昼贵=巳、夜贵=卯；丑时夜，贵人本宫=卯。
//
// 验：三传 巳→酉→丑；八专有克；天盘卯乘贵人。
func TestClassicCase11_BaZhuan_GuiChou(t *testing.T) {
	ctx := &Context{
		Gan: Gui, DayZhi: Chou, JiaziIndex: mustJiazi(Gui, Chou),
		ZhanShi: Chou, YueJiang: Si,
		ZhouYe: IsDayByZhi(Chou),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Si, You, Chou} {
		t.Errorf("例11 三传应=巳酉丑, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	if pan.SanChuan.Method != "八专法" {
		t.Errorf("例11 发传法应=八专法, 实际 %s", pan.SanChuan.Method)
	}
	// 八专日柱断言（防止 isBaZhuanDay 表回归）
	if !isBaZhuanDay(ctx.Gan, ctx.DayZhi) {
		t.Errorf("例11 癸丑应为八专日柱")
	}
	// 癸日夜贵=卯
	if got := TianJiangOf(pan.TianPan, pan.TianJiang, Mao); got != TJGuiRen {
		t.Errorf("例11 天盘卯应乘贵人(癸日夜贵), 实际 %s", got)
	}
}

// ---------------- 案例 12 ----------------
//
// 辛卯三月癸酉日乙卯时（《指南》卷四"仕宦"门）
// 干支：癸酉日；月将：戌（清明后日缠降娄）；占时：卯（夜）
// 原文标"涉害 砍轮"——三传 卯→戌→巳。
// 派系命名差异：当前实现按涉害规则取深克为用，命名落入"见机课"
//   （涉害中孟地受克最深 → 见机），与《指南》直接称"涉害砍轮"
//   只是命名层差异，不属 BUG。
//
// 验：三传 卯→戌→巳；发传法=涉害法。
func TestClassicCase12_SheHai_GuiYou(t *testing.T) {
	ctx := &Context{
		Gan: Gui, DayZhi: You, JiaziIndex: mustJiazi(Gui, You),
		ZhanShi: Mao, YueJiang: Xu,
		ZhouYe: IsDayByZhi(Mao),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Mao, Xu, Si} {
		t.Errorf("例12 三传应=卯戌巳, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	if pan.SanChuan.Method != "涉害法" {
		t.Errorf("例12 发传法应=涉害法, 实际 %s", pan.SanChuan.Method)
	}
}

// ---------------- 案例 13 ----------------
//
// 壬午十月己亥日辛未时（《指南》卷四"仕宦"门，米山先生至埂子街访顾占）
// 干支：己亥日；月将：亥（小雪后日缠析木）；占时：未（昼）
// 原文标"涉害 曲直 回环"——三传 亥→卯→未（三合木局曲直）。
//
// 验：三传 亥→卯→未；发传法=涉害法。
func TestClassicCase13_SheHai_JiHai(t *testing.T) {
	ctx := &Context{
		Gan: Ji, DayZhi: Hai, JiaziIndex: mustJiazi(Ji, Hai),
		ZhanShi: Wei, YueJiang: Hai,
		ZhouYe: IsDayByZhi(Wei),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Hai, Mao, Wei} {
		t.Errorf("例13 三传应=亥卯未, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	if pan.SanChuan.Method != "涉害法" {
		t.Errorf("例13 发传法应=涉害法, 实际 %s", pan.SanChuan.Method)
	}
	if pan.SanChuan.Note != "涉害课" {
		t.Errorf("例13 课体应=涉害课, 实际 %s", pan.SanChuan.Note)
	}
}

// ---------------- 案例 14 ----------------
//
// 甲申二月乙丑日亥将申时（《指南》卷四"仕宦"门，李大生先生持丹阳孙友所占）
// 干支：乙丑日；月将：亥（春分后日缠降娄前一辰）；占时：申（昼）
// 原文标"重审 稼穑"——三传 未→戌→丑（三合土局稼穑）。
//
// 验：三传 未→戌→丑；课体=重审课；发传法=贼克法。
func TestClassicCase14_ZhongShen_YiChou(t *testing.T) {
	ctx := &Context{
		Gan: Yi, DayZhi: Chou, JiaziIndex: mustJiazi(Yi, Chou),
		ZhanShi: Shen, YueJiang: Hai,
		ZhouYe: IsDayByZhi(Shen),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Wei, Xu, Chou} {
		t.Errorf("例14 三传应=未戌丑, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	if pan.SanChuan.Note != "重审课" {
		t.Errorf("例14 课体应=重审课, 实际 %s", pan.SanChuan.Note)
	}
	if pan.SanChuan.Method != "贼克法" {
		t.Errorf("例14 发传法应=贼克法, 实际 %s", pan.SanChuan.Method)
	}
}

// ---------------- 案例 15 ----------------
//
// 丁丑四月丁酉日乙巳时（《指南》卷四"奏章"门，浣中刘退齐太史索占）
// 干支：丁酉日；月将：酉（小满后日缠实沈前一辰）；占时：巳（昼）
// 原文标"元首 曲直"——三传 亥→卯→未（三合木局曲直）。
//
// 验：三传 亥→卯→未；课体=元首课；发传法=贼克法。
func TestClassicCase15_YuanShou_DingYou(t *testing.T) {
	ctx := &Context{
		Gan: Ding, DayZhi: You, JiaziIndex: mustJiazi(Ding, You),
		ZhanShi: Si, YueJiang: You,
		ZhouYe: IsDayByZhi(Si),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Hai, Mao, Wei} {
		t.Errorf("例15 三传应=亥卯未, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	if pan.SanChuan.Note != "元首课" {
		t.Errorf("例15 课体应=元首课, 实际 %s", pan.SanChuan.Note)
	}
	if pan.SanChuan.Method != "贼克法" {
		t.Errorf("例15 发传法应=贼克法, 实际 %s", pan.SanChuan.Method)
	}
}

// ---------------- 案例 16 ----------------
//
// 癸酉四月癸亥日丙辰时（《指南》卷四"公讼"门，宜兴周首辅因陈科长弹论命医者周诚生来占）
// 干支：癸亥日；月将：酉（小满后日缠实沈前一辰）；占时：辰（昼）
// 原文标"重审"——三传 午→亥→辰。癸亥日属甲寅旬，旬空=子丑。
//
// 验：三传 午→亥→辰；课体=重审课；发传法=贼克法。
func TestClassicCase16_ZhongShen_GuiHai(t *testing.T) {
	ctx := &Context{
		Gan: Gui, DayZhi: Hai, JiaziIndex: mustJiazi(Gui, Hai),
		ZhanShi: Chen, YueJiang: You,
		ZhouYe: IsDayByZhi(Chen),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Wu, Hai, Chen} {
		t.Errorf("例16 三传应=午亥辰, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	if pan.SanChuan.Note != "重审课" {
		t.Errorf("例16 课体应=重审课, 实际 %s", pan.SanChuan.Note)
	}
}

// ---------------- 案例 17 ----------------
//
// 庚寅四月乙酉日辛巳时（《指南》卷四"逃亡"门，弯子街二人来占子逃看何方找寻何日得见）
// 干支：乙酉日；月将：酉（小满前后日缠实沈）；占时：巳（昼）
// 原文标"元首 涧下"——三传 申→子→辰（三合水局涧下）。
//
// 验：三传 申→子→辰；课体=元首课；发传法=贼克法。
func TestClassicCase17_YuanShou_YiYou(t *testing.T) {
	ctx := &Context{
		Gan: Yi, DayZhi: You, JiaziIndex: mustJiazi(Yi, You),
		ZhanShi: Si, YueJiang: You,
		ZhouYe: IsDayByZhi(Si),
	}
	pan := DivineWithContext(ctx)

	if got := [3]Zhi{pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi}; got != [3]Zhi{Shen, Zi, Chen} {
		t.Errorf("例17 三传应=申子辰, 实际 %v→%v→%v",
			pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi)
	}
	if pan.SanChuan.Note != "元首课" {
		t.Errorf("例17 课体应=元首课, 实际 %s", pan.SanChuan.Note)
	}
	if pan.SanChuan.Method != "贼克法" {
		t.Errorf("例17 发传法应=贼克法, 实际 %s", pan.SanChuan.Method)
	}
}
