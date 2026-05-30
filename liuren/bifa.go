package liuren

// BiFaEntry 毕法赋一条断语
type BiFaEntry struct {
	Number int    // 条目序号 1..100
	Title  string // 条目题名（6 字或 7 字）
	Text   string // 断语原文（通常即标题本身）
	Note   string // 释义
}

// BiFaRule 毕法赋匹配规则：给定盘面返回是否命中
type BiFaRule struct {
	Entry   BiFaEntry
	Matches func(*Pan) bool
}

// MatchBiFa 返回当前盘面命中的全部毕法赋条目（自动匹配的机械规则）
func MatchBiFa(pan *Pan) []BiFaEntry {
	var out []BiFaEntry
	for _, r := range biFaRules {
		if r.Matches(pan) {
			out = append(out, r.Entry)
		}
	}
	return out
}

// BiFaCatalog 返回《毕法赋》100 条的完整文本（知识库）
//
// 用途：作为 AI 提示词的附录注入，便于大模型结合古籍断语对盘面进行判读。
// 其中部分条目（MatchBiFa 会命中的那些）也会被自动标记为"已命中"。
func BiFaCatalog() []BiFaEntry {
	return biFaCatalog
}

// ---------- 机械规则 ----------

var biFaRules = []BiFaRule{
	// 第 5 条：六阳数足 —— 日干支 + 四课上下 + 三传 共 N 位皆阳
	{
		Entry: lookupCatalog(5),
		Matches: func(p *Pan) bool {
			cnt := 0
			zhis := collectAllZhi(p)
			for _, z := range zhis {
				if ZhiYang[z] {
					cnt++
				}
			}
			// 简化：三传 + 四课上神 + 日支 >= 7/8 皆阳
			return cnt >= len(zhis)-1
		},
	},
	// 第 6 条：六阴相继尽昏迷
	{
		Entry: lookupCatalog(6),
		Matches: func(p *Pan) bool {
			cnt := 0
			zhis := collectAllZhi(p)
			for _, z := range zhis {
				if !ZhiYang[z] {
					cnt++
				}
			}
			return cnt >= len(zhis)-1
		},
	},
	// 第 7 条：旺禄临身 —— 日干之禄神临干上神
	{
		Entry: lookupCatalog(7),
		Matches: func(p *Pan) bool {
			ganUpper := p.TianPan[GanJiGong[p.Ctx.Gan]]
			return ganUpper == GanLuZhi[p.Ctx.Gan]
		},
	},
	// 第 8 条：权摄不正禄临支 —— 日禄临日支之上
	{
		Entry: lookupCatalog(8),
		Matches: func(p *Pan) bool {
			zhiUpper := p.TianPan[p.Ctx.DayZhi]
			return zhiUpper == GanLuZhi[p.Ctx.Gan]
		},
	},
	// 第 11/13 条合并：三传皆作日鬼
	{
		Entry: lookupCatalog(11),
		Matches: func(p *Pan) bool {
			return allTraanAre(p, LQGuanGui)
		},
	},
	// 第 14 条：传财太旺 —— 三传皆妻财
	{
		Entry: lookupCatalog(14),
		Matches: func(p *Pan) bool { return allTraanAre(p, LQQiCai) },
	},
	// 第 15 条：脱上逢脱（日干生上神）—— 干上神为子孙
	{
		Entry: lookupCatalog(15),
		Matches: func(p *Pan) bool {
			ganUpper := p.TianPan[GanJiGong[p.Ctx.Gan]]
			return LiuQinOfZhiByGan(ganUpper, p.Ctx.Gan) == LQZiSun
		},
	},
	// 第 16 条：空上乘空 —— 干上神落旬空，乘天空
	{
		Entry: lookupCatalog(16),
		Matches: func(p *Pan) bool {
			ganUpper := p.TianPan[GanJiGong[p.Ctx.Gan]]
			if !IsXunKong(ganUpper, p.Ctx.JiaziIndex) {
				return false
			}
			return TianJiangOf(p.TianPan, p.TianJiang, ganUpper) == TJTianKong
		},
	},
	// 第 17 条：进茹空亡 —— 连茹顺行 + 三传皆空
	{
		Entry: lookupCatalog(17),
		Matches: func(p *Pan) bool {
			fw, _ := isLianRu(p.SanChuan.Chu.Zhi, p.SanChuan.Zhong.Zhi, p.SanChuan.Mo.Zhi)
			if !fw {
				return false
			}
			return p.SanChuan.Chu.IsKong && p.SanChuan.Zhong.IsKong && p.SanChuan.Mo.IsKong
		},
	},
	// 第 18 条：踏脚空亡 —— 连茹逆行 + 三传皆空
	{
		Entry: lookupCatalog(18),
		Matches: func(p *Pan) bool {
			_, bw := isLianRu(p.SanChuan.Chu.Zhi, p.SanChuan.Zhong.Zhi, p.SanChuan.Mo.Zhi)
			if !bw {
				return false
			}
			return p.SanChuan.Chu.IsKong && p.SanChuan.Zhong.IsKong && p.SanChuan.Mo.IsKong
		},
	},
	// 第 22 条：上下皆合 —— 干上神与日干六合；支上神与日支六合
	{
		Entry: lookupCatalog(22),
		Matches: func(p *Pan) bool {
			ganJi := GanJiGong[p.Ctx.Gan]
			ganUp := p.TianPan[ganJi]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			return isLiuHe(ganJi, ganUp) && isLiuHe(p.Ctx.DayZhi, zhiUp)
		},
	},
	// 第 23 条：彼求我事支传干 —— 初传=支上神，末传=干上神
	{
		Entry: lookupCatalog(23),
		Matches: func(p *Pan) bool {
			supper := p.TianPan[p.Ctx.DayZhi]
			gupper := p.TianPan[GanJiGong[p.Ctx.Gan]]
			return p.SanChuan.Chu.Zhi == supper && p.SanChuan.Mo.Zhi == gupper
		},
	},
	// 第 24 条：我求彼事干传支 —— 初传=干上神，末传=支上神
	{
		Entry: lookupCatalog(24),
		Matches: func(p *Pan) bool {
			gupper := p.TianPan[GanJiGong[p.Ctx.Gan]]
			supper := p.TianPan[p.Ctx.DayZhi]
			return p.SanChuan.Chu.Zhi == gupper && p.SanChuan.Mo.Zhi == supper
		},
	},
	// 第 31 条：三传递生 —— 五行相生
	{
		Entry: lookupCatalog(31),
		Matches: func(p *Pan) bool {
			c, z, m := p.SanChuan.Chu.Zhi, p.SanChuan.Zhong.Zhi, p.SanChuan.Mo.Zhi
			return WuXingGenerates(ZhiWuXing[c], ZhiWuXing[z]) && WuXingGenerates(ZhiWuXing[z], ZhiWuXing[m])
		},
	},
	// 第 32 条：三传互克 —— 五行相克
	{
		Entry: lookupCatalog(32),
		Matches: func(p *Pan) bool {
			c, z, m := p.SanChuan.Chu.Zhi, p.SanChuan.Zhong.Zhi, p.SanChuan.Mo.Zhi
			return WuXingOvercomes(ZhiWuXing[c], ZhiWuXing[z]) && WuXingOvercomes(ZhiWuXing[z], ZhiWuXing[m])
		},
	},
	// 第 40 条：后合占婚 —— 支干上神乘天后或六合
	{
		Entry: lookupCatalog(40),
		Matches: func(p *Pan) bool {
			return anyGanZhiUpperHasTianJiang(p, TJTianHou) || anyGanZhiUpperHasTianJiang(p, TJLiuHe)
		},
	},
	// 第 41 条：富贵干支逢禄马 —— 干上乘支之驿马，支上乘干之禄
	{
		Entry: lookupCatalog(41),
		Matches: func(p *Pan) bool {
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			return ganUp == yiMa(p.Ctx.DayZhi) && zhiUp == GanLuZhi[p.Ctx.Gan]
		},
	},
	// 第 44 条：课传俱贵 —— 四课三传皆乘贵人（简化：任一上神乘贵即未必；需多处乘贵）
	{
		Entry: lookupCatalog(44),
		Matches: func(p *Pan) bool {
			// 统计四课上神 + 三传 共 7 处乘贵人的次数
			cnt := 0
			for _, ke := range p.SiKe {
				if TianJiangOf(p.TianPan, p.TianJiang, ke.Upper) == TJGuiRen {
					cnt++
				}
			}
			for _, ce := range chuanList(p) {
				if ce.TianJiang == TJGuiRen {
					cnt++
				}
			}
			return cnt >= 3
		},
	},
	// 第 50 条：二贵皆空 —— 旦暮贵人位皆旬空
	{
		Entry: lookupCatalog(50),
		Matches: func(p *Pan) bool {
			day := GuiRenByGan[p.Ctx.Gan][0]
			night := GuiRenByGan[p.Ctx.Gan][1]
			return IsXunKong(day, p.Ctx.JiaziIndex) && IsXunKong(night, p.Ctx.JiaziIndex)
		},
	},
	// 第 51 条：魁度天门 —— 戌加亥发用
	{
		Entry: lookupCatalog(51),
		Matches: func(p *Pan) bool {
			return p.TianPan[Hai] == Xu && p.SanChuan.Chu.Zhi == Xu
		},
	},
	// 第 52 条：罡塞鬼户 —— 辰加寅
	{
		Entry: lookupCatalog(52),
		Matches: func(p *Pan) bool {
			return p.TianPan[Yin] == Chen
		},
	},
	// 第 55 条：天罗地网 —— 干上乘干寄宫前一辰；支上乘支前一辰
	{
		Entry: lookupCatalog(55),
		Matches: func(p *Pan) bool {
			ganJi := GanJiGong[p.Ctx.Gan]
			ganUp := p.TianPan[ganJi]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			return int(ganUp) == (int(ganJi)+1)%12 && int(zhiUp) == (int(p.Ctx.DayZhi)+1)%12
		},
	},
	// 第 63 条：彼此全伤 —— 干上神克日干 且 支上神克日支
	{
		Entry: lookupCatalog(63),
		Matches: func(p *Pan) bool {
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			return WuXingOvercomes(ZhiWuXing[ganUp], GanWuXing[p.Ctx.Gan]) &&
				WuXingOvercomes(ZhiWuXing[zhiUp], ZhiWuXing[p.Ctx.DayZhi])
		},
	},
	// 第 64 条：夫妇芜淫 —— 干克支上神 且 支克干上神
	{
		Entry: lookupCatalog(64),
		Matches: func(p *Pan) bool {
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			return WuXingOvercomes(GanWuXing[p.Ctx.Gan], ZhiWuXing[zhiUp]) &&
				WuXingOvercomes(ZhiWuXing[p.Ctx.DayZhi], ZhiWuXing[ganUp])
		},
	},
	// 第 74 条：三传皆空
	{
		Entry: lookupCatalog(74),
		Matches: func(p *Pan) bool {
			return p.SanChuan.Chu.IsKong && p.SanChuan.Zhong.IsKong && p.SanChuan.Mo.IsKong
		},
	},
	// 第 75 条：宾主不投刑在上 —— 干上支上见三刑
	{
		Entry: lookupCatalog(75),
		Matches: func(p *Pan) bool {
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			// 二字刑：子卯
			if (ganUp == Zi && zhiUp == Mao) || (ganUp == Mao && zhiUp == Zi) {
				return true
			}
			// 三字刑：寅巳申、丑戌未（三传成刑组）
			return isThreeXingSet(p.SanChuan.Chu.Zhi, p.SanChuan.Zhong.Zhi, p.SanChuan.Mo.Zhi)
		},
	},
	// 第 76 条：彼此猜忌害相随 —— 干支上神互为六害
	{
		Entry: lookupCatalog(76),
		Matches: func(p *Pan) bool {
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			return isLiuHai(ganUp, zhiUp)
		},
	},
	// 第 77 条：互生俱生 —— 干上神生支 且 支上神生干
	{
		Entry: lookupCatalog(77),
		Matches: func(p *Pan) bool {
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			return WuXingGenerates(ZhiWuXing[ganUp], ZhiWuXing[p.Ctx.DayZhi]) &&
				WuXingGenerates(ZhiWuXing[zhiUp], GanWuXing[p.Ctx.Gan])
		},
	},
	// 第 78 条：互旺皆旺 —— 干上乘干之帝旺，支上乘支之帝旺（简化）
	{
		Entry: lookupCatalog(78),
		Matches: func(p *Pan) bool {
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			return ganUp == GanWangZhi[p.Ctx.Gan]
		},
	},
	// 第 79 条：干支值绝 —— 干上乘干之绝地
	{
		Entry: lookupCatalog(79),
		Matches: func(p *Pan) bool {
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			return ganUp == GanJueZhi[p.Ctx.Gan]
		},
	},
	// 第 80 条：人宅皆死 —— 干上或支上乘日干之死气
	{
		Entry: lookupCatalog(80),
		Matches: func(p *Pan) bool {
			death := GanDeathZhi[GanWuXing[p.Ctx.Gan]]
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			return ganUp == death && zhiUp == death
		},
	},
	// 第 81 条：传墓入墓 —— 末传为日干之墓
	{
		Entry: lookupCatalog(81),
		Matches: func(p *Pan) bool {
			return p.SanChuan.Mo.Zhi == GanMuZhi[p.Ctx.Gan]
		},
	},
	// 第 82 条：不行传者考初时 —— 中传与末传俱空
	{
		Entry: lookupCatalog(82),
		Matches: func(p *Pan) bool {
			return p.SanChuan.Zhong.IsKong && p.SanChuan.Mo.IsKong && !p.SanChuan.Chu.IsKong
		},
	},
	// 第 83 条：三六合 —— 三合课 + 干支上另有字与三合中间一字六合
	{
		Entry: lookupCatalog(83),
		Matches: func(p *Pan) bool {
			if !isSanHeTrio(p.SanChuan.Chu.Zhi, p.SanChuan.Zhong.Zhi, p.SanChuan.Mo.Zhi) {
				return false
			}
			mid := p.SanChuan.Zhong.Zhi
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			return isLiuHe(mid, ganUp) || isLiuHe(mid, zhiUp)
		},
	},
	// 第 84 条：合中犯煞 —— 三合课 + 干支上有字与三合中一字作刑冲害
	{
		Entry: lookupCatalog(84),
		Matches: func(p *Pan) bool {
			if !isSanHeTrio(p.SanChuan.Chu.Zhi, p.SanChuan.Zhong.Zhi, p.SanChuan.Mo.Zhi) {
				return false
			}
			trio := []Zhi{p.SanChuan.Chu.Zhi, p.SanChuan.Zhong.Zhi, p.SanChuan.Mo.Zhi}
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			for _, t := range trio {
				for _, u := range []Zhi{ganUp, zhiUp} {
					if isLiuChong(t, u) || isLiuHai(t, u) {
						return true
					}
				}
			}
			return false
		},
	},
	// 第 85 条：初遭夹克 —— 初传被两侧（中传与干/支上神）所克（简化）
	{
		Entry: lookupCatalog(85),
		Matches: func(p *Pan) bool {
			chu := p.SanChuan.Chu.Zhi
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			attacked := 0
			for _, n := range []Zhi{ganUp, zhiUp, p.SanChuan.Zhong.Zhi} {
				if WuXingOvercomes(ZhiWuXing[n], ZhiWuXing[chu]) {
					attacked++
				}
			}
			return attacked >= 2
		},
	},
	// 第 89 条：任信丁马 —— 伏吟课且三传或干支上有驿马
	{
		Entry: lookupCatalog(89),
		Matches: func(p *Pan) bool {
			if p.SanChuan.Method != "伏吟法" {
				return false
			}
			ym := yiMa(p.Ctx.DayZhi)
			for _, ce := range chuanList(p) {
				if ce.Zhi == ym {
					return true
				}
			}
			return false
		},
	},
	// 第 90 条：来去俱空 —— 返吟且三传皆空
	{
		Entry: lookupCatalog(90),
		Matches: func(p *Pan) bool {
			if p.SanChuan.Method != "返吟法" {
				return false
			}
			return p.SanChuan.Chu.IsKong && p.SanChuan.Zhong.IsKong && p.SanChuan.Mo.IsKong
		},
	},
	// 第 91 条：虎临干鬼 —— 干上神为官鬼，且乘白虎
	{
		Entry: lookupCatalog(91),
		Matches: func(p *Pan) bool {
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			if LiuQinOfZhiByGan(ganUp, p.Ctx.Gan) != LQGuanGui {
				return false
			}
			return TianJiangOf(p.TianPan, p.TianJiang, ganUp) == TJBaiHu
		},
	},
	// 第 92 条：龙加生气 —— 青龙乘日干父母之神
	{
		Entry: lookupCatalog(92),
		Matches: func(p *Pan) bool {
			for _, ce := range chuanList(p) {
				if ce.TianJiang == TJQingLong && ce.LiuQin == LQFuMu {
					return true
				}
			}
			return false
		},
	},

	// 保留原 20 条里仍有价值的几条（部分已被合并，未重复）：
	// 第 5 条末传空（=毕法第 5 条末传逢空事难终；在文献中对应 ~第 16 条或 74 条细目）
	// 这里把"三传空亡一二条"的条目已映射到第 17/18/74/82 条。

	// 额外：初传逢空（非正式第 N 条，但常用）：并入第 16 条。

	// 第 25 条：金日逢丁 —— 庚辛日，三传含丁神
	{
		Entry: lookupCatalog(25),
		Matches: func(p *Pan) bool {
			if p.Ctx.Gan != Geng && p.Ctx.Gan != Xin {
				return false
			}
			for _, ce := range chuanList(p) {
				if IsDingShen(ce.Zhi, p.Ctx.JiaziIndex) {
					return true
				}
			}
			return false
		},
	},
	// 第 26 条：水日逢丁 —— 壬癸日，三传含丁神
	{
		Entry: lookupCatalog(26),
		Matches: func(p *Pan) bool {
			if p.Ctx.Gan != Ren && p.Ctx.Gan != Gui {
				return false
			}
			for _, ce := range chuanList(p) {
				if IsDingShen(ce.Zhi, p.Ctx.JiaziIndex) {
					return true
				}
			}
			return false
		},
	},
	// 第 69 条：虎乘遁鬼 —— 白虎所乘地支之遁干与日干相克（为鬼）
	{
		Entry: lookupCatalog(69),
		Matches: func(p *Pan) bool {
			// 查哪个天盘地支乘白虎
			for i, tj := range p.TianJiang {
				if tj != TJBaiHu {
					continue
				}
				hu := p.TianPan[i] // 白虎所乘地支
				dg := DunGan(hu, p.Ctx.JiaziIndex)
				if dg < 0 {
					continue
				}
				// 判断 dg 是否克日干（即为"鬼"）
				if WuXingOvercomes(GanWuXing[dg], GanWuXing[p.Ctx.Gan]) {
					return true
				}
			}
			return false
		},
	},

	// 第 53 条：两蛇夹墓 —— 丙戌日戌加巳（极特殊；这里用条件判）
	{
		Entry: lookupCatalog(53),
		Matches: func(p *Pan) bool {
			if p.Ctx.Gan != Bing || p.Ctx.DayZhi != Xu {
				return false
			}
			return p.TianPan[Si] == Xu
		},
	},
	// 第 61 条：干乘墓虎 —— 六辛日丑加戌且乘白虎
	{
		Entry: lookupCatalog(61),
		Matches: func(p *Pan) bool {
			if p.Ctx.Gan != Xin {
				return false
			}
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			if ganUp != GanMuZhi[Xin] {
				return false
			}
			return TianJiangOf(p.TianPan, p.TianJiang, ganUp) == TJBaiHu
		},
	},
	// 第 59 条：华盖覆日 —— 干上神为日干之墓
	{
		Entry: lookupCatalog(59),
		Matches: func(p *Pan) bool {
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			return ganUp == GanMuZhi[p.Ctx.Gan]
		},
	},
	// 第 9 条：避难逃生须弃旧 —— 三传皆不利（官鬼或空）而干上神吉
	{
		Entry: lookupCatalog(9),
		Matches: func(p *Pan) bool {
			// 三传全部为官鬼或空
			for _, ce := range chuanList(p) {
				if !ce.IsKong && ce.LiuQin != LQGuanGui {
					return false
				}
			}
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			lq := LiuQinOfZhiByGan(ganUp, p.Ctx.Gan)
			return lq == LQFuMu || lq == LQZiSun
		},
	},
	// 第 94 条：喜惧空亡 —— 盘中有任何字空亡（总提示）
	{
		Entry: lookupCatalog(94),
		Matches: func(p *Pan) bool {
			for _, ce := range chuanList(p) {
				if ce.IsKong {
					return true
				}
			}
			return false
		},
	},
	// 第 100 条：已灾凶逃返 —— 三传含官鬼 + 已有空亡解之（启示语）
	{
		Entry: lookupCatalog(100),
		Matches: func(p *Pan) bool {
			hasRui := false
			for _, ce := range chuanList(p) {
				if ce.LiuQin == LQGuanGui {
					hasRui = true
					break
				}
			}
			hasKong := false
			for _, ce := range chuanList(p) {
				if ce.IsKong {
					hasKong = true
					break
				}
			}
			return hasRui && hasKong
		},
	},

	// ========== 四期新增 10 条 ==========

	// 第 10 条：朽木难雕 —— 初传为卯且卯旬空
	{
		Entry: lookupCatalog(10),
		Matches: func(p *Pan) bool {
			return p.SanChuan.Chu.Zhi == Mao && IsXunKong(Mao, p.Ctx.JiaziIndex)
		},
	},
	// 第 19 条：胎财生气妻怀孕 —— 日干胎神作妻财
	{
		Entry: lookupCatalog(19),
		Matches: func(p *Pan) bool {
			tai := GanTaiZhi[p.Ctx.Gan]
			return LiuQinOfZhiByGan(tai, p.Ctx.Gan) == LQQiCai
		},
	},
	// 第 20 条：胎财死气损胎 —— 日干胎神作妻财，且胎神处死气
	{
		Entry: lookupCatalog(20),
		Matches: func(p *Pan) bool {
			tai := GanTaiZhi[p.Ctx.Gan]
			if LiuQinOfZhiByGan(tai, p.Ctx.Gan) != LQQiCai {
				return false
			}
			death := GanDeathZhi[GanWuXing[p.Ctx.Gan]]
			return tai == death
		},
	},
	// 第 42 条：尊崇传内遇三奇 —— 三传遁干全是甲戊庚 或 乙丙丁
	{
		Entry: lookupCatalog(42),
		Matches: func(p *Pan) bool {
			set := map[Gan]bool{}
			for _, ce := range chuanList(p) {
				dg := DunGan(ce.Zhi, p.Ctx.JiaziIndex)
				if dg < 0 {
					return false
				}
				set[dg] = true
			}
			// 甲戊庚三奇
			if set[Jia] && set[Wu1] && set[Geng] {
				return true
			}
			// 乙丙丁三奇
			if set[Yi] && set[Bing] && set[Ding] {
				return true
			}
			return false
		},
	},
	// 第 65 条：干墓并关 —— 日干之墓与四季关神同为初传
	{
		Entry: lookupCatalog(65),
		Matches: func(p *Pan) bool {
			monthZhi := monthZhiOf(p.Ctx)
			guan := guanShen(monthZhi)
			return p.SanChuan.Chu.Zhi == GanMuZhi[p.Ctx.Gan] &&
				p.SanChuan.Chu.Zhi == guan
		},
	},
	// 第 66 条：支坟财并旅程稽 —— 日支之墓恰为日干之财，且作初传
	{
		Entry: lookupCatalog(66),
		Matches: func(p *Pan) bool {
			mu := zhiMu(p.Ctx.DayZhi)
			if mu < 0 || mu == p.Ctx.DayZhi {
				return false
			}
			if LiuQinOfZhiByGan(mu, p.Ctx.Gan) != LQQiCai {
				return false
			}
			return p.SanChuan.Chu.Zhi == mu
		},
	},
	// 第 71 条：病符克宅 —— 太岁前一位为病符，临日支且克日支
	{
		Entry: lookupCatalog(71),
		Matches: func(p *Pan) bool {
			if p.Ctx.Lunar == nil {
				return false
			}
			// 太岁：当前干支年支
			yz := ParseZhi(p.Ctx.Lunar.GetYearZhiExact())
			if yz < 0 {
				return false
			}
			// 病符 = 旧太岁 = 前一年地支 = 太岁前一位
			bf := Zhi((int(yz) - 1 + 12) % 12)
			// 病符在天盘中所临地盘位，是否为日支位，且病符五行克日支五行
			for i, up := range p.TianPan {
				if up == bf && Zhi(i) == p.Ctx.DayZhi {
					return WuXingOvercomes(ZhiWuXing[bf], ZhiWuXing[p.Ctx.DayZhi])
				}
			}
			return false
		},
	},
	// 第 72 条：丧吊全逢 —— 岁前二辰(丧门)与岁后二辰(吊客)都落在干支上
	{
		Entry: lookupCatalog(72),
		Matches: func(p *Pan) bool {
			if p.Ctx.Lunar == nil {
				return false
			}
			yz := ParseZhi(p.Ctx.Lunar.GetYearZhiExact())
			if yz < 0 {
				return false
			}
			sangMen := Zhi((int(yz) + 2) % 12)   // 岁前二辰
			diaoKe := Zhi((int(yz) - 2 + 12) % 12) // 岁后二辰
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			hasSang := ganUp == sangMen || zhiUp == sangMen
			hasDiao := ganUp == diaoKe || zhiUp == diaoKe
			return hasSang && hasDiao
		},
	},
	// 第 87 条：人宅坐墓 —— 日干寄宫与日支的"地盘位"同时为自身之墓（即干支都落在墓神地盘位上）
	{
		Entry: lookupCatalog(87),
		Matches: func(p *Pan) bool {
			// 地盘是固定的子丑寅卯…，"坐于地盘墓上"指某神落在地盘辰/戌/丑/未位
			// 此处简化：日干寄宫地盘位本身是墓支，且日支地盘位也是墓支
			// 即 GanJiGong[Gan] 或 DayZhi 本身是 辰戌丑未
			isMu := func(z Zhi) bool { return z == Chen || z == Xu || z == Chou || z == Wei }
			// 古法："天盘干支皆坐于地盘墓上"——用干上神和支上神是否落在地盘墓位判
			for i, up := range p.TianPan {
				if up == GanJiGong[p.Ctx.Gan] && isMu(Zhi(i)) {
					for j, up2 := range p.TianPan {
						if up2 == p.Ctx.DayZhi && isMu(Zhi(j)) {
							return true
						}
					}
				}
			}
			return false
		},
	},
	// 第 88 条：干支乘墓 —— 日干上神为日干之墓；日支上神为日支之墓
	{
		Entry: lookupCatalog(88),
		Matches: func(p *Pan) bool {
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			zm := zhiMu(p.Ctx.DayZhi)
			return ganUp == GanMuZhi[p.Ctx.Gan] && zhiUp == zm && zm > 0
		},
	},

	// ---------- v3 P1 新增 9 条规则 ----------

	// 第 1 条：前后引从升迁吉 —— 初传居干/支前位 + 末传居干/支后位
	//   "前位" 含义：地支序号 -1（前）/ +1（后）
	{
		Entry: lookupCatalog(1),
		Matches: func(p *Pan) bool {
			gz := GanJiGong[p.Ctx.Gan]
			zz := p.Ctx.DayZhi
			c := p.SanChuan.Chu.Zhi
			m := p.SanChuan.Mo.Zhi
			isPrev := func(a, b Zhi) bool { return (int(a)+1)%12 == int(b) }
			isNext := func(a, b Zhi) bool { return (int(a)+11)%12 == int(b) }
			// 初传居干前/支前 + 末传居干后/支后
			chuQian := isPrev(c, gz) || isPrev(c, zz)
			moHou := isNext(m, gz) || isNext(m, zz)
			return chuQian && moHou
		},
	},

	// 第 2 条：首尾相见始终宜 —— 干上有旬尾、支上有旬首（或反之）
	//   旬首即旬开始的地支（如甲子旬首为子），旬尾即旬最后地支（甲子旬尾为酉）
	{
		Entry: lookupCatalog(2),
		Matches: func(p *Pan) bool {
			xunStart := (p.Ctx.JiaziIndex / 10) * 10
			xunHead := Zhi(xunStart % 12)        // 旬首地支
			xunTail := Zhi((xunStart + 9) % 12)  // 旬尾地支
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			return (ganUp == xunTail && zhiUp == xunHead) ||
				(ganUp == xunHead && zhiUp == xunTail)
		},
	},

	// 第 3 条：帘幕贵人高甲第 —— 昼占夜贵 / 夜占昼贵 临年命或日干
	{
		Entry: lookupCatalog(3),
		Matches: func(p *Pan) bool {
			// 反贵：昼占取夜贵、夜占取昼贵
			var fanGui Zhi
			if p.Ctx.ZhouYe {
				fanGui = GuiRenByGan[p.Ctx.Gan][1]
			} else {
				fanGui = GuiRenByGan[p.Ctx.Gan][0]
			}
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			if ganUp == fanGui {
				return true
			}
			// 临年命
			if p.NianMing != nil && p.NianMing.BenMing != nil && p.NianMing.BenMing.Upper == fanGui {
				return true
			}
			if p.NianMing != nil && p.NianMing.XingNian != nil && p.NianMing.XingNian.Upper == fanGui {
				return true
			}
			return false
		},
	},

	// 第 4 条：催官使者赴官期 —— 日鬼乘白虎 + 临日干或年命
	//   "日鬼" = 克日干的地支
	{
		Entry: lookupCatalog(4),
		Matches: func(p *Pan) bool {
			isGui := func(z Zhi) bool {
				return LiuQinOfZhiByGan(z, p.Ctx.Gan) == LQGuanGui
			}
			// 取乘白虎的天盘地支
			var hu Zhi = -1
			for i, t := range p.TianJiang {
				if t == TJBaiHu {
					hu = p.TianPan[i]
					break
				}
			}
			if hu < 0 || !isGui(hu) {
				return false
			}
			// 临日干（即干上神 = 该地支）
			if p.TianPan[GanJiGong[p.Ctx.Gan]] == hu {
				return true
			}
			// 临年命
			if p.NianMing != nil {
				if p.NianMing.BenMing != nil && p.NianMing.BenMing.Upper == hu {
					return true
				}
				if p.NianMing.XingNian != nil && p.NianMing.XingNian.Upper == hu {
					return true
				}
			}
			return false
		},
	},

	// 第 12 条：狐假虎威仪 —— 日干被克但日支克其鬼神（支救干）
	{
		Entry: lookupCatalog(12),
		Matches: func(p *Pan) bool {
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			// 日干被上神克
			if !WuXingOvercomes(ZhiWuXing[ganUp], GanWuXing[p.Ctx.Gan]) {
				return false
			}
			// 日支克日干上神
			return WuXingOvercomes(ZhiWuXing[p.Ctx.DayZhi], ZhiWuXing[ganUp])
		},
	},

	// 第 13 条：鬼贼当时无畏忌 —— 三传作日鬼但本月令日干当令
	//   即三传皆为官鬼 + 月令五行与日干同
	{
		Entry: lookupCatalog(13),
		Matches: func(p *Pan) bool {
			if !allTraanAre(p, LQGuanGui) {
				return false
			}
			if p.Ctx.Lunar == nil {
				return false
			}
			mz := monthZhiOf(p.Ctx)
			if mz < 0 {
				return false
			}
			return ZhiWuXing[mz] == GanWuXing[p.Ctx.Gan]
		},
	},

	// 第 21 条：交车相合交关利 —— 一三课交叉作六合
	//   即 一课上 与 三课下 合 + 一课下 与 三课上 合
	{
		Entry: lookupCatalog(21),
		Matches: func(p *Pan) bool {
			if len(p.SiKe) < 4 {
				return false
			}
			k1 := p.SiKe[0]
			k3 := p.SiKe[2]
			return zhiLiuhe(k1.Upper, k3.Lower) && zhiLiuhe(k1.Lower, k3.Upper)
		},
	},

	// 第 27 条：富贵未及催官 —— 干支两贵 + 三传贵人临值
	{
		Entry: lookupCatalog(27),
		Matches: func(p *Pan) bool {
			day := GuiRenByGan[p.Ctx.Gan][0]
			nig := GuiRenByGan[p.Ctx.Gan][1]
			ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
			zhiUp := p.TianPan[p.Ctx.DayZhi]
			twoGui := (ganUp == day || ganUp == nig) && (zhiUp == day || zhiUp == nig)
			if !twoGui {
				return false
			}
			// 三传中至少一传乘贵人
			for _, ce := range chuanList(p) {
				if ce.TianJiang == TJGuiRen {
					return true
				}
			}
			return false
		},
	},

	// 第 30 条：进退维谷难抉择 —— 三传连茹 + 遁干自冲
	//   实务中"进退维谷"近义判定：三传顺连 + 末传冲日干寄宫
	{
		Entry: lookupCatalog(30),
		Matches: func(p *Pan) bool {
			fw, _ := isLianRu(p.SanChuan.Chu.Zhi, p.SanChuan.Zhong.Zhi, p.SanChuan.Mo.Zhi)
			if !fw {
				return false
			}
			gz := GanJiGong[p.Ctx.Gan]
			return (int(p.SanChuan.Mo.Zhi)-int(gz)+12)%12 == 6
		},
	},
}

