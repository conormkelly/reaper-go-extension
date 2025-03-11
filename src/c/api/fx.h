/**
 * FX-related bridge functions for REAPER API
 */
#ifndef REAPER_EXT_API_FX_H
#define REAPER_EXT_API_FX_H

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// Structure to hold parameter data
typedef struct {
    char name[256];
    double value;
    double min;
    double max;
    char formatted[256];
} fx_param_t;

// Structure to hold parameter formatting request
typedef struct {
    int fx_index;
    int param_index;
    double value;
    char formatted[256];
} fx_param_format_t;

// Structure to hold parameter change data
typedef struct {
    int fx_index;
    int param_index;
    double value;
} fx_param_change_t;

// FX-related API functions
int plugin_bridge_call_track_fx_get_count(void* func_ptr, void* track);
void plugin_bridge_call_track_fx_get_name(void* func_ptr, void* track, int fx_idx, char* buf, int buf_size);
int plugin_bridge_call_track_fx_get_param_count(void* func_ptr, void* track, int fx_idx);
void plugin_bridge_call_track_fx_get_param_name(void* func_ptr, void* track, int fx_idx, int param_idx, char* buf, int buf_size);
double plugin_bridge_call_track_fx_get_param(void* func_ptr, void* track, int fx_idx, int param_idx, double* minval, double* maxval);
void plugin_bridge_call_track_fx_get_param_formatted(void* func_ptr, void* track, int fx_idx, int param_idx, char* buf, int buf_size);
bool plugin_bridge_call_track_fx_set_param(void* func_ptr, void* track, int fx_idx, int param_idx, double val);
void plugin_bridge_call_track_fx_format_param_value(void* func_ptr, void* track, int fx_idx, int param_idx, double value, char* buf, int buf_size);

// Batch operations for parameters
bool plugin_bridge_batch_get_fx_parameters(void* track, int fx_idx, fx_param_t* params, int max_params, int* out_param_count);
bool plugin_bridge_batch_format_fx_parameters(void* track, fx_param_format_t* params, int param_count);
bool plugin_bridge_batch_set_fx_parameters(void* track, fx_param_change_t* changes, int change_count);

#ifdef __cplusplus
}
#endif

#endif // REAPER_EXT_API_FX_H
