// This implements a bridge between REAPER's C/C++ API and Go code.
// It provides safe function pointer handling and type conversion between
// the two language environments.
//
// The bridge pattern allows Go code to call REAPER API functions and
// allows REAPER to call back into Go code through registered callbacks.

#include "reaper_plugin_bridge.h"
#include "reaper_ext_logging.h"

// Implementation of the bridge functions

/**
 * REAPER's GetFunc to retrieve an API function pointer by name
 * This is the fundamental bootstrap mechanism for accessing REAPER's API
 */
void* plugin_bridge_call_get_func(void* get_func_ptr, const char* name) {
    LOG_DEBUG("Called with get_func_ptr=%p, name=%s", get_func_ptr, name ? name : "NULL");
    
    // Verify input pointers aren't NULL
    if (!get_func_ptr || !name) {
        LOG_ERROR("Invalid parameters: get_func_ptr=%p, name=%p", get_func_ptr, name);
        return NULL;
    }
    
    void* (*get_func)(const char*) = (void* (*)(const char*))get_func_ptr;
    void* result = get_func(name);
    
    LOG_DEBUG("Result: %p for function %s", result, name);
    return result;
}

/**
 * REAPER's ShowConsoleMsg function to log messages
 */
void plugin_bridge_call_show_console_msg(void* func_ptr, const char* message) {
    LOG_DEBUG("Called with func_ptr=%p, message=%s", func_ptr, message ? message : "NULL");
    
    // Verify input pointers aren't NULL
    if (!func_ptr || !message) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, message=%p", func_ptr, message);
        return;
    }
    
    void (*show_console_msg)(const char*) = (void (*)(const char*))func_ptr;
    LOG_DEBUG("Calling ShowConsoleMsg with message");
    show_console_msg(message);
    LOG_DEBUG("ShowConsoleMsg call completed");
}

/**
 * REAPER's Register function to register actions, hooks, etc.
 */
int plugin_bridge_call_register(void* register_func_ptr, const char* name, void* info) {
    LOG_DEBUG("Called with register_func_ptr=%p, name=%s, info=%p", 
              register_func_ptr, name ? name : "NULL", info);
    
    // Verify input pointers aren't NULL
    if (!register_func_ptr || !name) {
        LOG_ERROR("Invalid parameters: register_func_ptr=%p, name=%p", register_func_ptr, name);
        return -1; // Error code indicating failure
    }
    
    int (*register_func)(const char*, void*) = (int (*)(const char*, void*))register_func_ptr;
    LOG_DEBUG("Calling Register with name: %s", name);
    int result = register_func(name, info);
    LOG_DEBUG("Register call completed with result: %d", result);
    
    return result;
}

/**
 * REAPER's GetSelectedTrack function
 */
void* plugin_bridge_call_get_selected_track(void* func_ptr, int proj, int seltrackidx) {
    LOG_DEBUG("Called with func_ptr=%p, proj=%d, seltrackidx=%d", func_ptr, proj, seltrackidx);
    
    // Verify input pointer isn't NULL
    if (!func_ptr) {
        LOG_ERROR("Invalid parameter: func_ptr is NULL");
        return NULL;
    }
    
    void* (*get_selected_track)(int, int) = (void* (*)(int, int))func_ptr;
    LOG_DEBUG("Calling GetSelectedTrack with proj=%d, seltrackidx=%d", proj, seltrackidx);
    void* result = get_selected_track(proj, seltrackidx);
    LOG_DEBUG("GetSelectedTrack call completed with result: %p", result);
    
    return result;
}

/**
 * REAPER's TrackFX_GetCount function
 */
int plugin_bridge_call_track_fx_get_count(void* func_ptr, void* track) {
    LOG_DEBUG("Called with func_ptr=%p, track=%p", func_ptr, track);
    
    // Verify input pointers aren't NULL
    if (!func_ptr || !track) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, track=%p", func_ptr, track);
        return 0; // Return 0 FX on error
    }
    
    int (*track_fx_get_count)(void*) = (int (*)(void*))func_ptr;
    LOG_DEBUG("Calling TrackFX_GetCount with track=%p", track);
    int result = track_fx_get_count(track);
    LOG_DEBUG("TrackFX_GetCount call completed with result: %d", result);
    
    return result;
}

/**
 * REAPER's TrackFX_GetFXName function
 */
