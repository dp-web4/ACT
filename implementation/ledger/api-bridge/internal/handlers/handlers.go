package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"api-bridge/internal/blockchain"
	"api-bridge/internal/config"
	"api-bridge/internal/events"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

// Handler handles HTTP requests
type Handler struct {
	config     *config.Config
	logger     zerolog.Logger
	blockchain *blockchain.Client
	upgrader   websocket.Upgrader
	eventQueue *events.EventQueue
}

// New creates a new handler instance
func New(cfg *config.Config, logger zerolog.Logger) (*Handler, error) {
	// Create blockchain client using REST endpoint
	bcClient, err := blockchain.NewClient(cfg.Blockchain.RESTEndpoint, logger)
	if err != nil {
		return nil, err
	}

	// Create WebSocket upgrader
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now
		},
	}

	// Create event queue if enabled
	var eventQueue *events.EventQueue
	if cfg.Events.Enabled {
		eventQueue = events.NewEventQueue(cfg.Events.Endpoints, cfg.Events.MaxRetries, time.Duration(cfg.Events.RetryDelay)*time.Second, logger)
	}

	return &Handler{
		config:     cfg,
		logger:     logger,
		blockchain: bcClient,
		upgrader:   upgrader,
		eventQueue: eventQueue,
	}, nil
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "api-bridge",
	})
}

// BlockchainStatus handles blockchain connection status requests
func (h *Handler) BlockchainStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Test blockchain connection
	status := h.blockchain.TestConnection(ctx)

	c.JSON(http.StatusOK, gin.H{
		"blockchain_status": status,
		"timestamp":         time.Now().Unix(),
		"service":           "api-bridge",
	})
}

