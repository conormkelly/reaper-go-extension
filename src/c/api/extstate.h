/**
 * ExtState-related bridge functions for REAPER API
 */
#ifndef REAPER_EXT_API_EXTSTATE_H
#define REAPER_EXT_API_EXTSTATE_H

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// ExtState-related API functions
const char* plugin_bridge_call_get_ext_state(void* func_ptr, const char* section, const char* key);
void plugin_bridge_call_set_ext_state(void* func_ptr, const char* section, const char* key, 
                                     const char* value, int persist);
bool plugin_bridge_call_has_ext_state(void* func_ptr, const char* section, const char* key);
void plugin_bridge_call_delete_ext_state(void* func_ptr, const char* section, const char* key);

#ifdef __cplusplus
}
#endif

#endif // REAPER_EXT_API_EXTSTATE_H
