package liuren

import (
	"fmt"
	"strings"
)

// InterpretGuide 断课指引：六壬基本取材顺序与要点
func InterpretGuide(pan *Pan) string {
	var b strings.Builder
	b.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	b.WriteString("【断课指引】\n")
	b.WriteString("  · 日干为「我」，日支为「所问之事」。\n")
	b.WriteString("  · 初传 = 事之因、动机；中传 = 事之过程、经历；末传 = 事之结局。\n")
	b.WriteString("  · 四课三传之神乘何天将，兼看所乘之凶吉（贵/龙/合/常 为吉，蛇/雀/勾/虎/武 多凶）。\n")
	b.WriteString("  · 六亲以日干为我：生我者父母、同我者兄弟、我生者子孙、我克者妻财、克我者官鬼。\n")
	kong := pan.Ctx.XunKongPair()
	b.WriteString(fmt.Sprintf("  · 旬空（%s、%s）：所乘之神落空者事主虚浮，不实之应。\n", kong[0], kong[1]))
	if pan.NianMing != nil {
		b.WriteString("  · 年命救应：按课传吉凶结合年命应机判最终走向。\n")
	} else {
		b.WriteString("  · 年命上神可救凶济吉（若方便提供本命，解读更贴合）。\n")
	}
	b.WriteString("\n")

	// 关键爻提示
	b.WriteString("【盘面要点】\n")
	// 初传逢空
	if pan.SanChuan.Chu.IsKong {
		b.WriteString("  · 初传落空亡：事起头未实，动机虚浮或迟延。\n")
	}
	if pan.SanChuan.Mo.IsKong {
		b.WriteString("  · 末传落空亡：事之结局难成，宜早抽身。\n")
	}
	// 三传天将凶吉扼要
	for _, ce := range []ChuanEntry{pan.SanChuan.Chu, pan.SanChuan.Zhong, pan.SanChuan.Mo} {
		b.WriteString(fmt.Sprintf("  · %s %s 乘 %s（%s）：%s\n",
			ce.Name, ce.Zhi, ce.TianJiang, ce.LiuQin,
			TianJiangMeaning[ce.TianJiang]))
	}
	// 年命应机
	if pan.NianMing != nil {
		if bm := pan.NianMing.BenMing; bm != nil {
			b.WriteString(fmt.Sprintf("  · 本命 %s 乘 %s → %s：%s\n",
				bm.Zhi, bm.Upper, bm.Ying, bm.Ying.Desc()))
		}
		if xn := pan.NianMing.XingNian; xn != nil {
			b.WriteString(fmt.Sprintf("  · 行年 %s 乘 %s → %s：%s\n",
				xn.Zhi, xn.Upper, xn.Ying, xn.Ying.Desc()))
		}
	}
	return b.String()
}
