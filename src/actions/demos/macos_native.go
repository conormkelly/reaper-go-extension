package demo

import (
	"fmt"
	"go-reaper/src/pkg/logger"
	"go-reaper/src/reaper"
	"runtime"
	"time"
	"unsafe"
)

// This file implements a complete native macOS window with controls

/*
#cgo darwin CFLAGS: -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include "../../c/logging/logging.h"

#import <Cocoa/Cocoa.h>

// Use our core logging system
static void log_to_reaper(LogLevel level, const char* message) {
    log_message_v(level, "windowUI", message);
}

// Forward declaration of our controller class
@interface GoWindowController : NSObject
- (void)buttonClicked:(id)sender;
- (void)formSubmitted:(id)sender;
@end

// Global references to prevent garbage collection
static NSWindow* g_window = nil;
static GoWindowController* g_controller = nil;
static NSTextField* g_nameField = nil;
static NSTextField* g_descField = nil;
static NSTextField* g_resultsField = nil;

// Window context for passing data
typedef struct {
    const char* title;
    bool success;
} WindowContext;

// Implementation of our controller class
@implementation GoWindowController

- (void)buttonClicked:(id)sender {
    log_to_reaper(LOG_DEBUG, "Close button clicked");

    if (g_window != nil) {
        [g_window close];
        // Don't set g_window to nil here to allow checking if window was closed
    }
}

- (void)formSubmitted:(id)sender {
    log_to_reaper(LOG_DEBUG, "Form submitted");

    if (g_nameField == nil || g_descField == nil || g_resultsField == nil) {
        log_to_reaper(LOG_ERROR, "Fields not initialized");
        return;
    }

    NSString* name = [g_nameField stringValue];
    NSString* desc = [g_descField stringValue];

    NSString* results = [NSString stringWithFormat:@"Results:\n\nName: %@\nDescription: %@", name, desc];
    [g_resultsField setStringValue:results];
}

@end

// Function to create and show a window (to be executed on main thread)
void show_window_on_main_thread(void* context) {
    log_to_reaper(LOG_DEBUG, "Entering show_window_on_main_thread");

    WindowContext* ctx = (WindowContext*)context;
    if (!ctx) {
        log_to_reaper(LOG_ERROR, "Null context in show_window_on_main_thread");
        return;
    }

    @autoreleasepool {
        @try {
            // Check if window already exists
            if (g_window != nil) {
                log_to_reaper(LOG_DEBUG, "Window already exists, bringing to front");
                [g_window makeKeyAndOrderFront:nil];
                ctx->success = true;
                return;
            }

            log_to_reaper(LOG_DEBUG, "Creating controller...");

            // Create controller if needed
            if (g_controller == nil) {
                g_controller = [[GoWindowController alloc] init];
                if (!g_controller) {
                    log_to_reaper(LOG_ERROR, "Failed to create controller");
                    ctx->success = false;
                    return;
                }
            }

            log_to_reaper(LOG_DEBUG, "Creating window...");

            // Create the window
            NSRect frame = NSMakeRect(100, 100, 500, 400);
            NSWindow* window = [[NSWindow alloc]
                initWithContentRect:frame
                styleMask:NSWindowStyleMaskTitled|NSWindowStyleMaskClosable|NSWindowStyleMaskResizable
                backing:NSBackingStoreBuffered
                defer:NO];

            if (window == nil) {
                log_to_reaper(LOG_ERROR, "Failed to create window");
                ctx->success = false;
                return;
            }

            // Set window properties
            [window setTitle:[NSString stringWithUTF8String:ctx->title]];
            [window setReleasedWhenClosed:NO]; // Important: Don't release on close

            // Get the content view
            NSView* contentView = [window contentView];

            // Create a heading label
            NSTextField* headingLabel = [[NSTextField alloc] initWithFrame:NSMakeRect(20, 340, 460, 30)];
            [headingLabel setStringValue:@"REAPER Go Extension - Native UI Demo"];
            [headingLabel setFont:[NSFont boldSystemFontOfSize:18]];
            [headingLabel setBezeled:NO];
            [headingLabel setDrawsBackground:NO];
            [headingLabel setEditable:NO];
            [headingLabel setSelectable:NO];
            [contentView addSubview:headingLabel];

            // Create name label and field
            NSTextField* nameLabel = [[NSTextField alloc] initWithFrame:NSMakeRect(20, 300, 100, 24)];
            [nameLabel setStringValue:@"Name:"];
            [nameLabel setBezeled:NO];
            [nameLabel setDrawsBackground:NO];
            [nameLabel setEditable:NO];
            [nameLabel setSelectable:NO];
            [contentView addSubview:nameLabel];

            NSTextField* nameField = [[NSTextField alloc] initWithFrame:NSMakeRect(130, 300, 350, 24)];
            [nameField setPlaceholderString:@"Enter your name"];
            g_nameField = nameField; // Store reference
            [contentView addSubview:nameField];

            // Create description label and field
            NSTextField* descLabel = [[NSTextField alloc] initWithFrame:NSMakeRect(20, 260, 100, 24)];
            [descLabel setStringValue:@"Description:"];
            [descLabel setBezeled:NO];
            [descLabel setDrawsBackground:NO];
            [descLabel setEditable:NO];
            [descLabel setSelectable:NO];
            [contentView addSubview:descLabel];

            NSTextField* descField = [[NSTextField alloc] initWithFrame:NSMakeRect(130, 260, 350, 24)];
            [descField setPlaceholderString:@"Enter a description"];
            g_descField = descField; // Store reference
            [contentView addSubview:descField];

            // Create submit button
            NSButton* submitButton = [[NSButton alloc] initWithFrame:NSMakeRect(130, 220, 100, 32)];
            [submitButton setTitle:@"Submit"];
            [submitButton setBezelStyle:NSBezelStyleRounded];
            [submitButton setTarget:g_controller];
            [submitButton setAction:@selector(formSubmitted:)];
            [contentView addSubview:submitButton];

            // Create close button
            NSButton* closeButton = [[NSButton alloc] initWithFrame:NSMakeRect(240, 220, 100, 32)];
            [closeButton setTitle:@"Close"];
            [closeButton setBezelStyle:NSBezelStyleRounded];
            [closeButton setTarget:g_controller];
            [closeButton setAction:@selector(buttonClicked:)];
            [contentView addSubview:closeButton];

            // Create results area
            NSTextField* resultsField = [[NSTextField alloc] initWithFrame:NSMakeRect(20, 20, 460, 180)];
            [resultsField setStringValue:@"Results will appear here."];
            [resultsField setBezeled:YES];
            [resultsField setDrawsBackground:YES];
            [resultsField setEditable:NO];
            [resultsField setSelectable:YES];
            g_resultsField = resultsField; // Store reference
            [contentView addSubview:resultsField];

            // Store window reference
            g_window = window;

            // Show the window
            [window center];
            [window makeKeyAndOrderFront:nil];

            log_to_reaper(LOG_INFO, "Window created and shown successfully");
            ctx->success = true;
        }
        @catch (NSException *exception) {
            log_to_reaper(LOG_ERROR, "EXCEPTION creating window");
            NSLog(@"Exception: %@", exception);
            ctx->success = false;
        }
    }
}

// Function to close window on main thread
void close_window_on_main_thread(void* context) {
    log_to_reaper(LOG_DEBUG, "Entering close_window_on_main_thread");

    @try {
        if (g_window != nil) {
            [g_window close];
            g_window = nil;
            log_to_reaper(LOG_INFO, "Window closed successfully");
        } else {
            log_to_reaper(LOG_DEBUG, "Window already nil when trying to close");
        }
    }
    @catch (NSException *exception) {
        log_to_reaper(LOG_ERROR, "EXCEPTION closing window");
        NSLog(@"Exception: %@", exception);
    }
}

// Execute function on main thread
bool execute_on_main_thread(void (*func)(void*), void* context) {
    log_to_reaper(LOG_DEBUG, "Entering execute_on_main_thread");

    // Check if already on main thread
    if ([NSThread isMainThread]) {
        log_to_reaper(LOG_DEBUG, "Already on main thread, executing directly");
        func(context);
        return true;
    }

    // Execute on main thread
    log_to_reaper(LOG_DEBUG, "Dispatching to main thread...");

    __block bool completed = false;

    dispatch_sync(dispatch_get_main_queue(), ^{
        log_to_reaper(LOG_DEBUG, "Now on main thread via dispatch");
        func(context);
        completed = true;
        log_to_reaper(LOG_DEBUG, "Main thread execution completed");
    });

    log_to_reaper(LOG_DEBUG, "After dispatch_sync");
    return completed;
}

// Show a native window, handling thread requirements
bool show_native_window(const char* title) {
    log_to_reaper(LOG_INFO, "Entering show_native_window");

    // Create context for thread-safe parameter passing
    WindowContext ctx;
    ctx.title = title;
    ctx.success = false;

    // Execute on main thread
    bool executed = execute_on_main_thread(show_window_on_main_thread, &ctx);

    if (!executed) {
        log_to_reaper(LOG_ERROR, "Failed to execute on main thread");
        return false;
    }

    log_to_reaper(LOG_INFO, ctx.success ? "Window displayed successfully" : "Window display failed");
    return ctx.success;
}

// Close the window if it exists
void close_native_window(void) {
    log_to_reaper(LOG_DEBUG, "Entering close_native_window");

    if (g_window == nil) {
        log_to_reaper(LOG_DEBUG, "No window to close");
        return;
    }

    // Execute on main thread using a function pointer (not a block)
    execute_on_main_thread(close_window_on_main_thread, NULL);
}

// Check if window exists and is visible
bool is_window_visible(void) {
    return (g_window != nil);
}

// Check if we're on the main thread
bool is_main_thread(void) {
    return [NSThread isMainThread];
}

// Get the app's main thread state
const char* get_thread_state(void) {
    static char buffer[256];
    NSThread* mainThread = [NSThread mainThread];
    BOOL isMain = [NSThread isMainThread];

    sprintf(buffer, "isMainThread=%s, mainThread=%p, currentThread=%p",
            isMain ? "YES" : "NO",
            (void*)mainThread,
            (void*)[NSThread currentThread]);

    return buffer;
}
*/
import "C"

