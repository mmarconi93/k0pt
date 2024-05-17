#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status

# Variables
REPO_URL="https://k0pt.akstest.tech/k0pt"
KEYRING_DIR="/etc/apt/keyrings"
KEYRING_FILE="$KEYRING_DIR/k0pt-stable-amd64.gpg"
SOURCE_LIST_FILE="/etc/apt/sources.list.d/k0pt-repo.list"

# Create the keyring directory
sudo mkdir -p -m 755 "$KEYRING_DIR"

# Download the GPG key
wget -qO- "$REPO_URL/k0pt-stable-amd64.gpg" | sudo tee "$KEYRING_FILE" >/dev/null
sudo chmod go+r "$KEYRING_FILE"

# Add the repository to the sources list
echo "deb [signed-by=$KEYRING_FILE] $REPO_URL/ stable main" | sudo tee "$SOURCE_LIST_FILE" >/dev/null

# Update the package list and install the k0pt package
sudo apt update
sudo apt install -y k0pt

echo "k0pt installed successfully."