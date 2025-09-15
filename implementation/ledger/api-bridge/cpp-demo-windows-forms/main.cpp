#include <windows.h>
#include <commctrl.h>
#include <iostream>
#include <memory>
#include "MainForm.h"

// Global variables
HINSTANCE g_hInstance = nullptr;
std::unique_ptr<MainForm> g_mainForm = nullptr;

// Function declarations
int WINAPI WinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, LPSTR lpCmdLine, int nCmdShow);
LRESULT CALLBACK WndProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam);
void InitializeApplication();
void CleanupApplication();
void ShowError(const std::string& message);

int WINAPI WinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, LPSTR lpCmdLine, int nCmdShow) {
    // Store instance handle
    g_hInstance = hInstance;
    
    // Initialize application
    InitializeApplication();
    
    try {
        // Create main form
        g_mainForm = std::make_unique<MainForm>();
        
        if (!g_mainForm->Initialize(hInstance, nCmdShow)) {
            ShowError("Failed to initialize main form");
            return -1;
        }
        
        // Run message loop
        int result = g_mainForm->RunMessageLoop();
        
        // Cleanup
        CleanupApplication();
        
        return result;
        
    } catch (const std::exception& e) {
        ShowError("Application error: " + std::string(e.what()));
        CleanupApplication();
        return -1;
    } catch (...) {
        ShowError("Unknown application error");
        CleanupApplication();
        return -1;
    }
}

void InitializeApplication() {
    // Initialize common controls
    INITCOMMONCONTROLSEX icex;
    icex.dwSize = sizeof(INITCOMMONCONTROLSEX);
    icex.dwICC = ICC_WIN95_CLASSES | ICC_TAB_CLASSES | ICC_BAR_CLASSES | 
                 ICC_LISTVIEW_CLASSES | ICC_TREEVIEW_CLASSES | ICC_BUTTON_CLASSES |
                 ICC_EDIT_CLASSES | ICC_STATIC_CLASSES | ICC_COMBOBOX_CLASSES |
                 ICC_PROGRESS_CLASS | ICC_HOTKEY_CLASS | ICC_ANIMATE_CLASS |
                 ICC_DATE_CLASSES | ICC_USEREX_CLASSES | ICC_LINK_CLASS;
    
    if (!InitCommonControlsEx(&icex)) {
        throw std::runtime_error("Failed to initialize common controls");
    }
    
    // Set process DPI awareness (for high DPI displays)
    SetProcessDPIAware();
    
    // Set application icon
    // LoadIcon(g_hInstance, MAKEINTRESOURCE(IDI_APPLICATION));
}

void CleanupApplication() {
    // Cleanup main form
    g_mainForm.reset();
    
    // Additional cleanup if needed
}

void ShowError(const std::string& message) {
    std::wstring wmessage(message.begin(), message.end());
    MessageBox(nullptr, wmessage.c_str(), L"Web4 API Bridge Demo Error", 
               MB_OK | MB_ICONERROR);
}

// Additional utility functions for debugging
#ifdef _DEBUG
void DebugOutput(const std::string& message) {
    OutputDebugStringA(message.c_str());
}

void DebugOutputW(const std::wstring& message) {
    OutputDebugStringW(message.c_str());
}
#endif 