#!/bin/bash

ARCH=$(uname -m)
OS="$(uname -s)"
REPO="error-404-h/tg-keys"

if [ "$OS" = "Linux" ]; then
    if [ "$ARCH" = "x86_64" ]; then
        FILE="tg-keys-linux64"
    elif [ "$ARCH" = "i686" ]; then
        FILE="tg-keys-linux32"
    elif [ "$ARCH" = "aarch64" ]; then
        echo "For Termux/Android, please use: apt install tg-keys"
        exit 1
    fi
elif [ "$OS" = "Darwin" ]; then
    FILE="tg-keys-mac"
else
    echo "Unsupported OS"
    exit 1
fi

curl -L "https://github.com/$REPO/releases/latest/download/$FILE" -o tg-keys
chmod +x tg-keys
echo "Installation complete. Run it with: ./tg-keys"
