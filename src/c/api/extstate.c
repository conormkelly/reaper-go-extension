/**
 * Implementation of ExtState-related bridge functions for REAPER API
 */
#include "extstate.h"
#include "../bridge.h"
#include "../logging/logging.h"
#include <stdlib.h>

/**
 * Function to get extended state
 */
const char* plugin_bridge_call_get_ext_state(void* func_ptr, const char* section, const char* key) {
    LOG_DEBUG("Called with func_ptr=%p, section=%s, key=%s", 
              func_ptr, section ? section : "NULL", key ? key : "NULL");
    
    // Verify input parameters
    if (!func_ptr || !section || !key) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, section=%p, key=%p", 
                 func_ptr, section, key);
        return NULL;
    }
    
    // Cast function pointer to correct type
    const char* (*get_ext_state)(const char*, const char*) = 
        (const char* (*)(const char*, const char*))func_ptr;
    
    // Call GetExtState
    LOG_DEBUG("Calling GetExtState with section=%s, key=%s", section, key);
    const char* result = get_ext_state(section, key);
    LOG_DEBUG("GetExtState call completed with result: %s", result ? result : "NULL");
    
    return result;
}

/**
 * Function to set extended state
 */
void plugin_bridge_call_set_ext_state(void* func_ptr, const char* section, const char* key, 
                                     const char* value, int persist) {
    LOG_DEBUG("Called with func_ptr=%p, section=%s, key=%s, value=%s, persist=%d", 
              func_ptr, section ? section : "NULL", key ? key : "NULL", 
              value ? value : "NULL", persist);
    
    // Verify input parameters
    if (!func_ptr || !section || !key || !value) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, section=%p, key=%p, value=%p", 
                 func_ptr, section, key, value);
        return;
    }
    
    // Cast function pointer to correct type
    void (*set_ext_state)(const char*, const char*, const char*, int) = 
        (void (*)(const char*, const char*, const char*, int))func_ptr;
    
    // Call SetExtState
    LOG_DEBUG("Calling SetExtState with section=%s, key=%s, value=%s, persist=%d", 
             section, key, value, persist);
    set_ext_state(section, key, value, persist);
    LOG_DEBUG("SetExtState call completed");
}

/**
 * Function to check if extended state exists
 */
bool plugin_bridge_call_has_ext_state(void* func_ptr, const char* section, const char* key) {
    LOG_DEBUG("Called with func_ptr=%p, section=%s, key=%s", 
              func_ptr, section ? section : "NULL", key ? key : "NULL");
    
    // Verify input parameters
    if (!func_ptr || !section || !key) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, section=%p, key=%p", 
                 func_ptr, section, key);
        return false;
    }
    
    // Cast function pointer to correct type
    int (*has_ext_state)(const char*, const char*) = 
        (int (*)(const char*, const char*))func_ptr;
    
    // Call HasExtState
    LOG_DEBUG("Calling HasExtState with section=%s, key=%s", section, key);
    int result = has_ext_state(section, key);
    LOG_DEBUG("HasExtState call completed with result: %d", result);
    
    return result != 0;
}

/**
 * Function to delete extended state
 */
void plugin_bridge_call_delete_ext_state(void* func_ptr, const char* section, const char* key) {
    LOG_DEBUG("Called with func_ptr=%p, section=%s, key=%s", 
              func_ptr, section ? section : "NULL", key ? key : "NULL");
    
    // Verify input parameters
    if (!func_ptr || !section || !key) {
        LOG_ERROR("Invalid parameters: func_ptr=%p, section=%p, key=%p", 
                 func_ptr, section, key);
        return;
    }
    
    // Cast function pointer to correct type
    void (*delete_ext_state)(const char*, const char*) = 
        (void (*)(const char*, const char*))func_ptr;
    
    // Call DeleteExtState
    LOG_DEBUG("Calling DeleteExtState with section=%s, key=%s", section, key);
    delete_ext_state(section, key);
    LOG_DEBUG("DeleteExtState call completed");
}
