package liuren

// PlaceTianJiang 十二天将排布
//
//	规则（《六壬大全》《大六壬指南》《六壬粹言》主流派）：
//	  1. 由日干+昼夜查贵人应临的地支（"甲戊庚牛羊"口诀）
//	  2. 找贵人应临支在天盘中的位置，即为贵人所乘的地盘位
//	  3. 贵人加于"卯辰巳午未申"（日昼弧）则顺布；
//	     加于"酉戌亥子丑寅"（夜弧）则逆布
//	  4. 从贵人位起，按「贵腾朱六勾青空白常玄阴后」次第布 12 天将
//
// 返回数组：索引=地盘位，值=该位上的天将
func PlaceTianJiang(ctx *Context, tianpan [12]Zhi) [12]TianJiang {
	guiZhi := GuiRenByGan[ctx.Gan][0]
	if !ctx.ZhouYe {
		guiZhi = GuiRenByGan[ctx.Gan][1]
	}
	// 找贵人应临支在天盘中的位置
	guiPos := 0
	for i, up := range tianpan {
		if up == guiZhi {
			guiPos = i
			break
		}
	}
	clockwise := isClockwise(Zhi(guiPos))
	var tj [12]TianJiang
	for i := 0; i < 12; i++ {
		pos := (guiPos + i) % 12
		if !clockwise {
			pos = (guiPos - i + 12) % 12
		}
		tj[pos] = TianJiang(i)
	}
	return tj
}

// isClockwise 贵人临地盘卯辰巳午未申（日昼弧）→ 顺布 true；
// 临酉戌亥子丑寅（夜弧）→ 逆布 false。
//
// 出自《六壬大全》卷一神图："贵人加于卯辰巳午未申之上者，顺行；
// 加于酉戌亥子丑寅之上者，逆行"。
func isClockwise(diPos Zhi) bool {
	switch diPos {
	case Mao, Chen, Si, Wu, Wei, Shen:
		return true
	default:
		return false
	}
}

// TianJiangOf 返回某天盘地支所乘的天将
//
// 实现：先查该地支在天盘中的地盘位，再读 tj 数组。
func TianJiangOf(tianpan [12]Zhi, tj [12]TianJiang, upperZhi Zhi) TianJiang {
	for i, up := range tianpan {
		if up == upperZhi {
			return tj[i]
		}
	}
	return TJGuiRen
}

// TianJiangTaboo 单条乘临禁忌
type TianJiangTaboo struct {
	TianJiang TianJiang // 哪个天将
	DiZhi     Zhi       // 临到了哪个地盘位
	Note      string    // 古法说明
}

// CheckTianJiangTaboos 按《大全》卷一 p498「**十二天将惟贵神天空不乘辰戌，玄武六合不乘丑未，其餘無神不乘**」
// 检查本盘是否触犯禁忌。返回所有命中的禁忌列表（空切片表示无禁忌）。
//
// 不修改盘面，仅作信号传递给输出层和提示词。
func CheckTianJiangTaboos(tj [12]TianJiang) []TianJiangTaboo {
	var out []TianJiangTaboo
	for diPos, t := range tj {
		dz := Zhi(diPos)
		switch t {
		case TJGuiRen:
			if dz == Chen || dz == Xu {
				out = append(out, TianJiangTaboo{t, dz, "贵神不乘辰戌（地狱、天牢之位），主贵人失位、扶助减力"})
			}
		case TJTianKong:
			if dz == Chen || dz == Xu {
				out = append(out, TianJiangTaboo{t, dz, "天空不乘辰戌，主奴婢小人之诈伪难辨"})
			}
		case TJXuanWu:
			if dz == Chou || dz == Wei {
				out = append(out, TianJiangTaboo{t, dz, "玄武不乘丑未，主盗贼藏匿之机失常"})
			}
		case TJLiuHe:
			if dz == Chou || dz == Wei {
				out = append(out, TianJiangTaboo{t, dz, "六合不乘丑未，主婚姻交易之合见疏"})
			}
		}
	}
	return out
}
