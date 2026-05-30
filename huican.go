package main

import (
	"fmt"
	"strings"
	"time"

	"zhouyi/liuren"
	"zhouyi/qimen"
)

// HuCanResult 三式互参结果：同一时刻的周易 × 大六壬 × 奇门遁甲三盘
type HuCanResult struct {
	Zhouyi    DivinationResult
	LiuRenPan *liuren.Pan
	QimenPan  *qimen.Pan
	Question  string
	QType     QuestionType
}

// HuCanDivine 以同一时刻并行起周易卦（铜钱法）、大六壬课、奇门时家局
func HuCanDivine(t time.Time, question string, qt QuestionType) (*HuCanResult, error) {
	// 周易：铜钱法
	zy := DivineByCoins()
	zy.Time = t
	zy.QuestionType = qt

	// 大六壬
	lrPan, err := liuren.Divine(t)
	if err != nil {
		return nil, fmt.Errorf("六壬起课失败: %w", err)
	}

	// 奇门遁甲
	qmPan, err := qimen.BuildPan(t)
	if err != nil {
		return nil, fmt.Errorf("奇门起局失败: %w", err)
	}

	return &HuCanResult{
		Zhouyi:    zy,
		LiuRenPan: lrPan,
		QimenPan:  qmPan,
		Question:  question,
		QType:     qt,
	}, nil
}

// 八卦 → 五行（用于与六壬日干/三传的五行对比）
var trigramWuXing = map[string]liuren.WuXing{
	"乾": liuren.Jin, "兑": liuren.Jin,
	"震": liuren.Mu, "巽": liuren.Mu,
	"坤": liuren.Tu, "艮": liuren.Tu,
	"坎": liuren.Shui,
	"离": liuren.Huo,
}

// wuXingRelation 返回两五行关系的通俗描述
func wuXingRelation(a, b liuren.WuXing) string {
	switch {
	case a == b:
		return "比和（同气相求）"
	case liuren.WuXingGenerates(a, b):
		return fmt.Sprintf("%s 生 %s（相生扶助）", a, b)
	case liuren.WuXingGenerates(b, a):
		return fmt.Sprintf("%s 生 %s（反哺）", b, a)
	case liuren.WuXingOvercomes(a, b):
		return fmt.Sprintf("%s 克 %s（相战相制）", a, b)
	case liuren.WuXingOvercomes(b, a):
		return fmt.Sprintf("%s 克 %s（被制）", b, a)
	}
	return "无直接生克"
}

// hexOverallBias 粗略判断卦象吉凶倾向：变卦之卦辞、变爻数、关键吉凶字眼合取
// 返回 +1 偏吉 / 0 平 / -1 偏凶。仅做互参信号用，非严格断语。
func hexOverallBias(r DivinationResult) int {
	if r.MainHex == nil {
		return 0
	}
	score := 0
	words := r.MainHex.Judgment
	if r.ChangeHex != nil {
		words += r.ChangeHex.Judgment
	}
	positive := []string{"元亨", "利贞", "大吉", "吉", "无咎", "利涉大川", "利见大人", "亨"}
	negative := []string{"凶", "厉", "悔", "吝", "不利", "无攸利"}
	for _, k := range positive {
		if strings.Contains(words, k) {
			score++
		}
	}
	for _, k := range negative {
		if strings.Contains(words, k) {
			score--
		}
	}
	switch {
	case score >= 2:
		return 1
	case score <= -2:
		return -1
	}
	return 0
}

