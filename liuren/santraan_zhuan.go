package liuren

// 地支刑：寅刑巳、巳刑申、申刑寅（恃势之刑）；
//
//	丑刑戌、戌刑未、未刑丑（无恩之刑）；
//	子刑卯、卯刑子（无礼之刑）；
//	辰、午、酉、亥自刑。
var xingOf = map[Zhi]Zhi{
	Yin:  Si,
	Si:   Shen,
	Shen: Yin,
	Chou: Xu,
	Xu:   Wei,
	Wei:  Chou,
	Zi:   Mao,
	Mao:  Zi,
	Chen: Chen,
	Wu:   Wu,
	You:  You,
	Hai:  Hai,
}

// isSelfXing 是否自刑
func isSelfXing(z Zhi) bool {
	return z == Chen || z == Wu || z == You || z == Hai
}

// fuYinMethod 伏吟：月将=占时，天盘各神居本位。
//
//	有克：按贼克法取上神
//	无克：
//	  阳日（刚干）取干上神 → 自任课
//	  阴日（柔干）取支上神 → 自信课
//	中传：初传所刑的地支；末传：中传所刑的地支（遇自刑取支上神/冲位为救 → 杜传）
func fuYinMethod(ctx *Context, tianpan [12]Zhi, ke [4]Ke) SanChuan {
	// 1. 优先看克
	keks := findKes(ke, ctx.Gan)
	thieves := filterByType(keks, "下贼上")
	clashes := filterByType(keks, "上克下")
	var chu Zhi
	note := ""
	switch {
	case len(thieves) >= 1:
		chu = thieves[0].Upper
		note = "伏吟重审"
	case len(clashes) >= 1:
		chu = clashes[0].Upper
		note = "伏吟元首"
	default:
		if GanYang[ctx.Gan] {
			chu = tianpan[GanJiGong[ctx.Gan]]
			note = "自任课"
		} else {
			chu = tianpan[ctx.DayZhi]
			note = "自信课"
		}
	}

	zhong := xingNext(chu, ctx, true)
	mo := xingNext(zhong, ctx, false)
	if isSelfXing(chu) || isSelfXing(zhong) {
		note = "杜传课"
	}

	return SanChuan{
		Chu:    buildChuanEntry("初传", chu),
		Zhong:  buildChuanEntry("中传", zhong),
		Mo:     buildChuanEntry("末传", mo),
		Method: "伏吟法",
		Note:   note,
	}
}

// xingNext 在伏吟中：取 z 的刑支。遇自刑：中传取支上神（或干上神），末传取中传之冲
func xingNext(z Zhi, ctx *Context, isZhong bool) Zhi {
	if isSelfXing(z) {
		// 中传救：阳日取支上神，阴日取干上神
		if isZhong {
			if GanYang[ctx.Gan] {
				return ctx.DayZhi // 自刑时"取支"
			}
			return GanJiGong[ctx.Gan]
		}
		// 末传取冲
		return Zhi((int(z) + 6) % 12)
	}
	if nx, ok := xingOf[z]; ok {
		if nx == z {
			// 兜底：理论不应到达（自刑已处理）
			return Zhi((int(z) + 6) % 12)
		}
		return nx
	}
	return z
}

// fanYinMethod 返吟：月将与占时相冲（差6）。
//
//	有克：按贼克法常规走
//	无克：
//	  取日支之驿马为初传（无依课）
//	  中传：支上神；末传：干上神
func fanYinMethod(ctx *Context, tianpan [12]Zhi, ke [4]Ke) SanChuan {
	keks := findKes(ke, ctx.Gan)
	thieves := filterByType(keks, "下贼上")
	clashes := filterByType(keks, "上克下")

	if len(thieves) == 1 {
		sc := zeiKeMethod(tianpan, ke, thieves[0], ctx.Gan, "返吟有克")
		sc.Method = "返吟法"
		return sc
	}
	if len(thieves) > 1 {
		sc := resolveFromCandidates(ctx, tianpan, ke, thieves, "下贼上")
		sc.Method = "返吟法"
		sc.Note = "返吟有克"
		return sc
	}
	if len(clashes) == 1 {
		sc := zeiKeMethod(tianpan, ke, clashes[0], ctx.Gan, "返吟有克")
		sc.Method = "返吟法"
		return sc
	}
	if len(clashes) > 1 {
		sc := resolveFromCandidates(ctx, tianpan, ke, clashes, "上克下")
		sc.Method = "返吟法"
		sc.Note = "返吟有克"
		return sc
	}
	// 无克：分两种情境
	//   1. 井栏课（《大全》卷一 p474 + 《粹言》卷一 p30）：
	//      丁/己/辛日 + 日支为丑或未时（"丑未同干"），固定取申子辰为三传
	//   2. 其它情形：无依课，初传取日支驿马、中传支上、末传干上
	if isJingLanCase(ctx) {
		return SanChuan{
			Chu:    buildChuanEntry("初传", Shen),
			Zhong:  buildChuanEntry("中传", Zi),
			Mo:     buildChuanEntry("末传", Chen),
			Method: "返吟法",
			Note:   "井栏课",
		}
	}
	chu := yiMa(ctx.DayZhi)
	zhong := tianpan[ctx.DayZhi]
	mo := tianpan[GanJiGong[ctx.Gan]]
	return SanChuan{
		Chu:    buildChuanEntry("初传", chu),
		Zhong:  buildChuanEntry("中传", zhong),
		Mo:     buildChuanEntry("末传", mo),
		Method: "返吟法",
		Note:   "无依课",
	}
}

