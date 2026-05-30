package main

import (
	"fmt"
	"strings"
)

// GenerateAIPrompt 根据占卜结果和问题生成 AI 解卦提示词
func GenerateAIPrompt(r DivinationResult, question string) string {
	if r.MainHex == nil {
		return ""
	}

	var b strings.Builder

	b.WriteString("你是一位精通《周易》的易学顾问，请根据以下卦象为我解卦。\n")
	b.WriteString("**重要注意**：\n")
	b.WriteString("1. 以下卦象素材是给你的原始数据，请直接**解读其含义**，不要在回答中大段复述素材原文（可简短引用关键句）。\n")
	b.WriteString("2. 每个结论请**标注依据**，形如「（依据：九五当位中正 + 二五正应）」或「（依据：变卦乾→坤，外象转阴）」。\n")
	b.WriteString("3. 周易长于**明义理、辨吉凶之所以然、示进退取舍之道**——请发挥其所长，重在阐发卦爻象意与处世之理；至于精确的应期时辰、具体方位，非周易所擅长，若无明确卦象支撑则不必强断，可点到为止或建议另以六壬、奇门细推。\n")
	b.WriteString("4. 素材末尾的「## 本卦关键信号」是代码预先抽取的本盘结论，可直接采用；其余小节为支撑材料。\n\n")

	// 问题
	if question != "" {
		b.WriteString("## 所问之事\n")
		b.WriteString(question + "\n\n")
	}

	// 起卦方式
	methodName := map[DivinationMethod]string{
		CoinMethod:   "铜钱法（金钱卦）",
		YarrowMethod: "蓍草法（大衍揲蓍）",
		NumberMethod: "数字起卦法",
	}[r.Method]
	b.WriteString("## 起卦方式\n")
	b.WriteString(methodName + "\n\n")

	// 起卦时间（干支、节气、时令消息卦）
	var timing *TimingInfo
	if !r.Time.IsZero() {
		timing = CaptureTiming(r.Time)
		b.WriteString(FormatTimingSection(timing))
	}

	// 六爻原始数值
	b.WriteString("## 六爻数值（初爻→上爻）\n")
	lineNames := []string{"初", "二", "三", "四", "五", "上"}
	typeNames := map[int]string{6: "老阴●", 7: "少阳—", 8: "少阴--", 9: "老阳○"}
	symbols := map[int]string{6: "─×─（变）", 7: "───", 8: "── ──", 9: "─○─（变）"}
	for i := 0; i < 6; i++ {
		yy := "九"
		li := lineInfo(r.Lines[i])
		if !li.IsYang {
			yy = "六"
		}
		b.WriteString(fmt.Sprintf("  第%d爻（%s%s）：值=%d %s %s\n",
			i+1, lineNames[i], yy, r.Lines[i], typeNames[r.Lines[i]], symbols[r.Lines[i]]))
	}
	b.WriteString("\n")

	// 本卦
	h := r.MainHex
	b.WriteString("## 本卦\n")
	b.WriteString(fmt.Sprintf("**第%d卦 · %s卦** %s\n", h.Number, h.Name, h.Symbol))
	b.WriteString(fmt.Sprintf("- 上卦（外卦）：%s %s（象征%s）\n",
		Trigrams[h.Upper].Symbol, h.Upper, Trigrams[h.Upper].Nature))
	b.WriteString(fmt.Sprintf("- 下卦（内卦）：%s %s（象征%s）\n",
		Trigrams[h.Lower].Symbol, h.Lower, Trigrams[h.Lower].Nature))
	b.WriteString(fmt.Sprintf("- 卦辞：%s\n", h.Judgment))
	b.WriteString(fmt.Sprintf("- 象辞：%s\n", h.Image))
	b.WriteString("- 各爻爻辞：\n")
	for i := 0; i < 6; i++ {
		line := h.Lines[i]
		changeMark := ""
		for _, pos := range r.ChangingPos {
			if pos == i+1 {
				changeMark = "（变）"
				break
			}
		}
		b.WriteString(fmt.Sprintf("  · %s：%s%s\n", line.Type, line.Text, changeMark))
	}
	b.WriteString("\n")

	// 变爻
	if len(r.ChangingPos) > 0 {
		b.WriteString("## 变爻\n")
		for _, pos := range r.ChangingPos {
			line := h.Lines[pos-1]
			marker := "老阳○"
			if r.Lines[pos-1] == 6 {
				marker = "老阴●"
			}
			b.WriteString(fmt.Sprintf("- 第%d爻（%s，%s）：%s · %s\n",
				pos, lineNames[pos-1], marker, line.Type, line.Text))
		}
		b.WriteString("\n")
	}

	// 变卦
	if r.ChangeHex != nil {
		ch := r.ChangeHex
		b.WriteString("## 变卦（之卦）\n")
		b.WriteString(fmt.Sprintf("**第%d卦 · %s卦** %s\n", ch.Number, ch.Name, ch.Symbol))
		b.WriteString(fmt.Sprintf("- 上卦：%s %s（%s）\n",
			Trigrams[ch.Upper].Symbol, ch.Upper, Trigrams[ch.Upper].Nature))
		b.WriteString(fmt.Sprintf("- 下卦：%s %s（%s）\n",
			Trigrams[ch.Lower].Symbol, ch.Lower, Trigrams[ch.Lower].Nature))
		b.WriteString(fmt.Sprintf("- 卦辞：%s\n", ch.Judgment))
		b.WriteString(fmt.Sprintf("- 象辞：%s\n", ch.Image))
		b.WriteString("- 各爻爻辞：\n")
		for i := 0; i < 6; i++ {
			line := ch.Lines[i]
			b.WriteString(fmt.Sprintf("  · %s：%s\n", line.Type, line.Text))
		}
		b.WriteString("\n")
	}

	// 衍生卦：互卦、错卦、综卦
	derived := DeriveHexagrams(r.Lines)
	b.WriteString("## 衍生卦象（辅助参考）\n")
	if derived.Mutual != nil {
		b.WriteString(fmt.Sprintf("- **互卦**（取二、三、四爻为下，三、四、五爻为上）：第%d卦 %s卦 %s\n",
			derived.Mutual.Number, derived.Mutual.Name, derived.Mutual.Symbol))
		b.WriteString(fmt.Sprintf("  · 卦辞：%s\n", derived.Mutual.Judgment))
		b.WriteString("  · 含义：揭示事件内部潜在的运作机制与中间过程。\n")
	}
	if derived.Opposite != nil {
		b.WriteString(fmt.Sprintf("- **错卦**（六爻阴阳全反）：第%d卦 %s卦 %s\n",
			derived.Opposite.Number, derived.Opposite.Name, derived.Opposite.Symbol))
		b.WriteString(fmt.Sprintf("  · 卦辞：%s\n", derived.Opposite.Judgment))
		b.WriteString("  · 含义：事情的对立面、背面；当前所忽视或排斥的视角。\n")
	}
	if derived.Reverse != nil {
		isSelf := r.MainHex != nil && derived.Reverse.Number == r.MainHex.Number
		if isSelf {
			b.WriteString(fmt.Sprintf("- **综卦**（上下颠倒）：与本卦相同（第%d卦 %s卦），属上下对称之卦，正反一体。\n",
				derived.Reverse.Number, derived.Reverse.Name))
		} else {
			b.WriteString(fmt.Sprintf("- **综卦**（上下颠倒）：第%d卦 %s卦 %s\n",
				derived.Reverse.Number, derived.Reverse.Name, derived.Reverse.Symbol))
			b.WriteString(fmt.Sprintf("  · 卦辞：%s\n", derived.Reverse.Judgment))
			b.WriteString("  · 含义：换位思考——从对方或相反立场审视同一件事。\n")
		}
	}
	b.WriteString("\n")

	// 卦序前后邻卦（《序卦传》脉络）
	b.WriteString(FormatAdjacentSection(r.MainHex))

	// 时令消息卦与本卦的关系
	if timing != nil {
		if rel := FormatTimingHexRelation(timing, r.MainHex, r.Lines); rel != "" {
			b.WriteString("## 时令与卦象的关系\n")
			b.WriteString(rel)
			b.WriteString("\n")
		}
	}

	// 爻位关系分析
	b.WriteString("## 爻位关系分析\n")
	b.WriteString("（当位：阳爻居1/3/5位、阴爻居2/4/6位为当位；中：二爻为下卦之中，五爻为上卦之中；应：初-四、二-五、三-上相对，阴阳相反为有应）\n")
	b.WriteString(FormatLinePositionAnalysis(r.Lines))
	b.WriteString("\n")

	// 本卦关键信号（代码预抽取的本盘专属结论）
	if sig := ExtractZhouyiSignals(r, timing); sig != "" {
		b.WriteString("## 本卦关键信号（代码预抽取，请直接采用）\n")
		b.WriteString(sig)
		b.WriteString("\n")
	}

	// 解卦规则（朱熹《易学启蒙·考变占》七条占法）
	b.WriteString("## 解卦规则（朱熹《易学启蒙·考变占》）\n")
	n := len(r.ChangingPos)
	switch n {
	case 0:
		b.WriteString("无变爻：以本卦卦辞（彖辞）断，参以象辞；卦象稳定，局势不变。\n")
	case 1:
		b.WriteString(fmt.Sprintf("一爻变（第%d爻）：以本卦该变爻爻辞断。\n", r.ChangingPos[0]))
	case 2:
		b.WriteString(fmt.Sprintf("二爻变（第%d、%d爻）：以本卦两变爻爻辞断，仍以上爻（第%d爻）为主。\n",
			r.ChangingPos[0], r.ChangingPos[1], r.ChangingPos[1]))
	case 3:
		b.WriteString(fmt.Sprintf("三爻变（第%v爻）：占本卦与变卦之卦辞，**本卦为贞**（主，往、先），**变卦为悔**（次，来、后）。\n",
			r.ChangingPos))
	case 4:
		nonChanging := nonChangingPositions(r)
		b.WriteString(fmt.Sprintf("四爻变：以变卦中两个不变爻（第%d、%d爻）爻辞断，仍以下爻（第%d爻）为主。\n",
			nonChanging[0], nonChanging[1], nonChanging[0]))
	case 5:
		nonChanging := nonChangingPositions(r)
		b.WriteString(fmt.Sprintf("五爻变：以变卦中唯一不变的第%d爻爻辞断。\n", nonChanging[0]))
	case 6:
		switch h.Name {
		case "乾":
			b.WriteString("六爻皆变（乾卦）：占「用九：见群龙无首，吉」。\n")
		case "坤":
			b.WriteString("六爻皆变（坤卦）：占「用六：利永贞」。\n")
		default:
			b.WriteString("六爻皆变：占变卦卦辞，本卦为背景参考。\n")
		}
	}
	b.WriteString("\n")

	// 问题类型与解卦侧重
	b.WriteString(FormatQuestionTypeSection(r.QuestionType))

	// 解卦要求
	b.WriteString("## 请按以下结构解卦（每一步结论必须标注依据）\n")
	idx := 1
	b.WriteString(fmt.Sprintf("%d. **卦象总论**：用2-3句话概括此卦对所问之事的整体象意（依据：本卦关键信号或卦辞）。\n", idx))
	idx++
	b.WriteString(fmt.Sprintf("%d. **卦象分析**：分析上下卦（内外卦）的五行、象意及其相互关系，并结合爻位关系（当位/中/应）指出关键爻。\n", idx))
	idx++
	if len(r.ChangingPos) > 0 {
		b.WriteString(fmt.Sprintf("%d. **变爻解析**：逐一解释变爻爻辞，结合所问之事给出具体含义。\n", idx))
		idx++
		b.WriteString(fmt.Sprintf("%d. **变卦趋势**：说明由本卦变至变卦意味着局势如何演变。\n", idx))
		idx++
	}
	b.WriteString(fmt.Sprintf("%d. **衍生卦参照**：结合互卦（内部过程）、错卦（对立面）、综卦（换位视角）补充解读，指出被忽视的层面或潜在转机。不必逐卦详述，择其与所问之事最相关者展开。\n", idx))
	idx++
	b.WriteString(fmt.Sprintf("%d. **义理与处世之道**：本卦最核心的进退取舍之道是什么？结合所问之事，给出具体的行动建议或注意事项，重在「该以什么心态、循什么原则去应对」。\n", idx))
	idx++
	b.WriteString(fmt.Sprintf("%d. **一句断语**：用一句话给出最终判断（吉/凶/平，宜/忌何事）。若涉及精确应期或方位而卦象无明确支撑，宜如实说明周易于此非所长、不强断。\n", idx))
	b.WriteString("\n")
	b.WriteString("请用中文回答，语言兼顾传统易学用语与现代通俗表达，深入浅出。\n")

	return b.String()
}

