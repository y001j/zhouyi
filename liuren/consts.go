// Package liuren 实现大六壬起课、排盘与解课。
//
// 基础约定：
//   - 地支索引 0..11 依次为 子丑寅卯辰巳午未申酉戌亥
//   - 天干索引 0..9  依次为 甲乙丙丁戊己庚辛壬癸
//   - 五行 Jin/Mu/Shui/Huo/Tu 分别为 金木水火土
package liuren

// ============ 地支 ============

type Zhi int

const (
	Zi Zhi = iota
	Chou
	Yin
	Mao
	Chen
	Si
	Wu
	Wei
	Shen
	You
	Xu
	Hai
)

// ZhiNames 十二地支中文名（索引 = Zhi 值）
var ZhiNames = [12]string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

// ZhiBieMing 十二地支别名（用于月将显示）
var ZhiBieMing = [12]string{"神后", "大吉", "功曹", "太冲", "天罡", "太乙", "胜光", "小吉", "传送", "从魁", "河魁", "登明"}

// ZhiShengXiao 十二生肖
var ZhiShengXiao = [12]string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}

func (z Zhi) String() string { return ZhiNames[z] }

// ParseZhi 由中文名解析地支，未知返回 -1
func ParseZhi(s string) Zhi {
	for i, n := range ZhiNames {
		if n == s {
			return Zhi(i)
		}
	}
	return -1
}

// ============ 天干 ============

type Gan int

const (
	Jia Gan = iota
	Yi
	Bing
	Ding
	Wu1 // 戊（避免与地支 Wu 重名）
	Ji
	Geng
	Xin
	Ren
	Gui
)

var GanNames = [10]string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}

func (g Gan) String() string { return GanNames[g] }

func ParseGan(s string) Gan {
	for i, n := range GanNames {
		if n == s {
			return Gan(i)
		}
	}
	return -1
}

// ============ 五行 ============

type WuXing int

const (
	Jin WuXing = iota // 金
	Mu                // 木
	Shui              // 水
	Huo               // 火
	Tu                // 土
)

var WuXingNames = [5]string{"金", "木", "水", "火", "土"}

func (w WuXing) String() string { return WuXingNames[w] }

// ZhiWuXing 十二地支五行：子水 丑土 寅木 卯木 辰土 巳火 午火 未土 申金 酉金 戌土 亥水
var ZhiWuXing = [12]WuXing{Shui, Tu, Mu, Mu, Tu, Huo, Huo, Tu, Jin, Jin, Tu, Shui}

// GanWuXing 十天干五行：甲乙木 丙丁火 戊己土 庚辛金 壬癸水
var GanWuXing = [10]WuXing{Mu, Mu, Huo, Huo, Tu, Tu, Jin, Jin, Shui, Shui}

// ZhiYang 十二地支阴阳：子寅辰午申戌为阳，丑卯巳未酉亥为阴
var ZhiYang = [12]bool{true, false, true, false, true, false, true, false, true, false, true, false}

// GanYang 十天干阴阳：甲丙戊庚壬为阳，乙丁己辛癸为阴
var GanYang = [10]bool{true, false, true, false, true, false, true, false, true, false}

// ============ 十干寄宫 ============

// GanJiGong 十天干寄宫表：
//
//	甲→寅 乙→辰 丙→巳 丁→未 戊→巳 己→未 庚→申 辛→戌 壬→亥 癸→丑
//
// 子午卯酉为四正位，不作寄宫。
var GanJiGong = [10]Zhi{Yin, Chen, Si, Wei, Si, Wei, Shen, Xu, Hai, Chou}

// ============ 干合（六合） ============

// GanLiuhe 十干六合：甲己合化土、乙庚合化金、丙辛合化水、丁壬合化木、戊癸合化火
//
// 别责法阳日发用要查"日干六合所对应阴干的寄宫"上之天盘神。
var GanLiuhe = [10]Gan{
	Jia:  Ji,
	Yi:   Geng,
	Bing: Xin,
	Ding: Ren,
	Wu1:  Gui,
	Ji:   Jia,
	Geng: Yi,
	Xin:  Bing,
	Ren:  Ding,
	Gui:  Wu1,
}

