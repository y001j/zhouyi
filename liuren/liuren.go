package liuren

import "time"

// Divine 完整起课：给定时刻生成一副盘面。
func Divine(t time.Time) (*Pan, error) {
	ctx, err := BuildContext(t)
	if err != nil {
		return nil, err
	}
	return DivineWithContext(ctx), nil
}

// DivineWithContext 使用给定上下文起课（便于测试）
func DivineWithContext(ctx *Context) *Pan {
	tianpan := BuildTianPan(ctx.YueJiang, ctx.ZhanShi)
	dipan := BuildDiPan()
	ke := BuildSiKe(ctx, tianpan)
	san := DeriveSanChuan(ctx, tianpan, ke)
	tj := PlaceTianJiang(ctx, tianpan)

	// 为三传补齐：天将、六亲、空亡
	fill := func(ce *ChuanEntry) {
		ce.TianJiang = TianJiangOf(tianpan, tj, ce.Zhi)
		ce.LiuQin = LiuQinOfZhiByGan(ce.Zhi, ctx.Gan)
		ce.IsKong = IsXunKong(ce.Zhi, ctx.JiaziIndex)
	}
	fill(&san.Chu)
	fill(&san.Zhong)
	fill(&san.Mo)

	pan := &Pan{
		Ctx: ctx, TianPan: tianpan, DiPan: dipan,
		SiKe: ke, SanChuan: san, TianJiang: tj,
	}
	pan.KeTi = ResolveKeTi(pan)
	pan.Tags = KeTiTags(pan)
	pan.ShenSha = ComputeShenSha(ctx)

	// 年命/行年先于课格识别：因 KeTiGejuMore "繁昌课" / MatchBiFa "帘幕贵人/催官使者" 等规则
	// 都依赖 pan.NianMing。
	if ctx.BenMing != nil || ctx.BirthYear > 0 {
		info := &BenMingInfo{}
		if ctx.BenMing != nil {
			nm := BuildNianMing(pan, *ctx.BenMing)
			info.BenMing = &nm
		}
		if ctx.BirthYear > 0 && ctx.Lunar != nil {
			// 虚岁 = 当前年 - 出生年 + 1（简化：按公历年；精确可加立春判断）
			va := ctx.Time.Year() - ctx.BirthYear + 1
			if va < 1 {
				va = 1
			}
			xn := ComputeXingNian(ctx.Gender, va)
			nm := BuildNianMing(pan, xn)
			info.XingNian = &nm
		}
		// 本命/行年互动（《大全》卷三 心印赋）
		if info.BenMing != nil && info.XingNian != nil {
			info.BMXNRel, info.BMXNDesc = ClassifyBenMingXingNian(
				info.BenMing.Zhi, info.XingNian.Zhi)
		}
		pan.NianMing = info
	}

	// 课格识别（含基础格 + 8 五行成局 + 6 高频卷六课格 + 5 次频课格）
	pan.Tags = append(pan.Tags, KeTiGeju(pan)...)
	pan.Tags = append(pan.Tags, KeTiGejuMore(pan)...)
	pan.Tags = append(pan.Tags, KeTiGejuExtra(pan)...)
	pan.Tags = append(pan.Tags, KeTiGejuV4(pan)...)
	pan.BiFa = MatchBiFa(pan)
	pan.Taboos = CheckTianJiangTaboos(tj)

	return pan
}
