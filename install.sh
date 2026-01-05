#!/usr/bin/env bash
set -euo pipefail

REPO_URL=${REPO_URL:-"https://github.com/your-org/zivpn-installer.git"}
INSTALL_DIR=${INSTALL_DIR:-"/opt/zivpn-installer"}

if [[ $EUID -ne 0 ]]; then
  echo "Please run as root (sudo)."
  exit 1
fi

echo "[1/5] Installing dependencies..."
apt-get update -y
apt-get install -y git curl wget tar gzip golang-go

echo "[2/5] Fetching installer repository..."
if [[ -d "$INSTALL_DIR/.git" ]]; then
  git -C "$INSTALL_DIR" pull --ff-only
else
  rm -rf "$INSTALL_DIR"
  git clone "$REPO_URL" "$INSTALL_DIR"
fi

cd "$INSTALL_DIR"

echo "[3/5] Building installer..."
go build -o zivpn-installer .

echo "[4/5] Applying TCP tuning (BBR + FQ + socket buffers)..."
modprobe tcp_bbr || true

cat <<SYSCTL >/etc/sysctl.d/99-zivpn.conf
net.core.default_qdisc = fq
net.ipv4.tcp_congestion_control = bbr
net.ipv4.tcp_fastopen = 3
net.core.rmem_max = 16777216
net.core.wmem_max = 16777216
net.ipv4.tcp_rmem = 4096 87380 16777216
net.ipv4.tcp_wmem = 4096 65536 16777216
SYSCTL

sysctl --system >/dev/null

echo "[5/5] Running installer..."
./zivpn-installer install
