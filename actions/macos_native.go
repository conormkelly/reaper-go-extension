package actions

import (
	"fmt"
	"go-reaper/core"
	"go-reaper/reaper"
	"runtime"
	"time"
	"unsafe"
)

// This file implements a macOS UI that properly handles thread issues

/*
#cgo darwin CFLAGS: -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include "../reaper_ext_logging.h"

#import <Cocoa/Cocoa.h>

// Use our existing logging system
static void log_to_reaper(LogLevel level, const char* message) {
    log_message_v(level, "cocoaUI", message);
}

// Flag to track if we've initialized the main thread handler
static bool main_thread_initialized = false;

// Execute a function on the main thread synchronously
bool execute_on_main_thread(void (*func)(void*), void* context) {
    log_to_reaper(LOG_DEBUG, "Entering execute_on_main_thread");

    // Check if we're on the main thread already
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

// Global alert to keep it from being garbage collected
static NSAlert* g_alert = nil;

// Context struct for passing data to UI functions
typedef struct {
    const char* title;
    const char* message;
    bool success;
} AlertContext;

// Function to create and show an alert (to be executed on main thread)
void show_alert_on_main_thread(void* context) {
    log_to_reaper(LOG_DEBUG, "Entering show_alert_on_main_thread");

    AlertContext* ctx = (AlertContext*)context;
    if (!ctx) {
        log_to_reaper(LOG_ERROR, "Null context in show_alert_on_main_thread");
        return;
    }

    @autoreleasepool {
        @try {
            log_to_reaper(LOG_DEBUG, "Creating NSAlert...");

            // Create alert if not already created
            if (g_alert == nil) {
                g_alert = [[NSAlert alloc] init];
            }

            // Check if alloc/init succeeded
            if (g_alert == nil) {
                log_to_reaper(LOG_ERROR, "Failed to create NSAlert");
                ctx->success = false;
                return;
            }

            // Configure the alert
            log_to_reaper(LOG_DEBUG, "Setting alert properties...");
            [g_alert setMessageText:[NSString stringWithUTF8String:ctx->title]];
            [g_alert setInformativeText:[NSString stringWithUTF8String:ctx->message]];
            [g_alert addButtonWithTitle:@"OK"];

            // Show the alert
            log_to_reaper(LOG_DEBUG, "About to show alert modally...");
            [g_alert runModal];

            log_to_reaper(LOG_DEBUG, "Alert displayed and closed successfully");
            ctx->success = true;
        }
        @catch (NSException *exception) {
            log_to_reaper(LOG_ERROR, "EXCEPTION in show_alert_on_main_thread");
            NSLog(@"Exception: %@", exception);
            ctx->success = false;
        }
    }
}

// Main entry point for showing an alert - handles threading properly
bool show_ui_alert(const char* title, const char* message) {
    log_to_reaper(LOG_INFO, "Entering show_ui_alert");

    // Create context to pass to the main thread
    AlertContext ctx;
    ctx.title = title;
    ctx.message = message;
    ctx.success = false;

    // Execute on main thread
    bool executed = execute_on_main_thread(show_alert_on_main_thread, &ctx);

    if (!executed) {
        log_to_reaper(LOG_ERROR, "Failed to execute on main thread");
        return false;
    }

    log_to_reaper(LOG_INFO, ctx.success ? "Alert displayed successfully" : "Alert display failed");
    return ctx.success;
}

// Check if we're on the main thread
bool is_main_thread(void) {
    return [NSThread isMainThread];
}

// Get information about all running threads
const char* get_all_thread_info(void) {
    static char buffer[1024];
    NSArray *threads = [NSThread callStackSymbols];
    int count = (int)[threads count];

    sprintf(buffer, "Thread count: %d, Main thread: %s",
            count, [NSThread isMainThread] ? "YES" : "NO");

    return buffer;
}
*/
import "C"

// RegisterNativeDemo registers the native UI demo action
func RegisterNativeDemo() error {
	core.LogInfo("Registering Main Thread UI demo action")

	actionID, err := reaper.RegisterMainAction("GO_FYNE_DEMO", "Go: Main Thread UI Demo")
	if err != nil {
		core.LogError(fmt.Sprintf("Failed to register UI demo action: %v", err))
		return fmt.Errorf("failed to register UI demo action: %v", err)
	}

	logMsg := fmt.Sprintf("Main Thread UI demo action registered with ID: %d", actionID)
	reaper.ConsoleLog(logMsg)
	core.LogInfo(logMsg)

	reaper.SetActionHandler("GO_FYNE_DEMO", handleMainThreadUIDemo)
	return nil
}

// handleMainThreadUIDemo shows a native alert ensuring it runs on the main thread
func handleMainThreadUIDemo() {
	// Log action triggered
	reaper.ConsoleLog("Main Thread UI Demo action triggered!")
	core.LogInfo("Main Thread UI Demo action triggered!")

	// Ensure we're running on macOS
	if runtime.GOOS != "darwin" {
		reaper.MessageBox("This demo is currently only implemented for macOS", "Main Thread UI Demo")
		return
	}

	// Check thread information
	isMainThread := bool(C.is_main_thread())
	threadInfo := C.GoString(C.get_all_thread_info())

	reaper.ConsoleLog(fmt.Sprintf("Thread info: %s", threadInfo))
	reaper.ConsoleLog(fmt.Sprintf("Is main thread: %v", isMainThread))

	core.LogInfo(fmt.Sprintf("Thread info: %s", threadInfo))
	core.LogInfo(fmt.Sprintf("Is main thread: %v", isMainThread))

	// Prepare the alert content
	title := C.CString("REAPER Go Extension")
	defer C.free(unsafe.Pointer(title))

	message := C.CString("This is a native UI alert that ensures it runs on the main thread.")
	defer C.free(unsafe.Pointer(message))

	// Show the alert with proper thread handling
	reaper.ConsoleLog("About to show alert with proper thread handling...")
	core.LogInfo("About to show alert with proper thread handling...")

	// Display immediately
	result := C.show_ui_alert(title, message)

	if bool(result) {
		reaper.ConsoleLog("Alert displayed successfully")
		core.LogInfo("Alert displayed successfully")
	} else {
		reaper.ConsoleLog("Failed to display alert")
		core.LogError("Failed to display alert")
		reaper.MessageBox("Failed to display native alert. See log for details.", "Main Thread UI Demo")
	}

	// Keep extension running briefly to allow async operations to complete
	time.Sleep(100 * time.Millisecond)

	reaper.ConsoleLog("Main Thread UI Demo action handler completed")
	core.LogInfo("Main Thread UI Demo action handler completed")
}
