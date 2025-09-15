#include <iostream>
#include <string>
#include <vector>
#include <memory>
#include <chrono>
#include <thread>

// Include our API interfaces
#include "RESTClient.h"
#include "GRPCClient.h"
#include "DemoUI.h"

using namespace std;

class APIBridgeDemo {
private:
    unique_ptr<RESTClient> restClient;
    unique_ptr<GRPCClient> grpcClient;
    unique_ptr<DemoUI> ui;
    
    string restEndpoint;
    string grpcEndpoint;
    bool grpcAvailable;

public:
    APIBridgeDemo() : restEndpoint("http://localhost:8080"), 
                      grpcEndpoint("localhost:9092"),
                      grpcAvailable(false) {
        ui = make_unique<DemoUI>();
    }

    void initialize() {
        cout << "=== Web4 Race Car Battery Management API Bridge Demo ===" << endl;
        cout << "Initializing clients..." << endl;
        
        // Initialize REST client
        restClient = make_unique<RESTClient>(restEndpoint);
        
        // Try to initialize gRPC client
        try {
            grpcClient = make_unique<GRPCClient>(grpcEndpoint);
            grpcAvailable = true;
            cout << "✓ gRPC client initialized successfully" << endl;
        } catch (const exception& e) {
            cout << "⚠ gRPC client not available: " << e.what() << endl;
            cout << "   Only REST interface will be available" << endl;
        }
        
        cout << "✓ REST client initialized successfully" << endl;
        cout << endl;
    }

    void run() {
        while (true) {
            ui->showMainMenu(grpcAvailable);
            int choice = ui->getUserChoice();
            
            switch (choice) {
                case 1:
                    testAccountManagement();
                    break;
                case 2:
                    testComponentRegistry();
                    break;
                case 3:
                    testPrivacyFeatures();
                    break;
                case 4:
                    testLCTManagement();
                    break;
                case 5:
                    testPairingProcess();
                    break;
                case 6:
                    testPairingQueue();
                    break;
                case 7:
                    testTrustTensor();
                    break;
                case 8:
                    testEnergyOperations();
                    break;
                case 9:
                    if (grpcAvailable) {
                        testStreaming();
                    } else {
                        cout << "gRPC not available for streaming" << endl;
                    }
                    break;
                case 10:
                    comparePerformance();
                    break;
                case 11:
                    showSystemInfo();
                    break;
                case 0:
                    cout << "Exiting demo..." << endl;
                    return;
                default:
                    cout << "Invalid choice. Please try again." << endl;
            }
            
            cout << "\nPress Enter to continue...";
            cin.ignore(numeric_limits<streamsize>::max(), '\n');
            cin.get();
        }
    }

private:
    void testAccountManagement() {
        cout << "\n=== Account Management Test ===" << endl;
        
        // Test REST
        cout << "\n--- REST API Test ---" << endl;
        try {
            auto accounts = restClient->getAccounts();
            cout << "REST: Found " << accounts.size() << " accounts" << endl;
            for (const auto& account : accounts) {
                cout << "  - " << account.name << " (" << account.keyType << "): " << account.address << endl;
            }
        } catch (const exception& e) {
            cout << "REST Error: " << e.what() << endl;
        }
        
        // Test gRPC
        if (grpcAvailable) {
            cout << "\n--- gRPC API Test ---" << endl;
            try {
                auto accounts = grpcClient->getAccounts();
                cout << "gRPC: Found " << accounts.size() << " accounts" << endl;
                for (const auto& account : accounts) {
                    cout << "  - " << account.name << " (" << account.keyType << "): " << account.address << endl;
                }
            } catch (const exception& e) {
                cout << "gRPC Error: " << e.what() << endl;
            }
        }
    }