// 八卦 → 五行（用于动态信号的生克判断）
var trigramWuXingName = map[string]string{
	"乾": "金", "兑": "金",
	"震": "木", "巽": "木",
	"坤": "土", "艮": "土",
	"坎": "水",
	"离": "火",
}

// ExtractZhouyiSignals 从本卦+变卦+爻位+时令中抽取本盘专属结论（3-6 条）
// 返回多行 markdown，放入"## 本卦关键信号"小节。
func ExtractZhouyiSignals(r DivinationResult, timing *TimingInfo) string {
	if r.MainHex == nil {
		return ""
	}
	h := r.MainHex
	lineNames := []string{"初", "二", "三", "四", "五", "上"}
	positions := AnalyzeLinePositions(r.Lines)

	var b strings.Builder

	// —— 1. 上下卦五行关系
	upWX := trigramWuXingName[h.Upper]
	loWX := trigramWuXingName[h.Lower]
	if upWX != "" && loWX != "" {
		rel := trigramRelation(h.Upper, h.Lower)
		b.WriteString(fmt.Sprintf("- **上下卦五行关系**：上卦 %s（%s，外/他）对下卦 %s（%s，内/我）—— %s\n",
			h.Upper, upWX, h.Lower, loWX, rel))
	}

	// —— 2. 关键爻（集齐当位+居中+有应者最旺）
	type scored struct {
		pos   int
		score int
		attrs []string
	}
	var scoredList []scored
	for i, p := range positions {
		s := 0
		var attrs []string
		if p.IsProper {
			s++
			attrs = append(attrs, "当位")
		}
		if p.IsCentral {
			s++
			attrs = append(attrs, "居中")
		}
		if strings.Contains(p.Relation, "有应") {
			s++
			attrs = append(attrs, "有应")
		}
		if s > 0 {
			scoredList = append(scoredList, scored{pos: i + 1, score: s, attrs: attrs})
		}
	}
	// 找最高分
	best := scored{score: -1}
	for _, sc := range scoredList {
		if sc.score > best.score {
			best = sc
		}
	}
	if best.score >= 2 {
		b.WriteString(fmt.Sprintf("- **关键爻**：第%d爻（%s爻）兼 %s —— 为全卦最有力之爻，所言多为本事之主。\n",
			best.pos, lineNames[best.pos-1], strings.Join(best.attrs, "、")))
	}

	// —— 3. 变爻位置点评
	if len(r.ChangingPos) > 0 {
		posMeanings := map[int]string{
			1: "事之始、根基初动", 2: "臣/内、中正自守", 3: "过渡、危险之地",
			4: "近君/外交、进退之位", 5: "君/主、权位中正", 6: "事之终、过亢之极",
		}
		var notes []string
		for _, p := range r.ChangingPos {
			pi := positions[p-1]
			tag := ""
			if pi.IsProper && pi.IsCentral {
				tag = "，中正当位之动→权位主动"
			} else if pi.IsCentral {
				tag = "，居中而不当位→身在其位而德未配"
			} else if pi.IsProper {
				tag = "，当位而不居中→循分守位"
			} else {
				tag = "，不当不中→位置不稳之变"
			}
			notes = append(notes, fmt.Sprintf("第%d爻(%s)%s", p, posMeanings[p], tag))
		}
		b.WriteString(fmt.Sprintf("- **变爻之位**：%s。\n", strings.Join(notes, "；")))
	}

	// —— 4. 本卦 → 变卦：五行走向
	if r.ChangeHex != nil {
		ch := r.ChangeHex
		mainUP := trigramWuXingName[h.Upper]
		chgUP := trigramWuXingName[ch.Upper]
		if mainUP != "" && chgUP != "" {
			dir := wuXingNameRelation(mainUP, chgUP)
			b.WriteString(fmt.Sprintf("- **变卦走向**：本卦上卦 %s(%s) → 变卦上卦 %s(%s) —— %s\n",
				h.Upper, mainUP, ch.Upper, chgUP, dir))
		}
	} else {
		b.WriteString("- **变卦走向**：无变爻 → 局势稳定，以本卦彖辞为定局。\n")
	}

	// —— 5. 时令与本卦
	if timing != nil && timing.MonthlyHex != nil && timing.MonthlyHex.Number != h.Number {
		mhWX := trigramWuXingName[timing.MonthlyHex.Upper]
		hUpWX := trigramWuXingName[h.Upper]
		if mhWX != "" && hUpWX != "" {
			b.WriteString(fmt.Sprintf("- **时令卦气**：当下时令消息卦为「%s」（%s属%s）对本卦上卦%s(%s)—— %s；得时则势顺，失时则势阻。\n",
				timing.MonthlyHex.Name, timing.MonthlyHex.Upper, mhWX, h.Upper, hUpWX,
				wuXingNameRelation(mhWX, hUpWX)))
		}
	}

	// —— 6. 阴阳比例
	yangCount := 0
	for _, v := range r.Lines {
		if v == 7 || v == 9 {
			yangCount++
		}
	}
	var yinYangNote string
	switch {
	case yangCount >= 5:
		yinYangNote = "阳盛（主动、刚进、外显）"
	case yangCount <= 1:
		yinYangNote = "阴盛（主静、柔退、内藏）"
	case yangCount == 3:
		yinYangNote = "阴阳平衡（动静有度）"
	default:
		yinYangNote = fmt.Sprintf("阳 %d / 阴 %d（%s略占）", yangCount, 6-yangCount, ternary(yangCount > 3, "阳", "阴"))
	}
	b.WriteString(fmt.Sprintf("- **阴阳格局**：%s。\n", yinYangNote))

	return b.String()
}

