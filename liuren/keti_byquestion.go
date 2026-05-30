package liuren

// 课体按问题类型分维度断辞（《六壬大全》卷五课经一）
//
// 卷五每一种九宗门课都给出"占官 / 占婚 / 占财 / 占病 / 占讼 / 占行人 / 占盗"等
// 5-7 个分类断辞。本表为这种结构化形态的一个轻量版：
//   - 维度键采用主包 QuestionType（career/wealth/relation/health/decision/timing）
//   - 仅录入 6 个最常见课体（元首/重审/比用/涉害/伏吟/返吟），其余课体按需后续补
//
// 提示词层会先取按问题类型的断辞；若该课体或该问题类型未录入，回退到
// keti.go 的 ketiSummaryTable 通用 summary。

// QTypeKey 与主包 questiontype.go 的 QuestionType 字符串保持一致
type QTypeKey string

const (
	QTKCareer   QTypeKey = "career"
	QTKWealth   QTypeKey = "wealth"
	QTKRelation QTypeKey = "relation"
	QTKHealth   QTypeKey = "health"
	QTKDecision QTypeKey = "decision"
	QTKTiming   QTypeKey = "timing"
)

// KeTiByQuestion[课体名][问题类型] = 断辞
//
// 来源：《大全》卷五课经一各课首段断辞 + 卷六占类辑要。
// 每条 30-60 字，可与 ketiSummaryTable 互相呼应。
var KeTiByQuestion = map[string]map[QTypeKey]string{
	"元首课": {
		QTKCareer:   "元首正大，**占官**遇唐虞之君，新职升迁皆可，求名得贵；唯须看官鬼是否坐空。",
		QTKWealth:   "元首一上克下，**占财**用客胜主，市贾出色、求名利皆超鸢；财虽得，宜早不宜迟。",
		QTKRelation: "元首主**占婚**男出色、女嫁荣，名分相当；唯不利下犯上者（夫妇争先）。",
		QTKHealth:   "元首**占病**主病轻易愈，由外邪而非内伤；惟凶将乘鬼者另当别论。",
		QTKDecision: "元首主君臣顺序、上下分明；**抉择**当顺势而行，进取无忧。",
		QTKTiming:   "元首**时运**当前，名利双收之兆；唯防年命受克则吉中藏险。",
	},
	"重审课": {
		QTKCareer:   "重审一下贼上，**占官**事起于下/由女人或小人凌侵；占讼则后告者胜。",
		QTKWealth:   "重审**占财**主先难后得；下犯上之课，宜由暗中斡旋。",
		QTKRelation: "重审**占婚**女嫁妇之兆，事多由女方主动；伴朱雀勾陈则有诉讼。",
		QTKHealth:   "重审**占病**主病由内起、女人或子孙而生；忌见兄弟相侵。",
		QTKDecision: "重审下贼上 → **抉择**当审慎，先内后外、先退后进；勿先发难。",
		QTKTiming:   "重审**时运**先抑后扬，初阶受压、终得伸展；宜守不宜攻。",
	},
	"知一课": {
		QTKCareer:   "知一（比用）多克并见，**占官**事有两端，宜择善而从；与日干阴阳同者优先。",
		QTKWealth:   "知一**占财**多途求财，宜择一而行、勿贪多；同气者最易得。",
		QTKRelation: "知一**占婚**有两可之象，宜从亲近近隔者；异侧则疏远。",
		QTKHealth:   "知一**占病**多端杂症，须先分主次；近者易治、远者难愈。",
		QTKDecision: "知一象**抉择**两端，比近者亲、比远者疏；选与本气合者。",
		QTKTiming:   "知一**时运**两端俱动，须从所近：临身者先验，远位者后验。",
	},
	"涉害课": {
		QTKCareer:   "涉害**占官**事多波折、受克最深处即症结；宜深思审断、不可轻动。",
		QTKWealth:   "涉害**占财**财在重重阻碍中得，须经历曲折；孟深仲浅季当休。",
		QTKRelation: "涉害**占婚**婚事多家人意见相左、冷暖间隔；末传得地方可成。",
		QTKHealth:   "涉害**占病**病势缠绵、深入脏腑；须看末传是否得救应。",
		QTKDecision: "涉害**抉择**当避深就浅、舍重取轻；孟深则不宜行、季休则可缓图。",
		QTKTiming:   "涉害**时运**举步维艰、四面受克；惟旺相得地者尚可行。",
	},
	"伏吟课": {
		QTKCareer:   "伏吟天盘各居本位，**占官**事机不动、宜静守；动则愈乱，行年入用方有变化。",
		QTKWealth:   "伏吟**占财**财藏不露、宜静候时机；用神逢冲方动。",
		QTKRelation: "伏吟**占婚**婚意已定但不发声；宜待月日相合之期再动。",
		QTKHealth:   "伏吟**占病**旧疾复发、内伤暗藏；需冲动方显症状。",
		QTKDecision: "伏吟主静、不主动；**抉择**当守不当攻，待时而动方有出口。",
		QTKTiming:   "伏吟**时运**伏藏不显，外象平静而内有积压；宜韬光养晦。",
	},
	"返吟课": {
		QTKCareer:   "返吟天地反覆，**占官**主调动/反覆；事虽变动但主客明，宜应变。",
		QTKWealth:   "返吟**占财**财来财去、得而复失；交易反覆、宜速决勿留滞。",
		QTKRelation: "返吟**占婚**婚事反覆、合而又离；忌见冲破之神。",
		QTKHealth:   "返吟**占病**病势反覆、时好时坏；旧疾复发或新疾交侵。",
		QTKDecision: "返吟主反复，**抉择**所定之事易翻盘；宜留余地、勿一击决断。",
		QTKTiming:   "返吟**时运**变动不居、阴晴不定；宜随机应变。",
	},
	"井栏课": {
		QTKCareer:   "井栏返吟无克之特例，**占官**事如井栏临渊，宜谨慎；丁己辛日逢丑未占之。",
		QTKWealth:   "井栏**占财**财在井中、可见而难取；忌妄进、宜稳守。",
		QTKRelation: "井栏**占婚**婚意虽有但循环未决；宜静待。",
		QTKHealth:   "井栏**占病**病机伏而未发，须细察。",
		QTKDecision: "井栏**抉择**进退维谷、当慎勿坠。",
		QTKTiming:   "井栏**时运**临边不进，宜守业勿求新。",
	},
}

// KeTiSummaryByQuestion 取课体的按问题类型断辞；若未录入则返回空串（调用方应回退到通用 summary）
func KeTiSummaryByQuestion(ketiName, questionType string) string {
	m, ok := KeTiByQuestion[ketiName]
	if !ok {
		return ""
	}
	if v, ok := m[QTypeKey(questionType)]; ok {
		return v
	}
	return ""
}
