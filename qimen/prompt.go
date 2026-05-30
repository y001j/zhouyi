package qimen

import (
	"fmt"
	"strings"
)

// GenerateAIPrompt 为奇门遁甲盘面生成 AI 解局提示词（阶段 2 完整版）。
//
//	question     所问之事（可空）
//	focusGuide   问题类型侧重文字（由主包传入，如 周易的 FocusGuide；可空）
//	qType        问题类型键（career/wealth/...；用于生成类神直指。可空）
//
// 提示词结构：
//  1. 顶部注意事项（证据 / 不复述 / 预抽取信号可直采）
//  2. 所问之事
//  3. 起局时间 + 四柱 + 节气 + 阴阳遁 + 局数
//  4. 值符 值使
//  5. 九宫盘（3×3 markdown 表格）
//  6. 九宫逐格详情（含五行/旺衰/门迫/入墓/击刑）
//  7. 当月五行旺衰
//  8. 本局关键信号（代码预抽取）
//  9. 命中格局（按吉凶分组）
//  10. 类神直指（按问题类型）
//  11. 奇门取象原则
//  12. 解局侧重（由 focusGuide 给出，可空）
//  13. 请按以下结构解局（7 步 + 每步依据 + 信度自评）
func GenerateAIPrompt(pan *Pan, question, focusGuide, qType string) string {
	if pan == nil || pan.Ctx == nil {
		return ""
	}
	var b strings.Builder

	// ======= 顶部指令 =======
	b.WriteString("你是一位精通奇门遁甲（三式之一，兵家谋事择方之要）的术数顾问，请根据以下盘面为我判读。\n")
	b.WriteString("**重要注意**：\n")
	b.WriteString("1. 以下盘面为给你的原始数据，请直接**解读其含义**，不要大段复述盘面原文（关键字句可简短引用）。\n")
	b.WriteString("2. 每个结论请**标注依据**，形如「（依据：值符天心落乾六宫开门得令）」或「（依据：青龙返首于艮八宫）」。\n")
	b.WriteString("3. 奇门主**天时地利之争**——长于**谋大局、定方位、择时机、布奇正**，宜用于决策布局、趋吉避凶、用兵择方。请充分发挥这一所长，把「往哪个方向、在什么时辰、如何排兵布阵」讲透；至于幽微人心、细碎情态，非奇门所擅长，可略带一笔即可，避免当成六爻 / 六壬替代品去断人事细务。\n")
	b.WriteString("4. 素材中的「本局关键信号」「命中格局」「类神直指」为代码预抽取的结论，可直接采用；其余为支撑材料。\n\n")

	// ======= 所问之事 =======
	if question != "" {
		b.WriteString("## 所问之事\n" + question + "\n\n")
	}

	// ======= 起局信息 =======
	ctx := pan.Ctx
	b.WriteString("## 起局时间与四柱\n")
	b.WriteString("- 阳历：" + ctx.Time.Format("2006-01-02 15:04") + "\n")
	b.WriteString(fmt.Sprintf("- 四柱：%s %s %s %s\n", ctx.YearGZ, ctx.MonthGZ, ctx.DayGZ, ctx.HourGZ))
	b.WriteString(fmt.Sprintf("- 节气：%s\n", ctx.JieQi))
	b.WriteString(fmt.Sprintf("- 阴阳遁：%s · %s · 第%d局\n", ctx.Dun, ctx.Yuan, ctx.Ju))
	b.WriteString(fmt.Sprintf("- 时柱旬首：%s（遁干 %s）\n", ctx.Xunshou, ctx.Dungan))
	b.WriteString(fmt.Sprintf("- 旬空：%s、%s\n", ctx.XunKong[0], ctx.XunKong[1]))
	b.WriteString(fmt.Sprintf("- 时支驿马：%s\n\n", pan.YiMaZhi))

	// ======= 值符值使 =======
	b.WriteString("## 值符 · 值使\n")
	b.WriteString(fmt.Sprintf("- 值符星：**%s** 落 **%s**（跟时干，主事态之机轴）\n", pan.ZhiFuStar, pan.ZhiFuPalace))
	b.WriteString(fmt.Sprintf("- 值使门：**%s** 落 **%s**（跟时支，主人事之动向）\n\n", pan.ZhiShiGate, pan.ZhiShiPalace))

	// ======= 九宫盘 =======
	b.WriteString("## 九宫盘\n")
	b.WriteString("按传统 3×3 布局（上南下北、左东右西）：\n\n")
	b.WriteString(renderMarkdownPan(pan))
	b.WriteString("\n")

	// ======= 九宫逐格详情 =======
	b.WriteString("## 九宫逐格详情（按飞星序 坎1 → 离9）\n")
	b.WriteString("| 宫位 | 地干 | 天干 | 星(旺衰) | 门 | 神 | 门宫关系 | 标记 |\n")
	b.WriteString("|---|---|---|---|---|---|---|---|\n")
	for i := 0; i < 9; i++ {
		c := pan.Cells[i]
		heaven := dashIfEmpty(c.HeavenStem)
		door := dashIfEmpty(c.Door)
		god := dashIfEmpty(c.God)
		starCell := c.Star
		if c.StarWangShuai != "" {
			starCell = fmt.Sprintf("%s(%s)", c.Star, c.StarWangShuai)
		}
		doorRel := "—"
		if c.DoorPalaceRel != "" {
			doorRel = c.DoorPalaceRel
			if c.IsDoorPo {
				doorRel += "·**迫**"
			}
			if c.IsDoorSheng {
				doorRel += "·生"
			}
		}
		marks := []string{}
		if c.IsVoid {
			marks = append(marks, "空")
		}
		if c.IsYima {
			marks = append(marks, "马")
		}
		if c.IsTianStemMu {
			marks = append(marks, "天干墓")
		}
		if c.IsEarthStemMu {
			marks = append(marks, "地干墓")
		}
		if c.IsJiXing {
			marks = append(marks, "击刑")
		}
		mark := dashIfEmpty(strings.Join(marks, "/"))
		b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s |\n",
			c.PalaceName, c.EarthStem, heaven, starCell, door, god, doorRel, mark))
	}
	b.WriteString("\n")

	// ======= 当月五行旺衰 =======
	if pan.WangXiangXS[0] != "" {
		b.WriteString(fmt.Sprintf("## 当月五行旺衰（月支 %s）\n", ctx.MonthZhi))
		b.WriteString(fmt.Sprintf("- 旺：%s ｜ 相：%s ｜ 休：%s ｜ 囚：%s ｜ 死：%s\n\n",
			pan.WangXiangXS[0], pan.WangXiangXS[1], pan.WangXiangXS[2], pan.WangXiangXS[3], pan.WangXiangXS[4]))
	}

	// ======= 本局关键信号 =======
	if sig := ExtractQimenSignals(pan); sig != "" {
		b.WriteString("## 本局关键信号（代码预抽取，请直接采用）\n")
		b.WriteString(sig)
		b.WriteString("\n")
	}

	// ======= 命中格局 =======
	if hits := DetectPatterns(pan); len(hits) > 0 {
		b.WriteString("## 命中格局（按吉凶分组）\n")
		var ji, xiong, shen []PatternHit
		for _, h := range hits {
			switch h.Category {
			case "吉格":
				ji = append(ji, h)
			case "凶格":
				xiong = append(xiong, h)
			default:
				shen = append(shen, h)
			}
		}
		writeHits := func(title string, list []PatternHit) {
			if len(list) == 0 {
				return
			}
			b.WriteString("### " + title + "\n")
			for _, h := range list {
				b.WriteString(fmt.Sprintf("- **%s**（%s）：%s —— *%s*\n",
					h.Name, auspiceLabel(h.AuspiceScore), h.Summary, h.Classic))
			}
		}
		writeHits("🟢 吉格", ji)
		writeHits("🔴 凶格", xiong)
		writeHits("⚪ 其他", shen)
		b.WriteString("\n")
	} else {
		b.WriteString("## 命中格局\n")
		b.WriteString("本盘无强格局命中，请侧重值符值使落宫与三奇分布判读。\n\n")
	}

	// ======= 门加九宫克应（八宫各自的门宫组合） =======
	// 《奇门宝鉴》卷三逐宫断辞，每盘会自动命中 8 条（每个非中5宫的门宫组合一条）
	if mpHits := DetectDoorPalaceHits(pan); len(mpHits) > 0 {
		b.WriteString("## 门加九宫克应（奇门宝鉴古法，本盘 8 条逐宫断辞）\n")
		for _, h := range mpHits {
			auc := AuspiceLabelShort(h.Auspice)
			subj := ""
			if len(h.Subject) > 0 {
				subj = "（" + strings.Join(h.Subject, "/") + "）"
			}
			b.WriteString(fmt.Sprintf("- **%s**%s：%s\n", auc, subj, h.Summary))
		}
		b.WriteString("\n")
	}

	// ======= 类神直指 =======
	if qType != "" {
		if text := LeiShenDirective(pan, qType); text != "" {
			b.WriteString("## 类神直指（按问题类型锁定）\n")
			b.WriteString(text)
			b.WriteString("\n")
		}
	}

	// ======= 奇门取象原则 =======
	b.WriteString("## 奇门遁甲基础取象原则\n")
	b.WriteString("- **日干为我**，**时干为他**。甲遁于旬首六仪，找甲即找旬首所在宫。\n")
	b.WriteString("- **值符**为天盘机枢，主事之大势；**值使**为人盘先锋，主事之动向。\n")
	b.WriteString("- **三奇（乙丙丁）**落吉宫见吉门则大吉；**六仪（戊己庚辛壬癸）**中庚为天盗，六仪落墓/击刑/空亡皆忌。\n")
	b.WriteString("- **八门**：开休生为三吉门，伤杜景死惊各主所专；门迫（宫克门）力减，门得生（宫生门）力强。\n")
	b.WriteString("- **九星月令**：旺相有气、休囚死无力；九星五行与当月五行生克定其力。\n")
	b.WriteString("- **八神**：值符、六合、太阴、九天为吉；腾蛇、白虎、玄武、九地依场景辨。\n\n")

	// ======= 解局侧重 =======
	if focusGuide != "" {
		b.WriteString("## 解局侧重\n")
		b.WriteString(focusGuide + "\n\n")
	}

	// ======= 解局框架 =======
	b.WriteString("## 请按以下结构解局（每一步结论必须标注依据）\n")
	b.WriteString("1. **整体格局**（2-3 句）：本盘的吉凶大势（依据：值符值使落宫、命中格局、九星月令）。\n")
	b.WriteString("2. **用神分析**：围绕「类神直指」锁定的用神，说明其落宫、乘星门神、空亡入墓与否，判其得失。\n")
	b.WriteString("3. **关键信号解读**：从「本局关键信号」中挑 2-3 条最决定性的展开（如日干落宫、值符宫五行、三奇入墓等）。\n")
	b.WriteString("4. **动静判断**：值符值使落吉宫且无空亡入墓，利动；反之利静；结合三吉门与三吉将分布。\n")
	b.WriteString("5. **方位与时机**：指出对求测者最有利的**方位**（按九宫后天方位）和**时辰**（按值符值使所临地支）。\n")
	b.WriteString("6. **行动建议**（3-5 条）：分「宜」与「忌」，结合命中格局给出具体可行建议。\n")
	b.WriteString("7. **一句断语** + 信度自评（高/中/低）：用一句话给出最终判断，附信度。\n\n")
	b.WriteString("请用中文回答，兼顾传统术语与现代通俗表达；重解读轻复述。\n")

	return b.String()
}

