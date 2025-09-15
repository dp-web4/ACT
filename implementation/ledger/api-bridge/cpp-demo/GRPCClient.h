#pragma once

#include <string>
#include <vector>
#include <memory>
#include <functional>
#include <thread>
#include <atomic>

// Forward declarations for gRPC
namespace grpc {
    class Channel;
}

// Data structures for API responses (same as REST for consistency)
struct Account {
    std::string name;
    std::string address;
    std::string keyType;
};

struct ComponentRegistrationResult {
    std::string componentId;
    std::string componentIdentity;
    std::string componentData;
    std::string context;
    std::string creator;
    std::string lctId;
    std::string status;
    std::string txHash;
};

struct LCTResult {
    std::string lctId;
    std::string componentA;
    std::string componentB;
    std::string context;
    std::string proxyId;
    std::string status;
    int64_t createdAt;
    std::string creator;
    std::string txHash;
    std::string lctKeyHalf;
    std::string deviceKeyHalf;
};

struct PairingInitiateResult {
    std::string challengeId;
    std::string componentA;
    std::string componentB;
    std::string operationalContext;
    std::string proxyId;
    bool forceImmediate;
    std::string status;
    int64_t createdAt;
    std::string creator;
    std::string txHash;
};

struct PairingCompleteResult {
    std::string lctId;
    std::string sessionKeys;
    std::string trustSummary;
    std::string txHash;
    std::string splitKeyA;
    std::string splitKeyB;
};

struct TrustTensorResult {
    std::string tensorId;
    double score;
    std::string status;
    std::string txHash;
};

struct EnergyOperationResult {
    std::string operationId;
    std::string operationType;
    double amount;
    std::string status;
    std::string txHash;
};

struct BatteryStatusUpdate {
    std::string componentId;
    double voltage;
    double current;
    double temperature;
    double stateOfCharge;
    std::string status;
    int64_t timestamp;
};

class GRPCClient {
private:
    std::shared_ptr<grpc::Channel> channel;
    std::string serverAddress;
    std::atomic<bool> streamingActive;
    std::unique_ptr<std::thread> streamingThread;

public:
    GRPCClient(const std::string& endpoint);
    ~GRPCClient();

    // Account Management
    std::vector<Account> getAccounts();
    Account createAccount(const std::string& name);

    // Component Registry
    ComponentRegistrationResult registerComponent(const std::string& creator, 
                                                 const std::string& componentData, 
                                                 const std::string& context);
    ComponentRegistrationResult getComponent(const std::string& componentId);
    ComponentRegistrationResult getComponentIdentity(const std::string& componentId);
    ComponentRegistrationResult verifyComponent(const std::string& verifier, 
                                               const std::string& componentId, 
                                               const std::string& context);

    // LCT Management
    LCTResult createLCT(const std::string& creator, 
                       const std::string& componentA, 
                       const std::string& componentB, 
                       const std::string& context, 
                       const std::string& proxyId);
    LCTResult getLCT(const std::string& lctId);
    LCTResult updateLCTStatus(const std::string& creator, 
                             const std::string& lctId, 
                             const std::string& status, 
                             const std::string& context);

    // Pairing
    PairingInitiateResult initiatePairing(const std::string& creator, 
                                         const std::string& componentA, 
                                         const std::string& componentB, 
                                         const std::string& operationalContext, 
                                         const std::string& proxyId, 
                                         bool forceImmediate);
    PairingCompleteResult completePairing(const std::string& creator, 
                                         const std::string& challengeId, 
                                         const std::string& componentAAuth, 
                                         const std::string& componentBAuth, 
                                         const std::string& sessionContext);
    std::string revokePairing(const std::string& creator, 
                             const std::string& lctId, 
                             const std::string& reason, 
                             bool notifyOffline);
    std::string getPairingStatus(const std::string& challengeId);

    // Trust Tensor
    TrustTensorResult createTrustTensor(const std::string& creator, 
                                       const std::string& componentA, 
                                       const std::string& componentB, 
                                       const std::string& context, 
                                       double initialScore);
    TrustTensorResult getTrustTensor(const std::string& tensorId);
    TrustTensorResult updateTrustScore(const std::string& creator, 
                                      const std::string& tensorId, 
                                      double score, 
                                      const std::string& context);

    // Energy Operations
    EnergyOperationResult createEnergyOperation(const std::string& creator, 
                                               const std::string& componentA, 
                                               const std::string& componentB, 
                                               const std::string& operationType, 
                                               double amount, 
                                               const std::string& context);
    EnergyOperationResult executeEnergyTransfer(const std::string& creator, 
                                               const std::string& operationId, 
                                               double amount, 
                                               const std::string& context);
    double getEnergyBalance(const std::string& componentId);

    // Streaming for real-time updates
    void streamBatteryStatus(const std::string& componentId, 
                            int updateIntervalSeconds,
                            std::function<void(const BatteryStatusUpdate&)> callback);
    void stopStreaming();

    // Health check
    bool isConnected() const;

private:
    std::string makeRequest(const std::string& service, 
                           const std::string& method, 
                           const std::string& request);
    void streamingWorker(const std::string& componentId, 
                        int updateIntervalSeconds,
                        std::function<void(const BatteryStatusUpdate&)> callback);
}; 