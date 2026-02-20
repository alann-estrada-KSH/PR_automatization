#!/usr/bin/env bash
# ─────────────────────────────────────────────────────────────────────────────
# prgen installer — macOS & Linux
# Usage: bash scripts/install.sh
# ─────────────────────────────────────────────────────────────────────────────
set -e

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION_FILE="$REPO_ROOT/VERSION"
BIN_NAME="prgen"

# Read version
VERSION="$(cat "$VERSION_FILE" 2>/dev/null | tr -d '[:space:]')"
if [ -z "$VERSION" ]; then VERSION="dev"; fi

BUILD_DATE="$(date -u +%Y-%m-%d)"

echo ""
echo " 🔧 Building prgen v$VERSION..."

# ── Build ─────────────────────────────────────────────────────────────────────
cd "$REPO_ROOT"
go build \
  -ldflags "-s -w -X github.com/ksh/prgen/internal/version.Version=$VERSION -X github.com/ksh/prgen/internal/version.BuildDate=$BUILD_DATE" \
  -o "$BIN_NAME" \
  ./cmd/prgen

echo " ✅ Build exitoso"

# ── Install location ──────────────────────────────────────────────────────────
if [ -w "/usr/local/bin" ]; then
  INSTALL_DIR="/usr/local/bin"
else
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
fi

cp "$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"
chmod +x "$INSTALL_DIR/$BIN_NAME"
echo " 📦 Binario instalado en: $INSTALL_DIR/$BIN_NAME"

# ── PATH configuration ────────────────────────────────────────────────────────
add_to_path() {
  local FILE="$1"
  local LINE="export PATH=\"$INSTALL_DIR:\$PATH\""
  if [ -f "$FILE" ] && ! grep -qF "$INSTALL_DIR" "$FILE"; then
    echo "" >> "$FILE"
    echo "# prgen" >> "$FILE"
    echo "$LINE" >> "$FILE"
    echo " 🛤️  PATH actualizado en $FILE"
  fi
}

if [ "$INSTALL_DIR" != "/usr/local/bin" ]; then
  add_to_path "$HOME/.bashrc"
  add_to_path "$HOME/.zshrc"
  # Fish
  FISH_CONFIG="$HOME/.config/fish/config.fish"
  if [ -f "$FISH_CONFIG" ] && ! grep -qF "$INSTALL_DIR" "$FISH_CONFIG"; then
    echo "" >> "$FISH_CONFIG"
    echo "# prgen" >> "$FISH_CONFIG"
    echo "fish_add_path $INSTALL_DIR" >> "$FISH_CONFIG"
    echo " 🐟 PATH actualizado en $FISH_CONFIG"
  fi
fi

# ── Default config ────────────────────────────────────────────────────────────
CONFIG_DIR="$HOME/.prgen"
mkdir -p "$CONFIG_DIR"

if [ ! -f "$CONFIG_DIR/config.yaml" ]; then
  cp "$REPO_ROOT/config.yaml" "$CONFIG_DIR/config.yaml"
  echo " ⚙️  Config copiada a $CONFIG_DIR/config.yaml"
fi

# Cleanup build artifact
rm -f "$REPO_ROOT/$BIN_NAME"

# ── Verify ────────────────────────────────────────────────────────────────────
echo ""
echo " 🎉 ¡prgen instalado!"
echo ""

# Try to run version (might need to reload PATH)
if command -v "$BIN_NAME" &>/dev/null; then
  "$BIN_NAME" version
else
  echo " ⚠️  Reinicia tu terminal o ejecuta:"
  echo "    source ~/.bashrc   (bash)"
  echo "    source ~/.zshrc    (zsh)"
  echo ""
  echo "    Luego: prgen version"
fi
echo ""
