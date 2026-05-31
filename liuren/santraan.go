package liuren

// DeriveSanChuan 从四课发三传。按九宗门次第判定（伏吟/返吟优先，八专次之，后依有克/无克走贼克-比用-涉害-遥克-昴星-别责）。
func DeriveSanChuan(ctx *Context, tianpan [12]Zhi, ke [4]Ke) SanChuan {
	// 1. 伏吟：月将=占时，天盘各居本位
	if ctx.YueJiang == ctx.ZhanShi {
		return fuYinMethod(ctx, tianpan, ke)
	}
	// 2. 返吟：月将与占时相冲（差 6 位）
	if (int(ctx.YueJiang)-int(ctx.ZhanShi)+12)%12 == 6 {
		return fanYinMethod(ctx, tianpan, ke)
	}
	// 3. 八专：干支同行的六个日柱（《六壬大全》卷一）
	//    甲寅、丁未、己未、庚申、癸丑——寄宫=日支
	//    辛酉——寄宫戌、日支酉，干支同属金行紧邻，亦入八专
	if isBaZhuanDay(ctx.Gan, ctx.DayZhi) {
		return baZhuanMethod(ctx, tianpan, ke)
	}

	// 4. 四课中有克（上克下或下贼上）
	keks := findKes(ke, ctx.Gan)
	if len(keks) > 0 {
		// 有下贼上者优先取下贼上
		thieves := filterByType(keks, "下贼上")
		if len(thieves) == 1 {
			return zeiKeMethod(tianpan, ke, thieves[0], ctx.Gan, "重审课")
		}
		if len(thieves) > 1 {
			return resolveFromCandidates(ctx, tianpan, ke, thieves, "下贼上")
		}
		// 无下贼上，仅有上克下
		clashes := filterByType(keks, "上克下")
		if len(clashes) == 1 {
			return zeiKeMethod(tianpan, ke, clashes[0], ctx.Gan, "元首课")
		}
		return resolveFromCandidates(ctx, tianpan, ke, clashes, "上克下")
	}

	// 5. 无克：遥克
	if sc, ok := yaoKeMethod(ctx, tianpan, ke); ok {
		return sc
	}

	// 6. 无克无遥：昴星
	// 判别课数是否齐备：若四课上神去重后只有 3 个则为别责；齐 4 个则为昴星
	if isSanKeNotFour(ke) {
		return bieZeMethod(ctx, tianpan, ke)
	}
	return maoXingMethod(ctx, tianpan, ke)
}

// findKes 返回四课中存在克/贼关系的课
func findKes(ke [4]Ke, g Gan) []Ke {
	var out []Ke
	for _, k := range ke {
		rel := relationUpperLowerForKe(k, g)
		if rel == "上克下" || rel == "下贼上" {
			k.Relation = rel
			out = append(out, k)
		}
	}
	return out
}

// relationUpperLowerForKe 一课下神用日干五行，其余用地支五行
func relationUpperLowerForKe(k Ke, g Gan) string {
	if k.Index == 1 {
		return relationOfFirstKe(g, k.Upper)
	}
	return RelationOfZhi(k.Upper, k.Lower)
}

func filterByType(kes []Ke, t string) []Ke {
	var out []Ke
	for _, k := range kes {
		if k.Relation == t {
			out = append(out, k)
		}
	}
	return out
}

// isSanKeNotFour 四课上神去重后恰为 3 → 课不全（别责课）
//
// 古法"别责课"严格触发条件：四课中恰有两课的上神相同（通常为二与一、或四与三），
// 即实际上只有三个不同上神；且无克贼、无遥克。
// 去重数若 < 3（两组重叠，仅余两个或一个不同上神）属八专/伏吟等已被前置分支拦截的特殊盘，
// 不应归别责；故此处严格判 == 3。
func isSanKeNotFour(ke [4]Ke) bool {
	seen := map[Zhi]bool{}
	for _, k := range ke {
		seen[k.Upper] = true
	}
	return len(seen) == 3
}

// chainNextUpper 由一个天神 z 找其在地盘对应位置的"上神"（即 tianpan[z]）
func chainNextUpper(tianpan [12]Zhi, z Zhi) Zhi {
	return tianpan[z]
}

// buildChuanEntry 构造三传中的一传（填充天将/六亲/空亡的职责留给最终渲染阶段）
func buildChuanEntry(name string, z Zhi) ChuanEntry {
	return ChuanEntry{Name: name, Zhi: z}
}

// chainFromInitial 由初传神沿天盘递推中传、末传
func chainFromInitial(tianpan [12]Zhi, chu Zhi) (ChuanEntry, ChuanEntry, ChuanEntry) {
	zhong := chainNextUpper(tianpan, chu)
	mo := chainNextUpper(tianpan, zhong)
	return buildChuanEntry("初传", chu),
		buildChuanEntry("中传", zhong),
		buildChuanEntry("末传", mo)
}
