package main

import (
	"crypto/rand"
	"math/big"
	"time"
)

// DivinationMethod 起卦方法
type DivinationMethod int

const (
	CoinMethod    DivinationMethod = iota // 铜钱法（金钱卦）
	YarrowMethod                          // 蓍草法（揲蓍法）
	NumberMethod                          // 数字起卦法
)

// LineValue 爻的数值
// 6=老阴（变），7=少阳，8=少阴，9=老阳（变）
type LineValue int

const (
	OldYin   LineValue = 6 // 老阴，三枚反面，变爻 ×
	YoungYang LineValue = 7 // 少阳，两反一正
	YoungYin  LineValue = 8 // 少阴，两正一反
	OldYang   LineValue = 9 // 老阳，三枚正面，变爻 ○
)

// DivinationResult 占卜结果
type DivinationResult struct {
	Lines        [6]int    // 六爻数值（6/7/8/9）
	MainHex      *Hexagram // 本卦
	ChangeHex    *Hexagram // 变卦（无变爻时为nil）
	ChangingPos  []int     // 变爻位置（1-6）
	Method       DivinationMethod
	Time         time.Time    // 起卦时刻
	QuestionType QuestionType // 问题类型
}

// LineInfo 单爻信息
type LineInfo struct {
	Value    int
	IsYang   bool
	IsChange bool
	Symbol   string // ─── 阳 / ── –– 阴 / ○ 老阳变 / × 老阴变
}

func lineInfo(v int) LineInfo {
	switch v {
	case 9:
		return LineInfo{9, true, true, "─○─"}
	case 7:
		return LineInfo{7, true, false, "───"}
	case 8:
		return LineInfo{8, false, false, "── ──"}
	case 6:
		return LineInfo{6, false, true, "─×─"}
	}
	return LineInfo{}
}

// secureRandN 生成 [0, n) 的随机整数
func secureRandN(n int64) int64 {
	max := big.NewInt(n)
	v, _ := rand.Int(rand.Reader, max)
	return v.Int64()
}

// CoinThrow 模拟三枚铜钱投掷一次，返回爻值
// 规则：正面=3（阳），反面=2（阴）；三枚之和
//   6 = 2+2+2 → 老阴（变）
//   7 = 2+2+3 → 少阳
//   8 = 2+3+3 → 少阴
//   9 = 3+3+3 → 老阳（变）
func CoinThrow() int {
	sum := 0
	for i := 0; i < 3; i++ {
		if secureRandN(2) == 0 {
			sum += 2 // 反面
		} else {
			sum += 3 // 正面
		}
	}
	return sum
}


// DivineByCoins 铜钱法起六爻
func DivineByCoins() DivinationResult {
	var lines [6]int
	for i := 0; i < 6; i++ {
		lines[i] = CoinThrow()
	}
	return buildResult(lines, CoinMethod)
}

// DivineByYarrow 蓍草法起六爻
func DivineByYarrow() DivinationResult {
	var lines [6]int
	for i := 0; i < 6; i++ {
		lines[i] = yarrowOneLine()
	}
	return buildResult(lines, YarrowMethod)
}

// yarrowOneLine 蓍草法一爻三变
func yarrowOneLine() int {
	stalks := 49
	remainders := [3]int{}

	for v := 0; v < 3; v++ {
		split := int(secureRandN(int64(stalks-2))) + 1
		left := split
		right := stalks - split - 1 // 挂一

		leftRem := left % 4
		if leftRem == 0 {
			leftRem = 4
		}
		rightRem := right % 4
		if rightRem == 0 {
			rightRem = 4
		}
		remainder := 1 + leftRem + rightRem
		remainders[v] = remainder
		stalks -= remainder
	}

	total := remainders[0] + remainders[1] + remainders[2]
	switch total {
	case 13:
		return 9 // 老阳
	case 17:
		return 8 // 少阴
	case 21:
		return 7 // 少阳
	case 25:
		return 6 // 老阴
	default:
		// 极少概率落不到标准值，返回少阳
		return 7
	}
}

// DivineByNumber 数字起卦（以数字为基础，适合现代使用）
// 传入两个整数：上卦数、下卦数（取模8得1-8）
func DivineByNumber(upper, lower, changingLine int) DivinationResult {
	trigramOrder := []string{"乾", "兑", "离", "震", "巽", "坎", "艮", "坤"}
	upperIdx := ((upper % 8) + 7) % 8
	lowerIdx := ((lower % 8) + 7) % 8

	upperTrigram := trigramOrder[upperIdx]
	lowerTrigram := trigramOrder[lowerIdx]

	hex := FindHexagram(upperTrigram, lowerTrigram)
	if hex == nil {
		// fallback
		return DivineByCoins()
	}

	// 根据上下卦构建爻（7=少阳, 8=少阴）
	var lines [6]int
	lt := Trigrams[lowerTrigram]
	ut := Trigrams[upperTrigram]
	for i := 0; i < 3; i++ {
		if lt.Lines[i] == 1 {
			lines[i] = 7
		} else {
			lines[i] = 8
		}
	}
	for i := 0; i < 3; i++ {
		if ut.Lines[i] == 1 {
			lines[i+3] = 7
		} else {
			lines[i+3] = 8
		}
	}

	// 设置变爻
	if changingLine >= 1 && changingLine <= 6 {
		idx := changingLine - 1
		if lines[idx] == 7 {
			lines[idx] = 9
		} else {
			lines[idx] = 6
		}
	}

	return buildResult(lines, NumberMethod)
}

