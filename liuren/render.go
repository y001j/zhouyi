package liuren

import (
	"fmt"
	"strings"
)

// Render 生成完整的终端风格盘面文字
func Render(pan *Pan) string {
	var b strings.Builder
	b.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	b.WriteString("◎ 大六壬起课\n")
	b.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	b.WriteString("  " + pan.Ctx.Summary() + "\n")
	kong := pan.Ctx.XunKongPair()
	b.WriteString(fmt.Sprintf("  旬空：%s %s\n\n", kong[0], kong[1]))

	b.WriteString(renderTianDiPan(pan))
	b.WriteString("\n")
	b.WriteString(renderSiKe(pan))
	b.WriteString("\n")
	b.WriteString(renderSanChuan(pan))
	b.WriteString("\n")
	b.WriteString("【课体】" + pan.KeTi.Name + " —— " + pan.KeTi.Summary + "\n")

	// 附加课格标签
	if len(pan.Tags) > 0 {
		b.WriteString("\n【附加课格】\n")
		for _, tg := range pan.Tags {
			b.WriteString("  · " + tg.Name + " —— " + tg.Summary + "\n")
		}
	}

	// 神煞落位
	if len(pan.ShenSha) > 0 {
		b.WriteString("\n【神煞落位】\n")
		for _, ss := range pan.ShenSha {
			if ss.Zhi < 0 {
				continue
			}
			upperTJ := ""
			if pan.TianJiang[ss.Zhi] >= 0 {
				upperTJ = "（" + pan.TianJiang[ss.Zhi].String() + "位）"
			}
			b.WriteString(fmt.Sprintf("  · %s → 地盘 %s%s · %s\n", ss.Name, ss.Zhi, upperTJ, ss.Desc))
		}
	}

	// 年命
	if pan.NianMing != nil {
		b.WriteString("\n【年命】\n")
		if pan.NianMing.BenMing != nil {
			bm := pan.NianMing.BenMing
			kg := ""
			if bm.IsKong {
				kg = " · 空"
			}
			b.WriteString(fmt.Sprintf("  · 本命 %s：乘 %s %s（%s）%s\n",
				bm.Zhi, bm.Upper, bm.TianJiang, bm.LiuQin, kg))
			b.WriteString(fmt.Sprintf("    【%s】%s\n", bm.Ying, bm.Ying.Desc()))
		}
		if pan.NianMing.XingNian != nil {
			xn := pan.NianMing.XingNian
			kg := ""
			if xn.IsKong {
				kg = " · 空"
			}
			b.WriteString(fmt.Sprintf("  · 行年 %s：乘 %s %s（%s）%s\n",
				xn.Zhi, xn.Upper, xn.TianJiang, xn.LiuQin, kg))
			b.WriteString(fmt.Sprintf("    【%s】%s\n", xn.Ying, xn.Ying.Desc()))
		}
	}

	// 毕法赋匹配
	if len(pan.BiFa) > 0 {
		b.WriteString("\n【毕法赋匹配】\n")
		for _, e := range pan.BiFa {
			b.WriteString(fmt.Sprintf("  · 第%d · %s：%s\n    释：%s\n", e.Number, e.Title, e.Text, e.Note))
		}
	}
	return b.String()
}

// renderTianDiPan 表格形式：地盘、天盘、天将
func renderTianDiPan(pan *Pan) string {
	var b strings.Builder
	b.WriteString("【天地盘】\n")
	b.WriteString("  地盘  ")
	for i := 0; i < 12; i++ {
		b.WriteString(fmt.Sprintf(" %s ", Zhi(i)))
	}
	b.WriteString("\n  天盘  ")
	for i := 0; i < 12; i++ {
		b.WriteString(fmt.Sprintf(" %s ", pan.TianPan[i]))
	}
	b.WriteString("\n  天将  ")
	for i := 0; i < 12; i++ {
		b.WriteString(fmt.Sprintf(" %s ", pan.TianJiang[i].Short()))
	}
	b.WriteString("\n")
	return b.String()
}

func renderSiKe(pan *Pan) string {
	var b strings.Builder
	b.WriteString("【四课】（右起一课、二课、三课、四课）\n")
	// 一行：上神
	b.WriteString("  上神   ")
	for i := 3; i >= 0; i-- {
		b.WriteString(fmt.Sprintf(" %s ", pan.SiKe[i].Upper))
	}
	b.WriteString("\n  下神   ")
	for i := 3; i >= 0; i-- {
		// 第一课下神实为日干（以寄宫地支显示并括注日干）
		if pan.SiKe[i].Index == 1 {
			b.WriteString(fmt.Sprintf("%s(%s)", pan.SiKe[i].Lower, pan.Ctx.Gan))
		} else {
			b.WriteString(fmt.Sprintf(" %s  ", pan.SiKe[i].Lower))
		}
	}
	b.WriteString("\n  关系   ")
	for i := 3; i >= 0; i-- {
		b.WriteString(fmt.Sprintf(" %s ", pan.SiKe[i].Relation))
	}
	b.WriteString("\n")
	return b.String()
}

func renderSanChuan(pan *Pan) string {
	var b strings.Builder
	b.WriteString("【三传】（发传法：" + pan.SanChuan.Method + "）\n")
	rows := []ChuanEntry{pan.SanChuan.Chu, pan.SanChuan.Zhong, pan.SanChuan.Mo}
	for _, r := range rows {
		kong := ""
		if r.IsKong {
			kong = " · 空"
		}
		b.WriteString(fmt.Sprintf("  %s  %s · %s · %s%s\n",
			r.Name, r.Zhi, r.TianJiang, r.LiuQin, kong))
	}
	return b.String()
}

// RenderSummary 单行摘要（供菜单/日志）
func RenderSummary(pan *Pan) string {
	return fmt.Sprintf("%s｜%s｜初%s 中%s 末%s｜%s",
		pan.Ctx.Summary(),
		pan.KeTi.Name,
		pan.SanChuan.Chu.Zhi, pan.SanChuan.Zhong.Zhi, pan.SanChuan.Mo.Zhi,
		pan.SanChuan.Method,
	)
}
