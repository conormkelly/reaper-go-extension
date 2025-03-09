#ifndef REAPER_EXT_LOGGING_H
#define REAPER_EXT_LOGGING_H

#include <stdio.h>
#include <stdbool.h>
#include <stdarg.h>

// Define log levels
typedef enum {
    LOG_ERROR = 0,    // Serious errors that prevent operation
    LOG_WARNING = 1,  // Issues that might affect operation but don't prevent it
    LOG_INFO = 2,     // General information about extension operation
    LOG_DEBUG = 3,    // Detailed information useful for debugging
    LOG_TRACE = 4     // Very detailed tracing information
} LogLevel;

// Core logging functions
void log_init(void);
void log_cleanup(void);

// For direct calls to log_message without format from Go code
void log_message(LogLevel level, const char* func, const char* message);
// For variadic format string calls from C
void log_message_v(LogLevel level, const char* func, const char* format, ...);

// Configuration functions
void log_set_path(const char* path);
void log_set_enabled(bool enabled);
void log_set_level(LogLevel level);
LogLevel log_get_level(void);
bool log_is_enabled(void);

// Convenience macros for different log levels
#define LOG_ERROR_ENABLED (log_is_enabled() && log_get_level() >= LOG_ERROR)
#define LOG_WARNING_ENABLED (log_is_enabled() && log_get_level() >= LOG_WARNING)
#define LOG_INFO_ENABLED (log_is_enabled() && log_get_level() >= LOG_INFO)
#define LOG_DEBUG_ENABLED (log_is_enabled() && log_get_level() >= LOG_DEBUG)
#define LOG_TRACE_ENABLED (log_is_enabled() && log_get_level() >= LOG_TRACE)

// Conditional logging macros that only call the function if that level is enabled
// These use the variadic log_message_v function
#define LOG_ERROR(format, ...) \
    do { if (LOG_ERROR_ENABLED) log_message_v(LOG_ERROR, __func__, format, ##__VA_ARGS__); } while(0)

#define LOG_WARNING(format, ...) \
    do { if (LOG_WARNING_ENABLED) log_message_v(LOG_WARNING, __func__, format, ##__VA_ARGS__); } while(0)

#define LOG_INFO(format, ...) \
    do { if (LOG_INFO_ENABLED) log_message_v(LOG_INFO, __func__, format, ##__VA_ARGS__); } while(0)

#define LOG_DEBUG(format, ...) \
    do { if (LOG_DEBUG_ENABLED) log_message_v(LOG_DEBUG, __func__, format, ##__VA_ARGS__); } while(0)

#define LOG_TRACE(format, ...) \
    do { if (LOG_TRACE_ENABLED) log_message_v(LOG_TRACE, __func__, format, ##__VA_ARGS__); } while(0)

#endif // REAPER_EXT_LOGGING_H
