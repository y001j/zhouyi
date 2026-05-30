package qimen

import (
	"fmt"
	"strings"
)

// 本文件实现"类神直指"：按问题类型预先定义一组类神（九星/八门/八神/天干），
// 在本盘里自动定位这些类神落宫、乘神、是否空亡/入墓/击刑。
//
// 与六壬 LeiShenDirective 同构，让 LLM 不用自己去盘面里找用神。

// LeishenKind 类神种类
type LeishenKind string

const (
	KindStar   LeishenKind = "star"
	KindDoor   LeishenKind = "door"
	KindGod    LeishenKind = "god"
	KindStem   LeishenKind = "stem"
	KindDayGan LeishenKind = "dayGan" // 日干（代表"我"）
	KindHourGan LeishenKind = "hourGan" // 时干（代表所问之事/他）
)

// LeishenSpec 一条类神规格
type LeishenSpec struct {
	Label string       // 人读标签，如 "求财类神"
	Kind  LeishenKind  // 按哪种维度定位
	Names []string     // 候选名（星/门/神/干），任一命中即算
	Note  string       // 说明（对断占的意义）
}

// QimenLeishenByQType 按问题类型返回该问题应查找的类神规格组
//
// 依据：《奇门宝鉴》"九星所主 / 八门所主 / 八神类神 / 十干类神" + 张志春《神奇之门》实战
var QimenLeishenByQType = map[string][]LeishenSpec{
	"career": {
		{Label: "官印/职位类神", Kind: KindDoor, Names: []string{"开门"}, Note: "开门主讨谋图望、入官见贵、迁官进爵"},
		{Label: "决策/首脑类神", Kind: KindStar, Names: []string{"天心"}, Note: "天心为医为谋，主上级/决策之机"},
		{Label: "官贵/提携神", Kind: KindGod, Names: []string{"值符"}, Note: "值符在传/在宫皆主贵人扶持"},
		{Label: "文书/考核类神", Kind: KindStar, Names: []string{"天辅"}, Note: "天辅主文曲文化，考核文书应验"},
		{Label: "文书/信息类神", Kind: KindDoor, Names: []string{"景门"}, Note: "景门主上书献策、文书公告"},
		{Label: "我（求测者）", Kind: KindDayGan, Note: "日干代表求测者本人"},
		{Label: "所问之事", Kind: KindHourGan, Note: "时干代表所问之事/对方"},
	},
	"wealth": {
		{Label: "求财类神·吉", Kind: KindDoor, Names: []string{"生门"}, Note: "生门主生旺、求财第一门"},
		{Label: "财星/储粮类神", Kind: KindStar, Names: []string{"天任"}, Note: "天任为左辅，主稳重田产资财"},
		{Label: "金银财宝·类神", Kind: KindGod, Names: []string{"六合"}, Note: "六合主婚姻交易，合财之神"},
		{Label: "正财干", Kind: KindStem, Names: []string{"戊"}, Note: "戊为天之阳土，主田产财货"},
		{Label: "求财忌神·玄武", Kind: KindGod, Names: []string{"玄武"}, Note: "玄武主盗贼暗耗"},
		{Label: "求财忌神·天空", Kind: KindGod, Names: []string{"九地"}, Note: "九地（或天空）主虚诈空耗"},
		{Label: "我（求测者）", Kind: KindDayGan, Note: "日干主我之财力"},
		{Label: "所求之财", Kind: KindHourGan, Note: "时干代表所求之物"},
	},
	"relation": {
		{Label: "婚姻和合神", Kind: KindGod, Names: []string{"六合"}, Note: "六合主媒合缔交"},
		{Label: "阴私/情感神", Kind: KindGod, Names: []string{"太阴"}, Note: "太阴主暗助、柔情"},
		{Label: "情欲/桃花类神", Kind: KindDoor, Names: []string{"休门", "杜门"}, Note: "休主和美、杜主闭密"},
		{Label: "文曲/雅情类神", Kind: KindStar, Names: []string{"天辅"}, Note: "天辅主文雅之情"},
		{Label: "女方类神（乙）", Kind: KindStem, Names: []string{"乙"}, Note: "乙为柔木女奇"},
		{Label: "男方类神（庚）", Kind: KindStem, Names: []string{"庚"}, Note: "庚为刚金男象"},
		{Label: "我（求测者）", Kind: KindDayGan},
		{Label: "对方", Kind: KindHourGan},
	},
	"health": {
		{Label: "疾病/病位类神", Kind: KindStar, Names: []string{"天芮"}, Note: "天芮巨门，主疾病师徒之星，病位所在"},
		{Label: "医药/解病类神", Kind: KindStar, Names: []string{"天心"}, Note: "天心武曲，主医药谋略"},
		{Label: "病邪·白虎", Kind: KindGod, Names: []string{"白虎"}, Note: "白虎主刀伤血光疾病"},
		{Label: "惊疑怪异", Kind: KindGod, Names: []string{"腾蛇"}, Note: "腾蛇主惊恐怪异、虚火"},
		{Label: "忌门·杜死", Kind: KindDoor, Names: []string{"杜门", "死门"}, Note: "杜为闭塞、死为丧亡，占病大忌"},
		{Label: "喜门·生门", Kind: KindDoor, Names: []string{"生门"}, Note: "生门为生气所在，病得生则愈"},
		{Label: "我（病者）", Kind: KindDayGan},
	},
	"decision": {
		{Label: "决断类神·天心", Kind: KindStar, Names: []string{"天心"}, Note: "天心主谋断"},
		{Label: "思虑类神·天辅", Kind: KindStar, Names: []string{"天辅"}, Note: "天辅主文曲思虑"},
		{Label: "贵人指引", Kind: KindGod, Names: []string{"值符"}, Note: "值符临吉宫则有明人可请"},
		{Label: "阻滞之象·勾陈", Kind: KindGod, Names: []string{"九地"}, Note: "勾陈/九地主拖延阻滞"},
		{Label: "果断之门·开门", Kind: KindDoor, Names: []string{"开门"}, Note: "开门果断可行"},
		{Label: "闭塞之门·杜门", Kind: KindDoor, Names: []string{"杜门"}, Note: "杜门主停止、不宜进"},
		{Label: "我（决策者）", Kind: KindDayGan},
		{Label: "所决之事", Kind: KindHourGan},
	},
	"timing": {
		{Label: "吉将集合", Kind: KindGod, Names: []string{"值符", "六合", "太阴"}, Note: "三吉将落三传主运势顺遂"},
		{Label: "凶将集合", Kind: KindGod, Names: []string{"白虎", "玄武", "腾蛇"}, Note: "凶将多聚则运势受阻"},
		{Label: "吉门·开休生", Kind: KindDoor, Names: []string{"开门", "休门", "生门"}, Note: "三吉门主事事顺遂"},
		{Label: "凶门·死惊伤", Kind: KindDoor, Names: []string{"死门", "惊门", "伤门"}, Note: "凶门主阻滞"},
		{Label: "命主（我）", Kind: KindDayGan},
	},
}