// panOverallBias 依末传六亲、三传空亡、课体粗断课势
func panOverallBias(p *liuren.Pan) int {
	if p == nil {
		return 0
	}
	score := 0
	// 末传为吉神（青龙/六合/太常/天后/贵人），凶神（白虎/玄武/勾陈/螣蛇）
	tj := p.SanChuan.Mo.TianJiang.String()
	switch tj {
	case "青龙", "六合", "太常", "天后", "贵人":
		score++
	case "白虎", "玄武", "勾陈", "螣蛇":
		score--
	}
	// 末传空亡 → 事无终
	if p.SanChuan.Mo.IsKong {
		score--
	}
	// 三传全不空、元首课/重审等常见吉课体名中含"元首/重审"
	if strings.Contains(p.KeTi.Name, "元首") || strings.Contains(p.KeTi.Name, "重审") {
		score++
	}
	if strings.Contains(p.KeTi.Name, "涉害") || strings.Contains(p.KeTi.Name, "伏吟") || strings.Contains(p.KeTi.Name, "返吟") {
		score--
	}
	switch {
	case score >= 2:
		return 1
	case score <= -2:
		return -1
	}
	return score
}

func biasLabel(b int) string {
	switch {
	case b > 0:
		return "偏吉"
	case b < 0:
		return "偏凶"
	}
	return "中平"
}

// qimenOverallBias 根据奇门盘的命中格局 + 值符/值使落宫情况粗断吉凶倾向
func qimenOverallBias(p *qimen.Pan) int {
	if p == nil {
		return 0
	}
	score := 0
	// 格局吉凶分累积
	for _, h := range qimen.DetectPatterns(p) {
		score += h.AuspiceScore
	}
	// 值符值使落宫的"门迫/空亡/入墓"进一步压分
	for _, c := range p.Cells {
		if c.PalaceName == p.ZhiFuPalace || c.PalaceName == p.ZhiShiPalace {
			if c.IsVoid {
				score--
			}
			if c.IsDoorPo {
				score--
			}
			if c.IsJiXing {
				score--
			}
			if c.IsDoorSheng {
				score++
			}
		}
	}
	switch {
	case score >= 3:
		return 1
	case score <= -3:
		return -1
	}
	if score > 0 {
		return 1
	}
	if score < 0 {
		return -1
	}
	return 0
}

