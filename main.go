package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"zhouyi/liuren"
)

const banner = `
╔═══════════════════════════════════════════════════════════╗
║              周 易 算 卦 · 易经占卜系统                    ║
║                                                           ║
║   "易有太极，是生两仪，两仪生四象，四象生八卦"              ║
║                          ——《周易·系辞上》                 ║
╚═══════════════════════════════════════════════════════════╝
`

const help = `
【周易算卦使用说明】

▌ 起卦方式：
  1. 铜钱法   — 模拟三枚铜钱投掷六次，最接近传统的简易方法
  2. 蓍草法   — 模拟古法大衍揲蓍，精确还原《系辞》所载方法
  3. 数字起卦 — 输入两个数字（上卦/下卦）及变爻位置
  4. 查卦     — 直接查询某一卦的卦辞爻辞
  5. 大六壬   — 以月将加时起课，出四课三传十二天将，三式之一，论人事最精

▌ 命令：
  coin          铜钱法起卦
  yarrow        蓍草法起卦
  number        数字起卦
  liuren        大六壬起课（以当前时刻）
  huican        周易 × 大六壬互参占
  lookup <编号> 查询指定卦（1-64）
  list          列出全部六十四卦
  prompt        重新显示上一次起卦的 AI 解卦提示词
  help          显示此帮助
  quit / exit   退出

▌ 爻的符号说明：
  ───     阳爻（少阳，不变）
  ── ──   阴爻（少阴，不变）
  ─○─     老阳（变爻，将变为阴爻）
  ─×─     老阴（变爻，将变为阳爻）

▌ 解卦要诀（朱熹《易学启蒙·考变占》）：
  • 无变爻：以本卦卦辞（彖辞）断
  • 一爻变：以本卦该变爻爻辞断
  • 二爻变：以本卦两变爻爻辞断，仍以上爻为主
  • 三爻变：占本卦与变卦之卦辞，本卦为贞（主），变卦为悔（次）
  • 四爻变：以变卦中两个不变爻之爻辞断，仍以下爻为主
  • 五爻变：以变卦中唯一不变之爻辞断
  • 六爻皆变：乾用"用九"、坤用"用六"，余卦占变卦卦辞
`

// lastResult 保存上一次起卦结果，供 prompt 命令重新显示
var lastResult *DivinationResult
var lastQuestion string

func main() {
	// 子命令：cast 非交互式起盘（供 agent / skill 调用，输出 JSON）
	if len(os.Args) > 1 && os.Args[1] == "cast" {
		os.Exit(runCast(os.Args[2:]))
	}

	// 子命令：serve 启动 HTTP 服务
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		addr := ":8080"
		if len(os.Args) > 2 {
			addr = os.Args[2]
			if !strings.Contains(addr, ":") {
				addr = ":" + addr
			}
		}
		fmt.Print(banner)
		runServer(addr)
		return
	}

	fmt.Print(banner)
	fmt.Println("  输入 help 查看使用说明，输入 quit 退出")
	fmt.Println("  （要启动 Web 服务，请使用： zhouyi serve [port]）")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("周易> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		cmd := strings.ToLower(parts[0])

		switch cmd {
		case "quit", "exit", "q":
			fmt.Println("\n  阴阳流转，生生不息。再会！")
			return

		case "help", "h", "?":
			fmt.Print(help)

		case "coin", "铜钱", "金钱":
			doCoinDivination(reader)

		case "yarrow", "蓍草":
			doYarrowDivination(reader)

		case "number", "数字":
			doNumberDivination(reader)

		case "liuren", "六壬", "大六壬":
			doLiuRenDivination(reader)

		case "huican", "互参":
			doHuCanDivination(reader)

		case "prompt", "提示词", "ai":
			if lastResult == nil {
				fmt.Println("  尚未起卦，请先使用 coin / yarrow / number 起卦")
			} else {
				PrintAIPrompt(*lastResult, lastQuestion)
			}

		case "lookup", "查", "查卦":
			if len(parts) < 2 {
				fmt.Println("  用法：lookup <卦号1-64>")
				continue
			}
			n, e := strconv.Atoi(parts[1])
			if e != nil || n < 1 || n > 64 {
				fmt.Println("  请输入 1-64 的卦号")
				continue
			}
			printHexagramDetail(&Hexagrams[n-1])

		case "list", "列表":
			printHexagramList()

		default:
			// 尝试解析为数字，直接查卦
			if n, e := strconv.Atoi(input); e == nil && n >= 1 && n <= 64 {
				printHexagramDetail(&Hexagrams[n-1])
			} else {
				fmt.Println("  未知命令，输入 help 查看帮助")
			}
		}
	}
}

