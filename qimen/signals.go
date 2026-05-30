package qimen

import (
	"fmt"
	"sort"
	"strings"
)

// ExtractQimenSignals 从本盘抽取 6-8 条本局专属关键信号（非通用口诀）。
// 返回可直接作为"## 本局关键信号"小节内容的 markdown 列表字符串。
func ExtractQimenSignals(pan *Pan) string {
	if pan == nil {
		return ""
	}
	var b strings.Builder

	// —— 1. 日干落宫（"我"之状态）
	b.WriteString(signalStemLocation(pan, pan.Ctx.DayGan, "日干（我）"))

	// —— 2. 时干落宫（"他 / 所问"之状态）——避免与日干重复时仍要输出
	if pan.Ctx.HourGan != pan.Ctx.DayGan {
		b.WriteString(signalStemLocation(pan, pan.Ctx.HourGan, "时干（所问）"))
	}

	// —— 3. 值符宫格详情（机枢之所在）
	b.WriteString(signalZhiFuDetail(pan))

	// —— 4. 值使宫格详情
	b.WriteString(signalZhiShiDetail(pan))

	// —— 5. 三奇落宫分布
	b.WriteString(signalThreeAuspiceDistribution(pan))

	// —— 6. 九星月令旺衰（挑出最旺/最衰各一颗）
	b.WriteString(signalStarWangShuai(pan))

	// —— 7. 空亡与击刑警示
	b.WriteString(signalVoidAndJiXing(pan))

	// —— 8. 阴阳/动静格局
	b.WriteString(signalDunOverall(pan))

	return b.String()
}

// ============ 各信号生成器 ============

func signalStemLocation(pan *Pan, stem, label string) string {
	// 甲遁于旬首六仪，不直接出现在盘上——改查旬首所在宫
	lookup := stem
	noteOfJia := ""
	if stem == "甲" && pan.Ctx != nil && pan.Ctx.Dungan != "" {
		lookup = pan.Ctx.Dungan
		noteOfJia = fmt.Sprintf("（甲遁于旬首 %s）", lookup)
	}
	palFei := findCellByHeavenStem(pan, []string{lookup})
	kind := "天盘"
	if palFei < 0 {
		palFei = findCellByEarthStem(pan, []string{lookup})
		kind = "地盘"
	}
	if palFei < 0 {
		return fmt.Sprintf("- **%s · %s**：本盘未现（罕见，请重核盘面）\n", label, stem)
	}
	c := pan.Cells[palFei]
	parts := []string{}
	if c.Star != "" {
		ws := ""
		if c.StarWangShuai != "" {
			ws = "(" + c.StarWangShuai + ")"
		}
		parts = append(parts, c.Star+ws)
	}
	if c.Door != "" {
		parts = append(parts, c.Door)
	}
	if c.God != "" {
		parts = append(parts, c.God)
	}
	flags := cellFlagsText(c, pan)
	return fmt.Sprintf("- **%s · %s**%s：%s %s 宫，乘 %s%s\n",
		label, stem, noteOfJia, kind, c.PalaceName, strings.Join(parts, "·"), flags)
}

func signalZhiFuDetail(pan *Pan) string {
	// 找值符宫（按 pan.ZhiFuPalace 名字；中五宫则记为中5）
	palFei := -1
	if pan.ZhiFuPalace == "中五宫" {
		palFei = 1 // 寄坤2
	} else {
		for i, c := range pan.Cells {
			if c.PalaceName == pan.ZhiFuPalace {
				palFei = i
				break
			}
		}
	}
	if palFei < 0 {
		return ""
	}
	c := pan.Cells[palFei]
	ws := ""
	if c.StarWangShuai != "" {
		ws = "(" + c.StarWangShuai + ")"
	}
	flags := cellFlagsText(c, pan)
	return fmt.Sprintf("- **值符 %s 落宫**：%s · 天%s地%s · %s · %s%s\n",
		pan.ZhiFuStar, pan.ZhiFuPalace, c.HeavenStem, c.EarthStem, c.Star+ws, c.Door+"·"+c.God, flags)
}

func signalZhiShiDetail(pan *Pan) string {
	palFei := -1
	for i, c := range pan.Cells {
		if c.PalaceName == pan.ZhiShiPalace {
			palFei = i
			break
		}
	}
	if palFei < 0 {
		return ""
	}
	c := pan.Cells[palFei]
	flags := cellFlagsText(c, pan)
	doorPo := ""
	if c.IsDoorPo {
		doorPo = "（门迫 · " + c.DoorPalaceRel + "）"
	} else if c.IsDoorSheng {
		doorPo = "（门得生 · " + c.DoorPalaceRel + "）"
	}
	return fmt.Sprintf("- **值使 %s 落宫**：%s · 天%s地%s · 所乘 %s · %s%s%s\n",
		pan.ZhiShiGate, pan.ZhiShiPalace,
		c.HeavenStem, c.EarthStem, c.Star, c.God,
		doorPo, flags)
}

