# REAPER Go Extension

A Go-based extension for REAPER using CGO.

This project is an experiment to create REAPER extensions using Go with CGO to bridge to the C++ REAPER API.

## Setup Instructions

1. Clone this repository:

   ```sh
   git clone https://github.com/conormkelly/reaper-golang-extension.git
   cd reaper-golang-extension
   ```

2. Build the extension:

   ```sh
   make install
   ```

## Project Structure

```txt
reaper-go-extension/
├── actions/              # Package for all action handlers
│   ├── fx_prototype.go   # FX prototype-specific functionality
│   └── registry.go       # Central registry for action registration
├── core/                 # Core extension functionality 
│   └── bridge.go         # Core initialization and plugin entry logic
├── reaper/               # REAPER API wrappers
│   ├── actions.go        # Action registration and handling
│   ├── api.go            # Core API initialization
│   ├── console.go        # Console logging functions
│   ├── fx.go             # FX-related functions
│   ├── tracks.go         # Track-related functions
│   └── types.go          # Type definitions
├── main.go               # Main entry point and plugin export
├── c_go_config.go        # CGO configuration
├── reaper_plugin_bridge.c # C bridge code for REAPER API
└── reaper_plugin_bridge.h # C bridge header
```

### Key Components

- **main.go**: The entry point for the extension, exporting the `GoReaperPluginEntry` function which REAPER calls when loading the plugin.

