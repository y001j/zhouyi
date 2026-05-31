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
- **真太阳时校正**：按经度校正起盘时刻。Web / Skill / cast 三端默认启用，未提供经度时回退默认经度（北京 116.4°E）；cast 可用 `--lon` 指定经度、`--no-truesolar` 关闭
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

## AI 解卦（调用 LLM API）

每一式起盘后都会生成一份结构化「解卦提示词」。在 `config.json` 配好 `llm` 段后，程序可**直接调用大模型 API** 给出解读，无需手动复制提示词。

复制 `config.example.json` 为 `config.json`，填入 `llm` 段：

```json
{
  "adminPassword": "你的管理员密码",
  "llm": {
    "enabled": true,
    "provider": "openai",
    "baseURL": "https://api.openai.com/v1",
    "apiKey": "sk-...",
    "model": "gpt-4o-mini",
    "temperature": 0.7,
    "maxTokens": 2048,
    "timeoutSec": 120
  }
}
```

不配（或 `enabled:false` / `apiKey` 为空）则**优雅降级**为「仅出提示词，自行复制到 AI」，不影响起盘。

**两种协议**（由 `provider` 决定，一套配置可对接绝大多数模型）：

| provider | 接口 | 适配 |
|---|---|---|
| `openai`（默认） | `/chat/completions` | OpenAI、DeepSeek、月之暗面 Kimi、智谱、通义、Ollama、各类中转网关 |
| `anthropic` | `/v1/messages` | Claude 官方或兼容网关 |

`baseURL` / `model` 留空时按 `provider` 取默认值。

**两个入口：**

- **交互式 CLI**：起卦 → 出提示词 → 询问「是否直接调用 AI 解卦」，回车即在终端打印模型解读。
- **Web 服务**：周易 / 奇门 / 六壬 / 互参各页起盘后，提示词下方会出现「**召 AI 解卦**」按钮，点击即在页面内显示模型解读。底层为 `POST /api/interpret`，请求体 `{"prompt": "解卦提示词"}`，返回 `{"interpretation": "..."}`。该端点不额外消耗访问码（解卦是起卦的延续），以「服务端是否配置 LLM」为天然闸门——未配置即返回 503，前端自动提示改用「复制提示词」。

**环境变量覆盖**（便于部署不落盘密钥）：
`ZHOUYI_LLM_ENABLED`、`ZHOUYI_LLM_PROVIDER`、`ZHOUYI_LLM_BASE_URL`、`ZHOUYI_LLM_API_KEY`、`ZHOUYI_LLM_MODEL`、`ZHOUYI_LLM_MAX_TOKENS`。

## 作为 Claude Code Skill 使用

本项目配套一个 Claude Code skill，让你在 Claude Code 里直接说「帮我算一卦」「用奇门看方位」即可自动起盘解卦。该 skill 已独立为发行仓库：

> **👉 https://github.com/y001j/zhouyi-divination-skill**

**两种安装方式：**

- **克隆即用（在线，推荐，无需装 Go）**：
  ```bash
  git clone https://github.com/y001j/zhouyi-divination-skill.git ~/.claude/skills/zhouyi-divination
  ```
  首次占卜时会自动从 Release 下载对应平台二进制并缓存。
- **离线 zip**：到 skill 仓库的 [Releases](https://github.com/y001j/zhouyi-divination-skill/releases) 下载 `zhouyi-divination-skill-offline.zip`，解压后把 `zhouyi-divination/` 放进 `~/.claude/skills/`，完全不联网。

本仓库（Go 源码）是该 skill 二进制的**上游**：skill 仓库的 `scripts/release.sh` 会回到这里交叉编译多平台二进制并发到 Release。详见 skill 仓库的 README。

## 项目结构

```
.
├── main.go            CLI 入口（交互 / cast / serve）
├── cast.go            非交互起盘子命令（供 skill/agent 调用）
├── config.go          统一配置加载（含 llm 段 + 环境变量覆盖）
├── llm.go             AI 解卦客户端（OpenAI 兼容 / Anthropic 双协议）
├── divination.go      周易起卦核心
├── hexagrams.go       六十四卦数据
├── prompt.go          周易解卦提示词生成
├── huican.go          三式互参
├── server.go          Web HTTP 服务
├── qimen/             奇门遁甲（起局/排盘/格局/提示词）
├── liuren/            大六壬（起课/四课三传/天将/毕法赋/提示词）
└── web/               Web 前端静态资源
```

> Claude Code skill 已独立至 [y001j/zhouyi-divination-skill](https://github.com/y001j/zhouyi-divination-skill)，本仓库只保留 Go 源码作为其二进制上游。

## 发布 skill 的多平台二进制

skill 的二进制由其独立仓库的发布脚本产出（会回到本源码目录交叉编译）：

```bash
cd ../zhouyi-divination-skill
git tag vX.Y.Z && git push --tags
bash scripts/release.sh --src /path/to/zhouyi
```
会交叉编译 darwin-arm64/amd64、linux-amd64/arm64、windows-amd64 五个平台的二进制，打离线 zip，并作为附件发到 skill 仓库的 Release。

## 致谢

实现参考了《周易》《易学启蒙》《奇门遁甲统宗大全》《奇门宝鉴》《大六壬大全》《大六壬毕法赋》等传统典籍。感谢 [lunar-go](https://github.com/6tail/lunar-go) 提供的农历/干支/节气历法支持。

## 许可证

[Apache License 2.0](LICENSE)。

传统术数典籍本身属于公有领域文化遗产；本项目的代码实现以 Apache 2.0 开源。占卜内容仅供参考，使用者需自行判断，作者不对任何依据本工具所做决策的后果负责。
