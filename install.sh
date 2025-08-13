#!/bin/bash
set -e

OS=$(uname -s | tr 'A-Z' 'a-z')

ARCH=$(uname -m)
case "$ARCH" in
	x86_64)
		ARCH=amd64
		;;
	aarch64|arm64)
		ARCH=arm64
		;;
	*)
		echo "Unsupported architecture: $ARCH" >&2
		exit 1
		;;
esac

echo "Resolving latest version..."

VERSION=$(curl -sL https://api.github.com/repos/coalaura/idk/releases/latest | grep -Po '"tag_name": *"\K.*?(?=")')

if ! printf '%s\n' "$VERSION" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+$'; then
	echo "Error: '$VERSION' is not in vMAJOR.MINOR.PATCH format" >&2
	exit 1
fi

rm -f /tmp/idk

BIN="idk_${VERSION}_${OS}_${ARCH}"
URL="https://github.com/coalaura/idk/releases/download/${VERSION}/${BIN}"

echo "Downloading ${BIN}..."

if ! curl -sL "$URL" -o /tmp/idk; then
	echo "Error: failed to download $URL" >&2
	exit 1
fi

trap 'rm -f /tmp/idk' EXIT

chmod +x /tmp/idk

echo "Installing to /usr/local/bin/idk requires sudo"

if ! sudo install -m755 /tmp/idk /usr/local/bin/idk; then
	echo "Error: install failed" >&2
	exit 1
fi

echo "idk $VERSION installed to /usr/local/bin/idk"