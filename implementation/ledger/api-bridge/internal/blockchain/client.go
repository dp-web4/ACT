package blockchain

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

// Client represents a blockchain client
type Client struct {
	restClient *RESTClient
	logger     zerolog.Logger
}

// NewClient creates a new blockchain client
func NewClient(endpoint string, logger zerolog.Logger) (*Client, error) {
	// For now, use REST client instead of gRPC
	restClient := NewRESTClient(endpoint, logger)

	return &Client{
		restClient: restClient,
		logger:     logger,
	}, nil
}

// Close closes the blockchain connection
func (c *Client) Close() error {
	// No connection to close for REST client
	return nil
}

// RegisterComponent registers a component on the blockchain
func (c *Client) RegisterComponent(ctx context.Context, creator, componentData, context string) (map[string]interface{}, error) {
	return c.restClient.RegisterComponent(ctx, creator, componentData, context)
}

// GetComponent retrieves a component from the blockchain
func (c *Client) GetComponent(ctx context.Context, componentID string) (map[string]interface{}, error) {
	return c.restClient.GetComponent(ctx, componentID)
}

// GetComponentIdentity retrieves component identity from the blockchain
func (c *Client) GetComponentIdentity(ctx context.Context, componentID string) (map[string]interface{}, error) {
	return c.restClient.GetComponentIdentity(ctx, componentID)
}

// VerifyComponent verifies a component on the blockchain
func (c *Client) VerifyComponent(ctx context.Context, verifier, componentID, context string) (map[string]interface{}, error) {
	return c.restClient.VerifyComponent(ctx, verifier, componentID, context)
}

// Privacy-focused methods for anonymous component operations

// RegisterAnonymousComponent registers a component anonymously using hashes
func (c *Client) RegisterAnonymousComponent(ctx context.Context, creator, realComponentID, manufacturerID, componentType, context string) (map[string]interface{}, error) {
	return c.restClient.RegisterAnonymousComponent(ctx, creator, realComponentID, manufacturerID, componentType, context)
}

// VerifyComponentPairingWithHashes verifies component pairing using hashes
func (c *Client) VerifyComponentPairingWithHashes(ctx context.Context, verifier, componentHashA, componentHashB, context string) (map[string]interface{}, error) {
	return c.restClient.VerifyComponentPairingWithHashes(ctx, verifier, componentHashA, componentHashB, context)
}

// CreateAnonymousPairingAuthorization creates anonymous pairing authorization
func (c *Client) CreateAnonymousPairingAuthorization(ctx context.Context, creator, componentHashA, componentHashB, ruleHash, trustScoreRequirement, authorizationLevel string) (map[string]interface{}, error) {
	return c.restClient.CreateAnonymousPairingAuthorization(ctx, creator, componentHashA, componentHashB, ruleHash, trustScoreRequirement, authorizationLevel)
}

// CreateAnonymousRevocationEvent creates anonymous revocation event
func (c *Client) CreateAnonymousRevocationEvent(ctx context.Context, creator, targetHash, revocationType, urgencyLevel, reasonCategory, context string) (map[string]interface{}, error) {
	return c.restClient.CreateAnonymousRevocationEvent(ctx, creator, targetHash, revocationType, urgencyLevel, reasonCategory, context)
}

// GetAnonymousComponentMetadata retrieves anonymous component metadata
func (c *Client) GetAnonymousComponentMetadata(ctx context.Context, requester, componentHash string) (map[string]interface{}, error) {
	return c.restClient.GetAnonymousComponentMetadata(ctx, requester, componentHash)
}

// InitiatePairing initiates a pairing between components
func (c *Client) InitiatePairing(ctx context.Context, creator, componentA, componentB, operationalContext, proxyID string, forceImmediate bool) (map[string]interface{}, error) {
	return c.restClient.InitiatePairing(ctx, creator, componentA, componentB, operationalContext, proxyID, forceImmediate)
}

// CompletePairing completes a pairing between components
func (c *Client) CompletePairing(ctx context.Context, creator, challengeID, componentAAuth, componentBAuth, sessionContext string) (map[string]interface{}, error) {
	return c.restClient.CompletePairing(ctx, creator, challengeID, componentAAuth, componentBAuth, sessionContext)
}

