#!/bin/bash
#
# Docker Installation Script
#
set -e

GREEN='\033[0;32m'
NC='\033[0m'

print_success() { echo -e "${GREEN}[OK]${NC} $1"; }

if command -v docker &>/dev/null; then
    print_success "Docker is already installed"
    exit 0
fi

echo "Installing Docker..."
sudo apt update
sudo apt install -y docker.io
sudo systemctl start docker
sudo systemctl enable docker

print_success "Docker installed"