// isJingLanCase 是否为井栏课触发条件
//
// 《大全》卷一 p474："返吟有克亦为用，无克别有井栏名，丑未同干，已辛日登明亥…"
// 综合《粹言》卷一 p30 解读：丁/己/辛日 + 日支丑或未时，无返吟之克 → 井栏课
func isJingLanCase(ctx *Context) bool {
	dayInJingLanGan := ctx.Gan == Ding || ctx.Gan == Ji || ctx.Gan == Xin
	dayInChouWei := ctx.DayZhi == Chou || ctx.DayZhi == Wei
	return dayInJingLanGan && dayInChouWei
}

// yiMa 驿马：申子辰马在寅，寅午戌马在申，巳酉丑马在亥，亥卯未马在巳
func yiMa(z Zhi) Zhi {
	switch z {
	case Shen, Zi, Chen:
		return Yin
	case Yin, Wu, Xu:
		return Shen
	case Si, You, Chou:
		return Hai
	case Hai, Mao, Wei:
		return Si
	}
	return z
}

// isBaZhuanDay 是否为八专日：甲寅、丁未、己未、庚申、辛酉、癸丑六日。
//
// 前五日寄宫=日支；辛酉日寄宫戌、日支酉，干支同属金行紧邻，
// 古法仍归八专（《六壬大全》卷一"刚三柔六共九课"中辛酉为柔日之一）。
func isBaZhuanDay(g Gan, dz Zhi) bool {
	switch {
	case g == Jia && dz == Yin:
		return true
	case g == Ding && dz == Wei:
		return true
	case g == Ji && dz == Wei:
		return true
	case g == Geng && dz == Shen:
		return true
	case g == Xin && dz == You:
		return true
	case g == Gui && dz == Chou:
		return true
	}
	return false
}

// baZhuanMethod 八专：日干寄宫=日支（六个日子：甲寅、丁未、己未、庚申、辛酉、癸丑）。
//
// 据《六壬大全》卷一"刚三柔六共九课"、《六壬粹言》卷一 p.29
// "刚日从干阳遁顺数三神，柔日从支阴遁逆数三神"。
//
//	有克：按贼克法
//	无克：
//	  阳日（刚日）：从干上阳神（一课上神）顺数三位（含起点 +2）为初传
//	  阴日（柔日）：从第四课上神（=tianpan[tianpan[日支]]）逆数三位（-2）为初传
//	  中、末传：皆取干上神（一课上神）
func baZhuanMethod(ctx *Context, tianpan [12]Zhi, ke [4]Ke) SanChuan {
	keks := findKes(ke, ctx.Gan)
	if len(keks) == 1 {
		sc := zeiKeMethod(tianpan, ke, keks[0], ctx.Gan, "八专有克")
		sc.Method = "八专法"
		return sc
	}
	if len(keks) > 1 {
		// 八专日干支同位，四课实占两位（一二课、三四课各重叠），
		// 多克极罕见；此处沿用一般比用/涉害择克为简化处理，结果仍取「有克即用」之上神。
		sc := resolveFromCandidates(ctx, tianpan, ke, keks, "")
		sc.Method = "八专法"
		sc.Note = "八专有克"
		return sc
	}

	ganUpper := tianpan[GanJiGong[ctx.Gan]]
	var chu Zhi
	if GanYang[ctx.Gan] {
		// 阳日：从干上阳神顺数三位（含起点）= +2
		chu = Zhi((int(ganUpper) + 2) % 12)
	} else {
		// 阴日：从「第四课上神」逆数三位 = -2。
		// 八专干支同位（寄宫=日支），四课塌缩成两课：
		//   三课上神（支上神）= tianpan[日支]；
		//   四课上神           = tianpan[三课上神] = tianpan[tianpan[日支]]（故为两跳）。
		// 《大六壬指南》《大全》：「柔日从第四课上神逆数三神为发用」，取四课上神。
		zhiYinShen := tianpan[tianpan[ctx.DayZhi]]
		chu = Zhi((int(zhiYinShen) - 2 + 12) % 12)
	}
	return SanChuan{
		Chu:    buildChuanEntry("初传", chu),
		Zhong:  buildChuanEntry("中传", ganUpper),
		Mo:     buildChuanEntry("末传", ganUpper),
		Method: "八专法",
		Note:   "八专课",
	}
}