// RevokePairing revokes a pairing
func (c *Client) RevokePairing(ctx context.Context, creator, lctID, reason string, notifyOffline bool) (map[string]interface{}, error) {
	return c.restClient.RevokePairing(ctx, creator, lctID, reason, notifyOffline)
}

// GetPairingStatus gets the status of a pairing
func (c *Client) GetPairingStatus(ctx context.Context, challengeID string) (map[string]interface{}, error) {
	return c.restClient.GetPairingStatus(ctx, challengeID)
}

// CreateLCT creates a Linked Context Token
func (c *Client) CreateLCT(ctx context.Context, creator, componentA, componentB, context, proxyID string) (map[string]interface{}, error) {
	return c.restClient.CreateLCT(ctx, creator, componentA, componentB, context, proxyID)
}

// GetLCT retrieves a Linked Context Token
func (c *Client) GetLCT(ctx context.Context, lctID string) (map[string]interface{}, error) {
	return c.restClient.GetLCT(ctx, lctID)
}

// UpdateLCTStatus updates the status of a Linked Context Token
func (c *Client) UpdateLCTStatus(ctx context.Context, creator, lctID, status, context string) (map[string]interface{}, error) {
	return c.restClient.UpdateLCTStatus(ctx, creator, lctID, status, context)
}

// CreateTrustTensor creates a trust tensor
func (c *Client) CreateTrustTensor(ctx context.Context, creator, componentA, componentB, context string, initialScore float64) (map[string]interface{}, error) {
	return c.restClient.CreateTrustTensor(ctx, creator, componentA, componentB, context, initialScore)
}

// GetTrustTensor retrieves a trust tensor
func (c *Client) GetTrustTensor(ctx context.Context, tensorID string) (map[string]interface{}, error) {
	return c.restClient.GetTrustTensor(ctx, tensorID)
}

// UpdateTrustScore updates the trust score
func (c *Client) UpdateTrustScore(ctx context.Context, creator, tensorID string, score float64, context string) (map[string]interface{}, error) {
	return c.restClient.UpdateTrustScore(ctx, creator, tensorID, score, context)
}

// CreateEnergyOperation creates an energy operation
func (c *Client) CreateEnergyOperation(ctx context.Context, creator, componentA, componentB, operationType string, amount float64, context string) (map[string]interface{}, error) {
	return c.restClient.CreateEnergyOperation(ctx, creator, componentA, componentB, operationType, amount, context)
}

// ExecuteEnergyTransfer executes an energy transfer
func (c *Client) ExecuteEnergyTransfer(ctx context.Context, creator, operationID string, amount float64, context string) (map[string]interface{}, error) {
	return c.restClient.ExecuteEnergyTransfer(ctx, creator, operationID, amount, context)
}

// GetEnergyBalance gets the energy balance for a component
func (c *Client) GetEnergyBalance(ctx context.Context, componentID string) (map[string]interface{}, error) {
	return c.restClient.GetEnergyBalance(ctx, componentID)
}

// Queue Management Methods

// QueuePairingRequest queues a pairing request
func (c *Client) QueuePairingRequest(ctx context.Context, componentA, componentB, operationalContext, proxyID string) (map[string]interface{}, error) {
	return c.restClient.QueuePairingRequest(ctx, componentA, componentB, operationalContext, proxyID)
}

// GetQueueStatus gets the status of a queue
func (c *Client) GetQueueStatus(ctx context.Context, componentID string) (map[string]interface{}, error) {
	return c.restClient.GetQueueStatus(ctx, componentID)
}

// ProcessOfflineQueue processes offline operations for a component
func (c *Client) ProcessOfflineQueue(ctx context.Context, componentID string) (map[string]interface{}, error) {
	return c.restClient.ProcessOfflineQueue(ctx, componentID)
}

// CancelRequest cancels a queued request
func (c *Client) CancelRequest(ctx context.Context, requestID, reason string) (map[string]interface{}, error) {
	return c.restClient.CancelRequest(ctx, requestID, reason)
}

// GetQueuedRequests gets all queued requests for a component
func (c *Client) GetQueuedRequests(ctx context.Context, componentID string) (map[string]interface{}, error) {
	return c.restClient.GetQueuedRequests(ctx, componentID)
}

