#!/usr/bin/env bash
# 一键产出「周易三式占卜」独立可分发 skill 包。
# 流程：交叉编译五平台二进制 → 放入 skill/bin → 打成 zip。
# 产物：dist/zhouyi-divination-skill.zip，解压即用，对方无需装 Go。
#
# 用法：bash build_skill_dist.sh
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"

SKILL="$ROOT/.claude/skills/zhouyi-divination"
BIN="$SKILL/bin"
OUT="$ROOT/dist"
ZIP="$OUT/zhouyi-divination-skill.zip"

mkdir -p "$BIN" "$OUT"

echo "==> 交叉编译五平台二进制（最新代码）"
build() {
  local os=$1 arch=$2 ext=${3:-}
  local out="$BIN/zhouyi-${os}-${arch}${ext}"
  GOOS="$os" GOARCH="$arch" go build -trimpath -ldflags="-s -w" -o "$out" .
  printf "    ✅ %-22s %s\n" "${os}-${arch}${ext}" "$(du -h "$out" | cut -f1)"
}
build darwin  arm64
build darwin  amd64
build linux   amd64
build linux   arm64
build windows amd64 .exe

echo "==> 打包 zip"
rm -f "$ZIP"
# 用子 shell 进入 skills 父目录，让 zip 内的顶层就是 zhouyi-divination/
( cd "$SKILL/.." && zip -rq "$ZIP" zhouyi-divination \
    -x '*.DS_Store' )

echo ""
echo "==> 完成"
echo "    分发包：$ZIP"
echo "    大小：  $(du -h "$ZIP" | cut -f1)"
echo "    内容："
unzip -l "$ZIP" | awk 'NR>3 && $4!="" {print "      "$4}' | grep -v '/$' || true
echo ""
echo "    把这个 zip 发给别人，对方解压后将 zhouyi-divination/ 放进"
echo "    ~/.claude/skills/ 或 <项目>/.claude/skills/ 即可使用。"
