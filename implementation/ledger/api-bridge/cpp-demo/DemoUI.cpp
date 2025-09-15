#include "DemoUI.h"
#include <iostream>
#include <iomanip>
#include <sstream>
#include <limits>
#include <algorithm>

#ifdef _WIN32
#include <windows.h>
#else
#include <cstdlib>
#endif

void DemoUI::showMainMenu(bool grpcAvailable) {
    clearScreen();
    showHeader();
    
    std::cout << "\n=== Main Menu ===" << std::endl;
    std::cout << "1. Account Management" << std::endl;
    std::cout << "2. Component Registry (Legacy)" << std::endl;
    std::cout << "3. Privacy-Focused Features" << std::endl;
    std::cout << "4. LCT Management" << std::endl;
    std::cout << "5. Pairing Process" << std::endl;
    std::cout << "6. Pairing Queue Operations" << std::endl;
    std::cout << "7. Trust Tensor" << std::endl;
    std::cout << "8. Energy Operations" << std::endl;
    
    if (grpcAvailable) {
        std::cout << "9. Real-time Streaming (gRPC)" << std::endl;
    }
    
    std::cout << "10. Performance Comparison" << std::endl;
    std::cout << "11. System Information" << std::endl;
    std::cout << "0. Exit" << std::endl;
    
    printSeparator();
    std::cout << "gRPC Available: " << (grpcAvailable ? "✓ Yes" : "✗ No") << std::endl;
    std::cout << "Privacy Features: ✓ Enabled" << std::endl;
    std::cout << "Real Blockchain: ✓ Connected" << std::endl;
    printSeparator();
}

int DemoUI::getUserChoice() {
    std::cout << "\nEnter your choice (0-11): ";
    int choice;
    
    while (!(std::cin >> choice) || choice < 0 || choice > 11) {
        std::cout << "Invalid choice. Please enter a number between 0 and 11: ";
        std::cin.clear();
        std::cin.ignore(std::numeric_limits<std::streamsize>::max(), '\n');
    }
    
    return choice;
}

void DemoUI::showHeader() {
    printSeparator('=', 60);
    printCentered("Web4 Race Car Battery Management API Bridge Demo", 60);
    printCentered("C++ Client Reference Implementation", 60);
    printCentered("Privacy-Focused Features Enabled", 60);
    printCentered("Compatible with RAD Studio", 60);
    printSeparator('=', 60);
}

void DemoUI::showFooter() {
    printSeparator('=', 60);
    printCentered("Press Enter to continue...", 60);
    printSeparator('=', 60);
}

void DemoUI::clearScreen() {
#ifdef _WIN32
    system("cls");
#else
    system("clear");
#endif
}

void DemoUI::showLoading(const std::string& message) {
    std::cout << message << " ";
    std::cout.flush();
}

void DemoUI::showSuccess(const std::string& message) {
    std::cout << "✓ " << message << std::endl;
}

void DemoUI::showError(const std::string& message) {
    std::cout << "✗ " << message << std::endl;
}

void DemoUI::showInfo(const std::string& message) {
    std::cout << "ℹ " << message << std::endl;
}

void DemoUI::showProgressBar(int current, int total, const std::string& label) {
    const int barWidth = 50;
    float progress = static_cast<float>(current) / total;
    int pos = static_cast<int>(barWidth * progress);
    
    std::cout << "\r" << label << " [";
    for (int i = 0; i < barWidth; ++i) {
        if (i < pos) std::cout << "=";
        else if (i == pos) std::cout << ">";
        else std::cout << " ";
    }
    std::cout << "] " << static_cast<int>(progress * 100.0) << "%" << std::flush;
    
    if (current == total) {
        std::cout << std::endl;
    }
}

void DemoUI::showAccountMenu() {
    std::cout << "\n=== Account Management ===" << std::endl;
    std::cout << "1. List Accounts" << std::endl;
    std::cout << "2. Create Account" << std::endl;
    std::cout << "3. Get Account Details" << std::endl;
    std::cout << "0. Back to Main Menu" << std::endl;
}

void DemoUI::showComponentMenu() {
    std::cout << "\n=== Component Registry (Legacy) ===" << std::endl;
    std::cout << "1. Register Component" << std::endl;
    std::cout << "2. Get Component" << std::endl;
    std::cout << "3. Get Component Identity" << std::endl;
    std::cout << "4. Verify Component" << std::endl;
    std::cout << "0. Back to Main Menu" << std::endl;
}

