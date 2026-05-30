package qimen

// 本文件实现奇门遁甲格局自动识别。
//
// 依据：《烟波钓叟歌》「乙奇加甲戌甲午…」等原文 +《奇门宝鉴·凡例/卷三》格局总汇。
//
// 每个格局定义为一个函数 (p *Pan) → []PatternHit；同一个格局可能在多宫命中（如伏吟）。
// 主流程 Patterns(p) 依次调用所有规则，收集全部命中。

// PatternHit 一次格局命中
type PatternHit struct {
	Name         string // 格局名，如 "青龙返首"
	Category     string // "吉格" / "凶格" / "神格"
	PalaceFei    int    // 主要命中宫（飞星索引 0..8）；-1 表示全盘型（伏吟/反吟）
	Classic      string // 古籍原文/口诀（一句）
	Summary      string // 白话一句概括
	AuspiceScore int    // 吉凶分：+2大吉 +1吉 0平 -1凶 -2大凶
}

// DetectPatterns 扫描本盘命中的所有格局
func DetectPatterns(p *Pan) []PatternHit {
	if p == nil {
		return nil
	}
	var hits []PatternHit

	hits = append(hits, detectQingLongFanShou(p)...)
	hits = append(hits, detectFeiNiaoDieXue(p)...)
	hits = append(hits, detectYuNvShouMen(p)...)
	hits = append(hits, detectSanQiDeShi(p)...)
	hits = append(hits, detectSanQiShengDian(p)...)
	hits = append(hits, detectSanDun(p)...)     // 天/地/人/神/鬼遁 五格
	hits = append(hits, detectLongHuFengYun(p)...) // 龙/虎/风/云遁
	hits = append(hits, detectDaGe(p)...)
	hits = append(hits, detectXiaoGe(p)...)
	hits = append(hits, detectXingGe(p)...)
	hits = append(hits, detectTaiBaiRuYing(p)...)
	hits = append(hits, detectYingRuTaiBai(p)...)
	hits = append(hits, detectQingLongTaoZou(p)...)
	hits = append(hits, detectBaiHuChangKuang(p)...)
	hits = append(hits, detectZhuQueTouJiang(p)...)
	hits = append(hits, detectTengSheYaoJiao(p)...)
	hits = append(hits, detectSanQiRuMu(p)...)
	hits = append(hits, detectLiuYiJiXing(p)...)
	hits = append(hits, detectFuYin(p)...)
	hits = append(hits, detectFanYin(p)...)

	return hits
}

// ============ 基础辅助 ============

// cellWith 在九宫中找到"天盘干 heaven 加地盘干 earth"的格子。返回飞星索引，找不到返回 -1。
func cellWith(p *Pan, heaven, earth string) int {
	for i, c := range p.Cells {
		if c.HeavenStem == heaven && c.EarthStem == earth {
			return i
		}
	}
	return -1
}

// cellsWithHeaven 返回所有天盘干为 heaven 的格子索引
func cellsWithHeaven(p *Pan, heaven string) []int {
	var out []int
	for i, c := range p.Cells {
		if c.HeavenStem == heaven {
			out = append(out, i)
		}
	}
	return out
}

// ============ 格局规则 ============

// 青龙返首：天盘甲加地盘丙（甲遁于戊/旬首六仪；实际判定为天盘戊加地盘丙的组合）
// 古籍口诀："丙加甲，事事大吉"——注意甲遁于六仪，所以以"六仪加丙"作近似判定
// 主流识别法：天盘甲（等同于戊 + 旬首）加地盘丙；或用六仪遁于地盘
func detectQingLongFanShou(p *Pan) []PatternHit {
	var hits []PatternHit
	// 甲遁于 dungan；所以天盘 dungan 加地盘 丙 即视作甲加丙
	for i, c := range p.Cells {
		if c.HeavenStem == p.Ctx.Dungan && c.EarthStem == "丙" {
			hits = append(hits, PatternHit{
				Name: "青龙返首", Category: "吉格",
				PalaceFei:    i,
				Classic:      "丙加甲，事事大吉",
				Summary:      "旬首加丙奇于" + c.PalaceName + "：贵人相助，事事得力",
				AuspiceScore: 2,
			})
		}
	}
	return hits
}

// 飞鸟跌穴：丙加甲（天盘丙加地盘之旬首六仪）
func detectFeiNiaoDieXue(p *Pan) []PatternHit {
	var hits []PatternHit
	for i, c := range p.Cells {
		if c.HeavenStem == "丙" && c.EarthStem == p.Ctx.Dungan {
			hits = append(hits, PatternHit{
				Name: "飞鸟跌穴", Category: "吉格",
				PalaceFei:    i,
				Classic:      "甲加丙，百事大吉",
				Summary:      "丙奇加旬首于" + c.PalaceName + "：百事大吉，贵人迎合",
				AuspiceScore: 2,
			})
		}
	}
	return hits
}