// LeishenLocation 单条类神在本盘中的定位结果
type LeishenLocation struct {
	Spec       LeishenSpec
	Found      bool
	PalaceName string   // 落宫名
	PalaceFei  int      // 落宫飞星索引
	Context    string   // 所乘/所处：如 "乘天心、开门、值符" 或 "落坎一宫"
	Flags      []string // 旗标：空亡 / 入墓 / 击刑 / 门迫 / 得令 / 失令 等
}

// LocateLeishen 定位单条类神（多落点时取第一个命中；门/星/神都是唯一落点）
func LocateLeishen(pan *Pan, spec LeishenSpec) LeishenLocation {
	loc := LeishenLocation{Spec: spec}
	if pan == nil {
		return loc
	}

	// 确定本条类神在哪一格
	var palFei int = -1
	switch spec.Kind {
	case KindStar:
		palFei = findCellByStar(pan, spec.Names)
	case KindDoor:
		palFei = findCellByDoor(pan, spec.Names)
	case KindGod:
		palFei = findCellByGod(pan, spec.Names)
	case KindStem:
		// 先查天盘
		palFei = findCellByHeavenStem(pan, spec.Names)
		if palFei < 0 {
			palFei = findCellByEarthStem(pan, spec.Names)
		}
	case KindDayGan:
		palFei = findDutyStemCell(pan, pan.Ctx.DayGan)
	case KindHourGan:
		palFei = findDutyStemCell(pan, pan.Ctx.HourGan)
	}

	if palFei < 0 {
		return loc
	}
	loc.Found = true
	loc.PalaceFei = palFei
	c := pan.Cells[palFei]
	loc.PalaceName = c.PalaceName

	// 组装 context：该格乘哪些（星/门/神/天盘干/地盘干）
	parts := make([]string, 0, 5)
	if c.HeavenStem != "" {
		parts = append(parts, "天盘"+c.HeavenStem)
	}
	if c.EarthStem != "" {
		parts = append(parts, "地盘"+c.EarthStem)
	}
	if c.Star != "" {
		if c.StarWangShuai != "" {
			parts = append(parts, fmt.Sprintf("%s(%s)", c.Star, c.StarWangShuai))
		} else {
			parts = append(parts, c.Star)
		}
	}
	if c.Door != "" {
		parts = append(parts, c.Door)
	}
	if c.God != "" {
		parts = append(parts, c.God)
	}
	loc.Context = strings.Join(parts, " · ")

	// flags
	if c.IsVoid {
		loc.Flags = append(loc.Flags, "空亡")
	}
	if c.IsTianStemMu {
		loc.Flags = append(loc.Flags, "天盘干入墓")
	}
	if c.IsEarthStemMu {
		loc.Flags = append(loc.Flags, "地盘干入墓")
	}
	if c.IsJiXing {
		loc.Flags = append(loc.Flags, "击刑")
	}
	if c.IsDoorPo {
		loc.Flags = append(loc.Flags, "门迫")
	}
	if c.IsDoorSheng {
		loc.Flags = append(loc.Flags, "门得生")
	}
	if c.IsYima {
		loc.Flags = append(loc.Flags, "驿马")
	}
	// 值符/值使落宫
	if c.PalaceName == pan.ZhiFuPalace {
		loc.Flags = append(loc.Flags, "值符宫")
	}
	if c.PalaceName == pan.ZhiShiPalace {
		loc.Flags = append(loc.Flags, "值使宫")
	}
	return loc
}

