/**
 * Implementation of FX-related bridge functions for REAPER API
 */
#include "fx.h"
#include "../bridge.h"
#include "../logging/logging.h"
#include <stdlib.h>
#include <string.h>

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

/**
 * REAPER's TrackFX_FormatParamValue function - formats a value without changing the parameter
 */
void plugin_bridge_call_track_fx_format_param_value(void* func_ptr, void* track, int fx_idx, int param_idx, double value, char* buf, int buf_size) {
    LOG_DEBUG("Called with func_ptr=%p, track=%p, fx_idx=%d, param_idx=%d, value=%f, buf=%p, buf_size=%d",
              func_ptr, track, fx_idx, param_idx, value, buf, buf_size);
    
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
    
    void (*track_fx_format_param_value)(void*, int, int, double, char*, int) =
        (void (*)(void*, int, int, double, char*, int))func_ptr;
    
    LOG_DEBUG("Calling TrackFX_FormatParamValue with track=%p, fx_idx=%d, param_idx=%d, value=%f",
              track, fx_idx, param_idx, value);
    track_fx_format_param_value(track, fx_idx, param_idx, value, buf, buf_size);
    LOG_DEBUG("TrackFX_FormatParamValue call completed with result: %s", buf);
}

/**
 * Function to batch retrieve all FX parameters in a single call
 * This reduces the number of C-Go crossings dramatically
 */
bool plugin_bridge_batch_get_fx_parameters(void* track, int fx_idx, fx_param_t* params, 
    int max_params, int* out_param_count) {
    LOG_DEBUG("Called with track=%p, fx_idx=%d, params=%p, max_params=%d", 
    track, fx_idx, params, max_params);

    // Verify input pointers
    if (!track || !params || !out_param_count || max_params <= 0) {
        LOG_ERROR("Invalid parameters: track=%p, params=%p, out_param_count=%p, max_params=%d",
        track, params, out_param_count, max_params);
        return false;
    }

    // Get the GetFunc function using our bridge
    void* getFuncPtr = plugin_bridge_get_get_func();
    if (!getFuncPtr) {
        LOG_ERROR("Failed to get GetFunc pointer");
        return false;
    }

    // Get the number of parameters function
    void* getParamCountFunc = NULL;
    {
        char funcName[64] = "TrackFX_GetNumParams";
        getParamCountFunc = plugin_bridge_call_get_func(getFuncPtr, funcName);
        if (!getParamCountFunc) {
            LOG_ERROR("Failed to get TrackFX_GetNumParams function pointer");
            return false;
        }
    }

    // Get the parameter functions we need
    void* getParamNameFunc = NULL;
    void* getParamFunc = NULL;
    void* getFormattedFunc = NULL;

    {
        char funcName[64] = "TrackFX_GetParamName";
        getParamNameFunc = plugin_bridge_call_get_func(getFuncPtr, funcName);
        if (!getParamNameFunc) {
            LOG_ERROR("Failed to get TrackFX_GetParamName function pointer");
            return false;
        }
    }

    {
        char funcName[64] = "TrackFX_GetParam";
        getParamFunc = plugin_bridge_call_get_func(getFuncPtr, funcName);
        if (!getParamFunc) {
            LOG_ERROR("Failed to get TrackFX_GetParam function pointer");
            return false;
        }
    }

    {
        char funcName[64] = "TrackFX_GetFormattedParamValue";
        getFormattedFunc = plugin_bridge_call_get_func(getFuncPtr, funcName);
        if (!getFormattedFunc) {
            LOG_ERROR("Failed to get TrackFX_GetFormattedParamValue function pointer");
            return false;
        }
    }

    // Cast the function pointers to their proper types
    int (*track_fx_get_param_count)(void*, int) = 
    (int (*)(void*, int))getParamCountFunc;

    void (*track_fx_get_param_name)(void*, int, int, char*, int) = 
    (void (*)(void*, int, int, char*, int))getParamNameFunc;

    double (*track_fx_get_param)(void*, int, int, double*, double*) = 
    (double (*)(void*, int, int, double*, double*))getParamFunc;

    void (*track_fx_get_param_formatted)(void*, int, int, char*, int) = 
    (void (*)(void*, int, int, char*, int))getFormattedFunc;

    // Get parameter count
    int param_count = track_fx_get_param_count(track, fx_idx);
    LOG_DEBUG("FX parameter count: %d", param_count);

    if (param_count <= 0) {
        LOG_WARNING("FX has no parameters (count=%d)", param_count);
        *out_param_count = 0;
        return true; // Not an error, just no parameters
    }

    // Limit to max_params
    if (param_count > max_params) {
        LOG_WARNING("Parameter count (%d) exceeds max_params (%d), limiting to max_params", 
        param_count, max_params);
        param_count = max_params;
    }

    // Get all parameter data
    for (int i = 0; i < param_count; i++) {
        // Get parameter name
        track_fx_get_param_name(track, fx_idx, i, params[i].name, sizeof(params[i].name));

        // Get parameter value with min/max
        double min = 0, max = 0;
        params[i].value = track_fx_get_param(track, fx_idx, i, &min, &max);
        params[i].min = min;
        params[i].max = max;

        // Get formatted value
        track_fx_get_param_formatted(track, fx_idx, i, params[i].formatted, sizeof(params[i].formatted));

        LOG_DEBUG("Parameter %d: name=%s, value=%f, min=%f, max=%f, formatted=%s", 
        i, params[i].name, params[i].value, params[i].min, params[i].max, params[i].formatted);
    }

    // Return parameter count
    *out_param_count = param_count;
    LOG_DEBUG("Successfully retrieved %d parameters", param_count);

    return true;
}