func doCoinDivination(reader *bufio.Reader) {
	fmt.Println("\n  ◎ 铜钱法起卦")
	question, qtype := askQuestion(reader)
	fmt.Println("  请凝神静气，心中默念所问之事，片刻后按 Enter 开始摇卦...")
	reader.ReadString('\n')

	fmt.Print("  正在摇卦")
	for i := 0; i < 6; i++ {
		time.Sleep(200 * time.Millisecond)
		fmt.Print(".")
	}
	fmt.Println()

	result := DivineByCoins()
	result.QuestionType = qtype
	printDivinationProcess(result, "铜钱法")
	fmt.Println(InterpretResult(result))
	saveAndOfferPrompt(result, question, reader)
}

func doYarrowDivination(reader *bufio.Reader) {
	fmt.Println("\n  ◎ 蓍草法（大衍揲蓍）起卦")
	fmt.Println("  古法以五十根蓍草，用四十九根，经十八变成卦")
	question, qtype := askQuestion(reader)
	fmt.Println("  请凝神静气，心中默念所问之事，片刻后按 Enter 开始...")
	reader.ReadString('\n')

	fmt.Print("  十八变揲蓍中")
	for i := 0; i < 18; i++ {
		time.Sleep(100 * time.Millisecond)
		if i%3 == 2 {
			fmt.Print("·")
		}
	}
	fmt.Println()

	result := DivineByYarrow()
	result.QuestionType = qtype
	printDivinationProcess(result, "蓍草法")
	fmt.Println(InterpretResult(result))
	saveAndOfferPrompt(result, question, reader)
}

func doNumberDivination(reader *bufio.Reader) {
	fmt.Println("\n  ◎ 数字起卦法")
	fmt.Println("  请依次输入：上卦数字、下卦数字（任意整数均可，系统自动取模）")
	fmt.Println("  变爻位置（1-6，输入0则无变爻）")
	fmt.Println()

	question, qtype := askQuestion(reader)

	fmt.Print("  上卦数字: ")
	s1, _ := reader.ReadString('\n')
	upper, err1 := strconv.Atoi(strings.TrimSpace(s1))

	fmt.Print("  下卦数字: ")
	s2, _ := reader.ReadString('\n')
	lower, err2 := strconv.Atoi(strings.TrimSpace(s2))

	fmt.Print("  变爻位置(0=无): ")
	s3, _ := reader.ReadString('\n')
	changing, err3 := strconv.Atoi(strings.TrimSpace(s3))

	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println("  输入有误，请输入整数")
		return
	}

	result := DivineByNumber(upper, lower, changing)
	result.QuestionType = qtype
	printDivinationProcess(result, "数字法")
	fmt.Println(InterpretResult(result))
	saveAndOfferPrompt(result, question, reader)
}

// askQuestion 询问所问之事与问题类型，允许留空
func askQuestion(reader *bufio.Reader) (string, QuestionType) {
	fmt.Print("  所问之事（可直接回车跳过）: ")
	s, _ := reader.ReadString('\n')
	q := strings.TrimSpace(s)
	if q != "" {
		fmt.Printf("  已记录：%s\n", q)
	}

	fmt.Print(QuestionTypeMenu())
	fmt.Print("  选择: ")
	t, _ := reader.ReadString('\n')
	qt := ParseQuestionType(t)
	if strings.TrimSpace(t) != "" {
		fmt.Printf("  已记录问题类型：%s\n", QuestionTypeLabel(qt))
	}
	fmt.Println()
	return q, qt
}

