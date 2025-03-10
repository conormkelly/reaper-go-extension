#ifndef KRBRIDGE_H
#define KRBRIDGE_H

#include <stdbool.h>

// Context structure for passing data between Go and Objective-C
typedef struct {
    const char* title;
    bool key_exists;
    const char* message;
    bool success;
} KRContext;

// Function declarations that will be called from Go
bool kr_show_window(const char* title, bool key_exists, const char* message);
bool kr_update_message(bool key_exists, const char* message);
void kr_close_window(void);
bool kr_window_exists(void);

// Callback from Objective-C to Go
extern void go_process_keyring_key(char* key);

#endif /* KRBRIDGE_H */
