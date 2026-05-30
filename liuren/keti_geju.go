package liuren

// 《六壬大全》卷一 p493「课目」高频课格识别
//
// 本文件实现 8 个高频格局（按 audit 报告 P1-3 优先级）：
//   1. 稼穑课：三传辰戌丑未全（土局）
//   2. 从革课：三传巳酉丑全（金局）
//   3. 润下课：三传申子辰全（水局）
//   4. 炎上课：三传寅午戌全（火局）
//   5. 曲直课：三传亥卯未全（木局）
//   6. 绝嗣课：四课全为下贼上 → 大凶之象（《大全》卷一 p494「绝嗣四下贼上」）
//   7. 铸印课：戌加巳发用 → 主迁官改革（《大全》卷一 p494「铸印发用」）
//   8. 三奇课：三传中含天上三奇（甲戊庚 / 乙丙丁 / 辛壬癸 之一组）
//
// 命中后追加到 pan.Tags。
func KeTiGeju(pan *Pan) []KeTi {
	var out []KeTi
	c := pan.SanChuan.Chu.Zhi
	z := pan.SanChuan.Zhong.Zhi
	m := pan.SanChuan.Mo.Zhi

	// 1-5 五行成局
	if isSetEqual([]Zhi{c, z, m}, []Zhi{Chen, Xu, Chou, Wei}) {
		out = append(out, KeTi{Name: "稼穑课", Summary: "三传辰戌丑未全为土局，主厚载、积聚，事缓而稳。"})
	}
	if isSetEqual([]Zhi{c, z, m}, []Zhi{Si, You, Chou}) {
		out = append(out, KeTi{Name: "从革课", Summary: "三传巳酉丑全为金局，主刑革、肃杀，宜断不宜留。"})
	}
	if isSetEqual([]Zhi{c, z, m}, []Zhi{Shen, Zi, Chen}) {
		out = append(out, KeTi{Name: "润下课", Summary: "三传申子辰全为水局，主流通、智巧，事如水赴海。"})
	}
	if isSetEqual([]Zhi{c, z, m}, []Zhi{Yin, Wu, Xu}) {
		out = append(out, KeTi{Name: "炎上课", Summary: "三传寅午戌全为火局，主光明、显达，事疾而亨。"})
	}
	if isSetEqual([]Zhi{c, z, m}, []Zhi{Hai, Mao, Wei}) {
		out = append(out, KeTi{Name: "曲直课", Summary: "三传亥卯未全为木局，主仁爱、生发，事顺而长。"})
	}

	// 6 绝嗣课：四课均为下贼上
	all := true
	for _, k := range pan.SiKe {
		rel := relationUpperLowerForKe(k, pan.Ctx.Gan)
		if rel != "下贼上" {
			all = false
			break
		}
	}
	if all {
		out = append(out, KeTi{Name: "绝嗣课", Summary: "四课全为下贼上，纲纪倒置、大凶之象，事终不可成。"})
	}

	// 7 铸印课：戌加巳发用
	// 即天盘戌临地盘巳，且初传 = 该上神（戌）
	if c == Xu && pan.TianPan[Si] == Xu {
		out = append(out, KeTi{Name: "铸印课", Summary: "戌加巳发用，主迁官改革，宜进取功名。"})
	}

	// 8 三奇课：三传地支对应的"天上三奇"
	if isTianShangSanQi(c, z, m) {
		out = append(out, KeTi{Name: "三奇课", Summary: "三传带天上三奇，主吉庆、贵人扶持，凶事可解。"})
	}

	return out
}

// isSetEqual 三传是否构成给定五行局（局支可能 3 或 4 个，三传必须全部落在局内且不同）
func isSetEqual(three []Zhi, group []Zhi) bool {
	in := func(z Zhi) bool {
		for _, g := range group {
			if z == g {
				return true
			}
		}
		return false
	}
	for _, z := range three {
		if !in(z) {
			return false
		}
	}
	// 至少 3 个不同（土局允许 4 选 3）
	seen := map[Zhi]bool{}
	for _, z := range three {
		seen[z] = true
	}
	return len(seen) == 3
}

