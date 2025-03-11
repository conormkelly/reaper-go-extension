/**
 * Track-related bridge functions for REAPER API
 */
#ifndef REAPER_EXT_API_TRACKS_H
#define REAPER_EXT_API_TRACKS_H

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// Track-related API functions
void* plugin_bridge_call_get_selected_track(void* func_ptr, int proj, int seltrackidx);
double plugin_bridge_call_track_get_info_value(void* func_ptr, void* track, const char* param);
bool plugin_bridge_call_get_track_name(void* func_ptr, void* track, char* buf, int buf_size, int* flags);

#ifdef __cplusplus
}
#endif

#endif // REAPER_EXT_API_TRACKS_H
