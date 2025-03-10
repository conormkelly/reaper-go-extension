// reaper_plugin_bridge.h
#ifndef REAPER_PLUGIN_BRIDGE_H
#define REAPER_PLUGIN_BRIDGE_H

// Forward declarations for types that aren't defined in our main includes
typedef struct MIDI_eventlist MIDI_eventlist;
typedef struct accelerator_register_t accelerator_register_t;
typedef struct ProjectStateContext ProjectStateContext;
typedef struct PCM_source PCM_source;
typedef struct project_config_extension_t project_config_extension_t;

#include <stdlib.h>
#include <stdbool.h>
#include "../sdk/reaper_plugin.h"

// Helper functions for the Go code to call
#ifdef __cplusplus
extern "C" {
#endif

void* plugin_bridge_call_get_func(void* get_func_ptr, const char* name);
void plugin_bridge_call_show_console_msg(void* func_ptr, const char* message);
int plugin_bridge_call_register(void* register_func_ptr, const char* name, void* info);
void* plugin_bridge_call_get_selected_track(void* func_ptr, int proj, int seltrackidx);
int plugin_bridge_call_track_fx_get_count(void* func_ptr, void* track);
void plugin_bridge_call_track_fx_get_name(void* func_ptr, void* track, int fx_idx, char* buf, int buf_size);
int plugin_bridge_call_track_fx_get_param_count(void* func_ptr, void* track, int fx_idx);
void plugin_bridge_call_track_fx_get_param_name(void* func_ptr, void* track, int fx_idx, int param_idx, char* buf, int buf_size);
double plugin_bridge_call_track_fx_get_param(void* func_ptr, void* track, int fx_idx, int param_idx, double* minval, double* maxval);
void plugin_bridge_call_track_fx_get_param_formatted(void* func_ptr, void* track, int fx_idx, int param_idx, char* buf, int buf_size);
bool plugin_bridge_call_track_fx_set_param(void* func_ptr, void* track, int fx_idx, int param_idx, double val);

// Track information functions
double plugin_bridge_call_track_get_info_value(void* func_ptr, void* track, const char* param);
bool plugin_bridge_call_get_track_name(void* func_ptr, void* track, char* buf, int buf_size, int* flags);

// GetUserInputs - Simple form dialog
bool plugin_bridge_call_get_user_inputs(void* func_ptr, const char* title, int num_inputs, 
    const char* captions, char* values, int values_sz);

// ShowMessageBox - Standard message box
int plugin_bridge_call_show_message_box(void* func_ptr, const char* text, const char* title, int type);

void plugin_bridge_set_get_func(void* get_func_ptr);
void* plugin_bridge_get_get_func();

// Structure to hold parameter data
typedef struct {
    char name[256];
    double value;
    double min;
    double max;
    char formatted[256];
} fx_param_t;

// Function to get all FX parameters in a single call
bool plugin_bridge_batch_get_fx_parameters(void* track, int fx_idx, fx_param_t* params, 
                                        int max_params, int* out_param_count);


// GetExtState
const char* plugin_bridge_call_get_ext_state(void* func_ptr, const char* section, const char* key);

// SetExtState
void plugin_bridge_call_set_ext_state(void* func_ptr, const char* section, const char* key, 
                                    const char* value, int persist);

// HasExtState
bool plugin_bridge_call_has_ext_state(void* func_ptr, const char* section, const char* key);

// DeleteExtState
void plugin_bridge_call_delete_ext_state(void* func_ptr, const char* section, const char* key);

// Forward declaration of the Go functions

// GoReaperPluginEntry is the entry point called by REAPER. This function bridges between 
// REAPER's C API and our Go code. hInstance is the module handle, rec contains REAPER API functions
extern int GoReaperPluginEntry(void* hInstance, void* rec);

// Command hook callbacks
extern int goHookCommandProc(int commandId, int flag);
extern int goHookCommandProc2(void* section, int commandId, int val, int valhw, int relmode, void* hwnd, void* proj);

#ifdef __cplusplus
}
#endif

#endif // REAPER_PLUGIN_BRIDGE_H
