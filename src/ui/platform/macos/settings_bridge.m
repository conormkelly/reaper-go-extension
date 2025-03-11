#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include "../../../c/logging/logging.h"
#import <Cocoa/Cocoa.h>
#include "settings_bridge.h"

// Use our core logging system
static void settings_log_to_reaper(LogLevel level, const char* message) {
    log_message_v(level, "settingsUI", message);
}

// Forward declaration of our controller class
@interface RPRSettingsWindowController : NSObject
- (void)cancelButtonClicked:(id)sender;
- (void)saveButtonClicked:(id)sender;
@end

// Global references for window and controls
static NSWindow* settings_window = nil;
static RPRSettingsWindowController* settings_controller = nil;
static NSTextField* settings_api_key_field = nil;
static NSTextField* settings_model_field = nil;
static NSTextField* settings_temperature_field = nil;
static NSButton* settings_save_button = nil;
static NSButton* settings_cancel_button = nil;

// Implementation of our controller class
@implementation RPRSettingsWindowController

- (void)cancelButtonClicked:(id)sender {
    settings_log_to_reaper(LOG_DEBUG, "Cancel button clicked");

    if (settings_window != nil) {
        [settings_window close];
        settings_window = nil;
    }
}

- (void)saveButtonClicked:(id)sender {
    settings_log_to_reaper(LOG_DEBUG, "Save button clicked");

    if (settings_api_key_field == nil || settings_model_field == nil || settings_temperature_field == nil) {
        settings_log_to_reaper(LOG_ERROR, "Input fields not initialized");
        return;
    }

    NSString* apiKey = [settings_api_key_field stringValue];
    NSString* model = [settings_model_field stringValue];
    NSString* tempStr = [settings_temperature_field stringValue];
    
    // Convert temperature string to double
    double temperature = 0.7; // Default
    if ([tempStr length] > 0) {
        temperature = [tempStr doubleValue];
        
        // Clamp temperature to valid range
        if (temperature < 0.0) {
            temperature = 0.0;
        } else if (temperature > 1.0) {
            temperature = 1.0;
        }
    }
    
    // Get the values as C strings
    const char* apiKeyStr = [apiKey UTF8String];
    const char* modelStr = [model UTF8String];
    
    // Create persistent copies
    char* apiKeyCopy = (char*)malloc(strlen(apiKeyStr) + 1);
    char* modelCopy = (char*)malloc(strlen(modelStr) + 1);
    
    if (apiKeyCopy && modelCopy) {
        strcpy(apiKeyCopy, apiKeyStr);
        strcpy(modelCopy, modelStr);
        
        // First close the window to prevent UI freeze
        settings_log_to_reaper(LOG_DEBUG, "About to close settings window");
        if (settings_window != nil) {
            [settings_window orderOut:nil]; // Hide window immediately
            [settings_window close];
            settings_window = nil;
        }
        settings_log_to_reaper(LOG_DEBUG, "Window closed");
        
        // Create a semaphore for synchronization
        dispatch_semaphore_t semaphore = dispatch_semaphore_create(0);
        
        // Use the main queue for better REAPER API compatibility
        settings_log_to_reaper(LOG_DEBUG, "Dispatching to main queue");
        dispatch_async(dispatch_get_main_queue(), ^{
            settings_log_to_reaper(LOG_DEBUG, "Processing settings on main queue - BEFORE go_process_settings");
            
            // Log timestamp before Go function call
            NSDate *beforeGoCall = [NSDate date];
            NSDateFormatter *formatter = [[NSDateFormatter alloc] init];
            [formatter setDateFormat:@"yyyy-MM-dd HH:mm:ss.SSS"];
            NSString *startTimeStr = [NSString stringWithFormat:@"Go function call starting at: %@", 
                                    [formatter stringFromDate:beforeGoCall]];
            settings_log_to_reaper(LOG_DEBUG, [startTimeStr UTF8String]);
            
            // Call the Go function
            go_process_settings(apiKeyCopy, modelCopy, temperature);
            
            // Log timestamp after Go function returns
            NSDate *afterGoCall = [NSDate date];
            NSTimeInterval elapsedTime = [afterGoCall timeIntervalSinceDate:beforeGoCall];
            
            NSString *finishTimeStr = [NSString stringWithFormat:@"Go function returned after %.3f seconds at: %@", 
                                     elapsedTime,
                                     [formatter stringFromDate:afterGoCall]];
            settings_log_to_reaper(LOG_DEBUG, [finishTimeStr UTF8String]);
            
            // Free memory allocated for strings
            free(apiKeyCopy);
            free(modelCopy);
            
            // Signal completion
            settings_log_to_reaper(LOG_DEBUG, "About to signal semaphore");
            dispatch_semaphore_signal(semaphore);
            
            settings_log_to_reaper(LOG_DEBUG, "Settings processing completed, semaphore signaled");
        });
        
        // Wait for a reasonable timeout to ensure processing completes
        // Use a background thread for waiting to avoid blocking UI
        settings_log_to_reaper(LOG_DEBUG, "Starting background thread for semaphore wait");
        dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
            settings_log_to_reaper(LOG_DEBUG, "Waiting for settings processing to complete");
            
            // Wait up to 5 seconds for processing to complete
            NSDate *waitStart = [NSDate date];
            long result = dispatch_semaphore_wait(semaphore, dispatch_time(DISPATCH_TIME_NOW, 5 * NSEC_PER_SEC));
            NSTimeInterval waitTime = [[NSDate date] timeIntervalSinceDate:waitStart];
            
            NSString *waitResultStr;
            if (result == 0) {
                waitResultStr = [NSString stringWithFormat:@"Settings processing completed successfully (waited %.3f seconds)", waitTime];
            } else {
                waitResultStr = [NSString stringWithFormat:@"Settings processing timed out after %.3f seconds", waitTime];
            }
            settings_log_to_reaper(result == 0 ? LOG_DEBUG : LOG_WARNING, [waitResultStr UTF8String]);
        });
    } else {
        settings_log_to_reaper(LOG_ERROR, "Failed to allocate memory for settings");
        if (apiKeyCopy) free(apiKeyCopy);
        if (modelCopy) free(modelCopy);
        
        // Show error without closing window
        NSAlert *alert = [[NSAlert alloc] init];
        [alert setMessageText:@"Memory Error"];
        [alert setInformativeText:@"Failed to allocate memory for settings"];
        [alert addButtonWithTitle:@"OK"];
        [alert runModal];
    }
}

