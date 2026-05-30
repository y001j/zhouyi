package liuren

// ShenSha 神煞定义与求取
//
// 六壬常用神煞按"起例基"分两类：
//   - 以日支三合局起例：驿马、桃花（咸池）、劫煞、华盖、将星、亡神
//   - 以月建（月支）起例：天马、天医、天喜、月德、月将（另行处理）
//   - 以日干起例：天乙贵人（已在贵人部分处理）
//
// MVP 取最常用的 9 种。

type ShenShaEntry struct {
	Name string // 神煞名
	Zhi  Zhi    // 落在哪个地支
	Desc string // 吉凶与象意简注
}

// ComputeShenSha 生成本课盘面相关神煞（落在哪个地支上）
func ComputeShenSha(ctx *Context) []ShenShaEntry {
	monthZhi := monthZhiOf(ctx)
	dayZhi := ctx.DayZhi

	list := []ShenShaEntry{
		{Name: "驿马", Zhi: yiMa(dayZhi), Desc: "主奔走动变、迁徙出行"},
		{Name: "天马", Zhi: tianMa(monthZhi), Desc: "主贵动、远行之喜"},
		{Name: "桃花", Zhi: taoHua(dayZhi), Desc: "又名咸池，主情色、异性、酒乐"},
		{Name: "劫煞", Zhi: jieSha(dayZhi), Desc: "主劫夺破财、突发灾祸"},
		{Name: "华盖", Zhi: huaGai(dayZhi), Desc: "主孤独清高、艺术宗教"},
		{Name: "将星", Zhi: jiangXing(dayZhi), Desc: "主权势、带兵掌事"},
		{Name: "亡神", Zhi: wangShen(dayZhi), Desc: "主耗散失脱、虚惊"},
		{Name: "天医", Zhi: tianYi(monthZhi), Desc: "主医药康复、病解"},
		{Name: "天喜", Zhi: tianXi(monthZhi), Desc: "主喜庆婚姻"},
		{Name: "天德", Zhi: tianDe(monthZhi), Desc: "解凶救应之吉神，主庇护"},
		{Name: "月德", Zhi: yueDe(monthZhi), Desc: "解凶救应之吉神，与天德同功"},
		{Name: "解神", Zhi: jieShen(monthZhi), Desc: "主纷扰之解、官非之释"},
		{Name: "三奇", Zhi: sanQi(ctx.Gan), Desc: "天上三奇所聚，主大吉、贵人扶持"},
		{Name: "岁破", Zhi: suiPo(ctx), Desc: "太岁所冲之位，主破败、失意"},
		{Name: "月破", Zhi: yuePo(monthZhi), Desc: "月支所冲之位，主当月破败、克应在身"},
	}
	return list
}

// tianDe 天德：以月建（节气月）起例
//
//	歌诀：正月丁、二月坤、三月壬、四月辛、五月乾、六月甲、
//	      七月癸、八月艮、九月丙、十月乙、冬月巽、腊月庚（地支化简见下）
//
// 大六壬只取十二支位（不论干、不论坤艮乾巽四隅卦），故对应化简为：
//
//	正→丁(简化为午)、二→申、三→壬(亥)、四→辛(酉)、五→亥、六→甲(寅)、
//	七→癸(子)、八→寅、九→丙(巳)、十→乙(辰)、十一→巳、十二→庚(申)
//
// 注：派系略有不同，此处采《大全》卷一神煞章主流取法。
var tianDeByMonthZhi = map[Zhi]Zhi{
	Yin: Wu, Mao: Shen, Chen: Hai, Si: You, Wu: Hai, Wei: Yin,
	Shen: Zi, You: Yin, Xu: Si, Hai: Chen, Zi: Si, Chou: Shen,
}

func tianDe(monthZhi Zhi) Zhi {
	if v, ok := tianDeByMonthZhi[monthZhi]; ok {
		return v
	}
	return -1
}

// yueDe 月德：寅午戌月在丙(巳)，申子辰月在壬(亥)，巳酉丑月在庚(申)，亥卯未月在甲(寅)
func yueDe(monthZhi Zhi) Zhi {
	switch sanHeGroup(monthZhi) {
	case 0:
		return Hai // 申子辰 → 壬寄亥
	case 1:
		return Yin // 亥卯未 → 甲寄寅
	case 2:
		return Si // 寅午戌 → 丙寄巳
	case 3:
		return Shen // 巳酉丑 → 庚寄申
	}
	return -1
}

// jieShen 解神：申月起戌，逆行十二位
//
//	实用化口诀：正月在申、二月酉、三月戌……（即随月顺行 +6 之类多家不同）
//
// 此处采《大全》主流："正二在申、三四在戌、五六在子、七八在寅、九十在辰、十一二在午"
var jieShenByMonthZhi = map[Zhi]Zhi{
	Yin: Shen, Mao: Shen, Chen: Xu, Si: Xu,
	Wu: Zi, Wei: Zi, Shen: Yin, You: Yin,
	Xu: Chen, Hai: Chen, Zi: Wu, Chou: Wu,
}

func jieShen(monthZhi Zhi) Zhi {
	if v, ok := jieShenByMonthZhi[monthZhi]; ok {
		return v
	}
	return -1
}

