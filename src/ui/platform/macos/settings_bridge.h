#ifndef SETTINGS_BRIDGE_H
#define SETTINGS_BRIDGE_H

#include <stdbool.h>

// Context structure for passing data between Go and Objective-C
typedef struct {
    const char* title;
    const char* api_key;
    const char* model;
    double temperature;
    bool success;
} SettingsContext;

// Function declarations for bridge
bool settings_show_window(const char* title, const char* api_key, const char* model, double temperature);
void settings_close_window(void);
bool settings_window_exists(void);

// Callback from Objective-C to Go
extern void go_process_settings(char* api_key, char* model, double temperature);

#endif /* SETTINGS_BRIDGE_H */
