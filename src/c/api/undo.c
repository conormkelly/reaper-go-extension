/**
 * Implementation of Undo-related bridge functions for REAPER API
 */
#include "undo.h"
#include "../bridge.h"
#include "../logging/logging.h"
#include <stdlib.h>

/**
 * REAPER's Undo_BeginBlock function to start an undo block
 */
void plugin_bridge_call_undo_begin_block(void* func_ptr) {
    LOG_DEBUG("Called with func_ptr=%p", func_ptr);
    
    // Verify input pointer isn't NULL
    if (!func_ptr) {
        LOG_ERROR("Invalid parameter: func_ptr is NULL");
        return;
    }
    
    void (*undo_begin_block)() = (void (*)())func_ptr;
    LOG_DEBUG("Calling Undo_BeginBlock");
    undo_begin_block();
    LOG_DEBUG("Undo_BeginBlock call completed");
}

/**
 * REAPER's Undo_BeginBlock2 function to start an undo block with project reference
 */
void plugin_bridge_call_undo_begin_block2(void* func_ptr, void* proj) {
    LOG_DEBUG("Called with func_ptr=%p, proj=%p", func_ptr, proj);
    
    // Verify input pointer isn't NULL
    if (!func_ptr) {
        LOG_ERROR("Invalid parameter: func_ptr is NULL");
        return;
    }
    
    void (*undo_begin_block2)(void*) = (void (*)(void*))func_ptr;
    LOG_DEBUG("Calling Undo_BeginBlock2");
    undo_begin_block2(proj);
    LOG_DEBUG("Undo_BeginBlock2 call completed");
}

/**
 * REAPER's Undo_EndBlock function to end an undo block
 */
void plugin_bridge_call_undo_end_block(void* func_ptr, const char* description, int flags) {
    LOG_DEBUG("Called with func_ptr=%p, description=%s, flags=%d", 
              func_ptr, description ? description : "NULL", flags);
    
    // Verify input pointers aren't NULL
    if (!func_ptr || !description) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, description=%p", func_ptr, description);
        return;
    }
    
    void (*undo_end_block)(const char*, int) = (void (*)(const char*, int))func_ptr;
    LOG_DEBUG("Calling Undo_EndBlock with description=%s, flags=%d", description, flags);
    undo_end_block(description, flags);
    LOG_DEBUG("Undo_EndBlock call completed");
}

/**
 * REAPER's Undo_EndBlock2 function to end an undo block with project reference
 */
void plugin_bridge_call_undo_end_block2(void* func_ptr, void* proj, const char* description, int flags) {
    LOG_DEBUG("Called with func_ptr=%p, proj=%p, description=%s, flags=%d", 
              func_ptr, proj, description ? description : "NULL", flags);
    
    // Verify input pointers aren't NULL
    if (!func_ptr || !description) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, description=%p", func_ptr, description);
        return;
    }
    
    void (*undo_end_block2)(void*, const char*, int) = (void (*)(void*, const char*, int))func_ptr;
    LOG_DEBUG("Calling Undo_EndBlock2 with proj=%p, description=%s, flags=%d", proj, description, flags);
    undo_end_block2(proj, description, flags);
    LOG_DEBUG("Undo_EndBlock2 call completed");
}