void DemoUI::showPrivacyMenu() {
    std::cout << "\n=== Privacy-Focused Features ===" << std::endl;
    std::cout << "1. Register Anonymous Component" << std::endl;
    std::cout << "2. Verify Pairing with Hashes" << std::endl;
    std::cout << "3. Create Anonymous Pairing Authorization" << std::endl;
    std::cout << "4. Create Anonymous Revocation Event" << std::endl;
    std::cout << "5. Get Anonymous Component Metadata" << std::endl;
    std::cout << "6. Full Privacy Demo" << std::endl;
    std::cout << "0. Back to Main Menu" << std::endl;
}

void DemoUI::showLCTMenu() {
    std::cout << "\n=== LCT Management ===" << std::endl;
    std::cout << "1. Create LCT" << std::endl;
    std::cout << "2. Get LCT" << std::endl;
    std::cout << "3. Update LCT Status" << std::endl;
    std::cout << "4. List LCTs" << std::endl;
    std::cout << "0. Back to Main Menu" << std::endl;
}

void DemoUI::showPairingMenu() {
    std::cout << "\n=== Pairing Process ===" << std::endl;
    std::cout << "1. Initiate Pairing" << std::endl;
    std::cout << "2. Complete Pairing" << std::endl;
    std::cout << "3. Revoke Pairing" << std::endl;
    std::cout << "4. Get Pairing Status" << std::endl;
    std::cout << "5. Full Pairing Flow Demo" << std::endl;
    std::cout << "0. Back to Main Menu" << std::endl;
}

void DemoUI::showPairingQueueMenu() {
    std::cout << "\n=== Pairing Queue Operations ===" << std::endl;
    std::cout << "1. Queue Pairing Request" << std::endl;
    std::cout << "2. Get Queue Status" << std::endl;
    std::cout << "3. Get Queued Requests" << std::endl;
    std::cout << "4. List Proxy Queue" << std::endl;
    std::cout << "5. Process Offline Queue" << std::endl;
    std::cout << "6. Cancel Request" << std::endl;
    std::cout << "7. Full Queue Demo" << std::endl;
    std::cout << "0. Back to Main Menu" << std::endl;
}

void DemoUI::showTrustMenu() {
    std::cout << "\n=== Trust Tensor ===" << std::endl;
    std::cout << "1. Create Trust Tensor" << std::endl;
    std::cout << "2. Get Trust Tensor" << std::endl;
    std::cout << "3. Update Trust Score" << std::endl;
    std::cout << "4. List Trust Tensors" << std::endl;
    std::cout << "0. Back to Main Menu" << std::endl;
}

void DemoUI::showEnergyMenu() {
    std::cout << "\n=== Energy Operations ===" << std::endl;
    std::cout << "1. Create Energy Operation" << std::endl;
    std::cout << "2. Execute Energy Transfer" << std::endl;
    std::cout << "3. Get Energy Balance" << std::endl;
    std::cout << "4. List Energy Operations" << std::endl;
    std::cout << "0. Back to Main Menu" << std::endl;
}

void DemoUI::showPerformanceMenu() {
    std::cout << "\n=== Performance Comparison ===" << std::endl;
    std::cout << "1. REST vs gRPC Speed Test" << std::endl;
    std::cout << "2. Concurrent Request Test" << std::endl;
    std::cout << "3. Memory Usage Test" << std::endl;
    std::cout << "4. Network Latency Test" << std::endl;
    std::cout << "0. Back to Main Menu" << std::endl;
}

std::string DemoUI::getStringInput(const std::string& prompt) {
    std::string input;
    std::cout << prompt << ": ";
    std::cin.ignore(std::numeric_limits<std::streamsize>::max(), '\n');
    std::getline(std::cin, input);
    return input;
}

int DemoUI::getIntInput(const std::string& prompt, int min, int max) {
    int value;
    do {
        std::cout << prompt << " (" << min << "-" << max << "): ";
        while (!(std::cin >> value)) {
            std::cout << "Invalid input. Please enter a number: ";
            std::cin.clear();
            std::cin.ignore(std::numeric_limits<std::streamsize>::max(), '\n');
        }
    } while (value < min || value > max);
    
    return value;
}

