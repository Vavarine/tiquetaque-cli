#!/usr/bin/env bash

set -e

REPO="Vavarine/tiquetaque-cli"   # troque pelo seu repo
BINARY="ttq"              # nome do binário
VERSION="${VERSION:-latest}"

# Detecta sistema e arquitetura
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64 | arm64) ARCH="arm64" ;;
    *) echo "Arquitetura não suportada: $ARCH"; exit 1 ;;
esac

echo https://api.github.com/repos/$REPO/releases/latest

# Pega a URL da versão
if [ "$VERSION" = "latest" ]; then
    DOWNLOAD_URL=$(curl -s https://api.github.com/repos/$REPO/releases/latest \
      | grep "browser_download_url" \
      | cut -d '"' -f 4)
else
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/${BINARY}_${OS}_${ARCH}"
fi

echo "Baixando $BINARY de $DOWNLOAD_URL ..."

curl -L "$DOWNLOAD_URL" -o "/tmp/$BINARY"
chmod +x "/tmp/$BINARY"

# Instala em /usr/local/bin (pode precisar de sudo)
echo "Instalando em /usr/local/bin ..."
sudo mv "/tmp/$BINARY" "/usr/local/bin/$BINARY"

echo "✅ Instalado com sucesso! Execute '$BINARY --help'"
