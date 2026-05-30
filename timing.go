package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/6tail/lunar-go/calendar"
)

// TimingInfo 起卦时间的完整信息（阳历、阴历、干支历、节气、时令消息卦）
//
// 注意：阴历（朔望历）以农历正月初一切年，干支历以立春切年，两者可能相差数天乃至一个月。
// 本结构将两套体系分开表达，避免混用语义。
type TimingInfo struct {
	SolarTime     time.Time // 原始阳历时刻
	LunarDesc     string    // 纯阴历表达：阴历年（生肖）+ 月 + 日，如「癸卯兔年正月初一」
	GanZhiYear    string    // 干支纪年（立春分界）
	GanZhiMonth   string    // 干支纪月（节气分界）
	GanZhiDay     string    // 干支纪日
	GanZhiHour    string    // 干支纪时
	GanZhiSummary string    // 干支历摘要：年 月 日 时
	JieQiName     string    // 当前所处节气区间名称（上一个节气到下一个节气之间）
	JieQiDay      string    // 距上一节气的日数描述
	NextJieQiName string    // 下一节气
	NextJieQiDate string    // 下一节气阳历日期
	MonthlyHex    *Hexagram // 时令消息卦（十二辟卦中对应当前月份者）
	MonthlyHexNote string   // 消息卦的象意摘要
}

// 十二消息卦：以地支（月建）为键
// 子=复（一阳生，冬至），丑=临，寅=泰（立春），卯=大壮，辰=夬，巳=乾，
// 午=姤（一阴生，夏至），未=遁，申=否（立秋），酉=观，戌=剥，亥=坤
var monthlyHexByZhi = map[string]struct {
	Num  int
	Name string
	Note string
}{
	"子": {24, "复", "一阳来复，冬至阳气始生，宜静候萌动"},
	"丑": {19, "临", "二阳浸长，万物将出，事有起势"},
	"寅": {11, "泰", "三阳开泰，天地交而万物通，吉利亨通"},
	"卯": {34, "大壮", "四阳盛壮，春分雷动，力强而需节制"},
	"辰": {43, "夬", "五阳决阴，清明刚决，宜果断而防躁进"},
	"巳": {1, "乾", "六阳纯阳，立夏盛极，慎防亢龙有悔"},
	"午": {44, "姤", "一阴初生，夏至阴气始起，盛中有衰兆"},
	"未": {33, "遁", "二阴浸长，小暑当退避守静，不宜强进"},
	"申": {12, "否", "三阴闭塞，立秋天地不交，宜守不宜动"},
	"酉": {20, "观", "四阴盛长，秋分宜观察省思，不宜妄动"},
	"戌": {23, "剥", "五阴剥阳，寒露将尽，当防小人防损失"},
	"亥": {2, "坤", "六阴纯阴，立冬收藏，厚德载物，静以待阳"},
}

// CaptureTiming 根据给定时刻计算完整的起卦时间信息
func CaptureTiming(t time.Time) *TimingInfo {
	lunar := calendar.NewLunarFromDate(t)

	// 阴历（朔望历）：年份用 GetYearInGanZhi，以农历正月初一切换，与月日语义一致
	lunarYearGZ := lunar.GetYearInGanZhi()
	shengXiao := lunar.GetYearShengXiao()
	lunarDesc := fmt.Sprintf("%s%s年%s月%s",
		lunarYearGZ, shengXiao,
		lunar.GetMonthInChinese(),
		lunar.GetDayInChinese())

	// 干支历：年月日时均以节气/时辰精确切换
	gzYear := lunar.GetYearInGanZhiExact()
	gzMonth := lunar.GetMonthInGanZhiExact()
	gzDay := lunar.GetDayInGanZhiExact()
	gzHour := lunar.GetTimeInGanZhi()

	info := &TimingInfo{
		SolarTime:     t,
		LunarDesc:     lunarDesc,
		GanZhiYear:    gzYear,
		GanZhiMonth:   gzMonth,
		GanZhiDay:     gzDay,
		GanZhiHour:    gzHour,
		GanZhiSummary: fmt.Sprintf("%s年 %s月 %s日 %s时", gzYear, gzMonth, gzDay, gzHour),
	}

	// 节气：获取上一节气与下一节气（含精确时刻），计算当前处于节气区间第几天
	prev := lunar.GetPrevJieQi()
	next := lunar.GetNextJieQi()
	if prev != nil {
		info.JieQiName = prev.GetName()
		prevSolar := prev.GetSolar()
		// 使用节气精确时刻（含时分秒），避免把"节气当日但未到节气点"误判为"已入新节气"
		prevTime := time.Date(
			prevSolar.GetYear(), time.Month(prevSolar.GetMonth()), prevSolar.GetDay(),
			prevSolar.GetHour(), prevSolar.GetMinute(), prevSolar.GetSecond(),
			0, t.Location())
		// 按自然日差（两个日期的日历日差）计算"第几天"：节气当日为第1天
		dayDiff := daysBetween(prevTime, t)
		info.JieQiDay = fmt.Sprintf("%s后第%d天", prev.GetName(), dayDiff+1)
	}
	if next != nil {
		info.NextJieQiName = next.GetName()
		ns := next.GetSolar()
		info.NextJieQiDate = fmt.Sprintf("%d-%02d-%02d",
			ns.GetYear(), ns.GetMonth(), ns.GetDay())
	}

	// 时令消息卦：以月支定（节气分月）
	monthZhi := lunar.GetMonthZhiExact()
	if mh, ok := monthlyHexByZhi[monthZhi]; ok {
		if mh.Num >= 1 && mh.Num <= 64 {
			info.MonthlyHex = &Hexagrams[mh.Num-1]
			info.MonthlyHexNote = mh.Note
		}
	}

	return info
}