// ============ 月将（太阳躔次） ============

// 中气 → 月将 对照：
//
//	雨水后 → 亥（登明）     处暑后 → 巳（太乙）
//	春分后 → 戌（河魁）     秋分后 → 辰（天罡）
//	谷雨后 → 酉（从魁）     霜降后 → 卯（太冲）
//	小满后 → 申（传送）     小雪后 → 寅（功曹）
//	夏至后 → 未（小吉）     冬至后 → 丑（大吉）
//	大暑后 → 午（胜光）     大寒后 → 子（神后）
var QiToYueJiang = map[string]Zhi{
	"雨水": Hai, "春分": Xu, "谷雨": You, "小满": Shen,
	"夏至": Wei, "大暑": Wu, "处暑": Si, "秋分": Chen,
	"霜降": Mao, "小雪": Yin, "冬至": Chou, "大寒": Zi,
}

// ============ 贵人（天乙） ============

// GuiRenByGan 贵人表：日干 → [昼贵, 夜贵]
//
//	甲戊庚牛羊 → 昼丑 夜未
//	乙己鼠猴乡 → 昼子 夜申
//	丙丁猪鸡位 → 昼亥 夜酉
//	壬癸蛇兔藏 → 昼巳 夜卯
//	六辛逢马虎 → 昼午 夜寅
//
// 派系备注（《大全》卷一 p495 神图）：
//   - 主流口诀（本表所采）：「丙丁猪鸡位」即丙丁同取昼亥夜酉
//   - 四库本异说：将"丁"归"鼠猴"组（即 Ding 同 Yi/Ji 取昼子夜申）；此说不普及
//   - 通行《指南》《粹言》《心镜》皆与本表一致
// 故本表保留主流派；如需四库本派可在调用时覆盖 GuiRenByGan[Ding]。
var GuiRenByGan = [10][2]Zhi{
	Jia:  {Chou, Wei},
	Yi:   {Zi, Shen},
	Bing: {Hai, You},
	Ding: {Hai, You}, // 主流派；四库本异说为 {Zi, Shen}
	Wu1:  {Chou, Wei},
	Ji:   {Zi, Shen},
	Geng: {Chou, Wei},
	Xin:  {Wu, Yin},
	Ren:  {Si, Mao},
	Gui:  {Si, Mao},
}

// ============ 十二天将 ============

type TianJiang int

const (
	TJGuiRen   TianJiang = iota // 贵人
	TJTengShe                   // 腾蛇
	TJZhuQue                    // 朱雀
	TJLiuHe                     // 六合
	TJGouChen                   // 勾陈
	TJQingLong                  // 青龙
	TJTianKong                  // 天空
	TJBaiHu                     // 白虎
	TJTaiChang                  // 太常
	TJXuanWu                    // 玄武
	TJTaiYin                    // 太阴
	TJTianHou                   // 天后
)

var TianJiangNames = [12]string{"贵人", "腾蛇", "朱雀", "六合", "勾陈", "青龙", "天空", "白虎", "太常", "玄武", "太阴", "天后"}
var TianJiangShort = [12]string{"贵", "蛇", "雀", "合", "勾", "龙", "空", "虎", "常", "武", "阴", "后"}

func (t TianJiang) String() string { return TianJiangNames[t] }
func (t TianJiang) Short() string  { return TianJiangShort[t] }

// TianJiangMeaning 天将象意（断课基础）
var TianJiangMeaning = [12]string{
	"首领、长辈、贵人扶助，主贵气",
	"惊恐怪异、文书忧虑，性多疑",
	"口舌是非、文书信息，主火急",
	"婚姻和合、交易媒介，主合作",
	"勾连牵制、田土官讼，主拖延",
	"财喜庆贺、升迁喜事，主发达",
	"虚诈欺瞒、奴仆小人，主空耗",
	"疾病丧服、刀兵血光，主凶急",
	"宴席服饰、长辈之物，主平和",
	"盗贼暗昧、奸私欺骗，主暗失",
	"暗中扶助、女子阴私，主隐蔽",
	"女子后妃、阴私婚姻，主柔暗",
}