// zhiLiuhe 地支六合：子丑、寅亥、卯戌、辰酉、巳申、午未
func zhiLiuhe(a, b Zhi) bool {
	pairs := [][2]Zhi{
		{Zi, Chou}, {Yin, Hai}, {Mao, Xu},
		{Chen, You}, {Si, Shen}, {Wu, Wei},
	}
	for _, p := range pairs {
		if (a == p[0] && b == p[1]) || (a == p[1] && b == p[0]) {
			return true
		}
	}
	return false
}

// ---------- 辅助工具 ----------

// chuanList 三传列表
func chuanList(p *Pan) []ChuanEntry {
	return []ChuanEntry{p.SanChuan.Chu, p.SanChuan.Zhong, p.SanChuan.Mo}
}

// collectAllZhi 搜集盘中所有地支（日支、干上、支上、四课上下、三传）
func collectAllZhi(p *Pan) []Zhi {
	zs := []Zhi{p.Ctx.DayZhi, GanJiGong[p.Ctx.Gan]}
	for _, ke := range p.SiKe {
		zs = append(zs, ke.Upper, ke.Lower)
	}
	for _, ce := range chuanList(p) {
		zs = append(zs, ce.Zhi)
	}
	return zs
}

// allTraanAre 三传六亲皆为某类
func allTraanAre(p *Pan, lq LiuQin) bool {
	for _, ce := range chuanList(p) {
		if ce.LiuQin != lq {
			return false
		}
	}
	return true
}

// anyGanZhiUpperHasTianJiang 干/支上神是否乘某天将
func anyGanZhiUpperHasTianJiang(p *Pan, tj TianJiang) bool {
	ganUp := p.TianPan[GanJiGong[p.Ctx.Gan]]
	zhiUp := p.TianPan[p.Ctx.DayZhi]
	return TianJiangOf(p.TianPan, p.TianJiang, ganUp) == tj ||
		TianJiangOf(p.TianPan, p.TianJiang, zhiUp) == tj
}

// lookupCatalog 根据序号查找目录条目
func lookupCatalog(n int) BiFaEntry {
	for _, e := range biFaCatalog {
		if e.Number == n {
			return e
		}
	}
	return BiFaEntry{Number: n, Title: "(未知)"}
}