@end

// Function to create and show the settings window - internal
static void settings_show_window_on_main_thread(void* context) {
    settings_log_to_reaper(LOG_DEBUG, "Entering settings_show_window_on_main_thread");

    SettingsContext* ctx = (SettingsContext*)context;
    if (!ctx) {
        settings_log_to_reaper(LOG_ERROR, "Null context in settings_show_window_on_main_thread");
        return;
    }

    @autoreleasepool {
        @try {
            // Check if window already exists
            if (settings_window != nil) {
                settings_log_to_reaper(LOG_DEBUG, "Window already exists, bringing to front");
                [settings_window makeKeyAndOrderFront:nil];
                
                // Update field values if window exists
                if (settings_api_key_field != nil && ctx->api_key != NULL) {
                    [settings_api_key_field setStringValue:[NSString stringWithUTF8String:ctx->api_key]];
                }
                
                if (settings_model_field != nil && ctx->model != NULL) {
                    [settings_model_field setStringValue:[NSString stringWithUTF8String:ctx->model]];
                }
                
                if (settings_temperature_field != nil) {
                    [settings_temperature_field setStringValue:[NSString stringWithFormat:@"%.1f", ctx->temperature]];
                }
                
                ctx->success = true;
                return;
            }

            // Create controller if needed
            if (settings_controller == nil) {
                settings_controller = [[RPRSettingsWindowController alloc] init];
                if (!settings_controller) {
                    settings_log_to_reaper(LOG_ERROR, "Failed to create controller");
                    ctx->success = false;
                    return;
                }
            }

            // Create the window
            NSRect frame = NSMakeRect(100, 100, 450, 240);
            NSWindow* window = [[NSWindow alloc]
                initWithContentRect:frame
                styleMask:NSWindowStyleMaskTitled|NSWindowStyleMaskClosable
                backing:NSBackingStoreBuffered
                defer:NO];

            if (window == nil) {
                settings_log_to_reaper(LOG_ERROR, "Failed to create window");
                ctx->success = false;
                return;
            }

            // Set window properties
            [window setTitle:[NSString stringWithUTF8String:ctx->title]];
            [window setReleasedWhenClosed:NO]; // Important: Don't release on close

            // Get the content view
            NSView* contentView = [window contentView];

            // Create API Key label and field
            NSTextField* apiKeyLabel = [[NSTextField alloc] initWithFrame:NSMakeRect(20, 190, 120, 24)];
            [apiKeyLabel setStringValue:@"OpenAI API Key:"];
            [apiKeyLabel setBezeled:NO];
            [apiKeyLabel setDrawsBackground:NO];
            [apiKeyLabel setEditable:NO];
            [apiKeyLabel setSelectable:NO];
            [contentView addSubview:apiKeyLabel];

            NSTextField* apiKeyField = [[NSTextField alloc] initWithFrame:NSMakeRect(150, 190, 280, 24)];
            [apiKeyField setPlaceholderString:@"Enter your OpenAI API key"];
            if (ctx->api_key != NULL) {
                [apiKeyField setStringValue:[NSString stringWithUTF8String:ctx->api_key]];
            }
            settings_api_key_field = apiKeyField;
            [contentView addSubview:apiKeyField];

            // Create Model label and field
            NSTextField* modelLabel = [[NSTextField alloc] initWithFrame:NSMakeRect(20, 150, 120, 24)];
            [modelLabel setStringValue:@"Model:"];
            [modelLabel setBezeled:NO];
            [modelLabel setDrawsBackground:NO];
            [modelLabel setEditable:NO];
            [modelLabel setSelectable:NO];
            [contentView addSubview:modelLabel];

            NSTextField* modelField = [[NSTextField alloc] initWithFrame:NSMakeRect(150, 150, 280, 24)];
            [modelField setPlaceholderString:@"e.g., gpt-3.5-turbo, gpt-4"];
            if (ctx->model != NULL) {
                [modelField setStringValue:[NSString stringWithUTF8String:ctx->model]];
            }
            settings_model_field = modelField;
            [contentView addSubview:modelField];

            // Create Temperature label and field
            NSTextField* tempLabel = [[NSTextField alloc] initWithFrame:NSMakeRect(20, 110, 120, 24)];
            [tempLabel setStringValue:@"Temperature:"];
            [tempLabel setBezeled:NO];
            [tempLabel setDrawsBackground:NO];
            [tempLabel setEditable:NO];
            [tempLabel setSelectable:NO];
            [contentView addSubview:tempLabel];

            NSTextField* tempField = [[NSTextField alloc] initWithFrame:NSMakeRect(150, 110, 280, 24)];
            [tempField setPlaceholderString:@"0.0-1.0 (higher = more creative)"];
            [tempField setStringValue:[NSString stringWithFormat:@"%.1f", ctx->temperature]];
            settings_temperature_field = tempField;
            [contentView addSubview:tempField];
            
            // Create info label
            NSTextField* infoLabel = [[NSTextField alloc] initWithFrame:NSMakeRect(20, 70, 410, 30)];

            // Create save button
            NSButton* saveButton = [[NSButton alloc] initWithFrame:NSMakeRect(230, 5, 100, 32)];
            [saveButton setTitle:@"Save"];
            [saveButton setBezelStyle:NSBezelStyleRounded];
            [saveButton setTarget:settings_controller];
            [saveButton setAction:@selector(saveButtonClicked:)];
            [saveButton setKeyEquivalent:@"\r"]; // Enter key
            settings_save_button = saveButton;
            [contentView addSubview:saveButton];

            // Create cancel button
            NSButton* cancelButton = [[NSButton alloc] initWithFrame:NSMakeRect(340, 5, 100, 32)];
            [cancelButton setTitle:@"Cancel"];
            [cancelButton setBezelStyle:NSBezelStyleRounded];
            [cancelButton setTarget:settings_controller];
            [cancelButton setAction:@selector(cancelButtonClicked:)];
            [cancelButton setKeyEquivalent:@"\033"]; // Escape key
            settings_cancel_button = cancelButton;
            [contentView addSubview:cancelButton];

            // Store window reference
            settings_window = window;

            // Show the window
            [window center];
            [window makeKeyAndOrderFront:nil];

            settings_log_to_reaper(LOG_INFO, "Settings window created and shown successfully");
            ctx->success = true;
        }
        @catch (NSException *exception) {
            settings_log_to_reaper(LOG_ERROR, "EXCEPTION creating settings window");
            NSLog(@"Exception: %@", exception);
            ctx->success = false;
        }
    }
}