/**
 * Function to batch format parameter values in a single call
 * This reduces the number of C-Go crossings dramatically
 */
bool plugin_bridge_batch_format_fx_parameters(void* track, fx_param_format_t* params, int param_count) {
    LOG_DEBUG("Called with track=%p, params=%p, param_count=%d",
              track, params, param_count);
    
    // Verify input pointers
    if (!track || !params || param_count <= 0) {
        LOG_ERROR("Invalid parameters: track=%p, params=%p, param_count=%d",
                 track, params, param_count);
        return false;
    }
    
    // Get the GetFunc function using our bridge
    void* getFuncPtr = plugin_bridge_get_get_func();
    if (!getFuncPtr) {
        LOG_ERROR("Failed to get GetFunc pointer");
        return false;
    }
    
    // Get the TrackFX_FormatParamValue function
    void* formatValueFunc = NULL;
    {
        char funcName[64] = "TrackFX_FormatParamValue";
        formatValueFunc = plugin_bridge_call_get_func(getFuncPtr, funcName);
        if (!formatValueFunc) {
            LOG_ERROR("Failed to get TrackFX_FormatParamValue function pointer");
            return false;
        }
    }
    
    // Format all parameter values
    for (int i = 0; i < param_count; i++) {
        // Format the parameter value
        plugin_bridge_call_track_fx_format_param_value(
            formatValueFunc,
            track,
            params[i].fx_index,
            params[i].param_index,
            params[i].value,
            params[i].formatted,
            sizeof(params[i].formatted)
        );
        
        LOG_DEBUG("Parameter %d: fx_index=%d, param_index=%d, value=%f, formatted=%s",
                 i, params[i].fx_index, params[i].param_index, params[i].value, params[i].formatted);
    }
    
    LOG_DEBUG("Successfully formatted %d parameters", param_count);
    return true;
}

/**
 * Function to batch apply multiple parameter changes in a single call
 * This reduces the number of C-Go crossings dramatically
 */
