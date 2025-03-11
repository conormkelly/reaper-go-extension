// src/ui/platform/macos/ui.m - macOS UI implementation
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include "../../../c/logging/logging.h"
#import <Cocoa/Cocoa.h>

// Use our core logging system
static void ui_log_to_reaper(LogLevel level, const char* message) {
    log_message_v(level, "macosUI", message);
}

// Forward declaration of controller classes
@interface RPRWindowController : NSObject
- (void)buttonClicked:(id)sender;
- (void)sliderChanged:(id)sender;
@end

@interface RPRParamViewController : NSObject
- (void)sliderChanged:(id)sender;
@end

// Global references
static NSMutableDictionary* g_window_controllers = nil;
static NSMutableDictionary* g_param_controllers = nil;

// Button callback function type
typedef void (*ButtonCallback)(void* sender);

// Slider callback function type
typedef void (*SliderCallback)(void* sender, double value);

// Implementation of window controller
@implementation RPRWindowController {
    ButtonCallback buttonCallback;
    void* buttonContext;
}

- (id)init {
    self = [super init];
    if (self) {
        buttonCallback = NULL;
        buttonContext = NULL;
    }
    return self;
}

- (void)setButtonCallback:(ButtonCallback)callback context:(void*)context {
    buttonCallback = callback;
    buttonContext = context;
}

- (void)buttonClicked:(id)sender {
    ui_log_to_reaper(LOG_DEBUG, "Button clicked");
    
    if (buttonCallback != NULL) {
        buttonCallback(buttonContext);
    }
}

- (void)sliderChanged:(id)sender {
    ui_log_to_reaper(LOG_DEBUG, "Slider changed");
    
    // Get the slider value
    NSSlider* slider = (NSSlider*)sender;
    double value = [slider doubleValue];
    
    // This would call the registered slider callback
    // sliderCallback(sliderContext, value);
}

@end

// Implementation of parameter view controller
@implementation RPRParamViewController {
    NSSlider* slider;
    NSTextField* valueLabel;
    NSTextField* nameLabel;
    NSTextField* explanationLabel;
    NSTextField* originalLabel;
    SliderCallback sliderCallback;
    void* sliderContext;
}

- (id)initWithSlider:(NSSlider*)theSlider valueLabel:(NSTextField*)theValueLabel 
                nameLabel:(NSTextField*)theNameLabel explanationLabel:(NSTextField*)theExplanationLabel
                originalLabel:(NSTextField*)theOriginalLabel {
    self = [super init];
    if (self) {
        slider = theSlider;
        valueLabel = theValueLabel;
        nameLabel = theNameLabel;
        explanationLabel = theExplanationLabel;
        originalLabel = theOriginalLabel;
        sliderCallback = NULL;
        sliderContext = NULL;
    }
    return self;
}

- (void)setSliderCallback:(SliderCallback)callback context:(void*)context {
    sliderCallback = callback;
    sliderContext = context;
}

- (void)sliderChanged:(id)sender {
    ui_log_to_reaper(LOG_DEBUG, "Parameter slider changed");
    
    double value = [slider doubleValue];
    
    if (sliderCallback != NULL) {
        sliderCallback(sliderContext, value);
    }
}

@end

// Initialize global dictionaries
static void init_global_dictionaries() {
    if (g_window_controllers == nil) {
        g_window_controllers = [[NSMutableDictionary alloc] init];
    }
    
    if (g_param_controllers == nil) {
        g_param_controllers = [[NSMutableDictionary alloc] init];
    }
}

// Ensure function is executed on the main thread
bool macos_run_on_main_thread(void (*func)(void*), void* context) {
    if (func == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null function pointer in macos_run_on_main_thread");
        return false;
    }
    
    // If already on main thread, run directly
    if ([NSThread isMainThread]) {
        func(context);
        return true;
    }
    
    // Execute on main thread
    dispatch_sync(dispatch_get_main_queue(), ^{
        func(context);
    });
    
    return true;
}

// Check if current thread is the main thread
bool macos_is_main_thread() {
    return [NSThread isMainThread];
}

