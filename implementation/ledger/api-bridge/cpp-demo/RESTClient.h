#pragma once

#include <string>
#include <vector>
#include <memory>
#include <functional>

// Forward declarations for HTTP client
namespace httplib {
    class Client;
}

// Data structures for API responses
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

// Privacy-focused component structures
struct AnonymousComponentResult {
    std::string componentHash;
    std::string manufacturerHash;
    std::string categoryHash;
    std::string creator;
    std::string context;
    std::string status;
    std::string txHash;
};

struct PairingVerificationResult {
    std::string sourceHash;
    std::string targetHash;
    std::string context;
    std::string status;
    std::string txHash;
};

struct PairingAuthorizationResult {
    std::string authorizationId;
    std::string sourceHash;
    std::string targetHash;
    std::string context;
    std::string status;
    std::string txHash;
};

struct RevocationEventResult {
    std::string revocationId;
    std::string componentHash;
    std::string reason;
    std::string context;
    std::string status;
    std::string txHash;
};

struct ComponentMetadataResult {
    std::string componentHash;
    std::string metadata;
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

// Pairing Queue structures
struct PairingRequestResult {
    std::string requestId;
    std::string componentA;
    std::string componentB;
    std::string context;
    std::string status;
    int64_t createdAt;
    std::string creator;
    std::string txHash;
};

struct QueueStatusResult {
    std::string queueId;
    int pendingRequests;
    int processedRequests;
    std::string status;
    std::string txHash;
};

class RESTClient {
private:
    std::unique_ptr<httplib::Client> client;
    std::string baseUrl;

public:
    RESTClient(const std::string& endpoint);
    ~RESTClient();

    // Account Management
    std::vector<Account> getAccounts();
    Account createAccount(const std::string& name);

    // Component Registry (Legacy)
    ComponentRegistrationResult registerComponent(const std::string& creator, 
                                                 const std::string& componentData, 
                                                 const std::string& context);
    ComponentRegistrationResult getComponent(const std::string& componentId);
    ComponentRegistrationResult getComponentIdentity(const std::string& componentId);
    ComponentRegistrationResult verifyComponent(const std::string& verifier, 
                                               const std::string& componentId, 
                                               const std::string& context);

    // Privacy-focused Component Operations
    AnonymousComponentResult registerAnonymousComponent(const std::string& creator,
                                                       const std::string& realComponentId,
                                                       const std::string& manufacturerId,
                                                       const std::string& componentType,
                                                       const std::string& context);
    PairingVerificationResult verifyComponentPairingWithHashes(const std::string& verifier,
                                                              const std::string& sourceHash,
                                                              const std::string& targetHash,
                                                              const std::string& context);
    PairingAuthorizationResult createAnonymousPairingAuthorization(const std::string& creator,
                                                                  const std::string& sourceHash,
                                                                  const std::string& targetHash,
                                                                  const std::string& context);
    RevocationEventResult createAnonymousRevocationEvent(const std::string& creator,
                                                        const std::string& componentHash,
                                                        const std::string& reason,
                                                        const std::string& context);
    ComponentMetadataResult getAnonymousComponentMetadata(const std::string& componentHash);

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

    // Pairing Queue Operations
    PairingRequestResult queuePairingRequest(const std::string& creator,
                                            const std::string& componentA,
                                            const std::string& componentB,
                                            const std::string& context);
    QueueStatusResult getQueueStatus(const std::string& queueId);
    std::string processOfflineQueue(const std::string& processor,
                                   const std::string& queueId,
                                   const std::string& context);
    std::string cancelRequest(const std::string& creator,
                             const std::string& requestId,
                             const std::string& reason);
    std::vector<PairingRequestResult> getQueuedRequests(const std::string& queueId);
    std::vector<PairingRequestResult> listProxyQueue(const std::string& proxyId);

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

    // WebSocket for real-time updates
    void startWebSocket(const std::string& componentId, 
                       std::function<void(const BatteryStatusUpdate&)> callback);
    void stopWebSocket();

    // Health and Status
    std::string getHealthStatus();
    std::string getBlockchainStatus();

private:
    std::string makeRequest(const std::string& method, 
                           const std::string& endpoint, 
                           const std::string& body = "");
    std::string urlEncode(const std::string& str);
}; 