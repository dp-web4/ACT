#pragma once

#include <windows.h>
#include <commctrl.h>
#include <string>
#include <memory>
#include <vector>
#include <thread>
#include <atomic>

// Forward declarations
class RESTClient;
class GRPCClient;
class ConfigManager;
class LogManager;

// Main form class for the Windows Forms API Bridge Demo
class MainForm {
private:
    // Window handles
    HWND hMainWindow;
    HWND hTabControl;
    HWND hStatusBar;
    HWND hMenuBar;
    
    // Tab pages
    HWND hAccountTab;
    HWND hComponentTab;
    HWND hPrivacyTab;
    HWND hLCTTab;
    HWND hPairingTab;
    HWND hPairingQueueTab;
    HWND hTrustTab;
    HWND hEnergyTab;
    HWND hPerformanceTab;
    HWND hLogsTab;
    
    // Common controls
    HWND hConnectButton;
    HWND hDisconnectButton;
    HWND hRefreshButton;
    HWND hClearLogButton;
    HWND hExportLogButton;
    
    // Status indicators
    HWND hRestStatusLabel;
    HWND hGrpcStatusLabel;
    HWND hConnectionStatusLabel;
    
    // API clients
    std::unique_ptr<RESTClient> restClient;
    std::unique_ptr<GRPCClient> grpcClient;
    std::unique_ptr<ConfigManager> configManager;
    std::unique_ptr<LogManager> logManager;
    
    // Threading
    std::thread streamingThread;
    std::atomic<bool> streamingActive;
    std::atomic<bool> isConnected;
    
    // Window dimensions
    int windowWidth;
    int windowHeight;
    
    // Configuration
    std::string restEndpoint;
    std::string grpcEndpoint;
    bool grpcAvailable;

public:
    MainForm();
    ~MainForm();
    
    // Window management
    bool Initialize(HINSTANCE hInstance, int nCmdShow);
    static LRESULT CALLBACK WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam);
    
    // Event handlers
    void OnCreate(HWND hwnd);
    void OnSize(HWND hwnd, int width, int height);
    void OnCommand(HWND hwnd, int wmId, HWND hwndCtl);
    void OnNotify(HWND hwnd, int wmId, LPNMHDR pnmh);
    void OnClose(HWND hwnd);
    void OnDestroy(HWND hwnd);
    
    // UI creation methods
    void CreateMainWindow(HINSTANCE hInstance);
    void CreateTabControl();
    void CreateStatusBar();
    void CreateMenuBar();
    void CreateAccountTab();
    void CreateComponentTab();
    void CreatePrivacyTab();
    void CreateLCTTab();
    void CreatePairingTab();
    void CreatePairingQueueTab();
    void CreateTrustTab();
    void CreateEnergyTab();
    void CreatePerformanceTab();
    void CreateLogsTab();
    
    // API operations
    void ConnectToAPI();
    void DisconnectFromAPI();
    void TestAccountManagement();
    void TestComponentRegistry();
    void TestPrivacyFeatures();
    void TestLCTManagement();
    void TestPairingProcess();
    void TestPairingQueue();
    void TestTrustTensor();
    void TestEnergyOperations();
    void StartStreaming();
    void StopStreaming();
    void ComparePerformance();
    
    // UI update methods
    void UpdateStatusBar(const std::string& message);
    void UpdateConnectionStatus(bool connected);
    void UpdateRestStatus(bool available);
    void UpdateGrpcStatus(bool available);
    void AddLogMessage(const std::string& message, const std::string& level = "INFO");
    void ClearLogs();
    void ExportLogs();
    
    // Utility methods
    void ShowError(const std::string& title, const std::string& message);
    void ShowInfo(const std::string& title, const std::string& message);
    void ShowSuccess(const std::string& title, const std::string& message);
    std::string GetInputText(HWND hEdit);
    void SetInputText(HWND hEdit, const std::string& text);
    void EnableControl(HWND hControl, bool enable);
    
    // Message loop
    int RunMessageLoop();
    
    // Getters
    HWND GetMainWindow() const { return hMainWindow; }
    bool IsConnected() const { return isConnected; }
    bool IsGrpcAvailable() const { return grpcAvailable; }

private:
    // Helper methods
    void InitializeCommonControls();
    void LoadConfiguration();
    void SaveConfiguration();
    void CreateCommonControls();
    void LayoutControls();
    void HandleTabChange(int tabIndex);
    void UpdateTabContent(int tabIndex);
    
    // Threading helpers
    void StreamingWorker();
    void SafeUpdateUI(std::function<void()> updateFunc);
    
    // Constants
    static const int WINDOW_WIDTH = 1200;
    static const int WINDOW_HEIGHT = 800;
    static const int TAB_HEIGHT = 30;
    static const int STATUS_HEIGHT = 25;
    static const int MENU_HEIGHT = 30;
}; 