// 玉女守门：丁奇（时干为丁）且落值使门所在宫
// 古籍《烟波钓叟歌》p.34："玉女守门者，丁甲子时…"——丁奇与值使门同宫
func detectYuNvShouMen(p *Pan) []PatternHit {
	// 条件：丁在天盘某宫，该宫即值使门所在
	for _, i := range cellsWithHeaven(p, "丁") {
		c := p.Cells[i]
		if c.Door == p.ZhiShiGate && p.ZhiShiPalace == c.PalaceName {
			return []PatternHit{{
				Name: "玉女守门", Category: "吉格",
				PalaceFei:    i,
				Classic:      "丁奇与值使同宫",
				Summary:      "丁奇（玉女）守值使门于" + c.PalaceName + "：吉将守门，事有神助",
				AuspiceScore: 2,
			}}
		}
	}
	return nil
}

// 三奇得使：
//   乙奇加甲戌（己）或甲午（辛）
//   丙奇加甲子（戊）或甲申（庚）
//   丁奇加甲辰（壬）或甲寅（癸）
// 判定：天盘奇 加 地盘对应旬首六仪
func detectSanQiDeShi(p *Pan) []PatternHit {
	rules := map[string][2]string{
		"乙": {"己", "辛"},
		"丙": {"戊", "庚"},
		"丁": {"壬", "癸"},
	}
	var hits []PatternHit
	for qi, sixes := range rules {
		for _, six := range sixes {
			if i := cellWith(p, qi, six); i >= 0 {
				hits = append(hits, PatternHit{
					Name: "三奇得使", Category: "吉格",
					PalaceFei:    i,
					Classic:      qi + "奇加" + six + "（旬首六仪）",
					Summary:      qi + "奇得" + six + "之使于" + p.Cells[i].PalaceName + "：最吉之格，宜施为",
					AuspiceScore: 2,
				})
			}
		}
	}
	return hits
}

// 三奇升殿：
//   乙奇落震3（木之本宫）或巽4
//   丙奇落离9（火之本宫）
//   丁奇落兑7 — 主流说法为"丁奇居兑七"为升殿（存在分歧）
func detectSanQiShengDian(p *Pan) []PatternHit {
	rules := map[string][]int{
		"乙": {2, 3}, // 震3 巽4
		"丙": {8},    // 离9
		"丁": {6},    // 兑7
	}
	var hits []PatternHit
	for qi, pals := range rules {
		for _, i := range cellsWithHeaven(p, qi) {
			for _, want := range pals {
				if i == want {
					hits = append(hits, PatternHit{
						Name: "三奇升殿", Category: "吉格",
						PalaceFei:    i,
						Classic:      qi + "奇居本宫",
						Summary:      qi + "奇升殿于" + p.Cells[i].PalaceName + "：得地气之助",
						AuspiceScore: 1,
					})
					break
				}
			}
		}
	}
	return hits
}

// 三遁 + 神鬼遁：
//
//   天遁：丙+生门+天心同宫（有的说法是丙+生门+太阴，这里取奇门宝鉴的"丙生天心"版本）
//   地遁：乙+开门+己同宫
//   人遁：丁+休门+太阴同宫
//   神遁：丙+生门+九天同宫
//   鬼遁：丁+开门+九地同宫
func detectSanDun(p *Pan) []PatternHit {
	type rule struct {
		name, stem, door, starOrGod string
		isGod                       bool // true 表示匹配八神
	}
	rules := []rule{
		{"天遁", "丙", "生门", "天心", false},
		{"地遁", "乙", "开门", "己", false}, // 己是地盘干，后面特殊处理
		{"人遁", "丁", "休门", "太阴", true},
		{"神遁", "丙", "生门", "九天", true},
		{"鬼遁", "丁", "开门", "九地", true},
	}
	var hits []PatternHit
	for _, r := range rules {
		for i, c := range p.Cells {
			if c.HeavenStem != r.stem {
				continue
			}
			if c.Door != r.door {
				continue
			}
			ok := false
			switch {
			case r.name == "地遁":
				ok = c.EarthStem == r.starOrGod
			case r.isGod:
				ok = c.God == r.starOrGod
			default:
				ok = c.Star == r.starOrGod
			}
			if ok {
				hits = append(hits, PatternHit{
					Name: r.name, Category: "吉格",
					PalaceFei:    i,
					Classic:      r.stem + "+" + r.door + "+" + r.starOrGod,
					Summary:      r.name + "于" + c.PalaceName + "：大吉，利所问之事",
					AuspiceScore: 2,
				})
			}
		}
	}
	return hits
}