    void testComponentRegistry() {
        cout << "\n=== Component Registry Test (Legacy) ===" << endl;
        
        string creator = "demo-user";
        string componentData = "demo-battery-module-v1.0";
        string context = "demo-context";
        
        // Test REST
        cout << "\n--- REST API Test ---" << endl;
        try {
            auto result = restClient->registerComponent(creator, componentData, context);
            cout << "REST: Component registered successfully" << endl;
            cout << "  Component ID: " << result.componentId << endl;
            cout << "  Transaction Hash: " << result.txHash << endl;
            cout << "  Status: " << result.status << endl;
        } catch (const exception& e) {
            cout << "REST Error: " << e.what() << endl;
        }
        
        // Test gRPC
        if (grpcAvailable) {
            cout << "\n--- gRPC API Test ---" << endl;
            try {
                auto result = grpcClient->registerComponent(creator, componentData, context);
                cout << "gRPC: Component registered successfully" << endl;
                cout << "  Component ID: " << result.componentId << endl;
                cout << "  Transaction Hash: " << result.txHash << endl;
                cout << "  Status: " << result.status << endl;
            } catch (const exception& e) {
                cout << "gRPC Error: " << e.what() << endl;
            }
        }
    }

    void testPrivacyFeatures() {
        cout << "\n=== Privacy-Focused Features Test ===" << endl;
        
        string creator = "demo-user";
        string realComponentId = "battery-module-001";
        string manufacturerId = "tesla-motors";
        string componentType = "lithium-ion-battery";
        string context = "race-car-demo";
        
        // Test REST
        cout << "\n--- REST API Test ---" << endl;
        try {
            // 1. Register anonymous component
            cout << "1. Registering anonymous component..." << endl;
            auto anonResult = restClient->registerAnonymousComponent(creator, realComponentId, manufacturerId, componentType, context);
            cout << "   Component Hash: " << anonResult.componentHash << endl;
            cout << "   Manufacturer Hash: " << anonResult.manufacturerHash << endl;
            cout << "   Category Hash: " << anonResult.categoryHash << endl;
            cout << "   Transaction Hash: " << anonResult.txHash << endl;
            
            // 2. Verify pairing with hashes
            cout << "\n2. Verifying pairing with hashes..." << endl;
            auto verifyResult = restClient->verifyComponentPairingWithHashes("verifier-001", anonResult.componentHash, "motor-hash-001", context);
            cout << "   Source Hash: " << verifyResult.sourceHash << endl;
            cout << "   Target Hash: " << verifyResult.targetHash << endl;
            cout << "   Status: " << verifyResult.status << endl;
            
            // 3. Create anonymous pairing authorization
            cout << "\n3. Creating anonymous pairing authorization..." << endl;
            auto authResult = restClient->createAnonymousPairingAuthorization(creator, anonResult.componentHash, "motor-hash-001", context);
            cout << "   Authorization ID: " << authResult.authorizationId << endl;
            cout << "   Source Hash: " << authResult.sourceHash << endl;
            cout << "   Target Hash: " << authResult.targetHash << endl;
            cout << "   Status: " << authResult.status << endl;
            
            // 4. Create revocation event
            cout << "\n4. Creating revocation event..." << endl;
            auto revokeResult = restClient->createAnonymousRevocationEvent(creator, anonResult.componentHash, "component-failure", context);
            cout << "   Revocation ID: " << revokeResult.revocationId << endl;
            cout << "   Component Hash: " << revokeResult.componentHash << endl;
            cout << "   Reason: " << revokeResult.reason << endl;
            cout << "   Status: " << revokeResult.status << endl;
            
            // 5. Get anonymous component metadata
            cout << "\n5. Getting anonymous component metadata..." << endl;
            auto metadataResult = restClient->getAnonymousComponentMetadata(anonResult.componentHash);
            cout << "   Component Hash: " << metadataResult.componentHash << endl;
            cout << "   Metadata: " << metadataResult.metadata << endl;
            cout << "   Status: " << metadataResult.status << endl;
            
        } catch (const exception& e) {
            cout << "REST Error: " << e.what() << endl;
        }
        
        // Test gRPC (if available)
        if (grpcAvailable) {
            cout << "\n--- gRPC API Test ---" << endl;
            try {
                // Similar privacy tests for gRPC would go here
                cout << "gRPC privacy features test would be implemented here" << endl;
            } catch (const exception& e) {
                cout << "gRPC Error: " << e.what() << endl;
            }
        }
    }