// Create a window
void* macos_create_window(const char* title, int x, int y, int width, int height, bool resizable) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Window creation must be done on the main thread");
        return NULL;
    }
    
    init_global_dictionaries();
    
    @autoreleasepool {
        // Create style mask based on resizable flag
        NSWindowStyleMask styleMask = NSWindowStyleMaskTitled | NSWindowStyleMaskClosable;
        if (resizable) {
            styleMask |= NSWindowStyleMaskResizable;
        }
        
        // Create window
        NSRect frame = NSMakeRect(x, y, width, height);
        NSWindow* window = [[NSWindow alloc] initWithContentRect:frame
                                                       styleMask:styleMask
                                                         backing:NSBackingStoreBuffered
                                                           defer:NO];
        
        if (window == nil) {
            ui_log_to_reaper(LOG_ERROR, "Failed to create window");
            return NULL;
        }
        
        // Set title
        [window setTitle:[NSString stringWithUTF8String:title ? title : "Window"]];
        
        // Create window controller
        RPRWindowController* controller = [[RPRWindowController alloc] init];
        
        // Store controller in dictionary
        [g_window_controllers setObject:controller forKey:[NSValue valueWithPointer:window]];
        
        ui_log_to_reaper(LOG_INFO, "Window created successfully");
        
        // In non-ARC, we need to retain the window manually
        [window retain];
        return window;
    }
}

// Close a window
bool macos_close_window(void* window) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Window closing must be done on the main thread");
        return false;
    }
    
    if (window == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null window pointer in macos_close_window");
        return false;
    }
    
    @autoreleasepool {
        NSWindow* nsWindow = (NSWindow*)window;
        
        // Remove controller from dictionary
        [g_window_controllers removeObjectForKey:[NSValue valueWithPointer:nsWindow]];
        
        // Close window
        [nsWindow close];
        
        // Release the window
        [nsWindow release];
        
        ui_log_to_reaper(LOG_INFO, "Window closed successfully");
        return true;
    }
}

// Show a window
bool macos_show_window(void* window) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Window showing must be done on the main thread");
        return false;
    }
    
    if (window == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null window pointer in macos_show_window");
        return false;
    }
    
    @autoreleasepool {
        NSWindow* nsWindow = (NSWindow*)window;
        
        // Show window
        [nsWindow makeKeyAndOrderFront:nil];
        [nsWindow center]; // Center on screen
        
        ui_log_to_reaper(LOG_INFO, "Window shown successfully");
        return true;
    }
}

// Hide a window
bool macos_hide_window(void* window) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Window hiding must be done on the main thread");
        return false;
    }
    
    if (window == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null window pointer in macos_hide_window");
        return false;
    }
    
    @autoreleasepool {
        NSWindow* nsWindow = (NSWindow*)window;
        
        // Hide window
        [nsWindow orderOut:nil];
        
        ui_log_to_reaper(LOG_INFO, "Window hidden successfully");
        return true;
    }
}

// Check if a window is visible
bool macos_window_is_visible(void* window) {
    if (window == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null window pointer in macos_window_is_visible");
        return false;
    }
    
    @autoreleasepool {
        NSWindow* nsWindow = (NSWindow*)window;
        
        // Check if window is visible
        return [nsWindow isVisible];
    }
}

// Set window title
bool macos_set_window_title(void* window, const char* title) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Setting window title must be done on the main thread");
        return false;
    }
    
    if (window == NULL || title == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null pointer in macos_set_window_title");
        return false;
    }
    
    @autoreleasepool {
        NSWindow* nsWindow = (NSWindow*)window;
        
        // Set title
        [nsWindow setTitle:[NSString stringWithUTF8String:title]];
        
        return true;
    }
}