// 龙/虎/风/云四遁：
//   龙遁：乙+休门落坎1
//   虎遁：乙+生门落艮8
//   风遁：乙+开门落巽4
//   云遁：乙+开门落坤2（有说"杜门+巽"；此处取宝鉴版）
func detectLongHuFengYun(p *Pan) []PatternHit {
	type rule struct {
		name, door string
		palace     int
	}
	rules := []rule{
		{"龙遁", "休门", 0},
		{"虎遁", "生门", 7},
		{"风遁", "开门", 3},
		{"云遁", "开门", 1},
	}
	var hits []PatternHit
	for _, r := range rules {
		i := r.palace
		c := p.Cells[i]
		if c.HeavenStem == "乙" && c.Door == r.door {
			hits = append(hits, PatternHit{
				Name: r.name, Category: "吉格",
				PalaceFei:    i,
				Classic:      "乙+" + r.door + "+" + c.PalaceName,
				Summary:      r.name + "于" + c.PalaceName + "：吉，利特定之事",
				AuspiceScore: 1,
			})
		}
	}
	return hits
}

// 大格：天盘庚加地盘乙
func detectDaGe(p *Pan) []PatternHit {
	if i := cellWith(p, "庚", "乙"); i >= 0 {
		return []PatternHit{{
			Name: "大格", Category: "凶格",
			PalaceFei:    i,
			Classic:      "庚加乙为大格",
			Summary:      "庚加乙于" + p.Cells[i].PalaceName + "：百事皆凶，遇事阻隔",
			AuspiceScore: -2,
		}}
	}
	return nil
}

// 小格：天盘庚加地盘壬
func detectXiaoGe(p *Pan) []PatternHit {
	if i := cellWith(p, "庚", "壬"); i >= 0 {
		return []PatternHit{{
			Name: "小格", Category: "凶格",
			PalaceFei:    i,
			Classic:      "庚加壬为小格",
			Summary:      "庚加壬于" + p.Cells[i].PalaceName + "：事难成，人离财散",
			AuspiceScore: -1,
		}}
	}
	return nil
}

// 刑格：天盘庚加地盘己
func detectXingGe(p *Pan) []PatternHit {
	if i := cellWith(p, "庚", "己"); i >= 0 {
		return []PatternHit{{
			Name: "刑格", Category: "凶格",
			PalaceFei:    i,
			Classic:      "庚加己为刑格",
			Summary:      "庚加己于" + p.Cells[i].PalaceName + "：五行相刑，宜静忌动",
			AuspiceScore: -1,
		}}
	}
	return nil
}

// 太白入荧：庚加丙（主胜客、贼即来）
func detectTaiBaiRuYing(p *Pan) []PatternHit {
	if i := cellWith(p, "庚", "丙"); i >= 0 {
		return []PatternHit{{
			Name: "太白入荧", Category: "凶格",
			PalaceFei:    i,
			Classic:      "庚加丙，贼即来",
			Summary:      "庚加丙于" + p.Cells[i].PalaceName + "：主胜客，来者之势强，防侵扰",
			AuspiceScore: -1,
		}}
	}
	return nil
}

// 荧入太白：丙加庚（客胜主、贼即去）
func detectYingRuTaiBai(p *Pan) []PatternHit {
	if i := cellWith(p, "丙", "庚"); i >= 0 {
		return []PatternHit{{
			Name: "荧入太白", Category: "凶格",
			PalaceFei:    i,
			Classic:      "丙加庚，贼即去",
			Summary:      "丙加庚于" + p.Cells[i].PalaceName + "：客胜主，宜守忌攻",
			AuspiceScore: -1,
		}}
	}
	return nil
}

// 青龙逃走：乙加辛
func detectQingLongTaoZou(p *Pan) []PatternHit {
	if i := cellWith(p, "乙", "辛"); i >= 0 {
		return []PatternHit{{
			Name: "青龙逃走", Category: "凶格",
			PalaceFei:    i,
			Classic:      "乙加辛为青龙逃走",
			Summary:      "乙加辛于" + p.Cells[i].PalaceName + "：失财散信，所谋难成",
			AuspiceScore: -1,
		}}
	}
	return nil
}

// 白虎猖狂：辛加乙
func detectBaiHuChangKuang(p *Pan) []PatternHit {
	if i := cellWith(p, "辛", "乙"); i >= 0 {
		return []PatternHit{{
			Name: "白虎猖狂", Category: "凶格",
			PalaceFei:    i,
			Classic:      "辛加乙为白虎猖狂",
			Summary:      "辛加乙于" + p.Cells[i].PalaceName + "：刀兵刑狱，主人凶险",
			AuspiceScore: -2,
		}}
	}
	return nil
}