func signalThreeAuspiceDistribution(pan *Pan) string {
	// 三奇乙丙丁在天盘各落何宫
	parts := []string{}
	for _, qi := range []string{"乙", "丙", "丁"} {
		if i := findCellByHeavenStem(pan, []string{qi}); i >= 0 {
			c := pan.Cells[i]
			doorTag := c.Door
			if doorTag == "" {
				doorTag = "无门"
			}
			flagBits := []string{}
			if c.IsVoid {
				flagBits = append(flagBits, "空")
			}
			if c.IsTianStemMu {
				flagBits = append(flagBits, "入墓")
			}
			flag := ""
			if len(flagBits) > 0 {
				flag = "【" + strings.Join(flagBits, "/") + "】"
			}
			parts = append(parts, fmt.Sprintf("%s→%s(%s)%s", qi, c.PalaceName, doorTag, flag))
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return "- **三奇分布**：" + strings.Join(parts, " · ") + "\n"
}

func signalStarWangShuai(pan *Pan) string {
	// 挑选月令值下的九星旺相休囚死；列出"旺相"与"囚死"的星
	type starWs struct {
		star, ws string
	}
	seen := map[string]bool{}
	var lst []starWs
	for _, c := range pan.Cells {
		if c.Star == "" || seen[c.Star] {
			continue
		}
		seen[c.Star] = true
		lst = append(lst, starWs{c.Star, c.StarWangShuai})
	}
	order := map[string]int{"旺": 0, "相": 1, "休": 2, "囚": 3, "死": 4}
	sort.Slice(lst, func(i, j int) bool {
		return order[lst[i].ws] < order[lst[j].ws]
	})
	// 取前 2（旺/相）与后 2（囚/死）
	if len(lst) < 4 {
		return ""
	}
	good := []string{}
	bad := []string{}
	for _, x := range lst {
		if x.ws == "旺" || x.ws == "相" {
			good = append(good, fmt.Sprintf("%s(%s)", x.star, x.ws))
		} else if x.ws == "囚" || x.ws == "死" {
			bad = append(bad, fmt.Sprintf("%s(%s)", x.star, x.ws))
		}
	}
	return fmt.Sprintf("- **九星月令旺衰**（月支 %s）：得令 %s ｜ 失令 %s\n",
		pan.Ctx.MonthZhi, strings.Join(good, "·"), strings.Join(bad, "·"))
}

func signalVoidAndJiXing(pan *Pan) string {
	voidPalaces := []string{}
	var jiXing *Cell
	for i, c := range pan.Cells {
		if c.IsVoid {
			voidPalaces = append(voidPalaces, c.PalaceName)
		}
		if c.IsJiXing {
			copy := pan.Cells[i]
			jiXing = &copy
		}
	}
	if len(voidPalaces) == 0 && jiXing == nil {
		return ""
	}
	parts := []string{}
	if len(voidPalaces) > 0 {
		parts = append(parts, fmt.Sprintf("旬空地支（%s、%s）落在：%s",
			pan.Ctx.XunKong[0], pan.Ctx.XunKong[1], strings.Join(voidPalaces, "、")))
	}
	if jiXing != nil {
		parts = append(parts, fmt.Sprintf("旬首 %s 击刑位：%s",
			pan.Ctx.Dungan, jiXing.PalaceName))
	}
	return "- **空亡 / 击刑警示**：" + strings.Join(parts, "；") + "\n"
}

func signalDunOverall(pan *Pan) string {
	// 阳遁=主动进取；阴遁=主静守；结合值符宫五行与日干五行的关系
	var energy string
	if pan.Ctx.Dun == "阳遁" {
		energy = "主动势进、外显"
	} else {
		energy = "主阴势收、内藏"
	}
	// 日干五行 vs 值符宫五行
	dayWX := GanWuXing[pan.Ctx.DayGan]
	var zfPalWX string
	for _, c := range pan.Cells {
		if c.PalaceName == pan.ZhiFuPalace {
			zfPalWX = c.PalaceWuXing
			break
		}
	}
	rel := WuXingRelation(dayWX, zfPalWX)
	return fmt.Sprintf("- **格局走势**：%s · 第%d局（%s）；日干 %s(%s) 与值符宫 %s(%s) 关系：**%s**\n",
		pan.Ctx.Dun, pan.Ctx.Ju, energy,
		pan.Ctx.DayGan, dayWX, pan.ZhiFuPalace, zfPalWX, rel)
}

// cellFlagsText 把格子的常用标志拼成 "【空亡/门迫/…】" 文本，若无则返回空
func cellFlagsText(c Cell, pan *Pan) string {
	flags := []string{}
	if c.IsVoid {
		flags = append(flags, "空亡")
	}
	if c.IsTianStemMu {
		flags = append(flags, "天干入墓")
	}
	if c.IsEarthStemMu {
		flags = append(flags, "地干入墓")
	}
	if c.IsJiXing {
		flags = append(flags, "击刑")
	}
	if c.IsDoorPo {
		flags = append(flags, "门迫")
	}
	if c.IsYima {
		flags = append(flags, "驿马")
	}
	if len(flags) == 0 {
		return ""
	}
	return " **【" + strings.Join(flags, "/") + "】**"
}
