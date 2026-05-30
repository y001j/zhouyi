package liuren

// 日辰格识别（《六壬大全》卷三 p534 起 "日辰章"）
//
// 该章列 18 条日干（"日"）与日支（"辰"）通过各自上神发生互动的格局，
// 是断课的"语法"。每条只返回一段简短判语，由提示词层透传给 LLM。
//
// 缩写约定：
//   - 日 = 日干，五行用 ganWuxing
//   - 辰 = 日支
//   - 日上 = 日干寄宫之天盘神（一课上神）
//   - 辰上 = 日支之天盘神（三课上神）

// RiZhenGe 一条日辰格
type RiZhenGe struct {
	Name  string // 格名
	Judge string // 古法判语
}

// DetectRiZhenGe 识别本盘命中的日辰格（可命中多条）
func DetectRiZhenGe(pan *Pan) []RiZhenGe {
	var out []RiZhenGe

	gan := pan.Ctx.Gan
	zhi := pan.Ctx.DayZhi
	wDay := GanWuXing[gan]
	wZhi := ZhiWuXing[zhi]

	ganShang := pan.TianPan[GanJiGong[gan]] // 日上
	zhenShang := pan.TianPan[zhi]           // 辰上
	wGanShang := ZhiWuXing[ganShang]
	wZhenShang := ZhiWuXing[zhenShang]

	// 五行关系工具
	gen := WuXingGenerates // a 生 b
	ke := WuXingOvercomes  // a 克 b

	// 1. 日上生日 / 辰上生辰 → 百事吉
	if gen(wGanShang, wDay) && gen(wZhenShang, wZhi) {
		out = append(out, RiZhenGe{
			Name:  "日辰俱生",
			Judge: "日上生日、辰上生辰，百事吉、宅家平安。",
		})
	}

	// 2. 日上克日 → 百事不利
	if ke(wGanShang, wDay) {
		out = append(out, RiZhenGe{
			Name:  "日上克日",
			Judge: "日上克日，事将不利，宜守不宜进；防身受外压。",
		})
	}

	// 3. 日生上神 → 日费出（耗财耗神）
	if gen(wDay, wGanShang) {
		out = append(out, RiZhenGe{
			Name:  "日生上神",
			Judge: "日生上神，气从我出，事虽成而费力费财、夜鬼魅。",
		})
	}

	// 4. 日上神来生日、辰上神来生辰 → 两家顺利
	if gen(wGanShang, wDay) && gen(wZhenShang, wZhi) {
		// 同 1，避免重复就不再加
	}

	// 5. 日生上神，辰上神生日 → 两家俱顺利有生意
	if gen(wDay, wGanShang) && gen(wZhenShang, wDay) {
		out = append(out, RiZhenGe{
			Name:  "两家生意",
			Judge: "日生上神而辰上反生日，两相滋养，主有生意、合作之喜。",
		})
	}

	// 6. 日上之神去克辰、辰上之神去克日 → 两家俱不利
	if ke(wGanShang, wZhi) && ke(wZhenShang, wDay) {
		out = append(out, RiZhenGe{
			Name:  "两家俱伤",
			Judge: "日上克辰、辰上克日，内外交攻，两家俱不利、彼此防脱跷。",
		})
	}

	// 7. 日上脱辰、辰上脱日 → 彼此防脱
	if gen(wZhi, wGanShang) && gen(wDay, wZhenShang) {
		out = append(out, RiZhenGe{
			Name:  "彼此脱气",
			Judge: "日上脱辰、辰上脱日，两家俱被泄气；静则为禄、动则遭网。",
		})
	}

	// 8. 日临辰 / 辰临日 → 互为客主
	//    日"临"辰：日干寄宫的地盘位上之天盘神 = 日支
	//    辰"临"日：日支位上之天盘神 = 日干寄宫支
	if ganShang == zhi {
		out = append(out, RiZhenGe{
			Name:  "日临辰",
			Judge: "日临于辰，主自取上门，事在外凌侵于内、客主易位。",
		})
	}
	if zhenShang == GanJiGong[gan] {
		out = append(out, RiZhenGe{
			Name:  "辰临日",
			Judge: "辰临于日，主彼来犯我，事自外侵入、宜慎守门户。",
		})
	}

	// 9. 二者皆名乱首（既日临辰又辰临日）→ 父子兄弟各离析
	if ganShang == zhi && zhenShang == GanJiGong[gan] {
		out = append(out, RiZhenGe{
			Name:  "乱首",
			Judge: "日辰互临，纲纪倒置，名为乱首，主家中父子兄弟各离析、长幼失序。",
		})
	}

	// 10. 日比和、辰比和（同五行） → 比和
	if wGanShang == wDay && wZhenShang == wZhi {
		out = append(out, RiZhenGe{
			Name:  "日辰俱比",
			Judge: "日上、辰上俱与本气比和，事虽稳但缺生意；宜守成、勿求新。",
		})
	}

	// 11. 日上克日、辰上克辰 → 内外受压
	if ke(wGanShang, wDay) && ke(wZhenShang, wZhi) {
		out = append(out, RiZhenGe{
			Name:  "内外受克",
			Judge: "日上克日且辰上克辰，内外俱压，主家事身事并困、宜静观。",
		})
	}

	// 12. 日辰冲（日干寄宫与日支地支冲）
	if isChong(GanJiGong[gan], zhi) {
		out = append(out, RiZhenGe{
			Name:  "日辰相冲",
			Judge: "日干寄宫与日支相冲，主事必反复、人不能安。",
		})
	}

	// 13. 日上、辰上互冲 → 两端动摇
	if isChong(ganShang, zhenShang) {
		out = append(out, RiZhenGe{
			Name:  "日辰上神相冲",
			Judge: "日上与辰上互冲，两端动摇，事在变动、宜审进退。",
		})
	}

	// 14. 日上、辰上为同一神 → 日辰同位
	if ganShang == zhenShang {
		out = append(out, RiZhenGe{
			Name:  "日辰同上",
			Judge: "日上与辰上同神，事在一处、人事专一；利合而不利分。",
		})
	}

	// 15. 日辰各受上神脱（日上脱日、辰上脱辰）→ 静则禄动则网
	if gen(wDay, wGanShang) && gen(wZhi, wZhenShang) {
		out = append(out, RiZhenGe{
			Name:  "日辰俱脱",
			Judge: "日辰各受上神脱气，静则为禄、动则遭网。",
		})
	}

	// 16. 空亡相关：日上、辰上落空 → 虚浮不实
	kong := pan.Ctx.XunKongPair()
	isKong := func(z Zhi) bool { return z == kong[0] || z == kong[1] }
	if isKong(ganShang) {
		out = append(out, RiZhenGe{
			Name:  "日上落空",
			Judge: "日上神落旬空，事无实根、我之处境虚浮。",
		})
	}
	if isKong(zhenShang) {
		out = append(out, RiZhenGe{
			Name:  "辰上落空",
			Judge: "辰上神落旬空，家事不实、问事难成。",
		})
	}

	// 17. 日上、辰上俱旺 → 两家俱旺
	if isWang(wGanShang, pan) && isWang(wZhenShang, pan) {
		out = append(out, RiZhenGe{
			Name:  "日辰俱旺",
			Judge: "日上、辰上俱当令旺相，两家俱旺，事易成且有力。",
		})
	}

	// 18. 日上、辰上俱休囚 → 两家俱衰
	if isXiu(wGanShang, pan) && isXiu(wZhenShang, pan) {
		out = append(out, RiZhenGe{
			Name:  "日辰俱衰",
			Judge: "日上、辰上俱休囚无气，两家俱衰，事难起色。",
		})
	}

	return out
}

// isChong 两支是否相冲（差 6 位）
func isChong(a, b Zhi) bool {
	return (int(a)-int(b)+12)%12 == 6
}

// isWang 五行 w 是否当令旺相（按月支节气月）
func isWang(w WuXing, pan *Pan) bool {
	if pan.Ctx.Lunar == nil {
		return false
	}
	mz := monthZhiOf(pan.Ctx)
	if mz < 0 {
		return false
	}
	mw := ZhiWuXing[mz]
	// 旺：同我；相：生我
	return mw == w || WuXingGenerates(mw, w)
}

// isXiu 五行 w 是否休囚（被月令克或泄）
func isXiu(w WuXing, pan *Pan) bool {
	if pan.Ctx.Lunar == nil {
		return false
	}
	mz := monthZhiOf(pan.Ctx)
	if mz < 0 {
		return false
	}
	mw := ZhiWuXing[mz]
	// 休：我生月；囚：我克月；死：月克我
	return WuXingGenerates(w, mw) || WuXingOvercomes(w, mw) || WuXingOvercomes(mw, w)
}