// Add a label to a window
void* macos_add_label(void* window, const char* text, int x, int y, int width, int height, bool bold, double size) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Adding label must be done on the main thread");
        return NULL;
    }
    
    if (window == NULL || text == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null pointer in macos_add_label");
        return NULL;
    }
    
    @autoreleasepool {
        NSWindow* nsWindow = (NSWindow*)window;
        NSView* contentView = [nsWindow contentView];
        
        // Create label
        NSTextField* label = [[NSTextField alloc] initWithFrame:NSMakeRect(x, y, width, height)];
        [label setStringValue:[NSString stringWithUTF8String:text]];
        [label setBezeled:NO];
        [label setDrawsBackground:NO];
        [label setEditable:NO];
        [label setSelectable:NO];
        
        // Set font
        NSFont* font = nil;
        if (bold) {
            font = [NSFont boldSystemFontOfSize:size];
        } else {
            font = [NSFont systemFontOfSize:size];
        }
        [label setFont:font];
        
        // Add to window
        [contentView addSubview:label];
        
        // Retain and return
        return label;
    }
}

// Add a button to a window
void* macos_add_button(void* window, const char* text, int x, int y, int width, int height) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Adding button must be done on the main thread");
        return NULL;
    }
    
    if (window == NULL || text == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null pointer in macos_add_button");
        return NULL;
    }
    
    @autoreleasepool {
        NSWindow* nsWindow = (NSWindow*)window;
        NSView* contentView = [nsWindow contentView];
        
        // Get window controller
        RPRWindowController* controller = [g_window_controllers objectForKey:[NSValue valueWithPointer:nsWindow]];
        if (controller == nil) {
            ui_log_to_reaper(LOG_ERROR, "Window controller not found");
            return NULL;
        }
        
        // Create button
        NSButton* button = [[NSButton alloc] initWithFrame:NSMakeRect(x, y, width, height)];
        [button setTitle:[NSString stringWithUTF8String:text]];
        [button setBezelStyle:NSBezelStyleRounded];
        [button setTarget:controller];
        [button setAction:@selector(buttonClicked:)];
        
        // Add to window
        [contentView addSubview:button];
        
        // Retain and return
        return button;
    }
}

// Set button callback
bool macos_set_button_callback(void* button, ButtonCallback callback) {
    if (button == NULL || callback == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null pointer in macos_set_button_callback");
        return false;
    }
    
    @autoreleasepool {
        NSButton* nsButton = (NSButton*)button;
        
        // Get the window
        NSWindow* window = [nsButton window];
        if (window == nil) {
            ui_log_to_reaper(LOG_ERROR, "Button is not in a window");
            return false;
        }
        
        // Get window controller
        RPRWindowController* controller = [g_window_controllers objectForKey:[NSValue valueWithPointer:window]];
        if (controller == nil) {
            ui_log_to_reaper(LOG_ERROR, "Window controller not found");
            return false;
        }
        
        // Set callback
        [controller setButtonCallback:callback context:button];
        
        return true;
    }
}

// Add a slider to a window
void* macos_add_slider(void* window, double min, double max, double value, int x, int y, int width, int height) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Adding slider must be done on the main thread");
        return NULL;
    }
    
    if (window == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null window pointer in macos_add_slider");
        return NULL;
    }
    
    @autoreleasepool {
        NSWindow* nsWindow = (NSWindow*)window;
        NSView* contentView = [nsWindow contentView];
        
        // Get window controller
        RPRWindowController* controller = [g_window_controllers objectForKey:[NSValue valueWithPointer:nsWindow]];
        if (controller == nil) {
            ui_log_to_reaper(LOG_ERROR, "Window controller not found");
            return NULL;
        }
        
        // Create slider
        NSSlider* slider = [[NSSlider alloc] initWithFrame:NSMakeRect(x, y, width, height)];
        [slider setMinValue:min];
        [slider setMaxValue:max];
        [slider setDoubleValue:value];
        [slider setContinuous:YES];
        [slider setTarget:controller];
        [slider setAction:@selector(sliderChanged:)];
        
        // Add to window
        [contentView addSubview:slider];
        
        // Retain and return
        return slider;
    }
}