// RegisterNativeWindow registers the native window action
func RegisterNativeWindow() error {
	logger.Info("Registering Native Window action")

	actionID, err := reaper.RegisterMainAction("GO_NATIVE_WINDOW", "Go: Native Window Demo")
	if err != nil {
		logger.Error("Failed to register native window action: %v", err)
		return fmt.Errorf("failed to register native window action: %v", err)
	}

	logger.Info("Native Window action registered with ID: %d", actionID)

	reaper.SetActionHandler("GO_NATIVE_WINDOW", handleNativeWindow)
	return nil
}

// handleNativeWindow shows a native window with controls
func handleNativeWindow() {
	// Log action triggered
	logger.Info("Native Window action triggered!")

	// Ensure we're running on macOS
	if runtime.GOOS != "darwin" {
		reaper.MessageBox("This demo is currently only implemented for macOS", "Native Window Demo")
		return
	}

	// Check thread state
	isMainThread := bool(C.is_main_thread())
	threadState := C.GoString(C.get_thread_state())

	logger.Info("Thread state: %s", threadState)
	logger.Info("Is main thread: %v", isMainThread)

	// Show the native window
	title := C.CString("REAPER Go Extension")
	defer C.free(unsafe.Pointer(title))

	logger.Info("About to show native window...")

	// Show the window with proper thread handling
	result := C.show_native_window(title)

	if bool(result) {
		logger.Info("Window created/shown successfully")
	} else {
		logger.Error("Failed to create/show window")
		reaper.MessageBox("Failed to create/show native window. See log for details.", "Native Window Demo")
	}

	// Keep the action handler alive briefly to ensure UI operations complete
	time.Sleep(100 * time.Millisecond)

	logger.Info("Native Window action handler completed")
}

// CloseNativeWindow closes the native window if it exists
func CloseNativeWindow() {
	if runtime.GOOS == "darwin" {
		logger.Info("Closing native window...")

		C.close_native_window()

		logger.Info("Native window close request completed")
	}
}

// IsNativeWindowVisible checks if the native window is visible
func IsNativeWindowVisible() bool {
	if runtime.GOOS == "darwin" {
		return bool(C.is_window_visible())
	}
	return false
}
