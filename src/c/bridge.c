/**
 * Main bridge implementation for REAPER plugin
 *
 * This implements a bridge between REAPER's C/C++ API and Go code.
 * It provides safe function pointer handling and type conversion between
 * the two language environments.
 */

 #include "bridge.h"
 #include "logging/logging.h"
 
 // Basic bridge functions implementation
 
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