// isTianShangSanQi 三传是否含天上三奇之一组
//
// 《大全》卷一 p494 + 通行命理：天上三奇为「甲戊庚 / 乙丙丁 / 辛壬癸」。
// 此处以三传所乘天将的"日干轮值"间接判定有困难，简化为：
//   - 三传地支对应的"遁干"中是否含三奇组合
//
// 占课时三传地支均会被遁干，遁干表参见 ganzhi_aux.go。
func isTianShangSanQi(c, z, m Zhi) bool {
	// 简化：以三传地支本身为标识，看是否构成"三奇地支局"。
	// 经典三奇定位是干层的，地支层近似为：
	//   甲戊庚 → 寅辰申 / 寅戌申 等组合（甲寄寅、戊寄巳、庚寄申，但巳/寅常有冲突）
	//   乙丙丁 → 辰巳未（乙寄辰、丙寄巳、丁寄未）
	//   辛壬癸 → 戌亥丑（辛寄戌、壬寄亥、癸寄丑）
	groups := [][3]Zhi{
		{Yin, Si, Shen}, // 甲戊庚 寄宫
		{Chen, Si, Wei}, // 乙丙丁 寄宫
		{Xu, Hai, Chou}, // 辛壬癸 寄宫
	}
	for _, g := range groups {
		has := map[Zhi]bool{c: true, z: true, m: true}
		match := 0
		for _, x := range g {
			if has[x] {
				match++
			}
		}
		if match == 3 {
			return true
		}
	}
	return false
}

