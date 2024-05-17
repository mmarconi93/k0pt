#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status

# Variables
DIST_DIR="dists/stable/main/binary-amd64"
PACKAGE_NAME="k0pt-amd64"

# Change into the directory where the .deb files are located
cd "$DIST_DIR"

# Generate Packages files
apt-ftparchive packages . > Packages
bzip2 -c9 Packages > Packages.bz2
xz -c9 Packages > Packages.xz
xz -5fkev --format=lzma Packages > Packages.lzma
lz4 -c9 Packages > Packages.lz4
gzip -c9 Packages > Packages.gz
zstd -c19 Packages > Packages.zst

# Generate Contents files
apt-ftparchive contents . > Contents-"$PACKAGE_NAME"
bzip2 -c9 Contents-"$PACKAGE_NAME" > Contents-"$PACKAGE_NAME".bz2
xz -c9 Contents-"$PACKAGE_NAME" > Contents-"$PACKAGE_NAME".xz
xz -5fkev --format=lzma Contents-"$PACKAGE_NAME" > Contents-"$PACKAGE_NAME".lzma
lz4 -c9 Contents-"$PACKAGE_NAME" > Contents-"$PACKAGE_NAME".lz4
gzip -c9 Contents-"$PACKAGE_NAME" > Contents-"$PACKAGE_NAME".gz
zstd -c19 Contents-"$PACKAGE_NAME" > Contents-"$PACKAGE_NAME".zst

# Change back to the root of the repository
cd ../../..

# Generate Release files
grep -E "Origin:|Label:|Suite:|Version:|Codename:|Architectures:|Components:|Description:" Release > Base
apt-ftparchive release . > Release
cat Base Release > out && mv out Release

# Sign the Release file
GPG_KEY_ID="9D0FA86CF53F15831B6D8CB69633467FEBF46E4E"  # Replace with your GPG key ID
gpg -abs -u "$GPG_KEY_ID" -o Release.gpg Release

echo "Packages and Release files generated and signed successfully."