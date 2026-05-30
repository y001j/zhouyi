package main

import (
	"fmt"
	"strings"

	"zhouyi/liuren"
)

// QuestionType 问题类型
type QuestionType string

const (
	QTCareer   QuestionType = "career"   // 事业/工作
	QTWealth   QuestionType = "wealth"   // 财运/投资
	QTRelation QuestionType = "relation" // 感情/姻缘
	QTHealth   QuestionType = "health"   // 健康
	QTDecision QuestionType = "decision" // 抉择/选择
	QTTiming   QuestionType = "timing"   // 时运/吉凶
	QTOther    QuestionType = "other"    // 其他/不指定
)

// QuestionTypeLabel 问题类型的中文名
func QuestionTypeLabel(q QuestionType) string {
	return map[QuestionType]string{
		QTCareer:   "事业 / 工作",
		QTWealth:   "财运 / 投资",
		QTRelation: "感情 / 姻缘 / 人际",
		QTHealth:   "健康 / 身心",
		QTDecision: "抉择 / 两难",
		QTTiming:   "时运 / 整体吉凶",
		QTOther:    "其他",
	}[q]
}

// ParseQuestionType 从用户输入解析问题类型
func ParseQuestionType(s string) QuestionType {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "1", "career", "事业", "工作", "职业":
		return QTCareer
	case "2", "wealth", "财", "财运", "投资", "钱":
		return QTWealth
	case "3", "relation", "感情", "姻缘", "人际", "关系":
		return QTRelation
	case "4", "health", "健康", "身体", "病":
		return QTHealth
	case "5", "decision", "抉择", "选择", "决策":
		return QTDecision
	case "6", "timing", "时运", "吉凶", "运势":
		return QTTiming
	default:
		return QTOther
	}
}

// QuestionTypeMenu 问题类型菜单（用于交互提示）
func QuestionTypeMenu() string {
	var b strings.Builder
	b.WriteString("  请选择问题类型（可直接回车跳过）：\n")
	b.WriteString("    1) 事业 / 工作\n")
	b.WriteString("    2) 财运 / 投资\n")
	b.WriteString("    3) 感情 / 姻缘 / 人际\n")
	b.WriteString("    4) 健康 / 身心\n")
	b.WriteString("    5) 抉择 / 两难\n")
	b.WriteString("    6) 时运 / 整体吉凶\n")
	b.WriteString("    7) 其他 / 不指定\n")
	return b.String()
}