void plugin_bridge_call_track_fx_get_name(void* func_ptr, void* track, int fx_idx, char* buf, int buf_size) {
    LOG_DEBUG("Called with func_ptr=%p, track=%p, fx_idx=%d, buf=%p, buf_size=%d", 
              func_ptr, track, fx_idx, buf, buf_size);
    
    // Verify input pointers aren't NULL
    if (!func_ptr || !track || !buf || buf_size <= 0) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, track=%p, buf=%p, buf_size=%d", 
                  func_ptr, track, buf, buf_size);
        // If buffer is valid, make it an empty string for safety
        if (buf && buf_size > 0) {
            buf[0] = '\0';
            LOG_DEBUG("Buffer set to empty string for safety");
        }
        return;
    }
    
    void (*track_fx_get_name)(void*, int, char*, int) = 
        (void (*)(void*, int, char*, int))func_ptr;
    LOG_DEBUG("Calling TrackFX_GetFXName with track=%p, fx_idx=%d", track, fx_idx);
    track_fx_get_name(track, fx_idx, buf, buf_size);
    LOG_DEBUG("TrackFX_GetFXName call completed with result: %s", buf);
}

/**
 * REAPER's TrackFX_GetNumParams function
 */
int plugin_bridge_call_track_fx_get_param_count(void* func_ptr, void* track, int fx_idx) {
    LOG_DEBUG("Called with func_ptr=%p, track=%p, fx_idx=%d", func_ptr, track, fx_idx);
    
    // Verify input pointers aren't NULL
    if (!func_ptr || !track) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, track=%p", func_ptr, track);
        return 0; // Return 0 parameters on error
    }
    
    int (*track_fx_get_param_count)(void*, int) = (int (*)(void*, int))func_ptr;
    LOG_DEBUG("Calling TrackFX_GetNumParams with track=%p, fx_idx=%d", track, fx_idx);
    int result = track_fx_get_param_count(track, fx_idx);
    LOG_DEBUG("TrackFX_GetNumParams call completed with result: %d", result);
    
    return result;
}

/**
 * REAPER's TrackFX_GetParamName function
 */
void plugin_bridge_call_track_fx_get_param_name(void* func_ptr, void* track, int fx_idx, int param_idx, char* buf, int buf_size) {
    LOG_DEBUG("Called with func_ptr=%p, track=%p, fx_idx=%d, param_idx=%d, buf=%p, buf_size=%d", 
              func_ptr, track, fx_idx, param_idx, buf, buf_size);
    
    // Verify input pointers aren't NULL
    if (!func_ptr || !track || !buf || buf_size <= 0) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, track=%p, buf=%p, buf_size=%d", 
                  func_ptr, track, buf, buf_size);
        // If buffer is valid, make it an empty string for safety
        if (buf && buf_size > 0) {
            buf[0] = '\0';
            LOG_DEBUG("Buffer set to empty string for safety");
        }
        return;
    }
    
    void (*track_fx_get_param_name)(void*, int, int, char*, int) = 
        (void (*)(void*, int, int, char*, int))func_ptr;
    LOG_DEBUG("Calling TrackFX_GetParamName with track=%p, fx_idx=%d, param_idx=%d", 
              track, fx_idx, param_idx);
    track_fx_get_param_name(track, fx_idx, param_idx, buf, buf_size);
    LOG_DEBUG("TrackFX_GetParamName call completed with result: %s", buf);
}

/**
 * REAPER's TrackFX_GetParam function
 */
double plugin_bridge_call_track_fx_get_param(void* func_ptr, void* track, int fx_idx, int param_idx, double* minval, double* maxval) {
    LOG_DEBUG("Called with func_ptr=%p, track=%p, fx_idx=%d, param_idx=%d, minval=%p, maxval=%p", 
              func_ptr, track, fx_idx, param_idx, minval, maxval);
    
    // Verify input pointers aren't NULL
    if (!func_ptr || !track) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, track=%p", func_ptr, track);
        return 0.0; // Return 0.0 on error
    }
    
    double (*track_fx_get_param)(void*, int, int, double*, double*) = 
        (double (*)(void*, int, int, double*, double*))func_ptr;
    LOG_DEBUG("Calling TrackFX_GetParam with track=%p, fx_idx=%d, param_idx=%d", 
              track, fx_idx, param_idx);
    double result = track_fx_get_param(track, fx_idx, param_idx, minval, maxval);
    
    // Log min/max values if provided
    if (minval && maxval) {
        LOG_DEBUG("TrackFX_GetParam call completed with result: %f, min=%f, max=%f", 
                  result, *minval, *maxval);
    } else {
        LOG_DEBUG("TrackFX_GetParam call completed with result: %f", result);
    }
    
    return result;
}

