#pragma once

#include <string>
#include <map>
#include <memory>

// Forward declaration for JSON handling
namespace nlohmann {
    class json;
}

class ConfigManager {
private:
    std::string configFilePath;
    std::map<std::string, std::string> stringSettings;
    std::map<std::string, int> intSettings;
    std::map<std::string, bool> boolSettings;
    std::map<std::string, double> doubleSettings;
    
    bool configLoaded;
    bool autoSave;

public:
    ConfigManager(const std::string& filePath = "config.json");
    ~ConfigManager();
    
    // Configuration loading and saving
    bool LoadConfiguration();
    bool SaveConfiguration();
    bool ReloadConfiguration();
    
    // String settings
    void SetString(const std::string& key, const std::string& value);
    std::string GetString(const std::string& key, const std::string& defaultValue = "");
    
    // Integer settings
    void SetInt(const std::string& key, int value);
    int GetInt(const std::string& key, int defaultValue = 0);
    
    // Boolean settings
    void SetBool(const std::string& key, bool value);
    bool GetBool(const std::string& key, bool defaultValue = false);
    
    // Double settings
    void SetDouble(const std::string& key, double value);
    double GetDouble(const std::string& key, double defaultValue = 0.0);
    
    // API-specific settings
    void SetRestEndpoint(const std::string& endpoint);
    std::string GetRestEndpoint() const;
    
    void SetGrpcEndpoint(const std::string& endpoint);
    std::string GetGrpcEndpoint() const;
    
    void SetRestTimeout(int timeout);
    int GetRestTimeout() const;
    
    void SetGrpcTimeout(int timeout);
    int GetGrpcTimeout() const;
    
    void SetRetryAttempts(int attempts);
    int GetRetryAttempts() const;
    
    void SetRetryDelay(int delay);
    int GetRetryDelay() const;
    
    // UI settings
    void SetWindowWidth(int width);
    int GetWindowWidth() const;
    
    void SetWindowHeight(int height);
    int GetWindowHeight() const;
    
    void SetAutoSave(bool enabled);
    bool GetAutoSave() const;
    
    void SetLogLevel(const std::string& level);
    std::string GetLogLevel() const;
    
    void SetLogFile(const std::string& file);
    std::string GetLogFile() const;
    
    // Utility methods
    bool IsConfigLoaded() const { return configLoaded; }
    std::string GetConfigFilePath() const { return configFilePath; }
    void SetConfigFilePath(const std::string& path);
    
    // Validation
    bool ValidateConfiguration();
    std::string GetValidationErrors() const;

private:
    // Internal methods
    void InitializeDefaultSettings();
    bool ParseJsonConfig(const std::string& jsonContent);
    std::string GenerateJsonConfig() const;
    void CreateDefaultConfigFile();
    
    // Helper methods
    std::string EscapeString(const std::string& str) const;
    std::string UnescapeString(const std::string& str) const;
    
    // Default values
    static const std::string DEFAULT_REST_ENDPOINT;
    static const std::string DEFAULT_GRPC_ENDPOINT;
    static const int DEFAULT_REST_TIMEOUT;
    static const int DEFAULT_GRPC_TIMEOUT;
    static const int DEFAULT_RETRY_ATTEMPTS;
    static const int DEFAULT_RETRY_DELAY;
    static const int DEFAULT_WINDOW_WIDTH;
    static const int DEFAULT_WINDOW_HEIGHT;
    static const std::string DEFAULT_LOG_LEVEL;
    static const std::string DEFAULT_LOG_FILE;
}; 