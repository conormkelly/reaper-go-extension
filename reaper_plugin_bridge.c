// This implements a bridge between REAPER's C/C++ API and Go code.
// It provides safe function pointer handling and type conversion between
// the two language environments.
//
// The bridge pattern allows Go code to call REAPER API functions and
// allows REAPER to call back into Go code through registered callbacks.

#include "reaper_plugin_bridge.h"

// Implementation of the bridge functions

/**
 * REAPER's GetFunc to retrieve an API function pointer by name
 * This is the fundamental bootstrap mechanism for accessing REAPER's API
 */
void* plugin_bridge_call_get_func(void* get_func_ptr, const char* name) {
    // Verify input pointers aren't NULL
    if (!get_func_ptr || !name) {
        return NULL;
    }
    
    void* (*get_func)(const char*) = (void* (*)(const char*))get_func_ptr;
    return get_func(name);
}

/**
 * REAPER's ShowConsoleMsg function to log messages
 */
void plugin_bridge_call_show_console_msg(void* func_ptr, const char* message) {
    // Verify input pointers aren't NULL
    if (!func_ptr || !message) {
        return;
    }
    
    void (*show_console_msg)(const char*) = (void (*)(const char*))func_ptr;
    show_console_msg(message);
}

/**
 * REAPER's Register function to register actions, hooks, etc.
 */
int plugin_bridge_call_register(void* register_func_ptr, const char* name, void* info) {
    // Verify input pointers aren't NULL
    if (!register_func_ptr || !name) {
        return -1; // Error code indicating failure
    }
    
    int (*register_func)(const char*, void*) = (int (*)(const char*, void*))register_func_ptr;
    return register_func(name, info);
}

/**
 * REAPER's GetSelectedTrack function
 */
void* plugin_bridge_call_get_selected_track(void* func_ptr, int proj, int seltrackidx) {
    // Verify input pointer isn't NULL
    if (!func_ptr) {
        return NULL;
    }
    
    void* (*get_selected_track)(int, int) = (void* (*)(int, int))func_ptr;
    return get_selected_track(proj, seltrackidx);
}

/**
 * REAPER's TrackFX_GetCount function
 */
int plugin_bridge_call_track_fx_get_count(void* func_ptr, void* track) {
    // Verify input pointers aren't NULL
    if (!func_ptr || !track) {
        return 0; // Return 0 FX on error
    }
    
    int (*track_fx_get_count)(void*) = (int (*)(void*))func_ptr;
    return track_fx_get_count(track);
}

/**
 * REAPER's TrackFX_GetFXName function
 */
void plugin_bridge_call_track_fx_get_name(void* func_ptr, void* track, int fx_idx, char* buf, int buf_size) {
    // Verify input pointers aren't NULL
    if (!func_ptr || !track || !buf || buf_size <= 0) {
        // If buffer is valid, make it an empty string for safety
        if (buf && buf_size > 0) {
            buf[0] = '\0';
        }
        return;
    }
    
    void (*track_fx_get_name)(void*, int, char*, int) = 
        (void (*)(void*, int, char*, int))func_ptr;
    track_fx_get_name(track, fx_idx, buf, buf_size);
}

/**
 * REAPER's TrackFX_GetNumParams function
 */
int plugin_bridge_call_track_fx_get_param_count(void* func_ptr, void* track, int fx_idx) {
    // Verify input pointers aren't NULL
    if (!func_ptr || !track) {
        return 0; // Return 0 parameters on error
    }
    
    int (*track_fx_get_param_count)(void*, int) = (int (*)(void*, int))func_ptr;
    return track_fx_get_param_count(track, fx_idx);
}

/**
 * REAPER's TrackFX_GetParamName function
 */
void plugin_bridge_call_track_fx_get_param_name(void* func_ptr, void* track, int fx_idx, int param_idx, char* buf, int buf_size) {
    // Verify input pointers aren't NULL
    if (!func_ptr || !track || !buf || buf_size <= 0) {
        // If buffer is valid, make it an empty string for safety
        if (buf && buf_size > 0) {
            buf[0] = '\0';
        }
        return;
    }
    
    void (*track_fx_get_param_name)(void*, int, int, char*, int) = 
        (void (*)(void*, int, int, char*, int))func_ptr;
    track_fx_get_param_name(track, fx_idx, param_idx, buf, buf_size);
}

/**
 * REAPER's TrackFX_GetParam function
 */
double plugin_bridge_call_track_fx_get_param(void* func_ptr, void* track, int fx_idx, int param_idx, double* minval, double* maxval) {
    // Verify input pointers aren't NULL
    if (!func_ptr || !track) {
        return 0.0; // Return 0.0 on error
    }
    
    double (*track_fx_get_param)(void*, int, int, double*, double*) = 
        (double (*)(void*, int, int, double*, double*))func_ptr;
    return track_fx_get_param(track, fx_idx, param_idx, minval, maxval);
}

/**
 * REAPER's TrackFX_GetFormattedParamValue function
 */
void plugin_bridge_call_track_fx_get_param_formatted(void* func_ptr, void* track, int fx_idx, int param_idx, char* buf, int buf_size) {
    // Verify input pointers aren't NULL
    if (!func_ptr || !track || !buf || buf_size <= 0) {
        // If buffer is valid, make it an empty string for safety
        if (buf && buf_size > 0) {
            buf[0] = '\0';
        }
        return;
    }
    
    void (*track_fx_get_param_formatted)(void*, int, int, char*, int) = 
        (void (*)(void*, int, int, char*, int))func_ptr;
    track_fx_get_param_formatted(track, fx_idx, param_idx, buf, buf_size);
}

/**
 * REAPER's TrackFX_SetParam function
 */
bool plugin_bridge_call_track_fx_set_param(void* func_ptr, void* track, int fx_idx, int param_idx, double val) {
    // Verify input pointers aren't NULL
    if (!func_ptr || !track) {
        return false; // Return failure on error
    }
    
    bool (*track_fx_set_param)(void*, int, int, double) = 
        (bool (*)(void*, int, int, double))func_ptr;
    return track_fx_set_param(track, fx_idx, param_idx, val);
}

// Global storage for REAPER's GetFunc pointer
// This is a central lookup mechanism for all REAPER API functions
// It's accessed from multiple functions but is set only once during initialization
static void* s_GetFunc = NULL;

/**
 * Sets the global GetFunc pointer that's used to lookup REAPER functions
 * Called once during plugin initialization
 */
void plugin_bridge_set_get_func(void* get_func_ptr) {
    // Only store if it's a valid pointer
    if (get_func_ptr) {
        s_GetFunc = get_func_ptr;
    }
}

/**
 * Returns the stored GetFunc pointer
 * This is the bootstrap function used to access all other REAPER functions
 */
void* plugin_bridge_get_get_func() {
    return s_GetFunc;
}

/**
 * Main entry point called by REAPER when loading the plugin
 * This function forwards the call to the Go entry point via CGo
 */
REAPER_PLUGIN_DLL_EXPORT int ReaperPluginEntry(HINSTANCE hInstance, reaper_plugin_info_t* rec) {
    // Forward to Go entry point, using void* to simplify CGo binding
    return GoReaperPluginEntry((void*)hInstance, (void*)rec);
}
