# 周易三式占卜 · Claude Code Skill

一个把中式术数（**周易六爻 / 奇门遁甲 / 大六壬 / 三式互参**）封装成的 Claude Code skill。
你对它说「帮我算一卦」「用奇门看看方位」「六壬断断这个人」，它就会**起盘排卦**，并产出一份结构化的解卦材料，由 Claude 按传统解卦规则替你解读。

> ⚠️ 占卜结果仅供参考与启发，是换个角度思考问题的工具，**不构成医疗、法律、投资等任何专业建议**。

## 这个包是什么

一个**自包含、开箱即用**的 skill 目录：

```
zhouyi-divination/
├── SKILL.md                # 给 Claude 读的说明书（何时触发、怎么调用）
├── README.md               # 本文件，给人看的安装说明
├── scripts/
│   └── ensure_binary.sh    # 按你的平台挑出对应二进制
└── bin/                    # 五个平台的预编译二进制（无需装 Go）
    ├── zhouyi-darwin-arm64       # Apple 芯片 Mac
    ├── zhouyi-darwin-amd64       # Intel Mac
    ├── zhouyi-linux-amd64
    ├── zhouyi-linux-arm64
    └── zhouyi-windows-amd64.exe
```

**无需安装 Go、无需联网、无需配置**——二进制已随包附带，本地离线运行。

## 安装（三步）

1. **拷贝**整个 `zhouyi-divination/` 目录到你的 Claude Code skills 目录：
   - 项目级（只在该项目可用）：`<你的项目>/.claude/skills/zhouyi-divination/`
   - 用户级（所有项目可用）：`~/.claude/skills/zhouyi-divination/`

2. （macOS/Linux）给脚本和二进制可执行权限（通常拷贝后已保留；若不放心可执行）：
   ```bash
   chmod +x zhouyi-divination/scripts/ensure_binary.sh
   chmod +x zhouyi-divination/bin/*
   ```

3. 在 Claude Code 里直接说占卜需求即可，例如：
   - 「帮我算一卦，问问明年工作运势」
   - 「用奇门看看我该往哪个方向发力」
   - 「六壬断断这件事什么时候有结果」

   Claude 会自动识别并激活本 skill。

## 支持的平台

| 平台 | 二进制 |
|---|---|
| Apple 芯片 Mac (M1/M2/M3…) | `zhouyi-darwin-arm64` |
| Intel Mac | `zhouyi-darwin-amd64` |
| Linux x86_64 | `zhouyi-linux-amd64` |
| Linux ARM64 | `zhouyi-linux-arm64` |
| Windows x86_64 | `zhouyi-windows-amd64.exe` |

不在此列（如某些国产架构）？拿到 Go 源码后，在项目根 `go build` 自行编译，`ensure_binary.sh` 会自动回退到这条路径。

## 手动试一下（可选）

不经过 Claude，直接命令行验证它能跑：

```bash
cd zhouyi-divination
BIN=$(bash scripts/ensure_binary.sh)
"$BIN" cast -m zhouyi -q "今年事业如何" -t career
```

会输出一段 JSON，其中 `prompt` 字段就是解卦材料，`summary` 是一句话盘面摘要。

## 四种术怎么选

| 术 | 所长 | 适合 |
|---|---|---|
| 周易 `zhouyi` | 明义理、辨吉凶、示进退 | 该不该做、人生方向、心态决策 |
| 奇门 `qimen` | 谋大局、定方位、择时机 | 选址方位、出行择时、布局谋略 |
| 六壬 `liuren` | 看人心、断曲折、定应期 | 具体人事、对方心意、何时应验 |
| 互参 `huican` | 三式合参 | 重大或复杂之事，多角度印证 |

不指定时，Claude 会按问题性质替你挑。

## 常见问题

**Q：对方会泄露我的什么数据吗？**
A：完全本地运行，不联网、不上传。起盘只用系统当前时间（或你指定的时间）。

**Q：macOS 提示「无法验证开发者」打不开二进制？**
A：因为二进制未签名。可执行 `xattr -d com.apple.quarantine zhouyi-divination/bin/*` 解除隔离，或在「系统设置 → 隐私与安全性」里允许。

**Q：每次结果都一样吗？**
A：奇门/六壬/互参在同一时刻同一问题下结果确定；周易默认用铜钱法（随机起卦），每次不同——这符合「即时起卦」的传统。

## 许可与免责

本工具基于传统术数典籍实现，仅供学习、研究与娱乐。占卜不能预测未来，请勿用于迷信或重大决策。涉及健康、法律、财务等，请咨询相应领域的专业人士。