func trigramRelation(upper, lower string) string {
	u, ok1 := trigramWuXingName[upper]
	l, ok2 := trigramWuXingName[lower]
	if !ok1 || !ok2 {
		return ""
	}
	return wuXingNameRelation(u, l)
}

// 纯字符串五行关系（避免依赖 liuren 包）
func wuXingNameRelation(a, b string) string {
	if a == b {
		return "比和（同气相求、和合）"
	}
	gen := map[string]string{"金": "水", "水": "木", "木": "火", "火": "土", "土": "金"}
	ke := map[string]string{"金": "木", "木": "土", "土": "水", "水": "火", "火": "金"}
	switch {
	case gen[a] == b:
		return fmt.Sprintf("%s生%s（相生，外扶内 / 上助下）", a, b)
	case gen[b] == a:
		return fmt.Sprintf("%s生%s（反哺，内助外 / 下养上）", b, a)
	case ke[a] == b:
		return fmt.Sprintf("%s克%s（相战，外制内 / 上压下）", a, b)
	case ke[b] == a:
		return fmt.Sprintf("%s克%s（反克，内制外 / 下犯上）", b, a)
	}
	return "无直接生克"
}

func ternary(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

// nonChangingPositions 返回无变爻的位置（1-6）
func nonChangingPositions(r DivinationResult) []int {
	changing := map[int]bool{}
	for _, p := range r.ChangingPos {
		changing[p] = true
	}
	var result []int
	for i := 1; i <= 6; i++ {
		if !changing[i] {
			result = append(result, i)
		}
	}
	return result
}

// PrintAIPrompt 格式化输出提示词到终端
func PrintAIPrompt(r DivinationResult, question string) {
	prompt := GenerateAIPrompt(r, question)
	if prompt == "" {
		return
	}
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                  AI 解卦提示词（可直接复制）                ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println(prompt)
	fmt.Println("═══════════════════════════════════════════════════════════")
}
