package liuren

// zeiKeMethod 贼克法：单一克贼，取该课上神为初传
func zeiKeMethod(tianpan [12]Zhi, _ [4]Ke, picked Ke, _ Gan, note string) SanChuan {
	chu, zhong, mo := chainFromInitial(tianpan, picked.Upper)
	return SanChuan{
		Chu: chu, Zhong: zhong, Mo: mo,
		Method: "贼克法",
		Note:   note,
	}
}

// resolveFromCandidates 多克并存的情况：先用比用法；若不决再用涉害法
func resolveFromCandidates(ctx *Context, tianpan [12]Zhi, ke [4]Ke, kes []Ke, relation string) SanChuan {
	// 比用法：与日干阴阳相同的上神为用
	dayYang := GanYang[ctx.Gan]
	var matching []Ke
	for _, k := range kes {
		if ZhiYang[k.Upper] == dayYang {
			matching = append(matching, k)
		}
	}
	if len(matching) == 1 {
		chu, zhong, mo := chainFromInitial(tianpan, matching[0].Upper)
		return SanChuan{
			Chu: chu, Zhong: zhong, Mo: mo,
			Method: "比用法",
			Note:   "知一课",
		}
	}
	// 比用不决（俱比或俱不比）→ 涉害法
	return sheHaiMethod(ctx, tianpan, ke, kes, relation)
}

// sheHaiCand 涉害法候选
type sheHaiCand struct {
	ke    Ke
	deep  int
	place Zhi // 该上神当前所临地盘位
}

// sheHaiMethod 涉害法：取候选上神从天盘位置回归地盘本位过程中所受克最深者。
//
// 实现策略（主流《六壬大全》派）：
//  1. 对每个候选课之上神 X，从其在天盘上的当前地盘位置 P 出发；
//  2. 沿地盘顺推至 X 本位（即地盘中 X 所在位置），数沿途每步所在天盘神对 X 的克次数；
//  3. 克数多者为涉深，取其为初传；若仍并列，按"孟仲季"优先级（孟地寅申巳亥→仲子午卯酉→季辰戌丑未）；
//  4. 仍并列则刚（阳）日取干上神为用、柔（阴）日取支上神为用（缀瑕）。
func sheHaiMethod(ctx *Context, tianpan [12]Zhi, ke [4]Ke, kes []Ke, _ string) SanChuan {
	cs := make([]sheHaiCand, 0, len(kes))
	for _, k := range kes {
		deep, place := sheHaiCount(tianpan, k)
		cs = append(cs, sheHaiCand{k, deep, place})
	}
	maxDeep := -1
	for _, c := range cs {
		if c.deep > maxDeep {
			maxDeep = c.deep
		}
	}
	var winners []sheHaiCand
	for _, c := range cs {
		if c.deep == maxDeep {
			winners = append(winners, c)
		}
	}
	note := "涉害课"
	var pick sheHaiCand
	if len(winners) == 1 {
		pick = winners[0]
	} else {
		pick = prioritizeMengZhongJi(winners, &note)
		if pick.ke.Index == 0 {
			// 仍不决：缀瑕——阳日取干上神（一课上神）、阴日取支上神（三课上神）
			fallback := ke[0]
			if !GanYang[ctx.Gan] {
				fallback = ke[2]
			}
			pick = sheHaiCand{ke: fallback, place: fallback.Lower}
			note = "缀瑕课"
		}
	}
	chu, zhong, mo := chainFromInitial(tianpan, pick.ke.Upper)
	return SanChuan{
		Chu: chu, Zhong: zhong, Mo: mo,
		Method: "涉害法",
		Note:   note,
	}
}

// sheHaiCount 从上神 X 当前所临地盘位 P，"行来本家"顺行至本家位 X，数沿途受克次数。
//
// 《大全》卷一 p473：「**涉害行来本家止，路逢多克为用，取孟深仲浅季当休**」
// 月将加时是把天盘整体逆时针偏移到地盘上的，所以上神回归本家是顺时针方向。
// 沿途每经一个地盘位，看该位的天盘神若克 X，则记一次"涉害"，累积越多越深。
func sheHaiCount(tianpan [12]Zhi, k Ke) (int, Zhi) {
	x := k.Upper // 上神
	p := k.Lower // 上神当前所临地盘位
	count := 0
	// 从 p 顺行至 x 本位（最多走 11 步；本位不计）
	for step := 1; step <= 12; step++ {
		cur := Zhi((int(p) + step) % 12)
		if cur == x {
			break
		}
		up := tianpan[cur]
		if WuXingOvercomes(ZhiWuXing[up], ZhiWuXing[x]) {
			count++
		}
	}
	return count, p
}

// prioritizeMengZhongJi 孟（寅申巳亥）>仲（子午卯酉）>季（辰戌丑未）
func prioritizeMengZhongJi(winners []sheHaiCand, note *string) sheHaiCand {
	score := func(z Zhi) int {
		switch z {
		case Yin, Shen, Si, Hai:
			return 3
		case Zi, Wu, Mao, You:
			return 2
		default:
			return 1
		}
	}
	best := winners[0]
	bestScore := score(best.ke.Lower)
	same := 1
	for i := 1; i < len(winners); i++ {
		s := score(winners[i].ke.Lower)
		if s > bestScore {
			best = winners[i]
			bestScore = s
			same = 1
		} else if s == bestScore {
			same++
		}
	}
	if same == 1 {
		if bestScore == 3 {
			*note = "见机课"
		} else if bestScore == 2 {
			*note = "察微课"
		}
		return best
	}
	return sheHaiCand{}
}
