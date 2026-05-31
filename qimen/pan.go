package qimen

import (
	"fmt"
	"strings"
	"time"
)

// Cell 九宫单格的全部信息（飞星序索引 0..8 对应 坎1..离9）。
type Cell struct {
	PalaceFei  int    // 飞星序索引 0..8
	PalaceName string // 宫位名，如 "坎一宫"
	EarthStem  string // 地盘干（布仪奇所得）
	HeavenStem string // 天盘干（值符加时干后旋转所得）；中5宫为空
	Star       string // 九星；中5宫恒为"天禽"
	Door       string // 八门；中5宫为空
	God        string // 八神；中5宫为空

	// 五行相关
	PalaceWuXing string // 宫位五行（固定）
	HeavenWuXing string // 天盘干五行
	EarthWuXing  string // 地盘干五行
	StarWuXing   string // 九星五行
	DoorWuXing   string // 八门五行
	StarWangShuai string // 九星按月令的旺/相/休/囚/死

	// 门宫关系
	DoorPalaceRel  string // 门与宫的生克关系：比和 / 门生宫 / 宫生门 / 门克宫 / 宫克门
	IsDoorPo       bool   // 门迫（宫克门，即门被宫所制）
	IsDoorSheng    bool   // 门相生（宫生门 → 吉）

	// 常用标志
	IsVoid     bool     // 是否含旬空地支
	IsYima     bool     // 是否含时支驿马
	IsTianStemMu   bool // 天盘干入墓（落墓宫）
	IsEarthStemMu  bool // 地盘干入墓
	IsJiXing       bool // 六仪击刑（该宫是当前旬首所击之刑位）
	Branches   []string // 该宫所辖地支
}

// Pan 一盘奇门遁甲完整盘面。
type Pan struct {
	Ctx   *Context // 起局上下文（时间、四柱、局数、旬首等）
	Cells [9]Cell  // 九宫

	ZhiFuStar    string // 值符星（如 "天蓬"）；中宫寄宫时为 "天禽"
	ZhiFuPalace  string // 值符落宫名（可能为 "中五宫"；此时实际寄坤2）
	ZhiShiGate   string // 值使门
	ZhiShiPalace string // 值使落宫名

	YiMaZhi       string   // 时支对应的驿马地支
	WangXiangXS   [5]string // 本月月支对应 [旺,相,休,囚,死] 五行
}

// BuildPan 一次性完成起局并组装九宫盘。
func BuildPan(t time.Time) (*Pan, error) {
	ctx, err := BuildContext(t)
	if err != nil {
		return nil, err
	}
	return BuildPanWithContext(ctx), nil
}

// BuildPanWithContext 给定已构造好的 Context，直接布盘。
func BuildPanWithContext(ctx *Context) *Pan {
	// 1. 地盘
	earth := LayEarthStems(ctx.Dun, ctx.Ju)

	// 2. 值符落宫与值符星
	zfStar, zfPalFei, zfRawIsMid := LocateZhiFu(earth, ctx.Dungan, ctx.HourGan)

	// 3. 天盘干
	heaven := RotateHeavenStems(earth, ctx.Dungan, ctx.HourGan, zfPalFei)

	// 4. 九星
	stars := RotateStars(earth, ctx.Dungan, ctx.HourGan, zfPalFei)

	// 5. 八神
	gods := LayGods(ctx.Dun, zfPalFei)

	// 6. 值使门与落宫
	zsGate, zsPalFei := LocateZhiShi(earth, ctx.Dungan, ctx.Xunshou, ctx.HourGZ, ctx.Dun)
	doors := LayDoors(zsGate, zsPalFei, ctx.Dun)

	// 7. 旬空 & 驿马
	isVoidBranch := func(zhi string) bool {
		return zhi == ctx.XunKong[0] || zhi == ctx.XunKong[1]
	}
	yima := YiMa[ctx.HourZhi]

	// 8. 组装九宫
	jiXingPal := LiuyiJiXingPalace[ctx.Dungan] // 本旬首六仪所"击之刑位"（飞星索引）
	var cells [9]Cell
	for i := 0; i < 9; i++ {
		branches := PalaceToZhi[i]
		cell := Cell{
			PalaceFei:    i,
			PalaceName:   Palaces[i],
			EarthStem:    earth[i],
			HeavenStem:   heaven[i],
			Star:         stars[i],
			Door:         doors[i],
			God:          gods[i],
			Branches:     branches,
			PalaceWuXing: PalaceWuXingByFei[i],
		}
		// 五行
		if s := earth[i]; s != "" {
			cell.EarthWuXing = GanWuXing[s]
		}
		if s := heaven[i]; s != "" {
			cell.HeavenWuXing = GanWuXing[s]
		}
		if s := stars[i]; s != "" {
			cell.StarWuXing = StarWuXing[s]
			cell.StarWangShuai = WangShuaiByMonth(ctx.MonthZhi, cell.StarWuXing)
		}
		if s := doors[i]; s != "" {
			cell.DoorWuXing = DoorWuXing[s]
			cell.DoorPalaceRel = WuXingRelation(cell.DoorWuXing, cell.PalaceWuXing)
			// 门迫：宫克门（门被宫制）为门迫；宫生门为吉
			if WuXingKe(cell.PalaceWuXing, cell.DoorWuXing) {
				cell.IsDoorPo = true
			}
			if WuXingSheng(cell.PalaceWuXing, cell.DoorWuXing) {
				cell.IsDoorSheng = true
			}
		}
		// 空亡 / 驿马
		for _, b := range branches {
			if isVoidBranch(b) {
				cell.IsVoid = true
			}
			if b == yima {
				cell.IsYima = true
			}
		}
		// 天干入墓
		if heaven[i] != "" && IsStemInMu(heaven[i], i) {
			cell.IsTianStemMu = true
		}
		if earth[i] != "" && IsStemInMu(earth[i], i) {
			cell.IsEarthStemMu = true
		}
		// 六仪击刑：须旬首六仪【实际落于】其刑位才成立（与 detectLiuYiJiXing 口径一致）。
		// 仅"该宫恰为理论刑位"还不够——旬首干须真落此宫。
		if i == jiXingPal && earth[i] == ctx.Dungan {
			cell.IsJiXing = true
		}
		cells[i] = cell
	}

	// 9. 值符落宫名
	zfPalName := Palaces[zfPalFei]
	if zfRawIsMid {
		zfPalName = "中五宫" // 原始落中5，但实际寄坤2
	}
	zsPalName := Palaces[zsPalFei]

	// 10. 月支旺相休囚死（用于五行衡量）
	var wxxqs [5]string
	if v, ok := WangXiangXiuQiuSi[ctx.MonthZhi]; ok {
		wxxqs = v
	}

	return &Pan{
		Ctx:          ctx,
		Cells:        cells,
		ZhiFuStar:    zfStar,
		ZhiFuPalace:  zfPalName,
		ZhiShiGate:   zsGate,
		ZhiShiPalace: zsPalName,
		YiMaZhi:      yima,
		WangXiangXS:  wxxqs,
	}
}