    void testLCTManagement() {
        cout << "\n=== LCT Management Test ===" << endl;
        
        string creator = "demo-user";
        string componentA = "battery-001";
        string componentB = "motor-001";
        string context = "race-car-pairing";
        string proxyId = "proxy-001";
        
        // Test REST
        cout << "\n--- REST API Test ---" << endl;
        try {
            auto result = restClient->createLCT(creator, componentA, componentB, context, proxyId);
            cout << "REST: LCT created successfully" << endl;
            cout << "  LCT ID: " << result.lctId << endl;
            cout << "  Transaction Hash: " << result.txHash << endl;
            cout << "  LCT Key Half: " << result.lctKeyHalf << endl;
            cout << "  Device Key Half: " << result.deviceKeyHalf << endl;
        } catch (const exception& e) {
            cout << "REST Error: " << e.what() << endl;
        }
        
        // Test gRPC
        if (grpcAvailable) {
            cout << "\n--- gRPC API Test ---" << endl;
            try {
                auto result = grpcClient->createLCT(creator, componentA, componentB, context, proxyId);
                cout << "gRPC: LCT created successfully" << endl;
                cout << "  LCT ID: " << result.lctId << endl;
                cout << "  Transaction Hash: " << result.txHash << endl;
                cout << "  LCT Key Half: " << result.lctKeyHalf << endl;
                cout << "  Device Key Half: " << result.deviceKeyHalf << endl;
            } catch (const exception& e) {
                cout << "gRPC Error: " << e.what() << endl;
            }
        }
    }

    void testPairingProcess() {
        cout << "\n=== Pairing Process Test ===" << endl;
        
        string creator = "demo-user";
        string componentA = "battery-001";
        string componentB = "motor-001";
        string operationalContext = "race-car-operation";
        string proxyId = "proxy-001";
        
        // Test REST
        cout << "\n--- REST API Test ---" << endl;
        try {
            // Step 1: Initiate pairing
            auto initiateResult = restClient->initiatePairing(creator, componentA, componentB, operationalContext, proxyId, false);
            cout << "REST: Pairing initiated" << endl;
            cout << "  Challenge ID: " << initiateResult.challengeId << endl;
            cout << "  Transaction Hash: " << initiateResult.txHash << endl;
            
            // Step 2: Complete pairing
            auto completeResult = restClient->completePairing(creator, initiateResult.challengeId, "battery-auth", "motor-auth", "demo-session");
            cout << "REST: Pairing completed" << endl;
            cout << "  LCT ID: " << completeResult.lctId << endl;
            cout << "  Transaction Hash: " << completeResult.txHash << endl;
            cout << "  Split Key A: " << completeResult.splitKeyA << endl;
            cout << "  Split Key B: " << completeResult.splitKeyB << endl;
        } catch (const exception& e) {
            cout << "REST Error: " << e.what() << endl;
        }
        
        // Test gRPC
        if (grpcAvailable) {
            cout << "\n--- gRPC API Test ---" << endl;
            try {
                // Step 1: Initiate pairing
                auto initiateResult = grpcClient->initiatePairing(creator, componentA, componentB, operationalContext, proxyId, false);
                cout << "gRPC: Pairing initiated" << endl;
                cout << "  Challenge ID: " << initiateResult.challengeId << endl;
                cout << "  Transaction Hash: " << initiateResult.txHash << endl;
                
                // Step 2: Complete pairing
                auto completeResult = grpcClient->completePairing(creator, initiateResult.challengeId, "battery-auth", "motor-auth", "demo-session");
                cout << "gRPC: Pairing completed" << endl;
                cout << "  LCT ID: " << completeResult.lctId << endl;
                cout << "  Transaction Hash: " << completeResult.txHash << endl;
                cout << "  Split Key A: " << completeResult.splitKeyA << endl;
                cout << "  Split Key B: " << completeResult.splitKeyB << endl;
            } catch (const exception& e) {
                cout << "gRPC Error: " << e.what() << endl;
            }
        }
    }

