#include "MainForm.h"
#include "RESTClient.h"
#include "GRPCClient.h"
#include "ConfigManager.h"
#include "LogManager.h"
#include <commctrl.h>
#include <sstream>
#include <iomanip>
#include <chrono>

// Global instance for the main form
static MainForm* g_mainForm = nullptr;

// Control IDs
enum ControlIDs {
    ID_CONNECT_BUTTON = 1001,
    ID_DISCONNECT_BUTTON = 1002,
    ID_REFRESH_BUTTON = 1003,
    ID_CLEAR_LOG_BUTTON = 1004,
    ID_EXPORT_LOG_BUTTON = 1005,
    
    // Tab control
    ID_TAB_CONTROL = 2000,
    
    // Account tab controls
    ID_ACCOUNT_CREATE_BUTTON = 3001,
    ID_ACCOUNT_LIST_BUTTON = 3002,
    ID_ACCOUNT_DETAILS_BUTTON = 3003,
    
    // Component tab controls
    ID_COMPONENT_REGISTER_BUTTON = 4001,
    ID_COMPONENT_VERIFY_BUTTON = 4002,
    ID_COMPONENT_LIST_BUTTON = 4003,
    
    // Privacy tab controls
    ID_PRIVACY_REGISTER_BUTTON = 5001,
    ID_PRIVACY_VERIFY_BUTTON = 5002,
    ID_PRIVACY_AUTHORIZE_BUTTON = 5003,
    
    // LCT tab controls
    ID_LCT_CREATE_BUTTON = 6001,
    ID_LCT_TERMINATE_BUTTON = 6002,
    ID_LCT_LIST_BUTTON = 6003,
    
    // Pairing tab controls
    ID_PAIRING_INITIATE_BUTTON = 7001,
    ID_PAIRING_COMPLETE_BUTTON = 7002,
    ID_PAIRING_REVOKE_BUTTON = 7003,
    
    // Pairing Queue tab controls
    ID_QUEUE_REQUEST_BUTTON = 8001,
    ID_QUEUE_PROCESS_BUTTON = 8002,
    ID_QUEUE_CANCEL_BUTTON = 8003,
    
    // Trust tab controls
    ID_TRUST_CREATE_BUTTON = 9001,
    ID_TRUST_UPDATE_BUTTON = 9002,
    ID_TRUST_LIST_BUTTON = 9003,
    
    // Energy tab controls
    ID_ENERGY_CREATE_BUTTON = 10001,
    ID_ENERGY_EXECUTE_BUTTON = 10002,
    ID_ENERGY_LIST_BUTTON = 10003,
    
    // Performance tab controls
    ID_PERF_REST_TEST_BUTTON = 11001,
    ID_PERF_GRPC_TEST_BUTTON = 11002,
    ID_PERF_COMPARE_BUTTON = 11003,
    ID_PERF_STREAM_BUTTON = 11004,
    
    // Log tab controls
    ID_LOG_LISTBOX = 12001,
    ID_LOG_LEVEL_COMBO = 12002
};

MainForm::MainForm() 
    : hMainWindow(nullptr)
    , hTabControl(nullptr)
    , hStatusBar(nullptr)
    , hMenuBar(nullptr)
    , streamingActive(false)
    , isConnected(false)
    , windowWidth(WINDOW_WIDTH)
    , windowHeight(WINDOW_HEIGHT)
    , restEndpoint("http://localhost:8080")
    , grpcEndpoint("localhost:9092")
    , grpcAvailable(false) {
    
    g_mainForm = this;
}

MainForm::~MainForm() {
    if (streamingActive) {
        streamingActive = false;
        if (streamingThread.joinable()) {
            streamingThread.join();
        }
    }
    g_mainForm = nullptr;
}

bool MainForm::Initialize(HINSTANCE hInstance, int nCmdShow) {
    // Initialize common controls
    InitializeCommonControls();
    
    // Load configuration
    LoadConfiguration();
    
    // Create the main window
    CreateMainWindow(hInstance);
    
    if (!hMainWindow) {
        return false;
    }
    
    // Show the window
    ShowWindow(hMainWindow, nCmdShow);
    UpdateWindow(hMainWindow);
    
    return true;
}

