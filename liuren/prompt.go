package liuren

import (
	"fmt"
	"strings"
)

// GenerateAIPrompt 为大六壬盘面生成解课提示词。
//
//	question       所问之事（可空）
//	focusGuide     问题类型侧重文字（由主包传入，复用 questiontype.go 的 FocusGuide）
//	leishen        类神直指块（由主包传入：按问题类型锁定的类神落点；可空）
func GenerateAIPrompt(pan *Pan, question, focusGuide string, leishenBlock ...string) string {
	var b strings.Builder
	b.WriteString("你是一位精通大六壬（三式之一，论人事最精，尤长于应期与情态）的术数顾问，请根据以下盘面为我断课。\n")
	b.WriteString("**重要注意**：\n")
	b.WriteString("1. 以下盘面为给你的原始数据，请直接**解读其含义**，不要大段复述盘面原文（关键字句可简短引用）。\n")
	b.WriteString("2. 每个结论请**标注依据**，形如「（依据：末传乘白虎且空亡）」或「（依据：日上乘玄武入三传）」。\n")
	b.WriteString("3. 六壬主**人事情态之断**——长于**看人心向背、断事情曲折、定应期早晚、察彼我动静**，宜用于具体人事、谋望成败、寻人问事、情态吉凶。请充分发挥这一所长，把「人心如何、事情怎么走、何时应验」断细断准；至于宏观方位布局、国运大势，可点到为止，不必越俎代庖去抢奇门、太乙之长。\n")
	b.WriteString("4. 素材中「## 本课关键信号」「## 类神直指」为代码预抽取的结论，可直接采用；其余小节为支撑材料。\n\n")

	if question != "" {
		b.WriteString("## 所问之事\n")
		b.WriteString(question + "\n\n")
	}

	// 起课时间
	b.WriteString("## 起课时间\n")
	b.WriteString("- 阳历：" + pan.Ctx.Time.Format("2006-01-02 15:04") + "\n")
	b.WriteString(fmt.Sprintf("- 日柱：%s%s（六十甲子第 %d）\n", pan.Ctx.Gan, pan.Ctx.DayZhi, pan.Ctx.JiaziIndex+1))
	b.WriteString(fmt.Sprintf("- 占时：%s\n", pan.Ctx.ZhanShi))
	b.WriteString(fmt.Sprintf("- 月将：%s（%s 后，别称 %s）\n", pan.Ctx.YueJiang, pan.Ctx.QiName, ZhiBieMing[pan.Ctx.YueJiang]))
	zy := "夜占"
	if pan.Ctx.ZhouYe {
		zy = "昼占"
	}
	b.WriteString("- 昼夜：" + zy + "\n")
	kong := pan.Ctx.XunKongPair()
	b.WriteString(fmt.Sprintf("- 旬空：%s、%s\n\n", kong[0], kong[1]))

	// 天地盘
	b.WriteString("## 天地盘\n")
	b.WriteString("| 地盘 |")
	for i := 0; i < 12; i++ {
		b.WriteString(fmt.Sprintf(" %s |", Zhi(i)))
	}
	b.WriteString("\n|---" + strings.Repeat("|---", 12) + "|\n")
	b.WriteString("| 天盘 |")
	for i := 0; i < 12; i++ {
		b.WriteString(fmt.Sprintf(" %s |", pan.TianPan[i]))
	}
	b.WriteString("\n| 天将 |")
	for i := 0; i < 12; i++ {
		b.WriteString(fmt.Sprintf(" %s |", pan.TianJiang[i]))
	}
	b.WriteString("\n\n")

	// 四课
	b.WriteString("## 四课（右起为第一课）\n")
	b.WriteString("| 课 | 上神 | 下神 | 关系 |\n|---|---|---|---|\n")
	names := []string{"第一课", "第二课", "第三课", "第四课"}
	for i, ke := range pan.SiKe {
		lower := fmt.Sprintf("%s", ke.Lower)
		if ke.Index == 1 {
			lower = fmt.Sprintf("%s（%s 寄）", ke.Lower, pan.Ctx.Gan)
		}
		b.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", names[i], ke.Upper, lower, ke.Relation))
	}
	b.WriteString("\n")

	// 日辰格（《大全》卷三 p534 日辰章）
	if rzs := DetectRiZhenGe(pan); len(rzs) > 0 {
		b.WriteString("## 日辰格（古法日辰互动 · 《大全》卷三）\n")
		for _, rz := range rzs {
			b.WriteString(fmt.Sprintf("- **%s** —— %s\n", rz.Name, rz.Judge))
		}
		b.WriteString("\n")
	}

	// 三传
	b.WriteString(fmt.Sprintf("## 三传（发传法：%s · %s）\n", pan.SanChuan.Method, pan.KeTi.Name))
	b.WriteString("| 传 | 天神 | 天将 | 六亲 | 空亡 |\n|---|---|---|---|---|\n")
	for _, ce := range []ChuanEntry{pan.SanChuan.Chu, pan.SanChuan.Zhong, pan.SanChuan.Mo} {
		kg := "—"
		if ce.IsKong {
			kg = "空"
		}
		b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
			ce.Name, ce.Zhi, ce.TianJiang, ce.LiuQin, kg))
	}
	b.WriteString("\n")

	// 三传天将临宫断辞（《大全》卷二）
	b.WriteString("### 三传天将临宫（古法断辞）\n")
	for _, ce := range []ChuanEntry{pan.SanChuan.Chu, pan.SanChuan.Zhong, pan.SanChuan.Mo} {
		// 天将临的地盘位 = 天神 ce.Zhi 在天盘所临的地盘位
		var diPos Zhi = -1
		for i, up := range pan.TianPan {
			if up == ce.Zhi {
				diPos = Zhi(i)
				break
			}
		}
		if diPos < 0 {
			continue
		}
		line := TianJiangPalaceLine(ce.TianJiang, diPos)
		if line != "" {
			b.WriteString(fmt.Sprintf("- %s（%s 乘 %s 临地盘 %s）：%s\n",
				ce.Name, ce.Zhi, ce.TianJiang, diPos, line))
		}
	}
	b.WriteString("\n")

	// 日干/日支天将临宫断辞（占课"我"与"事"两端）
	ganDi := GanJiGong[pan.Ctx.Gan]
	zhiDi := pan.Ctx.DayZhi
	ganTJLine := TianJiangPalaceLine(pan.TianJiang[ganDi], ganDi)
	zhiTJLine := TianJiangPalaceLine(pan.TianJiang[zhiDi], zhiDi)
	if ganTJLine != "" || zhiTJLine != "" {
		b.WriteString("### 日干日支天将临宫\n")
		if ganTJLine != "" {
			b.WriteString(fmt.Sprintf("- 日干寄宫 %s 上 %s 临之：%s\n", ganDi, pan.TianJiang[ganDi], ganTJLine))
		}
		if zhiTJLine != "" {
			b.WriteString(fmt.Sprintf("- 日支 %s 上 %s 临之：%s\n", zhiDi, pan.TianJiang[zhiDi], zhiTJLine))
		}
		b.WriteString("\n")
	}

	// 课体
	b.WriteString("## 课体\n")
	b.WriteString(fmt.Sprintf("%s —— %s\n", pan.KeTi.Name, pan.KeTi.Summary))
	if pan.Ctx.QuestionType != "" {
		if byQ := KeTiSummaryByQuestion(pan.KeTi.Name, pan.Ctx.QuestionType); byQ != "" {
			b.WriteString(fmt.Sprintf("\n按问题类型断辞：%s\n", byQ))
		}
	}
	if len(pan.Tags) > 0 {
		b.WriteString("\n附加课格：\n")
		for _, tg := range pan.Tags {
			b.WriteString(fmt.Sprintf("- **%s** —— %s\n", tg.Name, tg.Summary))
		}
	}
	b.WriteString("\n")

	// 天将乘临禁忌
	if len(pan.Taboos) > 0 {
		b.WriteString("## 天将乘临禁忌（古法警示）\n")
		b.WriteString("《大全》卷一 p498：「**贵神天空不乘辰戌，玄武六合不乘丑未**」。本盘命中以下禁忌，主该天将之力变形：\n")
		for _, tb := range pan.Taboos {
			b.WriteString(fmt.Sprintf("- **%s 临 %s** —— %s\n", tb.TianJiang, tb.DiZhi, tb.Note))
		}
		b.WriteString("\n")
	}

	// 神煞
	if len(pan.ShenSha) > 0 {
		b.WriteString("## 神煞落位\n")
		for _, ss := range pan.ShenSha {
			if ss.Zhi < 0 {
				continue
			}
			b.WriteString(fmt.Sprintf("- %s 落 %s（%s）—— %s\n", ss.Name, ss.Zhi, pan.TianJiang[ss.Zhi], ss.Desc))
		}
		b.WriteString("\n")
	}

	// 年命
	if pan.NianMing != nil {
		b.WriteString("## 年命（救应/损应枢机）\n")
		if pan.NianMing.BenMing != nil {
			bm := pan.NianMing.BenMing
			kg := ""
			if bm.IsKong {
				kg = "（空）"
			}
			b.WriteString(fmt.Sprintf("- 本命 %s，乘神 %s（%s · %s）%s\n", bm.Zhi, bm.Upper, bm.TianJiang, bm.LiuQin, kg))
			b.WriteString(fmt.Sprintf("  · **%s**：%s\n", bm.Ying, bm.Ying.Desc()))
		}
		if pan.NianMing.XingNian != nil {
			xn := pan.NianMing.XingNian
			kg := ""
			if xn.IsKong {
				kg = "（空）"
			}
			b.WriteString(fmt.Sprintf("- 行年 %s，乘神 %s（%s · %s）%s\n", xn.Zhi, xn.Upper, xn.TianJiang, xn.LiuQin, kg))
			b.WriteString(fmt.Sprintf("  · **%s**：%s\n", xn.Ying, xn.Ying.Desc()))
		}
		if pan.NianMing.BMXNRel != "" {
			b.WriteString(fmt.Sprintf("- 本命×行年：**%s** —— %s\n", pan.NianMing.BMXNRel, pan.NianMing.BMXNDesc))
		}
		b.WriteString("\n古法：课传吉而年命凶则吉中藏险；课传凶而年命救应，则逢凶化吉。年命为最终定夺之权柄。\n\n")
	}

	// 毕法赋匹配（自动命中的条目）
	if len(pan.BiFa) > 0 {
		b.WriteString("## 毕法赋匹配（本盘自动命中）\n")
		for _, e := range pan.BiFa {
			b.WriteString(fmt.Sprintf("- **第%d · %s**：%s（%s）\n", e.Number, e.Title, e.Text, e.Note))
		}
		b.WriteString("\n")
	}

	// 毕法赋百条知识库（给 AI 作参考字典）
	b.WriteString("## 毕法赋百条参考（凌福之《大六壬毕法赋》）\n")
	b.WriteString("以下 100 条为古籍断语全文，供判读时查阅；如本盘具备条件但未被自动匹配，可自行引用：\n\n")
	for _, e := range BiFaCatalog() {
		b.WriteString(fmt.Sprintf("- 第%d · %s：%s\n", e.Number, e.Title, e.Note))
	}
	b.WriteString("\n")

	// 云霄赋（《大全》卷四）— 综合断课警句知识库
	b.WriteString("## 云霄赋参考（《大全》卷四）\n")
	b.WriteString("以下为综合断课的古法警句，可与本盘对照：\n\n")
	for _, e := range YunXiaoFu {
		b.WriteString(fmt.Sprintf("- **%s**：%s\n", e.Title, e.Note))
	}
	b.WriteString("\n")

	// 心印赋（《大全》卷三）— 日辰/年命/月将关系总纲
	b.WriteString("## 心印赋参考（《大全》卷三）\n")
	b.WriteString("以下为日辰、年命、月将关系的纲领断辞，是断课的『语法』：\n\n")
	for _, e := range XinYinFu {
		b.WriteString(fmt.Sprintf("- **%s**：%s\n", e.Title, e.Note))
	}
	b.WriteString("\n")

	// 括囊赋（《大全》卷四）— 起课操作手册
	b.WriteString("## 括囊赋参考（《大全》卷四）\n")
	b.WriteString("以下为起课与断课的操作手册，铺整个流程：\n\n")
	for _, e := range KuoNangFu {
		b.WriteString(fmt.Sprintf("- **%s**：%s\n", e.Title, e.Note))
	}
	b.WriteString("\n")

	// 本课关键信号（代码预抽取的本盘结论）
	if sig := ExtractLiuRenSignals(pan); sig != "" {
		b.WriteString("## 本课关键信号（代码预抽取，请直接采用）\n")
		b.WriteString(sig)
		b.WriteString("\n")
	}

	// 类神直指（由调用方按问题类型预先锁定的类神落点）
	if len(leishenBlock) > 0 && leishenBlock[0] != "" {
		b.WriteString("## 类神直指（按问题类型锁定）\n")
		b.WriteString(leishenBlock[0])
		b.WriteString("\n")
	}

	// 解课侧重
	if focusGuide != "" {
		b.WriteString("## 解课侧重\n")
		b.WriteString(focusGuide + "\n\n")
	}

	// 解课框架
	b.WriteString("## 请按以下结构断课（每一步结论必须标注依据）\n")
	b.WriteString("1. **课象总论**：用 2-3 句概括此盘对所问之事的整体象意（依据：本课关键信号或课体）。\n")
	b.WriteString("2. **四课分析**：以日干为我，日支为事，四课上神为环境；指出关键一课并说明理由。\n")
	b.WriteString("3. **三传解析**：初传之因、中传之经过、末传之结局，结合所乘天将与六亲；注意三传五行流向。\n")
	b.WriteString("4. **类神深究**：围绕「类神直指」锁定的类神，说明其得失、动静、落空与否，得出吉凶。\n")
	b.WriteString("5. **旬空与神煞**：若关键之神落空，说明虚浮或不实之应。\n")
	b.WriteString("6. **时令与卦气**：结合节气、月将的太阳躔次说明时势。\n")
	b.WriteString("7. **综合建议与一句断语**：给出具体行动建议并用一句话断吉凶宜忌，附信度自评（高/中/低）。\n\n")
	b.WriteString("请用中文回答，兼顾传统术语（类神、乘将、涉害、见机）与现代通俗表达。\n")
	return b.String()
}