// KeTiGejuMore 卷六课经三高频课格（v3 新增 6 个）
//
//   引従课 / 亨通课 / 繁昌课 / 荣华课 / 合欢课 / 盘珠课
func KeTiGejuMore(pan *Pan) []KeTi {
	var out []KeTi
	c := pan.SanChuan.Chu.Zhi
	z := pan.SanChuan.Zhong.Zhi
	m := pan.SanChuan.Mo.Zhi
	gan := pan.Ctx.Gan
	zhi := pan.Ctx.DayZhi
	gz := GanJiGong[gan]

	// 引従课：初传居干前/支前 + 末传居干后/支后
	//   "前" = 地支序号 -1，"后" = 地支序号 +1
	{
		isPrev := func(a, b Zhi) bool { return (int(a)+1)%12 == int(b) }
		isNext := func(a, b Zhi) bool { return (int(a)+11)%12 == int(b) }
		chuQian := isPrev(c, gz) || isPrev(c, zhi)
		moHou := isNext(m, gz) || isNext(m, zhi)
		if chuQian && moHou {
			out = append(out, KeTi{
				Name:    "引従课",
				Summary: "初传居干支前位、末传居干支后位，引前従后，主升迁/迁居之吉。",
			})
		}
	}

	// 亨通课：三传与日干互生 + 三传内部相生
	//   即 三传都生日 / 日生三传 / 中末递生
	{
		wDay := GanWuXing[gan]
		wC, wZ, wM := ZhiWuXing[c], ZhiWuXing[z], ZhiWuXing[m]
		// 三传遁生日干（三传五行皆生日干）
		allShengDay := WuXingGenerates(wC, wDay) &&
			WuXingGenerates(wZ, wDay) &&
			WuXingGenerates(wM, wDay)
		// 三传递生且末传生日干
		dijiSheng := WuXingGenerates(wC, wZ) && WuXingGenerates(wZ, wM) && WuXingGenerates(wM, wDay)
		if allShengDay || dijiSheng {
			out = append(out, KeTi{
				Name:    "亨通课",
				Summary: "三传与日干互生贯气，事如百川赴海、亨通无阻；占求事顺成。",
			})
		}
	}

	// 繁昌课：旺相之德临命发用
	//   即 月德/天德/年德 与 本命/行年同位 + 该位入三传
	{
		if pan.NianMing != nil {
			ssMap := map[Zhi]string{}
			for _, s := range pan.ShenSha {
				if s.Name == "天德" || s.Name == "月德" {
					ssMap[s.Zhi] = s.Name
				}
			}
			hits := []Zhi{}
			if pan.NianMing.BenMing != nil {
				if _, ok := ssMap[pan.NianMing.BenMing.Zhi]; ok {
					hits = append(hits, pan.NianMing.BenMing.Zhi)
				}
			}
			if pan.NianMing.XingNian != nil {
				if _, ok := ssMap[pan.NianMing.XingNian.Zhi]; ok {
					hits = append(hits, pan.NianMing.XingNian.Zhi)
				}
			}
			for _, h := range hits {
				if h == c || h == z || h == m {
					out = append(out, KeTi{
						Name:    "繁昌课",
						Summary: "天德/月德临本命行年并入三传，主家道兴隆、生育繁昌、婚姻圆满。",
					})
					break
				}
			}
		}
	}

	// 荣华课：禄马贵人临干上
	//   即 干上神 = 禄神/驿马/贵人 之一
	{
		ganUp := pan.TianPan[gz]
		var hit string
		if ganUp == GanLuZhi[gan] {
			hit = "禄"
		} else if ganUp == yiMa(zhi) {
			hit = "马"
		} else {
			day := GuiRenByGan[gan][0]
			nig := GuiRenByGan[gan][1]
			if ganUp == day || ganUp == nig {
				hit = "贵"
			}
		}
		if hit != "" {
			out = append(out, KeTi{
				Name:    "荣华课",
				Summary: "禄马贵人之一临干上神（" + hit + "神临干），主荣华富贵、求名求官皆吉。",
			})
		}
	}

	// 合欢课：日干日支六合 + 干上支上亦合
	//   即 GanJiGong[gan] 与 zhi 六合 OR 干上神与支上神六合
	{
		ganUp := pan.TianPan[gz]
		zhiUp := pan.TianPan[zhi]
		if zhiLiuhe(gz, zhi) || zhiLiuhe(ganUp, zhiUp) {
			out = append(out, KeTi{
				Name:    "合欢课",
				Summary: "干支六合或干支上神六合，主婚姻喜成、合作两悦、人和事谐。",
			})
		}
	}

	// 盘珠课：三传与日时同气贯通
	//   即 三传 + 日干 + 占时 共五位中至少 3 位五行相同
	{
		ws := []WuXing{
			GanWuXing[gan],
			ZhiWuXing[zhi],
			ZhiWuXing[pan.Ctx.ZhanShi],
			ZhiWuXing[c], ZhiWuXing[z], ZhiWuXing[m],
		}
		count := map[WuXing]int{}
		for _, w := range ws {
			count[w]++
		}
		max := 0
		for _, n := range count {
			if n > max {
				max = n
			}
		}
		if max >= 5 {
			out = append(out, KeTi{
				Name:    "盘珠课",
				Summary: "日干日支占时三传同气贯通，如珠盘流转，主圆满之象、事必周全。",
			})
		}
	}

	return out
}

