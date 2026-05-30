package liuren

// yaoKeMethod 遥克法：四课上神与日干之间的克（非四课内部克）
//
// 《六壬大全》卷一入式法："神遥赴日曰蒿矢，日遥赴神曰弹射"
// 《大六壬指南》卷一心印赋："遥克日曰蒿矢课，以彼能遥伤於我而似矢...
//
//	若遥克彼者取以为用名曰弹射课"
//
// 故定名：
//
//	四课上神遥克日干  → 蒿矢课（彼来克我，似矢射来；优先取用）
//	日干遥克四课上神  → 弹射课（我去克彼）
//
// 多个候选时与日干阴阳同者优先。
func yaoKeMethod(ctx *Context, tianpan [12]Zhi, ke [4]Ke) (SanChuan, bool) {
	wDay := GanWuXing[ctx.Gan]
	dayYang := GanYang[ctx.Gan]

	type cand struct {
		ke         Ke
		shenKeDay  bool // true: 上神克日干（蒿矢） false: 日干克上神（弹射）
	}
	var cs []cand
	for _, k := range ke {
		wu := ZhiWuXing[k.Upper]
		if WuXingOvercomes(wu, wDay) {
			cs = append(cs, cand{k, true}) // 上神克日 → 蒿矢
		} else if WuXingOvercomes(wDay, wu) {
			cs = append(cs, cand{k, false}) // 日克上神 → 弹射
		}
	}
	if len(cs) == 0 {
		return SanChuan{}, false
	}
	// 优先取"上神克日"（蒿矢）——彼来克我更紧，主事；若无再取"日克上神"（弹射）
	var pool []cand
	for _, c := range cs {
		if c.shenKeDay {
			pool = append(pool, c)
		}
	}
	note := "蒿矢课"
	if len(pool) == 0 {
		pool = cs
		note = "弹射课"
	}
	// 多者取与日干阴阳同者
	var sameYin []cand
	for _, c := range pool {
		if ZhiYang[c.ke.Upper] == dayYang {
			sameYin = append(sameYin, c)
		}
	}
	pick := pool[0]
	if len(sameYin) >= 1 {
		pick = sameYin[0]
	}
	chu, zhong, mo := chainFromInitial(tianpan, pick.ke.Upper)
	return SanChuan{
		Chu: chu, Zhong: zhong, Mo: mo,
		Method: "遥克法",
		Note:   note,
	}, true
}

// maoXingMethod 昴星法：无克无遥时以地盘酉（昴）为枢
//
// 《六壬大全》卷一入式法："五昴星法。无遥无昴星霜阳仰俯阴位中传初也。
//
//	刚日先辰而后日，柔日先日而后辰"
//
// 即：刚（阳）日中传取辰上（支上神）、末传取日上（干上神）；
//     柔（阴）日中传取日上（干上神）、末传取辰上（支上神）。
//
// 初传按"虎视/冬蛇"分别从地盘酉位取：
//
//	阳日（虎视转蓬课）：取天盘酉所临的地盘位为初传；中=支上、末=干上。
//	阴日（冬蛇掩目课）：取地盘酉位上方的天盘神为初传；中=干上、末=支上。
func maoXingMethod(ctx *Context, tianpan [12]Zhi, _ [4]Ke) SanChuan {
	var chu, zhong, mo Zhi
	var note string
	ganZhi := GanJiGong[ctx.Gan]
	ganUpper := tianpan[ganZhi]
	dayUpper := tianpan[ctx.DayZhi]

	if GanYang[ctx.Gan] {
		// 阳日（刚日）：先辰后日 → 中=支上、末=干上
		for i, up := range tianpan {
			if up == You {
				chu = Zhi(i)
				break
			}
		}
		zhong = dayUpper // 支上神
		mo = ganUpper    // 干上神
		note = "虎视转蓬课"
	} else {
		// 阴日（柔日）：先日后辰 → 中=干上、末=支上
		chu = tianpan[You]
		zhong = ganUpper // 干上神
		mo = dayUpper    // 支上神
		note = "冬蛇掩目课"
	}
	return SanChuan{
		Chu:    buildChuanEntry("初传", chu),
		Zhong:  buildChuanEntry("中传", zhong),
		Mo:     buildChuanEntry("末传", mo),
		Method: "昴星法",
		Note:   note,
	}
}

// bieZeMethod 别责法：四课不全（上神只有三个），无克无遥无昴之对称
//
// 据《六壬大全》卷一入式法："刚日干合上头神，柔日支前三合取"；
// 《六壬指南》"甲己庚乙丙辛丁壬戊癸六合也"；
// 《六壬粹言》卷一"刚日别贵课，取干合上神为用；柔日别责课，取支三合前位之神为用"。
//
//	阳日（刚日别贵）：取日干六合（如甲↔己、丙↔辛...）的阴干寄宫上之天盘神为初传
//	阴日（柔日别责）：取日支三合前一位地支在天盘上对应的神为初传
//	中、末传：皆取干上神（日干寄宫上神）
func bieZeMethod(ctx *Context, tianpan [12]Zhi, _ [4]Ke) SanChuan {
	ganUpper := tianpan[GanJiGong[ctx.Gan]]
	var chu Zhi
	if GanYang[ctx.Gan] {
		// 阳日：日干六合 → 该阴干寄宫 → 天盘上之神
		heGan := GanLiuhe[ctx.Gan]
		chu = tianpan[GanJiGong[heGan]]
	} else {
		// 阴日：日支三合之前一位 → 天盘上之神
		chu = tianpan[sanHeNext(ctx.DayZhi)]
	}
	return SanChuan{
		Chu:    buildChuanEntry("初传", chu),
		Zhong:  buildChuanEntry("中传", ganUpper),
		Mo:     buildChuanEntry("末传", ganUpper),
		Method: "别责法",
		Note:   "别责课",
	}
}

// sanHeNext 三合局中 z 的前一位
//
//	申子辰（水局）、寅午戌（火局）、巳酉丑（金局）、亥卯未（木局）
//	前一位含义：沿局按顺序 A→B→C→A 的"前进"。
func sanHeNext(z Zhi) Zhi {
	groups := [4][3]Zhi{
		{Shen, Zi, Chen},
		{Yin, Wu, Xu},
		{Si, You, Chou},
		{Hai, Mao, Wei},
	}
	for _, g := range groups {
		for i, v := range g {
			if v == z {
				return g[(i+1)%3]
			}
		}
	}
	return z
}
