package liuren

// NianMingInfo 年命上神信息
type NianMingInfo struct {
	Zhi       Zhi       // 本命或行年的地支
	Upper     Zhi       // 在天盘上的乘神（即该地盘位之上方天神）
	TianJiang TianJiang // 所乘天将
	LiuQin    LiuQin    // 相对日干的六亲
	IsKong    bool      // 是否旬空
	Ying      YingType  // 应机类型（对日干的生克关系）
}

// YingType 年命上神相对日干的应机类型
type YingType int

const (
	YingBiHe YingType = iota // 比应（同五行）
	YingSheng                // 救应（生我）
	YingKe                   // 损应（克我）
	YingTuo                  // 脱应（我生彼，泄气）
	YingZhi                  // 制应（我克彼）
)

var yingNames = [5]string{"比应", "救应", "损应", "脱应", "制应"}
var yingDesc = [5]string{
	"年命上神与日干同五行，比和帮身，性本和顺",
	"年命上神生日干，得扶助之力，凶事可解、吉事更增",
	"年命上神克日干，外力加身，凶事加祸、吉事减分",
	"日干生年命上神，我气泄耗，事虽成而费力",
	"日干克年命上神，我制于事，利主动进取",
}

func (y YingType) String() string { return yingNames[y] }
func (y YingType) Desc() string   { return yingDesc[y] }

// ClassifyYing 按日干与年命上神五行关系判定应机类型
func ClassifyYing(dayG Gan, upper Zhi) YingType {
	dw := GanWuXing[dayG]
	uw := ZhiWuXing[upper]
	switch {
	case dw == uw:
		return YingBiHe
	case WuXingGenerates(uw, dw):
		return YingSheng
	case WuXingOvercomes(uw, dw):
		return YingKe
	case WuXingGenerates(dw, uw):
		return YingTuo
	case WuXingOvercomes(dw, uw):
		return YingZhi
	}
	return YingBiHe
}

// BenMingInfo 盘面附着的本命/行年两路信息
type BenMingInfo struct {
	BenMing  *NianMingInfo // 本命（基于生肖/出生年地支）
	XingNian *NianMingInfo // 行年（按虚岁男顺女逆推）
	BMXNRel  string        // 本命与行年的互动关系：合/冲/刑/同位/比和/生/克/无关；空表未推
	BMXNDesc string        // 关系释义（《大全》卷三 心印赋"日命相合为福祥"）
}

// ClassifyBenMingXingNian 判定本命 vs 行年的互动关系
//
// 据《大全》卷三 p558 心印赋：「日命相合为福祥」「命若克日见灾秧」
// 此处广义用之于本命与行年的关系：
//
//	同位/比和：本命=行年 或同五行     → 命行一体，主稳
//	六合：    子丑、寅亥、卯戌、辰酉、巳申、午未 → 大吉，福祥
//	六冲：    差 6 位                 → 反复动荡
//	三刑：    寅巳申 / 丑戌未 / 子卯  → 刑伤之兆
//	相生：    本命生行年（命生行）→ 命扶行；行生命 → 行济命
//	相克：    一者克另一者     → 灾秧/压制
//	无关：    其他              → 平淡
func ClassifyBenMingXingNian(bm, xn Zhi) (string, string) {
	if bm < 0 || xn < 0 {
		return "", ""
	}
	if bm == xn {
		return "同位", "本命与行年同位，命运一体，主流年与生肖大方向一致；吉凶看其上神。"
	}
	if ZhiWuXing[bm] == ZhiWuXing[xn] {
		return "比和", "本命与行年五行比和，气性相投，主稳定无大变；助力可借而不强。"
	}
	if zhiLiuhe(bm, xn) {
		return "六合", "本命与行年六合，大吉之兆，主福祥贵助、人际亨通（心印赋『日命相合为福祥』）。"
	}
	if (int(bm)+6)%12 == int(xn) {
		return "六冲", "本命与行年六冲，主流年反复、心神不宁、宜守不宜攻。"
	}
	// 三刑：寅刑巳、巳刑申、申刑寅；丑刑戌、戌刑未、未刑丑；子刑卯、卯刑子；辰午酉亥自刑
	xingPairs := [][2]Zhi{
		{Yin, Si}, {Si, Shen}, {Shen, Yin},
		{Chou, Xu}, {Xu, Wei}, {Wei, Chou},
		{Zi, Mao}, {Mao, Zi},
	}
	for _, p := range xingPairs {
		if (bm == p[0] && xn == p[1]) || (bm == p[1] && xn == p[0]) {
			return "三刑", "本命与行年三刑，主刑伤之兆、是非缠身、宜避锋芒。"
		}
	}
	bw, xw := ZhiWuXing[bm], ZhiWuXing[xn]
	if WuXingGenerates(bw, xw) {
		return "命生行", "本命生行年，本命扶持流年，主自身付出、有所积累；耗气而非受益。"
	}
	if WuXingGenerates(xw, bw) {
		return "行生命", "行年生本命，流年补给本命，主受外援之助、机缘自来。"
	}
	if WuXingOvercomes(bw, xw) {
		return "命克行", "本命克行年，主自身压制流年，需破障而行；占求事虽阻终成。"
	}
	if WuXingOvercomes(xw, bw) {
		return "行克命", "行年克本命，主流年压本命，灾秧加身、宜韬光养晦（心印赋『命若克日见灾秧』）。"
	}
	return "无关", "本命与行年无明显生克合冲关系，流年平淡。"
}

// ComputeXingNian 推算行年地支。
//
//	男自 1 岁起寅顺行：1→寅、2→卯、3→辰 ... 12→丑、13→寅...
//	女自 1 岁起申逆行：1→申、2→未、3→午 ... 12→酉、13→申...
//	虚岁：出生年算 1，以当前农历年减出生年 +1。
func ComputeXingNian(gender string, virtualAge int) Zhi {
	if virtualAge < 1 {
		virtualAge = 1
	}
	step := (virtualAge - 1) % 12
	if gender == "女" || gender == "female" || gender == "f" || gender == "F" {
		// 申起逆行
		return Zhi(((int(Shen)-step)%12 + 12) % 12)
	}
	// 默认男：寅起顺行
	return Zhi((int(Yin) + step) % 12)
}

// BuildNianMing 以给定地支查询其在盘面上的乘神
func BuildNianMing(pan *Pan, z Zhi) NianMingInfo {
	upper := pan.TianPan[z]
	return NianMingInfo{
		Zhi:       z,
		Upper:     upper,
		TianJiang: pan.TianJiang[z],
		LiuQin:    LiuQinOfZhiByGan(upper, pan.Ctx.Gan),
		IsKong:    IsXunKong(upper, pan.Ctx.JiaziIndex),
		Ying:      ClassifyYing(pan.Ctx.Gan, upper),
	}
}
