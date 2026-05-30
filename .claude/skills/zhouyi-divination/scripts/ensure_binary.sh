#!/usr/bin/env bash
# 选出适配当前平台的 zhouyi 二进制，把其绝对路径输出到 stdout。
#
# 独立分发包模式：二进制随包附带在 ../bin/ 下，按 OS/架构挑选，开箱即用，无需装 Go。
# 仅当 bin/ 里没有当前平台的预编译产物时，才回退到「在项目根 go build」（需源码+Go）。
#
# 用法：BIN=$(bash ensure_binary.sh) && "$BIN" cast ...
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILL_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"   # skill 根目录
BIN_DIR="$SKILL_DIR/bin"

# —— 识别平台 ——
raw_os="$(uname -s)"
raw_arch="$(uname -m)"
case "$raw_os" in
  Darwin)            os="darwin" ;;
  Linux)             os="linux" ;;
  MINGW*|MSYS*|CYGWIN*) os="windows" ;;
  *)                 os="unknown" ;;
esac
case "$raw_arch" in
  x86_64|amd64)  arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *)             arch="unknown" ;;
esac

ext=""
[[ "$os" == "windows" ]] && ext=".exe"
candidate="$BIN_DIR/zhouyi-${os}-${arch}${ext}"

# —— 1. 优先用随包预编译二进制 ——
if [[ -f "$candidate" ]]; then
  chmod +x "$candidate" 2>/dev/null || true
  echo "$candidate"
  exit 0
fi

# —— 2. 回退：若包内嵌了源码且本机有 Go，则现场编译 ——
# （独立分发包默认不含源码，此分支主要服务于「连源码一起拿到」的开发者）
PROJECT_ROOT="$(cd "$SKILL_DIR/../../.." && pwd)"   # 若 skill 仍嵌在项目内
if [[ -f "$PROJECT_ROOT/go.mod" ]] && command -v go >/dev/null 2>&1; then
  built="$PROJECT_ROOT/zhouyi"
  echo "[ensure_binary] 未找到 ${os}-${arch} 预编译二进制，尝试用本机 Go 编译 ..." >&2
  if ( cd "$PROJECT_ROOT" && go build -o "$built" . ) >&2; then
    echo "$built"
    exit 0
  fi
fi

# —— 3. 都不行：明确报错 ——
echo "[ensure_binary] 无法获得可用的 zhouyi 二进制。" >&2
echo "[ensure_binary] 当前平台：${os}-${arch}（uname: $raw_os / $raw_arch）" >&2
echo "[ensure_binary] 包内 bin/ 提供的平台：" >&2
ls "$BIN_DIR" 2>/dev/null | sed 's/^/  - /' >&2 || echo "  （bin/ 为空或不存在）" >&2
echo "[ensure_binary] 解决办法：" >&2
echo "  1) 若你的平台不在上表，请向分发者索取对应平台的二进制；或" >&2
echo "  2) 拿到完整 Go 源码后，在项目根执行 go build 自行编译。" >&2
exit 1