/**
 * REAPER's TrackFX_GetFormattedParamValue function
 */
void plugin_bridge_call_track_fx_get_param_formatted(void* func_ptr, void* track, int fx_idx, int param_idx, char* buf, int buf_size) {
    LOG_DEBUG("Called with func_ptr=%p, track=%p, fx_idx=%d, param_idx=%d, buf=%p, buf_size=%d", 
              func_ptr, track, fx_idx, param_idx, buf, buf_size);
    
    // Verify input pointers aren't NULL
    if (!func_ptr || !track || !buf || buf_size <= 0) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, track=%p, buf=%p, buf_size=%d", 
                  func_ptr, track, buf, buf_size);
        // If buffer is valid, make it an empty string for safety
        if (buf && buf_size > 0) {
            buf[0] = '\0';
            LOG_DEBUG("Buffer set to empty string for safety");
        }
        return;
    }
    
    void (*track_fx_get_param_formatted)(void*, int, int, char*, int) = 
        (void (*)(void*, int, int, char*, int))func_ptr;
    LOG_DEBUG("Calling TrackFX_GetFormattedParamValue with track=%p, fx_idx=%d, param_idx=%d", 
              track, fx_idx, param_idx);
    track_fx_get_param_formatted(track, fx_idx, param_idx, buf, buf_size);
    LOG_DEBUG("TrackFX_GetFormattedParamValue call completed with result: %s", buf);
}

/**
 * REAPER's TrackFX_SetParam function
 */
bool plugin_bridge_call_track_fx_set_param(void* func_ptr, void* track, int fx_idx, int param_idx, double val) {
    LOG_DEBUG("Called with func_ptr=%p, track=%p, fx_idx=%d, param_idx=%d, val=%f", 
              func_ptr, track, fx_idx, param_idx, val);
    
    // Verify input pointers aren't NULL
    if (!func_ptr || !track) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, track=%p", func_ptr, track);
        return false; // Return failure on error
    }
    
    bool (*track_fx_set_param)(void*, int, int, double) = 
        (bool (*)(void*, int, int, double))func_ptr;
    LOG_DEBUG("Calling TrackFX_SetParam with track=%p, fx_idx=%d, param_idx=%d, val=%f", 
              track, fx_idx, param_idx, val);
    bool result = track_fx_set_param(track, fx_idx, param_idx, val);
    LOG_DEBUG("TrackFX_SetParam call completed with result: %d", result);
    
    return result;
}

// Get track information value
double plugin_bridge_call_track_get_info_value(void* func_ptr, void* track, const char* param) {
    LOG_DEBUG("Called with func_ptr=%p, track=%p, param=%s", 
              func_ptr, track, param ? param : "NULL");
    
    if (!func_ptr || !track || !param) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, track=%p, param=%p", 
                  func_ptr, track, param);
        return 0.0;
    }
    
    double (*get_track_info)(void*, const char*) = (double (*)(void*, const char*))func_ptr;
    LOG_DEBUG("Calling GetMediaTrackInfo_Value with track=%p, param=%s", track, param);
    double result = get_track_info(track, param);
    LOG_DEBUG("GetMediaTrackInfo_Value call completed with result: %f", result);
    
    return result;
}

// Get track name
bool plugin_bridge_call_get_track_name(void* func_ptr, void* track, char* buf, int buf_size, int* flags) {
    LOG_DEBUG("Called with func_ptr=%p, track=%p, buf=%p, buf_size=%d, flags=%p", 
              func_ptr, track, buf, buf_size, flags);
    
    if (!func_ptr || !track || !buf || buf_size <= 0) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, track=%p, buf=%p, buf_size=%d", 
                  func_ptr, track, buf, buf_size);
        // If buffer is valid, make it an empty string for safety
        if (buf && buf_size > 0) {
            buf[0] = '\0';
            LOG_DEBUG("Buffer set to empty string for safety");
        }
        return false;
    }
    
    bool (*get_track_name)(void*, char*, int, int*) = 
        (bool (*)(void*, char*, int, int*))func_ptr;
    
    LOG_DEBUG("Calling GetTrackName with track=%p", track);
    bool result = get_track_name(track, buf, buf_size, flags);
    
    // Log flags if provided
    if (flags) {
        LOG_DEBUG("GetTrackName call completed with result: %d, name=%s, flags=%d", 
                  result, buf, *flags);
    } else {
        LOG_DEBUG("GetTrackName call completed with result: %d, name=%s", 
                  result, buf);
    }
    
    return result;
}