// ============ 六亲（以日干为我） ============

type LiuQin int

const (
	LQFuMu   LiuQin = iota // 父母：生我者
	LQXiongDi               // 兄弟：同我者
	LQZiSun                 // 子孙：我生者
	LQQiCai                 // 妻财：我克者
	LQGuanGui               // 官鬼：克我者
)

var LiuQinNames = [5]string{"父母", "兄弟", "子孙", "妻财", "官鬼"}

func (l LiuQin) String() string { return LiuQinNames[l] }

// ============ 旬空 ============

// XunKongByXunIndex 旬首序号（0,10,20,30,40,50）→ 空亡两支
//
//	甲子旬空戌亥；甲戌旬空申酉；甲申旬空午未；
//	甲午旬空辰巳；甲辰旬空寅卯；甲寅旬空子丑。
var XunKongByXunIndex = map[int][2]Zhi{
	0:  {Xu, Hai},
	10: {Shen, You},
	20: {Wu, Wei},
	30: {Chen, Si},
	40: {Yin, Mao},
	50: {Zi, Chou},
}

// ============ 生克辅助 ============

// WuXingOvercomes w1 是否克 w2：金克木、木克土、土克水、水克火、火克金
func WuXingOvercomes(w1, w2 WuXing) bool {
	switch w1 {
	case Jin:
		return w2 == Mu
	case Mu:
		return w2 == Tu
	case Shui:
		return w2 == Huo
	case Huo:
		return w2 == Jin
	case Tu:
		return w2 == Shui
	}
	return false
}

// WuXingGenerates w1 是否生 w2：金生水、水生木、木生火、火生土、土生金
func WuXingGenerates(w1, w2 WuXing) bool {
	switch w1 {
	case Jin:
		return w2 == Shui
	case Shui:
		return w2 == Mu
	case Mu:
		return w2 == Huo
	case Huo:
		return w2 == Tu
	case Tu:
		return w2 == Jin
	}
	return false
}

// RelationOfZhi 天神（上）与地神（下）的关系描述，用于四课注记
//
//	返回值：上克下 / 下贼上 / 上生下 / 下生上 / 比和
func RelationOfZhi(upper, lower Zhi) string {
	wu := ZhiWuXing[upper]
	wl := ZhiWuXing[lower]
	switch {
	case WuXingOvercomes(wu, wl):
		return "上克下"
	case WuXingOvercomes(wl, wu):
		return "下贼上"
	case WuXingGenerates(wu, wl):
		return "上生下"
	case WuXingGenerates(wl, wu):
		return "下生上"
	default:
		return "比和"
	}
}

// LiuQinOfZhiByGan 以日干为我，判断一个地支的六亲属性（按五行）
func LiuQinOfZhiByGan(z Zhi, g Gan) LiuQin {
	wm := GanWuXing[g]
	wz := ZhiWuXing[z]
	switch {
	case wz == wm:
		return LQXiongDi
	case WuXingGenerates(wz, wm):
		return LQFuMu
	case WuXingGenerates(wm, wz):
		return LQZiSun
	case WuXingOvercomes(wm, wz):
		return LQQiCai
	case WuXingOvercomes(wz, wm):
		return LQGuanGui
	}
	return LQXiongDi
}

// IsXunKong 判断某支对给定旬（以日干支的六十甲子序号定旬）是否旬空
func IsXunKong(z Zhi, jiaziIndex int) bool {
	xunStart := (jiaziIndex / 10) * 10
	pair, ok := XunKongByXunIndex[xunStart]
	if !ok {
		return false
	}
	return z == pair[0] || z == pair[1]
}
