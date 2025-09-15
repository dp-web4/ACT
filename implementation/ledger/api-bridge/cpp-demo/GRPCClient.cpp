#include "GRPCClient.h"
#include <grpcpp/grpcpp.h>
#include <grpcpp/create_channel.h>
#include <grpcpp/security/credentials.h>
#include <nlohmann/json.hpp>
#include <chrono>
#include <thread>
#include <random>

using json = nlohmann::json;
using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;

// Forward declarations for generated protobuf services
namespace apibridge {
    class APIBridgeService;
    class APIBridgeServiceStub;
}

GRPCClient::GRPCClient(const std::string& endpoint) : serverAddress(endpoint), streamingActive(false) {
    channel = grpc::CreateChannel(endpoint, grpc::InsecureChannelCredentials());
    
    if (!channel) {
        throw std::runtime_error("Failed to create gRPC channel to " + endpoint);
    }
}

GRPCClient::~GRPCClient() {
    stopStreaming();
}

std::vector<Account> GRPCClient::getAccounts() {
    std::string response = makeRequest("APIBridgeService", "GetAccounts", "{}");
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

Account GRPCClient::createAccount(const std::string& name) {
    json request = {
        {"name", name}
    };
    
    std::string response = makeRequest("APIBridgeService", "CreateAccount", request.dump());
    json j = json::parse(response);
    
    Account account;
    account.name = j.value("name", "");
    account.address = j.value("address", "");
    account.keyType = j.value("key_type", "");
    
    return account;
}

ComponentRegistrationResult GRPCClient::registerComponent(const std::string& creator, 
                                                         const std::string& componentData, 
                                                         const std::string& context) {
    json request = {
        {"creator", creator},
        {"component_data", componentData},
        {"context", context}
    };
    
    std::string response = makeRequest("APIBridgeService", "RegisterComponent", request.dump());
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

ComponentRegistrationResult GRPCClient::getComponent(const std::string& componentId) {
    json request = {
        {"component_id", componentId}
    };
    
    std::string response = makeRequest("APIBridgeService", "GetComponent", request.dump());
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

ComponentRegistrationResult GRPCClient::getComponentIdentity(const std::string& componentId) {
    json request = {
        {"component_id", componentId}
    };
    
    std::string response = makeRequest("APIBridgeService", "GetComponentIdentity", request.dump());
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

ComponentRegistrationResult GRPCClient::verifyComponent(const std::string& verifier, 
                                                       const std::string& componentId, 
                                                       const std::string& context) {
    json request = {
        {"verifier", verifier},
        {"component_id", componentId},
        {"context", context}
    };
    
    std::string response = makeRequest("APIBridgeService", "VerifyComponent", request.dump());
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

LCTResult GRPCClient::createLCT(const std::string& creator, 
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
    
    std::string response = makeRequest("APIBridgeService", "CreateLCT", request.dump());
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

LCTResult GRPCClient::getLCT(const std::string& lctId) {
    json request = {
        {"lct_id", lctId}
    };
    
    std::string response = makeRequest("APIBridgeService", "GetLCT", request.dump());
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

LCTResult GRPCClient::updateLCTStatus(const std::string& creator, 
                                     const std::string& lctId, 
                                     const std::string& status, 
                                     const std::string& context) {
    json request = {
        {"creator", creator},
        {"lct_id", lctId},
        {"status", status},
        {"context", context}
    };
    
    std::string response = makeRequest("APIBridgeService", "UpdateLCTStatus", request.dump());
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

PairingInitiateResult GRPCClient::initiatePairing(const std::string& creator, 
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
    
    std::string response = makeRequest("APIBridgeService", "InitiatePairing", request.dump());
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

PairingCompleteResult GRPCClient::completePairing(const std::string& creator, 
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
    
    std::string response = makeRequest("APIBridgeService", "CompletePairing", request.dump());
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

std::string GRPCClient::revokePairing(const std::string& creator, 
                                     const std::string& lctId, 
                                     const std::string& reason, 
                                     bool notifyOffline) {
    json request = {
        {"creator", creator},
        {"lct_id", lctId},
        {"reason", reason},
        {"notify_offline", notifyOffline}
    };
    
    std::string response = makeRequest("APIBridgeService", "RevokePairing", request.dump());
    json j = json::parse(response);
    
    return j.value("status", "");
}

std::string GRPCClient::getPairingStatus(const std::string& challengeId) {
    json request = {
        {"challenge_id", challengeId}
    };
    
    std::string response = makeRequest("APIBridgeService", "GetPairingStatus", request.dump());
    json j = json::parse(response);
    
    return j.value("status", "");
}

TrustTensorResult GRPCClient::createTrustTensor(const std::string& creator, 
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
    
    std::string response = makeRequest("APIBridgeService", "CreateTrustTensor", request.dump());
    json j = json::parse(response);
    
    TrustTensorResult result;
    result.tensorId = j.value("tensor_id", "");
    result.score = j.value("score", 0.0);
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

TrustTensorResult GRPCClient::getTrustTensor(const std::string& tensorId) {
    json request = {
        {"tensor_id", tensorId}
    };
    
    std::string response = makeRequest("APIBridgeService", "GetTrustTensor", request.dump());
    json j = json::parse(response);
    
    TrustTensorResult result;
    result.tensorId = j.value("tensor_id", "");
    result.score = j.value("score", 0.0);
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

TrustTensorResult GRPCClient::updateTrustScore(const std::string& creator, 
                                              const std::string& tensorId, 
                                              double score, 
                                              const std::string& context) {
    json request = {
        {"creator", creator},
        {"tensor_id", tensorId},
        {"score", score},
        {"context", context}
    };
    
    std::string response = makeRequest("APIBridgeService", "UpdateTrustScore", request.dump());
    json j = json::parse(response);
    
    TrustTensorResult result;
    result.tensorId = j.value("tensor_id", "");
    result.score = j.value("score", 0.0);
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

EnergyOperationResult GRPCClient::createEnergyOperation(const std::string& creator, 
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
    
    std::string response = makeRequest("APIBridgeService", "CreateEnergyOperation", request.dump());
    json j = json::parse(response);
    
    EnergyOperationResult result;
    result.operationId = j.value("operation_id", "");
    result.operationType = j.value("operation_type", "");
    result.amount = j.value("amount", 0.0);
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

EnergyOperationResult GRPCClient::executeEnergyTransfer(const std::string& creator, 
                                                       const std::string& operationId, 
                                                       double amount, 
                                                       const std::string& context) {
    json request = {
        {"creator", creator},
        {"operation_id", operationId},
        {"amount", amount},
        {"context", context}
    };
    
    std::string response = makeRequest("APIBridgeService", "ExecuteEnergyTransfer", request.dump());
    json j = json::parse(response);
    
    EnergyOperationResult result;
    result.operationId = j.value("operation_id", "");
    result.operationType = j.value("operation_type", "");
    result.amount = j.value("amount", 0.0);
    result.status = j.value("status", "");
    result.txHash = j.value("tx_hash", "");
    
    return result;
}

double GRPCClient::getEnergyBalance(const std::string& componentId) {
    json request = {
        {"component_id", componentId}
    };
    
    std::string response = makeRequest("APIBridgeService", "GetEnergyBalance", request.dump());
    json j = json::parse(response);
    
    return j.value("balance", 0.0);
}

void GRPCClient::streamBatteryStatus(const std::string& componentId, 
                                    int updateIntervalSeconds,
                                    std::function<void(const BatteryStatusUpdate&)> callback) {
    if (streamingActive) {
        stopStreaming();
    }
    
    streamingActive = true;
    streamingThread = std::make_unique<std::thread>(
        &GRPCClient::streamingWorker, this, componentId, updateIntervalSeconds, callback
    );
}

void GRPCClient::stopStreaming() {
    streamingActive = false;
    if (streamingThread && streamingThread->joinable()) {
        streamingThread->join();
    }
}

bool GRPCClient::isConnected() const {
    return channel && channel->GetState(false) == GRPC_CHANNEL_READY;
}

std::string GRPCClient::makeRequest(const std::string& service, 
                                   const std::string& method, 
                                   const std::string& request) {
    // For demo purposes, we'll simulate gRPC calls by making HTTP requests to the gRPC gateway
    // In a real implementation, you would use the generated protobuf stubs
    
    std::string url = "http://" + serverAddress + "/" + service + "/" + method;
    
    // Create a simple HTTP client for the demo
    // In production, use the actual gRPC stubs
    httplib::Client httpClient(serverAddress);
    httpClient.set_connection_timeout(10);
    httpClient.set_read_timeout(30);
    
    httplib::Headers headers = {
        {"Content-Type", "application/json"},
        {"Accept", "application/json"}
    };
    
    auto result = httpClient.Post("/" + service + "/" + method, headers, request, "application/json");
    
    if (!result) {
        throw std::runtime_error("gRPC request failed: " + std::to_string(result.error()));
    }
    
    if (result->status != 200) {
        throw std::runtime_error("gRPC error " + std::to_string(result->status) + ": " + result->body);
    }
    
    return result->body;
}

void GRPCClient::streamingWorker(const std::string& componentId, 
                                int updateIntervalSeconds,
                                std::function<void(const BatteryStatusUpdate&)> callback) {
    std::random_device rd;
    std::mt19937 gen(rd());
    std::uniform_real_distribution<> voltageDist(3.0, 4.2);
    std::uniform_real_distribution<> currentDist(-50.0, 50.0);
    std::uniform_real_distribution<> tempDist(15.0, 45.0);
    std::uniform_real_distribution<> socDist(0.0, 100.0);
    
    std::vector<std::string> statuses = {"normal", "charging", "discharging", "warning", "error"};
    std::uniform_int_distribution<> statusDist(0, statuses.size() - 1);
    
    while (streamingActive) {
        BatteryStatusUpdate update;
        update.componentId = componentId;
        update.voltage = voltageDist(gen);
        update.current = currentDist(gen);
        update.temperature = tempDist(gen);
        update.stateOfCharge = socDist(gen);
        update.status = statuses[statusDist(gen)];
        update.timestamp = std::chrono::duration_cast<std::chrono::milliseconds>(
            std::chrono::system_clock::now().time_since_epoch()
        ).count();
        
        callback(update);
        
        std::this_thread::sleep_for(std::chrono::seconds(updateIntervalSeconds));
    }
} 