// saveAndOfferPrompt 保存结果并询问是否生成 AI 提示词
func saveAndOfferPrompt(result DivinationResult, question string, reader *bufio.Reader) {
	lastResult = &result
	lastQuestion = question

	fmt.Print("\n  是否生成 AI 解卦提示词？(y/n，默认 y): ")
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(strings.ToLower(ans))
	if ans == "" || ans == "y" || ans == "yes" {
		PrintAIPrompt(result, question)
		offerLLMInterpret(reader, GenerateAIPrompt(result, question))
	} else {
		fmt.Println("  （可随时输入 prompt 命令重新显示提示词）")
	}
}

// offerLLMInterpret 在已生成提示词后，询问并调用大模型 API 直接给出解读。
// 仅当 config.json 的 llm 段已配置（apiKey 非空）时才提示；否则给一句温和的引导。
// prompt 为对应术数已生成好的解卦提示词文本。
func offerLLMInterpret(reader *bufio.Reader, prompt string) {
	if strings.TrimSpace(prompt) == "" {
		return
	}
	cfg := LoadConfig().LLM
	_, usable := cfg.resolved()
	if !usable {
		fmt.Println("\n  （提示：在 config.json 的 llm 段填入 apiKey 并设 enabled=true，即可让程序直接调用 AI 解卦）")
		return
	}

	fmt.Print("\n  是否直接调用 AI 解卦？(y/n，默认 y): ")
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(strings.ToLower(ans))
	if ans != "" && ans != "y" && ans != "yes" {
		return
	}

	fmt.Println("\n  正在请大模型解卦，请稍候……")
	text, err := Interpret(context.Background(), cfg, prompt)
	if err != nil {
		fmt.Printf("  解卦失败：%v\n", err)
		fmt.Println("  （可复制上面的提示词到任意 AI 自行解读）")
		return
	}
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    AI 解卦结果                              ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println(text)
	fmt.Println("═══════════════════════════════════════════════════════════")
}

func printDivinationProcess(r DivinationResult, method string) {
	lineNames := []string{"初", "二", "三", "四", "五", "上"}
	fmt.Printf("\n  【%s · 六爻详情】\n", method)
	fmt.Println("  爻位   投掷值   类型        爻形")
	fmt.Println("  ─────────────────────────────────")
	for i := 5; i >= 0; i-- {
		li := lineInfo(r.Lines[i])
		yy := "九"
		if !li.IsYang {
			yy = "六"
		}
		typeName := map[int]string{
			6: "老阴（变）",
			7: "少阳",
			8: "少阴",
			9: "老阳（变）",
		}[r.Lines[i]]
		change := ""
		if li.IsChange {
			change = " ★"
		}
		fmt.Printf("  %s%s  [%d]   %-10s  %s%s\n",
			lineNames[i], yy, r.Lines[i], typeName, li.Symbol, change)
	}
	fmt.Println()
}