void MainForm::InitializeCommonControls() {
    INITCOMMONCONTROLSEX icex;
    icex.dwSize = sizeof(INITCOMMONCONTROLSEX);
    icex.dwICC = ICC_TAB_CLASSES | ICC_BAR_CLASSES | ICC_LISTVIEW_CLASSES | 
                 ICC_TREEVIEW_CLASSES | ICC_BUTTON_CLASSES | ICC_EDIT_CLASSES |
                 ICC_STATIC_CLASSES | ICC_COMBOBOX_CLASSES;
    InitCommonControlsEx(&icex);
}

void MainForm::CreateMainWindow(HINSTANCE hInstance) {
    // Register window class
    WNDCLASSEX wc = {};
    wc.cbSize = sizeof(WNDCLASSEX);
    wc.style = CS_HREDRAW | CS_VREDRAW;
    wc.lpfnWndProc = WindowProc;
    wc.hInstance = hInstance;
    wc.hIcon = LoadIcon(nullptr, IDI_APPLICATION);
    wc.hCursor = LoadCursor(nullptr, IDC_ARROW);
    wc.hbrBackground = (HBRUSH)(COLOR_WINDOW + 1);
    wc.lpszClassName = L"APIBridgeDemoMainForm";
    wc.hIconSm = LoadIcon(nullptr, IDI_APPLICATION);
    
    RegisterClassEx(&wc);
    
    // Create the main window
    hMainWindow = CreateWindowEx(
        0,
        L"APIBridgeDemoMainForm",
        L"Web4 Race Car Battery Management API Bridge Demo",
        WS_OVERLAPPEDWINDOW,
        CW_USEDEFAULT, CW_USEDEFAULT,
        windowWidth, windowHeight,
        nullptr, nullptr, hInstance, nullptr
    );
}

void MainForm::OnCreate(HWND hwnd) {
    // Create status bar
    CreateStatusBar();
    
    // Create menu bar
    CreateMenuBar();
    
    // Create tab control
    CreateTabControl();
    
    // Create common controls
    CreateCommonControls();
    
    // Layout controls
    LayoutControls();
    
    // Initialize API clients
    configManager = std::make_unique<ConfigManager>();
    logManager = std::make_unique<LogManager>();
    
    // Update initial status
    UpdateStatusBar("Ready - Click Connect to start");
    UpdateConnectionStatus(false);
    UpdateRestStatus(false);
    UpdateGrpcStatus(false);
}

void MainForm::CreateStatusBar() {
    hStatusBar = CreateWindowEx(
        0, STATUSCLASSNAME, nullptr,
        WS_CHILD | WS_VISIBLE | SBARS_SIZEGRIP,
        0, 0, 0, 0,
        hMainWindow, nullptr, GetModuleHandle(nullptr), nullptr
    );
    
    // Create status bar parts
    int parts[] = { 200, 400, 600, -1 };
    SendMessage(hStatusBar, SB_SETPARTS, 4, (LPARAM)parts);
    
    // Set initial text
    SendMessage(hStatusBar, SB_SETTEXT, 0, (LPARAM)L"Ready");
    SendMessage(hStatusBar, SB_SETTEXT, 1, (LPARAM)L"REST: Disconnected");
    SendMessage(hStatusBar, SB_SETTEXT, 2, (LPARAM)L"gRPC: Disconnected");
    SendMessage(hStatusBar, SB_SETTEXT, 3, (LPARAM)L"");
}

void MainForm::CreateMenuBar() {
    hMenuBar = CreateWindowEx(
        0, L"STATIC", L"",
        WS_CHILD | WS_VISIBLE | SS_CENTER,
        0, 0, 0, MENU_HEIGHT,
        hMainWindow, nullptr, GetModuleHandle(nullptr), nullptr
    );
}

void MainForm::CreateTabControl() {
    hTabControl = CreateWindowEx(
        0, WC_TABCONTROL, nullptr,
        WS_CHILD | WS_VISIBLE | WS_CLIPSIBLINGS,
        0, 0, 0, 0,
        hMainWindow, (HMENU)ID_TAB_CONTROL, GetModuleHandle(nullptr), nullptr
    );
    
    // Add tabs
    TCITEM tie = {};
    tie.mask = TCIF_TEXT;
    
    tie.pszText = L"Account";
    TabCtrl_InsertItem(hTabControl, 0, &tie);
    
    tie.pszText = L"Component";
    TabCtrl_InsertItem(hTabControl, 1, &tie);
    
    tie.pszText = L"Privacy";
    TabCtrl_InsertItem(hTabControl, 2, &tie);
    
    tie.pszText = L"LCT";
    TabCtrl_InsertItem(hTabControl, 3, &tie);
    
    tie.pszText = L"Pairing";
    TabCtrl_InsertItem(hTabControl, 4, &tie);
    
    tie.pszText = L"Pairing Queue";
    TabCtrl_InsertItem(hTabControl, 5, &tie);
    
    tie.pszText = L"Trust";
    TabCtrl_InsertItem(hTabControl, 6, &tie);
    
    tie.pszText = L"Energy";
    TabCtrl_InsertItem(hTabControl, 7, &tie);
    
    tie.pszText = L"Performance";
    TabCtrl_InsertItem(hTabControl, 8, &tie);
    
    tie.pszText = L"Logs";
    TabCtrl_InsertItem(hTabControl, 9, &tie);
}