    void testPairingQueue() {
        cout << "\n=== Pairing Queue Test ===" << endl;
        
        string creator = "demo-user";
        string componentA = "battery-001";
        string componentB = "motor-001";
        string context = "race-car-queue";
        
        // Test REST
        cout << "\n--- REST API Test ---" << endl;
        try {
            // 1. Queue pairing request
            cout << "1. Queuing pairing request..." << endl;
            auto queueResult = restClient->queuePairingRequest(creator, componentA, componentB, context);
            cout << "   Request ID: " << queueResult.requestId << endl;
            cout << "   Component A: " << queueResult.componentA << endl;
            cout << "   Component B: " << queueResult.componentB << endl;
            cout << "   Status: " << queueResult.status << endl;
            cout << "   Transaction Hash: " << queueResult.txHash << endl;
            
            // 2. Get queue status
            cout << "\n2. Getting queue status..." << endl;
            auto statusResult = restClient->getQueueStatus("default-queue");
            cout << "   Queue ID: " << statusResult.queueId << endl;
            cout << "   Pending Requests: " << statusResult.pendingRequests << endl;
            cout << "   Processed Requests: " << statusResult.processedRequests << endl;
            cout << "   Status: " << statusResult.status << endl;
            
            // 3. Get queued requests
            cout << "\n3. Getting queued requests..." << endl;
            auto requests = restClient->getQueuedRequests("default-queue");
            cout << "   Found " << requests.size() << " queued requests" << endl;
            for (const auto& req : requests) {
                cout << "     - Request ID: " << req.requestId << " (" << req.status << ")" << endl;
            }
            
            // 4. List proxy queue
            cout << "\n4. Listing proxy queue..." << endl;
            auto proxyRequests = restClient->listProxyQueue("proxy-001");
            cout << "   Found " << proxyRequests.size() << " proxy requests" << endl;
            for (const auto& req : proxyRequests) {
                cout << "     - Request ID: " << req.requestId << " (" << req.status << ")" << endl;
            }
            
            // 5. Process offline queue
            cout << "\n5. Processing offline queue..." << endl;
            auto processResult = restClient->processOfflineQueue("processor-001", "default-queue", context);
            cout << "   Process Result: " << processResult << endl;
            
            // 6. Cancel request
            cout << "\n6. Canceling request..." << endl;
            auto cancelResult = restClient->cancelRequest(creator, queueResult.requestId, "user-cancellation");
            cout << "   Cancel Result: " << cancelResult << endl;
            
        } catch (const exception& e) {
            cout << "REST Error: " << e.what() << endl;
        }
        
        // Test gRPC (if available)
        if (grpcAvailable) {
            cout << "\n--- gRPC API Test ---" << endl;
            try {
                cout << "gRPC pairing queue test would be implemented here" << endl;
            } catch (const exception& e) {
                cout << "gRPC Error: " << e.what() << endl;
            }
        }
    }

    void testTrustTensor() {
        cout << "\n=== Trust Tensor Test ===" << endl;
        
        string creator = "demo-user";
        string componentA = "battery-001";
        string componentB = "motor-001";
        string context = "race-car-trust";
        double initialScore = 0.8;
        
        // Test REST
        cout << "\n--- REST API Test ---" << endl;
        try {
            auto result = restClient->createTrustTensor(creator, componentA, componentB, context, initialScore);
            cout << "REST: Trust tensor created successfully" << endl;
            cout << "  Tensor ID: " << result.tensorId << endl;
            cout << "  Initial Score: " << result.score << endl;
            cout << "  Transaction Hash: " << result.txHash << endl;
            cout << "  Status: " << result.status << endl;
        } catch (const exception& e) {
            cout << "REST Error: " << e.what() << endl;
        }
        
        // Test gRPC
        if (grpcAvailable) {
            cout << "\n--- gRPC API Test ---" << endl;
            try {
                auto result = grpcClient->createTrustTensor(creator, componentA, componentB, context, initialScore);
                cout << "gRPC: Trust tensor created successfully" << endl;
                cout << "  Tensor ID: " << result.tensorId << endl;
                cout << "  Initial Score: " << result.score << endl;
                cout << "  Transaction Hash: " << result.txHash << endl;
                cout << "  Status: " << result.status << endl;
            } catch (const exception& e) {
                cout << "gRPC Error: " << e.what() << endl;
            }
        }
    }

