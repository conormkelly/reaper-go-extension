#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include "../c/logging/logging.h"
#import <Cocoa/Cocoa.h>
#include "krbridge.h"

// Use our core logging system
static void kr_log_to_reaper(LogLevel level, const char* message) {
    log_message_v(level, "keyringUI", message);
}

// Forward declaration of our controller class
@interface RPRKeyringController : NSObject
- (void)closeButtonClicked:(id)sender;
- (void)saveButtonClicked:(id)sender;
@end

// Global references for window and controls
static NSWindow* kr_window = nil;
static RPRKeyringController* kr_controller = nil;
static NSTextField* kr_message_field = nil;
static NSTextField* kr_input_field = nil;
static NSButton* kr_save_button = nil;
static NSButton* kr_close_button = nil;

// Implementation of our controller class
@implementation RPRKeyringController

- (void)closeButtonClicked:(id)sender {
    kr_log_to_reaper(LOG_DEBUG, "Close button clicked");

    if (kr_window != nil) {
        [kr_window close];
        kr_window = nil;
    }
}

- (void)saveButtonClicked:(id)sender {
    kr_log_to_reaper(LOG_DEBUG, "Save button clicked");

    if (kr_input_field == nil) {
        kr_log_to_reaper(LOG_ERROR, "Input field not initialized");
        return;
    }

    NSString* key = [kr_input_field stringValue];
    
    if ([key length] == 0) {
        // Show error for empty key
        if (kr_message_field != nil) {
            [kr_message_field setStringValue:@"Please enter a key value."];
        }
        return;
    }
    
    // Get the key as C string
    const char* keyStr = [key UTF8String];
    
    // Create a persistent copy of the key
    char* keyCopy = (char*)malloc(strlen(keyStr) + 1);
    if (keyCopy) {
        strcpy(keyCopy, keyStr);
        
        // Call the Go function directly
        kr_log_to_reaper(LOG_DEBUG, "Calling go_process_keyring_key function");
        go_process_keyring_key(keyCopy);
        
        // Free memory allocated for key (Go will have copied it)
        free(keyCopy);
    } else {
        kr_log_to_reaper(LOG_ERROR, "Failed to allocate memory for key");
        if (kr_message_field != nil) {
            [kr_message_field setStringValue:@"Memory error processing key"];
        }
    }
}

@end

