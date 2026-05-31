package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"zhouyi/liuren"
	"zhouyi/qimen"
)

// cast 子命令：非交互式起盘，供外部 agent / skill 调用。
// 一次调用即起即灭，输出 JSON（含解卦提示词），不进入 REPL，不需鉴权。
//
// 用法：
//   zhouyi cast --method <zhouyi|qimen|liuren|huican> --question "所问之事" \
//               [--type career] [--coin] [--upper N --lower N --changing N] \
//               [--time RFC3339] [--lon 经度] [--mode prompt|full]
//
// --method 说明：
//   zhouyi  周易六爻（默认铜钱法；给 --upper/--lower/--changing 则用数字起卦法，可复现）
//   qimen   奇门遁甲
//   liuren  大六壬
//   huican  三式互参（周易×六壬×奇门）
//
// --mode 说明：
//   prompt  （默认）只返回起盘数据 + 解卦提示词，解读交给调用方
//   full    额外附带程序渲染的盘面/解卦文字（interpret 字段），供调用方参考

// castOutput 是 cast 子命令的统一 JSON 输出结构。
type castOutput struct {
	OK            bool   `json:"ok"`
	Error         string `json:"error,omitempty"`
	Method        string `json:"method"`                  // zhouyi|qimen|liuren|huican
	Question      string `json:"question,omitempty"`      // 所问之事
	QuestionType  string `json:"questionType,omitempty"`  // 规范化后的类型键
	QuestionLabel string `json:"questionLabel,omitempty"` // 类型中文标签
	Time          string `json:"time"`                    // 起盘时刻（RFC3339）
	Summary       string `json:"summary,omitempty"`       // 一句话盘面摘要（卦名/课体/局数）
	Prompt        string `json:"prompt"`                  // ⭐ 解卦提示词（核心交付物）
	Interpret     string `json:"interpret,omitempty"`     // full 模式：渲染好的盘面/解卦文字
}

// runCast 解析参数并执行一次非交互起盘，把 JSON 写到 stdout。
// 返回进程退出码（0 成功，2 参数错误，1 起盘失败）。
func runCast(args []string) int {
	var (
		method   = "zhouyi"
		question string
		qtypeStr string
		coin     bool
		yarrow   bool
		upper    = -1
		lower    = -1
		changing = -1
		timeStr     string
		lon         float64
		noTrueSolar bool
		mode        = "prompt"
		// 六壬可选：年命/行年（不填则按纯时间起课）
		benmingStr string
		gender     string
		birthYear  int
	)

	// 极简参数解析：--key value / --flag
	for i := 0; i < len(args); i++ {
		a := args[i]
		next := func() string {
			if i+1 < len(args) {
				i++
				return args[i]
			}
			return ""
		}
		switch a {
		case "--method", "-m":
			method = strings.ToLower(next())
		case "--question", "-q":
			question = next()
		case "--type", "-t":
			qtypeStr = next()
		case "--coin":
			coin = true
		case "--yarrow", "--shi", "--shicao":
			yarrow = true
		case "--upper":
			upper, _ = strconv.Atoi(next())
		case "--lower":
			lower, _ = strconv.Atoi(next())
		case "--changing":
			changing, _ = strconv.Atoi(next())
		case "--time":
			timeStr = next()
		case "--lon", "--longitude":
			lon, _ = strconv.ParseFloat(next(), 64)
		case "--no-truesolar", "--no-true-solar":
			noTrueSolar = true
		case "--mode":
			mode = strings.ToLower(next())
		case "--benming", "--shengxiao", "--zodiac":
			benmingStr = next()
		case "--gender", "--sex":
			gender = next()
		case "--birthyear", "--birth-year":
			birthYear, _ = strconv.Atoi(next())
		case "-h", "--help":
			fmt.Print(castUsage)
			return 0
		default:
			fmt.Fprintf(os.Stderr, "未知参数：%s\n\n%s", a, castUsage)
			return 2
		}
	}

	// 起盘时刻：默认当前；--time 支持 RFC3339（含时区）。
	t := time.Now()
	if timeStr != "" {
		parsed, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return emitCastError(method, fmt.Sprintf("--time 解析失败（需 RFC3339，如 2026-05-30T14:30:00+08:00）：%v", err), 2)
		}
		t = parsed
	}
	// 真太阳时校正（与 Web 端一致）：默认即启用。
	// 未传 --lon（lon==0）时由 ResolveLongitude 回退到兜底经度（北京 116.4°E）；
	// 传 --lon 则用用户经度。如需关闭校正，传 --no-truesolar。
	if !noTrueSolar {
		t = ApplyTrueSolarTime(t, ResolveLongitude(lon))
	}

	qt := ParseQuestionType(qtypeStr)
	full := mode == "full"
	nm := liurenNianMing{benming: benmingStr, gender: gender, birthYear: birthYear}

	switch method {
	case "zhouyi", "周易", "liuyao", "六爻":
		return castZhouyi(question, qt, coin, yarrow, upper, lower, changing, t, full)
	case "qimen", "奇门", "奇门遁甲":
		return castQimen(question, qt, t, full)
	case "liuren", "六壬", "大六壬":
		return castLiuren(question, qt, t, full, nm)
	case "huican", "互参", "三式":
		return castHuican(question, qt, t, full, nm)
	default:
		return emitCastError(method, fmt.Sprintf("未知 method：%q（可选 zhouyi|qimen|liuren|huican）", method), 2)
	}
}