// LeiShenDirective 按问题类型生成"类神直指"文本块（markdown 列表）
func LeiShenDirective(pan *Pan, qType string) string {
	specs, ok := QimenLeishenByQType[qType]
	if !ok || len(specs) == 0 {
		return ""
	}
	var b strings.Builder
	for _, spec := range specs {
		loc := LocateLeishen(pan, spec)
		if !loc.Found {
			b.WriteString(fmt.Sprintf("- **%s**：本盘未现 —— %s\n", spec.Label, spec.Note))
			continue
		}
		flags := ""
		if len(loc.Flags) > 0 {
			flags = " **【" + strings.Join(loc.Flags, "/") + "】**"
		}
		note := ""
		if spec.Note != "" {
			note = " —— " + spec.Note
		}
		b.WriteString(fmt.Sprintf("- **%s**：落 %s · %s%s%s\n",
			spec.Label, loc.PalaceName, loc.Context, flags, note))
	}
	return b.String()
}

// ========== 辅助：在九宫中按各字段查找第一个命中格 ==========

func findCellByStar(p *Pan, names []string) int {
	set := sliceToSet(names)
	for i, c := range p.Cells {
		if set[c.Star] {
			return i
		}
	}
	return -1
}
func findCellByDoor(p *Pan, names []string) int {
	set := sliceToSet(names)
	for i, c := range p.Cells {
		if set[c.Door] {
			return i
		}
	}
	return -1
}
func findCellByGod(p *Pan, names []string) int {
	set := sliceToSet(names)
	for i, c := range p.Cells {
		if set[c.God] {
			return i
		}
	}
	return -1
}
func findCellByHeavenStem(p *Pan, names []string) int {
	set := sliceToSet(names)
	for i, c := range p.Cells {
		if set[c.HeavenStem] {
			return i
		}
	}
	return -1
}
func findCellByEarthStem(p *Pan, names []string) int {
	set := sliceToSet(names)
	for i, c := range p.Cells {
		if set[c.EarthStem] {
			return i
		}
	}
	return -1
}

// findDutyStemCell 找一个天干在盘上的落宫。
// 若是"甲"，它不直接出现在盘面，而是遁于旬首六仪——返回旬首六仪所在的宫。
func findDutyStemCell(p *Pan, stem string) int {
	if stem == "甲" && p.Ctx != nil && p.Ctx.Dungan != "" {
		stem = p.Ctx.Dungan
	}
	if i := findCellByHeavenStem(p, []string{stem}); i >= 0 {
		return i
	}
	return findCellByEarthStem(p, []string{stem})
}

func sliceToSet(names []string) map[string]bool {
	set := make(map[string]bool, len(names))
	for _, n := range names {
		set[n] = true
	}
	return set
}