// Function to close window on main thread - internal
static void settings_close_window_on_main_thread(void* context) {
    settings_log_to_reaper(LOG_DEBUG, "Entering settings_close_window_on_main_thread");

    @try {
        if (settings_window != nil) {
            [settings_window close];
            settings_window = nil;
            settings_api_key_field = nil;
            settings_model_field = nil;
            settings_temperature_field = nil;
            settings_save_button = nil;
            settings_cancel_button = nil;
            settings_log_to_reaper(LOG_INFO, "Settings window closed successfully");
        } else {
            settings_log_to_reaper(LOG_DEBUG, "Settings window already nil when trying to close");
        }
    }
    @catch (NSException *exception) {
        settings_log_to_reaper(LOG_ERROR, "EXCEPTION closing settings window");
        NSLog(@"Exception: %@", exception);
    }
}

// Execute function on main thread - internal
static bool settings_execute_on_main_thread(void (*func)(void*), void* context) {
    settings_log_to_reaper(LOG_DEBUG, "Entering settings_execute_on_main_thread");

    // Check if already on main thread
    if ([NSThread isMainThread]) {
        settings_log_to_reaper(LOG_DEBUG, "Already on main thread, executing directly");
        func(context);
        return true;
    }

    // Execute on main thread
    settings_log_to_reaper(LOG_DEBUG, "Dispatching to main thread...");

    __block bool completed = false;

    dispatch_sync(dispatch_get_main_queue(), ^{
        settings_log_to_reaper(LOG_DEBUG, "Now on main thread via dispatch");
        func(context);
        completed = true;
        settings_log_to_reaper(LOG_DEBUG, "Main thread execution completed");
    });

    settings_log_to_reaper(LOG_DEBUG, "After dispatch_sync");
    return completed;
}

