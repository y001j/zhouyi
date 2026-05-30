package qimen

// 本文件集中五行相关的常量与小工具，用于：
//   - 每格天盘干/地盘干的五行
//   - 九星五行（做门迫、月令旺衰判断）
//   - 八门五行
//   - 宫位（地支）五行
//   - 干入墓表、地支击刑表

// ============ 五行名 ============
// 五行统一用单字字符串："金", "木", "水", "火", "土"

// GanWuXing 十天干五行
var GanWuXing = map[string]string{
	"甲": "木", "乙": "木",
	"丙": "火", "丁": "火",
	"戊": "土", "己": "土",
	"庚": "金", "辛": "金",
	"壬": "水", "癸": "水",
}

// ZhiWuXing 十二地支五行
var ZhiWuXing = map[string]string{
	"寅": "木", "卯": "木",
	"巳": "火", "午": "火",
	"申": "金", "酉": "金",
	"亥": "水", "子": "水",
	"辰": "土", "戌": "土", "丑": "土", "未": "土",
}

// StarWuXing 九星五行
//
//	天蓬水、天芮土、天冲木、天辅木、天禽土、
//	天心金、天柱金、天任土、天英火
var StarWuXing = map[string]string{
	"天蓬": "水",
	"天芮": "土",
	"天冲": "木",
	"天辅": "木",
	"天禽": "土",
	"天心": "金",
	"天柱": "金",
	"天任": "土",
	"天英": "火",
}

// DoorWuXing 八门五行
//
//	坎 休门 水、艮 生门 土、震 伤门 木、巽 杜门 木、
//	离 景门 火、坤 死门 土、兑 惊门 金、乾 开门 金
var DoorWuXing = map[string]string{
	"休门": "水",
	"生门": "土",
	"伤门": "木",
	"杜门": "木",
	"景门": "火",
	"死门": "土",
	"惊门": "金",
	"开门": "金",
}

// PalaceWuXingByFei 宫位（飞星序 0..8）的五行
var PalaceWuXingByFei = [9]string{
	"水", // 坎一
	"土", // 坤二
	"木", // 震三
	"木", // 巽四
	"土", // 中五
	"金", // 乾六
	"金", // 兑七
	"土", // 艮八
	"火", // 离九
}

// ============ 生克关系 ============

// WuXingSheng a 生 b 吗？（五行相生：木生火 / 火生土 / 土生金 / 金生水 / 水生木）
func WuXingSheng(a, b string) bool {
	return map[string]string{"木": "火", "火": "土", "土": "金", "金": "水", "水": "木"}[a] == b
}

// WuXingKe a 克 b 吗？（五行相克：木克土 / 土克水 / 水克火 / 火克金 / 金克木）
func WuXingKe(a, b string) bool {
	return map[string]string{"木": "土", "土": "水", "水": "火", "火": "金", "金": "木"}[a] == b
}

// WuXingRelation 返回 a 对 b 的关系描述
//
//	"比和" / "a生b" / "b生a" / "a克b" / "b克a"
func WuXingRelation(a, b string) string {
	switch {
	case a == b:
		return "比和"
	case WuXingSheng(a, b):
		return a + "生" + b
	case WuXingSheng(b, a):
		return b + "生" + a
	case WuXingKe(a, b):
		return a + "克" + b
	case WuXingKe(b, a):
		return b + "克" + a
	}
	return "无直接生克"
}

// ============ 月令旺衰（当月五行 → 对象五行 得到旺/相/休/囚/死） ============

// WangShuaiByMonth 给定月支与对象的五行，返回旺衰名。
//
//	旺：对象五行 == 当月值五行
//	相：当月值五行生对象（我所生为相）—— 修正：实际口诀是"**对象生当月**为休、**当月生对象**为相"
//	休：对象生当月值（退居休闲）
//	囚：对象克当月值（受制为囚）
//	死：当月值克对象（被克为死）
//
// 实现参考 WangXiangXiuQiuSi[月支] = [旺,相,休,囚,死] 的五行顺序直接反查。
func WangShuaiByMonth(monthZhi, objWX string) string {
	arr, ok := WangXiangXiuQiuSi[monthZhi]
	if !ok {
		return ""
	}
	names := [5]string{"旺", "相", "休", "囚", "死"}
	for i, wx := range arr {
		if wx == objWX {
			return names[i]
		}
	}
	return ""
}

// ============ 入墓 ============
//
// 五行入墓法（四库说）：
//
//	水土入辰（水墓辰；土也入辰为"水土同墓"）
//	木入未
//	火入戌
//	金入丑
//
// 主流采用："丙丁火入戌、庚辛金入丑、壬癸水入辰、甲乙木入未、戊己土入戌"（土从火）
// 另一派"戊己土入辰"（土从水）。本实现采用主流说法。
var StemMuZhi = map[string]string{
	"甲": "未", "乙": "未",
	"丙": "戌", "丁": "戌",
	"戊": "戌", "己": "戌",
	"庚": "丑", "辛": "丑",
	"壬": "辰", "癸": "辰",
}

// PalaceContainsZhi 某飞星索引宫是否含有指定地支
func PalaceContainsZhi(palaceFei int, zhi string) bool {
	if palaceFei < 0 || palaceFei >= 9 {
		return false
	}
	for _, z := range PalaceToZhi[palaceFei] {
		if z == zhi {
			return true
		}
	}
	return false
}

// IsStemInMu 判断天干 stem 在飞星索引 palaceFei 所落之宫是否入墓
func IsStemInMu(stem string, palaceFei int) bool {
	mu, ok := StemMuZhi[stem]
	if !ok {
		return false
	}
	return PalaceContainsZhi(palaceFei, mu)
}

// ============ 六仪击刑 ============
//
// 古法（《奇门宝鉴》）：
//
//	甲子旬首戊 → 击刑加坎一宫（子自刑——实际是"戊加子为击刑"的古注变体，
//	                     主流写作甲子戊临子宫自刑）
//	甲戌旬首己 → 击刑加艮八宫（丑寅所在；己加未）
//	甲申旬首庚 → 击刑加艮八宫（寅）
//	甲午旬首辛 → 击刑加离九宫（午自刑）
//	甲辰旬首壬 → 击刑加巽四宫（辰自刑；辰位）
//	甲寅旬首癸 → 击刑加巽四宫（巳）
//
// 本表按"旬首六仪 → 被刑之飞星索引"给出（源自奇门宝鉴 v2 报告）。
var LiuyiJiXingPalace = map[string]int{
	"戊": 0, // 甲子戊 加坎1（子刑）
	"己": 7, // 甲戌己 加艮8
	"庚": 7, // 甲申庚 加艮8（寅）
	"辛": 8, // 甲午辛 加离9（午自刑）
	"壬": 3, // 甲辰壬 加巽4（辰自刑的变形）
	"癸": 3, // 甲寅癸 加巽4（巳）
}
