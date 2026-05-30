package qimen

// 本文件实现奇门遁甲"布盘"的全部核心算法：
//
//   LayEarthStems      布地盘三奇六仪
//   LocateZhiFu        求值符落宫与值符星
//   RotateHeavenStems  天盘干旋转（值符加时干）
//   RotateStars        九星旋转（随天盘）
//   LayGods            布八神（阳顺阴逆，起于值符落宫）
//   LocateZhiShi       求值使门与值使落宫（值使加时支）
//   LayDoors           布八门（起于值使落宫）
//
// 所有输入/输出均以 "飞星序索引 0..8" 表示宫位（下标 0 对应坎1，4 对应中5，8 对应离9）。

// ======================== 1. 地盘三奇六仪 ========================

// LayEarthStems 按 dun 和 ju 布地盘三奇六仪。
//
//   阳遁：戊起于 ju 宫，顺走 1→2→3→4→5→6→7→8→9
//   阴遁：戊起于 ju 宫，逆走 9→8→7→6→5→4→3→2→1
//
// 填入顺序固定为 "戊己庚辛壬癸丁丙乙"（三奇丁丙乙为逆序）。
//
// 返回 9 元素数组：下标是飞星序索引（0..8），值是该宫的地盘干。
func LayEarthStems(dun string, ju int) [9]string {
	var arr [9]string
	start := ju - 1 // 飞星序索引（0..8）
	step := 1
	if dun == "阴遁" {
		step = -1
	}
	for i := 0; i < 9; i++ {
		idx := ((start + i*step) % 9 + 9) % 9
		arr[idx] = SanqiLiuyiSeq[i]
	}
	return arr
}

// ======================== 2. 值符落宫（天盘值符星） ========================

// LocateZhiFu 返回 (值符星, 值符落宫飞星序索引)。
//
// 算法（《奇门宝鉴·直符加时干法》）：
//  1. 找遁干（旬首六仪）在地盘的宫 —— 该宫的九星即为值符星。
//     若遁干落中5（地盘 earthStems[4]=旬首），则值符星 = 天禽，寄坤2。
//  2. 找时干在地盘的宫 —— 即为值符落宫。
//     若时干 = "甲"，甲遁于旬首六仪，取遁干的宫。
//     若时干落中5（earthStems[4]=hourGan），值符寄坤2。
//
// 返回的 palaceFei 始终是 0..8 的有效索引（寄宫时返回坤2的索引1）。
// 额外返回 rawIsMiddle 标记原始算出的宫是否为中5（调用方需要时可用于展示）。
func LocateZhiFu(earthStems [9]string, dungan, hourGan string) (star string, palaceFei int, rawIsMiddle bool) {
	// 1. 遁干（六仪）在地盘的宫
	dunPalace := findStemPalace(earthStems, dungan)
	if dunPalace == 4 {
		star = "天禽"
	} else {
		star = PalaceToStar[dunPalace]
	}

	// 2. 时干在地盘的宫；甲时取遁干
	var hourPalace int
	if hourGan == "甲" {
		hourPalace = dunPalace
	} else {
		hourPalace = findStemPalace(earthStems, hourGan)
	}

	rawIsMiddle = hourPalace == 4
	if rawIsMiddle {
		// 寄坤2（飞星序索引 1）
		palaceFei = 1
	} else {
		palaceFei = hourPalace
	}
	return
}

// findStemPalace 在地盘 earthStems（飞星序 9 元素数组）中找到指定天干的宫（飞星序索引 0..8）。
// 若未找到返回 -1（不应发生，除非输入异常）。
func findStemPalace(earthStems [9]string, stem string) int {
	for i, s := range earthStems {
		if s == stem {
			return i
		}
	}
	return -1
}

// ======================== 3. 天盘干旋转（值符加时干） ========================

// RotateHeavenStems 生成天盘干（9 元素，飞星序，中5 宫置空）。
//
// 方法：把地盘按"转盘链 [1,8,3,4,9,2,7,6]"重排成 8 元素序列；
// 从"遁干在此链上的位置"开始，将序列平移到"值符落宫在此链上的位置"。
// 这实现了"直符常遣加时干，值符带动天盘整体旋转"的几何操作。
//
// 参数 zhiFuPalaceFei 是值符落宫的飞星索引（由 LocateZhiFu 给出，已处理中5寄坤2）。
// 参数 dungan、hourGan 用于确定在转盘链上的起点与终点（与 LocateZhiFu 的内部计算一致）。
func RotateHeavenStems(earthStems [9]string, dungan, hourGan string, zhiFuPalaceFei int) [9]string {
	var heaven [9]string

	// 把地盘按转盘链取出（8 元素）
	earthZhuan := make([]string, 8)
	for i := 0; i < 8; i++ {
		earthZhuan[i] = earthStems[ZhuanToFei[i]-1]
	}

	// 求转盘链上的起点索引：遁干在转盘链上的位置
	fromIdx := findStemPalaceInZhuan(earthStems, dungan)
	// 求终点索引：值符落宫在转盘链上的位置
	toIdx := FeiToZhuan[zhiFuPalaceFei+1] // zhiFuPalaceFei 0..8 → FeiToZhuan 索引 +1

	// 安全兜底：中5 的 FeiToZhuan=-1，强制为坤2 转盘索引
	if toIdx < 0 {
		toIdx = FeiToZhuan[2] // 坤2 对应 FeiToZhuan[2]=7
	}

	shift := (toIdx - fromIdx + 8) % 8

	// 沿转盘链整体平移
	for i := 0; i < 8; i++ {
		srcZhuanIdx := (fromIdx + i) % 8
		dstZhuanIdx := (fromIdx + i + shift) % 8
		// dstZhuanIdx 在飞星序的宫号：ZhuanToFei[dstZhuanIdx]
		destFei := ZhuanToFei[dstZhuanIdx] - 1
		heaven[destFei] = earthZhuan[srcZhuanIdx]
	}
	// 中5 宫：天盘不直接写；保持空串（调用方按需"寄坤2"展示）

	return heaven
}