// Add a text field to a window
void* macos_add_text_field(void* window, const char* placeholder, int x, int y, int width, int height) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Adding text field must be done on the main thread");
        return NULL;
    }
    
    if (window == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null window pointer in macos_add_text_field");
        return NULL;
    }
    
    @autoreleasepool {
        NSWindow* nsWindow = (NSWindow*)window;
        NSView* contentView = [nsWindow contentView];
        
        // Create text field
        NSTextField* textField = [[NSTextField alloc] initWithFrame:NSMakeRect(x, y, width, height)];
        if (placeholder != NULL) {
            [textField setPlaceholderString:[NSString stringWithUTF8String:placeholder]];
        }
        
        // Add to window
        [contentView addSubview:textField];
        
        // Retain and return
        return textField;
    }
}

// Show an alert dialog
int macos_show_alert(const char* title, const char* message, int style) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Showing alert must be done on the main thread");
        return 0;
    }
    
    if (title == NULL || message == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null pointer in macos_show_alert");
        return 0;
    }
    
    @autoreleasepool {
        NSAlert* alert = [[NSAlert alloc] init];
        [alert setMessageText:[NSString stringWithUTF8String:title]];
        [alert setInformativeText:[NSString stringWithUTF8String:message]];
        
        // Set buttons based on style
        switch (style) {
            case 0: // OK
                [alert addButtonWithTitle:@"OK"];
                break;
            case 1: // Yes/No
                [alert addButtonWithTitle:@"Yes"];
                [alert addButtonWithTitle:@"No"];
                break;
            case 2: // OK/Cancel
                [alert addButtonWithTitle:@"OK"];
                [alert addButtonWithTitle:@"Cancel"];
                break;
            default:
                [alert addButtonWithTitle:@"OK"];
                break;
        }
        
        // Show alert
        NSModalResponse response = [alert runModal];
        
        // Clean up
        [alert release];
        
        // Map response to return value
        if (style == 1) { // Yes/No
            return (response == NSAlertFirstButtonReturn) ? 1 : 0;
        } else if (style == 2) { // OK/Cancel
            return (response == NSAlertFirstButtonReturn) ? 1 : 0;
        } else {
            return 1; // OK
        }
    }
}

// Show a dialog with input fields
bool macos_get_user_inputs(const char* title, int num_inputs, const char* captions, char* values, int values_sz) {
    // This would implement the input dialog functionality
    // For now, it's a stub
    return false;
}