// KeTiGejuExtra 卷六课经三次频课格（v3 P2 新增 5 个）
//
//   赘塔 / 凌犯 / 淫泆 / 龙虎交战 / 蕪淫
func KeTiGejuExtra(pan *Pan) []KeTi {
	var out []KeTi
	c := pan.SanChuan.Chu.Zhi
	z := pan.SanChuan.Zhong.Zhi
	m := pan.SanChuan.Mo.Zhi
	gan := pan.Ctx.Gan
	zhi := pan.Ctx.DayZhi
	gz := GanJiGong[gan]

	// 赘塔课：日干就支临、静身在外（《大全》卷六 p646）
	//   定义：干临支位（即日支位上的天盘神 = 日干寄宫支）+ 干上神被克
	{
		zhiUp := pan.TianPan[zhi]
		if zhiUp == gz {
			ganUp := pan.TianPan[gz]
			// 干上神被日支或上神克
			if WuXingOvercomes(ZhiWuXing[zhi], ZhiWuXing[ganUp]) {
				out = append(out, KeTi{
					Name:    "赘塔课",
					Summary: "干就支临、静身在外，男子如赘婿入女家，求财奔他乡；占婚则女家强势。",
				})
			}
		}
	}

	// 凌犯课：上克下又下贼上（《大全》卷六 p645）
	//   定义：四课中至少 2 课为下贼上，且至少 1 课为上克下
	{
		var thieves, clashes int
		for _, k := range pan.SiKe {
			rel := relationUpperLowerForKe(k, gan)
			switch rel {
			case "下贼上":
				thieves++
			case "上克下":
				clashes++
			}
		}
		if thieves >= 2 && clashes >= 1 {
			out = append(out, KeTi{
				Name:    "凌犯课",
				Summary: "上下交克，名侵下凌犯；占讼主两造皆有理、官司难决；占争主以下犯上。",
			})
		}
	}

	// 淫泆课：六合或天后入用 + 三传带桃花/沐浴位（《大全》卷六 p648-649）
	{
		hasLiuHe := pan.SanChuan.Chu.TianJiang == TJLiuHe ||
			pan.SanChuan.Zhong.TianJiang == TJLiuHe ||
			pan.SanChuan.Mo.TianJiang == TJLiuHe
		hasTianHou := pan.SanChuan.Chu.TianJiang == TJTianHou ||
			pan.SanChuan.Zhong.TianJiang == TJTianHou ||
			pan.SanChuan.Mo.TianJiang == TJTianHou
		taoHuaPos := taoHua(zhi)
		hasTaoHua := c == taoHuaPos || z == taoHuaPos || m == taoHuaPos
		if (hasLiuHe || hasTianHou) && hasTaoHua {
			out = append(out, KeTi{
				Name:    "淫泆课",
				Summary: "六合天后入用且三传带桃花，主男女私情、阴私不正；占婚不利，占病多由色生。",
			})
		}
	}

	// 龙虎交战课：青龙白虎相加为用（《大全》卷六 p650）
	//   定义：三传中同时有青龙、白虎乘临
	{
		hasLong := pan.SanChuan.Chu.TianJiang == TJQingLong ||
			pan.SanChuan.Zhong.TianJiang == TJQingLong ||
			pan.SanChuan.Mo.TianJiang == TJQingLong
		hasHu := pan.SanChuan.Chu.TianJiang == TJBaiHu ||
			pan.SanChuan.Zhong.TianJiang == TJBaiHu ||
			pan.SanChuan.Mo.TianJiang == TJBaiHu
		if hasLong && hasHu {
			out = append(out, KeTi{
				Name:    "龙虎交战课",
				Summary: "青龙白虎同入三传，吉凶相搏；占争主两强相持、占病主吉凶难决、占讼主官事翻覆。",
			})
		}
	}

	// 蕪淫课：四课不备（上神去重 < 4）+ 三传带桃花或玄武入用（《大全》卷七 p652）
	{
		seen := map[Zhi]bool{}
		for _, k := range pan.SiKe {
			seen[k.Upper] = true
		}
		notFull := len(seen) < 4
		hasYin := pan.SanChuan.Chu.TianJiang == TJXuanWu ||
			pan.SanChuan.Zhong.TianJiang == TJXuanWu ||
			pan.SanChuan.Mo.TianJiang == TJXuanWu ||
			pan.SanChuan.Chu.TianJiang == TJTaiYin ||
			pan.SanChuan.Zhong.TianJiang == TJTaiYin ||
			pan.SanChuan.Mo.TianJiang == TJTaiYin
		if notFull && hasYin {
			out = append(out, KeTi{
				Name:    "蕪淫课",
				Summary: "四课不备且玄武太阴入用，家中阴私交杂；占家事不正、占婚多失节。",
			})
		}
	}

	return out
}