// ExtractLiuRenSignals 从本盘抽取 4-6 条本课专属关键信号（非通用口诀）
func ExtractLiuRenSignals(pan *Pan) string {
	if pan == nil {
		return ""
	}
	var b strings.Builder

	// —— 1. 日上神：我之状态
	ganJi := GanJiGong[pan.Ctx.Gan]
	ganUpper := pan.TianPan[ganJi]
	ganTJ := pan.TianJiang[ganJi]
	ganLQ := LiuQinOfZhiByGan(ganUpper, pan.Ctx.Gan)
	kongPair := pan.Ctx.XunKongPair()
	isKong := func(z Zhi) bool { return z == kongPair[0] || z == kongPair[1] }
	ganKong := ""
	if isKong(ganUpper) {
		ganKong = "（空亡）"
	}
	b.WriteString(fmt.Sprintf("- **日上神（我之状态）**：日干 %s 寄宫 %s，上乘 **%s**，乘 %s，对我为 %s%s。",
		pan.Ctx.Gan, ganJi, ganUpper, ganTJ, ganLQ, ganKong))
	switch ganLQ {
	case LQFuMu:
		b.WriteString("父母生我 → 得长辈/文书/庇护之助。\n")
	case LQXiongDi:
		b.WriteString("兄弟同气 → 平辈相助或竞争夹杂。\n")
	case LQZiSun:
		b.WriteString("子孙泄我 → 解厄、吐故纳新，亦主耗费精神。\n")
	case LQQiCai:
		b.WriteString("妻财为我所制 → 主财帛、所欲之物到手。\n")
	case LQGuanGui:
		b.WriteString("官鬼克我 → 压力、官非、疾病或责任在身。\n")
	}

	// —— 2. 三传流向
	c := pan.SanChuan.Chu.Zhi
	z := pan.SanChuan.Zhong.Zhi
	m := pan.SanChuan.Mo.Zhi
	wc, wz, wm := ZhiWuXing[c], ZhiWuXing[z], ZhiWuXing[m]
	var flow string
	switch {
	case WuXingGenerates(wc, wz) && WuXingGenerates(wz, wm):
		flow = "三传递生 → 事势顺延、层层推进、有始有终"
	case WuXingOvercomes(wc, wz) && WuXingOvercomes(wz, wm):
		flow = "三传递克 → 节节受阻、事多波折"
	case WuXingGenerates(wc, wz) && WuXingOvercomes(wz, wm):
		flow = "先生后克 → 初顺而终败，虎头蛇尾"
	case WuXingOvercomes(wc, wz) && WuXingGenerates(wz, wm):
		flow = "先克后生 → 初阻而终成，否极泰来"
	case wc == wm:
		flow = "首尾同气 → 事情循环回复"
	default:
		flow = fmt.Sprintf("流向不规整（%s→%s→%s）→ 变数多，须看关键一传", wc, wz, wm)
	}
	b.WriteString(fmt.Sprintf("- **三传流向**：%s（初 %s%s · 中 %s%s · 末 %s%s）\n",
		flow, c, wc, z, wz, m, wm))

	// —— 3. 末传吉凶信号
	moTJ := pan.SanChuan.Mo.TianJiang
	moKong := ""
	if pan.SanChuan.Mo.IsKong {
		moKong = "（空亡，事无终或虚应）"
	}
	b.WriteString(fmt.Sprintf("- **末传归宿**：末传 %s 乘 **%s**，对我为 %s%s。",
		m, moTJ, pan.SanChuan.Mo.LiuQin, moKong))
	switch moTJ {
	case TJQingLong, TJLiuHe, TJTaiChang, TJTianHou, TJGuiRen:
		b.WriteString("末乘吉将 → 结局趋吉。\n")
	case TJBaiHu, TJXuanWu, TJGouChen, TJTengShe:
		b.WriteString("末乘凶将 → 结局有隐患。\n")
	default:
		b.WriteString("\n")
	}

	// —— 4. 贵人位置与顺逆
	guis := GuiRenByGan[pan.Ctx.Gan]
	guiZhi := guis[1]
	zy := "夜"
	if pan.Ctx.ZhouYe {
		guiZhi = guis[0]
		zy = "昼"
	}
	// 贵人在天盘的位置（地盘索引，即其临哪宫）
	guiDi := -1
	for i := 0; i < 12; i++ {
		if pan.TianPan[i] == guiZhi {
			guiDi = i
			break
		}
	}
	// 顺治：贵人临地盘亥子丑寅卯辰（0..5 对应子丑…亥；此处用地支序号）；
	// 简化：贵人所临地支如属"亥子丑寅卯辰"六阳位为顺、其余为逆
	shunZhi := map[Zhi]bool{Hai: true, Zi: true, Chou: true, Yin: true, Mao: true, Chen: true}
	var shunNi string
	if guiDi >= 0 {
		if shunZhi[Zhi(guiDi)] {
			shunNi = "顺治（贵人得位，助力有力）"
		} else {
			shunNi = "逆治（贵人失位，助力减半）"
		}
		b.WriteString(fmt.Sprintf("- **贵人位势**：%s占，贵人为 %s，临地盘 %s 宫 —— %s。\n",
			zy, guiZhi, Zhi(guiDi), shunNi))
	}

	// —— 5. 空亡信号汇总
	var kongs []string
	if pan.SanChuan.Chu.IsKong {
		kongs = append(kongs, "初传")
	}
	if pan.SanChuan.Zhong.IsKong {
		kongs = append(kongs, "中传")
	}
	if pan.SanChuan.Mo.IsKong {
		kongs = append(kongs, "末传")
	}
	if len(kongs) > 0 {
		b.WriteString(fmt.Sprintf("- **空亡警示**：三传中 %s 落空亡（%s、%s）→ 对应位置之事虚浮不实。\n",
			strings.Join(kongs, "、"), kongPair[0], kongPair[1]))
	}

	// —— 6. 课体定性
	ktName := pan.KeTi.Name
	var ktBias string
	switch {
	case strings.Contains(ktName, "元首") || strings.Contains(ktName, "重审"):
		ktBias = "常规课体 → 事有主次明确"
	case strings.Contains(ktName, "涉害"):
		ktBias = "涉害课 → 事多波折，须深思审断"
	case strings.Contains(ktName, "伏吟"):
		ktBias = "伏吟课 → 事机不动，宜静守"
	case strings.Contains(ktName, "返吟"):
		ktBias = "返吟课 → 事有反复、动荡难安"
	case strings.Contains(ktName, "八专"):
		ktBias = "八专课 → 事涉专一、少变化"
	case strings.Contains(ktName, "别责"):
		ktBias = "别责课 → 事有偏枯，主独立"
	case strings.Contains(ktName, "昴星"):
		ktBias = "昴星课 → 阴谋阻隔，须防暗算"
	}
	if ktBias != "" {
		b.WriteString(fmt.Sprintf("- **课体定性**：%s —— %s。\n", ktName, ktBias))
	}

	return b.String()
}