double DemoUI::getDoubleInput(const std::string& prompt, double min, double max) {
    double value;
    do {
        std::cout << prompt << " (" << min << "-" << max << "): ";
        while (!(std::cin >> value)) {
            std::cout << "Invalid input. Please enter a number: ";
            std::cin.clear();
            std::cin.ignore(std::numeric_limits<std::streamsize>::max(), '\n');
        }
    } while (value < min || value > max);
    
    return value;
}

bool DemoUI::getYesNoInput(const std::string& prompt) {
    std::string input;
    do {
        std::cout << prompt << " (y/n): ";
        std::cin >> input;
        std::transform(input.begin(), input.end(), input.begin(), ::tolower);
    } while (input != "y" && input != "n" && input != "yes" && input != "no");
    
    return (input == "y" || input == "yes");
}

void DemoUI::displayAccount(const std::string& name, const std::string& address, const std::string& keyType) {
    std::cout << std::left << std::setw(20) << name 
              << std::setw(45) << address 
              << std::setw(15) << keyType << std::endl;
}

void DemoUI::displayComponent(const std::string& id, const std::string& data, const std::string& status) {
    std::cout << std::left << std::setw(25) << id 
              << std::setw(30) << data 
              << std::setw(15) << status << std::endl;
}

void DemoUI::displayLCT(const std::string& id, const std::string& componentA, const std::string& componentB, const std::string& status) {
    std::cout << std::left << std::setw(35) << id 
              << std::setw(20) << componentA 
              << std::setw(20) << componentB 
              << std::setw(15) << status << std::endl;
}

void DemoUI::displayPairing(const std::string& challengeId, const std::string& componentA, const std::string& componentB, const std::string& status) {
    std::cout << std::left << std::setw(35) << challengeId 
              << std::setw(20) << componentA 
              << std::setw(20) << componentB 
              << std::setw(15) << status << std::endl;
}

void DemoUI::displayTrustTensor(const std::string& id, double score, const std::string& status) {
    std::cout << std::left << std::setw(35) << id 
              << std::setw(10) << std::fixed << std::setprecision(3) << score 
              << std::setw(15) << status << std::endl;
}

void DemoUI::displayEnergyOperation(const std::string& id, const std::string& type, double amount, const std::string& status) {
    std::cout << std::left << std::setw(35) << id 
              << std::setw(15) << type 
              << std::setw(10) << std::fixed << std::setprecision(2) << amount 
              << std::setw(15) << status << std::endl;
}

void DemoUI::displayBatteryStatus(const std::string& componentId, double voltage, double current, double temperature, double soc, const std::string& status) {
    std::cout << std::left << std::setw(20) << componentId 
              << std::setw(8) << std::fixed << std::setprecision(2) << voltage << "V"
              << std::setw(10) << std::fixed << std::setprecision(2) << current << "A"
              << std::setw(8) << std::fixed << std::setprecision(1) << temperature << "°C"
              << std::setw(8) << std::fixed << std::setprecision(1) << soc << "%"
              << std::setw(15) << status << std::endl;
}

void DemoUI::printSeparator(char character, int length) {
    std::cout << std::string(length, character) << std::endl;
}

void DemoUI::printCentered(const std::string& text, int width) {
    int padding = (width - text.length()) / 2;
    std::cout << std::string(padding, ' ') << text << std::endl;
}

void DemoUI::printLeftAligned(const std::string& text, int width) {
    std::cout << std::left << std::setw(width) << text << std::endl;
}

void DemoUI::printRightAligned(const std::string& text, int width) {
    std::cout << std::right << std::setw(width) << text << std::endl;
}

std::string DemoUI::formatBytes(size_t bytes) {
    const char* units[] = {"B", "KB", "MB", "GB", "TB"};
    int unit = 0;
    double size = static_cast<double>(bytes);
    
    while (size >= 1024.0 && unit < 4) {
        size /= 1024.0;
        unit++;
    }
    
    std::ostringstream oss;
    oss << std::fixed << std::setprecision(2) << size << " " << units[unit];
    return oss.str();
}

std::string DemoUI::formatDuration(int64_t milliseconds) {
    if (milliseconds < 1000) {
        return std::to_string(milliseconds) + "ms";
    } else if (milliseconds < 60000) {
        return std::to_string(milliseconds / 1000.0) + "s";
    } else {
        return std::to_string(milliseconds / 60000.0) + "m";
    }
}

std::string DemoUI::formatPercentage(double value) {
    std::ostringstream oss;
    oss << std::fixed << std::setprecision(2) << (value * 100.0) << "%";
    return oss.str();
} 