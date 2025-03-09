#ifndef REAPER_EXT_LOGGING_H
#define REAPER_EXT_LOGGING_H

#include <stdio.h>
#include <time.h>
#include <string.h>
#include <stdarg.h>

// Function declarations
void log_message(const char* tag, const char* func, const char* format, ...);
FILE* open_log_file(void);

// Simple macros for different log levels
#define LOG_DEBUG(format, ...) \
    log_message("DEBUG", __func__, format, ##__VA_ARGS__)

#define LOG_INFO(format, ...) \
    log_message("INFO", __func__, format, ##__VA_ARGS__)

#define LOG_WARNING(format, ...) \
    log_message("WARNING", __func__, format, ##__VA_ARGS__)

#define LOG_ERROR(format, ...) \
    log_message("ERROR", __func__, format, ##__VA_ARGS__)

#endif // REAPER_EXT_LOGGING_H
