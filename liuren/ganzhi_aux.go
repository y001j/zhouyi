package liuren

// 天干寄行（长生十二宫/禄/旺/墓/胎/绝/死气）辅助表。
//
// 以五行十二宫（长生→沐浴→冠带→临官→帝旺→衰→病→死→墓→绝→胎→养）为模，阳干顺行、阴干逆行。
// 这里只取常用的几项。

// GanLuZhi 干禄：甲禄寅、乙禄卯、丙戊禄巳、丁己禄午、庚禄申、辛禄酉、壬禄亥、癸禄子
var GanLuZhi = [10]Zhi{Yin, Mao, Si, Wu, Si, Wu, Shen, You, Hai, Zi}

// GanWangZhi 干旺（帝旺）：甲旺卯、乙旺寅、丙戊旺午、丁己旺巳、庚旺酉、辛旺申、壬旺子、癸旺亥
var GanWangZhi = [10]Zhi{Mao, Yin, Wu, Si, Wu, Si, You, Shen, Zi, Hai}

// GanMuZhi 干墓：甲墓未、乙墓戌、丙戊墓戌、丁己墓丑、庚墓丑、辛墓辰、壬墓辰、癸墓未
//
//	（以五行阴阳十二宫："墓"位）
var GanMuZhi = [10]Zhi{Wei, Xu, Xu, Chou, Xu, Chou, Chou, Chen, Chen, Wei}

// GanZhangShengZhi 干长生：甲长生亥、乙长生午、丙戊长生寅、丁己长生酉、庚长生巳、辛长生子、壬长生申、癸长生卯
var GanZhangShengZhi = [10]Zhi{Hai, Wu, Yin, You, Yin, You, Si, Zi, Shen, Mao}

// GanTaiZhi 干胎：甲胎酉、乙胎申、丙戊胎子、丁己胎亥、庚胎卯、辛胎寅、壬胎午、癸胎巳
var GanTaiZhi = [10]Zhi{You, Shen, Zi, Hai, Zi, Hai, Mao, Yin, Wu, Si}

// GanJueZhi 干绝：甲绝申、乙绝酉、丙戊绝亥、丁己绝子、庚绝寅、辛绝卯、壬绝巳、癸绝午
var GanJueZhi = [10]Zhi{Shen, You, Hai, Zi, Hai, Zi, Yin, Mao, Si, Wu}

// WuXingDeathZhi 月内死气（按日干五行）：
//
//	金日（庚辛）死气子；木日（甲乙）死气午；水土日（壬癸戊己）死气卯；火日（丙丁）死气酉
var GanDeathZhi = map[WuXing]Zhi{Jin: Zi, Mu: Wu, Shui: Mao, Huo: You, Tu: Mao}

// 地支三刑
//
//	三字刑（三合）：寅巳申、丑戌未
//	二字刑：子卯
//	自刑：辰辰、午午、酉酉、亥亥
func isThreeXingSet(a, b, c Zhi) bool {
	set := map[Zhi]bool{a: true, b: true, c: true}
	g1 := map[Zhi]bool{Yin: true, Si: true, Shen: true}
	g2 := map[Zhi]bool{Chou: true, Xu: true, Wei: true}
	match := 0
	for z := range g1 {
		if set[z] {
			match++
		}
	}
	if match == 3 {
		return true
	}
	match = 0
	for z := range g2 {
		if set[z] {
			match++
		}
	}
	return match == 3
}

// 六害：子未、丑午、寅巳、卯辰、申亥、酉戌
var liuHaiPairs = map[Zhi]Zhi{
	Zi: Wei, Wei: Zi,
	Chou: Wu, Wu: Chou,
	Yin: Si, Si: Yin,
	Mao: Chen, Chen: Mao,
	Shen: Hai, Hai: Shen,
	You: Xu, Xu: You,
}

func isLiuHai(a, b Zhi) bool { return liuHaiPairs[a] == b }

// 六合：子丑、寅亥、卯戌、辰酉、巳申、午未
var liuHePairs = map[Zhi]Zhi{
	Zi: Chou, Chou: Zi,
	Yin: Hai, Hai: Yin,
	Mao: Xu, Xu: Mao,
	Chen: You, You: Chen,
	Si: Shen, Shen: Si,
	Wu: Wei, Wei: Wu,
}

func isLiuHe(a, b Zhi) bool { return liuHePairs[a] == b }

// 六冲：子午、丑未、寅申、卯酉、辰戌、巳亥
func isLiuChong(a, b Zhi) bool { return (int(a)-int(b)+12)%12 == 6 }

// 四季关神：春丑、夏辰、秋未、冬戌（以干支月支或阳历月简化判定）
//
//	正/二/三月 (寅卯辰) 为春 → 关丑
//	四/五/六月 (巳午未) 为夏 → 关辰
//	七/八/九月 (申酉戌) 为秋 → 关未
//	十/冬/腊月 (亥子丑) 为冬 → 关戌
func guanShen(monthZhi Zhi) Zhi {
	switch monthZhi {
	case Yin, Mao, Chen:
		return Chou
	case Si, Wu, Wei:
		return Chen
	case Shen, You, Xu:
		return Wei
	case Hai, Zi, Chou:
		return Xu
	}
	return -1
}

// 地支之墓：辰为水墓、戌为火墓、丑为金墓、未为木墓
//
//	（支的"墓库"与天干墓略有差异，支墓以其所属五行的库地支为准）
func zhiMu(z Zhi) Zhi {
	switch ZhiWuXing[z] {
	case Shui:
		return Chen
	case Huo:
		return Xu
	case Jin:
		return Chou
	case Mu:
		return Wei
	case Tu:
		// 土无独立墓，本支自旺自藏
		return z
	}
	return -1
}

// 连茹：三传成等差 ±1（顺连/逆连）
func isLianRu(a, b, c Zhi) (forward, backward bool) {
	d1 := (int(b) - int(a) + 12) % 12
	d2 := (int(c) - int(b) + 12) % 12
	return d1 == 1 && d2 == 1, d1 == 11 && d2 == 11
}

// liuQinFromWuXing 以五行判断六亲（相对日干）
func liuQinFromWuXingRelation(dayW, targetW WuXing) LiuQin {
	switch {
	case targetW == dayW:
		return LQXiongDi
	case WuXingGenerates(targetW, dayW):
		return LQFuMu
	case WuXingGenerates(dayW, targetW):
		return LQZiSun
	case WuXingOvercomes(dayW, targetW):
		return LQQiCai
	case WuXingOvercomes(targetW, dayW):
		return LQGuanGui
	}
	return LQXiongDi
}

// 求某类六亲对应的五行：给定日干，找出哪个五行属某六亲
func wuXingForLiuQin(dayG Gan, lq LiuQin) WuXing {
	dayW := GanWuXing[dayG]
	switch lq {
	case LQXiongDi:
		return dayW
	case LQFuMu:
		// 生我者
		for w := Jin; w <= Tu; w++ {
			if WuXingGenerates(w, dayW) {
				return w
			}
		}
	case LQZiSun:
		for w := Jin; w <= Tu; w++ {
			if WuXingGenerates(dayW, w) {
				return w
			}
		}
	case LQQiCai:
		for w := Jin; w <= Tu; w++ {
			if WuXingOvercomes(dayW, w) {
				return w
			}
		}
	case LQGuanGui:
		for w := Jin; w <= Tu; w++ {
			if WuXingOvercomes(w, dayW) {
				return w
			}
		}
	}
	return dayW
}
