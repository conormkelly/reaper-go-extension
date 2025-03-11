/**
 * Implementation of track-related bridge functions for REAPER API
 */
#include "tracks.h"
#include "../bridge.h"
#include "../logging/logging.h"
#include <stdlib.h>

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
 * Get track information value
 */
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

/**
 * Get track name
 */
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