func printHexagramDetail(h *Hexagram) {
	fmt.Printf("\n  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  第%d卦  %s  %s卦\n", h.Number, h.Symbol, h.Name)
	fmt.Printf("  上卦：%s %s（%s）  下卦：%s %s（%s）\n",
		Trigrams[h.Upper].Symbol, h.Upper, Trigrams[h.Upper].Nature,
		Trigrams[h.Lower].Symbol, h.Lower, Trigrams[h.Lower].Nature)
	fmt.Printf("  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("  卦辞：%s\n\n", h.Judgment)
	fmt.Printf("  象辞：%s\n\n", h.Image)
	fmt.Printf("  【爻辞】\n")
	for i := 5; i >= 0; i-- {
		line := h.Lines[i]
		fmt.Printf("  %s：%s\n", line.Type, line.Text)
	}
	fmt.Println()
}

func printHexagramList() {
	fmt.Println("\n  【六十四卦一览表】")
	fmt.Println("  ┌──────────────────────────────────────────────────────────────┐")
	for i, h := range Hexagrams {
		upper := Trigrams[h.Upper]
		lower := Trigrams[h.Lower]
		fmt.Printf("  │ %2d. %s %-4s  %s%s/%s%s", h.Number, h.Symbol, h.Name,
			upper.Symbol, h.Upper, lower.Symbol, h.Lower)
		if (i+1)%2 == 0 {
			fmt.Println()
		} else {
			fmt.Printf("    ")
		}
	}
	if len(Hexagrams)%2 != 0 {
		fmt.Println()
	}
	fmt.Println("  └──────────────────────────────────────────────────────────────┘")
	fmt.Println("  输入 lookup <编号> 查看详情")
	fmt.Println()
}

// doLiuRenDivination 大六壬起课交互
func doLiuRenDivination(reader *bufio.Reader) {
	fmt.Println("\n  ◎ 大六壬起课（以当前时刻）")
	question, qtype := askQuestion(reader)

	// 二期：可选本命 / 出生年 / 性别
	fmt.Print("  本命生肖（鼠/牛/虎/.../猪，可回车跳过）: ")
	s, _ := reader.ReadString('\n')
	sxStr := strings.TrimSpace(s)
	var benMingZhi *liuren.Zhi
	if sxStr != "" {
		if z := parseShengXiao(sxStr); z >= 0 {
			benMingZhi = &z
		}
	}
	fmt.Print("  出生公历年份（用于推行年，可回车跳过）: ")
	byStr, _ := reader.ReadString('\n')
	birthYear, _ := strconv.Atoi(strings.TrimSpace(byStr))
	gender := ""
	if birthYear > 0 {
		fmt.Print("  性别（男/女，默认男）: ")
		gs, _ := reader.ReadString('\n')
		gender = strings.TrimSpace(gs)
		if gender == "" {
			gender = "男"
		}
	}

	fmt.Println("  请凝神静气，心中默念所问之事，片刻后按 Enter 开始起课...")
	reader.ReadString('\n')

	fmt.Print("  月将加时中")
	for i := 0; i < 6; i++ {
		time.Sleep(150 * time.Millisecond)
		fmt.Print(".")
	}
	fmt.Println()

	// 用上下文构造，以支持本命/行年
	ctx, err := liuren.BuildContext(time.Now())
	if err != nil {
		fmt.Printf("  起课失败：%v\n", err)
		return
	}
	if benMingZhi != nil {
		ctx.BenMing = benMingZhi
	}
	if birthYear > 0 {
		ctx.BirthYear = birthYear
		ctx.Gender = gender
	}
	ctx.QuestionType = string(qtype)
	pan := liuren.DivineWithContext(ctx)
	_ = pan
	// 后续保持原逻辑
	fmt.Println()
	fmt.Println(liuren.Render(pan))
	fmt.Println(liuren.InterpretGuide(pan))

	// 询问是否生成 AI 提示词
	fmt.Print("\n  是否生成 AI 解课提示词？(y/n，默认 y): ")
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(strings.ToLower(ans))
	if ans == "" || ans == "y" || ans == "yes" {
		focus := FocusGuideLiuRen(qtype)
		leishen := LeiShenDirective(pan, qtype)
		prompt := liuren.GenerateAIPrompt(pan, question, focus, leishen)
		fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
		fmt.Println("║                AI 断课提示词（可直接复制）                  ║")
		fmt.Println("╚═══════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println(prompt)
		fmt.Println("═══════════════════════════════════════════════════════════")
		offerLLMInterpret(reader, prompt)
	}
}

// doHuCanDivination 周易 × 大六壬互参占
func doHuCanDivination(reader *bufio.Reader) {
	fmt.Println("\n  ◎ 周易 × 大六壬 互参占")
	question, qtype := askQuestion(reader)
	fmt.Println("  请凝神静气，心中默念所问之事，片刻后按 Enter 开始...")
	reader.ReadString('\n')

	fmt.Print("  同时摇卦与起课")
	for i := 0; i < 6; i++ {
		time.Sleep(150 * time.Millisecond)
		fmt.Print(".")
	}
	fmt.Println()

	r, err := HuCanDivine(time.Now(), question, qtype)
	if err != nil {
		fmt.Printf("  互参失败：%v\n", err)
		return
	}
	fmt.Println()
	fmt.Println(HuCanText(r))

	fmt.Print("\n  是否生成 AI 互参提示词？(y/n，默认 y): ")
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(strings.ToLower(ans))
	if ans == "" || ans == "y" || ans == "yes" {
		prompt := HuCanPrompt(r)
		fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
		fmt.Println("║              AI 互参解读提示词（可直接复制）                ║")
		fmt.Println("╚═══════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println(prompt)
		fmt.Println("═══════════════════════════════════════════════════════════")
		offerLLMInterpret(reader, prompt)
	}
}
