#!/bin/bash
#
# Docker Installation Script
#
set -e

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/../lib/colors.sh"

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
