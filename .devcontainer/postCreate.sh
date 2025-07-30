#!/bin/bash

set -e -o pipefail

echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | sudo tee /etc/apt/sources.list.d/goreleaser.list
sudo --preserve-env=http_proxy,https_proxy apt update
sudo --preserve-env=http_proxy,https_proxy apt install -y goreleaser