// daysBetween 按本地自然日计算两个时刻相差的整数天数（忽略时分秒）
// 结果：from 日到 to 日的日历天数差。from == to 同一天返回 0。
func daysBetween(from, to time.Time) int {
	fromDay := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	toDay := time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, to.Location())
	return int(toDay.Sub(fromDay).Hours() / 24)
}

// FormatTimingSection 格式化时间信息为提示词段落
func FormatTimingSection(ti *TimingInfo) string {
	if ti == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString("## 起卦时间\n")
	b.WriteString(fmt.Sprintf("- 阳历：%s\n", ti.SolarTime.Format("2006-01-02 15:04")))
	b.WriteString(fmt.Sprintf("- 阴历：%s\n", ti.LunarDesc))
	b.WriteString(fmt.Sprintf("- 干支历：%s（年与月均以节气交接为界）\n", ti.GanZhiSummary))
	if ti.JieQiName != "" {
		b.WriteString(fmt.Sprintf("- 节气：%s", ti.JieQiDay))
		if ti.NextJieQiName != "" {
			b.WriteString(fmt.Sprintf("，距「%s」（%s）\n", ti.NextJieQiName, ti.NextJieQiDate))
		} else {
			b.WriteString("\n")
		}
	}
	if ti.MonthlyHex != nil {
		b.WriteString(fmt.Sprintf("- 时令消息卦：第%d卦 %s卦 %s —— %s\n",
			ti.MonthlyHex.Number, ti.MonthlyHex.Name, ti.MonthlyHex.Symbol, ti.MonthlyHexNote))
	}
	b.WriteString("\n")
	return b.String()
}

// FormatTimingHexRelation 指出时令消息卦与本卦的关系（本卦即时令、互为错综等）
// 如果没有明显关系则返回空字符串。
func FormatTimingHexRelation(ti *TimingInfo, mainHex *Hexagram, lines [6]int) string {
	if ti == nil || ti.MonthlyHex == nil || mainHex == nil {
		return ""
	}
	mh := ti.MonthlyHex
	if mh.Number == mainHex.Number {
		return fmt.Sprintf("**本卦即当下时令消息卦**（%s卦）——卦气与天时相应，当顺势而为，把握节令所主之机。\n",
			mh.Name)
	}
	if opp := OppositeHexagram(lines); opp != nil && opp.Number == mh.Number {
		return fmt.Sprintf("**本卦与时令消息卦（%s卦）互为错卦**——所问之事处于时令的反面或背面，宜反思当前方向是否逆天时而行。\n",
			mh.Name)
	}
	if rev := ReverseHexagram(lines); rev != nil && rev.Number == mh.Number {
		return fmt.Sprintf("**本卦与时令消息卦（%s卦）互为综卦**——时势正从对立角度作用，宜换位体察。\n",
			mh.Name)
	}
	if mut := MutualHexagram(lines); mut != nil && mut.Number == mh.Number {
		return fmt.Sprintf("**本卦互卦为时令消息卦（%s卦）**——时令之气潜藏于事件内部运作之中。\n",
			mh.Name)
	}
	return ""
}

// AdjacentHexagrams 本卦在六十四卦序中的前一卦与后一卦
type AdjacentHexagrams struct {
	Prev *Hexagram
	Next *Hexagram
}

// GetAdjacent 返回本卦在卦序上的前后邻卦（《序卦传》逻辑所依据）
func GetAdjacent(h *Hexagram) AdjacentHexagrams {
	if h == nil {
		return AdjacentHexagrams{}
	}
	var adj AdjacentHexagrams
	if h.Number >= 2 {
		adj.Prev = &Hexagrams[h.Number-2]
	}
	if h.Number <= 63 {
		adj.Next = &Hexagrams[h.Number]
	}
	return adj
}

// FormatAdjacentSection 格式化卦序前后邻卦段落
func FormatAdjacentSection(h *Hexagram) string {
	if h == nil {
		return ""
	}
	adj := GetAdjacent(h)
	var b strings.Builder
	b.WriteString("## 卦序前后（《序卦传》脉络）\n")
	if adj.Prev != nil {
		b.WriteString(fmt.Sprintf("- 前卦（第%d卦）：%s卦 %s —— 事之所从来\n",
			adj.Prev.Number, adj.Prev.Name, adj.Prev.Symbol))
	} else {
		b.WriteString("- 前卦：无（本卦为六十四卦之首）\n")
	}
	b.WriteString(fmt.Sprintf("- 本卦（第%d卦）：%s卦 %s\n", h.Number, h.Name, h.Symbol))
	if adj.Next != nil {
		b.WriteString(fmt.Sprintf("- 后卦（第%d卦）：%s卦 %s —— 事之所当往\n",
			adj.Next.Number, adj.Next.Name, adj.Next.Symbol))
	} else {
		b.WriteString("- 后卦：无（本卦为六十四卦之末）\n")
	}
	b.WriteString("\n")
	return b.String()
}
