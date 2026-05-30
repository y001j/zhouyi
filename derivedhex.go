package main

import (
	"fmt"
	"strings"
)

// lineYang 将爻值（6/7/8/9）归一化为 0（阴）/1（阳）
func lineYang(v int) int {
	if v == 7 || v == 9 {
		return 1
	}
	return 0
}

// linesToYang 将六爻数值数组转为阴阳数组（0/1）
func linesToYang(lines [6]int) [6]int {
	var out [6]int
	for i, v := range lines {
		out[i] = lineYang(v)
	}
	return out
}

// yangToLines 将阴阳数组（0/1）转回可查询的六爻数值（7/8），无变爻
func yangToLines(yy [6]int) [6]int {
	var out [6]int
	for i, v := range yy {
		if v == 1 {
			out[i] = 7
		} else {
			out[i] = 8
		}
	}
	return out
}

// MutualHexagram 互卦
// 取本卦二、三、四爻为下卦（即新卦初、二、三爻），
// 三、四、五爻为上卦（即新卦四、五、六爻）。
// 揭示事情内部潜在的走向。
func MutualHexagram(lines [6]int) *Hexagram {
	yy := linesToYang(lines)
	var m [6]int
	// 新下卦 = 原二、三、四爻（索引1,2,3）
	m[0], m[1], m[2] = yy[1], yy[2], yy[3]
	// 新上卦 = 原三、四、五爻（索引2,3,4）
	m[3], m[4], m[5] = yy[2], yy[3], yy[4]
	return FindHexagramByLines(yangToLines(m))
}

// OppositeHexagram 错卦（旁通）
// 六爻阴阳全部反转，揭示对立面、背面的视角。
func OppositeHexagram(lines [6]int) *Hexagram {
	yy := linesToYang(lines)
	var o [6]int
	for i, v := range yy {
		o[i] = 1 - v
	}
	return FindHexagramByLines(yangToLines(o))
}

// ReverseHexagram 综卦（反对）
// 整卦上下颠倒（初↔上，二↔五，三↔四），换个立场看同一件事。
// 对于上下对称的卦（如乾、坤、坎、离、大过、颐、小过、中孚），综卦即本卦自身。
func ReverseHexagram(lines [6]int) *Hexagram {
	yy := linesToYang(lines)
	var r [6]int
	for i := 0; i < 6; i++ {
		r[i] = yy[5-i]
	}
	return FindHexagramByLines(yangToLines(r))
}

// LinePosition 单爻的爻位属性
type LinePosition struct {
	Pos       int    // 1-6
	IsYang    bool   // 是否阳爻
	IsProper  bool   // 当位（阳爻居奇位 / 阴爻居偶位）
	IsCentral bool   // 居中（二爻、五爻）
	Relation  string // 与对应爻（初-四、二-五、三-上）的承乘应比关系摘要
}

// AnalyzeLinePositions 分析六爻的爻位关系（当位、中、应）
// lines[0]=初爻 ... lines[5]=上爻
func AnalyzeLinePositions(lines [6]int) [6]LinePosition {
	var out [6]LinePosition
	yy := linesToYang(lines)
	for i := 0; i < 6; i++ {
		pos := i + 1
		isYang := yy[i] == 1
		// 当位：奇位(1,3,5)阳 / 偶位(2,4,6)阴
		isProper := (pos%2 == 1 && isYang) || (pos%2 == 0 && !isYang)
		isCentral := pos == 2 || pos == 5
		out[i] = LinePosition{
			Pos:       pos,
			IsYang:    isYang,
			IsProper:  isProper,
			IsCentral: isCentral,
		}
	}
	// 应爻关系：初-四、二-五、三-上。阴阳相反为"有应"，相同为"无应/敌应"
	pairs := [3][2]int{{0, 3}, {1, 4}, {2, 5}}
	for _, p := range pairs {
		a, b := p[0], p[1]
		aYang := yy[a] == 1
		bYang := yy[b] == 1
		var desc string
		if aYang != bYang {
			desc = fmt.Sprintf("与第%d爻阴阳相应（有应）", b+1)
		} else {
			desc = fmt.Sprintf("与第%d爻同性不应（敌应）", b+1)
		}
		out[a].Relation = desc
		// 反向描述
		if aYang != bYang {
			out[b].Relation = fmt.Sprintf("与第%d爻阴阳相应（有应）", a+1)
		} else {
			out[b].Relation = fmt.Sprintf("与第%d爻同性不应（敌应）", a+1)
		}
	}
	return out
}

// FormatLinePositionAnalysis 生成爻位分析的文本描述（用于提示词）
func FormatLinePositionAnalysis(lines [6]int) string {
	lineNames := []string{"初", "二", "三", "四", "五", "上"}
	positions := AnalyzeLinePositions(lines)
	var b strings.Builder
	for i := 0; i < 6; i++ {
		p := positions[i]
		yy := "九"
		if !p.IsYang {
			yy = "六"
		}
		attrs := []string{}
		if p.IsProper {
			attrs = append(attrs, "当位")
		} else {
			attrs = append(attrs, "不当位")
		}
		if p.IsCentral {
			attrs = append(attrs, "居中")
		}
		attrs = append(attrs, p.Relation)
		b.WriteString(fmt.Sprintf("  · 第%d爻（%s%s）：%s\n",
			p.Pos, lineNames[i], yy, strings.Join(attrs, "，")))
	}
	return b.String()
}

// DerivedHexagrams 本卦衍生出的三类卦
type DerivedHexagrams struct {
	Mutual   *Hexagram // 互卦
	Opposite *Hexagram // 错卦
	Reverse  *Hexagram // 综卦
}

// DeriveHexagrams 计算本卦的互、错、综三卦
func DeriveHexagrams(lines [6]int) DerivedHexagrams {
	return DerivedHexagrams{
		Mutual:   MutualHexagram(lines),
		Opposite: OppositeHexagram(lines),
		Reverse:  ReverseHexagram(lines),
	}
}