// extractHuCanSignals 从本盘抽取三式动态互参信号（非通用口诀）
func extractHuCanSignals(r *HuCanResult) string {
	var b strings.Builder

	// —— 1. 三系吉凶倾向对照（矩阵判断）
	zyBias := hexOverallBias(r.Zhouyi)
	lrBias := panOverallBias(r.LiuRenPan)
	qmBias := qimenOverallBias(r.QimenPan)
	b.WriteString(fmt.Sprintf("- **三系吉凶对照**：周易%s ｜ 大六壬%s ｜ 奇门%s。",
		biasLabel(zyBias), biasLabel(lrBias), biasLabel(qmBias)))
	b.WriteString(threeSystemVerdict(zyBias, lrBias, qmBias))
	b.WriteString("\n")

	// —— 2. 五行相与：周易本卦（上/下卦）vs 六壬日干
	if r.Zhouyi.MainHex != nil && r.LiuRenPan != nil {
		upper := r.Zhouyi.MainHex.Upper
		lower := r.Zhouyi.MainHex.Lower
		wu, upOK := trigramWuXing[upper]
		wl, loOK := trigramWuXing[lower]
		ganWX := liuren.GanWuXing[r.LiuRenPan.Ctx.Gan]
		if upOK && loOK {
			b.WriteString(fmt.Sprintf("- **五行相与**：本卦上卦 %s（%s）· 下卦 %s（%s）；日干 %s 属 %s。",
				upper, wu, lower, wl, r.LiuRenPan.Ctx.Gan, ganWX))
			b.WriteString("上卦为外象/他方 → " + wuXingRelation(wu, ganWX) + "；")
			b.WriteString("下卦为内象/自身 → " + wuXingRelation(wl, ganWX) + "。\n")
		}
	}

	// —— 3. 变卦方向 vs 末传趋势：卦由 A 变 B，看 A→B 五行方向；末传为事之终局
	if r.Zhouyi.ChangeHex != nil && r.LiuRenPan != nil {
		mainWX, ok1 := trigramWuXing[r.Zhouyi.MainHex.Upper]
		chgWX, ok2 := trigramWuXing[r.Zhouyi.ChangeHex.Upper]
		moWX := liuren.ZhiWuXing[r.LiuRenPan.SanChuan.Mo.Zhi]
		if ok1 && ok2 {
			b.WriteString(fmt.Sprintf("- **趋势共振**：本卦上卦 %s → 变卦上卦 %s（%s→%s）；六壬末传 %s 属 %s。",
				r.Zhouyi.MainHex.Upper, r.Zhouyi.ChangeHex.Upper, mainWX, chgWX,
				r.LiuRenPan.SanChuan.Mo.Zhi, moWX))
			// 方向一致：变卦五行 与 末传五行 生克方向同
			if chgWX == moWX {
				b.WriteString("**同归一象 → 结局明确**。\n")
			} else if liuren.WuXingGenerates(chgWX, moWX) || liuren.WuXingGenerates(moWX, chgWX) {
				b.WriteString("两者相生 → 结局虽非同物而能**相辅相承**。\n")
			} else if liuren.WuXingOvercomes(chgWX, moWX) || liuren.WuXingOvercomes(moWX, chgWX) {
				b.WriteString("两者相战 → 卦壬所示之终局**有别**，须辨何者为真。\n")
			} else {
				b.WriteString("各行其道，可参酌月将、类神定取舍。\n")
			}
		}
	} else if len(r.Zhouyi.ChangingPos) == 0 {
		b.WriteString("- **趋势共振**：周易无变爻，以本卦之象为定局；六壬末传仍在演化，**六壬主动、周易主静**，以六壬末传所临之神为事之最终归趋。\n")
	}

	// —— 4. 变爻位置 vs 三传关键位
	if len(r.Zhouyi.ChangingPos) > 0 {
		b.WriteString(fmt.Sprintf("- **关键位共振**：周易变爻落在第 %v 爻（", r.Zhouyi.ChangingPos))
		posMeanings := map[int]string{
			1: "事之始、根基", 2: "臣/内、中正", 3: "过渡、危地",
			4: "近君/外、谨慎", 5: "君/主、权位", 6: "事之终、过亢",
		}
		var names []string
		for _, p := range r.Zhouyi.ChangingPos {
			if m, ok := posMeanings[p]; ok {
				names = append(names, fmt.Sprintf("%d爻=%s", p, m))
			}
		}
		b.WriteString(strings.Join(names, "；"))
		b.WriteString(fmt.Sprintf("）。对照六壬三传：初传 %s、中传 %s、末传 %s —— 取变爻之位义与三传之发展次第相互印证。\n",
			r.LiuRenPan.SanChuan.Chu.Zhi, r.LiuRenPan.SanChuan.Zhong.Zhi, r.LiuRenPan.SanChuan.Mo.Zhi))
	}

	// —— 5. 空亡信号
	anyKong := r.LiuRenPan.SanChuan.Chu.IsKong || r.LiuRenPan.SanChuan.Zhong.IsKong || r.LiuRenPan.SanChuan.Mo.IsKong
	if anyKong {
		var kongs []string
		if r.LiuRenPan.SanChuan.Chu.IsKong {
			kongs = append(kongs, "初传空")
		}
		if r.LiuRenPan.SanChuan.Zhong.IsKong {
			kongs = append(kongs, "中传空")
		}
		if r.LiuRenPan.SanChuan.Mo.IsKong {
			kongs = append(kongs, "末传空")
		}
		b.WriteString("- **空亡警示**：六壬三传中 " + strings.Join(kongs, "、") +
			"，若周易变爻所指之事与空传所主相合 → 虚浮不实、事易中辍；若无关，则空亡只减损其一侧之信度，另一侧仍可据。\n")
	}

	// —— 6. 贵人 / 昼夜与周易阳刚
	ganBias := liuren.GuiRenByGan[r.LiuRenPan.Ctx.Gan]
	zy := "夜占"
	gui := ganBias[1]
	if r.LiuRenPan.Ctx.ZhouYe {
		zy = "昼占"
		gui = ganBias[0]
	}
	yangCount := 0
	for _, v := range r.Zhouyi.Lines {
		if v == 7 || v == 9 {
			yangCount++
		}
	}
	b.WriteString(fmt.Sprintf("- **昼夜与刚柔**：占时属 %s，昼夜贵人为 %s；周易本卦阳爻 %d / 阴爻 %d。",
		zy, gui, yangCount, 6-yangCount))
	if r.LiuRenPan.Ctx.ZhouYe && yangCount >= 4 {
		b.WriteString("昼占而卦阳盛 → **天时人势俱顺**，宜主动出击。\n")
	} else if !r.LiuRenPan.Ctx.ZhouYe && yangCount <= 2 {
		b.WriteString("夜占而卦阴盛 → **静守藏器**，不宜张扬。\n")
	} else {
		b.WriteString("刚柔与昼夜不相应 → 于此察「天时与人谋」之分歧。\n")
	}

	// —— 7. 奇门命中格局摘要（跨系核心信号）
	if r.QimenPan != nil {
		hits := qimen.DetectPatterns(r.QimenPan)
		if len(hits) > 0 {
			var jiNames, xiongNames []string
			for _, h := range hits {
				if h.AuspiceScore >= 1 {
					jiNames = append(jiNames, h.Name)
				} else if h.AuspiceScore <= -1 {
					xiongNames = append(xiongNames, h.Name)
				}
			}
			parts := []string{}
			if len(jiNames) > 0 {
				parts = append(parts, "吉格 "+strings.Join(jiNames, "、"))
			}
			if len(xiongNames) > 0 {
				parts = append(parts, "凶格 "+strings.Join(xiongNames, "、"))
			}
			if len(parts) > 0 {
				b.WriteString("- **奇门格局快览**：" + strings.Join(parts, " ｜ ") + " —— 奇门所示之方位时机线索。\n")
			}
		}
	}

	// —— 8. 奇门日干落宫 vs 六壬日上神 vs 周易内卦（三系"我"之对照）
	if r.QimenPan != nil {
		qmDayPal, qmDayState := qimenDayGanLocation(r.QimenPan)
		lrDayUp := lrDayUpperState(r.LiuRenPan)
		zyInner := ""
		if r.Zhouyi.MainHex != nil {
			zyInner = r.Zhouyi.MainHex.Lower
		}
		b.WriteString(fmt.Sprintf("- **三系「我」之对照**：奇门日干 %s 落 %s（%s）｜ 六壬日上神 %s ｜ 周易下卦 %s。",
			r.QimenPan.Ctx.DayGan, qmDayPal, qmDayState, lrDayUp, zyInner))
		b.WriteString("三象同明者信度最高；若一系示吉而他系示凶，须辨表里。\n")
	}

	// —— 9. 奇门值使落宫方位 vs 六壬末传地支（方位时机对照）
	if r.QimenPan != nil {
		zsPal := r.QimenPan.ZhiShiPalace
		zsDir := palaceDirection(zsPal)
		moZhi := r.LiuRenPan.SanChuan.Mo.Zhi.String()
		b.WriteString(fmt.Sprintf("- **方位时机**：奇门值使 %s 落 %s（方位%s）｜ 六壬末传 %s。——奇门主「去哪、何时动」，以值使宫位为用事方向；六壬末传为事之归宿时。\n",
			r.QimenPan.ZhiShiGate, zsPal, zsDir, moZhi))
	}

	return b.String()
}

