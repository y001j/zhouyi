# 周易三式 · 中式术数推演系统

> 易有太极，是生两仪，两仪生四象，四象生八卦。——《周易·系辞上》

一个用 Go 实现的中式术数起盘与解读系统，涵盖 **周易六爻、奇门遁甲、大六壬** 三式，以及三式互参合断。它不仅能排出准确的卦象/盘面，还会为每一盘生成一份**结构化的解卦提示词**，交给 AI（如 Claude）按传统典籍的解卦规则做深入解读。

提供三种使用形态：**命令行交互**、**Web 网页**、以及 **Claude Code Skill**（让任意 AI agent 一句话起卦）。

> ⚠️ 本项目基于传统术数典籍实现，仅供学习、研究与娱乐。占卜不能预测未来，**不构成医疗、法律、投资等任何专业建议**，请勿用于迷信或重大决策。

## 界面预览

Web 形态采用古典宣纸风格，繁简白话三体可切换：

| 周易筮占 | 奇门遁甲 |
|---|---|
| ![周易筮占](assets/screenshot-zhouyi.jpg) | ![奇门遁甲](assets/screenshot-qimen.jpg) |

## 特性

- **三式齐备**
  - 🔮 **周易六爻**：铜钱法 / 蓍草法（大衍揲蓍）/ 数字起卦，含变卦、互卦、错卦、综卦、爻位中正应比分析
  - 🧭 **奇门遁甲**：时家奇门，排九宫盘（星/门/神/干）、值符值使、格局命中、类神直指
  - 🎴 **大六壬**：月将加时起四课三传、十二天将、毕法赋断语、年命行年
  - ☯️ **三式互参**：同一时刻三盘合一，跨系互证
- **真太阳时校正**：按经度校正起盘时刻
- **结构化解卦提示词**：每盘自动生成给 AI 的解卦任务书，含卦象材料 + 解卦规则 + 问题侧重
- **因术制宜**：周易主义理、奇门主方位时机、六壬主人事应期，各术的提示词定位与解卦结构各有侧重
- **三种形态**：CLI / Web / Claude Code Skill

## 快速开始

需要 [Go 1.26+](https://go.dev/dl/)。

```bash
git clone https://github.com/y001j/zhouyi.git
cd zhouyi
go build -o zhouyi .
```

### 1. 命令行交互

```bash
./zhouyi
```
进入交互界面后输入 `coin`（铜钱起卦）、`liuren`（六壬）、`huican`（互参）、`help` 等。

### 2. 非交互起盘（输出 JSON，供脚本/agent 调用）

```bash
# 周易
./zhouyi cast -m zhouyi -q "今年事业如何" -t career
# 奇门
./zhouyi cast -m qimen  -q "新店该往哪个方向" -t career
# 六壬
./zhouyi cast -m liuren -q "这件事何时有结果" -t timing
# 三式互参
./zhouyi cast -m huican -q "明年整体运势" -t timing
```
输出 JSON，`prompt` 字段即解卦提示词。运行 `./zhouyi cast --help` 看全部参数。

输出示例（`prompt` 较长，此处省略其正文）：

```json
{
  "ok": true,
  "method": "zhouyi",
  "question": "今年事业发展如何",
  "questionLabel": "事业 / 工作",
  "summary": "第37卦 家人卦 → 之 第13卦 同人卦",
  "prompt": "你是一位精通《周易》的易学顾问，请根据以下卦象为我解卦……（完整解卦提示词，约 2700 字）"
}
```

AI（如 Claude）拿到 `prompt` 后，会按其中「请按以下结构解卦」的步骤逐条解读，并为每个结论标注卦象依据。

### 3. Web 网页

```bash
./zhouyi serve 8080
```
浏览器打开 http://localhost:8080/ ，含周易 / 奇门 / 六壬 / 互参各页面。
（Web 的管理后台需设置 `ADMIN_PASSWORD` 环境变量或 `config.json`，仅生成访问码时用到。）

## 作为 Claude Code Skill 使用

本项目内置一个 Claude Code skill（`/.claude/skills/zhouyi-divination/`），让你在 Claude Code 里直接说「帮我算一卦」「用奇门看方位」即可自动起盘解卦。

**两种获取方式：**

- **下载预编译包（推荐，无需装 Go）**：到本仓库 [Releases](https://github.com/y001j/zhouyi/releases) 下载 `zhouyi-divination-skill.zip`，解压后把 `zhouyi-divination/` 放进 `~/.claude/skills/` 即可。
- **从源码自带**：克隆本仓库后，skill 已在 `.claude/skills/zhouyi-divination/`，运行 `bash build_skill_dist.sh` 会交叉编译多平台二进制并打包。

详见 [skill 的 README](.claude/skills/zhouyi-divination/README.md)。

## 项目结构

```
.
├── main.go            CLI 入口（交互 / cast / serve）
├── cast.go            非交互起盘子命令（供 skill/agent 调用）
├── divination.go      周易起卦核心
├── hexagrams.go       六十四卦数据
├── prompt.go          周易解卦提示词生成
├── huican.go          三式互参
├── server.go          Web HTTP 服务
├── qimen/             奇门遁甲（起局/排盘/格局/提示词）
├── liuren/            大六壬（起课/四课三传/天将/毕法赋/提示词）
├── web/               Web 前端静态资源
└── .claude/skills/zhouyi-divination/   Claude Code Skill
```

## 构建多平台分发包

```bash
bash build_skill_dist.sh
```
会交叉编译 darwin-arm64/amd64、linux-amd64/arm64、windows-amd64 五个平台的二进制，打包成 `dist/zhouyi-divination-skill.zip`。

## 致谢

实现参考了《周易》《易学启蒙》《奇门遁甲统宗大全》《奇门宝鉴》《大六壬大全》《大六壬毕法赋》等传统典籍。感谢 [lunar-go](https://github.com/6tail/lunar-go) 提供的农历/干支/节气历法支持。

## 许可证

[Apache License 2.0](LICENSE)。

传统术数典籍本身属于公有领域文化遗产；本项目的代码实现以 Apache 2.0 开源。占卜内容仅供参考，使用者需自行判断，作者不对任何依据本工具所做决策的后果负责。
