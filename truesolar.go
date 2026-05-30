package main

import (
	"math"
	"time"
)

// 默认经度：北京（东经 116.4°）。前端未提供时使用。
const defaultLongitude = 116.4

// equationOfTime 计算给定日期的均时差（EOT），单位分钟。
//
// 真太阳与平太阳的钟表时间差，源于：
//  1. 地球轨道椭圆（开普勒第二定律），近日点附近真太阳走得快
//  2. 黄赤交角 23.5°，太阳在赤道方向的投影分量随季节变化
//
// 全年在 −14（2 月中旬）~ +16 分钟（11 月初）之间波动。
//
// 采用 NOAA 推荐的近似公式（Spencer 简化形式），全年误差 < 30 秒：
//
//	B = 2π × (dayOfYear − 81) / 365
//	EOT = 9.87·sin(2B) − 7.53·cos(B) − 1.5·sin(B)
//
// 入参用 t 的 YearDay()——同一日期 EOT 当作常数，不随时分秒变化（精度足够）。
func equationOfTime(t time.Time) float64 {
	d := float64(t.YearDay() - 81)
	B := 2 * math.Pi * d / 365.0
	return 9.87*math.Sin(2*B) - 7.53*math.Cos(B) - 1.5*math.Sin(B)
}

// ApplyTrueSolarTime 把钟表时 t 校正到本地"真太阳时"。
//
// 两步校正：
//  1. 经度差：t 所在时区的基准经度（offset×15°）与本地经度的差，每度 4 分钟
//  2. 均时差（EOT）：地球轨道与黄赤交角导致的全年 ±16 分钟波动
//
// 公式：t' = t + (本地经度 − 时区基准经度) × 4 分钟 + EOT(日期)
//
// 时区基准经度由 t.Location() 的 UTC 偏移自动推导。
// 例：北京 UTC+8 → 基准 120°E；纽约 UTC−5 → 基准 −75°W。
//
// longitude 单位度，东经为正、西经为负；越界按默认值处理。
func ApplyTrueSolarTime(t time.Time, longitude float64) time.Time {
	if longitude < -180 || longitude > 180 {
		longitude = defaultLongitude
	}
	_, offsetSec := t.Zone()
	tzBaseLon := float64(offsetSec) / 3600.0 * 15.0
	deltaMinutes := (longitude-tzBaseLon)*4.0 + equationOfTime(t)
	return t.Add(time.Duration(deltaMinutes * float64(time.Minute)))
}

// ResolveLongitude 从请求字段（可能为 0 表示未提供）和默认值之间挑选一个有效的经度。
// 0 视为"未提供"——前端真要表达赤道经度 0° 也极少用到，按默认值处理可接受。
func ResolveLongitude(reqLon float64) float64 {
	if reqLon == 0 || reqLon < -180 || reqLon > 180 {
		return defaultLongitude
	}
	return reqLon
}
