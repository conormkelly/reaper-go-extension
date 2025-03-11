/**
 * UI-related bridge functions for REAPER API
 */
#ifndef REAPER_EXT_API_UI_H
#define REAPER_EXT_API_UI_H

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// UI-related API functions
bool plugin_bridge_call_get_user_inputs(void* func_ptr, const char* title, int num_inputs, 
    const char* captions, char* values, int values_sz);
int plugin_bridge_call_show_message_box(void* func_ptr, const char* text, const char* title, int type);

#ifdef __cplusplus
}
#endif

#endif // REAPER_EXT_API_UI_H