    void testEnergyOperations() {
        cout << "\n=== Energy Operations Test ===" << endl;
        
        string creator = "demo-user";
        string componentA = "battery-001";
        string componentB = "motor-001";
        string operationType = "energy-transfer";
        double amount = 100.5;
        string context = "race-car-energy";
        
        // Test REST
        cout << "\n--- REST API Test ---" << endl;
        try {
            auto result = restClient->createEnergyOperation(creator, componentA, componentB, operationType, amount, context);
            cout << "REST: Energy operation created successfully" << endl;
            cout << "  Operation ID: " << result.operationId << endl;
            cout << "  Operation Type: " << result.operationType << endl;
            cout << "  Amount: " << result.amount << endl;
            cout << "  Transaction Hash: " << result.txHash << endl;
            cout << "  Status: " << result.status << endl;
        } catch (const exception& e) {
            cout << "REST Error: " << e.what() << endl;
        }
        
        // Test gRPC
        if (grpcAvailable) {
            cout << "\n--- gRPC API Test ---" << endl;
            try {
                auto result = grpcClient->createEnergyOperation(creator, componentA, componentB, operationType, amount, context);
                cout << "gRPC: Energy operation created successfully" << endl;
                cout << "  Operation ID: " << result.operationId << endl;
                cout << "  Operation Type: " << result.operationType << endl;
                cout << "  Amount: " << result.amount << endl;
                cout << "  Transaction Hash: " << result.txHash << endl;
                cout << "  Status: " << result.status << endl;
            } catch (const exception& e) {
                cout << "gRPC Error: " << e.what() << endl;
            }
        }
    }

    void testStreaming() {
        cout << "\n=== Real-time Streaming Test (gRPC) ===" << endl;
        
        if (!grpcAvailable) {
            cout << "gRPC not available for streaming" << endl;
            return;
        }
        
        try {
            cout << "Starting battery status stream for 10 seconds..." << endl;
            
            auto startTime = chrono::steady_clock::now();
            bool streamActive = true;
            
            grpcClient->streamBatteryStatus("battery-001", 10, [&](const BatteryStatusUpdate& update) {
                auto now = chrono::steady_clock::now();
                auto elapsed = chrono::duration_cast<chrono::seconds>(now - startTime).count();
                
                cout << "[" << elapsed << "s] Battery Status Update:" << endl;
                cout << "  Component ID: " << update.componentId << endl;
                cout << "  Voltage: " << update.voltage << "V" << endl;
                cout << "  Current: " << update.current << "A" << endl;
                cout << "  Temperature: " << update.temperature << "°C" << endl;
                cout << "  State of Charge: " << update.stateOfCharge << "%" << endl;
                cout << "  Status: " << update.status << endl;
                cout << "  Timestamp: " << update.timestamp << endl;
                cout << endl;
                
                if (elapsed >= 10) {
                    streamActive = false;
                }
            });
            
            while (streamActive) {
                this_thread::sleep_for(chrono::milliseconds(100));
            }
            
            cout << "Streaming test completed." << endl;
            
        } catch (const exception& e) {
            cout << "Streaming Error: " << e.what() << endl;
        }
    }