bool plugin_bridge_batch_set_fx_parameters(void* track, fx_param_change_t* changes, int change_count) {
    LOG_DEBUG("Called with track=%p, changes=%p, change_count=%d",
              track, changes, change_count);
    
    // Verify input pointers
    if (!track || !changes || change_count <= 0) {
        LOG_ERROR("Invalid parameters: track=%p, changes=%p, change_count=%d",
                 track, changes, change_count);
        return false;
    }
    
    // Get the GetFunc function using our bridge
    void* getFuncPtr = plugin_bridge_get_get_func();
    if (!getFuncPtr) {
        LOG_ERROR("Failed to get GetFunc pointer");
        return false;
    }
    
    // Get the TrackFX_SetParam function
    void* setParamFunc = NULL;
    {
        char funcName[64] = "TrackFX_SetParam";
        setParamFunc = plugin_bridge_call_get_func(getFuncPtr, funcName);
        if (!setParamFunc) {
            LOG_ERROR("Failed to get TrackFX_SetParam function pointer");
            return false;
        }
    }
    
    // Apply all parameter changes
    bool all_success = true;
    for (int i = 0; i < change_count; i++) {
        // Get the change data
        int fx_idx = changes[i].fx_index;
        int param_idx = changes[i].param_index;
        double value = changes[i].value;
        
        // Apply the parameter change
        bool success = plugin_bridge_call_track_fx_set_param(
            setParamFunc,
            track,
            fx_idx,
            param_idx,
            value
        );
        
        if (!success) {
            LOG_ERROR("Failed to set parameter: fx_index=%d, param_index=%d, value=%f",
                     fx_idx, param_idx, value);
            all_success = false;
        } else {
            LOG_DEBUG("Parameter set: fx_index=%d, param_index=%d, value=%f",
                     fx_idx, param_idx, value);
        }
    }
    
    LOG_DEBUG("Applied %d parameter changes, success=%d", change_count, all_success);
    return all_success;
}

/**
 * Function to get parameters from multiple tracks and FX in a single call
 * This dramatically reduces CGo transitions when working with multiple tracks
*/
bool plugin_bridge_batch_get_multi_track_fx_parameters(
    void** tracks, int track_count, 
    int** fx_indices, int* fx_counts,
    fx_param_t** param_buffers, int* param_counts) {
    
    LOG_DEBUG("Called with tracks=%p, track_count=%d", tracks, track_count);
    
    // Verify input pointers
    if (!tracks || track_count <= 0 || !fx_indices || !fx_counts || !param_buffers || !param_counts) {
        LOG_ERROR("Invalid parameters in batch_get_multi_track_fx_parameters");
        return false;
    }
    
    // Get the GetFunc function
    void* getFuncPtr = plugin_bridge_get_get_func();
    if (!getFuncPtr) {
        LOG_ERROR("Failed to get GetFunc pointer");
        return false;
    }
    
    // Get function pointers we'll need
    void* getParamCountFunc = NULL;
    void* getParamNameFunc = NULL;
    void* getParamFunc = NULL;
    void* getFormattedFunc = NULL;
    
    {
        char funcName[64] = "TrackFX_GetNumParams";
        getParamCountFunc = plugin_bridge_call_get_func(getFuncPtr, funcName);
        if (!getParamCountFunc) {
            LOG_ERROR("Failed to get TrackFX_GetNumParams function pointer");
            return false;
        }
    }
    
    {
        char funcName[64] = "TrackFX_GetParamName";
        getParamNameFunc = plugin_bridge_call_get_func(getFuncPtr, funcName);
        if (!getParamNameFunc) {
            LOG_ERROR("Failed to get TrackFX_GetParamName function pointer");
            return false;
        }
    }
    
    {
        char funcName[64] = "TrackFX_GetParam";
        getParamFunc = plugin_bridge_call_get_func(getFuncPtr, funcName);
        if (!getParamFunc) {
            LOG_ERROR("Failed to get TrackFX_GetParam function pointer");
            return false;
        }
    }
    
    {
        char funcName[64] = "TrackFX_GetFormattedParamValue";
        getFormattedFunc = plugin_bridge_call_get_func(getFuncPtr, funcName);
        if (!getFormattedFunc) {
            LOG_ERROR("Failed to get TrackFX_GetFormattedParamValue function pointer");
            return false;
        }
    }
    
    // Cast function pointers to their proper types
    int (*track_fx_get_param_count)(void*, int) = 
        (int (*)(void*, int))getParamCountFunc;
    
    void (*track_fx_get_param_name)(void*, int, int, char*, int) = 
        (void (*)(void*, int, int, char*, int))getParamNameFunc;
    
    double (*track_fx_get_param)(void*, int, int, double*, double*) = 
        (double (*)(void*, int, int, double*, double*))getParamFunc;
    
    void (*track_fx_get_param_formatted)(void*, int, int, char*, int) = 
        (void (*)(void*, int, int, char*, int))getFormattedFunc;
    
    // Process each track
    for (int t = 0; t < track_count; t++) {
        void* track = tracks[t];
        if (!track) {
            LOG_WARNING("Null track pointer at index %d, skipping", t);
            continue;
        }
        
        int fx_count = fx_counts[t];
        if (fx_count <= 0) {
            LOG_WARNING("No FX to process for track %d", t);
            continue;
        }
        
        LOG_DEBUG("Processing track %d with %d FX", t, fx_count);
        
        // Process each FX
        for (int f = 0; f < fx_count; f++) {
            int fx_idx = fx_indices[t][f];
            
            // Calculate the buffer index for this FX's parameters
            // This assumes param_buffers is laid out as a 2D array where each row corresponds to an FX
            int buffer_idx = 0;
            for (int i = 0; i < t; i++) {
                buffer_idx += fx_counts[i];
            }
            buffer_idx += f;
            
            // Get parameter count for this FX
            int param_count = track_fx_get_param_count(track, fx_idx);
            if (param_count <= 0) {
                LOG_WARNING("FX %d on track %d has no parameters", fx_idx, t);
                param_counts[buffer_idx] = 0;
                continue;
            }
            
            LOG_DEBUG("FX %d on track %d has %d parameters", fx_idx, t, param_count);
            
            // Get all parameters for this FX
            fx_param_t* params = param_buffers[buffer_idx];
            
            for (int p = 0; p < param_count; p++) {
                // Get parameter name
                track_fx_get_param_name(track, fx_idx, p, params[p].name, sizeof(params[p].name));
                
                // Get parameter value with min/max
                double min = 0, max = 0;
                params[p].value = track_fx_get_param(track, fx_idx, p, &min, &max);
                params[p].min = min;
                params[p].max = max;
                
                // Get formatted value
                track_fx_get_param_formatted(track, fx_idx, p, params[p].formatted, sizeof(params[p].formatted));
                
                LOG_DEBUG("Parameter %d: name=%s, value=%f, min=%f, max=%f, formatted=%s",
                         p, params[p].name, params[p].value, params[p].min, params[p].max, params[p].formatted);
            }
            
            // Store parameter count
            param_counts[buffer_idx] = param_count;
        }
    }
    
    LOG_DEBUG("Successfully processed %d tracks", track_count);
    return true;
}