func buildResult(lines [6]int, method DivinationMethod) DivinationResult {
	mainHex := FindHexagramByLines(lines)
	changeHex := FindChangedHexagram(lines)
	changingPos := CountChangingLines(lines)

	return DivinationResult{
		Lines:       lines,
		MainHex:     mainHex,
		ChangeHex:   changeHex,
		ChangingPos: changingPos,
		Method:      method,
		Time:        time.Now(),
	}
}

// InterpretResult 生成解卦文字
func InterpretResult(r DivinationResult) string {
	if r.MainHex == nil {
		return "起卦失败，请重试。"
	}

	out := ""
	out += renderHexagramArt(r)
	out += "\n"
	out += "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	out += formatHexagramInfo(r.MainHex, "本卦")
	out += "\n"

	// 变爻解释
	if len(r.ChangingPos) > 0 {
		out += "【变爻】\n"
		for _, pos := range r.ChangingPos {
			line := r.MainHex.Lines[pos-1]
			marker := "○" // 老阳
			if r.Lines[pos-1] == 6 {
				marker = "×" // 老阴
			}
			out += formatChangingLine(pos, line, marker, r.Lines[pos-1])
		}
		out += "\n"
	}

	// 变卦
	if r.ChangeHex != nil {
		out += "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
		out += formatHexagramInfo(r.ChangeHex, "变卦（之卦）")
		out += "\n"
	}

	// 解卦指引
	out += "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	out += interpretationGuide(r)

	return out
}

func renderHexagramArt(r DivinationResult) string {
	out := "\n"
	// 从上爻到初爻显示（视觉上上爻在顶部）
	for i := 5; i >= 0; i-- {
		li := lineInfo(r.Lines[i])
		pos := i + 1
		posName := []string{"初", "二", "三", "四", "五", "上"}[i]
		yinyang := "九"
		if !li.IsYang {
			yinyang = "六"
		}
		out += "  " + li.Symbol + "   " + posName + yinyang
		if li.IsChange {
			out += " ←变"
		}
		_ = pos
		out += "\n"
	}
	if r.MainHex != nil {
		out += "\n  " + r.MainHex.Symbol + " " + r.MainHex.Name + "卦（第" + intToStr(r.MainHex.Number) + "卦）"
		out += "  上" + r.MainHex.Upper + "（" + Trigrams[r.MainHex.Upper].Symbol + Trigrams[r.MainHex.Upper].Nature + "）"
		out += " 下" + r.MainHex.Lower + "（" + Trigrams[r.MainHex.Lower].Symbol + Trigrams[r.MainHex.Lower].Nature + "）\n"
	}
	return out
}

func formatHexagramInfo(h *Hexagram, title string) string {
	out := "【" + title + "】 " + h.Symbol + " " + h.Name + "卦（第" + intToStr(h.Number) + "卦）\n"
	out += "上卦：" + h.Upper + "（" + Trigrams[h.Upper].Symbol + " " + Trigrams[h.Upper].Nature + "）  "
	out += "下卦：" + h.Lower + "（" + Trigrams[h.Lower].Symbol + " " + Trigrams[h.Lower].Nature + "）\n\n"
	out += "卦辞：" + h.Judgment + "\n"
	out += "象辞：" + h.Image + "\n"
	return out
}

func formatChangingLine(pos int, line Line, marker string, val int) string {
	lineType := "老阳（三正）"
	if val == 6 {
		lineType = "老阴（三反）"
	}
	return "  第" + intToStr(pos) + "爻 " + marker + " " + lineType + "\n" +
		"  " + line.Type + "：" + line.Text + "\n"
}

func interpretationGuide(r DivinationResult) string {
	n := len(r.ChangingPos)
	out := "【解卦指引 · 朱熹《易学启蒙·考变占》】\n"
	switch n {
	case 0:
		out += "无变爻：以本卦卦辞（彖辞）断。卦象稳定，局势不变。\n"
	case 1:
		out += "一爻变：以本卦该变爻爻辞断。\n"
	case 2:
		out += "二爻变：以本卦两变爻爻辞断，仍以上爻为主。\n"
	case 3:
		out += "三爻变：占本卦与变卦之卦辞，本卦为贞（主，先）、变卦为悔（次，后）。\n"
	case 4:
		out += "四爻变：以变卦中两个不变爻之爻辞断，仍以下爻为主。\n"
	case 5:
		out += "五爻变：以变卦中唯一不变之爻辞断。\n"
	case 6:
		out += "六爻皆变：乾坤占用九、用六，余卦占变卦卦辞。\n"
		if r.MainHex != nil {
			switch r.MainHex.Name {
			case "乾":
				out += "  用九：见群龙无首，吉。\n"
			case "坤":
				out += "  用六：利永贞。\n"
			}
		}
	}
	return out
}

func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	digits := []byte{}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