/**
 * REAPER's GetUserInputs function for creating simple form dialogs
 */
bool plugin_bridge_call_get_user_inputs(void* func_ptr, const char* title, int num_inputs, 
    const char* captions, char* values, int values_sz) {
    LOG_DEBUG("Called with func_ptr=%p, title=%s, num_inputs=%d, captions=%s, values_sz=%d", 
              func_ptr, title ? title : "NULL", num_inputs, 
              captions ? captions : "NULL", values_sz);
    
    // Verify input pointers aren't NULL
    if (!func_ptr || !title || !captions || !values || values_sz <= 0) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, title=%p, captions=%p, values=%p, values_sz=%d", 
                  func_ptr, title, captions, values, values_sz);
        return false;
    }

    bool (*get_user_inputs)(const char*, int, const char*, char*, int) =
        (bool (*)(const char*, int, const char*, char*, int))func_ptr;

    LOG_DEBUG("Calling GetUserInputs with title=%s, num_inputs=%d", title, num_inputs);
    bool result = get_user_inputs(title, num_inputs, captions, values, values_sz);
    LOG_DEBUG("GetUserInputs call completed with result: %d, values=%s", result, values);
    
    return result;
}

/**
* REAPER's ShowMessageBox function for standard message boxes
*/
int plugin_bridge_call_show_message_box(void* func_ptr, const char* text, const char* title, int type) {
    LOG_DEBUG("Called with func_ptr=%p, text=%s, title=%s, type=%d", 
              func_ptr, text ? text : "NULL", title ? title : "NULL", type);
    
    // Verify input pointers aren't NULL
    if (!func_ptr || !text || !title) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, text=%p, title=%p", 
                  func_ptr, text, title);
        return 0; // Return IDOK (1) by default on error
    }

    int (*show_message_box)(const char*, const char*, int) =
        (int (*)(const char*, const char*, int))func_ptr;

    LOG_DEBUG("Calling ShowMessageBox with text='%s', title='%s', type=%d", text, title, type);
    int result = show_message_box(text, title, type);
    LOG_DEBUG("ShowMessageBox call completed with result: %d", result);
    
    return result;
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
    LOG_INFO("Setting global GetFunc pointer to %p", get_func_ptr);
    
    // Only store if it's a valid pointer
    if (get_func_ptr) {
        s_GetFunc = get_func_ptr;
        LOG_INFO("Global GetFunc pointer set successfully");
    } else {
        LOG_ERROR("Attempted to set NULL GetFunc pointer");
    }
}

/**
 * Returns the stored GetFunc pointer
 * This is the bootstrap function used to access all other REAPER functions
 */
void* plugin_bridge_get_get_func() {
    LOG_DEBUG("Retrieving global GetFunc pointer: %p", s_GetFunc);
    return s_GetFunc;
}

/**
 * Main entry point called by REAPER when loading the plugin
 * This function forwards the call to the Go entry point via CGo
 */
REAPER_PLUGIN_DLL_EXPORT int ReaperPluginEntry(HINSTANCE hInstance, reaper_plugin_info_t* rec) {
    LOG_INFO("REAPER plugin entry called with hInstance=%p, rec=%p", hInstance, rec);
    
    if (!rec) {
        LOG_INFO("rec is NULL, plugin is being unloaded");
        return 0;
    }
    
    // Log REAPER API version
    if (rec) {
        LOG_INFO("REAPER API version: 0x%X", rec->caller_version);
    }
    
    // Forward to Go entry point, using void* to simplify CGo binding
    LOG_INFO("Forwarding to Go entry point");
    int result = GoReaperPluginEntry((void*)hInstance, (void*)rec);
    LOG_INFO("Go entry point returned: %d", result);
    
    return result;
}
