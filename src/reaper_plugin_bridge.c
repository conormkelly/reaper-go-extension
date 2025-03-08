// reaper_plugin_bridge.c
#include "reaper_plugin_bridge.h"

// Implementation of the bridge functions
void* plugin_bridge_call_get_func(void* get_func_ptr, const char* name) {
    void* (*get_func)(const char*) = (void* (*)(const char*))get_func_ptr;
    return get_func(name);
}

void plugin_bridge_call_show_console_msg(void* func_ptr, const char* message) {
    void (*show_console_msg)(const char*) = (void (*)(const char*))func_ptr;
    show_console_msg(message);
}

int plugin_bridge_call_register(void* register_func_ptr, const char* name, void* info) {
    int (*register_func)(const char*, void*) = (int (*)(const char*, void*))register_func_ptr;
    return register_func(name, info);
}

// This function will be called by REAPER
REAPER_PLUGIN_DLL_EXPORT int ReaperPluginEntry(HINSTANCE hInstance, reaper_plugin_info_t* rec) {
    return GoReaperPluginEntry((void*)hInstance, (void*)rec);
}