// Render 生成 3×3 ASCII 九宫盘字符串，便于 CLI 调试与日志。
//
// 布局（按习惯顺序）：
//
//	巽四(左上)  离九(中上)  坤二(右上)
//	震三(左中)  中五(中)    兑七(右中)
//	艮八(左下)  坎一(中下)  乾六(右下)
func (p *Pan) Render() string {
	// 飞星索引到 3x3 网格位置
	layout := [3][3]int{
		{3, 8, 1}, // 巽4, 离9, 坤2
		{2, 4, 6}, // 震3, 中5, 兑7
		{7, 0, 5}, // 艮8, 坎1, 乾6
	}

	lines := make([]string, 0, 16)
	lines = append(lines, fmt.Sprintf("【%s · %s · %d局】%s",
		p.Ctx.Dun, p.Ctx.Yuan, p.Ctx.Ju, p.Ctx.Summary()))
	lines = append(lines, fmt.Sprintf("值符：%s落%s　值使：%s落%s　驿马：%s",
		p.ZhiFuStar, p.ZhiFuPalace, p.ZhiShiGate, p.ZhiShiPalace, p.YiMaZhi))
	lines = append(lines, "")

	cellW := 16
	border := strings.Repeat("─", cellW)
	sep := "┼" + border + "┼" + border + "┼" + border + "┼"
	top := "┌" + border + "┬" + border + "┬" + border + "┐"
	bot := "└" + border + "┴" + border + "┴" + border + "┘"

	lines = append(lines, top)
	for row := 0; row < 3; row++ {
		// 每格显示 3 行：神|星|门，天盘干|地盘干，宫位名
		var r1, r2, r3 strings.Builder
		r1.WriteString("│")
		r2.WriteString("│")
		r3.WriteString("│")
		for col := 0; col < 3; col++ {
			c := p.Cells[layout[row][col]]
			mark := ""
			if c.IsVoid {
				mark += "空"
			}
			if c.IsYima {
				mark += "马"
			}
			r1.WriteString(padCJK(fmt.Sprintf(" %s %s %s", c.God, c.Star, c.Door), cellW))
			r2.WriteString(padCJK(fmt.Sprintf(" 天%s 地%s %s", c.HeavenStem, c.EarthStem, mark), cellW))
			r3.WriteString(padCJK(" "+c.PalaceName, cellW))
			r1.WriteString("│")
			r2.WriteString("│")
			r3.WriteString("│")
		}
		lines = append(lines, r1.String())
		lines = append(lines, r2.String())
		lines = append(lines, r3.String())
		if row < 2 {
			lines = append(lines, sep)
		}
	}
	lines = append(lines, bot)
	return strings.Join(lines, "\n")
}

// padCJK 按可见宽度（中文 2 / ASCII 1）填充到指定列宽。超出则截断。
func padCJK(s string, width int) string {
	w := 0
	var b strings.Builder
	for _, r := range s {
		rw := 1
		if r >= 0x4E00 && r <= 0x9FFF {
			rw = 2
		} else if r >= 0x3000 && r <= 0x303F { // CJK 符号标点
			rw = 2
		}
		if w+rw > width {
			break
		}
		b.WriteRune(r)
		w += rw
	}
	for ; w < width; w++ {
		b.WriteByte(' ')
	}
	return b.String()
}
