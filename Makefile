# Variables
GOCMD = go
GOBUILD = $(GOCMD) build -ldflags="-s -w"
GO_FILE = $(PWD)/cmd/main.go
BINARY_FILE = k0pt
PKG_DIR = k0pt_amd64
BIN_DIR = /usr/local/bin

all: remove-binary install build-package

install:
	@echo "Generating binary in the pkg folder"
	$(GOBUILD) -o $(BINARY_FILE) $(GO_FILE)
	sudo cp $(BINARY_FILE) $(BIN_DIR)
	mv $(BINARY_FILE) ./$(PKG_DIR)$(BIN_DIR)/

build-package:
	@echo "Removing existing package..."
	rm -f *.deb
	@echo "Building Debian package..."
	dpkg-deb --build $(PKG_DIR)

remove-binary: 
	@echo "Removing existing binary..."
	sudo rm -f $(BINARY_FILE)
	sudo rm -f ./$(PKG_DIR)$(BIN_DIR)/$(BINARY_FILE)
	sudo rm -f $(BIN_DIR)/$(BINARY_FILE)

.PHONY: all remove-binary install build-package