// findStemPalaceInZhuan 在转盘链上找到指定天干的索引（0..7）。若在中5则返回坤2的转盘索引。
func findStemPalaceInZhuan(earthStems [9]string, stem string) int {
	palace := findStemPalace(earthStems, stem)
	if palace == 4 {
		return FeiToZhuan[2] // 坤2
	}
	return FeiToZhuan[palace+1] // palace 是 0..8 的飞星索引，+1 得宫号 1..9
}

// ======================== 4. 九星旋转（随天盘） ========================

// RotateStars 生成九星分布（9 元素，中5 为天禽，其余按转盘链旋转）。
//
// 和 RotateHeavenStems 同构：把"天蓬、天任、天冲、天辅、天英、天芮、天柱、天心"按
// 转盘链 [坎1,艮8,震3,巽4,离9,坤2,兑7,乾6] 的原位顺序排开，然后整体平移。
func RotateStars(earthStems [9]string, dungan, hourGan string, zhiFuPalaceFei int) [9]string {
	var stars [9]string
	stars[4] = "天禽" // 中5 恒为天禽

	fromIdx := findStemPalaceInZhuan(earthStems, dungan)
	toIdx := FeiToZhuan[zhiFuPalaceFei+1]
	if toIdx < 0 {
		toIdx = FeiToZhuan[2]
	}
	shift := (toIdx - fromIdx + 8) % 8

	// JiuxingSeq 与 ZhuanToFei 同序：转盘索引 i 对应九星 JiuxingSeq[i]
	for i := 0; i < 8; i++ {
		dstZhuanIdx := (i + shift) % 8
		destFei := ZhuanToFei[dstZhuanIdx] - 1
		stars[destFei] = JiuxingSeq[i]
	}
	return stars
}

// ======================== 5. 八神（阳顺阴逆，起于值符落宫） ========================

// LayGods 布八神（9 元素，中5 置空）。
// 起点是值符落宫（zhiFuPalaceFei 已处理中5寄坤2）；阳遁沿转盘链顺排，阴遁逆排。
func LayGods(dun string, zhiFuPalaceFei int) [9]string {
	var gods [9]string
	startZhuanIdx := FeiToZhuan[zhiFuPalaceFei+1]
	if startZhuanIdx < 0 {
		startZhuanIdx = FeiToZhuan[2] // 坤2 fallback
	}
	step := 1
	if dun == "阴遁" {
		step = -1
	}
	for i := 0; i < 8; i++ {
		zhuanIdx := ((startZhuanIdx + i*step) % 8 + 8) % 8
		destFei := ZhuanToFei[zhuanIdx] - 1
		gods[destFei] = BashenSeq[i]
	}
	return gods
}

// ======================== 6. 值使门（值使加时支） ========================

// LocateZhiShi 返回 (值使门, 值使落宫飞星序索引)。
//
// 算法（《奇门宝鉴·直使加时法》）：
//  1. 值使门名 = 遁干（旬首六仪）在地盘的宫所对应的门。
//     若遁干落中5，寄坤2（死门）。
//  2. 值使落宫 = 时支在地盘（固定 12 地支 → 宫）所对应的宫。
func LocateZhiShi(earthStems [9]string, dungan, hourZhi string) (gate string, palaceFei int) {
	dunPalace := findStemPalace(earthStems, dungan)
	if dunPalace == 4 {
		// 寄坤2
		gate = PalaceToGate[1]
	} else {
		gate = PalaceToGate[dunPalace]
	}

	palaceFei = ZhiToPalaceFei[hourZhi]
	return
}

// ======================== 7. 布八门 ========================

// LayDoors 以值使门为起点、按 "休生伤杜景死惊开" 顺序，
// 在 8 个非中5宫（转盘链）上布八门。阳遁顺走、阴遁逆走。
//
// zhiShiGate 是值使门名（休门/生门…）。
// zhiShiPalaceFei 是值使落宫的飞星序索引（0..8，值使落中5 时寄坤2 即 1）。
func LayDoors(zhiShiGate string, zhiShiPalaceFei int, dun string) [9]string {
	var doors [9]string

	startGateIdx := indexOfString(BamenSeq[:], zhiShiGate)

	// 值使落宫在转盘链的索引；若落中5，寄坤2
	startZhuanIdx := FeiToZhuan[zhiShiPalaceFei+1]
	if startZhuanIdx < 0 {
		startZhuanIdx = FeiToZhuan[2]
	}

	step := 1
	if dun == "阴遁" {
		step = -1
	}

	for i := 0; i < 8; i++ {
		gateIdx := (startGateIdx + i) % 8
		zhuanIdx := ((startZhuanIdx + i*step) % 8 + 8) % 8
		destFei := ZhuanToFei[zhuanIdx] - 1
		doors[destFei] = BamenSeq[gateIdx]
	}
	// 中5 无门
	return doors
}

// ======================== 辅助 ========================

func indexOfString(arr []string, s string) int {
	for i, v := range arr {
		if v == s {
			return i
		}
	}
	return -1
}