// ============ 辅助 ============

func dashIfEmpty(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

func auspiceLabel(s int) string {
	switch {
	case s >= 2:
		return "大吉"
	case s == 1:
		return "吉"
	case s == -1:
		return "凶"
	case s <= -2:
		return "大凶"
	}
	return "平"
}

// renderMarkdownPan 生成 3×3 markdown 表格形式的九宫盘。
func renderMarkdownPan(pan *Pan) string {
	layout := [3][3]int{
		{3, 8, 1}, // 巽4, 离9, 坤2
		{2, 4, 6}, // 震3, 中5, 兑7
		{7, 0, 5}, // 艮8, 坎1, 乾6
	}
	var b strings.Builder
	b.WriteString("| 东南（巽4） | 正南（离9） | 西南（坤2） |\n")
	b.WriteString("|---|---|---|\n")
	renderRow(&b, pan, layout[0])
	b.WriteString("\n| 正东（震3） | 中宫 | 正西（兑7） |\n")
	b.WriteString("|---|---|---|\n")
	renderRow(&b, pan, layout[1])
	b.WriteString("\n| 东北（艮8） | 正北（坎1） | 西北（乾6） |\n")
	b.WriteString("|---|---|---|\n")
	renderRow(&b, pan, layout[2])
	return b.String()
}

func renderRow(b *strings.Builder, pan *Pan, row [3]int) {
	b.WriteString("|")
	for _, idx := range row {
		c := pan.Cells[idx]
		heaven := dashIfEmpty(c.HeavenStem)
		door := dashIfEmpty(c.Door)
		god := dashIfEmpty(c.God)
		marks := []string{}
		if c.IsVoid {
			marks = append(marks, "空")
		}
		if c.IsYima {
			marks = append(marks, "马")
		}
		if c.IsTianStemMu || c.IsEarthStemMu {
			marks = append(marks, "墓")
		}
		if c.IsJiXing {
			marks = append(marks, "刑")
		}
		if c.IsDoorPo {
			marks = append(marks, "迫")
		}
		mark := ""
		if len(marks) > 0 {
			mark = "（" + strings.Join(marks, "·") + "）"
		}
		fmt.Fprintf(b, " %s %s %s<br>天%s 地%s%s<br>%s |",
			god, c.Star, door, heaven, c.EarthStem, mark, c.PalaceName)
	}
	b.WriteString("\n")
}
