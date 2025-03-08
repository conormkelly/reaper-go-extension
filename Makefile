# Makefile for building REAPER Go extension

GOOS=$(shell go env GOOS)
SDK_DIR=./sdk
SRC_DIR=./src
BUILD_DIR=./build

# Set extension based on platform
ifeq ($(GOOS),windows)
  EXT=.dll
else ifeq ($(GOOS),darwin)
  EXT=.dylib
else
  EXT=.so
endif

# Make sure build directory exists
$(shell mkdir -p $(BUILD_DIR))

all: $(BUILD_DIR)/reaper_hello_go$(EXT)

# First compile the Go code to a temporary archive
$(BUILD_DIR)/libgo_reaper.a: $(SRC_DIR)/main.go
	cd $(SRC_DIR) && go build -buildmode=c-archive -o ../$(BUILD_DIR)/libgo_reaper.a main.go

# Compile the bridge code
$(BUILD_DIR)/reaper_plugin_bridge.o: $(SRC_DIR)/reaper_plugin_bridge.c $(SRC_DIR)/reaper_plugin_bridge.h
	gcc -c -I$(SDK_DIR) -I$(SRC_DIR) $(SRC_DIR)/reaper_plugin_bridge.c -o $(BUILD_DIR)/reaper_plugin_bridge.o

# Link everything together
$(BUILD_DIR)/reaper_hello_go$(EXT): $(BUILD_DIR)/libgo_reaper.a $(BUILD_DIR)/reaper_plugin_bridge.o
	gcc -shared -o $(BUILD_DIR)/reaper_hello_go$(EXT) $(BUILD_DIR)/reaper_plugin_bridge.o $(BUILD_DIR)/libgo_reaper.a -framework CoreFoundation -lpthread

# Install the plugin to REAPER's plugin directory (macOS path shown)
install: $(BUILD_DIR)/reaper_hello_go$(EXT)
	cp $(BUILD_DIR)/reaper_hello_go$(EXT) "$(HOME)/Library/Application Support/REAPER/UserPlugins/"

clean:
	rm -rf $(BUILD_DIR)/*

.PHONY: all clean install
