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

# Find all C source files in src/c and subdirectories
C_SRC_FILES := $(shell find $(SRC_DIR)/c -name "*.c")
C_OBJ_FILES := $(patsubst $(SRC_DIR)/%.c,$(BUILD_DIR)/%.o,$(C_SRC_FILES))

# For macOS, also include Objective-C files from UI platform
ifeq ($(GOOS),darwin)
  OBJC_SRC_FILES := $(shell find $(SRC_DIR)/ui/platform/macos -name "*.m")
  OBJC_OBJ_FILES := $(patsubst $(SRC_DIR)/%.m,$(BUILD_DIR)/%.o,$(OBJC_SRC_FILES))
  
  # Also include krbridge.m from actions directory
  ACTIONS_OBJC_SRC := $(SRC_DIR)/actions/krbridge.m
  ACTIONS_OBJC_OBJ := $(BUILD_DIR)/actions/krbridge.o
  
  # Combined object files
  ALL_OBJ_FILES := $(C_OBJ_FILES) $(OBJC_OBJ_FILES) $(ACTIONS_OBJC_OBJ)
else
  ALL_OBJ_FILES := $(C_OBJ_FILES)
endif

# Make sure build directory exists
$(shell mkdir -p $(BUILD_DIR))
$(shell mkdir -p $(BUILD_DIR)/c/api)
$(shell mkdir -p $(BUILD_DIR)/ui/platform/macos)
$(shell mkdir -p $(BUILD_DIR)/actions)

all: $(BUILD_DIR)/reaper_hello_go$(EXT)

# First compile the Go code to a temporary archive
# This now depends on all Go source files
$(BUILD_DIR)/libgo_reaper.a: $(GO_SRC_FILES)
	go build -buildmode=c-archive -o $(BUILD_DIR)/libgo_reaper.a $(CMD_DIR)/main.go

# Compile individual C files
$(BUILD_DIR)/%.o: $(SRC_DIR)/%.c
	mkdir -p $(dir $@)
	gcc -c -I$(SDK_DIR) -I$(SRC_DIR) $< -o $@

# Compile Objective-C files (for macOS)
$(BUILD_DIR)/%.o: $(SRC_DIR)/%.m
	mkdir -p $(dir $@)
	gcc -c -x objective-c -I$(SDK_DIR) -I$(SRC_DIR) $< -o $@

# Link everything together - macOS specific version
ifeq ($(GOOS),darwin)
$(BUILD_DIR)/reaper_hello_go$(EXT): $(BUILD_DIR)/libgo_reaper.a $(ALL_OBJ_FILES)
	gcc -shared -o $(BUILD_DIR)/reaper_hello_go$(EXT) $(ALL_OBJ_FILES) $(BUILD_DIR)/libgo_reaper.a $(MACOS_LDFLAGS) -lpthread
else
# Link everything together - non-macOS version
$(BUILD_DIR)/reaper_hello_go$(EXT): $(BUILD_DIR)/libgo_reaper.a $(ALL_OBJ_FILES)
	gcc -shared -o $(BUILD_DIR)/reaper_hello_go$(EXT) $(ALL_OBJ_FILES) $(BUILD_DIR)/libgo_reaper.a -lpthread
endif

# Install the plugin to REAPER's plugin directory
install: $(BUILD_DIR)/reaper_hello_go$(EXT)
	cp $(BUILD_DIR)/reaper_hello_go$(EXT) $(INSTALL_PATH)

clean:
	rm -rf $(BUILD_DIR)/*

# Check logs for errors and important messages
check:
	@echo "Checking log file reaper-ext.log for issues..."
	@echo ""
	
	@echo "=== Checking for ERROR messages ==="
	@rg "ERROR" reaper-ext.log || echo "No errors found."
	@echo ""
	
	@echo "=== Checking for WARNING messages ==="
	@rg "WARNING" reaper-ext.log || echo "No warnings found."
	@echo ""
	
	@echo "=== Checking for successful plugin loading ==="
	@rg "plugin loaded successfully" reaper-ext.log || echo "Plugin load message not found!"
	@echo ""
	
	@echo "=== Checking for function call failures ==="
	@rg "failed to" reaper-ext.log || echo "No function call failures found."
	@echo ""
	
	@echo "=== Examining bridge initialization ==="
	@rg "REAPER plugin entry" reaper-ext.log || echo "Bridge initialization not found!"
	@echo ""
	
	@echo "=== Checking for action issues ==="
	@rg "action triggered" reaper-ext.log || echo "No actions triggered yet."
	@echo ""
	
	@echo "=== Verifying API function lookups ==="
	@rg "Failed to get .* function pointer" reaper-ext.log || echo "All API functions found successfully."
	@echo ""
	
	@echo "Log check complete!"

.PHONY: all clean install check