// Create a parameter view
void* macos_create_param_view(void* window, int x, int y, int width, int height, const char* name, 
                             double value, double min, double max, const char* formatted) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Creating parameter view must be done on the main thread");
        return NULL;
    }
    
    if (window == NULL || name == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null pointer in macos_create_param_view");
        return NULL;
    }
    
    @autoreleasepool {
        NSWindow* nsWindow = (NSWindow*)window;
        NSView* contentView = [nsWindow contentView];
        
        // Create a container view
        NSView* containerView = [[NSView alloc] initWithFrame:NSMakeRect(x, y, width, height)];
        
        // Create name label
        NSTextField* nameLabel = [[NSTextField alloc] initWithFrame:NSMakeRect(0, height-20, width/2, 20)];
        [nameLabel setStringValue:[NSString stringWithUTF8String:name]];
        [nameLabel setBezeled:NO];
        [nameLabel setDrawsBackground:NO];
        [nameLabel setEditable:NO];
        [nameLabel setSelectable:NO];
        [nameLabel setFont:[NSFont boldSystemFontOfSize:12]];
        [containerView addSubview:nameLabel];
        
        // Create value label
        NSTextField* valueLabel = [[NSTextField alloc] initWithFrame:NSMakeRect(width/2, height-20, width/2, 20)];
        [valueLabel setStringValue:formatted ? [NSString stringWithUTF8String:formatted] : @""];
        [valueLabel setBezeled:NO];
        [valueLabel setDrawsBackground:NO];
        [valueLabel setEditable:NO];
        [valueLabel setSelectable:NO];
        [valueLabel setAlignment:NSTextAlignmentRight];
        [containerView addSubview:valueLabel];
        
        // Create slider
        NSSlider* slider = [[NSSlider alloc] initWithFrame:NSMakeRect(0, height-45, width, 25)];
        [slider setMinValue:min];
        [slider setMaxValue:max];
        [slider setDoubleValue:value];
        [slider setContinuous:YES];
        [containerView addSubview:slider];
        
        // Create original value label
        NSTextField* originalLabel = [[NSTextField alloc] initWithFrame:NSMakeRect(0, height-70, width, 20)];
        [originalLabel setStringValue:@"Original: (not set)"];
        [originalLabel setBezeled:NO];
        [originalLabel setDrawsBackground:NO];
        [originalLabel setEditable:NO];
        [originalLabel setSelectable:NO];
        [originalLabel setFont:[NSFont systemFontOfSize:10]];
        [containerView addSubview:originalLabel];
        
        // Create explanation label
        NSTextField* explanationLabel = [[NSTextField alloc] initWithFrame:NSMakeRect(0, 0, width, height-70)];
        [explanationLabel setStringValue:@""];
        [explanationLabel setBezeled:NO];
        [explanationLabel setDrawsBackground:NO];
        [explanationLabel setEditable:NO];
        [explanationLabel setSelectable:YES];
        [explanationLabel setFont:[NSFont systemFontOfSize:11]];
        [explanationLabel setTextColor:[NSColor darkGrayColor]];
        [containerView addSubview:explanationLabel];
        
        // Create controller
        RPRParamViewController* controller = [[RPRParamViewController alloc] initWithSlider:slider
                                                                              valueLabel:valueLabel
                                                                               nameLabel:nameLabel
                                                                         explanationLabel:explanationLabel
                                                                           originalLabel:originalLabel];
        
        // Set slider target
        [slider setTarget:controller];
        [slider setAction:@selector(sliderChanged:)];
        
        // Store controller
        [g_param_controllers setObject:controller forKey:[NSValue valueWithPointer:containerView]];
        
        // Add to window
        [contentView addSubview:containerView];
        
        // Retain and return
        return containerView;
    }
}

// Set parameter value
bool macos_param_view_set_value(void* view, double value) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Setting parameter value must be done on the main thread");
        return false;
    }
    
    if (view == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null view pointer in macos_param_view_set_value");
        return false;
    }
    
    @autoreleasepool {
        NSView* containerView = (NSView*)view;
        
        // Find controller
        RPRParamViewController* controller = [g_param_controllers objectForKey:[NSValue valueWithPointer:containerView]];
        if (controller == nil) {
            ui_log_to_reaper(LOG_ERROR, "Parameter controller not found");
            return false;
        }
        
        // Find slider
        NSSlider* slider = nil;
        for (NSView* subview in [containerView subviews]) {
            if ([subview isKindOfClass:[NSSlider class]]) {
                slider = (NSSlider*)subview;
                break;
            }
        }
        
        if (slider == nil) {
            ui_log_to_reaper(LOG_ERROR, "Slider not found in parameter view");
            return false;
        }
        
        // Set value
        [slider setDoubleValue:value];
        
        return true;
    }
}

// Set formatted value
bool macos_param_view_set_formatted(void* view, const char* formatted) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Setting formatted value must be done on the main thread");
        return false;
    }
    
    if (view == NULL || formatted == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null pointer in macos_param_view_set_formatted");
        return false;
    }
    
    @autoreleasepool {
        NSView* containerView = (NSView*)view;
        
        // Find value label (second text field)
        NSTextField* valueLabel = nil;
        int textFieldCount = 0;
        for (NSView* subview in [containerView subviews]) {
            if ([subview isKindOfClass:[NSTextField class]]) {
                textFieldCount++;
                if (textFieldCount == 2) {
                    valueLabel = (NSTextField*)subview;
                    break;
                }
            }
        }
        
        if (valueLabel == nil) {
            ui_log_to_reaper(LOG_ERROR, "Value label not found in parameter view");
            return false;
        }
        
        // Set formatted value
        [valueLabel setStringValue:[NSString stringWithUTF8String:formatted]];
        
        return true;
    }
}