// FocusGuide 依问题类型返回**周易解卦**取象侧重（给 AI 看的）
//
// 注：六壬有独立的 FocusGuideLiuRen，奇门有独立的 qimen.FocusGuide；
// 主包按术数体系分别调用，避免侧重语言串线。
func FocusGuide(q QuestionType) string {
	switch q {
	case QTCareer:
		return "【取象侧重 · 事业】\n" +
			"- 重点关注：九五爻（君位，象征上级与机遇）、九三爻（进退之际，职位转折）、上下卦之和（内外环境）\n" +
			"- 关注刚柔得位（阳爻居位则能任事）、中位得正（二、五）则遇明主或居其职\n" +
			"- 震巽为动、乾为决断、坤为积累、坎为险阻、艮为止守——依卦象判断宜进宜退\n" +
			"- 变爻若在五爻：关系到权位与决策；若在三四爻：关系到人际或晋升\n"
	case QTWealth:
		return "【取象侧重 · 财运】\n" +
			"- 重点关注：本卦所象之财性（如泰则流通、否则阻塞、损益直接关乎得失、鼎革关乎变局）\n" +
			"- 卦象里坎为险/为水流（财如流水）、兑为悦/为口舌、艮为止（财停）、震为动（财动）\n" +
			"- 老阳变阴多主「由盛转衰宜收」、老阴变阳多主「由困转通宜进」\n" +
			"- 判断：宜守还是宜进？财的性质是投机、稳健、还是合作分润？风险点在哪一爻？\n"
	case QTRelation:
		return "【取象侧重 · 感情/人际】\n" +
			"- 重点关注：二爻与五爻的应与否（阴阳相应主两情相通，同性敌应主不和）\n" +
			"- 内卦为自己、外卦为对方；上下卦相合（如水火既济、天地交泰）主和谐，相背（如天水违行）主分歧\n" +
			"- 咸（感）、恒（久）、家人、睽（乖）、革（变）、归妹等卦在此类问题中常有直接启示\n" +
			"- 变爻的位置暗示转变由谁发起（内卦变=自己主动，外卦变=对方主动）\n"
	case QTHealth:
		return "【取象侧重 · 健康】\n" +
			"- 八卦所主身部（《说卦传》）：乾为首、坤为腹、震为足、巽为股、坎为耳、离为目、艮为手、兑为口\n" +
			"- 八卦所主脏腑（后天五行配属）：乾/兑属金主肺与大肠、坤/艮属土主脾胃、震/巽属木主肝胆、坎属水主肾与膀胱、离属火主心与小肠\n" +
			"- 变爻位置提示病位：初爻足下、二爻腹、三爻腰腹之交、四爻胸胁、五爻心胸头颈、上爻头面\n" +
			"- 卦象吉凶结合节气判断：逆时令者病重，顺时令者病轻\n" +
			"- 卦象推病性：坎主寒/湿/血、离主热/炎/心神、震主惊/筋、巽主风/气、艮主积滞、兑主口舌/肺燥\n"
	case QTDecision:
		return "【取象侧重 · 抉择】\n" +
			"- 重点关注：本卦与变卦之差——当前局势 vs 若行某选项后的局势\n" +
			"- 若有多个变爻：三爻变以上转折大，一二爻变则调整幅度有限\n" +
			"- 互卦揭示选项背后潜藏的真实走向（常与表面意图不同）\n" +
			"- 错卦给出「若选相反方向」的参考\n" +
			"- 最终回答必须明确给出倾向（不要模棱两可），并说明代价与时机\n"
	case QTTiming:
		return "【取象侧重 · 时运】\n" +
			"- 重点关注：本卦与时令消息卦的关系（相合则顺时、相悖则逆时）\n" +
			"- 卦序前后邻卦提示「从何处来、向何处去」的运势脉络\n" +
			"- 六爻整体阴阳配比：阳多主外显奋进之时，阴多主潜藏积累之时\n" +
			"- 变爻指出运势在哪一阶段将出现转折\n"
	default:
		return "【取象侧重 · 通用】\n" +
			"- 结合卦辞、象辞、爻辞的整体意象解读，不拘一派\n" +
			"- 注意时令、爻位、变爻的综合作用\n"
	}
}

// leishenSpec 一条类神规格：按天将或按六亲匹配
type leishenSpec struct {
	Label   string          // 人读标签：如 "求财类神"
	TianJ   []liuren.TianJiang // 天将（任一命中即算）
	LiuQin  []liuren.LiuQin    // 六亲（任一命中即算）
	Note    string          // 补充说明
}