void MainForm::CreateCommonControls() {
    // Create connection buttons
    hConnectButton = CreateWindowEx(
        0, L"BUTTON", L"Connect",
        WS_CHILD | WS_VISIBLE | BS_PUSHBUTTON,
        10, 10, 80, 25,
        hMainWindow, (HMENU)ID_CONNECT_BUTTON, GetModuleHandle(nullptr), nullptr
    );
    
    hDisconnectButton = CreateWindowEx(
        0, L"BUTTON", L"Disconnect",
        WS_CHILD | WS_VISIBLE | BS_PUSHBUTTON,
        100, 10, 80, 25,
        hMainWindow, (HMENU)ID_DISCONNECT_BUTTON, GetModuleHandle(nullptr), nullptr
    );
    
    hRefreshButton = CreateWindowEx(
        0, L"BUTTON", L"Refresh",
        WS_CHILD | WS_VISIBLE | BS_PUSHBUTTON,
        190, 10, 80, 25,
        hMainWindow, (HMENU)ID_REFRESH_BUTTON, GetModuleHandle(nullptr), nullptr
    );
    
    // Initially disable disconnect and refresh
    EnableWindow(hDisconnectButton, FALSE);
    EnableWindow(hRefreshButton, FALSE);
}

void MainForm::LayoutControls() {
    RECT clientRect;
    GetClientRect(hMainWindow, &clientRect);
    
    int clientWidth = clientRect.right - clientRect.left;
    int clientHeight = clientRect.bottom - clientRect.top;
    
    // Position status bar
    SetWindowPos(hStatusBar, nullptr, 0, clientHeight - STATUS_HEIGHT, 
                clientWidth, STATUS_HEIGHT, SWP_NOZORDER);
    
    // Position menu bar
    SetWindowPos(hMenuBar, nullptr, 0, 0, clientWidth, MENU_HEIGHT, SWP_NOZORDER);
    
    // Position tab control
    SetWindowPos(hTabControl, nullptr, 0, MENU_HEIGHT + 50, 
                clientWidth, clientHeight - MENU_HEIGHT - STATUS_HEIGHT - 50, SWP_NOZORDER);
    
    // Position connection buttons
    SetWindowPos(hConnectButton, nullptr, 10, MENU_HEIGHT + 10, 80, 25, SWP_NOZORDER);
    SetWindowPos(hDisconnectButton, nullptr, 100, MENU_HEIGHT + 10, 80, 25, SWP_NOZORDER);
    SetWindowPos(hRefreshButton, nullptr, 190, MENU_HEIGHT + 10, 80, 25, SWP_NOZORDER);
}

void MainForm::OnCommand(HWND hwnd, int wmId, HWND hwndCtl) {
    switch (wmId) {
        case ID_CONNECT_BUTTON:
            ConnectToAPI();
            break;
            
        case ID_DISCONNECT_BUTTON:
            DisconnectFromAPI();
            break;
            
        case ID_REFRESH_BUTTON:
            // Refresh current tab content
            int currentTab = TabCtrl_GetCurSel(hTabControl);
            UpdateTabContent(currentTab);
            break;
            
        case ID_ACCOUNT_CREATE_BUTTON:
            TestAccountManagement();
            break;
            
        case ID_COMPONENT_REGISTER_BUTTON:
            TestComponentRegistry();
            break;
            
        case ID_PRIVACY_REGISTER_BUTTON:
            TestPrivacyFeatures();
            break;
            
        case ID_LCT_CREATE_BUTTON:
            TestLCTManagement();
            break;
            
        case ID_PAIRING_INITIATE_BUTTON:
            TestPairingProcess();
            break;
            
        case ID_QUEUE_REQUEST_BUTTON:
            TestPairingQueue();
            break;
            
        case ID_TRUST_CREATE_BUTTON:
            TestTrustTensor();
            break;
            
        case ID_ENERGY_CREATE_BUTTON:
            TestEnergyOperations();
            break;
            
        case ID_PERF_COMPARE_BUTTON:
            ComparePerformance();
            break;
            
        case ID_PERF_STREAM_BUTTON:
            if (streamingActive) {
                StopStreaming();
            } else {
                StartStreaming();
            }
            break;
            
        case ID_CLEAR_LOG_BUTTON:
            ClearLogs();
            break;
            
        case ID_EXPORT_LOG_BUTTON:
            ExportLogs();
            break;
    }
}

