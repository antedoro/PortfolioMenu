#!/bin/bash

# ==========================================================
# PortfolioMenu - Build script per macOS
# Genera l'icona (.icns) e compila per
# Apple Silicon (arm64) e Intel (amd64).
# ==========================================================

set -euo pipefail

# --- Configurazione ---------------------------------------
APP_NAME="PortfolioMenu"
BUNDLE_ID="com.antedoro.Portfoliomenu"
SRC_ICON="icon.png"
VERSION="${1:-dev}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

OUT_DIR="build"
ICONSET="$OUT_DIR/${APP_NAME}.iconset"
ICNS="$OUT_DIR/AppIcon.icns"

echo "PortfolioMenu macOS build"
echo "=========================="

# --- 1. Generazione icona ---------------------------------
echo "[1/4] Generazione icona da $SRC_ICON ..."

if [ ! -f "$SRC_ICON" ]; then
  echo "Errore: $SRC_ICON non trovata"
  exit 1
fi

rm -rf "$ICONSET"
mkdir -p "$ICONSET"

# Dimensioni richieste da iconutil (con varianti @2x).
declare -a SIZES=(16 32 64 128 256 512 1024)

for s in "${SIZES[@]}"; do
  # versione 1x
  sips -z "$s" "$s" "$SRC_ICON" \
    --out "$ICONSET/icon_${s}x${s}.png" >/dev/null
  # versione @2x (doppia risoluzione)
  d=$((s * 2))
  sips -z "$d" "$d" "$SRC_ICON" \
    --out "$ICONSET/icon_${s}x${s}@2x.png" >/dev/null
done

iconutil --convert icns \
  --output "$ICNS" "$ICONSET"

echo "      -> $ICNS"

# --- 2. Compilazione --------------------------------------
echo "[2/4] Compilazione binari (CGO abilitato) ..."

mkdir -p "$OUT_DIR/arm64" "$OUT_DIR/amd64" "$OUT_DIR/universal"

# Apple Silicon
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 \
  go build -trimpath -ldflags "-s -w" \
  -o "$OUT_DIR/arm64/$APP_NAME" ./cmd/portfoliomenu

# Intel
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 \
  go build -trimpath -ldflags "-s -w" \
  -o "$OUT_DIR/amd64/$APP_NAME" ./cmd/portfoliomenu

# Binario universale
lipo -create -output \
  "$OUT_DIR/universal/$APP_NAME" \
  "$OUT_DIR/arm64/$APP_NAME" \
  "$OUT_DIR/amd64/$APP_NAME"

echo "      arm64 / amd64 / universal pronti"

# --- 3. Creazione .app bundle -----------------------------
make_app() {
  local arch="$1"      # arm64 | amd64 | universal
  local bin="$2"       # percorso binario
  local app="$OUT_DIR/${APP_NAME}-${arch}.app"

  echo "[3/4] Creazione $app ..."

  rm -rf "$app"
  mkdir -p "$app/Contents/MacOS/configs" \
           "$app/Contents/Resources"

  cp "$bin" "$app/Contents/MacOS/$APP_NAME"
  cp "$ICNS" "$app/Contents/Resources/AppIcon.icns"
  cp configs/portfoliomenu.toml \
     "$app/Contents/MacOS/configs/portfoliomenu.toml"

  cat > "$app/Contents/Info.plist" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
 "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleName</key>
  <string>${APP_NAME}</string>
  <key>CFBundleDisplayName</key>
  <string>${APP_NAME}</string>
  <key>CFBundleIdentifier</key>
  <string>${BUNDLE_ID}</string>
  <key>CFBundleExecutable</key>
  <string>${APP_NAME}</string>
  <key>CFBundleVersion</key>
  <string>${VERSION}</string>
  <key>CFBundleShortVersionString</key>
  <string>${VERSION}</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
  <key>CFBundleIconFile</key>
  <string>AppIcon</string>
  <key>LSMinimumSystemVersion</key>
  <string>10.13</string>
  <key>LSUIElement</key>
  <true/>
  <key>NSHighResolutionCapable</key>
  <true/>
</dict>
</plist>
PLIST

  chmod +x "$app/Contents/MacOS/$APP_NAME"
}

make_app "arm64"    "$OUT_DIR/arm64/$APP_NAME"
make_app "amd64"    "$OUT_DIR/amd64/$APP_NAME"
make_app "universal" "$OUT_DIR/universal/$APP_NAME"

# --- 4. Fine ----------------------------------------------
echo "[4/4] Fatto."
echo ""
echo "Bundle generati in $OUT_DIR/:"
ls -1 "$OUT_DIR"/*.app
echo ""
echo "Per avviare (menubar):"
echo "  open $OUT_DIR/${APP_NAME}-universal.app"