// leishenSpecsByQType 按问题类型返回需查找的类神规格组
//
// 据《六壬大全》卷二天将释 + 卷六占类章 + 毕法赋。每类给 6-8 条主辅类神，
// 覆盖：主类神（事之本体）、救应类神（顺成之机）、阻力类神（阻碍之象）、
//      六亲（按问题类型选关键 1-2 个）。
func leishenSpecsByQType(q QuestionType) []leishenSpec {
	switch q {
	case QTCareer:
		return []leishenSpec{
			{Label: "上级/赏识（贵人）", TianJ: []liuren.TianJiang{liuren.TJGuiRen}, Note: "天乙主上司赏识、扶助；登天门（亥）则吉中之吉，临辰戌（地狱/天牢）失位"},
			{Label: "职位/名位（太常）", TianJ: []liuren.TianJiang{liuren.TJTaiChang}, Note: "太常主衣冠/俸禄/官职稳固；带印信则受职"},
			{Label: "文书/考核（朱雀）", TianJ: []liuren.TianJiang{liuren.TJZhuQue}, Note: "朱雀主文书/捷报/考核公告；临巳午得地，临亥子失明"},
			{Label: "官鬼（职位/职责本体）", LiuQin: []liuren.LiuQin{liuren.LQGuanGui}, Note: "官鬼得地则职事稳；落空则虚名；与贵人合更佳"},
			{Label: "父母（印绶/任命书）", LiuQin: []liuren.LiuQin{liuren.LQFuMu}, Note: "父母为印绶，主任命/调令；空则虚信、入墓则压在抽屉"},
			{Label: "牵制/阻力（勾陈）", TianJ: []liuren.TianJiang{liuren.TJGouChen}, Note: "勾陈主官非牵连/审批拖延；持印则文书阻；入三传主事多波折"},
			{Label: "罢免/刑罚（白虎）", TianJ: []liuren.TianJiang{liuren.TJBaiHu}, Note: "白虎乘官鬼临年命 → 催官使者（赴任之期）；乘官鬼临日干 → 主刑罚加身"},
			{Label: "暗算/陷害（玄武）", TianJ: []liuren.TianJiang{liuren.TJXuanWu}, Note: "玄武入三传主小人暗陷；与朱雀合则文书有诈"},
		}
	case QTWealth:
		return []leishenSpec{
			{Label: "正财/喜庆（青龙）", TianJ: []liuren.TianJiang{liuren.TJQingLong}, Note: "青龙主正财/升迁之财；临寅卯木地大吉，临申酉折角失财"},
			{Label: "守财/产业（太常）", TianJ: []liuren.TianJiang{liuren.TJTaiChang}, Note: "太常主田园产业/俸禄之财；临丑未得位"},
			{Label: "妻财（所求之财本体）", LiuQin: []liuren.LiuQin{liuren.LQQiCai}, Note: "妻财得地有气则财可得；空则不实；入墓则财困；与日干合则财就我"},
			{Label: "禄神（俸禄/正财）", TianJ: []liuren.TianJiang{liuren.TJTaiChang, liuren.TJGuiRen}, Note: "禄神临干上 → 旺禄临身（毕法第7：宜守不可妄求）"},
			{Label: "盗贼/暗耗（玄武）", TianJ: []liuren.TianJiang{liuren.TJXuanWu}, Note: "玄武入三传主财被暗夺；临子（散发）盗机大盛"},
			{Label: "虚诈/空耗（天空）", TianJ: []liuren.TianJiang{liuren.TJTianKong}, Note: "天空乘财 → 财在指空划空；与脱气合则虚诈不实"},
			{Label: "破财/血光（白虎）", TianJ: []liuren.TianJiang{liuren.TJBaiHu}, Note: "白虎乘财临年命 → 财来带血光（如药费、官非赔偿）"},
			{Label: "兄弟（争夺/分润）", LiuQin: []liuren.LiuQin{liuren.LQXiongDi}, Note: "兄弟入用主人多分财；与劫煞合则被夺"},
		}
	case QTRelation:
		return []leishenSpec{
			{Label: "媒合/缔交（六合）", TianJ: []liuren.TianJiang{liuren.TJLiuHe}, Note: "六合为正媒之神；不乘丑未（古法禁忌）；入三传主婚易成"},
			{Label: "女方/正配（天后）", TianJ: []liuren.TianJiang{liuren.TJTianHou}, Note: "天后归亥（升殿）大吉；临辰巳午失位；临子（沐浴）淫泆之兆"},
			{Label: "男方/财礼（青龙）", TianJ: []liuren.TianJiang{liuren.TJQingLong}, Note: "青龙乘财临支 → 男方诚意/财礼之喜"},
			{Label: "暗助/私情（太阴）", TianJ: []liuren.TianJiang{liuren.TJTaiYin}, Note: "太阴主暗助、女子真情；与玄武合则阴私不正"},
			{Label: "淫奔/不明（玄武）", TianJ: []liuren.TianJiang{liuren.TJXuanWu}, Note: "玄武入婚 → 暗昧不明、可能桃色；与天后合则淫泆课"},
			{Label: "口舌/争议（朱雀）", TianJ: []liuren.TianJiang{liuren.TJZhuQue}, Note: "朱雀入婚 → 言语口角；乘鬼则婚约因文书生争"},
			{Label: "婚约阻力（勾陈）", TianJ: []liuren.TianJiang{liuren.TJGouChen}, Note: "勾陈主婚约牵连、长辈阻挠"},
			{Label: "父母（家长意见/婚书）", LiuQin: []liuren.LiuQin{liuren.LQFuMu}, Note: "父母为名分之证；空则婚书未立或长辈不准"},
			{Label: "妻财（妻位/聘礼）", LiuQin: []liuren.LiuQin{liuren.LQQiCai}, Note: "占夫则看官鬼，占妻则看妻财；妻财得地则婚姻和顺"},
		}
	case QTHealth:
		return []leishenSpec{
			{Label: "病邪本体（官鬼）", LiuQin: []liuren.LiuQin{liuren.LQGuanGui}, Note: "官鬼为病；旺则病重，墓则病臥，空则病退"},
			{Label: "病符/血光（白虎）", TianJ: []liuren.TianJiang{liuren.TJBaiHu}, Note: "白虎乘鬼临身 → 重病；白虎入墓 → 病入膏肓"},
			{Label: "惊悸/怪梦（腾蛇）", TianJ: []liuren.TianJiang{liuren.TJTengShe}, Note: "腾蛇入用 → 惊恐怪异/夜梦不安"},
			{Label: "药石/食疗（太常）", TianJ: []liuren.TianJiang{liuren.TJTaiChang}, Note: "太常主药石/食疗；临三传则有医方可寻"},
			{Label: "医者/扶助（贵人）", TianJ: []liuren.TianJiang{liuren.TJGuiRen}, Note: "贵人主医者；制鬼之位即良医之方（《大全》卷十）"},
			{Label: "丧讯/外丧（白虎+丧门吊客）", TianJ: []liuren.TianJiang{liuren.TJBaiHu}, Note: "白虎临年命+丧门 → 家有外丧；吊客带白虎 → 主丧门入宅"},
			{Label: "墓神（拘困/重病）", TianJ: []liuren.TianJiang{liuren.TJTaiChang}, Note: "辰戌丑未为五墓；病入墓神则臥床难起"},
			{Label: "子孙（解病/痊愈）", LiuQin: []liuren.LiuQin{liuren.LQZiSun}, Note: "子孙克鬼为药；子孙得地有气 → 病可解"},
		}
	case QTDecision:
		return []leishenSpec{
			{Label: "贵人指引", TianJ: []liuren.TianJiang{liuren.TJGuiRen}, Note: "贵人临三传/年命 → 有明人可请教"},
			{Label: "文书启示（朱雀）", TianJ: []liuren.TianJiang{liuren.TJZhuQue}, Note: "朱雀临干 → 决策有文书凭据"},
			{Label: "牵连阻力（勾陈）", TianJ: []liuren.TianJiang{liuren.TJGouChen}, Note: "勾陈主牵缠/拖延；入三传则决策受制"},
			{Label: "疑虑动摇（腾蛇）", TianJ: []liuren.TianJiang{liuren.TJTengShe}, Note: "腾蛇主反复犹疑；与朱雀合则言不一致"},
			{Label: "官鬼（事之压力）", LiuQin: []liuren.LiuQin{liuren.LQGuanGui}, Note: "官鬼临身 → 决策受外压；克日则被动"},
			{Label: "父母（凭据/规则）", LiuQin: []liuren.LiuQin{liuren.LQFuMu}, Note: "父母主决策依据；空则无凭"},
			{Label: "子孙（解忧/退路）", LiuQin: []liuren.LiuQin{liuren.LQZiSun}, Note: "子孙得地 → 有退路可循"},
		}
	case QTTiming:
		return []leishenSpec{
			{Label: "吉将组", TianJ: []liuren.TianJiang{liuren.TJQingLong, liuren.TJLiuHe, liuren.TJTaiChang, liuren.TJTianHou, liuren.TJGuiRen, liuren.TJTaiYin}, Note: "六吉将入三传越多越顺"},
			{Label: "凶将组", TianJ: []liuren.TianJiang{liuren.TJBaiHu, liuren.TJXuanWu, liuren.TJGouChen, liuren.TJTengShe, liuren.TJZhuQue, liuren.TJTianKong}, Note: "六凶将入三传越多越阻"},
			{Label: "妻财（所欲之物）", LiuQin: []liuren.LiuQin{liuren.LQQiCai}, Note: "妻财得地 → 欲望可遂"},
			{Label: "官鬼（外压）", LiuQin: []liuren.LiuQin{liuren.LQGuanGui}, Note: "官鬼克日 → 外压加身"},
			{Label: "父母（庇护）", LiuQin: []liuren.LiuQin{liuren.LQFuMu}, Note: "父母生日 → 得庇护"},
			{Label: "子孙（解厄）", LiuQin: []liuren.LiuQin{liuren.LQZiSun}, Note: "子孙制鬼 → 解厄之机"},
		}
	}
	return nil
}

