package liuren

// KeTiTags 在九宗门课体之外，进一步识别盘面具备的若干传统课格标签。
//
// 说明：一副盘面可能同时具备多个标签（例如元首课 + 三阳课 + 连珠课）。
// 本函数返回所有命中的标签，按重要度从高到低排列。
func KeTiTags(pan *Pan) []KeTi {
	var tags []KeTi
	c := pan.SanChuan.Chu.Zhi
	z := pan.SanChuan.Zhong.Zhi
	m := pan.SanChuan.Mo.Zhi

	// ===== 按三传阴阳构成 =====
	yangCnt := 0
	for _, x := range []Zhi{c, z, m} {
		if ZhiYang[x] {
			yangCnt++
		}
	}
	switch yangCnt {
	case 3:
		tags = append(tags, KeTi{Name: "三阳课", Summary: "三传皆阳，事主进取、显达、外露。"})
	case 0:
		tags = append(tags, KeTi{Name: "三阴课", Summary: "三传皆阴，事主潜藏、退守、阴私。"})
	}

	// ===== 三光 / 三明：三传所乘天将吉凶 =====
	// 简化判定：三传所乘皆为吉将（贵、龙、合、常、阴、后）→ 三光课
	goodTJ := map[TianJiang]bool{TJGuiRen: true, TJQingLong: true, TJLiuHe: true, TJTaiChang: true, TJTaiYin: true, TJTianHou: true}
	evilTJ := map[TianJiang]bool{TJTengShe: true, TJZhuQue: true, TJGouChen: true, TJTianKong: true, TJBaiHu: true, TJXuanWu: true}
	allGood, allEvil := true, true
	for _, ce := range []ChuanEntry{pan.SanChuan.Chu, pan.SanChuan.Zhong, pan.SanChuan.Mo} {
		if !goodTJ[ce.TianJiang] {
			allGood = false
		}
		if !evilTJ[ce.TianJiang] {
			allEvil = false
		}
	}
	if allGood {
		tags = append(tags, KeTi{Name: "三光课", Summary: "三传皆乘吉将，光明显达，喜事当前。"})
	}
	if allEvil {
		tags = append(tags, KeTi{Name: "昏暗课", Summary: "三传皆乘凶将，晦暗阻塞，行事多阻。"})
	}

	// ===== 连珠 / 间传 / 递生 / 递克（按三传地支数列）=====
	diffCZ := (int(z) - int(c) + 12) % 12
	diffZM := (int(m) - int(z) + 12) % 12
	if diffCZ == 1 && diffZM == 1 {
		tags = append(tags, KeTi{Name: "连珠课", Summary: "三传顺连如珠，节节相承，事顺而速。"})
	} else if diffCZ == 11 && diffZM == 11 {
		tags = append(tags, KeTi{Name: "逆连课", Summary: "三传逆连，势与时反，进反退。"})
	} else if diffCZ == diffZM && diffCZ != 0 {
		tags = append(tags, KeTi{Name: "间传课", Summary: "三传等距相间，事有节奏、阶段分明。"})
	}

	// ===== 递生 / 递克（按五行）=====
	wC := ZhiWuXing[c]
	wZ := ZhiWuXing[z]
	wM := ZhiWuXing[m]
	if WuXingGenerates(wC, wZ) && WuXingGenerates(wZ, wM) {
		tags = append(tags, KeTi{Name: "递生课", Summary: "三传相生而进，气脉顺承，事易成。"})
	}
	if WuXingOvercomes(wC, wZ) && WuXingOvercomes(wZ, wM) {
		tags = append(tags, KeTi{Name: "递克课", Summary: "三传相克而进，势如逼迫，慎之。"})
	}

	// ===== 三合局 =====
	if isSanHeTrio(c, z, m) {
		tags = append(tags, KeTi{Name: "三合课", Summary: "三传成三合局，合力同心，事主结聚。"})
	}

	// ===== 六仪：三传有两传得日干阴阳贵人所乘位 =====
	// 简化：贵人所乘初传或末传 → 六仪课
	if pan.SanChuan.Chu.TianJiang == TJGuiRen || pan.SanChuan.Mo.TianJiang == TJGuiRen {
		tags = append(tags, KeTi{Name: "六仪课", Summary: "贵人乘于初/末传，事得贵助、威仪俱足。"})
	}

	// ===== 龙德：青龙临日干或日支上神 =====
	ganUpper := pan.TianPan[GanJiGong[pan.Ctx.Gan]]
	zhiUpper := pan.TianPan[pan.Ctx.DayZhi]
	if TianJiangOf(pan.TianPan, pan.TianJiang, ganUpper) == TJQingLong ||
		TianJiangOf(pan.TianPan, pan.TianJiang, zhiUpper) == TJQingLong {
		tags = append(tags, KeTi{Name: "龙德课", Summary: "青龙临日干或日支上神，主财喜临门、吉庆之象。"})
	}

	// ===== 玄胎：三传含胎养之气（用五行长生/胎 12 宫位简化）=====
	// 胎位：寅胎于子、午胎于卯、申胎于午、子胎于酉（阳干胎方）——简化：三传含子午卯酉中的两个及以上
	centerCount := 0
	for _, x := range []Zhi{c, z, m} {
		if x == Zi || x == Wu || x == Mao || x == You {
			centerCount++
		}
	}
	if centerCount >= 2 {
		tags = append(tags, KeTi{Name: "玄胎课", Summary: "三传多居四仲（子午卯酉），孕妊胎养之象，机未发。"})
	}

	// ===== 空亡相关 =====
	kongCnt := 0
	for _, ce := range []ChuanEntry{pan.SanChuan.Chu, pan.SanChuan.Zhong, pan.SanChuan.Mo} {
		if ce.IsKong {
			kongCnt++
		}
	}
	if kongCnt == 3 {
		tags = append(tags, KeTi{Name: "空亡课", Summary: "三传皆空，事终成空，宜息心勿动。"})
	} else if kongCnt == 2 {
		tags = append(tags, KeTi{Name: "两空课", Summary: "三传两空，事多不实，须择实者而行。"})
	}

	return tags
}

// isSanHeTrio 三支是否构成三合局
func isSanHeTrio(a, b, c Zhi) bool {
	groups := [][3]Zhi{
		{Shen, Zi, Chen}, {Hai, Mao, Wei},
		{Yin, Wu, Xu}, {Si, You, Chou},
	}
	for _, g := range groups {
		has := map[Zhi]bool{a: true, b: true, c: true}
		match := 0
		for _, z := range g {
			if has[z] {
				match++
			}
		}
		if match == 3 {
			return true
		}
	}
	return false
}