// Show settings window, handling thread requirements - PUBLIC FUNCTION
bool settings_show_window(const char* title, const char* api_key, const char* model, double temperature) {
    settings_log_to_reaper(LOG_INFO, "Entering settings_show_window");

    // Create context for thread-safe parameter passing
    SettingsContext ctx;
    ctx.title = title;
    ctx.api_key = api_key;
    ctx.model = model;
    ctx.temperature = temperature;
    ctx.success = false;

    // Execute on main thread
    bool executed = settings_execute_on_main_thread(settings_show_window_on_main_thread, &ctx);

    if (!executed) {
        settings_log_to_reaper(LOG_ERROR, "Failed to execute on main thread");
        return false;
    }

    settings_log_to_reaper(LOG_INFO, ctx.success ? "Settings window displayed successfully" : "Settings window display failed");
    return ctx.success;
}

// Close the settings window if it exists - PUBLIC FUNCTION
void settings_close_window(void) {
    settings_log_to_reaper(LOG_DEBUG, "Entering settings_close_window");

    if (settings_window == nil) {
        settings_log_to_reaper(LOG_DEBUG, "No settings window to close");
        return;
    }

    // Execute on main thread
    settings_execute_on_main_thread(settings_close_window_on_main_thread, NULL);
}

// Check if settings window exists - PUBLIC FUNCTION
bool settings_window_exists(void) {
    return (settings_window != nil);
}
