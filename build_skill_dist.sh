#!/usr/bin/env bash
# ⚠️ 本脚本已停用：占卜 skill 已抽取为独立发行仓库。
#
# skill（SKILL.md / README / 脚本）现位于独立仓库：
#     https://github.com/y001j/zhouyi-divination-skill
# 本仓库只保留 Go 源码，作为该 skill 二进制的「上游」。
#
# 发布新版二进制的方法（在 skill 仓库里跑，它会回到本源码目录交叉编译）：
#     cd ../zhouyi-divination-skill
#     git tag vX.Y.Z && git push --tags                # 先打版本 tag
#     bash scripts/release.sh --src /path/to/zhouyi    # 编译五平台 + 打离线 zip + 传 Release
#
# 用户侧无需本脚本：skill 仓库的 ensure_binary.sh 会在首次运行时
# 从 GitHub Release 自动下载对应平台二进制并缓存。
set -euo pipefail
echo "本脚本已停用。占卜 skill 已迁至独立仓库 y001j/zhouyi-divination-skill。" >&2
echo "发布二进制请在该仓库执行 scripts/release.sh，详见本文件注释。" >&2
exit 1