// sanQi 三奇：日干推
//
//	甲日：戊庚之奇 → 取庚寄宫申
//	乙日：丙丁之奇 → 取丁寄宫未
//	丙日：丁戊 → 取戊寄宫巳
//	丁日：丙乙 → 取乙寄宫辰
//	戊日：庚甲 → 取庚寄宫申
//	己日：戊庚 → 取庚寄宫申
//	庚日：壬癸 → 取壬寄宫亥
//	辛日：壬癸 → 取壬寄宫亥
//	壬日：癸甲 → 取癸寄宫丑
//	癸日：壬辛 → 取辛寄宫戌
//
// 注：三奇本是干层的概念，落地支取其奇干寄宫位作为代表。
func sanQi(g Gan) Zhi {
	switch g {
	case Jia, Wu1, Ji:
		return Shen // 庚寄
	case Yi, Ding:
		return Wei // 丁寄
	case Bing:
		return Si // 戊寄
	case Geng, Xin:
		return Hai // 壬寄
	case Ren:
		return Chou // 癸寄
	case Gui:
		return Xu // 辛寄
	}
	return -1
}

// suiPo 岁破：太岁冲位
func suiPo(ctx *Context) Zhi {
	if ctx.Lunar == nil {
		return -1
	}
	yz := ParseZhi(ctx.Lunar.GetYearZhiExact())
	if yz < 0 {
		return -1
	}
	return Zhi((int(yz) + 6) % 12)
}

// yuePo 月破：月支冲位
func yuePo(monthZhi Zhi) Zhi {
	if monthZhi < 0 {
		return -1
	}
	return Zhi((int(monthZhi) + 6) % 12)
}

// monthZhiOf 取干支月的地支（以节气分月）
func monthZhiOf(ctx *Context) Zhi {
	if ctx.Lunar == nil {
		return ctx.DayZhi // 兜底
	}
	mz := ctx.Lunar.GetMonthZhiExact()
	return ParseZhi(mz)
}

// 三合局索引：申子辰(水)=0, 亥卯未(木)=1, 寅午戌(火)=2, 巳酉丑(金)=3
func sanHeGroup(z Zhi) int {
	switch z {
	case Shen, Zi, Chen:
		return 0
	case Hai, Mao, Wei:
		return 1
	case Yin, Wu, Xu:
		return 2
	case Si, You, Chou:
		return 3
	}
	return -1
}

// tianMa 天马：正月午、二月申、三月戌、四月子、五月寅、六月辰，七月后循环
//
//	歌诀：正七在午申，二八在申戌……（以节气月起）
var tianMaByMonthZhi = map[Zhi]Zhi{
	Yin: Wu, Mao: Shen, Chen: Xu, Si: Zi, Wu: Yin, Wei: Chen,
	Shen: Wu, You: Shen, Xu: Xu, Hai: Zi, Zi: Yin, Chou: Chen,
}

func tianMa(monthZhi Zhi) Zhi {
	if v, ok := tianMaByMonthZhi[monthZhi]; ok {
		return v
	}
	return -1
}

// taoHua 桃花（咸池）：三合局首位之前一位（即三合"沐浴位"）
//
//	申子辰在酉、寅午戌在卯、巳酉丑在午、亥卯未在子
func taoHua(z Zhi) Zhi {
	switch sanHeGroup(z) {
	case 0:
		return You
	case 1:
		return Zi
	case 2:
		return Mao
	case 3:
		return Wu
	}
	return -1
}

// jieSha 劫煞：三合局首位对冲之前一位（冲位+1），即三合绝位
//
//	申子辰在巳、亥卯未在申、寅午戌在亥、巳酉丑在寅
func jieSha(z Zhi) Zhi {
	switch sanHeGroup(z) {
	case 0:
		return Si
	case 1:
		return Shen
	case 2:
		return Hai
	case 3:
		return Yin
	}
	return -1
}

// huaGai 华盖：三合局末位（墓库）
//
//	申子辰在辰、亥卯未在未、寅午戌在戌、巳酉丑在丑
func huaGai(z Zhi) Zhi {
	switch sanHeGroup(z) {
	case 0:
		return Chen
	case 1:
		return Wei
	case 2:
		return Xu
	case 3:
		return Chou
	}
	return -1
}

// jiangXing 将星：三合局中位（帝旺）
//
//	申子辰在子、亥卯未在卯、寅午戌在午、巳酉丑在酉
func jiangXing(z Zhi) Zhi {
	switch sanHeGroup(z) {
	case 0:
		return Zi
	case 1:
		return Mao
	case 2:
		return Wu
	case 3:
		return You
	}
	return -1
}

// wangShen 亡神：三合局首位之后一位（即临官禄位）
//
//	申子辰在亥、亥卯未在寅、寅午戌在巳、巳酉丑在申
func wangShen(z Zhi) Zhi {
	switch sanHeGroup(z) {
	case 0:
		return Hai
	case 1:
		return Yin
	case 2:
		return Si
	case 3:
		return Shen
	}
	return -1
}

// tianYi 天医：月建前一位
//
//	正月丑、二月寅……十二月子
var tianYiByMonthZhi = map[Zhi]Zhi{
	Yin: Chou, Mao: Yin, Chen: Mao, Si: Chen, Wu: Si, Wei: Wu,
	Shen: Wei, You: Shen, Xu: You, Hai: Xu, Zi: Hai, Chou: Zi,
}

func tianYi(monthZhi Zhi) Zhi {
	if v, ok := tianYiByMonthZhi[monthZhi]; ok {
		return v
	}
	return -1
}

// tianXi 天喜：正月戌、二月亥……（月建+8）
func tianXi(monthZhi Zhi) Zhi {
	if monthZhi < 0 {
		return -1
	}
	return Zhi((int(monthZhi) + 8) % 12)
}
