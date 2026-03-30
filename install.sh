#!/bin/bash

echo "Installing mytunnel..."

OS=$(uname)
ARCH=$(uname -m)

if [ "$OS" = "Linux" ]; then
    URL="https://github.com/DpkRn/devtunnel/releases/latest/download/mytunnel-linux"
elif [ "$OS" = "Darwin" ]; then
    if [ "$ARCH" = "arm64" ]; then
        URL="https://github.com/DpkRn/devtunnel/releases/latest/download/mytunnel-mac-arm64"
    else
        URL="https://github.com/DpkRn/devtunnel/releases/latest/download/mytunnel-mac"
    fi
else
    echo "Unsupported OS: $OS $ARCH"
    exit 1
fi

curl -fSL --progress-bar "$URL" -o mytunnel </dev/tty

if file mytunnel | grep -qv 'text'; then
    chmod +x mytunnel
    sudo mv mytunnel /usr/local/bin/
else
    echo "❌ Download failed — file is not a binary:"
    cat mytunnel
    rm -f mytunnel
    exit 1
fi

echo "✅ Installed! Run: mytunnel http 3000"