// Function to create and show the keyring window - internal
static void kr_show_window_on_main_thread(void* context) {
    kr_log_to_reaper(LOG_DEBUG, "Entering kr_show_window_on_main_thread");

    KRContext* ctx = (KRContext*)context;
    if (!ctx) {
        kr_log_to_reaper(LOG_ERROR, "Null context in kr_show_window_on_main_thread");
        return;
    }

    @autoreleasepool {
        @try {
            // Check if window already exists
            if (kr_window != nil) {
                kr_log_to_reaper(LOG_DEBUG, "Window already exists, bringing to front");
                [kr_window makeKeyAndOrderFront:nil];
                ctx->success = true;
                return;
            }

            // Create controller if needed
            if (kr_controller == nil) {
                kr_controller = [[RPRKeyringController alloc] init];
                if (!kr_controller) {
                    kr_log_to_reaper(LOG_ERROR, "Failed to create controller");
                    ctx->success = false;
                    return;
                }
            }

            // Create the window
            NSRect frame = NSMakeRect(100, 100, 400, 200);
            NSWindow* window = [[NSWindow alloc]
                initWithContentRect:frame
                styleMask:NSWindowStyleMaskTitled|NSWindowStyleMaskClosable
                backing:NSBackingStoreBuffered
                defer:NO];

            if (window == nil) {
                kr_log_to_reaper(LOG_ERROR, "Failed to create window");
                ctx->success = false;
                return;
            }

            // Set window properties
            [window setTitle:[NSString stringWithUTF8String:ctx->title]];
            [window setReleasedWhenClosed:NO]; // Important: Don't release on close

            // Get the content view
            NSView* contentView = [window contentView];

            // Create a message label
            NSTextField* messageLabel = [[NSTextField alloc] initWithFrame:NSMakeRect(20, 140, 360, 40)];
            [messageLabel setStringValue:[NSString stringWithUTF8String:ctx->message]];
            [messageLabel setBezeled:NO];
            [messageLabel setDrawsBackground:NO];
            [messageLabel setEditable:NO];
            [messageLabel setSelectable:NO];
            kr_message_field = messageLabel;
            [contentView addSubview:messageLabel];

            // If key doesn't exist, show input field
            if (!ctx->key_exists) {
                // Create input field label
                NSTextField* inputLabel = [[NSTextField alloc] initWithFrame:NSMakeRect(20, 100, 100, 24)];
                [inputLabel setStringValue:@"Enter API Key:"];
                [inputLabel setBezeled:NO];
                [inputLabel setDrawsBackground:NO];
                [inputLabel setEditable:NO];
                [inputLabel setSelectable:NO];
                [contentView addSubview:inputLabel];

                // Create input field
                NSTextField* inputField = [[NSTextField alloc] initWithFrame:NSMakeRect(120, 100, 260, 24)];
                [inputField setPlaceholderString:@"Enter your API key"];
                kr_input_field = inputField;
                [contentView addSubview:inputField];

                // Create save button
                NSButton* saveButton = [[NSButton alloc] initWithFrame:NSMakeRect(120, 60, 100, 32)];
                [saveButton setTitle:@"Save"];
                [saveButton setBezelStyle:NSBezelStyleRounded];
                [saveButton setTarget:kr_controller];
                [saveButton setAction:@selector(saveButtonClicked:)];
                kr_save_button = saveButton;
                [contentView addSubview:saveButton];
            }

            // Create close button
            NSButton* closeButton = [[NSButton alloc] initWithFrame:NSMakeRect(230, 60, 100, 32)];
            [closeButton setTitle:@"OK"];
            [closeButton setBezelStyle:NSBezelStyleRounded];
            [closeButton setTarget:kr_controller];
            [closeButton setAction:@selector(closeButtonClicked:)];
            kr_close_button = closeButton;
            [contentView addSubview:closeButton];

            // Store window reference
            kr_window = window;

            // Show the window
            [window center];
            [window makeKeyAndOrderFront:nil];

            kr_log_to_reaper(LOG_INFO, "Keyring window created and shown successfully");
            ctx->success = true;
        }
        @catch (NSException *exception) {
            kr_log_to_reaper(LOG_ERROR, "EXCEPTION creating keyring window");
            NSLog(@"Exception: %@", exception);
            ctx->success = false;
        }
    }
}

// Function to update the message in the window - internal
static void kr_update_message_on_main_thread(void* context) {
    kr_log_to_reaper(LOG_DEBUG, "Entering kr_update_message_on_main_thread");

    KRContext* ctx = (KRContext*)context;
    if (!ctx) {
        kr_log_to_reaper(LOG_ERROR, "Null context in kr_update_message_on_main_thread");
        return;
    }

    @autoreleasepool {
        @try {
            // Check if window exists
            if (kr_window == nil || kr_message_field == nil) {
                kr_log_to_reaper(LOG_ERROR, "Window or message field doesn't exist");
                ctx->success = false;
                return;
            }

            // Update the message
            [kr_message_field setStringValue:[NSString stringWithUTF8String:ctx->message]];
            
            // Hide input field and save button if key was saved
            if (ctx->key_exists) {
                if (kr_input_field != nil) {
                    [kr_input_field setHidden:YES];
                }
                if (kr_save_button != nil) {
                    [kr_save_button setHidden:YES];
                }
            }

            ctx->success = true;
        }
        @catch (NSException *exception) {
            kr_log_to_reaper(LOG_ERROR, "EXCEPTION updating keyring message");
            NSLog(@"Exception: %@", exception);
            ctx->success = false;
        }
    }
}

