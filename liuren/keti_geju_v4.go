package liuren

// 《六壬大全》卷七-卷十 课目录补遗（v4 改进点）
//
// v2/v3 已实现 19 个课格（稼穑/从革/润下/炎上/曲直/绝嗣/铸印/三奇 +
// 引従/亨通/繁昌/荣华/合欢/盘珠 + 赘塔/凌犯/淫泆/龙虎交战/蕪淫）。
// v4 据《大全》卷七-卷十「凶课目录」补 8 个高频凶课：
//
//   1. 四绝课：节气前一日，五行绝气发用 ⭐⭐⭐
//   2. 鬼墓课：干上发用 + 三传带墓神 ⭐⭐⭐
//   3. 九醜课：戊子戊午戊辰己卯己酉乙卯乙酉辛卯辛酉九日，主不正之事 ⭐⭐⭐
//   4. 白虎入丧车：白虎乘驿马入用，主奔丧、暴病 ⭐⭐⭐
//   5. 天狱课：日干上乘官鬼临年命，主刑禁牢狱
//   6. 天网课：年命落空亡之处，又遇凶神成局
//   7. 励德课：日干被三传交克，励志反省之象
//   8. 无禄课：日上空亡 + 旬空发用，主财官两虚

// KeTiGejuV4 卷七-卷十 凶课补遗 8 格
func KeTiGejuV4(pan *Pan) []KeTi {
	var out []KeTi
	c := pan.SanChuan.Chu.Zhi
	z := pan.SanChuan.Zhong.Zhi
	m := pan.SanChuan.Mo.Zhi
	gan := pan.Ctx.Gan
	dayZhi := pan.Ctx.DayZhi
	gz := GanJiGong[gan]

	// 1. 四绝课：节气前一日（立春前木绝、立夏前火绝、立秋前金绝、立冬前水绝）
	//    简化判定：日干五行 = 月支前一节气所属五行的"绝位"，且发用为该绝气支
	//    此处采用通行口诀：
	//      春金绝（春月金绝） → 月支寅卯辰、初传 = 庚辛之绝位（寅）
	//      夏水绝 → 月支巳午未、初传 = 壬癸之绝位（巳）
	//      秋木绝 → 月支申酉戌、初传 = 甲乙之绝位（申）
	//      冬火绝 → 月支亥子丑、初传 = 丙丁之绝位（亥）
	//    《大全》卷七 p665「四绝四离课」
	{
		monthZhi := monthZhiOf(pan.Ctx)
		isSiJue := false
		var jueDesc string
		switch {
		case (monthZhi == Yin || monthZhi == Mao || monthZhi == Chen) && c == Yin:
			// 春月金绝（庚绝寅）
			isSiJue = true
			jueDesc = "春月金绝（庚绝于寅）"
		case (monthZhi == Si || monthZhi == Wu || monthZhi == Wei) && c == Si:
			// 夏月水绝（壬绝巳）
			isSiJue = true
			jueDesc = "夏月水绝（壬绝于巳）"
		case (monthZhi == Shen || monthZhi == You || monthZhi == Xu) && c == Shen:
			// 秋月木绝（甲绝申）
			isSiJue = true
			jueDesc = "秋月木绝（甲绝于申）"
		case (monthZhi == Hai || monthZhi == Zi || monthZhi == Chou) && c == Hai:
			// 冬月火绝（丙绝亥）
			isSiJue = true
			jueDesc = "冬月火绝（丙绝于亥）"
		}
		if isSiJue {
			out = append(out, KeTi{
				Name:    "四绝课",
				Summary: "" + jueDesc + "，五行临绝发用，万事终止；占行止主中道而废、占病主气脱难救。",
			})
		}
	}

	// 2. 鬼墓课：日干上神为官鬼且为日干墓神；或末传乘墓且为官鬼
	//    《大全》卷七 p669「鬼墓课」
	{
		ganUp := pan.TianPan[gz]
		muZhi := GanMuZhi[gan]
		isGuanGui := func(z Zhi) bool {
			return LiuQinOfZhiByGan(z, gan) == LQGuanGui
		}
		// 干上 = 墓 且为官鬼
		hit := false
		if ganUp == muZhi && isGuanGui(ganUp) {
			hit = true
		}
		// 或末传 = 墓 且为官鬼
		if m == muZhi && isGuanGui(m) {
			hit = true
		}
		if hit {
			out = append(out, KeTi{
				Name:    "鬼墓课",
				Summary: "官鬼乘墓临干或末传，鬼气入墓困身；占病凶险、占讼有刑、占求财财失。",
			})
		}
	}

	// 3. 九醜课：戊子/戊午/戊辰/己卯/己酉/乙卯/乙酉/辛卯/辛酉九日
	//    《大全》卷七 p670「九醜课」
	{
		type gz2 struct {
			g Gan
			z Zhi
		}
		jiuChouDays := map[gz2]bool{
			{Wu1, Zi}: true, {Wu1, Wu}: true, {Wu1, Chen}: true,
			{Ji, Mao}: true, {Ji, You}: true,
			{Yi, Mao}: true, {Yi, You}: true,
			{Xin, Mao}: true, {Xin, You}: true,
		}
		if jiuChouDays[gz2{gan, dayZhi}] {
			out = append(out, KeTi{
				Name:    "九醜课",
				Summary: "九醜日起课，神煞会聚不正之处；占婚不正、占官失体、占事多反复。",
			})
		}
	}

	// 4. 白虎入丧车：白虎乘驿马 + 临年命/三传 + 庚日加重
	//    《大全》卷八 p685「白虎入丧车格」
	{
		yiMaZhi := yiMa(dayZhi)
		// 找白虎所乘地支
		var baiHuOn Zhi = -1
		for i := Zi; i <= Hai; i++ {
			if pan.TianJiang[i] == TJBaiHu {
				baiHuOn = i
				break
			}
		}
		if baiHuOn >= 0 && baiHuOn == yiMaZhi {
			// 入三传或临年命
			inSan := baiHuOn == c || baiHuOn == z || baiHuOn == m
			inNianMing := false
			if pan.NianMing != nil {
				if pan.NianMing.BenMing != nil && baiHuOn == pan.NianMing.BenMing.Zhi {
					inNianMing = true
				}
				if pan.NianMing.XingNian != nil && baiHuOn == pan.NianMing.XingNian.Zhi {
					inNianMing = true
				}
			}
			if inSan || inNianMing {
				severity := ""
				if gan == Geng {
					severity = "（庚日尤凶）"
				}
				out = append(out, KeTi{
					Name:    "白虎入丧车格",
					Summary: "白虎乘驿马入用临命" + severity + "，急变奔丧、暴病横事；占行止主死生关头、占病主危。",
				})
			}
		}
	}

	// 5. 天狱课：日干上神乘官鬼 + 该上神临年命/支
	//    《大全》卷八 p688「天狱课」
	{
		ganUp := pan.TianPan[gz]
		if LiuQinOfZhiByGan(ganUp, gan) == LQGuanGui {
			critical := false
			// 临年命
			if pan.NianMing != nil {
				if pan.NianMing.BenMing != nil && ganUp == pan.NianMing.BenMing.Zhi {
					critical = true
				}
				if pan.NianMing.XingNian != nil && ganUp == pan.NianMing.XingNian.Zhi {
					critical = true
				}
			}
			// 或乘白虎/勾陈
			tjOn := pan.TianJiang[gz]
			if tjOn == TJBaiHu || tjOn == TJGouChen {
				critical = true
			}
			if critical {
				out = append(out, KeTi{
					Name:    "天狱课",
					Summary: "官鬼临干 + 凶将（虎/勾）或临命，刑禁牢狱之象；占讼凶、占病主缠绵、占行止勿动。",
				})
			}
		}
	}

	// 6. 天网课：年命落空亡 + 三传遇凶将（虎/蛇/勾）
	//    《大全》卷八 p691「天网课」
	{
		if pan.NianMing != nil {
			pair := pan.Ctx.XunKongPair()
			isKong := func(zh Zhi) bool {
				return zh == pair[0] || zh == pair[1]
			}
			nmKong := false
			if pan.NianMing.BenMing != nil && isKong(pan.NianMing.BenMing.Zhi) {
				nmKong = true
			}
			if pan.NianMing.XingNian != nil && isKong(pan.NianMing.XingNian.Zhi) {
				nmKong = true
			}
			hasXiongJiang := false
			for _, ce := range []ChuanEntry{pan.SanChuan.Chu, pan.SanChuan.Zhong, pan.SanChuan.Mo} {
				if ce.TianJiang == TJBaiHu || ce.TianJiang == TJTengShe || ce.TianJiang == TJGouChen {
					hasXiongJiang = true
					break
				}
			}
			if nmKong && hasXiongJiang {
				out = append(out, KeTi{
					Name:    "天网课",
					Summary: "年命落空 + 三传虎蛇勾陈，如鸟入网、求脱无路；占病主沉重、占讼无脱、占求事多阻。",
				})
			}
		}
	}

	// 7. 励德课：三传交克日干 + 末传不空、为官鬼
	//    《大全》卷八 p693「励德课」
	{
		wDay := GanWuXing[gan]
		clashCount := 0
		for _, zh := range []Zhi{c, z, m} {
			if WuXingOvercomes(ZhiWuXing[zh], wDay) {
				clashCount++
			}
		}
		if clashCount >= 2 && LiuQinOfZhiByGan(m, gan) == LQGuanGui && !pan.SanChuan.Mo.IsKong {
			out = append(out, KeTi{
				Name:    "励德课",
				Summary: "三传交克日干 + 末传官鬼实位，主磨难砥砺、励志方能脱困；占事必经苦劳、不可侥幸。",
			})
		}
	}

	// 8. 无禄课：日上神 = 旬空 + 干禄落空亡或克破
	//    《大全》卷十 p738「无禄课」
	{
		ganUp := pan.TianPan[gz]
		pair := pan.Ctx.XunKongPair()
		isKong := func(zh Zhi) bool { return zh == pair[0] || zh == pair[1] }
		luZhi := GanLuZhi[gan]
		ganUpKong := isKong(ganUp)
		luKong := isKong(luZhi)
		// 干禄被破：禄位在天盘上的神被克
		luOnTian := pan.TianPan[luZhi]
		luKe := WuXingOvercomes(ZhiWuXing[luOnTian], ZhiWuXing[luZhi])
		if ganUpKong && (luKong || luKe) {
			out = append(out, KeTi{
				Name:    "无禄课",
				Summary: "日上空亡 + 禄位破空，财官两虚、谋望不成；占求财求官皆失、占事虚惊。",
			})
		}
	}

	return out
}