/**
 * Function to format parameter values for multiple tracks in a single call
 */
bool plugin_bridge_batch_format_multi_fx_parameters(
    void** tracks, 
    fx_param_format_t* params, 
    int param_count) {
    
    LOG_DEBUG("Called with tracks=%p, params=%p, param_count=%d", tracks, params, param_count);
    
    // Verify input pointers
    if (!tracks || !params || param_count <= 0) {
        LOG_ERROR("Invalid parameters in batch_format_multi_fx_parameters");
        return false;
    }
    
    // Get the GetFunc function
    void* getFuncPtr = plugin_bridge_get_get_func();
    if (!getFuncPtr) {
        LOG_ERROR("Failed to get GetFunc pointer");
        return false;
    }
    
    // Get the TrackFX_FormatParamValue function
    void* formatValueFunc = NULL;
    {
        char funcName[64] = "TrackFX_FormatParamValue";
        formatValueFunc = plugin_bridge_call_get_func(getFuncPtr, funcName);
        if (!formatValueFunc) {
            LOG_ERROR("Failed to get TrackFX_FormatParamValue function pointer");
            return false;
        }
    }
    
    // Format all parameter values
    for (int i = 0; i < param_count; i++) {
        // Get the track index from the params structure
        int track_idx = params[i].track_index;
        
        // Make sure track index is valid
        if (track_idx < 0 || tracks[track_idx] == NULL) {
            LOG_ERROR("Invalid track index %d or NULL track pointer", track_idx);
            return false;
        }
        
        // Format the parameter value
        plugin_bridge_call_track_fx_format_param_value(
            formatValueFunc,
            tracks[track_idx],
            params[i].fx_index,
            params[i].param_index,
            params[i].value,
            params[i].formatted,
            sizeof(params[i].formatted)
        );
        
        LOG_DEBUG("Parameter %d: track=%d, fx_index=%d, param_index=%d, value=%f, formatted=%s",
                i, track_idx, params[i].fx_index, params[i].param_index, params[i].value, params[i].formatted);
    }
    
    LOG_DEBUG("Successfully formatted %d parameters", param_count);
    return true;
}