// 朱雀投江：丁加癸
func detectZhuQueTouJiang(p *Pan) []PatternHit {
	if i := cellWith(p, "丁", "癸"); i >= 0 {
		return []PatternHit{{
			Name: "朱雀投江", Category: "凶格",
			PalaceFei:    i,
			Classic:      "丁加癸为朱雀投江",
			Summary:      "丁加癸于" + p.Cells[i].PalaceName + "：文书遗失，信息不达",
			AuspiceScore: -1,
		}}
	}
	return nil
}

// 螣蛇夭矫：癸加丁
func detectTengSheYaoJiao(p *Pan) []PatternHit {
	if i := cellWith(p, "癸", "丁"); i >= 0 {
		return []PatternHit{{
			Name: "螣蛇夭矫", Category: "凶格",
			PalaceFei:    i,
			Classic:      "癸加丁为螣蛇夭矫",
			Summary:      "癸加丁于" + p.Cells[i].PalaceName + "：惊恐怪异，心神不宁",
			AuspiceScore: -1,
		}}
	}
	return nil
}

// 三奇入墓：
//   乙入未（坤2）
//   丙入戌（乾6）
//   丁入丑（艮8）
// 判定：天盘奇所落宫是对应墓宫
func detectSanQiRuMu(p *Pan) []PatternHit {
	rules := map[string]int{
		"乙": 1, // 坤2
		"丙": 5, // 乾6
		"丁": 7, // 艮8
	}
	var hits []PatternHit
	for qi, want := range rules {
		for _, i := range cellsWithHeaven(p, qi) {
			if i == want {
				hits = append(hits, PatternHit{
					Name: "三奇入墓", Category: "凶格",
					PalaceFei:    i,
					Classic:      qi + "入" + p.Cells[i].PalaceName,
					Summary:      qi + "奇入墓于" + p.Cells[i].PalaceName + "：虚浮不实，事易中辍",
					AuspiceScore: -1,
				})
			}
		}
	}
	return hits
}

// 六仪击刑：当前旬首六仪落到相应刑位
func detectLiuYiJiXing(p *Pan) []PatternHit {
	pal, ok := LiuyiJiXingPalace[p.Ctx.Dungan]
	if !ok {
		return nil
	}
	c := p.Cells[pal]
	// 判定：旬首六仪的确在这个宫
	if c.EarthStem != p.Ctx.Dungan {
		return nil
	}
	return []PatternHit{{
		Name: "六仪击刑", Category: "凶格",
		PalaceFei:    pal,
		Classic:      p.Ctx.Xunshou + p.Ctx.Dungan + "击刑于" + c.PalaceName,
		Summary:      "旬首" + p.Ctx.Dungan + "落击刑位" + c.PalaceName + "：所谋不成，有灾厄",
		AuspiceScore: -1,
	}}
}

// 伏吟：天盘干 = 地盘干（某宫天地同干），且值符落在自己的本宫
// 简化判定：有 ≥3 格天盘干 == 地盘干，即伏吟盘
func detectFuYin(p *Pan) []PatternHit {
	same := 0
	for _, c := range p.Cells {
		if c.HeavenStem != "" && c.HeavenStem == c.EarthStem {
			same++
		}
	}
	if same >= 6 { // 8 个非中5宫，≥6 格同干 = 几乎伏吟
		return []PatternHit{{
			Name: "伏吟", Category: "凶格",
			PalaceFei:    -1,
			Classic:      "天盘同地盘",
			Summary:      "全盘伏吟：事机未发，宜静忌动",
			AuspiceScore: -1,
		}}
	}
	return nil
}

// 反吟：天盘宫与地盘宫对冲（九宫对冲 1↔9, 2↔8, 3↔7, 4↔6）
// 简化判定：局数决定阴阳遁和起点；"反吟"古法定义为"阳遁转阴遁同局号 ± 3 节气"之一。
// 此处用粗略法：天盘干数与地盘干数差 4 以上的格子数 ≥6，提示反吟盘。
func detectFanYin(p *Pan) []PatternHit {
	opposite := map[int]int{0: 8, 8: 0, 1: 7, 7: 1, 2: 6, 6: 2, 3: 5, 5: 3}
	// 判据：某格天盘干 = 对冲宫的地盘干，且至少 ≥6 格成立
	match := 0
	for i, c := range p.Cells {
		if i == 4 {
			continue
		}
		op, ok := opposite[i]
		if !ok {
			continue
		}
		if c.HeavenStem != "" && c.HeavenStem == p.Cells[op].EarthStem {
			match++
		}
	}
	if match >= 6 {
		return []PatternHit{{
			Name: "反吟", Category: "凶格",
			PalaceFei:    -1,
			Classic:      "天盘对冲地盘",
			Summary:      "全盘反吟：事有反复、动荡难安",
			AuspiceScore: -1,
		}}
	}
	return nil
}