// Function to close window on main thread - internal
static void kr_close_window_on_main_thread(void* context) {
    kr_log_to_reaper(LOG_DEBUG, "Entering kr_close_window_on_main_thread");

    @try {
        if (kr_window != nil) {
            [kr_window close];
            kr_window = nil;
            kr_message_field = nil;
            kr_input_field = nil;
            kr_save_button = nil;
            kr_close_button = nil;
            kr_log_to_reaper(LOG_INFO, "Keyring window closed successfully");
        } else {
            kr_log_to_reaper(LOG_DEBUG, "Keyring window already nil when trying to close");
        }
    }
    @catch (NSException *exception) {
        kr_log_to_reaper(LOG_ERROR, "EXCEPTION closing keyring window");
        NSLog(@"Exception: %@", exception);
    }
}

// Execute function on main thread - internal
static bool kr_execute_on_main_thread(void (*func)(void*), void* context) {
    kr_log_to_reaper(LOG_DEBUG, "Entering kr_execute_on_main_thread");

    // Check if already on main thread
    if ([NSThread isMainThread]) {
        kr_log_to_reaper(LOG_DEBUG, "Already on main thread, executing directly");
        func(context);
        return true;
    }

    // Execute on main thread
    kr_log_to_reaper(LOG_DEBUG, "Dispatching to main thread...");

    __block bool completed = false;

    dispatch_sync(dispatch_get_main_queue(), ^{
        kr_log_to_reaper(LOG_DEBUG, "Now on main thread via dispatch");
        func(context);
        completed = true;
        kr_log_to_reaper(LOG_DEBUG, "Main thread execution completed");
    });

    kr_log_to_reaper(LOG_DEBUG, "After dispatch_sync");
    return completed;
}

// Show a keyring window, handling thread requirements - PUBLIC FUNCTION
bool kr_show_window(const char* title, bool key_exists, const char* message) {
    kr_log_to_reaper(LOG_INFO, "Entering kr_show_window");

    // Create context for thread-safe parameter passing
    KRContext ctx;
    ctx.title = title;
    ctx.key_exists = key_exists;
    ctx.message = message;
    ctx.success = false;

    // Execute on main thread
    bool executed = kr_execute_on_main_thread(kr_show_window_on_main_thread, &ctx);

    if (!executed) {
        kr_log_to_reaper(LOG_ERROR, "Failed to execute on main thread");
        return false;
    }

    kr_log_to_reaper(LOG_INFO, ctx.success ? "Keyring window displayed successfully" : "Keyring window display failed");
    return ctx.success;
}

// Update the keyring message - PUBLIC FUNCTION
bool kr_update_message(bool key_exists, const char* message) {
    kr_log_to_reaper(LOG_INFO, "Entering kr_update_message");

    // Create context for thread-safe parameter passing
    KRContext ctx;
    ctx.key_exists = key_exists;
    ctx.message = message;
    ctx.success = false;

    // Execute on main thread
    bool executed = kr_execute_on_main_thread(kr_update_message_on_main_thread, &ctx);

    if (!executed) {
        kr_log_to_reaper(LOG_ERROR, "Failed to execute on main thread");
        return false;
    }

    kr_log_to_reaper(LOG_INFO, ctx.success ? "Keyring message updated successfully" : "Keyring message update failed");
    return ctx.success;
}

// Close the keyring window if it exists - PUBLIC FUNCTION
void kr_close_window(void) {
    kr_log_to_reaper(LOG_DEBUG, "Entering kr_close_window");

    if (kr_window == nil) {
        kr_log_to_reaper(LOG_DEBUG, "No keyring window to close");
        return;
    }

    // Execute on main thread using a function pointer (not a block)
    kr_execute_on_main_thread(kr_close_window_on_main_thread, NULL);
}

// Check if keyring window exists - PUBLIC FUNCTION
bool kr_window_exists(void) {
    return (kr_window != nil);
}
