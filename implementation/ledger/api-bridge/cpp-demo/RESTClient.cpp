#include "RESTClient.h"
#include <httplib.h>
#include <nlohmann/json.hpp>
#include <sstream>
#include <iomanip>
#include <cctype>

using json = nlohmann::json;

RESTClient::RESTClient(const std::string& endpoint) : baseUrl(endpoint) {
    client = std::make_unique<httplib::Client>(endpoint);
    client->set_connection_timeout(10);
    client->set_read_timeout(30);
}

RESTClient::~RESTClient() = default;

std::vector<Account> RESTClient::getAccounts() {
    std::string response = makeRequest("GET", "/accounts");
    json j = json::parse(response);
    
    std::vector<Account> accounts;
    if (j.contains("accounts")) {
        for (const auto& account : j["accounts"]) {
            Account acc;
            acc.name = account.value("name", "");
            acc.address = account.value("address", "");
            acc.keyType = account.value("key_type", "");
            accounts.push_back(acc);
        }
    }
    
    return accounts;
}

Account RESTClient::createAccount(const std::string& name) {
    json request = {
        {"name", name}
    };
    
    std::string response = makeRequest("POST", "/accounts", request.dump());
    json j = json::parse(response);
    
    Account account;
    account.name = j.value("name", "");
    account.address = j.value("address", "");
    account.keyType = j.value("key_type", "");
    
    return account;
}