void MainForm::OnNotify(HWND hwnd, int wmId, LPNMHDR pnmh) {
    if (pnmh->idFrom == ID_TAB_CONTROL && pnmh->code == TCN_SELCHANGE) {
        int currentTab = TabCtrl_GetCurSel(hTabControl);
        HandleTabChange(currentTab);
    }
}

void MainForm::HandleTabChange(int tabIndex) {
    // Hide all tab pages
    ShowWindow(hAccountTab, SW_HIDE);
    ShowWindow(hComponentTab, SW_HIDE);
    ShowWindow(hPrivacyTab, SW_HIDE);
    ShowWindow(hLCTTab, SW_HIDE);
    ShowWindow(hPairingTab, SW_HIDE);
    ShowWindow(hPairingQueueTab, SW_HIDE);
    ShowWindow(hTrustTab, SW_HIDE);
    ShowWindow(hEnergyTab, SW_HIDE);
    ShowWindow(hPerformanceTab, SW_HIDE);
    ShowWindow(hLogsTab, SW_HIDE);
    
    // Show the selected tab
    switch (tabIndex) {
        case 0: ShowWindow(hAccountTab, SW_SHOW); break;
        case 1: ShowWindow(hComponentTab, SW_SHOW); break;
        case 2: ShowWindow(hPrivacyTab, SW_SHOW); break;
        case 3: ShowWindow(hLCTTab, SW_SHOW); break;
        case 4: ShowWindow(hPairingTab, SW_SHOW); break;
        case 5: ShowWindow(hPairingQueueTab, SW_SHOW); break;
        case 6: ShowWindow(hTrustTab, SW_SHOW); break;
        case 7: ShowWindow(hEnergyTab, SW_SHOW); break;
        case 8: ShowWindow(hPerformanceTab, SW_SHOW); break;
        case 9: ShowWindow(hLogsTab, SW_SHOW); break;
    }
    
    UpdateTabContent(tabIndex);
}

void MainForm::UpdateTabContent(int tabIndex) {
    // Update content based on current tab and connection status
    if (!isConnected) {
        AddLogMessage("Not connected to API - please connect first", "WARNING");
        return;
    }
    
    switch (tabIndex) {
        case 0: // Account tab
            AddLogMessage("Account tab selected - ready for account operations");
            break;
        case 1: // Component tab
            AddLogMessage("Component tab selected - ready for component operations");
            break;
        case 2: // Privacy tab
            AddLogMessage("Privacy tab selected - ready for privacy operations");
            break;
        case 3: // LCT tab
            AddLogMessage("LCT tab selected - ready for LCT operations");
            break;
        case 4: // Pairing tab
            AddLogMessage("Pairing tab selected - ready for pairing operations");
            break;
        case 5: // Pairing Queue tab
            AddLogMessage("Pairing Queue tab selected - ready for queue operations");
            break;
        case 6: // Trust tab
            AddLogMessage("Trust tab selected - ready for trust operations");
            break;
        case 7: // Energy tab
            AddLogMessage("Energy tab selected - ready for energy operations");
            break;
        case 8: // Performance tab
            AddLogMessage("Performance tab selected - ready for performance testing");
            break;
        case 9: // Logs tab
            AddLogMessage("Logs tab selected - viewing application logs");
            break;
    }
}

