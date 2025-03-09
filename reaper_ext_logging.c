#include "reaper_ext_logging.h"

// Open a dedicated log file with a hardcoded path that we control
FILE* open_log_file(void) {
    // Use a hardcoded path that we know we have write access to
    const char* log_path = "~/reaper-ext.log";
    FILE* file = fopen(log_path, "a");
    if (file == NULL) {
        // If we can't open the log file, we'll print an error to stderr
        // but won't try anything fancy that might cause issues
        fprintf(stderr, "Failed to open log file: %s\n", log_path);
    }
    return file;
}

// Log a message to our dedicated log file
void log_message(const char* tag, const char* func, const char* format, ...) {
    FILE* log_file = open_log_file();
    if (log_file == NULL) {
        return;
    }
    
    // Get current time
    time_t now = time(NULL);
    struct tm* timeinfo = localtime(&now);
    char timestamp[32];
    strftime(timestamp, sizeof(timestamp), "%Y-%m-%d %H:%M:%S", timeinfo);
    
    // Print the log header
    fprintf(log_file, "[%s] [%s] [%s] ", timestamp, tag, func);
    
    // Print the formatted message
    va_list args;
    va_start(args, format);
    vfprintf(log_file, format, args);
    va_end(args);
    
    // End with newline and flush immediately
    fprintf(log_file, "\n");
    fflush(log_file);
    
    // Close the file immediately to avoid keeping handles open
    fclose(log_file);
}
