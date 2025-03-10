# Makefile for building REAPER Go extension

GOOS=$(shell go env GOOS)
SDK_DIR=./sdk
SRC_DIR=./src
CMD_DIR=./cmd/reaper-ext
BUILD_DIR=./build

# Set extension based on platform
ifeq ($(GOOS),windows)
  EXT=.dll
  INSTALL_PATH="$(APPDATA)/REAPER/UserPlugins/"
else ifeq ($(GOOS),darwin)
  EXT=.dylib
  INSTALL_PATH="$(HOME)/Library/Application Support/REAPER/UserPlugins/"
  # Add macOS specific flags
    MACOS_LDFLAGS=-framework CoreFoundation \
                -framework Security \
                -framework Cocoa
else
  EXT=.so
  INSTALL_PATH="$(HOME)/.config/REAPER/UserPlugins/"
endif

# Find all Go source files
GO_SRC_FILES := $(shell find $(SRC_DIR) $(CMD_DIR) -name "*.go")

# Make sure build directory exists
$(shell mkdir -p $(BUILD_DIR))

all: $(BUILD_DIR)/reaper_hello_go$(EXT)

# First compile the Go code to a temporary archive
# This now depends on all Go source files
$(BUILD_DIR)/libgo_reaper.a: $(GO_SRC_FILES)
	go build -buildmode=c-archive -o $(BUILD_DIR)/libgo_reaper.a $(CMD_DIR)/main.go

# Compile the bridge code
$(BUILD_DIR)/bridge.o: $(SRC_DIR)/c/bridge.c $(SRC_DIR)/c/bridge.h
	gcc -c -I$(SDK_DIR) -I$(SRC_DIR) $(SRC_DIR)/c/bridge.c -o $(BUILD_DIR)/bridge.o

# Compile the logging code
$(BUILD_DIR)/logging.o: $(SRC_DIR)/c/logging.c $(SRC_DIR)/c/logging.h
	gcc -c -I$(SDK_DIR) -I$(SRC_DIR) $(SRC_DIR)/c/logging.c -o $(BUILD_DIR)/logging.o

# Compile the keyring bridge code (for macOS only)
ifeq ($(GOOS),darwin)
$(BUILD_DIR)/krbridge.o: $(SRC_DIR)/actions/krbridge.m $(SRC_DIR)/actions/krbridge.h
	gcc -c -x objective-c -I$(SDK_DIR) -I$(SRC_DIR) $(SRC_DIR)/actions/krbridge.m -o $(BUILD_DIR)/krbridge.o
endif

# Link everything together
ifeq ($(GOOS),darwin)
$(BUILD_DIR)/reaper_hello_go$(EXT): $(BUILD_DIR)/libgo_reaper.a $(BUILD_DIR)/bridge.o $(BUILD_DIR)/logging.o $(BUILD_DIR)/krbridge.o
	gcc -shared -o $(BUILD_DIR)/reaper_hello_go$(EXT) $(BUILD_DIR)/bridge.o $(BUILD_DIR)/logging.o $(BUILD_DIR)/krbridge.o $(BUILD_DIR)/libgo_reaper.a $(MACOS_LDFLAGS) -lpthread
else
$(BUILD_DIR)/reaper_hello_go$(EXT): $(BUILD_DIR)/libgo_reaper.a $(BUILD_DIR)/bridge.o $(BUILD_DIR)/logging.o
	gcc -shared -o $(BUILD_DIR)/reaper_hello_go$(EXT) $(BUILD_DIR)/bridge.o $(BUILD_DIR)/logging.o $(BUILD_DIR)/libgo_reaper.a -lpthread
endif

# Install the plugin to REAPER's plugin directory
install: $(BUILD_DIR)/reaper_hello_go$(EXT)
	cp $(BUILD_DIR)/reaper_hello_go$(EXT) $(INSTALL_PATH)

clean:
	rm -rf $(BUILD_DIR)/*

.PHONY: all clean install