// threeSystemVerdict 给出三系 bias 的交叉裁决文字
func threeSystemVerdict(zy, lr, qm int) string {
	allAgree := func(b int) bool { return zy == b && lr == b && qm == b }
	countBy := func(b int) int {
		c := 0
		if zy == b {
			c++
		}
		if lr == b {
			c++
		}
		if qm == b {
			c++
		}
		return c
	}
	switch {
	case allAgree(1):
		return " **三系皆吉 · 信度最高**——可放心依判，大吉之机。"
	case allAgree(-1):
		return " **三系皆凶 · 信度最高**——须严守戒备，勿强为。"
	case countBy(1) == 2:
		// 两吉一凶
		return " **两系偏吉、一系偏凶**：主势向吉，但有一隐伏之患——注意凶侧所指为事中之暗礁。"
	case countBy(-1) == 2:
		// 两凶一吉
		return " **两系偏凶、一系偏吉**：主势向凶，唯余一线生机——抓住吉侧所示之转机方可化解。"
	case countBy(1) == 1 && countBy(-1) == 1:
		return " **一吉一凶一平**：结局未定，取舍之间——以问题类型的主用神（类神）所临系统为准判。"
	}
	return " 三系倾向不明 → 以本盘关键信号与类神直指为准。"
}

// qimenDayGanLocation 返回 (宫位名, 状态描述)
func qimenDayGanLocation(p *qimen.Pan) (string, string) {
	if p == nil || p.Ctx == nil {
		return "", ""
	}
	dg := p.Ctx.DayGan
	lookup := dg
	if dg == "甲" {
		lookup = p.Ctx.Dungan
	}
	for _, c := range p.Cells {
		if c.HeavenStem == lookup || c.EarthStem == lookup {
			parts := []string{}
			if c.Star != "" {
				parts = append(parts, c.Star+c.StarWangShuai)
			}
			if c.Door != "" {
				parts = append(parts, c.Door)
			}
			if c.God != "" {
				parts = append(parts, c.God)
			}
			if c.IsDoorPo {
				parts = append(parts, "门迫")
			}
			if c.IsVoid {
				parts = append(parts, "空亡")
			}
			state := "乘 " + strings.Join(parts, "·")
			return c.PalaceName, state
		}
	}
	return "", ""
}

