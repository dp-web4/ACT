#pragma once

#include <string>
#include <iostream>

class DemoUI {
public:
    DemoUI() = default;
    ~DemoUI() = default;

    void showMainMenu(bool grpcAvailable);
    int getUserChoice();
    void showHeader();
    void showFooter();
    void clearScreen();
    void showLoading(const std::string& message);
    void showSuccess(const std::string& message);
    void showError(const std::string& message);
    void showInfo(const std::string& message);
    void showProgressBar(int current, int total, const std::string& label = "");
    
    // Specific menu sections
    void showAccountMenu();
    void showComponentMenu();
    void showPrivacyMenu();
    void showLCTMenu();
    void showPairingMenu();
    void showPairingQueueMenu();
    void showTrustMenu();
    void showEnergyMenu();
    void showPerformanceMenu();
    
    // Input helpers
    std::string getStringInput(const std::string& prompt);
    int getIntInput(const std::string& prompt, int min = 0, int max = 999);
    double getDoubleInput(const std::string& prompt, double min = 0.0, double max = 1000.0);
    bool getYesNoInput(const std::string& prompt);
    
    // Display helpers
    void displayAccount(const std::string& name, const std::string& address, const std::string& keyType);
    void displayComponent(const std::string& id, const std::string& data, const std::string& status);
    void displayAnonymousComponent(const std::string& componentHash, const std::string& manufacturerHash, const std::string& categoryHash, const std::string& status);
    void displayLCT(const std::string& id, const std::string& componentA, const std::string& componentB, const std::string& status);
    void displayPairing(const std::string& challengeId, const std::string& componentA, const std::string& componentB, const std::string& status);
    void displayPairingRequest(const std::string& requestId, const std::string& componentA, const std::string& componentB, const std::string& status);
    void displayTrustTensor(const std::string& id, double score, const std::string& status);
    void displayEnergyOperation(const std::string& id, const std::string& type, double amount, const std::string& status);
    void displayBatteryStatus(const std::string& componentId, double voltage, double current, double temperature, double soc, const std::string& status);

private:
    void printSeparator(char character = '-', int length = 50);
    void printCentered(const std::string& text, int width = 50);
    void printLeftAligned(const std::string& text, int width = 50);
    void printRightAligned(const std::string& text, int width = 50);
    std::string formatBytes(size_t bytes);
    std::string formatDuration(int64_t milliseconds);
    std::string formatPercentage(double value);
}; 