// Set explanation
bool macos_param_view_set_explanation(void* view, const char* explanation) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Setting explanation must be done on the main thread");
        return false;
    }
    
    if (view == NULL || explanation == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null pointer in macos_param_view_set_explanation");
        return false;
    }
    
    @autoreleasepool {
        NSView* containerView = (NSView*)view;
        
        // Find explanation label (last text field)
        NSTextField* explanationLabel = nil;
        for (NSView* subview in [containerView subviews]) {
            if ([subview isKindOfClass:[NSTextField class]]) {
                explanationLabel = (NSTextField*)subview;
            }
        }
        
        if (explanationLabel == nil) {
            ui_log_to_reaper(LOG_ERROR, "Explanation label not found in parameter view");
            return false;
        }
        
        // Set explanation
        [explanationLabel setStringValue:[NSString stringWithUTF8String:explanation]];
        
        return true;
    }
}

// Set original value
bool macos_param_view_set_original(void* view, double original, const char* formatted) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Setting original value must be done on the main thread");
        return false;
    }
    
    if (view == NULL || formatted == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null pointer in macos_param_view_set_original");
        return false;
    }
    
    @autoreleasepool {
        NSView* containerView = (NSView*)view;
        
        // Find original label (fourth text field)
        NSTextField* originalLabel = nil;
        int textFieldCount = 0;
        for (NSView* subview in [containerView subviews]) {
            if ([subview isKindOfClass:[NSTextField class]]) {
                textFieldCount++;
                if (textFieldCount == 4) {
                    originalLabel = (NSTextField*)subview;
                    break;
                }
            }
        }
        
        if (originalLabel == nil) {
            ui_log_to_reaper(LOG_ERROR, "Original label not found in parameter view");
            return false;
        }
        
        // Set original value
        NSString* text = [NSString stringWithFormat:@"Original: %s", formatted];
        [originalLabel setStringValue:text];
        
        return true;
    }
}

// Show parameter view
bool macos_param_view_show(void* view) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Showing parameter view must be done on the main thread");
        return false;
    }
    
    if (view == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null view pointer in macos_param_view_show");
        return false;
    }
    
    @autoreleasepool {
        NSView* containerView = (NSView*)view;
        
        // Show view
        [containerView setHidden:NO];
        
        return true;
    }
}

// Hide parameter view
bool macos_param_view_hide(void* view) {
    if (!macos_is_main_thread()) {
        ui_log_to_reaper(LOG_ERROR, "Hiding parameter view must be done on the main thread");
        return false;
    }
    
    if (view == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null view pointer in macos_param_view_hide");
        return false;
    }
    
    @autoreleasepool {
        NSView* containerView = (NSView*)view;
        
        // Hide view
        [containerView setHidden:YES];
        
        return true;
    }
}

// Get parameter value
double macos_param_view_get_value(void* view) {
    if (view == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null view pointer in macos_param_view_get_value");
        return 0.0;
    }
    
    @autoreleasepool {
        NSView* containerView = (NSView*)view;
        
        // Find slider
        NSSlider* slider = nil;
        for (NSView* subview in [containerView subviews]) {
            if ([subview isKindOfClass:[NSSlider class]]) {
                slider = (NSSlider*)subview;
                break;
            }
        }
        
        if (slider == nil) {
            ui_log_to_reaper(LOG_ERROR, "Slider not found in parameter view");
            return 0.0;
        }
        
        // Get value
        return [slider doubleValue];
    }
}

// Set parameter callback
bool macos_param_view_set_callback(void* view, SliderCallback callback) {
    if (view == NULL || callback == NULL) {
        ui_log_to_reaper(LOG_ERROR, "Null pointer in macos_param_view_set_callback");
        return false;
    }
    
    @autoreleasepool {
        NSView* containerView = (NSView*)view;
        
        // Find controller
        RPRParamViewController* controller = [g_param_controllers objectForKey:[NSValue valueWithPointer:containerView]];
        if (controller == nil) {
            ui_log_to_reaper(LOG_ERROR, "Parameter controller not found");
            return false;
        }
        
        // Set callback
        [controller setSliderCallback:callback context:view];
        
        return true;
    }
}