/**
 * Function to apply parameter changes to multiple tracks in a single call
 */
bool plugin_bridge_batch_set_multi_track_fx_parameters(
    void** tracks,
    fx_param_multi_change_t* changes,
    int change_count) {
    
    LOG_DEBUG("Called with tracks=%p, changes=%p, change_count=%d", tracks, changes, change_count);
    
    // Verify input pointers
    if (!tracks || !changes || change_count <= 0) {
        LOG_ERROR("Invalid parameters in batch_set_multi_track_fx_parameters");
        return false;
    }
    
    // Get the GetFunc function
    void* getFuncPtr = plugin_bridge_get_get_func();
    if (!getFuncPtr) {
        LOG_ERROR("Failed to get GetFunc pointer");
        return false;
    }
    
    // Get the TrackFX_SetParam function
    void* setParamFunc = NULL;
    {
        char funcName[64] = "TrackFX_SetParam";
        setParamFunc = plugin_bridge_call_get_func(getFuncPtr, funcName);
        if (!setParamFunc) {
            LOG_ERROR("Failed to get TrackFX_SetParam function pointer");
            return false;
        }
    }
    
    // Apply all parameter changes
    bool all_success = true;
    for (int i = 0; i < change_count; i++) {
        // Get the track index
        int track_idx = changes[i].track_index;
        
        // Make sure track index is valid
        if (track_idx < 0 || tracks[track_idx] == NULL) {
            LOG_ERROR("Invalid track index %d or NULL track pointer", track_idx);
            all_success = false;
            continue;
        }
        
        // Get the change data
        int fx_idx = changes[i].fx_index;
        int param_idx = changes[i].param_index;
        double value = changes[i].value;
        
        // Apply the parameter change
        bool success = plugin_bridge_call_track_fx_set_param(
            setParamFunc,
            tracks[track_idx],
            fx_idx,
            param_idx,
            value
        );
        
        if (!success) {
            LOG_ERROR("Failed to set parameter: track=%d, fx_index=%d, param_index=%d, value=%f",
                     track_idx, fx_idx, param_idx, value);
            all_success = false;
        } else {
            LOG_DEBUG("Parameter set: track=%d, fx_index=%d, param_index=%d, value=%f",
                     track_idx, fx_idx, param_idx, value);
        }
    }
    
    LOG_DEBUG("Applied %d parameter changes across multiple tracks, success=%d", change_count, all_success);
    return all_success;
}

/**
 * REAPER's TrackFX_GetParameterStepSizes function
 */
bool plugin_bridge_call_track_fx_get_parameter_step_sizes(void* func_ptr, void* track, int fx_idx, int param_idx, 
    double* step, double* small_step, double* large_step, bool* is_toggle) {
    LOG_DEBUG("Called with func_ptr=%p, track=%p, fx_idx=%d, param_idx=%d", 
    func_ptr, track, fx_idx, param_idx);

    // Verify input pointers aren't NULL
    if (!func_ptr || !track) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, track=%p", func_ptr, track);
        return false;
    }

    // Output pointers can be NULL if caller doesn't need those values

    bool (*track_fx_get_param_step_sizes)(void*, int, int, double*, double*, double*, bool*) = 
    (bool (*)(void*, int, int, double*, double*, double*, bool*))func_ptr;

    LOG_DEBUG("Calling TrackFX_GetParameterStepSizes with track=%p, fx_idx=%d, param_idx=%d", 
    track, fx_idx, param_idx);

    bool result = track_fx_get_param_step_sizes(track, fx_idx, param_idx, step, small_step, large_step, is_toggle);

    // Log results for debugging
    if (result) {
        // Log values only if the pointers are non-NULL
        LOG_DEBUG("TrackFX_GetParameterStepSizes call completed with result: %d", result);

        if (step) LOG_DEBUG("  step: %f", *step);
        if (small_step) LOG_DEBUG("  small_step: %f", *small_step);
        if (large_step) LOG_DEBUG("  large_step: %f", *large_step);
        if (is_toggle) LOG_DEBUG("  is_toggle: %d", *is_toggle);
    } else {
       LOG_DEBUG("TrackFX_GetParameterStepSizes call failed");
    }

    return result;
}
