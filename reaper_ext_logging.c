#include "reaper_ext_logging.h"
#include <stdio.h>
#include <stdlib.h>  // For getenv
#include <string.h>
#include <stdarg.h>
#include <time.h>
#include <ctype.h>  // For strcasecmp alternatives if needed

// Log level to control verbosity
static LogLevel current_log_level = LOG_INFO;
static bool logging_enabled = false;
static char custom_log_path[512] = {0};
static bool custom_path_set = false;

// Platform-specific directory functions
#ifdef _WIN32
#include <windows.h>
#include <shlobj.h>
#define PATH_SEPARATOR '\\'
// Windows doesn't have strcasecmp, it has _stricmp
#define strcasecmp _stricmp
#else
#include <unistd.h>
#include <pwd.h>
#include <sys/types.h>
#define PATH_SEPARATOR '/'
#endif

// Get path to user's home directory
static const char* get_home_directory() {
#ifdef _WIN32
    static char path[MAX_PATH];
    if (SHGetFolderPathA(NULL, CSIDL_PROFILE, NULL, 0, path) != S_OK) {
        return NULL;
    }
    return path;
#else
    const char* home = getenv("HOME");
    if (home) {
        return home;
    }
    
    // Fallback to using getpwuid if HOME is not available
    struct passwd* pw = getpwuid(getuid());
    if (pw) {
        return pw->pw_dir;
    }
    
    return NULL;
#endif
}

// Determine log file path based on REAPER paths, environment variables, or defaults
static const char* get_log_file_path() {
    static char log_path[512];
    
    // If a custom path was explicitly set, use it
    if (custom_path_set) {
        return custom_log_path;
    }
    
    // Check for environment variable override
    const char* env_path = getenv("REAPER_GO_LOG_PATH");
    if (env_path && strlen(env_path) > 0) {
        strncpy(log_path, env_path, sizeof(log_path) - 1);
        log_path[sizeof(log_path) - 1] = '\0';
        return log_path;
    }
    
    // Use home directory + .reaper folder
    const char* home_dir = get_home_directory();
    if (!home_dir) {
        // Last resort fallback to current directory
        strcpy(log_path, "reaper_go_ext.log");
        return log_path;
    }
    
#ifdef _WIN32
    snprintf(log_path, sizeof(log_path), "%s\\AppData\\Roaming\\REAPER\\go_ext.log", home_dir);
#elif defined(__APPLE__)
    snprintf(log_path, sizeof(log_path), "%s/Library/Application Support/REAPER/go_ext.log", home_dir);
#else // Linux/Unix
    snprintf(log_path, sizeof(log_path), "%s/.config/REAPER/go_ext.log", home_dir);
#endif

    return log_path;
}

// Initialize the logging system, called at plugin startup
void log_init() {
    // Check for environment variable to enable/disable logging
    const char* env_enabled = getenv("REAPER_GO_LOG_ENABLED");
    if (env_enabled && (strcmp(env_enabled, "1") == 0 || 
                      strcasecmp(env_enabled, "true") == 0 || 
                      strcasecmp(env_enabled, "yes") == 0)) {
        logging_enabled = true;
    }
    
    // Check for environment variable to set log level
    const char* env_level = getenv("REAPER_GO_LOG_LEVEL");
    if (env_level) {
        if (strcasecmp(env_level, "error") == 0) {
            current_log_level = LOG_ERROR;
        } else if (strcasecmp(env_level, "warning") == 0) {
            current_log_level = LOG_WARNING;
        } else if (strcasecmp(env_level, "info") == 0) {
            current_log_level = LOG_INFO;
        } else if (strcasecmp(env_level, "debug") == 0) {
            current_log_level = LOG_DEBUG;
        } else if (strcasecmp(env_level, "trace") == 0) {
            current_log_level = LOG_TRACE;
        }
    }
    
    // Log initialization message if enabled
    if (logging_enabled) {
        const char* path = get_log_file_path();
        // Clear log file on startup for a fresh log
        FILE* f = fopen(path, "w");
        if (f) {
            fprintf(f, "--- REAPER Go Extension Log Started ---\n");
            fclose(f);
        }
        
        // First log message
        log_message(LOG_INFO, "log_init", "Logging initialized");
    }
}

// Cleanup logging, called at plugin shutdown
void log_cleanup() {
    if (logging_enabled) {
        log_message(LOG_INFO, "log_cleanup", "Logging system shutting down");
    }
}

// Set custom log path
void log_set_path(const char* path) {
    if (path && strlen(path) > 0) {
        strncpy(custom_log_path, path, sizeof(custom_log_path) - 1);
        custom_log_path[sizeof(custom_log_path) - 1] = '\0';
        custom_path_set = true;
        
        // Log the path change
        if (logging_enabled) {
            log_message(LOG_INFO, "log_set_path", "Log path set to new location");
        }
    }
}

// Enable or disable logging at runtime
void log_set_enabled(bool enabled) {
    if (enabled != logging_enabled) {
        logging_enabled = enabled;
        
        if (enabled) {
            // Log that logging was enabled
            log_message(LOG_INFO, "log_set_enabled", "Logging enabled");
        }
    }
}

// Set the current log level
void log_set_level(LogLevel level) {
    if (level != current_log_level && level >= LOG_ERROR && level <= LOG_TRACE) {
        LogLevel old_level = current_log_level;
        current_log_level = level;
        
        if (logging_enabled) {
            log_message(LOG_INFO, "log_set_level", "Log level changed");
        }
    }
}

// Get current log level
LogLevel log_get_level() {
    return current_log_level;
}

// Check if logging is enabled
bool log_is_enabled() {
    return logging_enabled;
}

// Log a message directly without format - for Go code
void log_message(LogLevel level, const char* func, const char* message) {
    // Skip if logging is disabled or message level is more verbose than current setting
    if (!logging_enabled || level > current_log_level) {
        return;
    }
    
    // Get log file path
    const char* log_path = get_log_file_path();
    
    // Open the log file in append mode
    FILE* log_file = fopen(log_path, "a");
    if (log_file == NULL) {
        return;
    }
    
    // Get current time
    time_t now = time(NULL);
    struct tm* timeinfo = localtime(&now);
    char timestamp[32];
    strftime(timestamp, sizeof(timestamp), "%Y-%m-%d %H:%M:%S", timeinfo);
    
    // Convert level to string
    const char* level_str;
    switch (level) {
        case LOG_ERROR:   level_str = "ERROR"; break;
        case LOG_WARNING: level_str = "WARN"; break;
        case LOG_INFO:    level_str = "INFO"; break;
        case LOG_DEBUG:   level_str = "DEBUG"; break;
        case LOG_TRACE:   level_str = "TRACE"; break;
        default:          level_str = "UNKNOWN"; break;
    }
    
    // Print the log header and message
    fprintf(log_file, "[%s] [%s] [%s] %s\n", 
            timestamp, level_str, func, message);
    fflush(log_file);
    
    // Close the file immediately to avoid keeping handles open
    fclose(log_file);
}

// Log a message with format string - for C code
void log_message_v(LogLevel level, const char* func, const char* format, ...) {
    // Skip if logging is disabled or message level is more verbose than current setting
    if (!logging_enabled || level > current_log_level) {
        return;
    }
    
    // Format the message using vsnprintf
    char message[2048]; // Reasonable buffer size for most messages
    va_list args;
    va_start(args, format);
    vsnprintf(message, sizeof(message), format, args);
    va_end(args);
    
    // Call the standard log_message function with the formatted message
    log_message(level, func, message);
}