func castZhouyi(question string, qt QuestionType, coin, yarrow bool, upper, lower, changing int, t time.Time, full bool) int {
	var r DivinationResult
	// 起卦法优先级：
	//   1) 给齐 upper+lower（changing 可为 0=无变爻）→ 数字法，结果可复现（--coin 可强制跳过）
	//   2) --yarrow → 蓍草法（大衍揲蓍，随机，最正统）
	//   3) 默认 / --coin → 铜钱法（金钱卦，随机）
	switch {
	case !coin && upper > 0 && lower > 0 && changing >= 0:
		r = DivineByNumber(upper, lower, changing)
	case yarrow && !coin:
		r = DivineByYarrow()
	default:
		r = DivineByCoins()
	}
	r.Time = t
	r.QuestionType = qt

	out := castOutput{
		OK:            true,
		Method:        "zhouyi",
		Question:      question,
		QuestionType:  string(qt),
		QuestionLabel: QuestionTypeLabel(qt),
		Time:          t.Format(time.RFC3339),
		Prompt:        GenerateAIPrompt(r, question),
	}
	if r.MainHex != nil {
		out.Summary = fmt.Sprintf("第%d卦 %s卦", r.MainHex.Number, r.MainHex.Name)
		if r.ChangeHex != nil {
			out.Summary += fmt.Sprintf(" → 之 第%d卦 %s卦", r.ChangeHex.Number, r.ChangeHex.Name)
		}
	}
	if full {
		out.Interpret = InterpretResult(r)
	}
	return emitCast(out)
}

func castQimen(question string, qt QuestionType, t time.Time, full bool) int {
	pan, err := qimen.BuildPan(t)
	if err != nil {
		return emitCastError("qimen", fmt.Sprintf("奇门起局失败：%v", err), 1)
	}
	focus := qimen.FocusGuide(string(qt))
	out := castOutput{
		OK:            true,
		Method:        "qimen",
		Question:      question,
		QuestionType:  string(qt),
		QuestionLabel: QuestionTypeLabel(qt),
		Time:          t.Format(time.RFC3339),
		Summary:       fmt.Sprintf("%s%d局 · 值符%s落%s · 值使%s落%s", pan.Ctx.Dun, pan.Ctx.Ju, pan.ZhiFuStar, pan.ZhiFuPalace, pan.ZhiShiGate, pan.ZhiShiPalace),
		Prompt:        qimen.GenerateAIPrompt(pan, question, focus, string(qt)),
	}
	if full {
		out.Interpret = pan.Render()
	}
	return emitCast(out)
}

// liurenNianMing 收集六壬的可选「年命/行年」输入（原始字符串，未解析）。
// 三项皆空时即为纯时间起课；任一非空则尝试接入年命断法。
type liurenNianMing struct {
	benming   string // 本命：地支字或生肖字，如「亥」「猪」「属猪」
	gender    string // 性别：男/女（用于行年顺逆），可空
	birthYear int    // 出生公历年份（用于行年），0=不用
}

// applyTo 把年命输入解析后填入六壬 Context。
// 解析失败的本命会被忽略（仅打 stderr 提示），不阻断起课。
func (n liurenNianMing) applyTo(ctx *liuren.Context) {
	if n.benming != "" {
		if z, ok := liuren.ParseBenMing(n.benming); ok {
			ctx.BenMing = &z
		} else {
			fmt.Fprintf(os.Stderr, "[cast] 无法识别本命/生肖 %q，已忽略（可填地支「亥」或生肖「猪」）\n", n.benming)
		}
	}
	if n.gender != "" {
		ctx.Gender = n.gender
	}
	if n.birthYear > 0 {
		ctx.BirthYear = n.birthYear
	}
}