// RegisterComponent handles component registration
func (h *Handler) RegisterComponent(c *gin.Context) {
	var req struct {
		Creator       string `json:"creator" binding:"required"`
		ComponentData string `json:"component_data" binding:"required"`
		Context       string `json:"context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.RegisterComponent(ctx, req.Creator, req.ComponentData, req.Context)
	if err != nil {
		h.logger.Error().Err(err).Str("creator", req.Creator).Msg("Failed to register component")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to register component: %v", err)})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"component_id":   resp["component_id"],
			"creator":        req.Creator,
			"component_data": req.ComponentData,
			"context":        req.Context,
			"timestamp":      time.Now().Unix(),
			"tx_hash":        resp["txhash"],
		}
		h.eventQueue.Emit("component_registered", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// GetComponent handles component retrieval
func (h *Handler) GetComponent(c *gin.Context) {
	componentID := c.Param("id")
	if componentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Component ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	component, err := h.blockchain.GetComponent(ctx, componentID)
	if err != nil {
		h.logger.Error().Err(err).Str("component_id", componentID).Msg("Failed to get component")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get component"})
		return
	}

	c.JSON(http.StatusOK, component)
}

// GetComponentIdentity handles component identity retrieval
func (h *Handler) GetComponentIdentity(c *gin.Context) {
	componentID := c.Param("id")
	if componentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Component ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	identity, err := h.blockchain.GetComponentIdentity(ctx, componentID)
	if err != nil {
		h.logger.Error().Err(err).Str("component_id", componentID).Msg("Failed to get component identity")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get component identity"})
		return
	}

	c.JSON(http.StatusOK, identity)
}

// Privacy-focused handlers for anonymous component operations

// RegisterAnonymousComponent handles anonymous component registration
func (h *Handler) RegisterAnonymousComponent(c *gin.Context) {
	var req struct {
		Creator         string `json:"creator" binding:"required"`
		RealComponentID string `json:"real_component_id" binding:"required"`
		ManufacturerID  string `json:"manufacturer_id" binding:"required"`
		ComponentType   string `json:"component_type" binding:"required"`
		Context         string `json:"context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	// Use the new anonymous registration endpoint
	resp, err := h.blockchain.RegisterAnonymousComponent(ctx, req.Creator, req.RealComponentID, req.ManufacturerID, req.ComponentType, req.Context)
	if err != nil {
		h.logger.Error().Err(err).Str("creator", req.Creator).Msg("Failed to register anonymous component")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to register anonymous component: %v", err)})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"component_hash":    resp["component_hash"],
			"manufacturer_hash": resp["manufacturer_hash"],
			"category_hash":     resp["category_hash"],
			"creator":           req.Creator,
			"context":           req.Context,
			"timestamp":         time.Now().Unix(),
			"tx_hash":           resp["txhash"],
		}
		h.eventQueue.Emit("anonymous_component_registered", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// VerifyComponentPairingWithHashes handles hash-based component pairing verification
func (h *Handler) VerifyComponentPairingWithHashes(c *gin.Context) {
	var req struct {
		Verifier       string `json:"verifier" binding:"required"`
		ComponentHashA string `json:"component_hash_a" binding:"required"`
		ComponentHashB string `json:"component_hash_b" binding:"required"`
		Context        string `json:"context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.VerifyComponentPairingWithHashes(ctx, req.Verifier, req.ComponentHashA, req.ComponentHashB, req.Context)
	if err != nil {
		h.logger.Error().Err(err).Str("component_hash_a", req.ComponentHashA).Str("component_hash_b", req.ComponentHashB).Msg("Failed to verify component pairing with hashes")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify component pairing with hashes"})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"component_hash_a": req.ComponentHashA,
			"component_hash_b": req.ComponentHashB,
			"verifier":         req.Verifier,
			"can_pair":         resp["can_pair"],
			"reason":           resp["reason"],
			"trust_score":      resp["trust_score"],
			"context":          req.Context,
			"timestamp":        time.Now().Unix(),
			"tx_hash":          resp["txhash"],
		}
		h.eventQueue.Emit("component_pairing_verified_with_hashes", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// CreateAnonymousPairingAuthorization handles anonymous pairing authorization creation
func (h *Handler) CreateAnonymousPairingAuthorization(c *gin.Context) {
	var req struct {
		Creator               string `json:"creator" binding:"required"`
		ComponentHashA        string `json:"component_hash_a" binding:"required"`
		ComponentHashB        string `json:"component_hash_b" binding:"required"`
		RuleHash              string `json:"rule_hash" binding:"required"`
		TrustScoreRequirement string `json:"trust_score_requirement"`
		AuthorizationLevel    string `json:"authorization_level"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.CreateAnonymousPairingAuthorization(ctx, req.Creator, req.ComponentHashA, req.ComponentHashB, req.RuleHash, req.TrustScoreRequirement, req.AuthorizationLevel)
	if err != nil {
		h.logger.Error().Err(err).Str("creator", req.Creator).Msg("Failed to create anonymous pairing authorization")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create anonymous pairing authorization: %v", err)})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"auth_id":          resp["auth_id"],
			"component_hash_a": req.ComponentHashA,
			"component_hash_b": req.ComponentHashB,
			"creator":          req.Creator,
			"status":           resp["status"],
			"expires_at":       resp["expires_at"],
			"timestamp":        time.Now().Unix(),
			"tx_hash":          resp["txhash"],
		}
		h.eventQueue.Emit("anonymous_pairing_authorized", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// CreateAnonymousRevocationEvent handles anonymous revocation event creation
func (h *Handler) CreateAnonymousRevocationEvent(c *gin.Context) {
	var req struct {
		Creator        string `json:"creator" binding:"required"`
		TargetHash     string `json:"target_hash" binding:"required"`
		RevocationType string `json:"revocation_type" binding:"required"`
		UrgencyLevel   string `json:"urgency_level" binding:"required"`
		ReasonCategory string `json:"reason_category" binding:"required"`
		Context        string `json:"context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.CreateAnonymousRevocationEvent(ctx, req.Creator, req.TargetHash, req.RevocationType, req.UrgencyLevel, req.ReasonCategory, req.Context)
	if err != nil {
		h.logger.Error().Err(err).Str("creator", req.Creator).Str("target_hash", req.TargetHash).Msg("Failed to create anonymous revocation event")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create anonymous revocation event: %v", err)})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"revocation_id":   resp["revocation_id"],
			"target_hash":     req.TargetHash,
			"revocation_type": req.RevocationType,
			"urgency_level":   req.UrgencyLevel,
			"reason_category": req.ReasonCategory,
			"creator":         req.Creator,
			"context":         req.Context,
			"status":          resp["status"],
			"effective_at":    resp["effective_at"],
			"timestamp":       time.Now().Unix(),
			"tx_hash":         resp["txhash"],
		}
		h.eventQueue.Emit("anonymous_revocation_created", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// GetAnonymousComponentMetadata handles anonymous component metadata retrieval
func (h *Handler) GetAnonymousComponentMetadata(c *gin.Context) {
	componentHash := c.Param("hash")
	if componentHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Component hash is required"})
		return
	}

	var req struct {
		Requester string `json:"requester" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	metadata, err := h.blockchain.GetAnonymousComponentMetadata(ctx, req.Requester, componentHash)
	if err != nil {
		h.logger.Error().Err(err).Str("component_hash", componentHash).Msg("Failed to get anonymous component metadata")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get anonymous component metadata"})
		return
	}

	c.JSON(http.StatusOK, metadata)
}

// VerifyComponent handles component verification
func (h *Handler) VerifyComponent(c *gin.Context) {
	componentID := c.Param("id")
	if componentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Component ID is required"})
		return
	}

	var req struct {
		Verifier string `json:"verifier" binding:"required"`
		Context  string `json:"context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.VerifyComponent(ctx, req.Verifier, componentID, req.Context)
	if err != nil {
		h.logger.Error().Err(err).Str("component_id", componentID).Msg("Failed to verify component")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify component"})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"component_id": componentID,
			"verifier":     req.Verifier,
			"context":      req.Context,
			"timestamp":    time.Now().Unix(),
			"tx_hash":      resp["txhash"],
		}
		h.eventQueue.Emit("component_verified", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// InitiatePairing handles pairing initiation
func (h *Handler) InitiatePairing(c *gin.Context) {
	var req struct {
		Creator            string `json:"creator" binding:"required"`
		ComponentA         string `json:"component_a" binding:"required"`
		ComponentB         string `json:"component_b" binding:"required"`
		OperationalContext string `json:"operational_context"`
		ProxyID            string `json:"proxy_id"`
		ForceImmediate     bool   `json:"force_immediate"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.InitiatePairing(ctx, req.Creator, req.ComponentA, req.ComponentB, req.OperationalContext, req.ProxyID, req.ForceImmediate)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to initiate pairing")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to initiate pairing: %v", err)})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"challenge_id":        resp["challenge_id"],
			"creator":             req.Creator,
			"component_a":         req.ComponentA,
			"component_b":         req.ComponentB,
			"operational_context": req.OperationalContext,
			"proxy_id":            req.ProxyID,
			"force_immediate":     req.ForceImmediate,
			"timestamp":           time.Now().Unix(),
			"tx_hash":             resp["txhash"],
		}
		h.eventQueue.Emit("pairing_initiated", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// CompletePairing handles pairing completion
func (h *Handler) CompletePairing(c *gin.Context) {
	var req struct {
		Creator        string `json:"creator" binding:"required"`
		ChallengeID    string `json:"challenge_id" binding:"required"`
		ComponentAAuth string `json:"component_a_auth" binding:"required"`
		ComponentBAuth string `json:"component_b_auth" binding:"required"`
		SessionContext string `json:"session_context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.CompletePairing(ctx, req.Creator, req.ChallengeID, req.ComponentAAuth, req.ComponentBAuth, req.SessionContext)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to complete pairing")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete pairing"})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"challenge_id":    req.ChallengeID,
			"creator":         req.Creator,
			"session_context": req.SessionContext,
			"lct_id":          resp["lct_id"],
			"timestamp":       time.Now().Unix(),
			"tx_hash":         resp["txhash"],
		}
		h.eventQueue.Emit("pairing_completed", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// RevokePairing handles pairing revocation
func (h *Handler) RevokePairing(c *gin.Context) {
	var req struct {
		Creator       string `json:"creator" binding:"required"`
		LctID         string `json:"lct_id" binding:"required"`
		Reason        string `json:"reason"`
		NotifyOffline bool   `json:"notify_offline"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.RevokePairing(ctx, req.Creator, req.LctID, req.Reason, req.NotifyOffline)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to revoke pairing")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke pairing"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetPairingStatus handles pairing status retrieval
func (h *Handler) GetPairingStatus(c *gin.Context) {
	challengeID := c.Param("challenge_id")
	if challengeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Challenge ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	status, err := h.blockchain.GetPairingStatus(ctx, challengeID)
	if err != nil {
		h.logger.Error().Err(err).Str("challenge_id", challengeID).Msg("Failed to get pairing status")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pairing status"})
		return
	}

	c.JSON(http.StatusOK, status)
}

// CreateLCT handles LCT creation
func (h *Handler) CreateLCT(c *gin.Context) {
	var req struct {
		Creator    string `json:"creator" binding:"required"`
		ComponentA string `json:"component_a" binding:"required"`
		ComponentB string `json:"component_b" binding:"required"`
		Context    string `json:"context"`
		ProxyID    string `json:"proxy_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.CreateLCT(ctx, req.Creator, req.ComponentA, req.ComponentB, req.Context, req.ProxyID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create LCT")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create LCT: %v", err)})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"lct_id":      resp["lct_id"],
			"creator":     req.Creator,
			"component_a": req.ComponentA,
			"component_b": req.ComponentB,
			"context":     req.Context,
			"proxy_id":    req.ProxyID,
			"timestamp":   time.Now().Unix(),
			"tx_hash":     resp["txhash"],
		}
		h.eventQueue.Emit("lct_created", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// GetLCT handles LCT retrieval
func (h *Handler) GetLCT(c *gin.Context) {
	lctID := c.Param("id")
	if lctID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "LCT ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	lct, err := h.blockchain.GetLCT(ctx, lctID)
	if err != nil {
		h.logger.Error().Err(err).Str("lct_id", lctID).Msg("Failed to get LCT")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get LCT"})
		return
	}

	c.JSON(http.StatusOK, lct)
}

// UpdateLCTStatus handles LCT status updates
func (h *Handler) UpdateLCTStatus(c *gin.Context) {
	lctID := c.Param("id")
	if lctID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "LCT ID is required"})
		return
	}

	var req struct {
		Creator string `json:"creator" binding:"required"`
		Status  string `json:"status" binding:"required"`
		Context string `json:"context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.UpdateLCTStatus(ctx, req.Creator, lctID, req.Status, req.Context)
	if err != nil {
		h.logger.Error().Err(err).Str("lct_id", lctID).Msg("Failed to update LCT status")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update LCT status"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateTrustTensor handles trust tensor creation
func (h *Handler) CreateTrustTensor(c *gin.Context) {
	var req struct {
		Creator      string  `json:"creator" binding:"required"`
		ComponentA   string  `json:"component_a" binding:"required"`
		ComponentB   string  `json:"component_b" binding:"required"`
		Context      string  `json:"context"`
		InitialScore float64 `json:"initial_score"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.CreateTrustTensor(ctx, req.Creator, req.ComponentA, req.ComponentB, req.Context, req.InitialScore)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create trust tensor")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trust tensor"})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"tensor_id":     resp["tensor_id"],
			"creator":       req.Creator,
			"component_a":   req.ComponentA,
			"component_b":   req.ComponentB,
			"context":       req.Context,
			"initial_score": req.InitialScore,
			"timestamp":     time.Now().Unix(),
			"tx_hash":       resp["txhash"],
		}
		h.eventQueue.Emit("trust_tensor_created", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// GetTrustTensor handles trust tensor retrieval
func (h *Handler) GetTrustTensor(c *gin.Context) {
	tensorID := c.Param("id")
	if tensorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tensor ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	tensor, err := h.blockchain.GetTrustTensor(ctx, tensorID)
	if err != nil {
		h.logger.Error().Err(err).Str("tensor_id", tensorID).Msg("Failed to get trust tensor")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trust tensor"})
		return
	}

	c.JSON(http.StatusOK, tensor)
}

// UpdateTrustScore handles trust score updates
func (h *Handler) UpdateTrustScore(c *gin.Context) {
	tensorID := c.Param("id")
	if tensorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tensor ID is required"})
		return
	}

	var req struct {
		Creator string  `json:"creator" binding:"required"`
		Score   float64 `json:"score" binding:"required"`
		Context string  `json:"context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.UpdateTrustScore(ctx, req.Creator, tensorID, req.Score, req.Context)
	if err != nil {
		h.logger.Error().Err(err).Str("tensor_id", tensorID).Msg("Failed to update trust score")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trust score"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateEnergyOperation handles energy operation creation
func (h *Handler) CreateEnergyOperation(c *gin.Context) {
	var req struct {
		Creator       string  `json:"creator" binding:"required"`
		ComponentA    string  `json:"component_a" binding:"required"`
		ComponentB    string  `json:"component_b" binding:"required"`
		OperationType string  `json:"operation_type" binding:"required"`
		Amount        float64 `json:"amount" binding:"required"`
		Context       string  `json:"context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.CreateEnergyOperation(ctx, req.Creator, req.ComponentA, req.ComponentB, req.OperationType, req.Amount, req.Context)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create energy operation")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create energy operation"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ExecuteEnergyTransfer handles energy transfer execution
func (h *Handler) ExecuteEnergyTransfer(c *gin.Context) {
	var req struct {
		Creator     string  `json:"creator" binding:"required"`
		OperationID string  `json:"operation_id" binding:"required"`
		Amount      float64 `json:"amount" binding:"required"`
		Context     string  `json:"context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.ExecuteEnergyTransfer(ctx, req.Creator, req.OperationID, req.Amount, req.Context)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to execute energy transfer")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute energy transfer"})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"operation_id": req.OperationID,
			"creator":      req.Creator,
			"amount":       req.Amount,
			"context":      req.Context,
			"timestamp":    time.Now().Unix(),
			"tx_hash":      resp["txhash"],
		}
		h.eventQueue.Emit("energy_transfer", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// GetEnergyBalance handles energy balance retrieval
func (h *Handler) GetEnergyBalance(c *gin.Context) {
	componentID := c.Param("component_id")
	if componentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Component ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	balance, err := h.blockchain.GetEnergyBalance(ctx, componentID)
	if err != nil {
		h.logger.Error().Err(err).Str("component_id", componentID).Msg("Failed to get energy balance")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get energy balance"})
		return
	}

	c.JSON(http.StatusOK, balance)
}

// WebSocketHandler handles WebSocket connections for real-time events
func (h *Handler) WebSocketHandler(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to upgrade connection to WebSocket")
		return
	}
	defer conn.Close()

	h.logger.Info().Msg("WebSocket connection established")

	// Handle WebSocket connection
	for {
		// Read message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			h.logger.Error().Err(err).Msg("Failed to read WebSocket message")
			break
		}

		// Echo message back (for now)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			h.logger.Error().Err(err).Msg("Failed to write WebSocket message")
			break
		}
	}

	h.logger.Info().Msg("WebSocket connection closed")
}

// GetAccounts handles account listing
func (h *Handler) GetAccounts(c *gin.Context) {
	accountManager := h.blockchain.GetAccountManager()
	accounts := accountManager.ListAccounts()
	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
		"count":    len(accounts),
	})
}

// CreateAccount handles account creation
func (h *Handler) CreateAccount(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	accountManager := h.blockchain.GetAccountManager()
	account, err := accountManager.GetOrCreateAccount(ctx, req.Name)
	if err != nil {
		h.logger.Error().Err(err).Str("name", req.Name).Msg("Failed to create account")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create account: %v", err)})
		return
	}

	c.JSON(http.StatusOK, account)
}

// GetAccountInfo handles account information and usage
func (h *Handler) GetAccountInfo(c *gin.Context) {
	accountManager := h.blockchain.GetAccountManager()

	// Get default account
	defaultAccount := accountManager.GetDefaultAccount()

	// Get all accounts
	allAccounts := accountManager.ListAccounts()

	c.JSON(http.StatusOK, gin.H{
		"default_account": defaultAccount,
		"all_accounts":    allAccounts,
		"usage_info": gin.H{
			"message": "Transactions will use the best matching account for each creator",
			"examples": []gin.H{
				{"creator": "alice", "will_use": "alice account"},
				{"creator": "bob", "will_use": "bob account"},
				{"creator": "unknown", "will_use": "default account (alice)"},
			},
		},
	})
}

// TestIgniteCLI handles Ignite CLI testing
func (h *Handler) TestIgniteCLI(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Test Ignite CLI
	err := h.blockchain.TestIgniteCLI(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": "ignite_cli_unavailable",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ignite_cli_available",
		"message": "Ignite CLI is working correctly",
	})
}

// TestTransactionFormat handles transaction format testing
func (h *Handler) TestTransactionFormat(c *gin.Context) {
	// Create a test transaction
	message := map[string]interface{}{
		"@type":             "/racecarweb.componentregistry.v1.MsgRegisterComponent",
		"creator":           "alice",
		"component_id":      "TEST-001",
		"component_type":    "module",
		"manufacturer_data": "test-data",
	}

	// Create transaction file
	txFile, err := h.blockchain.CreateTransactionFile(message, "test-memo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer txFile.Close()

	// Read the file content
	content, err := os.ReadFile(txFile.Name())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction_file":    txFile.Name(),
		"transaction_content": string(content),
		"message":             message,
	})
}

// GetIgniteHelp handles Ignite CLI help information
func (h *Handler) GetIgniteHelp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	ignitePath := h.blockchain.GetIgnitePath()
	projectRoot := h.blockchain.GetProjectRoot()
	home := os.Getenv("HOME")
	gopath := os.Getenv("GOPATH")
	pathEnv := os.Getenv("PATH")

	cmd := exec.CommandContext(ctx, ignitePath, "tx", "--help")
	cmd.Env = append(os.Environ(),
		"HOME="+home,
		"GOPATH="+gopath,
		"PATH="+pathEnv,
	)
	cmd.Dir = projectRoot

	output, err := cmd.Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": "help_failed",
		})
		return
	}

	cmd2 := exec.CommandContext(ctx, ignitePath, "tx", "componentregistry", "--help")
	cmd2.Env = append(os.Environ(),
		"HOME="+home,
		"GOPATH="+gopath,
		"PATH="+pathEnv,
	)
	cmd2.Dir = projectRoot

	output2, err2 := cmd2.Output()
	if err2 != nil {
		c.JSON(http.StatusOK, gin.H{
			"tx_help":                string(output),
			"componentregistry_help": "Command not found or help not available",
			"status":                 "partial_help_available",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tx_help":                string(output),
		"componentregistry_help": string(output2),
		"status":                 "help_available",
	})
}

// GetBlockchainClient returns the blockchain client
func (h *Handler) GetBlockchainClient() *blockchain.Client {
	return h.blockchain
}

// Shutdown gracefully shuts down the handler
func (h *Handler) Shutdown() {
	if h.eventQueue != nil {
		h.eventQueue.Shutdown()
	}
}

// Queue Management Handlers

// QueuePairingRequest handles queue pairing request
func (h *Handler) QueuePairingRequest(c *gin.Context) {
	var req struct {
		ComponentA         string `json:"component_a" binding:"required"`
		ComponentB         string `json:"component_b" binding:"required"`
		OperationalContext string `json:"operational_context" binding:"required"`
		ProxyID            string `json:"proxy_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.QueuePairingRequest(ctx, req.ComponentA, req.ComponentB, req.OperationalContext, req.ProxyID)
	if err != nil {
		h.logger.Error().Err(err).Str("component_a", req.ComponentA).Str("component_b", req.ComponentB).Msg("Failed to queue pairing request")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to queue pairing request: %v", err)})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"request_id":          resp["request_id"],
			"component_a":         req.ComponentA,
			"component_b":         req.ComponentB,
			"operational_context": req.OperationalContext,
			"proxy_id":            req.ProxyID,
			"timestamp":           time.Now().Unix(),
			"tx_hash":             resp["txhash"],
		}
		h.eventQueue.Emit("pairing_request_queued", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// GetQueueStatus handles queue status retrieval
func (h *Handler) GetQueueStatus(c *gin.Context) {
	componentID := c.Param("id")
	if componentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Component ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	status, err := h.blockchain.GetQueueStatus(ctx, componentID)
	if err != nil {
		h.logger.Error().Err(err).Str("component_id", componentID).Msg("Failed to get queue status")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get queue status"})
		return
	}

	c.JSON(http.StatusOK, status)
}

// ProcessOfflineQueue handles offline queue processing
func (h *Handler) ProcessOfflineQueue(c *gin.Context) {
	componentID := c.Param("id")
	if componentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Component ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.ProcessOfflineQueue(ctx, componentID)
	if err != nil {
		h.logger.Error().Err(err).Str("component_id", componentID).Msg("Failed to process offline queue")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to process offline queue: %v", err)})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"component_id":       componentID,
			"processed_requests": resp["processed_requests"],
			"failed_requests":    resp["failed_requests"],
			"timestamp":          time.Now().Unix(),
			"tx_hash":            resp["txhash"],
		}
		h.eventQueue.Emit("offline_queue_processed", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// CancelRequest handles request cancellation
func (h *Handler) CancelRequest(c *gin.Context) {
	requestID := c.Param("id")
	if requestID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request ID is required"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.CancelRequest(ctx, requestID, req.Reason)
	if err != nil {
		h.logger.Error().Err(err).Str("request_id", requestID).Msg("Failed to cancel request")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to cancel request: %v", err)})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"request_id": requestID,
			"reason":     req.Reason,
			"timestamp":  time.Now().Unix(),
			"tx_hash":    resp["txhash"],
		}
		h.eventQueue.Emit("request_cancelled", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// GetQueuedRequests handles queued requests retrieval
func (h *Handler) GetQueuedRequests(c *gin.Context) {
	componentID := c.Param("id")
	if componentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Component ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	requests, err := h.blockchain.GetQueuedRequests(ctx, componentID)
	if err != nil {
		h.logger.Error().Err(err).Str("component_id", componentID).Msg("Failed to get queued requests")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get queued requests"})
		return
	}

	c.JSON(http.StatusOK, requests)
}

// ListProxyQueue handles proxy queue listing
func (h *Handler) ListProxyQueue(c *gin.Context) {
	proxyID := c.Param("id")
	if proxyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Proxy ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	queue, err := h.blockchain.ListProxyQueue(ctx, proxyID)
	if err != nil {
		h.logger.Error().Err(err).Str("proxy_id", proxyID).Msg("Failed to list proxy queue")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list proxy queue"})
		return
	}

	c.JSON(http.StatusOK, queue)
}

// Authorization Management Handlers

// CreatePairingAuthorization handles pairing authorization creation
func (h *Handler) CreatePairingAuthorization(c *gin.Context) {
	var req struct {
		ComponentA         string `json:"component_a" binding:"required"`
		ComponentB         string `json:"component_b" binding:"required"`
		OperationalContext string `json:"operational_context" binding:"required"`
		AuthorizationRules string `json:"authorization_rules" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.CreatePairingAuthorization(ctx, req.ComponentA, req.ComponentB, req.OperationalContext, req.AuthorizationRules)
	if err != nil {
		h.logger.Error().Err(err).Str("component_a", req.ComponentA).Str("component_b", req.ComponentB).Msg("Failed to create pairing authorization")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create pairing authorization: %v", err)})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"authorization_id":    resp["authorization_id"],
			"component_a":         req.ComponentA,
			"component_b":         req.ComponentB,
			"operational_context": req.OperationalContext,
			"authorization_rules": req.AuthorizationRules,
			"timestamp":           time.Now().Unix(),
			"tx_hash":             resp["txhash"],
		}
		h.eventQueue.Emit("pairing_authorization_created", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// GetComponentAuthorizations handles component authorizations retrieval
func (h *Handler) GetComponentAuthorizations(c *gin.Context) {
	componentID := c.Param("id")
	if componentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Component ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	authorizations, err := h.blockchain.GetComponentAuthorizations(ctx, componentID)
	if err != nil {
		h.logger.Error().Err(err).Str("component_id", componentID).Msg("Failed to get component authorizations")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get component authorizations"})
		return
	}

	c.JSON(http.StatusOK, authorizations)
}

// UpdateAuthorization handles authorization updates
func (h *Handler) UpdateAuthorization(c *gin.Context) {
	authorizationID := c.Param("id")
	if authorizationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization ID is required"})
		return
	}

	var req struct {
		OperationalContext string `json:"operational_context"`
		AuthorizationRules string `json:"authorization_rules"`
		Status             string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.OperationalContext != "" {
		updates["operational_context"] = req.OperationalContext
	}
	if req.AuthorizationRules != "" {
		updates["authorization_rules"] = req.AuthorizationRules
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.UpdateAuthorization(ctx, authorizationID, updates)
	if err != nil {
		h.logger.Error().Err(err).Str("authorization_id", authorizationID).Msg("Failed to update authorization")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update authorization: %v", err)})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"authorization_id": authorizationID,
			"updates":          updates,
			"timestamp":        time.Now().Unix(),
			"tx_hash":          resp["txhash"],
		}
		h.eventQueue.Emit("authorization_updated", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// RevokeAuthorization handles authorization revocation
func (h *Handler) RevokeAuthorization(c *gin.Context) {
	authorizationID := c.Param("id")
	if authorizationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization ID is required"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.RevokeAuthorization(ctx, authorizationID, req.Reason)
	if err != nil {
		h.logger.Error().Err(err).Str("authorization_id", authorizationID).Msg("Failed to revoke authorization")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to revoke authorization: %v", err)})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"authorization_id": authorizationID,
			"reason":           req.Reason,
			"timestamp":        time.Now().Unix(),
			"tx_hash":          resp["txhash"],
		}
		h.eventQueue.Emit("authorization_revoked", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// CheckPairingAuthorization handles pairing authorization checks
func (h *Handler) CheckPairingAuthorization(c *gin.Context) {
	componentA := c.Query("component_a")
	componentB := c.Query("component_b")
	operationalContext := c.Query("operational_context")

	if componentA == "" || componentB == "" || operationalContext == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "component_a, component_b, and operational_context are required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	result, err := h.blockchain.CheckPairingAuthorization(ctx, componentA, componentB, operationalContext)
	if err != nil {
		h.logger.Error().Err(err).Str("component_a", componentA).Str("component_b", componentB).Msg("Failed to check pairing authorization")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check pairing authorization"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Trust Tensor Enhanced Handlers

// CalculateRelationshipTrust handles relationship trust calculation
func (h *Handler) CalculateRelationshipTrust(c *gin.Context) {
	var req struct {
		ComponentA         string `json:"component_a" binding:"required"`
		ComponentB         string `json:"component_b" binding:"required"`
		OperationalContext string `json:"operational_context" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.CalculateRelationshipTrust(ctx, req.ComponentA, req.ComponentB, req.OperationalContext)
	if err != nil {
		h.logger.Error().Err(err).Str("component_a", req.ComponentA).Str("component_b", req.ComponentB).Msg("Failed to calculate relationship trust")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to calculate relationship trust: %v", err)})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"tensor_id":           resp["tensor_id"],
			"component_a":         req.ComponentA,
			"component_b":         req.ComponentB,
			"operational_context": req.OperationalContext,
			"trust_score":         resp["trust_score"],
			"timestamp":           time.Now().Unix(),
			"tx_hash":             resp["txhash"],
		}
		h.eventQueue.Emit("relationship_trust_calculated", eventData)
	}

	c.JSON(http.StatusOK, resp)
}

// GetRelationshipTensor handles relationship tensor retrieval
func (h *Handler) GetRelationshipTensor(c *gin.Context) {
	componentA := c.Query("component_a")
	componentB := c.Query("component_b")

	if componentA == "" || componentB == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "component_a and component_b are required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	tensor, err := h.blockchain.GetRelationshipTensor(ctx, componentA, componentB)
	if err != nil {
		h.logger.Error().Err(err).Str("component_a", componentA).Str("component_b", componentB).Msg("Failed to get relationship tensor")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get relationship tensor"})
		return
	}

	c.JSON(http.StatusOK, tensor)
}

// UpdateTensorScore handles tensor score updates
func (h *Handler) UpdateTensorScore(c *gin.Context) {
	var req struct {
		Creator    string  `json:"creator" binding:"required"`
		ComponentA string  `json:"component_a" binding:"required"`
		ComponentB string  `json:"component_b" binding:"required"`
		Score      float64 `json:"score" binding:"required"`
		Context    string  `json:"context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(h.config.Blockchain.Timeout)*time.Second)
	defer cancel()

	resp, err := h.blockchain.UpdateTensorScore(ctx, req.Creator, req.ComponentA, req.ComponentB, req.Score, req.Context)
	if err != nil {
		h.logger.Error().Err(err).Str("creator", req.Creator).Str("component_a", req.ComponentA).Str("component_b", req.ComponentB).Msg("Failed to update tensor score")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update tensor score: %v", err)})
		return
	}

	// Emit event if event queue is enabled
	if h.eventQueue != nil {
		eventData := map[string]interface{}{
			"tensor_id":   resp["tensor_id"],
			"creator":     req.Creator,
			"component_a": req.ComponentA,
			"component_b": req.ComponentB,
			"score":       req.Score,
			"context":     req.Context,
			"timestamp":   time.Now().Unix(),
			"tx_hash":     resp["txhash"],
		}
		h.eventQueue.Emit("tensor_score_updated", eventData)
	}

	c.JSON(http.StatusOK, resp)
}