// LeiShenDirective 按问题类型生成"类神直指"块：列出每种类神在本盘的落点、乘神、空亡与否
func LeiShenDirective(pan *liuren.Pan, q QuestionType) string {
	specs := leishenSpecsByQType(q)
	if pan == nil || len(specs) == 0 {
		return ""
	}
	kongPair := pan.Ctx.XunKongPair()
	isKong := func(z liuren.Zhi) bool { return z == kongPair[0] || z == kongPair[1] }

	// 准备：三传所占的地支集合
	inSanChuan := map[liuren.Zhi]string{
		pan.SanChuan.Chu.Zhi:   pan.SanChuan.Chu.Name,
		pan.SanChuan.Zhong.Zhi: pan.SanChuan.Zhong.Name,
		pan.SanChuan.Mo.Zhi:    pan.SanChuan.Mo.Name,
	}

	var b strings.Builder
	for _, spec := range specs {
		var hits []string

		// 按天将查（在 12 宫位里找天将落点，天盘上对应的地支）
		for _, tj := range spec.TianJ {
			for i := 0; i < 12; i++ {
				if pan.TianJiang[i] == tj {
					upper := pan.TianPan[i]
					kong := ""
					if isKong(upper) {
						kong = " **空亡**"
					}
					chuan := ""
					if name, ok := inSanChuan[upper]; ok {
						chuan = fmt.Sprintf(" · **入%s**", name)
					}
					hits = append(hits, fmt.Sprintf("%s 临地盘 %s 宫（天盘 %s）%s%s",
						tj, liuren.Zhi(i), upper, kong, chuan))
					break // 每个天将只占一位
				}
			}
		}

		// 按六亲查（在天盘 12 地支里找与日干构成此六亲的地支）
		for _, lq := range spec.LiuQin {
			for i := 0; i < 12; i++ {
				upper := pan.TianPan[i]
				if liuren.LiuQinOfZhiByGan(upper, pan.Ctx.Gan) == lq {
					kong := ""
					if isKong(upper) {
						kong = " **空亡**"
					}
					chuan := ""
					if name, ok := inSanChuan[upper]; ok {
						chuan = fmt.Sprintf(" · **入%s**", name)
					}
					tj := pan.TianJiang[i]
					hits = append(hits, fmt.Sprintf("%s（%s）临地盘 %s 宫 · 乘 %s%s%s",
						lq, upper, liuren.Zhi(i), tj, kong, chuan))
				}
			}
		}

		if len(hits) == 0 {
			b.WriteString(fmt.Sprintf("- **%s**：本盘未现 —— %s\n", spec.Label, spec.Note))
		} else {
			b.WriteString(fmt.Sprintf("- **%s**：%s —— %s\n",
				spec.Label, strings.Join(hits, "；"), spec.Note))
		}
	}

	return b.String()
}