// Privacy-focused Component Operations
AnonymousComponentResult RESTClient::registerAnonymousComponent(const std::string& creator,
                                                               const std::string& realComponentId,
                                                               const std::string& manufacturerId,
                                                               const std::string& componentType,
                                                               const std::string& context) {
    json request = {
        {"creator", creator},
        {"real_component_id", realComponentId},
        {"manufacturer_id", manufacturerId},
        {"component_type", componentType},
        {"context", context}
    };
    
    std::string response = makeRequest("POST", "/components/register-anonymous", request.dump());
    json j = json::parse(response);
    
    AnonymousComponentResult result;
    result.componentHash = j.value("component_hash", "");
    result.manufacturerHash = j.value("manufacturer_hash", "");
    result.categoryHash = j.value("category_hash", "");
    result.creator = j.value("creator", "");
    result.context = j.value("context", "");
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

PairingVerificationResult RESTClient::verifyComponentPairingWithHashes(const std::string& verifier,
                                                                      const std::string& sourceHash,
                                                                      const std::string& targetHash,
                                                                      const std::string& context) {
    json request = {
        {"verifier", verifier},
        {"source_hash", sourceHash},
        {"target_hash", targetHash},
        {"context", context}
    };
    
    std::string response = makeRequest("POST", "/components/verify-pairing-hashes", request.dump());
    json j = json::parse(response);
    
    PairingVerificationResult result;
    result.sourceHash = j.value("source_hash", "");
    result.targetHash = j.value("target_hash", "");
    result.context = j.value("context", "");
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

PairingAuthorizationResult RESTClient::createAnonymousPairingAuthorization(const std::string& creator,
                                                                          const std::string& sourceHash,
                                                                          const std::string& targetHash,
                                                                          const std::string& context) {
    json request = {
        {"creator", creator},
        {"source_hash", sourceHash},
        {"target_hash", targetHash},
        {"context", context}
    };
    
    std::string response = makeRequest("POST", "/components/create-pairing-authorization", request.dump());
    json j = json::parse(response);
    
    PairingAuthorizationResult result;
    result.authorizationId = j.value("authorization_id", "");
    result.sourceHash = j.value("source_hash", "");
    result.targetHash = j.value("target_hash", "");
    result.context = j.value("context", "");
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

RevocationEventResult RESTClient::createAnonymousRevocationEvent(const std::string& creator,
                                                                const std::string& componentHash,
                                                                const std::string& reason,
                                                                const std::string& context) {
    json request = {
        {"creator", creator},
        {"component_hash", componentHash},
        {"reason", reason},
        {"context", context}
    };
    
    std::string response = makeRequest("POST", "/components/create-revocation-event", request.dump());
    json j = json::parse(response);
    
    RevocationEventResult result;
    result.revocationId = j.value("revocation_id", "");
    result.componentHash = j.value("component_hash", "");
    result.reason = j.value("reason", "");
    result.context = j.value("context", "");
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

ComponentMetadataResult RESTClient::getAnonymousComponentMetadata(const std::string& componentHash) {
    std::string response = makeRequest("GET", "/components/anonymous/" + urlEncode(componentHash) + "/metadata");
    json j = json::parse(response);
    
    ComponentMetadataResult result;
    result.componentHash = j.value("component_hash", "");
    result.metadata = j.value("metadata", "");
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

// Pairing Queue Operations
PairingRequestResult RESTClient::queuePairingRequest(const std::string& creator,
                                                    const std::string& componentA,
                                                    const std::string& componentB,
                                                    const std::string& context) {
    json request = {
        {"creator", creator},
        {"component_a", componentA},
        {"component_b", componentB},
        {"context", context}
    };
    
    std::string response = makeRequest("POST", "/pairing/queue", request.dump());
    json j = json::parse(response);
    
    PairingRequestResult result;
    result.requestId = j.value("request_id", "");
    result.componentA = j.value("component_a", "");
    result.componentB = j.value("component_b", "");
    result.context = j.value("context", "");
    result.status = j.value("status", "");
    result.createdAt = j.value("created_at", 0);
    result.creator = j.value("creator", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

QueueStatusResult RESTClient::getQueueStatus(const std::string& queueId) {
    std::string response = makeRequest("GET", "/pairing/queue/" + urlEncode(queueId) + "/status");
    json j = json::parse(response);
    
    QueueStatusResult result;
    result.queueId = j.value("queue_id", "");
    result.pendingRequests = j.value("pending_requests", 0);
    result.processedRequests = j.value("processed_requests", 0);
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

std::string RESTClient::processOfflineQueue(const std::string& processor,
                                           const std::string& queueId,
                                           const std::string& context) {
    json request = {
        {"processor", processor},
        {"queue_id", queueId},
        {"context", context}
    };
    
    std::string response = makeRequest("POST", "/pairing/queue/process", request.dump());
    json j = json::parse(response);
    
    return j.value("result", "");
}

std::string RESTClient::cancelRequest(const std::string& creator,
                                     const std::string& requestId,
                                     const std::string& reason) {
    json request = {
        {"creator", creator},
        {"request_id", requestId},
        {"reason", reason}
    };
    
    std::string response = makeRequest("POST", "/pairing/queue/cancel", request.dump());
    json j = json::parse(response);
    
    return j.value("result", "");
}

std::vector<PairingRequestResult> RESTClient::getQueuedRequests(const std::string& queueId) {
    std::string response = makeRequest("GET", "/pairing/queue/" + urlEncode(queueId) + "/requests");
    json j = json::parse(response);
    
    std::vector<PairingRequestResult> requests;
    if (j.contains("requests")) {
        for (const auto& req : j["requests"]) {
            PairingRequestResult result;
            result.requestId = req.value("request_id", "");
            result.componentA = req.value("component_a", "");
            result.componentB = req.value("component_b", "");
            result.context = req.value("context", "");
            result.status = req.value("status", "");
            result.createdAt = req.value("created_at", 0);
            result.creator = req.value("creator", "");
            result.txHash = req.value("tx_hash", "");
            requests.push_back(result);
        }
    }
    
    return requests;
}

std::vector<PairingRequestResult> RESTClient::listProxyQueue(const std::string& proxyId) {
    std::string response = makeRequest("GET", "/pairing/queue/proxy/" + urlEncode(proxyId));
    json j = json::parse(response);
    
    std::vector<PairingRequestResult> requests;
    if (j.contains("requests")) {
        for (const auto& req : j["requests"]) {
            PairingRequestResult result;
            result.requestId = req.value("request_id", "");
            result.componentA = req.value("component_a", "");
            result.componentB = req.value("component_b", "");
            result.context = req.value("context", "");
            result.status = req.value("status", "");
            result.createdAt = req.value("created_at", 0);
            result.creator = req.value("creator", "");
            result.txHash = req.value("tx_hash", "");
            requests.push_back(result);
        }
    }
    
    return requests;
}

// Health and Status
std::string RESTClient::getHealthStatus() {
    return makeRequest("GET", "/health");
}

std::string RESTClient::getBlockchainStatus() {
    return makeRequest("GET", "/blockchain/status");
}

ComponentRegistrationResult RESTClient::registerComponent(const std::string& creator, 
                                                        const std::string& componentData, 
                                                        const std::string& context) {
    json request = {
        {"creator", creator},
        {"component_data", componentData},
        {"context", context}
    };
    
    std::string response = makeRequest("POST", "/components/register", request.dump());
    json j = json::parse(response);
    
    ComponentRegistrationResult result;
    result.componentId = j.value("component_id", "");
    result.componentIdentity = j.value("component_identity", "");
    result.componentData = j.value("component_data", "");
    result.context = j.value("context", "");
    result.creator = j.value("creator", "");
    result.lctId = j.value("lct_id", "");
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

ComponentRegistrationResult RESTClient::getComponent(const std::string& componentId) {
    std::string response = makeRequest("GET", "/components/" + urlEncode(componentId));
    json j = json::parse(response);
    
    ComponentRegistrationResult result;
    result.componentId = j.value("component_id", "");
    result.componentIdentity = j.value("component_identity", "");
    result.componentData = j.value("component_data", "");
    result.context = j.value("context", "");
    result.creator = j.value("creator", "");
    result.lctId = j.value("lct_id", "");
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

ComponentRegistrationResult RESTClient::getComponentIdentity(const std::string& componentId) {
    std::string response = makeRequest("GET", "/components/" + urlEncode(componentId) + "/identity");
    json j = json::parse(response);
    
    ComponentRegistrationResult result;
    result.componentId = j.value("component_id", "");
    result.componentIdentity = j.value("component_identity", "");
    result.componentData = j.value("component_data", "");
    result.context = j.value("context", "");
    result.creator = j.value("creator", "");
    result.lctId = j.value("lct_id", "");
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

ComponentRegistrationResult RESTClient::verifyComponent(const std::string& verifier, 
                                                       const std::string& componentId, 
                                                       const std::string& context) {
    json request = {
        {"verifier", verifier},
        {"component_id", componentId},
        {"context", context}
    };
    
    std::string response = makeRequest("POST", "/components/verify", request.dump());
    json j = json::parse(response);
    
    ComponentRegistrationResult result;
    result.componentId = j.value("component_id", "");
    result.componentIdentity = j.value("component_identity", "");
    result.componentData = j.value("component_data", "");
    result.context = j.value("context", "");
    result.creator = j.value("creator", "");
    result.lctId = j.value("lct_id", "");
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

LCTResult RESTClient::createLCT(const std::string& creator, 
                               const std::string& componentA, 
                               const std::string& componentB, 
                               const std::string& context, 
                               const std::string& proxyId) {
    json request = {
        {"creator", creator},
        {"component_a", componentA},
        {"component_b", componentB},
        {"context", context},
        {"proxy_id", proxyId}
    };
    
    std::string response = makeRequest("POST", "/lct/create", request.dump());
    json j = json::parse(response);
    
    LCTResult result;
    result.lctId = j.value("lct_id", "");
    result.componentA = j.value("component_a", "");
    result.componentB = j.value("component_b", "");
    result.context = j.value("context", "");
    result.proxyId = j.value("proxy_id", "");
    result.status = j.value("status", "");
    result.createdAt = j.value("created_at", 0);
    result.creator = j.value("creator", "");
    result.txHash = j.value("tx_hash", "");
    result.lctKeyHalf = j.value("lct_key_half", "");
    result.deviceKeyHalf = j.value("device_key_half", "");
    
    return result;
}

LCTResult RESTClient::getLCT(const std::string& lctId) {
    std::string response = makeRequest("GET", "/lct/" + urlEncode(lctId));
    json j = json::parse(response);
    
    LCTResult result;
    result.lctId = j.value("lct_id", "");
    result.componentA = j.value("component_a", "");
    result.componentB = j.value("component_b", "");
    result.context = j.value("context", "");
    result.proxyId = j.value("proxy_id", "");
    result.status = j.value("status", "");
    result.createdAt = j.value("created_at", 0);
    result.creator = j.value("creator", "");
    result.txHash = j.value("tx_hash", "");
    result.lctKeyHalf = j.value("lct_key_half", "");
    result.deviceKeyHalf = j.value("device_key_half", "");
    
    return result;
}

LCTResult RESTClient::updateLCTStatus(const std::string& creator, 
                                     const std::string& lctId, 
                                     const std::string& status, 
                                     const std::string& context) {
    json request = {
        {"creator", creator},
        {"lct_id", lctId},
        {"status", status},
        {"context", context}
    };
    
    std::string response = makeRequest("PUT", "/lct/" + urlEncode(lctId) + "/status", request.dump());
    json j = json::parse(response);
    
    LCTResult result;
    result.lctId = j.value("lct_id", "");
    result.componentA = j.value("component_a", "");
    result.componentB = j.value("component_b", "");
    result.context = j.value("context", "");
    result.proxyId = j.value("proxy_id", "");
    result.status = j.value("status", "");
    result.createdAt = j.value("created_at", 0);
    result.creator = j.value("creator", "");
    result.txHash = j.value("tx_hash", "");
    result.lctKeyHalf = j.value("lct_key_half", "");
    result.deviceKeyHalf = j.value("device_key_half", "");
    
    return result;
}

PairingInitiateResult RESTClient::initiatePairing(const std::string& creator, 
                                                 const std::string& componentA, 
                                                 const std::string& componentB, 
                                                 const std::string& operationalContext, 
                                                 const std::string& proxyId, 
                                                 bool forceImmediate) {
    json request = {
        {"creator", creator},
        {"component_a", componentA},
        {"component_b", componentB},
        {"operational_context", operationalContext},
        {"proxy_id", proxyId},
        {"force_immediate", forceImmediate}
    };
    
    std::string response = makeRequest("POST", "/pairing/initiate", request.dump());
    json j = json::parse(response);
    
    PairingInitiateResult result;
    result.challengeId = j.value("challenge_id", "");
    result.componentA = j.value("component_a", "");
    result.componentB = j.value("component_b", "");
    result.operationalContext = j.value("operational_context", "");
    result.proxyId = j.value("proxy_id", "");
    result.forceImmediate = j.value("force_immediate", false);
    result.status = j.value("status", "");
    result.createdAt = j.value("created_at", 0);
    result.creator = j.value("creator", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

PairingCompleteResult RESTClient::completePairing(const std::string& creator, 
                                                 const std::string& challengeId, 
                                                 const std::string& componentAAuth, 
                                                 const std::string& componentBAuth, 
                                                 const std::string& sessionContext) {
    json request = {
        {"creator", creator},
        {"challenge_id", challengeId},
        {"component_a_auth", componentAAuth},
        {"component_b_auth", componentBAuth},
        {"session_context", sessionContext}
    };
    
    std::string response = makeRequest("POST", "/pairing/complete", request.dump());
    json j = json::parse(response);
    
    PairingCompleteResult result;
    result.lctId = j.value("lct_id", "");
    result.sessionKeys = j.value("session_keys", "");
    result.trustSummary = j.value("trust_summary", "");
    result.txHash = j.value("tx_hash", "");
    result.splitKeyA = j.value("split_key_a", "");
    result.splitKeyB = j.value("split_key_b", "");
    
    return result;
}

std::string RESTClient::revokePairing(const std::string& creator, 
                                     const std::string& lctId, 
                                     const std::string& reason, 
                                     bool notifyOffline) {
    json request = {
        {"creator", creator},
        {"lct_id", lctId},
        {"reason", reason},
        {"notify_offline", notifyOffline}
    };
    
    std::string response = makeRequest("POST", "/pairing/revoke", request.dump());
    json j = json::parse(response);
    
    return j.value("status", "");
}

std::string RESTClient::getPairingStatus(const std::string& challengeId) {
    std::string response = makeRequest("GET", "/pairing/status/" + urlEncode(challengeId));
    json j = json::parse(response);
    
    return j.value("status", "");
}

TrustTensorResult RESTClient::createTrustTensor(const std::string& creator, 
                                               const std::string& componentA, 
                                               const std::string& componentB, 
                                               const std::string& context, 
                                               double initialScore) {
    json request = {
        {"creator", creator},
        {"component_a", componentA},
        {"component_b", componentB},
        {"context", context},
        {"initial_score", initialScore}
    };
    
    std::string response = makeRequest("POST", "/trust/create", request.dump());
    json j = json::parse(response);
    
    TrustTensorResult result;
    result.tensorId = j.value("tensor_id", "");
    result.score = j.value("score", 0.0);
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

TrustTensorResult RESTClient::getTrustTensor(const std::string& tensorId) {
    std::string response = makeRequest("GET", "/trust/" + urlEncode(tensorId));
    json j = json::parse(response);
    
    TrustTensorResult result;
    result.tensorId = j.value("tensor_id", "");
    result.score = j.value("score", 0.0);
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

TrustTensorResult RESTClient::updateTrustScore(const std::string& creator, 
                                              const std::string& tensorId, 
                                              double score, 
                                              const std::string& context) {
    json request = {
        {"creator", creator},
        {"tensor_id", tensorId},
        {"score", score},
        {"context", context}
    };
    
    std::string response = makeRequest("PUT", "/trust/" + urlEncode(tensorId) + "/score", request.dump());
    json j = json::parse(response);
    
    TrustTensorResult result;
    result.tensorId = j.value("tensor_id", "");
    result.score = j.value("score", 0.0);
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

EnergyOperationResult RESTClient::createEnergyOperation(const std::string& creator, 
                                                       const std::string& componentA, 
                                                       const std::string& componentB, 
                                                       const std::string& operationType, 
                                                       double amount, 
                                                       const std::string& context) {
    json request = {
        {"creator", creator},
        {"component_a", componentA},
        {"component_b", componentB},
        {"operation_type", operationType},
        {"amount", amount},
        {"context", context}
    };
    
    std::string response = makeRequest("POST", "/energy/create", request.dump());
    json j = json::parse(response);
    
    EnergyOperationResult result;
    result.operationId = j.value("operation_id", "");
    result.operationType = j.value("operation_type", "");
    result.amount = j.value("amount", 0.0);
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

EnergyOperationResult RESTClient::executeEnergyTransfer(const std::string& creator, 
                                                       const std::string& operationId, 
                                                       double amount, 
                                                       const std::string& context) {
    json request = {
        {"creator", creator},
        {"operation_id", operationId},
        {"amount", amount},
        {"context", context}
    };
    
    std::string response = makeRequest("POST", "/energy/transfer", request.dump());
    json j = json::parse(response);
    
    EnergyOperationResult result;
    result.operationId = j.value("operation_id", "");
    result.operationType = j.value("operation_type", "");
    result.amount = j.value("amount", 0.0);
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

double RESTClient::getEnergyBalance(const std::string& componentId) {
    std::string response = makeRequest("GET", "/energy/balance/" + urlEncode(componentId));
    json j = json::parse(response);
    
    return j.value("balance", 0.0);
}

void RESTClient::startWebSocket(const std::string& componentId, 
                               std::function<void(const BatteryStatusUpdate&)> callback) {
    // WebSocket implementation would go here
    // For demo purposes, we'll simulate with polling
    std::cout << "WebSocket connection started for component: " << componentId << std::endl;
}

void RESTClient::stopWebSocket() {
    std::cout << "WebSocket connection stopped" << std::endl;
}

std::string RESTClient::makeRequest(const std::string& method, 
                                   const std::string& endpoint, 
                                   const std::string& body) {
    httplib::Headers headers = {
        {"Content-Type", "application/json"},
        {"Accept", "application/json"}
    };
    
    httplib::Result result;
    
    if (method == "GET") {
        result = client->Get(endpoint, headers);
    } else if (method == "POST") {
        result = client->Post(endpoint, headers, body, "application/json");
    } else if (method == "PUT") {
        result = client->Put(endpoint, headers, body, "application/json");
    } else if (method == "DELETE") {
        result = client->Delete(endpoint, headers);
    } else {
        throw std::runtime_error("Unsupported HTTP method: " + method);
    }
    
    if (!result) {
        throw std::runtime_error("HTTP request failed: " + std::to_string(result.error()));
    }
    
    if (result->status != 200) {
        throw std::runtime_error("HTTP error " + std::to_string(result->status) + ": " + result->body);
    }
    
    return result->body;
}

std::string RESTClient::urlEncode(const std::string& str) {
    std::ostringstream escaped;
    escaped.fill('0');
    escaped << std::hex;
    
    for (char c : str) {
        if (isalnum(c) || c == '-' || c == '_' || c == '.' || c == '~') {
            escaped << c;
        } else {
            escaped << '%' << std::setw(2) << int((unsigned char)c);
        }
    }
    
    return escaped.str();
} 