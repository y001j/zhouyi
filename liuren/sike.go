package liuren

// BuildSiKe 起四课
//
//	一课：下 = 日干寄宫支，上 = 该位天盘神
//	二课：下 = 一课上神，上 = 该位天盘神
//	三课：下 = 日支，上 = 该位天盘神
//	四课：下 = 三课上神，上 = 该位天盘神
func BuildSiKe(ctx *Context, tianpan [12]Zhi) [4]Ke {
	var ke [4]Ke

	// 一课
	low1 := GanJiGong[ctx.Gan]
	up1 := UpperOf(tianpan, low1)
	ke[0] = Ke{Index: 1, Upper: up1, Lower: low1, Relation: relationOfFirstKe(ctx.Gan, up1)}

	// 二课
	low2 := up1
	up2 := UpperOf(tianpan, low2)
	ke[1] = Ke{Index: 2, Upper: up2, Lower: low2, Relation: RelationOfZhi(up2, low2)}

	// 三课
	low3 := ctx.DayZhi
	up3 := UpperOf(tianpan, low3)
	ke[2] = Ke{Index: 3, Upper: up3, Lower: low3, Relation: RelationOfZhi(up3, low3)}

	// 四课
	low4 := up3
	up4 := UpperOf(tianpan, low4)
	ke[3] = Ke{Index: 4, Upper: up4, Lower: low4, Relation: RelationOfZhi(up4, low4)}

	return ke
}

// relationOfFirstKe 一课下神视作日干（五行取自天干）
func relationOfFirstKe(g Gan, upper Zhi) string {
	wu := ZhiWuXing[upper]
	wl := GanWuXing[g]
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