// lrDayUpperState 六壬日上神及其所乘天将
func lrDayUpperState(p *liuren.Pan) string {
	if p == nil {
		return ""
	}
	gong := liuren.GanJiGong[p.Ctx.Gan]
	upper := p.TianPan[gong]
	tj := p.TianJiang[gong]
	return fmt.Sprintf("%s(乘%s)", upper, tj)
}

// palaceDirection 宫位名 → 后天方位
func palaceDirection(palace string) string {
	switch palace {
	case "坎一宫":
		return "正北"
	case "坤二宫":
		return "西南"
	case "震三宫":
		return "正东"
	case "巽四宫":
		return "东南"
	case "中五宫":
		return "中央"
	case "乾六宫":
		return "西北"
	case "兑七宫":
		return "正西"
	case "艮八宫":
		return "东北"
	case "离九宫":
		return "正南"
	}
	return ""
}

// HuCanPrompt 合并三个系统的完整素材 + 本盘互参信号，作为给 AI 的断占提示词
func HuCanPrompt(r *HuCanResult) string {
	var b strings.Builder
	b.WriteString("你是一位兼通**周易六爻**与**大六壬**、**奇门遁甲**等传统术数的资深顾问。\n")
	b.WriteString("以下同一时刻分别起周易卦、大六壬课、奇门时家局——请将**三副盘面视为一体**，互为印证、互相补全，对所问之事给出深入的互参解读。\n")
	b.WriteString("**重要注意**：\n")
	b.WriteString("1. 以下素材为给你的原始数据，请直接**解读其含义**，不要大段复述素材原文（关键字句可简短引用）。\n")
	b.WriteString("2. 每个结论请**标注依据**，形如「（依据：本卦变爻九五 + 末传乘青龙 + 奇门值符临坎）」。\n")
	b.WriteString("3. 素材中「关键信号」「类神直指」「命中格局」「三系互参信号」为代码预抽取的本盘结论，可直接采用；其余为支撑材料。\n")
	b.WriteString("4. **三式分工**：周易断「事理象义」、六壬断「人事细节」、奇门断「方位时机策略」——各擅胜场，不可混判。\n\n")

	if r.Question != "" {
		b.WriteString("## 所问之事\n" + r.Question + "\n\n")
	}

	// === 第一部分：完整周易盘 ===
	b.WriteString("# 第一部分 · 周易卦象（完整盘）\n\n")
	zyCopy := r.Zhouyi
	zyCopy.QuestionType = ""
	zyPrompt := GenerateAIPrompt(zyCopy, "")
	zyBody := stripPromptHeader(zyPrompt)
	zyBody = trimAtFirst(zyBody, "## 问题类型", "## 解卦侧重", "## 请按以下结构解卦")
	b.WriteString(zyBody)
	b.WriteString("\n")

	// === 第二部分：完整大六壬盘 ===
	b.WriteString("# 第二部分 · 大六壬盘面（完整盘）\n\n")
	leishen := LeiShenDirective(r.LiuRenPan, r.QType)
	lrPrompt := liuren.GenerateAIPrompt(r.LiuRenPan, "", "", leishen)
	lrBody := stripPromptHeader(lrPrompt)
	lrBody = trimAtFirst(lrBody, "## 解课侧重", "## 请按以下结构断课")
	b.WriteString(lrBody)
	b.WriteString("\n")

	// === 第三部分：完整奇门遁甲盘 ===
	b.WriteString("# 第三部分 · 奇门遁甲盘面（完整盘）\n\n")
	qmPrompt := qimen.GenerateAIPrompt(r.QimenPan, "", "", string(r.QType))
	qmBody := stripPromptHeader(qmPrompt)
	qmBody = trimAtFirst(qmBody, "## 解局侧重", "## 请按以下结构解局")
	b.WriteString(qmBody)
	b.WriteString("\n")

	// === 第四部分：三系本盘动态互参信号 ===
	b.WriteString("# 第四部分 · 三系本盘动态互参信号（由三副盘面自动抽取）\n")
	b.WriteString("以下信号为本次起占专属的跨系对照点，非通用口诀；请据此判定三系一致点（信度）与分歧点（取舍）。\n\n")
	b.WriteString(extractHuCanSignals(r))
	b.WriteString("\n")

	// === 第五部分：三式互参理法（通用原则） ===
	b.WriteString("# 第五部分 · 三式互参理法（通用原则）\n")
	b.WriteString("- **三式分工**：\n")
	b.WriteString("  - **周易**主「**象**」——揭示事理之微、人心之向（为什么 / 是什么）\n")
	b.WriteString("  - **大六壬**主「**式**」——揭示人事之动、类神得失（谁在动 / 怎么牵连）\n")
	b.WriteString("  - **奇门**主「**机**」——揭示方位时机、谋略布局（去哪 / 何时 / 怎么布局）\n")
	b.WriteString("- **内外之分**：周易下卦为我、上卦为他；六壬日干为我、日支为事；奇门日干为我、时干为所问。\n")
	b.WriteString("- **静动之分**：周易无变爻主静（定局），有变爻主动；六壬三传永远在动，末传即事之归宿；奇门值符值使所临定事之行止。\n")
	b.WriteString("- **信度分级**（三系 bias 对照）：\n")
	b.WriteString("  - 三系同指 → **信度最高**，可直判\n")
	b.WriteString("  - 两系同指、一系异 → **信度中等**，异系所指之方向为暗礁或转机\n")
	b.WriteString("  - 三系皆异 → **结论宜保守**，以类神直指的主用神所临系统为裁决\n")
	b.WriteString("- **应期合参**：周易变爻干支 + 六壬末传所值 + 奇门值使所临地支——三者合一之时应验最准。\n\n")

	// 问题类型侧重（如有）
	if r.QType != "" && r.QType != QTOther {
		b.WriteString("# 第六部分 · 问题类型侧重\n")
		b.WriteString(fmt.Sprintf("问题类型：%s\n\n", QuestionTypeLabel(r.QType)))
		b.WriteString("### 周易 · 解卦侧重\n")
		b.WriteString(FocusGuide(r.QType))
		b.WriteString("\n")
		b.WriteString("### 六壬 · 断课侧重\n")
		b.WriteString(FocusGuideLiuRen(r.QType))
		b.WriteString("\n")
		b.WriteString("### 奇门 · 解局侧重\n")
		b.WriteString(qimen.FocusGuide(string(r.QType)))
		b.WriteString("\n")
	}

	// === 最终结构化解读框架（7 步，强化三系合参） ===
	b.WriteString("# 请按以下结构给出三式互参解答\n")
	b.WriteString("1. **三系总论**（各 1-2 句）：分别用周易象、六壬式、奇门机三句话各自概括对所问之事的基本判断（不要混谈）。\n")
	b.WriteString("2. **三系一致点**（信度最高）：列出三系共同指向的 2-4 项结论，每条说明依据三系各自的证据（对照第四部分动态信号）。\n")
	b.WriteString("3. **分歧点**（需辨表里）：若有系统之间吉凶/方向/应期不一致，指出分歧所在，结合「体/式/机」分工判断何者为主象、何者为辅。\n")
	b.WriteString("4. **关键信号解析**：从变爻、三传、值符值使、格局命中中挑 2-3 个最决定性的集中展开，点出「决定性证据」。\n")
	b.WriteString("5. **时机与方位**：\n")
	b.WriteString("   - **时机**：综合周易变爻月令、六壬末传值神、奇门值使所临地支，给出发动/转折/了结的时间节点\n")
	b.WriteString("   - **方位**：主要取奇门——指出对求测者最有利的后天方位（基于值使 / 三吉门 / 三奇落宫）\n")
	b.WriteString("6. **行动建议**（3-5 条）：综合三系结论，分「宜」与「忌」，结合奇门命中格局给出具体布局。\n")
	b.WriteString("7. **一句断语** + **三系信度自评**：\n")
	b.WriteString("   - 一句断语给最终判断（吉/凶/平；宜/忌何事）\n")
	b.WriteString("   - 信度分三档：高（三系同指）/ 中（两系同指）/ 低（三系分歧）\n\n")
	b.WriteString("请用中文回答，兼顾传统术语（爻位中正应比、涉害见机、贵人乘将、类神所临、三奇得使、门迫入墓）与现代通俗表达；**重解读轻复述**。\n")
	return b.String()
}

