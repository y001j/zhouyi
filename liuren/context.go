package liuren

import (
	"fmt"
	"time"

	"github.com/6tail/lunar-go/calendar"
)

// Context 大六壬起课的输入要素：
//
//	时间、日干支、占时、月将、昼夜，以及可选的本命与性别。
type Context struct {
	Time       time.Time // 阳历起课时刻
	Gan        Gan       // 日干（六壬以"日干支日"为主，23:00 后换日按 Exact2）
	DayZhi     Zhi       // 日支
	JiaziIndex int       // 日柱的六十甲子序号 0..59（用于旬空）
	ZhanShi    Zhi       // 占时支（由时刻换算）
	YueJiang   Zhi       // 月将支（以中气换将）
	QiName     string    // 当前所处中气名，便于说明
	ZhouYe     bool      // true=昼（卯~申时），false=夜（酉~寅时）

	// 可选：
	BenMing  *Zhi   // 本命（生肖对应地支）
	Gender   string // "男" / "女"
	BirthYear int   // 出生公历年份（用于行年）
	QuestionType string // 问题类型（career/wealth/relation/health/decision/timing；空=未指定）

	// 原始 Lunar，便于需要时进一步取节气干支
	Lunar *calendar.Lunar
}

// jiaziOrder 把干支字符串（如 "甲子"）转成 0..59 的序号
func jiaziOrder(gz string) int {
	if len([]rune(gz)) < 2 {
		return -1
	}
	runes := []rune(gz)
	g := ParseGan(string(runes[0]))
	z := ParseZhi(string(runes[1]))
	if g < 0 || z < 0 {
		return -1
	}
	// 六十甲子：每 10 个一旬，天干重复 6 次；序号 n 满足 n%10=g, n%12=z
	for n := 0; n < 60; n++ {
		if n%10 == int(g) && n%12 == int(z) {
			return n
		}
	}
	return -1
}

// ResolveZhanShi 按地支时辰计算占时支：
//
//	子时 23:00-01:00，丑 01-03，……，亥 21-23
func ResolveZhanShi(t time.Time) Zhi {
	h := t.Hour()
	// (h+1)/2 mod 12 == 0 对应子时
	idx := ((h + 1) / 2) % 12
	return Zhi(idx)
}

// IsDayByZhi 卯~申（5-17 点）属昼，酉~寅（17 点到次日 5 点）属夜
func IsDayByZhi(z Zhi) bool {
	return z >= Mao && z <= Shen
}

// ResolveYueJiang 以最近一个已过的中气决定月将
func ResolveYueJiang(lunar *calendar.Lunar, loc *time.Location) (Zhi, string) {
	// 从上一个节气往前找，直至找到中气
	current := lunar
	for i := 0; i < 24; i++ {
		jq := current.GetPrevJieQi()
		if jq == nil {
			break
		}
		name := jq.GetName()
		if z, ok := QiToYueJiang[name]; ok {
			return z, name
		}
		// 往前再退到该节气发生前 1 分钟
		s := jq.GetSolar()
		prevTime := time.Date(s.GetYear(), time.Month(s.GetMonth()), s.GetDay(),
			s.GetHour(), s.GetMinute(), s.GetSecond(), 0, loc).Add(-time.Minute)
		current = calendar.NewLunarFromDate(prevTime)
	}
	// 兜底（不应到达）
	return Hai, "雨水"
}

// BuildContext 根据阳历时刻构造起课上下文
func BuildContext(t time.Time) (*Context, error) {
	lunar := calendar.NewLunarFromDate(t)

	ganStr := lunar.GetDayGanExact()
	zhiStr := lunar.GetDayZhiExact()
	gan := ParseGan(ganStr)
	zhi := ParseZhi(zhiStr)
	if gan < 0 || zhi < 0 {
		return nil, fmt.Errorf("解析日干支失败：%s%s", ganStr, zhiStr)
	}

	dayGanZhi := lunar.GetDayInGanZhiExact()
	ji := jiaziOrder(dayGanZhi)
	if ji < 0 {
		return nil, fmt.Errorf("六十甲子序号解析失败：%s", dayGanZhi)
	}

	zhanShi := ResolveZhanShi(t)
	yueJiang, qiName := ResolveYueJiang(lunar, t.Location())

	return &Context{
		Time:       t,
		Gan:        gan,
		DayZhi:     zhi,
		JiaziIndex: ji,
		ZhanShi:    zhanShi,
		YueJiang:   yueJiang,
		QiName:     qiName,
		ZhouYe:     IsDayByZhi(zhanShi),
		Lunar:      lunar,
	}, nil
}

// Summary 描述起课上下文（用于调试/CLI 回显）
func (c *Context) Summary() string {
	zhou := "夜"
	if c.ZhouYe {
		zhou = "昼"
	}
	return fmt.Sprintf("日柱 %s%s（六十甲子第%d）｜占时 %s｜月将 %s（%s）｜%s占",
		c.Gan, c.DayZhi, c.JiaziIndex+1, c.ZhanShi, c.YueJiang, c.QiName, zhou)
}

// XunKongPair 返回当前旬的空亡两支
func (c *Context) XunKongPair() [2]Zhi {
	xunStart := (c.JiaziIndex / 10) * 10
	if pair, ok := XunKongByXunIndex[xunStart]; ok {
		return pair
	}
	return [2]Zhi{-1, -1}
}
