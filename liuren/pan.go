package liuren

// Pan 大六壬盘面：起课结束后应具备的全部要素。
type Pan struct {
	Ctx       *Context
	TianPan   [12]Zhi        // 索引 i 表示地盘第 i 位上方的天盘地支
	DiPan     [12]Zhi        // 恒为 {子,丑,...,亥}
	SiKe      [4]Ke          // 四课（从第一课到第四课）
	SanChuan  SanChuan       // 初传、中传、末传
	TianJiang [12]TianJiang  // 地盘第 i 位的天将
	KeTi      KeTi           // 课体（发传法所定之主格）
	Tags      []KeTi         // 附加课格标签（三光/三阳/递生等，0..N 个）
	ShenSha   []ShenShaEntry // 神煞落位
	NianMing  *BenMingInfo   // 年命/行年（可选）
	BiFa      []BiFaEntry    // 毕法赋匹配条目
	Taboos    []TianJiangTaboo // 天将乘临禁忌（《大全》卷一 p498）
}

// Ke 单课（上神+下神）
type Ke struct {
	Index    int    // 1..4
	Upper    Zhi    // 天神
	Lower    Zhi    // 地神
	Relation string // 上克下 / 下贼上 / 上生下 / 下生上 / 比和
}

// ChuanEntry 三传中的一传
type ChuanEntry struct {
	Name      string    // 初传 / 中传 / 末传
	Zhi       Zhi       // 天神
	TianJiang TianJiang // 所乘天将
	LiuQin    LiuQin    // 对日干的六亲
	IsKong    bool      // 是否旬空
}

// SanChuan 三传
type SanChuan struct {
	Chu    ChuanEntry
	Zhong  ChuanEntry
	Mo     ChuanEntry
	Method string // 发传法：贼克/比用/涉害/遥克/昴星/别责/八专/伏吟/返吟
	Note   string // 补充说明（元首/重审等）
}

// KeTi 课体名与简释
type KeTi struct {
	Name    string
	Summary string
}

// BuildDiPan 返回固定地盘（子丑寅…亥）
func BuildDiPan() [12]Zhi {
	var d [12]Zhi
	for i := range d {
		d[i] = Zhi(i)
	}
	return d
}
