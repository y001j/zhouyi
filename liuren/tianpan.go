package liuren

// BuildTianPan 月将加时生成天盘。
//
// 算法：设月将为 M、占时为 S，要求天盘上的 M 压在地盘的 S 位置上方。
// 因此地盘第 i 位之上的天盘支为 (i + (M - S)) mod 12。
func BuildTianPan(yueJiang, zhanShi Zhi) [12]Zhi {
	offset := (int(yueJiang) - int(zhanShi) + 12) % 12
	var tian [12]Zhi
	for i := 0; i < 12; i++ {
		tian[i] = Zhi((i + offset) % 12)
	}
	return tian
}

// UpperOf 返回地盘 z 位置上方的天盘地支
func UpperOf(tianpan [12]Zhi, z Zhi) Zhi {
	return tianpan[z]
}