- **reaper/**: Contains Go wrappers for REAPER's C API, making it easier to work with REAPER from Go.
  
- **actions/**: Contains all the custom actions that this extension provides, with a central registry to handle action registration.

- **core/**: Handles core initialization and setup required for the extension to work with REAPER.

- **reaper_plugin_bridge.c/h**: A C bridge that connects Go code to REAPER's C API, handling function pointer conversions and memory management between the two languages.

## Adding New Actions

To add a new action to the extension:

1. Create a new file in the `actions/` directory (e.g., `actions/my_action.go`)
2. Define your action handler and registration function
3. Add your registration function to `actions/registry.go`

Example of a new action file:

```go
package actions

import (
    "fmt"
    "go-reaper/reaper"
)

// RegisterMyAction registers the new action
func RegisterMyAction() error {
    actionID, err := reaper.RegisterMainAction("GO_MY_ACTION", "Go: My New Action")
    if err != nil {
        return fmt.Errorf("failed to register action: %v", err)
    }
    
    reaper.SetActionHandler("GO_MY_ACTION", handleMyAction)
    return nil
}

// handleMyAction handles the action when triggered
func handleMyAction() {
    reaper.ConsoleLog("My action was triggered!")
    // Add your action's functionality here
}
```

Then update `actions/registry.go` to register your new action:

```go
// RegisterAll registers all actions
func RegisterAll() error {
    // Existing registrations...
    
    // Register your new action
    if err := RegisterMyAction(); err != nil {
        return err
    }
    
    return nil
}
```

## Working with REAPER's API

The `reaper/` package provides Go wrappers for common REAPER API functions. If you need to add support for additional REAPER functions:

1. Choose the appropriate file in `reaper/` based on functionality 
2. Add the new function wrapper
3. Update the C bridge in `reaper_plugin_bridge.c/h` if necessary


## Plugin Bridge

### Architecture

The REAPER Go extension is structured in layers:

```txt
┌───────────────────┐
│     REAPER DAW    │  Native C/C++ application
└─────────┬─────────┘
          │ C API calls
┌─────────▼─────────┐
│  Plugin Bridge    │  C code (reaper_plugin_bridge.c/h)
│  (C/Go Interface) │  Translates between C and Go
└─────────┬─────────┘
          │ CGo bindings
┌─────────▼─────────┐
│   REAPER Wrapper  │  Go code (reaper/reaper.go)
│    (Go API)       │  Provides Go-friendly API
└─────────┬─────────┘
          │ Go function calls
┌─────────▼─────────┐
│  Extension Logic  │  Go code (main.go)
│  (Actions & UI)   │  Implements plugin features
└───────────────────┘
```

### Purpose

The plugin bridge (reaper_plugin_bridge.c/h) serves as the critical interface layer between Go and REAPER's native C API. Its primary functions are:

1. **Entry Point Exposure**: Provides the `ReaperPluginEntry` function that REAPER calls when loading the plugin
2. **Function Pointer Handling**: Safely converts and passes C function pointers between REAPER and Go
3. **Memory Safety**: Manages interactions between different memory models (Go garbage collection vs C manual memory)
4. **Type Marshalling**: Converts between Go and C data types safely
5. **Callback Registration**: Registers Go functions as callbacks for REAPER events

### How It Works

#### Initialization

1. REAPER loads the compiled plugin (.dll/.so/.dylib)
2. REAPER calls the exported `ReaperPluginEntry` C function
3. The bridge forwards this call to the Go `GoReaperPluginEntry` function
4. The Go code initializes, storing function pointers for later use
5. The bridge registers Go callbacks with REAPER (for actions, etc.)

#### API Call

When Go code needs to call a REAPER function:

1. Go calls a wrapper function in reaper/reaper.go
2. The wrapper uses CGo to call a bridge function in reaper_plugin_bridge.c
3. The bridge function uses the stored function pointers to call into REAPER
4. Results flow back through the same chain

#### Callback

When REAPER triggers an action:

1. REAPER calls the registered C function in reaper_plugin_bridge.c
2. The C function forwards the call to the exported Go function (e.g., `goHookCommandProc`)
3. The Go function processes the action and returns a result
4. The result is passed back to REAPER

### Key Components

#### GetFunc Mechanism

The `plugin_bridge_get_get_func` and `plugin_bridge_set_get_func` functions manage the core "bootstrap" mechanism that allows dynamically looking up REAPER API functions. This is crucial because:

- REAPER doesn't provide a static library to link against
- Function addresses are only available at runtime
- This mechanism provides access to hundreds of REAPER functions

#### Function Wrappers

The bridge provides typed wrappers for REAPER functions to handle C function pointer casting safely:

```c
void plugin_bridge_call_show_console_msg(void* func_ptr, const char* message) {
    if (!func_ptr || !message) return; // Safety check
    void (*show_console_msg)(const char*) = (void (*)(const char*))func_ptr;
    show_console_msg(message);
}
```

These wrappers make CGo bindings cleaner and safer, as direct function pointer manipulation in Go is complex.

#### Callback Registration

The bridge registers Go functions as REAPER callbacks:

```c
C.plugin_bridge_call_register(registerFuncPtr, cHookCmd, 
                             unsafe.Pointer(C.goHookCommandProc))
```

This allows REAPER to call directly into Go code when actions are triggered.

### Safety Features

The plugin bridge includes several safety features to prevent crashes and ensure robustness:

1. **NULL pointer validation** on all function pointers and critical parameters
2. **Error handling** with appropriate return values for error cases
3. **String buffer safety** ensuring proper NULL-termination
4. **Memory cleanup** with `defer C.free()` after `C.CString()`
5. **Clear API boundaries** between C and Go code

### Thread Safety Considerations

- The bridge assumes REAPER's API is not thread-safe
- All REAPER API calls should happen on the main thread
- Go callbacks triggered by REAPER should not spawn goroutines that call back into REAPER

### Memory Management

Memory allocation spans two worlds:

1. **Go memory**: Managed by Go's garbage collector
2. **C memory**: Manually allocated and freed

Key rules:

- C strings created with `C.CString()` must be freed with `C.free()`
- Pointers to Go memory must not be stored by C code
- C memory must be properly freed to avoid leaks

### Design Patterns

The plugin bridge uses these design patterns:

1. **Adaptor Pattern**: Converts between the C API and Go-friendly interfaces
2. **Facade Pattern**: Simplifies the complex REAPER API into a more manageable Go API
3. **Singleton Pattern**: For global state like the GetFunc pointer
4. **Callback Pattern**: For registration and handling of REAPER events

### Extending the Bridge

When adding new REAPER API support:

1. Find the REAPER function signature in the SDK headers
2. Add a typed wrapper function in reaper_plugin_bridge.h/.c
3. Add the corresponding Go wrapper in reaper/reaper.go
4. Use the wrapper in your extension logic

### Common Pitfalls

1. **Thread Safety**: Never call REAPER functions from background goroutines
2. **Memory Leaks**: Always use `defer C.free()` after `C.CString()`
3. **Type Mismatches**: Ensure C struct definitions match exactly what REAPER expects
4. **String Handling**: Remember C strings are NULL-terminated while Go strings are not
5. **Function Pointer Safety**: Always check function pointers before calling them

## Development Notes

- Study the REAPER SDK headers and SWS extension code for reference

## Acknowledgments

This project wouldn't be possible without:

- Justin Frankel and Cockos Incorporated for creating REAPER and providing the SDK.
- The SWS Extension team for their incredible work, which served as inspiration and a valuable reference for understanding the REAPER API.

## License

The REAPER SDK and WDL are property of Justin Frankel / Cockos Incorporated and are used in accordance with their licensing terms.

The remaining Go code is provided under the [MIT License](LICENSE).