// FormatQuestionTypeSection 生成问题类型与侧重的提示词段落
func FormatQuestionTypeSection(q QuestionType) string {
	if q == "" || q == QTOther {
		return "## 解卦侧重\n" + FocusGuide(QTOther) + "\n"
	}
	var b strings.Builder
	b.WriteString("## 问题类型\n")
	b.WriteString(fmt.Sprintf("%s\n\n", QuestionTypeLabel(q)))
	b.WriteString("## 解卦侧重\n")
	b.WriteString(FocusGuide(q))
	b.WriteString("\n")
	return b.String()
}

// FocusGuideLiuRen 依问题类型返回**大六壬断课**取象侧重（给 AI 看的）
//
// 与周易的 FocusGuide 同构但语言完全六壬化：
// 围绕「日干（我）/ 日支（事）/ 三传 / 天将 / 类神 / 神煞 / 年命」展开，
// 不出现"九五爻、变爻、卦象"等周易概念。
func FocusGuideLiuRen(q QuestionType) string {
	switch q {
	case QTCareer:
		return "【六壬断课侧重 · 事业 / 工作】\n" +
			"- 主类神：**官鬼**（职位/职责）、**贵人**（上级/赏识）、**太常**（衣冠/名位）、**朱雀**（公文/考核）\n" +
			"- 阻力类神：**勾陈**（公门牵制）、**腾蛇**（猜疑/虚惊）、**白虎**（罢免/刑罚）、**玄武**（暗中陷害）\n" +
			"- 重点关注：日干（我）是否得日上神之生扶？官鬼是否得地有气、入三传？是否落空（职位虚浮）？\n" +
			"- 关键格局：**催官使者**（日鬼乘虎临干/年命，赴任之期）、**帘幕贵人**（昼占夜贵或反之，临年命主登科）、**荣华课**（禄马贵人临干上）\n" +
			"- 升迁判断：年命上神生日干 → 救应；行年/本命入三传与吉将合 → 高升；空亡乘官 → 虚名\n"
	case QTWealth:
		return "【六壬断课侧重 · 财运 / 投资】\n" +
			"- 主类神：**妻财**（六亲 · 我所求）、**青龙**（正财/喜庆）、**太常**（守财/产业）、**禄神**（俸禄）\n" +
			"- 破财类神：**玄武**（盗贼/暗耗）、**天空**（虚诈/欺骗）、**白虎**（破财/血光之耗）、**劫煞**（劫夺）\n" +
			"- 重点关注：财爻是否得地有气、是否入三传、是否落空？日干是否克财？财乘吉将/凶将？\n" +
			"- 关键格局：**财生气**（妻财乘月内生气，财在路）、**繁昌课**（旺相之德临命发用）、**传财太旺反财亏**（毕法第14：三传皆财而日干休囚）\n" +
			"- 投资判断：财得月将之气、初传财、末传归我 → 大吉；财空、墓库、被克 → 不可妄进\n"
	case QTRelation:
		return "【六壬断课侧重 · 感情 / 姻缘 / 人际】\n" +
			"- 主类神：**六合**（媒合/缔结）、**天后**（女方/正配）、**青龙**（男方/财礼）、**神后**（女方真情，子位）\n" +
			"- 阴私类神：**太阴**（暗助/私情）、**玄武**（不明/淫奔）、**朱雀**（口舌/争议）\n" +
			"- 阻力类神：**勾陈**（婚约阻挠）、**腾蛇**（疑虑/反复）、**月厌丁马**（破婚之煞）\n" +
			"- 重点关注：日干日支是否相生相合？六合是否入三传？天后/青龙是否乘旺？是否见冲破？\n" +
			"- 关键格局：**合欢课**（干支两六合）、**淫泆课**（六合天后入用，主男女私情）、**龙虎交战课**（青龙白虎相加，主夫妻反目）\n" +
			"- 应期判断：天喜+月将的支位 → 喜期；六合落空 → 媒合不成；孤辰寡宿入命 → 失偶之兆\n"
	case QTHealth:
		return "【六壬断课侧重 · 健康 / 身心】\n" +
			"- 主类神：**官鬼**（病邪本身）、**白虎**（病符/血光/伤损）、**腾蛇**（惊悸/虚火/怪梦）、**死气死神**（病重之兆）\n" +
			"- 救应类神：**太常**（药石/食疗）、**贵人**（医者/扶助）、**生气**（病解之机）、**天医**（医药）\n" +
			"- 重点关注：年命上神是最终救应之主；官鬼是否得地、入墓、入三传、落空？白虎乘何、临何宫？\n" +
			"- 病位推断：白虎临何地支即病在何方（《大全》卷二白虎章 12 宫断），三传五行决定病性\n" +
			"- 关键格局：**虎乘鬼**（白虎+官鬼临三传，重病）、**鬼空亡**（病退之机）、**墓神临身**（病臥难起）、**丧门吊客**（家有外丧讯）\n" +
			"- 病解判断：年命见生气、天医、贵人 → 可救；年命见五墓、白虎、死气 → 大凶\n"
	case QTDecision:
		return "【六壬断课侧重 · 抉择 / 两难】\n" +
			"- 主类神：**贵人**（指引方向）、**朱雀**（文书启示）、**勾陈**（牵连阻力）、**腾蛇**（疑虑动摇）\n" +
			"- 重点关注：初传之因（事起何方）、中传之经过（中途变化）、末传之结局（最终归宿）\n" +
			"- 三传流向：递生（顺成）/ 递克（节节阻）/ 先生后克（虎头蛇尾）/ 先克后生（否极泰来）\n" +
			"- 关键格局：**进茹空亡宜退步**（毕法17：三传顺连皆空 → 应退）、**踏脚空亡进用宜**（毕法18：三传逆连皆空 → 应进）\n" +
			"- 决策判断：年命入三传与吉将合 → 当行；末传乘凶将且空亡 → 当弃；伏吟主静守，返吟主反复\n"
	case QTTiming:
		return "【六壬断课侧重 · 时运 / 整体吉凶】\n" +
			"- 吉将组：**贵人 / 青龙 / 六合 / 太常 / 太阴 / 天后** —— 入三传越多越顺\n" +
			"- 凶将组：**白虎 / 玄武 / 勾陈 / 腾蛇 / 朱雀 / 天空** —— 入三传越多越阻\n" +
			"- 重点关注：月将（太阳躔次）与占时构成的格局；年命与日辰的旺衰互动\n" +
			"- 应期判断：斗建发时当月见、气首难过半月间、五日远者得时不出八刻中（《大全》卷三应期歌）\n" +
			"- 关键格局：**三阳课 / 三阴课 / 三光课 / 昏暗课 / 日辰俱旺 / 日辰俱衰**\n" +
			"- 整体判断：本命行年得德合、入吉将三传 → 顺时；遇岁月破、丧吊、白虎入命 → 逆时\n"
	default:
		return "【六壬断课侧重 · 通用】\n" +
			"- 以日干为我、日支为事、四课为表、三传为里、年命为最终救应权柄\n" +
			"- 三传所乘天将定吉凶之色、所临地盘定旺衰之地\n" +
			"- 关键信号：旬空（虚浮）、墓神（拘困）、刑冲（动摇）、合德（救应）\n" +
			"- 注意《毕法赋》100 条与卷二天将临 12 宫断辞，逐条对照本盘\n"
	}
}

// FormatQuestionTypeSectionLiuRen 生成六壬版问题类型段（同 FormatQuestionTypeSection 但用六壬侧重）
func FormatQuestionTypeSectionLiuRen(q QuestionType) string {
	if q == "" || q == QTOther {
		return "## 解课侧重\n" + FocusGuideLiuRen(QTOther) + "\n"
	}
	var b strings.Builder
	b.WriteString("## 问题类型\n")
	b.WriteString(fmt.Sprintf("%s\n\n", QuestionTypeLabel(q)))
	b.WriteString("## 解课侧重\n")
	b.WriteString(FocusGuideLiuRen(q))
	b.WriteString("\n")
	return b.String()
}