void MainForm::ConnectToAPI() {
    try {
        UpdateStatusBar("Connecting to API...");
        
        // Initialize REST client
        restClient = std::make_unique<RESTClient>(restEndpoint);
        
        // Try to initialize gRPC client
        try {
            grpcClient = std::make_unique<GRPCClient>(grpcEndpoint);
            grpcAvailable = true;
            AddLogMessage("gRPC client initialized successfully", "SUCCESS");
        } catch (const std::exception& e) {
            grpcAvailable = false;
            AddLogMessage("gRPC client not available: " + std::string(e.what()), "WARNING");
        }
        
        isConnected = true;
        
        // Update UI
        EnableWindow(hConnectButton, FALSE);
        EnableWindow(hDisconnectButton, TRUE);
        EnableWindow(hRefreshButton, TRUE);
        
        UpdateConnectionStatus(true);
        UpdateRestStatus(true);
        UpdateGrpcStatus(grpcAvailable);
        UpdateStatusBar("Connected to API - Ready for operations");
        
        AddLogMessage("Successfully connected to API", "SUCCESS");
        
    } catch (const std::exception& e) {
        ShowError("Connection Error", "Failed to connect to API: " + std::string(e.what()));
        UpdateStatusBar("Connection failed");
        AddLogMessage("Connection failed: " + std::string(e.what()), "ERROR");
    }
}

void MainForm::DisconnectFromAPI() {
    // Stop streaming if active
    if (streamingActive) {
        StopStreaming();
    }
    
    // Clean up clients
    restClient.reset();
    grpcClient.reset();
    
    isConnected = false;
    grpcAvailable = false;
    
    // Update UI
    EnableWindow(hConnectButton, TRUE);
    EnableWindow(hDisconnectButton, FALSE);
    EnableWindow(hRefreshButton, FALSE);
    
    UpdateConnectionStatus(false);
    UpdateRestStatus(false);
    UpdateGrpcStatus(false);
    UpdateStatusBar("Disconnected from API");
    
    AddLogMessage("Disconnected from API", "INFO");
}

void MainForm::UpdateStatusBar(const std::string& message) {
    std::wstring wmessage(message.begin(), message.end());
    SendMessage(hStatusBar, SB_SETTEXT, 0, (LPARAM)wmessage.c_str());
}

void MainForm::UpdateConnectionStatus(bool connected) {
    std::wstring status = connected ? L"Connected" : L"Disconnected";
    SendMessage(hStatusBar, SB_SETTEXT, 3, (LPARAM)status.c_str());
}

void MainForm::UpdateRestStatus(bool available) {
    std::wstring status = available ? L"REST: Connected" : L"REST: Disconnected";
    SendMessage(hStatusBar, SB_SETTEXT, 1, (LPARAM)status.c_str());
}

void MainForm::UpdateGrpcStatus(bool available) {
    std::wstring status = available ? L"gRPC: Connected" : L"gRPC: Disconnected";
    SendMessage(hStatusBar, SB_SETTEXT, 2, (LPARAM)status.c_str());
}

void MainForm::AddLogMessage(const std::string& message, const std::string& level) {
    auto now = std::chrono::system_clock::now();
    auto time_t = std::chrono::system_clock::to_time_t(now);
    auto tm = *std::localtime(&time_t);
    
    std::ostringstream oss;
    oss << std::put_time(&tm, "%H:%M:%S") << " [" << level << "] " << message;
    
    std::string logEntry = oss.str();
    
    // Add to log manager
    if (logManager) {
        logManager->AddLog(logEntry, level);
    }
    
    // Update log listbox if it exists
    if (hLogsTab) {
        // This would update the log display in the logs tab
        // Implementation depends on the specific log display control
    }
}

void MainForm::ShowError(const std::string& title, const std::string& message) {
    std::wstring wtitle(title.begin(), title.end());
    std::wstring wmessage(message.begin(), message.end());
    MessageBox(hMainWindow, wmessage.c_str(), wtitle.c_str(), MB_OK | MB_ICONERROR);
}

void MainForm::ShowInfo(const std::string& title, const std::string& message) {
    std::wstring wtitle(title.begin(), title.end());
    std::wstring wmessage(message.begin(), message.end());
    MessageBox(hMainWindow, wmessage.c_str(), wtitle.c_str(), MB_OK | MB_ICONINFORMATION);
}

void MainForm::ShowSuccess(const std::string& title, const std::string& message) {
    std::wstring wtitle(title.begin(), title.end());
    std::wstring wmessage(message.begin(), message.end());
    MessageBox(hMainWindow, wmessage.c_str(), wtitle.c_str(), MB_OK | MB_ICONINFORMATION);
}

void MainForm::LoadConfiguration() {
    // Load configuration from file
    // This would use the ConfigManager to load settings
}

