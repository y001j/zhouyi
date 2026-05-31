package qimen

import "testing"

// TestLocateZhiShi_Golden 值使门 golden test（转盘法／拆补法）。
//
// 正确规则（《烟波钓叟歌》"值符常遣加时干，值使顺逆遁宫去"）：
//  1. 值使门名 = 旬首遁干在地盘所落宫的固定配门；遁干寄中5者寄坤2、配死门。
//  2. 值使落宫 = 自「值使本宫」（旬首遁干落宫，寄中5者寄坤2）起，沿转盘链
//     [坎1,艮8,震3,巽4,离9,坤2,兑7,乾6] 前进「时辰在本旬内的序数 n」步
//     （旬首 n=0 即伏吟），阳遁顺行、阴遁逆行。
//
// 关键自洽性：值使与值符在转盘法中连动同转，故旬首时辰（n=0）值使必与值符同宫、
// 全盘伏吟。本批期望值由独立脚本按上法生成，覆盖各旬首、阴阳遁、寄中5、整旬步进。
//
// 注：旧实现误用「时支地盘宫」定落宫，仅在旬首伏吟点偶合，已修正（见 layout.go）。
func TestLocateZhiShi_Golden(t *testing.T) {
	type tc struct {
		name     string
		dun      string
		ju       int
		dungan   string // 旬首六仪
		xunshou  string // 旬首（甲X）
		hourGZ   string // 时柱干支
		wantGate string
		wantPal  int // 飞星索引 0..8
	}
	tests := []tc{
		// ===== 阳遁一局 · 甲子旬（旬首戊落坎1 → 休门）：整旬步进，n=0 伏吟坎1 =====
		{"阳一甲子时(伏吟)", "阳遁", 1, "戊", "甲子", "甲子", "休门", 0},
		{"阳一乙丑时(n=1)", "阳遁", 1, "戊", "甲子", "乙丑", "休门", 7}, // 坎1顺1→艮8
		{"阳一丙寅时(n=2)", "阳遁", 1, "戊", "甲子", "丙寅", "休门", 2}, // →震3
		{"阳一丁卯时(n=3)", "阳遁", 1, "戊", "甲子", "丁卯", "休门", 3}, // →巽4
		{"阳一戊辰时(n=4)", "阳遁", 1, "戊", "甲子", "戊辰", "休门", 8}, // →离9
		{"阳一癸酉时(n=9)", "阳遁", 1, "戊", "甲子", "癸酉", "休门", 7}, // 绕回艮8

		// ===== 阳遁二局 · 甲子旬（戊落坤2 → 死门）=====
		{"阳二甲子时(伏吟坤2)", "阳遁", 2, "戊", "甲子", "甲子", "死门", 1},
		{"阳二丁卯时(n=3)", "阳遁", 2, "戊", "甲子", "丁卯", "死门", 0}, // 坤2顺3→坎1

		// ===== 阴遁九局 · 甲子旬（戊落离9 → 景门）：逆行 =====
		{"阴九甲子时(伏吟离9)", "阴遁", 9, "戊", "甲子", "甲子", "景门", 8},
		{"阴九丙寅时(n=2逆)", "阴遁", 9, "戊", "甲子", "丙寅", "景门", 2}, // 离9逆2→震3
		{"阴九庚午时(n=6逆)", "阴遁", 9, "戊", "甲子", "庚午", "景门", 6}, // →兑7

		// ===== 阳遁五局 · 甲子旬（戊寄中5 → 寄坤2 → 死门，伏吟坤2）=====
		{"阳五甲子时(戊寄中5→坤2)", "阳遁", 5, "戊", "甲子", "甲子", "死门", 1},

		// ===== 各旬首 · 旬首时辰均应伏吟（落各自值使本宫）=====
		{"阳一甲戌时(己落坤2)", "阳遁", 1, "己", "甲戌", "甲戌", "死门", 1},
		{"阳一甲申时(庚落震3)", "阳遁", 1, "庚", "甲申", "甲申", "伤门", 2},
		{"阳一甲午时(辛落巽4)", "阳遁", 1, "辛", "甲午", "甲午", "杜门", 3},
		{"阳一甲辰时(壬寄中5→坤2)", "阳遁", 1, "壬", "甲辰", "甲辰", "死门", 1},
		{"阳一甲寅时(癸落乾6)", "阳遁", 1, "癸", "甲寅", "甲寅", "开门", 5},

		// ===== 真实复验盘：阳遁8局 甲申旬（旬首庚落坎1 → 休门）=====
		// 2026-05-31 15:30 即此盘的甲申时，值使应伏吟落坎1（与值符天蓬同宫）。
		{"阳八甲申时(伏吟坎1)", "阳遁", 8, "庚", "甲申", "甲申", "休门", 0},
		{"阳八乙酉时(n=1)", "阳遁", 8, "庚", "甲申", "乙酉", "休门", 7}, // →艮8
		{"阳八丙戌时(n=2)", "阳遁", 8, "庚", "甲申", "丙戌", "休门", 2}, // →震3
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			earth := LayEarthStems(tc.dun, tc.ju)
			gotGate, gotPal := LocateZhiShi(earth, tc.dungan, tc.xunshou, tc.hourGZ, tc.dun)
			if gotGate != tc.wantGate || gotPal != tc.wantPal {
				t.Errorf("LocateZhiShi(dun=%s, ju=%d, dungan=%s, xunshou=%s, hourGZ=%s):\n  got  (%s, %d)\n  want (%s, %d)",
					tc.dun, tc.ju, tc.dungan, tc.xunshou, tc.hourGZ, gotGate, gotPal, tc.wantGate, tc.wantPal)
			}
		})
	}
}

// TestLayDoors_Smoke 验证八门布盘合理性：值使门落 zhiShiPalaceFei、8 门齐全无重、中5空。
func TestLayDoors_Smoke(t *testing.T) {
	earth := LayEarthStems("阳遁", 1)
	gate, palFei := LocateZhiShi(earth, "戊", "甲子", "甲子", "阳遁") // 甲子时伏吟，休门加坎1
	doors := LayDoors(gate, palFei, "阳遁")

	if doors[palFei] != gate {
		t.Errorf("值使门落位错误：palFei=%d, doors[palFei]=%q, gate=%q", palFei, doors[palFei], gate)
	}
	if doors[4] != "" {
		t.Errorf("中5宫应为空: got %q", doors[4])
	}
	seen := map[string]bool{}
	for i := 0; i < 9; i++ {
		if i == 4 {
			continue
		}
		if doors[i] == "" {
			t.Errorf("飞星索引 %d 宫缺门", i)
			continue
		}
		if seen[doors[i]] {
			t.Errorf("门重复: %q", doors[i])
		}
		seen[doors[i]] = true
	}
	if len(seen) != 8 {
		t.Errorf("门数 %d ≠ 8", len(seen))
	}
}
