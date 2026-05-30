package liuren

// DunGan 遁干：给定地支在当前旬中所配之天干
//
// 例：甲子旬，子→甲、丑→乙、寅→丙、卯→丁、辰→戊、巳→己、午→庚、未→辛、申→壬、酉→癸；
// 戌、亥为旬空，无遁干（返回 -1）。
//
// 参数 jiaziIndex 为日干支的六十甲子序号 0..59。
func DunGan(z Zhi, jiaziIndex int) Gan {
	xunStart := (jiaziIndex / 10) * 10 // 旬首六十甲子序
	// 旬首地支：甲子旬首支=子；甲戌旬首支=戌；每旬"旬首地支"也递推
	xunStartZhi := Zhi(xunStart % 12) // 0→子, 10→戌, 20→申, 30→午, 40→辰, 50→寅
	// offset = (z - xunStartZhi + 12) % 12 给出该支在本旬中的"第几位"(0..11)
	off := (int(z) - int(xunStartZhi) + 12) % 12
	if off >= 10 {
		// 落旬空，无遁干
		return -1
	}
	return Gan(off) // 甲=0, 乙=1, ...
}

// IsDingShen 是否为六丁神（遁干为丁的地支）
//
// 毕法第 25、26、69 条皆以六丁神为断断语触发要素。
func IsDingShen(z Zhi, jiaziIndex int) bool {
	return DunGan(z, jiaziIndex) == Ding
}
