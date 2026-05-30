// 现代汉语 UI 字典（i18n-en.js 文件名沿用，内容已改为白话简体）
// 覆盖范围：UI 外壳——导航、按钮、表单字段、状态提示、模态框、页脚等
// 古典内容（卦辞、爻辞、神煞、原典引文等）仍保留古文，由 data-i18n-keep 标记

window.__I18N_EN__ = {
  // ===== 顶部子导航 =====
  'nav.home': '卜筮明心',
  'nav.zhouyi': '周易占卦',
  'nav.liuren': '大六壬起课',
  'nav.qimen': '奇门遁甲',
  'nav.huican': '三式同参',
  'nav.journal': '占卦笔记',

  // ===== 印心小卡（结果页顶部） =====
  'yinxin.title': '印　心',
  'yinxin.youask': '你问：',
  'yinxin.atwhen': '于',
  'yinxin.got': '，得',
  'yinxin.firstthought': '你看到此卦的第一个念头是——',
  'yinxin.thoughtph': '把这一念写下。它常常比后来的解卦更接近你的本心。',
  'yinxin.thoughthint': '（可以留空。封藏后还能继续编辑，会自动保存。）',
  'yinxin.btn.seal': '封藏此卦',
  'yinxin.sealed': '已封藏 ✓',
  'yinxin.justsealed': '已封藏入占卦笔记',
  'yinxin.saved': '已自动保存',
  'yinxin.btn.viewjournal': '去占卦笔记',
  'yinxin.noquestion': '（未明所问）',

  // ===== 占卦笔记页 =====
  'journal.title': '占卦笔记',
  'journal.motto.left': '所占者　所验者',
  'journal.motto.right': '行而验　验而学',
  'journal.header.title': '占　卦　笔　记',
  'journal.subtitle': '占　而　记　之　·　验　而　悟　之',
  'journal.section.intro': '关于占卦笔记',
  'journal.intro.p1': '「<strong>长期记录所占所验，积累为学问；不是一时的好奇</strong>」——首页《心法》如此说。',
  'journal.intro.p2': '凡是占完封藏的卦，都进入这本笔记。事后可以标记是否应验、写复盘；久而久之，自然能体会「易不过是我心的一面镜子」。',
  'journal.intro.hint': '※ 全部资料只存于本机浏览器，不会主动上传。可以随时导出 JSON 自己保存。',
  'journal.section.stats': '概览',
  'journal.section.list': '所占之卦',
  'journal.empty': '占卦笔记还是空的。当你封藏第一卦，它会出现在这里。',
  'journal.full': '占卦笔记已达 200 卦上限，请先在笔记页导出或清理后再封藏。',
  'journal.stat.total': '凡 卦',
  'journal.stat.fulfilled': '已 应',
  'journal.stat.partial': '半 应',
  'journal.stat.unfulfilled': '未 应',
  'journal.stat.pending': '待 验',
  'journal.btn.export': '导出 JSON',
  'journal.btn.clear': '清空笔记',
  'journal.kind.zhouyi': '周易',
  'journal.kind.liuren': '六壬',
  'journal.kind.qimen': '奇门',
  'journal.kind.huican': '互参',
  'journal.qlabel': '所　问',
  'journal.headline': '所　得',
  'journal.firstthought': '第一念',
  'journal.nothought': '（未记第一念）',
  'journal.noquestion': '（未明所问）',
  'journal.vf.fulfilled': '已 应',
  'journal.vf.unfulfilled': '未 应',
  'journal.vf.partial': '半 应',
  'journal.vf.note': '复盘',
  'journal.btn.verify': '标 应 验',
  'journal.btn.reverify': '修改应验',
  'journal.btn.editthought': '改第一念',
  'journal.btn.del': '删 除',
  'journal.confirm.del': '确认删除此卦？此操作不可恢复。',
  'journal.confirm.clear': '确认清空全部占卦笔记？此操作不可恢复。建议先「导出 JSON」备份。',
  'journal.confirm.clear2': '再次确认：将永久删除全部占卦笔记。',
  'journal.export.empty': '占卦笔记为空，没有可导出的内容。',
  'journal.prompt.editthought': '修改第一念（可空）：',
  'journal.verify.title': '标 记 应 验',
  'journal.verify.fulfilled': '已 应',
  'journal.verify.unfulfilled': '未 应',
  'journal.verify.partial': '半 应 / 待 验',
  'journal.verify.notelabel': '事后复盘（可选）',
  'journal.verify.noteph': '记下事情实际的进展、和卦象的对应或偏离……',
  'journal.verify.save': '保 存',
  'journal.footer': '所占者　心也　·　所验者　行也',

  // ===== 首页 =====
  'index.title': '卜筮明心 · 占前必读',
  'index.motto.left': '不诚不占　不疑不卜',
  'index.motto.right': '占卦须慎　不可儿戏',
  'index.subtitle': '占前必读 · 静心凝神 · 至诚则灵 · 敬而用之',
  'index.toc.title': '目录',
  'index.toc.preface': '序言 · 静心',
  'index.toc.cautions': '一 · 占卦六诫',
  'index.toc.ritual': '二 · 占卦流程',
  'index.toc.heart': '三 · 心法宜忌',
  'index.toc.disclaimer': '四 · 免责声明',
  'index.toc.faq': '五 · 常见问答',
  'index.toc.enter': '终 · 进入起占',
  'index.intro.summary': '本系统 · 简介与优势',
  'index.enter.btn': '已读 · 进入起占',
  'index.enter.hint': '点击上方按钮，表示你已阅读并同意上文的免责声明与占卦戒律。',
  'index.preface.title': '静心',
  'index.preface.subtitle': '—— 读完此篇，胜占十卦 ——',
  'index.cautions.title': '占卦六诫',
  'index.ritual.title': '占卦流程',
  'index.heart.title': '心法与禁忌',
  'index.heart.do': '宜',
  'index.heart.dont': '忌',
  'index.heart.final': '—— 结语 ——',
  'index.disclaimer.title': '免责声明',
  'index.faq.title': '常见问答',
  'index.enter.title': '进入起占',
  'index.enter.zhouyi.title': '周易占卦',
  'index.enter.zhouyi.sub': '事理走向 · 大方向',
  'index.enter.liuren.title': '大六壬',
  'index.enter.liuren.sub': '人事细节 · 应期',
  'index.enter.qimen.title': '奇门遁甲',
  'index.enter.qimen.sub': '方位时机 · 谋事',
  'index.enter.huican.title': '三式同参',
  'index.enter.huican.sub': '三盘同参 · 互证',

  // ===== 模态框（首次访问免责） =====
  'modal.title': '免责声明 · 使用须知',
  'modal.intro': '欢迎使用本系统。在继续之前，请仔细阅读以下条款；点击「同意并继续」表示你已知悉并接受这些条款。',
  'modal.t1': '本系统是<strong>传统文化研习与占卜体验</strong>工具，输出的卦象、课象、局象与 AI 解读<strong>仅供参考</strong>，不构成任何专业建议。',
  'modal.t2': '<strong>不替代医疗建议</strong>：涉及健康、疾病、用药、诊疗，请及时咨询执业医师，不要以占卜代替正规医疗判断。',
  'modal.t3': '<strong>不替代法律建议</strong>：涉及诉讼、合同、维权、赔偿，请咨询执业律师。',
  'modal.t4': '<strong>不替代投资建议</strong>：涉及金融、股票、投资、创业，请基于自身研究和风险承受能力做决定。本系统不为任何投资结果负责。',
  'modal.t5': '<strong>不构成命运定论</strong>：占卜重在启发思考、提示风险、辅助决策，最终选择和后果取决于你自己的行动。',
  'modal.t6': '<strong>用户自主决策</strong>：你使用本系统得出的任何结论，应由你自己独立判断并承担相应后果，本系统及开发者对由此产生的任何损失不承担责任。',
  'modal.t7': '<strong>禁止不当用途</strong>：禁止将本系统用于赌博、欺诈、损害他人或违反法律道德的目的；也不应在未获他人授权时代为占测其私事。',
  'modal.t8': '<strong>数据与隐私</strong>：你输入的问题和命主资料（姓名、性别、出生年月日时、出生地、现居地）只存在你本机浏览器里，不会主动上传到第三方；其中只有性别和出生年参与六壬「年命」推算，姓名、出生地、现居地、月日时辰等<strong>不参与卦象算法</strong>，只是起卦时的仪式化记录；如果使用 AI 提示词功能，请自行评估粘贴出去的信息有多敏感。',
  'modal.warn': '如果你不同意上述任一条款，请立即关闭本页面，停止使用本系统。',
  'modal.check': '我已仔细阅读并同意上述免责声明与使用条款',
  'modal.confirm': '同意并继续',

  // ===== 通用表单字段 =====
  'form.question.label': '想问的事（可留空）',
  'form.question.placeholder': '静下来想一想，写下心中所问……',
  'form.qtype.label': '问题类别',
  'form.qtype.other': '其他 / 不指定',
  'form.qtype.career': '事业 / 工作',
  'form.qtype.wealth': '财运 / 投资',
  'form.qtype.relation': '感情 / 婚姻',
  'form.qtype.health': '健康 / 身心',
  'form.qtype.decision': '抉择 / 两难',
  'form.qtype.timing': '时运 / 吉凶',
  'form.unspecified': '— 不指定 —',

  // ===== 周易页 =====
  'zhouyi.title': '周易占卦 · 起卦与解卦',
  'zhouyi.motto.left': '太极生两仪',
  'zhouyi.motto.right': '两仪生四象 四象生八卦',
  'zhouyi.subtitle': '两仪生四象　四象生八卦　八卦定吉凶',
  'zhouyi.intro.summary': '周易占卦 · 方法简介',
  'zhouyi.method.label': '起卦方法',
  'zhouyi.method.coin': '铜钱法',
  'zhouyi.method.coin.sub': '三钱六掷',
  'zhouyi.method.yarrow': '蓍草法',
  'zhouyi.method.yarrow.sub': '大衍揲蓍',
  'zhouyi.method.number': '数字法',
  'zhouyi.method.number.sub': '随机起卦',
  'zhouyi.number.label': '数字起卦参数',
  'zhouyi.number.upper.placeholder': '上卦数',
  'zhouyi.number.lower.placeholder': '下卦数',
  'zhouyi.number.changing.placeholder': '变爻 0-6',
  'zhouyi.number.hint': '任意整数自动取模为 1-8；变爻 0 表示无变爻。',
  'zhouyi.btn.divine': '起卦',
  'zhouyi.btn.divine.sub': '静心凝神，心诚则灵',
  'zhouyi.section.hex': '卦象',
  'zhouyi.section.lines': '六爻详情',
  'zhouyi.section.mainhex': '本卦',
  'zhouyi.section.changehex': '变卦（之卦）',
  'zhouyi.section.derived': '衍生卦象',
  'zhouyi.section.timing': '起卦时令',
  'zhouyi.section.guide': '解卦指引',
  'zhouyi.section.prompt': '解卦提示词',
  'zhouyi.prompt.hint': '把这段提示词贴给博学的 AI，可获详细解读。',
  'zhouyi.btn.copy': '复制提示词',
  'zhouyi.lines.col.position': '爻位',
  'zhouyi.lines.col.value': '值',
  'zhouyi.lines.col.type': '类型',
  'zhouyi.lines.col.sym': '爻形',
  'zhouyi.lines.col.text': '爻辞',
  'zhouyi.preface.tip': '占卦有法，不可轻用',
  'zhouyi.preface.tip.link': '请先阅读「占卜须知」',
  'zhouyi.preface.tip.sub': '无事不占　·　一事不二占　·　小事不占　·　邪事不占　·　已决不占　·　代占须慎',
  'zhouyi.footer': '阴阳流转　生生不息',
  'zhouyi.ritual.skip': '跳过仪式',
  'zhouyi.ritual.meditate': '静心凝神',
  'zhouyi.ritual.quote': '心之所系　筮之所应<br/>不诚不占　不疑不卜',
  'zhouyi.ritual.hint': '默念心中所问……',
  'zhouyi.ritual.shake': '摇卦',
  'zhouyi.ritual.shake.coin': '三钱入掌，六掷成卦……',
  'zhouyi.ritual.shake.yarrow': '五十蓍草，十八变而成卦……',
  'zhouyi.ritual.shake.number': '以数定象，上下相合……',
  'zhouyi.ritual.shake.title.coin': '摇卦',
  'zhouyi.ritual.shake.title.yarrow': '揲蓍',
  'zhouyi.ritual.shake.title.number': '布卦',
  'zhouyi.ritual.reveal': '卦象既成',
  'zhouyi.ritual.reveal.hint': '从初爻到上爻，依次显现……',

  // ===== 大六壬页 =====
  'liuren.title': '大六壬 · 起课与断课',
  'liuren.motto.left': '月将加时　四课三传',
  'liuren.motto.right': '神藏煞没　课体天机',
  'liuren.subtitle': '天地盘　十二天将　断占之法',
  'liuren.intro.summary': '大六壬 · 方法简介',
  'liuren.btn.divine': '起课',
  'liuren.btn.divine.sub': '月将加时',
  'liuren.btn.divining': '起课中…',
  'liuren.benming.label': '本命生肖（仅用于年命神煞，可不填）',
  'liuren.birthyear.label': '出生公历年（可不填）',
  'liuren.gender.label': '性别（可不填）',
  'liuren.gender.male': '男',
  'liuren.gender.female': '女',
  'liuren.querent.prefill': '※ 本命、出生年、性别已从「命主」设定自动带入；点右上角「命主」可随时修改。',
  'querent.tag.unset': '命主：未设',
  'querent.tag.label': '命主：',
  'querent.dialog.title.require': '命主设定（必填）',
  'querent.dialog.title.edit': '命主设定',
  'querent.dialog.intro.require': '占卦之前，请先以正心填录命主信息，作为仪式之凭。所填资料只存于本机浏览器，不会上传。',
  'querent.dialog.intro.edit': '所填资料只存于本机浏览器，不会上传。',
  'querent.dialog.hint': '※ 上述资料只存于本机浏览器，不主动上传。其中性别、出生年用于六壬「年命」推算；姓名、出生地、现居地、月日时辰<strong>不参与卦象算法</strong>，只是起卦时的仪式化标识。',
  'querent.field.name': '姓　名',
  'querent.field.gender': '性　别',
  'querent.field.birthdate': '出生公历',
  'querent.field.hour': '出生时辰',
  'querent.field.birthplace': '出　生　地',
  'querent.field.curplace': '现　居　地',
  'querent.field.overseascity.ph': '请输入国家／城市',
  'querent.opt.province': '省／地区',
  'querent.opt.city': '市／区',
  'querent.opt.empty': '—',
  'querent.opt.placeholder': '— 请选 —',
  'querent.opt.year': '年',
  'querent.opt.month': '月',
  'querent.opt.day': '日',
  'querent.btn.clear': '清除',
  'querent.btn.cancel': '取消',
  'querent.btn.save': '保存',
  'liuren.section.board': '盘面',
  'liuren.section.tdp': '天地盘 · 十二天将',
  'liuren.section.sike': '四课',
  'liuren.section.sanchuan': '三传 · 课体',
  'liuren.section.tags': '附加课格',
  'liuren.section.shensha': '神煞落位',
  'liuren.section.taboo': '乘临禁忌',
  'liuren.section.nianming': '年命',
  'liuren.section.bifa': '毕法赋 · 本盘命中',
  'liuren.section.guide': '断课指引',
  'liuren.section.prompt': '断课提示词',
  'liuren.btn.copy': '复制',
  'liuren.copied': '已复制',
  'liuren.bifa.expand': '▼ 展开《毕法赋》全文 100 条（知识库）',
  'liuren.taboo.empty': '（本盘无天将乘临禁忌）',
  'liuren.nianming.empty': '（未提供本命／行年信息）',
  'liuren.bifa.empty': '（本盘无自动命中条目，参见下方全文）',
  'liuren.fffa': '发传法',
  'liuren.keti': '课体',
  'liuren.footer': '天机不露　占验在心',

  // ===== 奇门遁甲页 =====
  'qimen.title': '奇门遁甲 · 时家起局',
  'qimen.motto.left': '三式之一　兵机之要',
  'qimen.motto.right': '阴阳顺逆　神藏鬼没',
  'qimen.subtitle': '时家起局　九宫八门　三奇六仪',
  'qimen.intro.summary': '奇门遁甲 · 方法简介',
  'qimen.btn.divine': '起局',
  'qimen.btn.divine.sub': '时家奇门',
  'qimen.section.board': '盘面',
  'qimen.section.ninepalaces': '九宫盘',
  'qimen.legend.zhifu': '值符落宫',
  'qimen.legend.zhishi': '值使落宫',
  'qimen.legend.void': '旬空',
  'qimen.legend.yima': '驿马',
  'qimen.section.patterns': '命中格局',
  'qimen.section.prompt': '解局提示词',
  'qimen.btn.copy': '复制',
  'qimen.copied': '已复制',
  'qimen.footer': '天机藏用　神鬼莫测',

  // ===== 三式同参页 =====
  'huican.title': '周易 × 大六壬 × 奇门遁甲 · 三式同参',
  'huican.motto.left': '同时起占　三式同参',
  'huican.motto.right': '卦显象义　壬主人事　奇定机方',
  'huican.subtitle': '一问三占　互为印证　三盘同参',
  'huican.intro.summary': '三式同参 · 方法简介与对照',
  'huican.btn.divine': '起占',
  'huican.btn.divine.sub': '三式同参',
  'huican.col.zhouyi': '周易',
  'huican.col.liuren': '大六壬',
  'huican.col.qimen': '奇门遁甲',
  'huican.section.prompt': '同参提示词',
  'huican.btn.copy': '复制',
  'huican.copied': '已复制',
  'huican.footer': '卦壬互证　吉凶自明',

  // ===== 管理后台 =====
  'admin.title': '管理员控制台',
  'admin.back': '← 返回起占',
  'admin.logout': '登出',
  'admin.login.title': '管理员登录',
  'admin.login.placeholder': '管理员密码',
  'admin.login.btn': '登录',
  'admin.codes.title': '访问码管理',
  'admin.codes.hint': '每个码自签发起一年内有效，使用次数由「次数」设定（默认 1 次）。同参占（周易 + 六壬 + 奇门）算一次消耗。',
  'admin.section.gen': '生成新码',
  'admin.field.count': '数量',
  'admin.field.maxuses': '次数',
  'admin.field.note': '备注',
  'admin.field.note.placeholder': '例如：给张三',
  'admin.btn.gen': '生成',
  'admin.btn.refresh': '刷新',
  'admin.section.filter': '查找与操作',
  'admin.field.kw': '关键词',
  'admin.field.kw.placeholder': '码 / 备注 / 使用 IP',
  'admin.field.status': '状态',
  'admin.status.all': '全部',
  'admin.status.unused': '仅未使用',
  'admin.status.partial': '仅使用中',
  'admin.status.used': '仅已用尽',
  'admin.status.expired': '仅已过期',
  'admin.btn.dlavail': '下载可用',
  'admin.btn.dlfilt': '下载筛选',
  'admin.btn.delexpired': '删除过期',
  'admin.btn.delused': '删除用尽',
  'admin.col.code': '访问码',
  'admin.col.status': '状态',
  'admin.col.uses': '次数',
  'admin.col.created': '创建时间',
  'admin.col.expires': '过期时间',
  'admin.col.lastused': '最后使用',
  'admin.col.note': '备注',
  'admin.col.action': '操作',
  'admin.btn.copy': '复制',
  'admin.copied': '已复制',
  'admin.stats.total': '共计',
  'admin.stats.avail': '可用',
  'admin.stats.partial': '使用中',
  'admin.stats.used': '已用尽',
  'admin.stats.expired': '已过期',
  'admin.stats.shown.prefix': '当前显示',
  'admin.stats.shown.suffix': '个',

  // ===== 首页页脚 =====
  'index.footer': '敬而用之　不戏不渎',

  // ===== 标题（页眉） =====
  'index.header.title': '卜　筮　明　心',
  'zhouyi.header.title': '周易占卦',
  'liuren.header.title': '大六壬',
  'qimen.header.title': '奇门遁甲',
  'huican.header.title': '周易 × 大六壬 × 奇门遁甲',

  // ===== 动态消息 =====
  'msg.cast.failed': '起卦失败',
  'msg.set.failed': '起局失败',
  'msg.network.error': '网络错误',
  'msg.access.required': '请输入访问码（向管理员索取）',
  'msg.access.invalid': '访问码无效或已使用',
  'msg.access.empty': '未输入访问码',
  'msg.copied': '已复制到剪贴板',
  'msg.copy.failed': '复制失败，请手动选择',
  'msg.qtype.load.fail': '加载问题类型失败',
  'msg.fill.upper.lower': '请填写上卦数与下卦数（整数即可）',

  // ===== 经纬度小标签 =====
  'loc.notset': '📍 经度未设',
  'loc.east': '📍 东经',
  'loc.west': '📍 西经',
  'loc.title': '点击修改起卦地经度（用于真太阳时校正）',
  'loc.prompt': '请输入起卦地经度（东经为正、西经为负，范围 -180 ~ 180）：\n例如：北京 116.4，上海 121.5，纽约 -74.0',
  'loc.invalid': '经度无效，已保留原值',

  // ===== 仪式动画文案 =====
  'ritual.liuren.title2': '月将加时',
  'ritual.liuren.hint2': '十二地支，依次轮转……',
  'ritual.qimen.title2': '布局起盘',
  'ritual.qimen.hint2': '三奇六仪，依次入宫……',
  'ritual.huican.title2': '三式同起',
  'ritual.huican.hint2': '钱 · 将 · 宫，三盘齐发……',
  'ritual.meditate': '静心凝神',
  'ritual.skip': '跳过仪式',
  'ritual.hint': '默念心中所问……',

  // ===== 按钮瞬时文字 =====
  'btn.casting': '起课中…',
  'btn.setting': '起局中…',
  'btn.divining': '起占中…',
  'btn.cast.short': '起卦',
  'btn.set.short': '起局',
  'btn.divine.short': '起占',
  'btn.course.short': '起课',

  // ===== admin 动态消息（新加） =====
  'admin.msg.session.expired': '会话已过期，请重新登录',
  'admin.msg.logged.out': '已登出',
  'admin.msg.confirm.logout': '确定要登出管理员会话吗？',
  'admin.msg.empty.codes': '尚无访问码，点击「生成」开始',
  'admin.msg.empty.match': '无匹配的访问码',

  // ===== 首页 · 序言 =====
  'index.preface.title': '静心',
  'index.preface.subtitle': '—— 读完这一篇，胜过十次占卦 ——',
  'index.preface.quote':
    '孔子说：「不占而已矣。」（《论语 · 子路》）<br/>' +
    '又说：「善为易者不占。」（《荀子 · 大略》）<br/>' +
    '<span class="src">—— 真正懂易的人，不靠占卜过日子。</span>',
  'index.preface.lead':
    '《易经》本是用来「沟通天地、辨识万物」的，不是用来玩、也不是替别人做决定的。<br/>' +
    '只有遇到真正难以判断的事情，发自至诚地去问，才会得到有意义的回应；把它当游戏，再清晰的卦也无法应验。<br/>' +
    '<strong>心存敬意才有灵验，轻慢则毫无意义。</strong>占卦的关键不在工具，而在人心的诚意。',

  // ===== 首页 · 六戒 =====
  'index.jie1.title': '没事别占卦',
  'index.jie1.p1': '心里没有疑惑、事情都明白，就<strong>不必占</strong>。占卦是用来解决「难以决断」的问题，不是用来预知一切的。',
  'index.jie1.p2': '有些人把占卦当游戏，看别人占自己也占，无事生事，只会扰乱心神。这样占出来的卦，反映的也只是你「无事可问」的茫然，没什么用。',
  'index.jie1.p3': '💡 如果心中没有具体困惑，先静下来想清楚问题再占。',

  'index.jie2.title': '一件事别占两次',
  'index.jie2.p1': '《易经》说：「初次问会告知；再三追问就是亵渎，亵渎了就不再告知。」第一卦就是答案；反复求问，说明你不信，神明就不会再回应。',
  'index.jie2.p2': '同一件事、同一个人、同一时段，只占一次。即使卦象不合心意，也要认真接受、深入思考；如果反复占直到得到满意答案，已经失去了占卦的本意，结果也不可信。',
  'index.jie2.p3': '💡 如果第一卦实在看不明白，可以过几天、调整心态再占；或者请懂行的人帮你解读，而不是连续追问。',

  'index.jie3.title': '小事别占',
  'index.jie3.p1': '琐事、闲事、口腹之欲、无关成败的事，<strong>不必动用卦象</strong>。',
  'index.jie3.p2': '「中午吃什么」「买哪件衣服」「下周要不要出去玩」——这种小事用常识、喜好、轻松决定就好。事事都占，会形成依赖，反而失去自主判断的能力。',
  'index.jie3.p3': '💡 占卦应该用在「重大抉择」「久拖不决」「影响深远」的事情上。轻重自己心里有数，别动不动就占。',

  'index.jie4.title': '不正当的事别占',
  'index.jie4.p1': '损人、害人、不伦、不义、违法、违德的事，<strong>不能占</strong>。',
  'index.jie4.p2': '《易经》追求顺应自然规律、合乎人心。用占卦为做坏事找理由，是借天道之名行私欲之恶——即使占到吉兆也是凶，占到凶兆也照样执迷。这种事根本不该做，何必占？',
  'index.jie4.p3': '💡 凡是心里有愧、不敢告诉别人的，就是「邪」。该停下来反思，而不是来占卜。',

  'index.jie5.title': '决定了的事别占',
  'index.jie5.p1': '心里其实已经有了主意，只是想找个「印证」才占的，<strong>不必占</strong>。',
  'index.jie5.p2': '这时候占到吉卦就高兴，占到凶卦就怀疑，甚至重占——卦象已经不是神明的指引，而成了你自己内心的镜子。不如直接去做，用结果检验。',
  'index.jie5.p3': '💡 真正该占的时候，是心中悬而未决、进退两难、左右都有道理的时候。',

  'index.jie6.title': '替别人占要慎重',
  'index.jie6.p1': '替别人占卦，必须<strong>得到对方的明确委托、清楚他要问什么</strong>，才可以代占。',
  'index.jie6.p2': '不能擅自给没有委托你的人占私事（比如暗中测他人吉凶、窥探隐私），这也属于亵渎。如果对方亲自委托，就以受托人的「年命」作为用神，凝神替他占。',
  'index.jie6.p3': '💡 代占时，问的人诚恳，占的人也才能诚恳；如果连受托人自己都漫不经心，这卦也没什么力量。',

  // ===== 首页 · 占卦流程 =====
  'index.ritual.mark.before': '占前',
  'index.ritual.mark.during': '占中',
  'index.ritual.mark.after': '占后',
  'index.ritual1.title': '一 · 净身静心',
  'index.ritual1.list':
    '<li>挑<strong>心境平静</strong>的时刻为好——不急不忙、不饿不撑、不怒不烦。</li>' +
    '<li>大醉、大怒、大悲、大喜的时候，心神散乱，不适合占卦。</li>' +
    '<li>条件允许的话，洗手、漱口、整衣，以示敬意。</li>' +
    '<li>找一个安静的地方，远离喧闹、电视、闲谈，独坐片刻。</li>',
  'index.ritual2.title': '二 · 想清楚要问什么',
  'index.ritual2.list':
    '<li>把要问的事化成<strong>一句明确的话</strong>，比如「这件事能不能做」「某人是否可信」「这个时间该动还是该静」。</li>' +
    '<li>问题别太宽泛（比如「我今年运势怎样」），也别太琐碎（比如「明天几点出门」），要有清楚的「对象 + 行为 + 判断」。</li>' +
    '<li>在心里默念问题三五遍，让念头清晰、专注。</li>' +
    '<li>问自己：这件事真的有疑虑吗？值得占吗？我此刻够诚恳吗？——任何一项不是，先暂缓。</li>',
  'index.ritual3.title': '三 · 起卦/起课',
  'index.ritual3.list':
    '<li>选方法：周易用铜钱 / 蓍草 / 数字；大六壬按当前时刻自动起课；奇门遁甲按当前时辰自动起局。</li>' +
    '<li>起卦时<strong>心念不能散</strong>——摇钱、拨蓍、按下「起占」键的瞬间，心里的问题应该还在。</li>' +
    '<li>自动起盘的方法（六壬、奇门、互参），按键前默念问题，按键后就停，别多按。</li>' +
    '<li>起卦的瞬间就是「天人感应」的时机，<strong>一击即定，不要反复</strong>。</li>',
  'index.ritual4.title': '四 · 看卦象、解卦',
  'index.ritual4.list':
    '<li>先看<strong>卦的总象</strong>（周易）/<strong>课的总象</strong>（六壬）/<strong>值符值使</strong>（奇门）——把握大意。</li>' +
    '<li>再看<strong>变爻</strong>（周易）/<strong>三传走势</strong>（六壬）/<strong>命中格局</strong>（奇门）——找到转折点。</li>' +
    '<li>最后看细节：爻辞象辞、天将六亲、三奇六仪。</li>' +
    '<li>可以用系统自动生成的 AI 提示词，请博学的 AI 帮忙解读；但最终判断还在你自己。</li>',
  'index.ritual5.title': '五 · 平心接受，不要执着',
  'index.ritual5.list':
    '<li>占到吉卦：别因高兴而骄傲、因骄傲而懈怠——吉的关键在「能去做」，不做就没有吉。</li>' +
    '<li>占到凶卦：别因惊慌而乱阵脚——凶的意义在「示警」，能避就能化凶为吉。</li>' +
    '<li>占到平卦：也要细想——往往是「时机未到」「宜守不宜进」的提示。</li>' +
    '<li><strong>卦象只是参考，决定权在你</strong>。易理的可贵，是启发人深思与担当，不是替你做决定。</li>',
  'index.ritual6.title': '六 · 行动并验证',
  'index.ritual6.list':
    '<li>占完之后<strong>按卦象的指引行动</strong>，过段时间再看是否应验。</li>' +
    '<li>无论应不应验，都是修学：应验则感受「天人之微」；不应验则反思「自己的盲点」。</li>' +
    '<li>不要在一件事占完后又反复重占求验证——那已经是亵渎，不是信任。</li>' +
    '<li>久而久之就能体会「易不过是我心的镜子」——至诚以问，卦象自然清晰。</li>',

  // ===== 首页 · 心法 =====
  'index.heart.do': '宜',
  'index.heart.dont': '忌',
  'index.xinfa.do':
    '<li><strong>敬</strong>：以对待先贤的敬意来面对，不嬉戏、不轻慢、不嘲弄。</li>' +
    '<li><strong>诚</strong>：心里所问就是心里所牵挂的；不诚的提问，必然得到不诚的答案。</li>' +
    '<li><strong>静</strong>：选安静的地方、安静的时刻、安静的心境去做。</li>' +
    '<li><strong>专</strong>：一时一问，思绪不要散乱。</li>' +
    '<li><strong>明</strong>：问得清楚、占得清楚、解得清楚；三明之后才能依卦而行。</li>' +
    '<li><strong>恒</strong>：长期记录占卦与应验情况，积累成学问；不是一时好奇。</li>',
  'index.xinfa.dont':
    '<li><strong>忌嬉戏</strong>：把卦当游戏、玩笑、谈资，一嬉则神散。</li>' +
    '<li><strong>忌亵渎</strong>：反复重占、一事多问，是亵渎。</li>' +
    '<li><strong>忌贪心</strong>：贪图吉象避开凶象而挑选着占；吉凶在自己的行动，不在卦象。</li>' +
    '<li><strong>忌怀疑</strong>：占完又不信，那不如不占。</li>' +
    '<li><strong>忌代占</strong>：未经当事人委托，不能擅自替他占。</li>' +
    '<li><strong>忌依赖</strong>：事事都占，会失去自主判断的能力。</li>',
  'index.xinfa.final.body':
    '《系辞》说：「<strong>君子平时观象玩辞，行动时观变玩占。</strong>」<br/>' +
    '真正的易学功夫，不在占了多少卦，而在平日里观察卦象、玩味爻辞，自省自修，用卦象养心。<br/>' +
    '占卦只是<strong>辅助</strong>——遇到难以判断的事，借天地之鉴稍作提点；如果能做到「善为易者不占」，才算真正得到了易的精髓。',
  'index.xinfa.final.sign': '—— 愿诸位慎用、敬用、诚用 ——',

  // ===== 首页 · 简介区 =====
  'index.intro.lead':
    '本系统把<strong>周易筮占</strong>、<strong>大六壬起课</strong>、<strong>奇门遁甲时家局</strong>三种古法集成到一起，并支持「三式互参」——一个问题三套同时起占、互相印证。' +
    '项目历经多年研制，<strong>对照历代原典反复推演与验证</strong>，力求每一步起卦、每一条断语都有出处、可考可信。' +
    '<span class="quote">—— 占法源于先贤，断语经由今人复校；汇集历代占家之精华，服务当下的实际问题。</span>',
  'index.intro.card1.title': '典籍考据',
  'index.intro.card1.list':
    '<li>周易：参照《周易本义》《易学启蒙》的变占法则，六十四卦的卦辞、爻辞、彖辞、象辞全部收录</li>' +
    '<li>大六壬：依据《六壬大全》《六壬指南》《毕法赋》全本，九宗十课、神煞、课体一应俱全</li>' +
    '<li>奇门：依据《烟波钓叟歌》《奇门遁甲统宗》，时家拆补起局、九星八门八神、格局自动识别</li>',
  'index.intro.card2.title': '校对严谨',
  'index.intro.card2.list':
    '<li>干支、节气、月将、旬空等全自动推算，采用真太阳时和精确历法</li>' +
    '<li>三式之间共享时刻与时令，互参时三盘同时、同源、可比对</li>' +
    '<li>经过长期迭代与实占校验，疑难卦例反复推敲，力求严谨</li>',
  'index.intro.card3.title': '三大优势',
  'index.intro.card3.list':
    '<li><strong>三式同堂</strong>：周易看大方向，六壬看人事，奇门看时机方位，三盘互补无死角</li>' +
    '<li><strong>类神直指</strong>：六类问题（事业/财运/感情/健康/抉择/时运）各有对应的用神，自动标记</li>' +
    '<li><strong>AI 解占就绪</strong>：每次起占自动生成结构化提示词，可直接交给 AI 深度解读</li>',
  'index.intro.card4.title': '适合谁用',
  'index.intro.card4.list':
    '<li>易学爱好者：想对照原典印证、复盘、推敲卦例</li>' +
    '<li>三式爱好者：想一问三占、交叉印证重大事宜</li>' +
    '<li>普通求问者：想在重大抉择时，借古法获得一份冷静的观照</li>',
  'index.intro.tips.title': '系统特色',
  'index.intro.tips.list':
    '<li><strong>原汁原味</strong>：卦辞、爻辞、课体、格局的名相，尽量保留古籍原文，不作浅白改写。</li>' +
    '<li><strong>结构清晰</strong>：每个盘分「起占过程 / 盘面 / 解读指引 / AI 提示词」四段呈现，新手老手都能上手。</li>' +
    '<li><strong>仪式感</strong>：起卦时保留摇钱、揲蓍、按时起局的仪式动画，凝神静心，一击即定。</li>' +
    '<li><strong>谨守占律</strong>：本系统秉持「占卦六戒」——没事不占、一事不二占、小事不占、邪事不占、已决不占、代占须慎，请大家敬而用之。</li>',

  // ===== 首页 · 免责声明 =====
  'index.disc.intro': '本系统是<strong>传统文化研习与占卜体验</strong>工具，所呈现的卦象、课象、局象及 AI 解读，仅供个人参考、自省与启发，<strong>不构成任何形式的专业建议</strong>。',
  'index.disc.list':
    '<li><strong>不替代医疗建议</strong>：涉及健康、疾病、用药、诊疗，请咨询执业医师，切勿以卦象代替正规医疗判断。</li>' +
    '<li><strong>不替代法律建议</strong>：涉及诉讼、合同、维权、赔偿，请咨询执业律师，卦象不能替代法律意见。</li>' +
    '<li><strong>不替代投资建议</strong>：涉及金融、股票、投资、创业，市场有风险，决策应基于自身研究与风险承受能力。</li>' +
    '<li><strong>不构成命运定论</strong>：易理重在「示警」「启发」「择时」，吉凶的最终结果还是取决于<strong>你自己的行动</strong>，不由卦象决定。</li>' +
    '<li><strong>用户自主决策</strong>：使用本系统得出的任何结论，应由<strong>你自己判断、自行承担相应后果</strong>，本系统及开发者不对由此引起的任何损失负责。</li>' +
    '<li><strong>不可用于不当用途</strong>：严禁将本系统用于赌博、损人、欺诈、违法违德的事；也不可在未经他人委托时为其占卜私事。</li>' +
    '<li><strong>数据与隐私</strong>：你输入的问题、命主资料（姓名、性别、出生年月日时、出生地、现居地）只存在本机浏览器，不会主动上传至第三方；其中只有性别与出生年参与六壬「年命」推算，姓名、出生地、现居地、月日时辰等<strong>不参与卦象算法</strong>，只是起卦时的仪式化记录；如使用 AI 提示词功能，请自行评估粘贴信息的敏感程度。</li>',
  'index.disc.final':
    '继续使用本系统，即表示你已<strong>知悉、理解并同意</strong>上述声明。<br/>' +
    '愿你以敬慎之心使用，把易理当作镜子，以自己的判断为主。',

  // ===== 首页 · FAQ =====
  'index.faq1.q': '现在 AI 这么强，为什么还需要这个工具？',
  'index.faq1.a1': 'AI 能<strong>解卦</strong>，但不能<strong>起卦</strong>。起卦讲究时、地、心、法——干支需要按真太阳时算，月将需要按节气来定，旬空、贵人、月建都有固定算法；让 AI 自己「想」出一个卦象，是无源之水。',
  'index.faq1.a2': '本系统的职责，是把<strong>古法起占的每一步</strong>严格还原出来：铜钱、揲蓍、时家拆补，都按原典推算；起出来的盘连同结构化提示词一起交给 AI，由 AI 来解。<strong>本系统负责「起」，AI 负责「解」</strong>，二者不是替代，而是分工。',
  'index.faq1.a3': '💡 如果直接让 AI「占一卦」，多半得到的是它编的故事；先用本系统起盘，再交给 AI 解读，才是正路。',

  'index.faq2.q': '起卦结果可信吗？是随机的还是有依据？',
  'index.faq2.a1': '周易筮占的铜钱、揲蓍法，本就是借<strong>偶然</strong>承接「天人感应」的契机——古人就是这样立法的，今人不必另寻他途。在这里使用随机数生成器<strong>正合古意</strong>，关键在于起卦时心念专一。',
  'index.faq2.a2': '大六壬、奇门遁甲<strong>完全由时刻起盘</strong>，没有随机成分：你按下起占键的那一刻，干支、月将、旬空、值符值使就唯一确定了，与本系统无关，也与占者无关。',
  'index.faq2.a3': '💡 「信」不在算法，而在你按下起占键时，心中所问是否清明、是否至诚。',

  'index.faq3.q': '同一件事，不同盘的结果不一样怎么办？',
  'index.faq3.a1': '「三式互参」的精髓<strong>正在于此</strong>。三盘各有所主：周易主<strong>大方向与义理</strong>，六壬主<strong>人事细节与应期</strong>，奇门主<strong>方位与谋划</strong>。三盘从不同角度切入，偶有出入是正常的。',
  'index.faq3.a2': '三盘都吉、都凶时最好判断；如果有分歧，就要看<strong>你问的事重在哪一面</strong>——问「该不该做」，看周易；问「什么时候能成、谁来帮忙」，看六壬；问「该往哪个方向、如何布局」，看奇门。',
  'index.faq3.a3': '💡 不要强求三盘完全一致。有分歧的地方，往往才是事情真正微妙之处。',

  'index.faq4.q': '我从未学过易，能直接用吗？',
  'index.faq4.a1': '可以。本系统为新手保留了<strong>完整的盘面与爻辞原文</strong>，并自动生成可直接喂给 AI 的提示词——不懂卦理也没关系，把提示词贴到任何 AI 对话框，就能得到通俗讲解。',
  'index.faq4.a2': '但要明白：AI 的解读<strong>是参考</strong>，不是终审。如果想深入理解、对应到自己的具体处境，仍需要自己读卦辞、想象义；本系统呈现的是<strong>原汁原味的古籍文本</strong>，可以边用边学。',
  'index.faq4.a3': '💡 与其追求「立刻看懂」，不如把每次起占当作一次和古人对话的机会。',

  'index.faq5.q': '我的问题会被上传或记录吗？',
  'index.faq5.a1': '你输入的问题、命主资料（姓名、性别、出生年月日时、出生地、现居地）<strong>只存在本机浏览器</strong>，不会主动上传到第三方。其中只有性别和出生年参与六壬「年命」推算，其余字段只作起卦记录与仪式化标识，不进入卦象算法。',
  'index.faq5.a2': '但如果你<strong>主动</strong>把提示词复制到外部 AI 服务（ChatGPT、Claude、文心一言等），那部分内容就进入了该服务的处理范围——是否敏感，请<strong>自行斟酌</strong>后再粘贴。',
  'index.faq5.a3': '💡 涉及他人姓名、隐私的事，贴入 AI 前可以先做匿名化处理（用「某甲」「某乙」代替）。',

  'index.faq6.q': '占到凶卦怎么办？会不会应验？',
  'index.faq6.a1': '凶卦的意义在<strong>「示警」而不是「定论」</strong>。古人讲「趋吉避凶」，正是因为凶象示警之后，<strong>人还有避开的机会</strong>——如果已经无法避免，何必示警？',
  'index.faq6.a2': '占到凶卦，要细察<strong>因何而凶</strong>：是时机不对？人不和？方位不利？自己心不诚？对症去做，往往凶能化为平、平能进为吉。<strong>真正可怕的不是凶卦，而是占到凶卦后还盲目行动</strong>。',
  'index.faq6.a3': '💡 千万不要因为占到凶卦就连着占求吉——那已经是「亵渎」，神明不再回应，只是自欺而已。',
};
