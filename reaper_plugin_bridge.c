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

void* plugin_bridge_call_get_selected_track(void* func_ptr, int proj, int seltrackidx) {
    void* (*get_selected_track)(int, int) = (void* (*)(int, int))func_ptr;
    return get_selected_track(proj, seltrackidx);
}

int plugin_bridge_call_track_fx_get_count(void* func_ptr, void* track) {
    int (*track_fx_get_count)(void*) = (int (*)(void*))func_ptr;
    return track_fx_get_count(track);
}

void plugin_bridge_call_track_fx_get_name(void* func_ptr, void* track, int fx_idx, char* buf, int buf_size) {
    void (*track_fx_get_name)(void*, int, char*, int) = 
        (void (*)(void*, int, char*, int))func_ptr;
    track_fx_get_name(track, fx_idx, buf, buf_size);
}

int plugin_bridge_call_track_fx_get_param_count(void* func_ptr, void* track, int fx_idx) {
    int (*track_fx_get_param_count)(void*, int) = (int (*)(void*, int))func_ptr;
    return track_fx_get_param_count(track, fx_idx);
}

void plugin_bridge_call_track_fx_get_param_name(void* func_ptr, void* track, int fx_idx, int param_idx, char* buf, int buf_size) {
    void (*track_fx_get_param_name)(void*, int, int, char*, int) = 
        (void (*)(void*, int, int, char*, int))func_ptr;
    track_fx_get_param_name(track, fx_idx, param_idx, buf, buf_size);
}

double plugin_bridge_call_track_fx_get_param(void* func_ptr, void* track, int fx_idx, int param_idx, double* minval, double* maxval) {
    double (*track_fx_get_param)(void*, int, int, double*, double*) = 
        (double (*)(void*, int, int, double*, double*))func_ptr;
    return track_fx_get_param(track, fx_idx, param_idx, minval, maxval);
}

void plugin_bridge_call_track_fx_get_param_formatted(void* func_ptr, void* track, int fx_idx, int param_idx, char* buf, int buf_size) {
    void (*track_fx_get_param_formatted)(void*, int, int, char*, int) = 
        (void (*)(void*, int, int, char*, int))func_ptr;
    track_fx_get_param_formatted(track, fx_idx, param_idx, buf, buf_size);
}

bool plugin_bridge_call_track_fx_set_param(void* func_ptr, void* track, int fx_idx, int param_idx, double val) {
    bool (*track_fx_set_param)(void*, int, int, double) = 
        (bool (*)(void*, int, int, double))func_ptr;
    return track_fx_set_param(track, fx_idx, param_idx, val);
}

static void* s_GetFunc = NULL;

// Function to set the GetFunc pointer
void plugin_bridge_set_get_func(void* get_func_ptr) {
    s_GetFunc = get_func_ptr;
}

// Function to access the GetFunc pointer
void* plugin_bridge_get_get_func() {
    return s_GetFunc;
}

// This function will be called by REAPER
REAPER_PLUGIN_DLL_EXPORT int ReaperPluginEntry(HINSTANCE hInstance, reaper_plugin_info_t* rec) {
    return GoReaperPluginEntry((void*)hInstance, (void*)rec);
}
