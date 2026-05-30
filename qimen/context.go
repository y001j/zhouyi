package qimen

import (
	"fmt"
	"time"

	"github.com/6tail/lunar-go/calendar"
)

// Context 奇门遁甲起局输入要素。
//
// 由时间戳派生而来，不需要额外输入（与六壬不同，后者需要占时之类）。
type Context struct {
	Time time.Time // 阳历起局时刻

	// 四柱（干支）
	YearGZ  string // 年柱干支，如 "甲辰"
	MonthGZ string // 月柱干支
	DayGZ   string // 日柱干支
	HourGZ  string // 时柱干支

	// 四柱的分量（避免调用方重复切字符串）
	YearGan, YearZhi   string
	MonthGan, MonthZhi string
	DayGan, DayZhi     string
	HourGan, HourZhi   string

	// 节气
	JieQi          string    // 当前所处节气名（最近一次已过的节气）
	JieQiStartTime time.Time // 当前节气开始的精确时刻
	NextJieQi      string
	NextJieQiStart time.Time

	// 奇门起局三要素
	Dun    string // "阳遁" 或 "阴遁"
	Yuan   string // "上元" / "中元" / "下元"
	Ju     int    // 1..9

	// 时柱旬首与遁干
	Xunshou string // 如 "甲子"
	Dungan  string // 旬首六仪（戊/己/庚/辛/壬/癸）

	// 旬空
	XunKong [2]string

	// 原始 lunar（便于调用方需要时进一步取农历信息）
	Lunar *calendar.Lunar
}

// BuildContext 根据阳历时刻构造奇门起局上下文（简化拆补法：不处理超神接气）。
func BuildContext(t time.Time) (*Context, error) {
	lunar := calendar.NewLunarFromDate(t)

	// 1) 四柱干支（精确到节气/时辰边界）
	yearGZ := lunar.GetYearInGanZhiExact()
	monthGZ := lunar.GetMonthInGanZhiExact()
	dayGZ := lunar.GetDayInGanZhiExact()
	hourGZ := lunar.GetTimeInGanZhi()
	if len([]rune(dayGZ)) < 2 || len([]rune(hourGZ)) < 2 {
		return nil, fmt.Errorf("四柱干支解析失败：day=%q hour=%q", dayGZ, hourGZ)
	}

	ctx := &Context{
		Time:    t,
		YearGZ:  yearGZ,
		MonthGZ: monthGZ,
		DayGZ:   dayGZ,
		HourGZ:  hourGZ,
		Lunar:   lunar,
	}
	ctx.YearGan, ctx.YearZhi = splitGZ(yearGZ)
	ctx.MonthGan, ctx.MonthZhi = splitGZ(monthGZ)
	ctx.DayGan, ctx.DayZhi = splitGZ(dayGZ)
	ctx.HourGan, ctx.HourZhi = splitGZ(hourGZ)

	// 2) 节气：取最近一次已过的节气（含"节"与"气"都算）
	prev := lunar.GetPrevJieQi()
	if prev == nil {
		return nil, fmt.Errorf("获取前一节气失败")
	}
	jieQiName := prev.GetName()
	// lunar-go 给的是包含24节气的完整序列，名字应与 JieqiDun 的 key 一致
	prevSolar := prev.GetSolar()
	ctx.JieQi = jieQiName
	ctx.JieQiStartTime = time.Date(
		prevSolar.GetYear(), time.Month(prevSolar.GetMonth()), prevSolar.GetDay(),
		prevSolar.GetHour(), prevSolar.GetMinute(), prevSolar.GetSecond(),
		0, t.Location())

	if next := lunar.GetNextJieQi(); next != nil {
		ctx.NextJieQi = next.GetName()
		ns := next.GetSolar()
		ctx.NextJieQiStart = time.Date(
			ns.GetYear(), time.Month(ns.GetMonth()), ns.GetDay(),
			ns.GetHour(), ns.GetMinute(), ns.GetSecond(),
			0, t.Location())
	}

	// 3) 阴阳遁
	dun, ok := JieqiDun[ctx.JieQi]
	if !ok {
		return nil, fmt.Errorf("未知节气：%s", ctx.JieQi)
	}
	ctx.Dun = dun

	// 4) 上中下元（按日干支查元 —— 简化拆补法，不处理超神接气）
	yuan, ok := UpperMiddleLowerYuan[ctx.DayGZ]
	if !ok {
		return nil, fmt.Errorf("未知日干支：%s", ctx.DayGZ)
	}
	ctx.Yuan = yuan

	// 5) 局数
	jushu, ok := JushuByJieqi[ctx.JieQi]
	if !ok {
		return nil, fmt.Errorf("节气 %s 无局数定义", ctx.JieQi)
	}
	ctx.Ju = jushu[YuanIndex[yuan]]

	// 6) 时柱旬首 & 遁干
	ctx.Xunshou = findXunshou(ctx.HourGZ)
	if ctx.Xunshou == "" {
		return nil, fmt.Errorf("未能为时柱 %s 查到旬首", ctx.HourGZ)
	}
	ctx.Dungan = XunshouToDungan[ctx.Xunshou]

	// 7) 旬空
	ctx.XunKong = XunKong[ctx.Xunshou]

	return ctx, nil
}

// splitGZ 把 "甲子" 拆成 ("甲", "子")。失败时返回两个空串。
func splitGZ(gz string) (string, string) {
	r := []rune(gz)
	if len(r) < 2 {
		return "", ""
	}
	return string(r[0]), string(r[1])
}

// findXunshou 返回给定干支所属的六旬旬首（如 "甲子"）。
// 不命中时返回空串。
func findXunshou(gz string) string {
	for xunshou, members := range XunshouToGanzhi {
		for _, m := range members {
			if m == gz {
				return xunshou
			}
		}
	}
	return ""
}

// Summary 简要描述本次起局上下文（CLI 调试用）。
func (c *Context) Summary() string {
	return fmt.Sprintf("四柱 %s %s %s %s ｜ %s · %s · %s %d局 ｜ 旬首 %s(遁干 %s) ｜ 旬空 %s %s",
		c.YearGZ, c.MonthGZ, c.DayGZ, c.HourGZ,
		c.JieQi, c.Yuan, c.Dun, c.Ju,
		c.Xunshou, c.Dungan,
		c.XunKong[0], c.XunKong[1])
}