// ListProxyQueue lists all operations for a proxy
func (c *Client) ListProxyQueue(ctx context.Context, proxyID string) (map[string]interface{}, error) {
	return c.restClient.ListProxyQueue(ctx, proxyID)
}

// Authorization Management Methods

// CreatePairingAuthorization creates a pairing authorization
func (c *Client) CreatePairingAuthorization(ctx context.Context, componentA, componentB, operationalContext, authorizationRules string) (map[string]interface{}, error) {
	return c.restClient.CreatePairingAuthorization(ctx, componentA, componentB, operationalContext, authorizationRules)
}

// GetComponentAuthorizations gets all authorizations for a component
func (c *Client) GetComponentAuthorizations(ctx context.Context, componentID string) (map[string]interface{}, error) {
	return c.restClient.GetComponentAuthorizations(ctx, componentID)
}

// UpdateAuthorization updates an authorization
func (c *Client) UpdateAuthorization(ctx context.Context, authorizationID string, updates map[string]interface{}) (map[string]interface{}, error) {
	return c.restClient.UpdateAuthorization(ctx, authorizationID, updates)
}

// RevokeAuthorization revokes an authorization
func (c *Client) RevokeAuthorization(ctx context.Context, authorizationID, reason string) (map[string]interface{}, error) {
	return c.restClient.RevokeAuthorization(ctx, authorizationID, reason)
}

// CheckPairingAuthorization checks if pairing is authorized
func (c *Client) CheckPairingAuthorization(ctx context.Context, componentA, componentB, operationalContext string) (map[string]interface{}, error) {
	return c.restClient.CheckPairingAuthorization(ctx, componentA, componentB, operationalContext)
}

// Trust Tensor Enhanced Methods

// CalculateRelationshipTrust calculates trust score for a relationship
func (c *Client) CalculateRelationshipTrust(ctx context.Context, componentA, componentB, operationalContext string) (map[string]interface{}, error) {
	return c.restClient.CalculateRelationshipTrust(ctx, componentA, componentB, operationalContext)
}

// GetRelationshipTensor gets a relationship tensor
func (c *Client) GetRelationshipTensor(ctx context.Context, componentA, componentB string) (map[string]interface{}, error) {
	return c.restClient.GetRelationshipTensor(ctx, componentA, componentB)
}

// UpdateTensorScore updates a tensor score
func (c *Client) UpdateTensorScore(ctx context.Context, creator, componentA, componentB string, score float64, context string) (map[string]interface{}, error) {
	return c.restClient.UpdateTensorScore(ctx, creator, componentA, componentB, score, context)
}

// GetAccountManager returns the account manager
func (c *Client) GetAccountManager() *AccountManager {
	return c.restClient.accountManager
}

// TestIgniteCLI tests Ignite CLI availability
func (c *Client) TestIgniteCLI(ctx context.Context) error {
	return c.restClient.testIgniteCLI(ctx)
}

// TestConnection tests the blockchain connection and returns status
func (c *Client) TestConnection(ctx context.Context) map[string]interface{} {
	// Test REST client connection
	status := map[string]interface{}{
		"connected": false,
		"errors":    []string{},
		"mode":      "real_blockchain",
	}

	// Test blockchain REST API connection
	if err := c.restClient.testBlockchainConnection(ctx); err != nil {
		status["errors"] = append(status["errors"].([]string), fmt.Sprintf("Blockchain REST API failed: %v", err))
	} else {
		status["rest_api_connected"] = true
	}

	// Test if Ignite CLI is accessible
	if err := c.restClient.testIgniteCLI(ctx); err != nil {
		status["errors"] = append(status["errors"].([]string), fmt.Sprintf("Ignite CLI test failed: %v", err))
	} else {
		status["ignite_available"] = true
		status["connected"] = true
	}

	// Add path information for debugging
	status["project_root"] = c.restClient.projectRoot
	status["racecar_cmd"] = c.restClient.racecarCmd
	status["ignite_path"] = c.restClient.ignitePath

	return status
}

// CreateTransactionFile creates a transaction file for testing
func (c *Client) CreateTransactionFile(message map[string]interface{}, memo string) (*os.File, error) {
	return c.restClient.createTransactionFile(message, memo)
}

func (c *Client) GetIgnitePath() string {
	return c.restClient.ignitePath
}

func (c *Client) GetProjectRoot() string {
	return c.restClient.projectRoot
}
