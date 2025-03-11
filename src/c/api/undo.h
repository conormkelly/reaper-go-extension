/**
 * Undo-related bridge functions for REAPER API
 */
#ifndef REAPER_EXT_API_UNDO_H
#define REAPER_EXT_API_UNDO_H

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// Undo-related API functions
void plugin_bridge_call_undo_begin_block(void* func_ptr);
void plugin_bridge_call_undo_begin_block2(void* func_ptr, void* proj);
void plugin_bridge_call_undo_end_block(void* func_ptr, const char* description, int flags);
void plugin_bridge_call_undo_end_block2(void* func_ptr, void* proj, const char* description, int flags);

#ifdef __cplusplus
}
#endif

#endif // REAPER_EXT_API_UNDO_H