// stripPromptHeader 剥去子提示词顶部的"你是XX…"+重要注意块+所问之事，
// 返回从第一个"行首以 '## ' 开头"的小节开始的剩余内容。
func stripPromptHeader(s string) string {
	if i := strings.Index(s, "\n## "); i >= 0 {
		return s[i+1:]
	}
	if strings.HasPrefix(s, "## ") {
		return s
	}
	return s
}

// trimAtFirst 在 s 中查找 markers 里最早出现的一个并截断到它之前（不含 marker）。
// 若都未出现则返回原串。
func trimAtFirst(s string, markers ...string) string {
	cut := -1
	for _, m := range markers {
		if i := strings.Index(s, m); i >= 0 && (cut == -1 || i < cut) {
			cut = i
		}
	}
	if cut >= 0 {
		return s[:cut]
	}
	return s
}

// HuCanText 终端风格的互参摘要（含奇门）
func HuCanText(r *HuCanResult) string {
	var b strings.Builder
	b.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	b.WriteString("◎ 周易 × 大六壬 × 奇门遁甲 三式互参\n")
	b.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	b.WriteString("  时刻：" + r.Zhouyi.Time.Format("2006-01-02 15:04") + "\n")
	if r.LiuRenPan != nil {
		b.WriteString("  " + r.LiuRenPan.Ctx.Summary() + "\n")
	}
	if r.QimenPan != nil {
		b.WriteString("  " + r.QimenPan.Ctx.Summary() + "\n")
	}
	b.WriteString("\n")
	b.WriteString(InterpretResult(r.Zhouyi))
	b.WriteString("\n\n")
	b.WriteString(liuren.Render(r.LiuRenPan))
	b.WriteString("\n\n")
	if r.QimenPan != nil {
		b.WriteString(r.QimenPan.Render())
	}
	return b.String()
}