func castLiuren(question string, qt QuestionType, t time.Time, full bool, nm liurenNianMing) int {
	ctx, err := liuren.BuildContext(t)
	if err != nil {
		return emitCastError("liuren", fmt.Sprintf("六壬起课失败：%v", err), 1)
	}
	ctx.QuestionType = string(qt)
	nm.applyTo(ctx)
	pan := liuren.DivineWithContext(ctx)
	focus := FocusGuideLiuRen(qt)
	leishen := LeiShenDirective(pan, qt)
	out := castOutput{
		OK:            true,
		Method:        "liuren",
		Question:      question,
		QuestionType:  string(qt),
		QuestionLabel: QuestionTypeLabel(qt),
		Time:          t.Format(time.RFC3339),
		Summary:       fmt.Sprintf("%s%s日 · %s", pan.Ctx.Gan.String(), pan.Ctx.DayZhi.String(), pan.KeTi.Name),
		Prompt:        liuren.GenerateAIPrompt(pan, question, focus, leishen),
	}
	if full {
		out.Interpret = liuren.Render(pan)
	}
	return emitCast(out)
}

func castHuican(question string, qt QuestionType, t time.Time, full bool, nm liurenNianMing) int {
	r, err := HuCanDivineNianMing(t, question, qt, nm)
	if err != nil {
		return emitCastError("huican", fmt.Sprintf("互参起盘失败：%v", err), 1)
	}
	out := castOutput{
		OK:            true,
		Method:        "huican",
		Question:      question,
		QuestionType:  string(qt),
		QuestionLabel: QuestionTypeLabel(qt),
		Time:          t.Format(time.RFC3339),
		Prompt:        HuCanPrompt(r),
	}
	if r.Zhouyi.MainHex != nil {
		out.Summary = fmt.Sprintf("周易%s卦 · 六壬%s · 奇门%s%d局",
			r.Zhouyi.MainHex.Name, r.LiuRenPan.KeTi.Name, r.QimenPan.Ctx.Dun, r.QimenPan.Ctx.Ju)
	}
	if full {
		out.Interpret = HuCanText(r)
	}
	return emitCast(out)
}

// emitCast 把成功结果序列化为 JSON 输出到 stdout。
func emitCast(out castOutput) int {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
	return 0
}

// emitCastError 输出错误 JSON（stdout 保持机器可读），返回退出码。
func emitCastError(method, msg string, code int) int {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	_ = enc.Encode(castOutput{OK: false, Method: method, Error: msg})
	return code
}

const castUsage = `周易占卜 · cast 非交互子命令（供 agent / skill 调用）

用法：
  zhouyi cast --method <zhouyi|qimen|liuren|huican> --question "所问之事" [选项]

选项：
  --method, -m   术数：zhouyi(周易六爻) | qimen(奇门) | liuren(六壬) | huican(三式互参)，默认 zhouyi
  --question,-q  所问之事
  --type, -t     问题类型：career|wealth|relation|health|decision|timing|other（也接受中文/数字）
  --time         起盘时刻，RFC3339，如 2026-05-30T14:30:00+08:00，默认当前时刻
  --lon          经度（东经为正），用于真太阳时校正；不传则回退默认经度（北京 116.4°E）
  --no-truesolar 关闭真太阳时校正（默认开启，与 Web 端一致）
  --mode         prompt(默认，仅起盘+提示词) | full(额外附渲染盘面文字)
  周易专属（起卦法）：
  --coin         强制用铜钱法（金钱卦，随机；默认即此法）
  --yarrow       用蓍草法（大衍揲蓍，随机，最正统）；与 --coin 同时给时以 --coin 为准
  --upper N      上卦数  --lower N 下卦数  --changing N 变爻位(0=无变爻)
                 （三者齐全且未加 --coin 时走数字起卦法，结果可复现；优先级高于 --yarrow）
  六壬/互参可选（年命断法，不填则按纯时间起课）：
  --benming      本命：地支「亥」或生肖「猪」「属猪」皆可（用于年命救应）
  --gender       性别：男|女（用于推行年顺逆，配合 --birthyear）
  --birthyear    出生公历年份，如 1990（用于推行年）

输出：JSON（stdout），核心字段 prompt 为解卦提示词，解读交由调用方完成。
说明：六壬起课本身只需时间；--benming/--gender/--birthyear 为可选增强，
      提供后会接入「年命/行年」断法，让吉凶救应判得更贴合，不提供不影响起课。
`