    void comparePerformance() {
        cout << "\n=== Performance Comparison Test ===" << endl;
        
        const int iterations = 10;
        string creator = "perf-test-user";
        string componentData = "perf-test-component";
        string context = "perf-test-context";
        
        // Test REST performance
        cout << "\n--- REST API Performance ---" << endl;
        auto restStart = chrono::high_resolution_clock::now();
        
        for (int i = 0; i < iterations; i++) {
            try {
                restClient->registerComponent(creator + to_string(i), componentData + to_string(i), context);
            } catch (const exception& e) {
                cout << "REST iteration " << i << " failed: " << e.what() << endl;
            }
        }
        
        auto restEnd = chrono::high_resolution_clock::now();
        auto restDuration = chrono::duration_cast<chrono::milliseconds>(restEnd - restStart).count();
        
        cout << "REST: " << iterations << " operations completed in " << restDuration << "ms" << endl;
        cout << "REST: Average " << (double)restDuration / iterations << "ms per operation" << endl;
        
        // Test gRPC performance
        if (grpcAvailable) {
            cout << "\n--- gRPC API Performance ---" << endl;
            auto grpcStart = chrono::high_resolution_clock::now();
            
            for (int i = 0; i < iterations; i++) {
                try {
                    grpcClient->registerComponent(creator + to_string(i), componentData + to_string(i), context);
                } catch (const exception& e) {
                    cout << "gRPC iteration " << i << " failed: " << e.what() << endl;
                }
            }
            
            auto grpcEnd = chrono::high_resolution_clock::now();
            auto grpcDuration = chrono::duration_cast<chrono::milliseconds>(grpcEnd - grpcStart).count();
            
            cout << "gRPC: " << iterations << " operations completed in " << grpcDuration << "ms" << endl;
            cout << "gRPC: Average " << (double)grpcDuration / iterations << "ms per operation" << endl;
            
            // Compare performance
            cout << "\n--- Performance Comparison ---" << endl;
            if (restDuration < grpcDuration) {
                cout << "REST is " << (double)grpcDuration / restDuration << "x faster than gRPC" << endl;
            } else {
                cout << "gRPC is " << (double)restDuration / grpcDuration << "x faster than REST" << endl;
            }
        }
    }

    void showSystemInfo() {
        cout << "\n=== System Information ===" << endl;
        
        cout << "API Bridge Configuration:" << endl;
        cout << "  REST Endpoint: " << restEndpoint << endl;
        cout << "  gRPC Endpoint: " << grpcEndpoint << endl;
        cout << "  gRPC Available: " << (grpcAvailable ? "Yes" : "No") << endl;
        
        cout << "\nTesting API Bridge Health..." << endl;
        try {
            string healthStatus = restClient->getHealthStatus();
            cout << "  Health Status: " << healthStatus << endl;
        } catch (const exception& e) {
            cout << "  Health Check Failed: " << e.what() << endl;
        }
        
        cout << "\nTesting Blockchain Status..." << endl;
        try {
            string blockchainStatus = restClient->getBlockchainStatus();
            cout << "  Blockchain Status: " << blockchainStatus << endl;
        } catch (const exception& e) {
            cout << "  Blockchain Status Check Failed: " << e.what() << endl;
        }
        
        cout << "\nAvailable Features:" << endl;
        cout << "  ✓ Account Management" << endl;
        cout << "  ✓ Component Registry (Legacy)" << endl;
        cout << "  ✓ Privacy-Focused Component Operations" << endl;
        cout << "  ✓ LCT Management" << endl;
        cout << "  ✓ Pairing Process" << endl;
        cout << "  ✓ Pairing Queue Operations" << endl;
        cout << "  ✓ Trust Tensor Operations" << endl;
        cout << "  ✓ Energy Operations" << endl;
        cout << "  " << (grpcAvailable ? "✓" : "✗") << " Real-time Streaming (gRPC)" << endl;
        cout << "  ✓ Performance Testing" << endl;
    }
};

int main() {
    try {
        APIBridgeDemo demo;
        demo.initialize();
        demo.run();
    } catch (const exception& e) {
        cerr << "Fatal error: " << e.what() << endl;
        return 1;
    }
    
    return 0;
} 