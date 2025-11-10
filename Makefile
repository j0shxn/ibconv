# --- Variables ---

# The output directory for binaries
BIN_DIR = ./bin

# The name of the binary
BINARY_NAME = ibconv

# Full paths for the executables
# Native binary (Linux)
BINARY_NATIVE = $(BIN_DIR)/$(BINARY_NAME)
# Windows binary
BINARY_WIN = $(BIN_DIR)/$(BINARY_NAME).exe

# Go commands
GO = go
GOBUILD = $(GO) build

# Find all .go files in the current directory
GO_FILES = $(wildcard *.go)

# --- Targets ---

# .PHONY ensures these targets run even if files with the same name exist
.PHONY: all build build-linux windows install clean tidy

# 'all' is the default target. `make` will run this if you type just 'make'.
# It now points to the 'build' target.
all: build

# 'build' is a generic target that builds for the current (native) OS.
# For you, this is Linux.
build: $(GO_FILES) go.mod
	@echo "Building for native OS (Linux)..."
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) -o $(BINARY_NATIVE) .
	@echo "Build complete: $(BINARY_NATIVE)"

# 'windows' target explicitly cross-compiles for Windows.
windows: $(GO_FILES) go.mod
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WIN) .
	@echo "Build complete: $(BINARY_WIN)"

# 'install' target installs the NATIVE binary to /usr/local/bin.
# It depends on 'build' to ensure the binary exists first.
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@echo "This requires administrator privileges."
	sudo install -m 0755 $(BINARY_NATIVE) /usr/local/bin/$(BINARY_NAME)
	@echo "Install complete."
	@echo "Run 'ibconv -h' to test."

# 'clean' target removes the entire bin directory.
clean:
	@echo "Cleaning up..."
	rm -rf $(BIN_DIR)

# 'tidy' target updates go.mod
tidy:
	@echo "Running go mod tidy..."
	$(GO) mod tidy