void MainForm::SaveConfiguration() {
    // Save configuration to file
    // This would use the ConfigManager to save settings
}

void MainForm::OnClose(HWND hwnd) {
    if (isConnected) {
        int result = MessageBox(hwnd, L"Are you sure you want to exit? This will disconnect from the API.", 
                               L"Confirm Exit", MB_YESNO | MB_ICONQUESTION);
        if (result == IDNO) {
            return;
        }
    }
    DestroyWindow(hwnd);
}

void MainForm::OnDestroy(HWND hwnd) {
    SaveConfiguration();
    PostQuitMessage(0);
}

void MainForm::OnSize(HWND hwnd, int width, int height) {
    windowWidth = width;
    windowHeight = height;
    LayoutControls();
}

int MainForm::RunMessageLoop() {
    MSG msg = {};
    while (GetMessage(&msg, nullptr, 0, 0)) {
        TranslateMessage(&msg);
        DispatchMessage(&msg);
    }
    return (int)msg.wParam;
}

// Window procedure
LRESULT CALLBACK MainForm::WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam) {
    if (g_mainForm) {
        switch (uMsg) {
            case WM_CREATE:
                g_mainForm->OnCreate(hwnd);
                return 0;
                
            case WM_COMMAND:
                g_mainForm->OnCommand(hwnd, LOWORD(wParam), (HWND)lParam);
                return 0;
                
            case WM_NOTIFY:
                g_mainForm->OnNotify(hwnd, LOWORD(wParam), (LPNMHDR)lParam);
                return 0;
                
            case WM_SIZE:
                g_mainForm->OnSize(hwnd, LOWORD(lParam), HIWORD(lParam));
                return 0;
                
            case WM_CLOSE:
                g_mainForm->OnClose(hwnd);
                return 0;
                
            case WM_DESTROY:
                g_mainForm->OnDestroy(hwnd);
                return 0;
        }
    }
    
    return DefWindowProc(hwnd, uMsg, wParam, lParam);
}

// Placeholder implementations for API test methods
void MainForm::TestAccountManagement() {
    AddLogMessage("Testing account management...", "INFO");
    // Implementation would go here
}

void MainForm::TestComponentRegistry() {
    AddLogMessage("Testing component registry...", "INFO");
    // Implementation would go here
}

void MainForm::TestPrivacyFeatures() {
    AddLogMessage("Testing privacy features...", "INFO");
    // Implementation would go here
}

void MainForm::TestLCTManagement() {
    AddLogMessage("Testing LCT management...", "INFO");
    // Implementation would go here
}

void MainForm::TestPairingProcess() {
    AddLogMessage("Testing pairing process...", "INFO");
    // Implementation would go here
}

void MainForm::TestPairingQueue() {
    AddLogMessage("Testing pairing queue...", "INFO");
    // Implementation would go here
}

void MainForm::TestTrustTensor() {
    AddLogMessage("Testing trust tensor...", "INFO");
    // Implementation would go here
}

void MainForm::TestEnergyOperations() {
    AddLogMessage("Testing energy operations...", "INFO");
    // Implementation would go here
}

void MainForm::StartStreaming() {
    if (!grpcAvailable) {
        ShowError("Streaming Error", "gRPC is not available for streaming");
        return;
    }
    
    streamingActive = true;
    streamingThread = std::thread(&MainForm::StreamingWorker, this);
    AddLogMessage("Started streaming thread", "INFO");
}

void MainForm::StopStreaming() {
    streamingActive = false;
    if (streamingThread.joinable()) {
        streamingThread.join();
    }
    AddLogMessage("Stopped streaming thread", "INFO");
}

void MainForm::StreamingWorker() {
    // Streaming implementation would go here
    while (streamingActive) {
        // Process streaming data
        std::this_thread::sleep_for(std::chrono::milliseconds(100));
    }
}

void MainForm::ComparePerformance() {
    AddLogMessage("Comparing REST vs gRPC performance...", "INFO");
    // Implementation would go here
}

void MainForm::ClearLogs() {
    if (logManager) {
        logManager->ClearLogs();
    }
    AddLogMessage("Logs cleared", "INFO");
}

void MainForm::ExportLogs() {
    if (logManager) {
        std::string filename = "apibridge_demo_log_" + 
                              std::to_string(std::time(nullptr)) + ".txt";
        logManager->ExportLogs(filename);
        AddLogMessage("Logs exported to " + filename, "INFO");
    }
} 