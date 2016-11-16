#!/usr/bin/env bash
set -eux

CACHE_DIR="$1"; shift
mkdir -p "$CACHE_DIR"
cd "$CACHE_DIR"

sudo apt-get remove --purge golang
sudo rm -rf '/usr/local/go/'

VERSION='1.7.3'
GO="go${VERSION}.linux-amd64.tar.gz"

wget --no-clobber "https://storage.googleapis.com/golang/$GO"
sudo tar -xzf "$GO" -C '/usr/local'
