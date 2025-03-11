/**
 * Implementation of UI-related bridge functions for REAPER API
 */
#include "ui.h"
#include "../bridge.h"
#include "../logging/logging.h"
#include <stdlib.h>

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
