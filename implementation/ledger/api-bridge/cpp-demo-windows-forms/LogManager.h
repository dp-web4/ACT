#pragma once

#include <string>
#include <vector>
#include <mutex>
#include <fstream>
#include <memory>
#include <chrono>

// Log entry structure
struct LogEntry {
    std::chrono::system_clock::time_point timestamp;
    std::string level;
    std::string message;
    std::string source;
    
    LogEntry(const std::string& msg, const std::string& lvl = "INFO", const std::string& src = "Main")
        : timestamp(std::chrono::system_clock::now())
        , level(lvl)
        , message(msg)
        , source(src) {}
};

class LogManager {
private:
    std::vector<LogEntry> logEntries;
    std::mutex logMutex;
    std::ofstream logFile;
    
    std::string logFilePath;
    std::string currentLogLevel;
    size_t maxEntries;
    size_t maxFileSize;
    int maxFiles;
    
    bool consoleOutput;
    bool fileOutput;
    bool autoFlush;

public:
    LogManager(const std::string& filePath = "apibridge_demo.log");
    ~LogManager();
    
    // Logging methods
    void AddLog(const std::string& message, const std::string& level = "INFO", const std::string& source = "Main");
    void AddInfo(const std::string& message, const std::string& source = "Main");
    void AddWarning(const std::string& message, const std::string& source = "Main");
    void AddError(const std::string& message, const std::string& source = "Main");
    void AddSuccess(const std::string& message, const std::string& source = "Main");
    void AddDebug(const std::string& message, const std::string& source = "Main");
    
    // Log retrieval
    std::vector<LogEntry> GetLogs() const;
    std::vector<LogEntry> GetLogsByLevel(const std::string& level) const;
    std::vector<LogEntry> GetLogsBySource(const std::string& source) const;
    std::vector<LogEntry> GetLogsByTimeRange(const std::chrono::system_clock::time_point& start, 
                                           const std::chrono::system_clock::time_point& end) const;
    std::vector<LogEntry> GetRecentLogs(size_t count) const;
    
    // Log management
    void ClearLogs();
    void ClearLogsByLevel(const std::string& level);
    void ClearLogsBySource(const std::string& source);
    void ClearOldLogs(const std::chrono::system_clock::time_point& cutoff);
    
    // File operations
    bool ExportLogs(const std::string& filePath);
    bool ExportLogsByLevel(const std::string& filePath, const std::string& level);
    bool ExportLogsBySource(const std::string& filePath, const std::string& source);
    bool ExportLogsByTimeRange(const std::string& filePath, 
                              const std::chrono::system_clock::time_point& start,
                              const std::chrono::system_clock::time_point& end);
    
    // Configuration
    void SetLogLevel(const std::string& level);
    std::string GetLogLevel() const { return currentLogLevel; }
    
    void SetLogFile(const std::string& filePath);
    std::string GetLogFile() const { return logFilePath; }
    
    void SetMaxEntries(size_t max);
    size_t GetMaxEntries() const { return maxEntries; }
    
    void SetMaxFileSize(size_t maxSize);
    size_t GetMaxFileSize() const { return maxFileSize; }
    
    void SetMaxFiles(int max);
    int GetMaxFiles() const { return maxFiles; }
    
    void SetConsoleOutput(bool enabled);
    bool GetConsoleOutput() const { return consoleOutput; }
    
    void SetFileOutput(bool enabled);
    bool GetFileOutput() const { return fileOutput; }
    
    void SetAutoFlush(bool enabled);
    bool GetAutoFlush() const { return autoFlush; }
    
    // Utility methods
    size_t GetLogCount() const;
    std::string GetLogLevels() const;
    std::string GetSources() const;
    
    // Statistics
    struct LogStats {
        size_t totalEntries;
        size_t infoCount;
        size_t warningCount;
        size_t errorCount;
        size_t successCount;
        size_t debugCount;
        std::chrono::system_clock::time_point firstEntry;
        std::chrono::system_clock::time_point lastEntry;
    };
    
    LogStats GetStatistics() const;

private:
    // Internal methods
    void InitializeLogFile();
    void RotateLogFile();
    void FlushLogFile();
    void WriteLogEntry(const LogEntry& entry);
    bool ShouldLog(const std::string& level) const;
    std::string FormatLogEntry(const LogEntry& entry) const;
    std::string FormatTimestamp(const std::chrono::system_clock::time_point& timestamp) const;
    
    // Helper methods
    int GetLogLevelPriority(const std::string& level) const;
    std::string GetLogLevelColor(const std::string& level) const;
    void TrimLogEntries();
    
    // Default values
    static const size_t DEFAULT_MAX_ENTRIES;
    static const size_t DEFAULT_MAX_FILE_SIZE;
    static const int DEFAULT_MAX_FILES;
    static const std::string DEFAULT_LOG_LEVEL;
}; 