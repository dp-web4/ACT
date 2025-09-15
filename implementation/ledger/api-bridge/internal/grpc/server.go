package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"api-bridge/internal/auth"
	"api-bridge/internal/blockchain"
	"api-bridge/internal/config"
	pb "api-bridge/proto"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedAPIBridgeServiceServer
	blockchainClient *blockchain.Client
	config           *config.Config
	logger           zerolog.Logger
	authInterceptor  *auth.GRPCAuthInterceptor
}

func NewServer(blockchainClient *blockchain.Client, config *config.Config) *Server {
	return &Server{
		blockchainClient: blockchainClient,
		config:           config,
	}
}

// SetLogger sets the logger for the server
func (s *Server) SetLogger(logger zerolog.Logger) {
	s.logger = logger
}

// SetAuthInterceptor sets the authentication interceptor
func (s *Server) SetAuthInterceptor(interceptor *auth.GRPCAuthInterceptor) {
	s.authInterceptor = interceptor
}

func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Create gRPC server with optional authentication interceptor
	var grpcServer *grpc.Server
	if s.authInterceptor != nil {
		grpcServer = grpc.NewServer(
			grpc.UnaryInterceptor(s.authInterceptor.UnaryInterceptor),
		)
	} else {
		grpcServer = grpc.NewServer()
	}

	log.Printf("Registering APIBridgeService on port %d", port)
	pb.RegisterAPIBridgeServiceServer(grpcServer, s)

	log.Printf("gRPC server starting on port %d", port)
	return grpcServer.Serve(lis)
}

// Account Management
func (s *Server) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	accountManager := s.blockchainClient.GetAccountManager()
	accounts := accountManager.ListAccounts()

	pbAccounts := make([]*pb.Account, len(accounts))
	for i, acc := range accounts {
		pbAccounts[i] = &pb.Account{
			Name:    acc.Name,
			Address: acc.Address,
			KeyType: acc.KeyType,
		}
	}

	return &pb.GetAccountsResponse{
		Accounts: pbAccounts,
		Count:    int32(len(pbAccounts)),
	}, nil
}

// Component Registry
func (s *Server) RegisterComponent(ctx context.Context, req *pb.RegisterComponentRequest) (*pb.RegisterComponentResponse, error) {
	result, err := s.blockchainClient.RegisterComponent(ctx, req.Creator, req.ComponentData, req.Context)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register component: %v", err)
	}

	// Extract values from map[string]interface{}
	componentID, _ := result["component_id"].(string)
	componentIdentity, _ := result["component_identity"].(string)
	componentData, _ := result["component_data"].(string)
	context, _ := result["context"].(string)
	creator, _ := result["creator"].(string)
	lctID, _ := result["lct_id"].(string)
	status, _ := result["status"].(string)
	txhash, _ := result["txhash"].(string)

	return &pb.RegisterComponentResponse{
		ComponentId:       componentID,
		ComponentIdentity: componentIdentity,
		ComponentData:     componentData,
		Context:           context,
		Creator:           creator,
		LctId:             lctID,
		Status:            status,
		Txhash:            txhash,
	}, nil
}

func (s *Server) GetComponent(ctx context.Context, req *pb.GetComponentRequest) (*pb.GetComponentResponse, error) {
	component, err := s.blockchainClient.GetComponent(ctx, req.ComponentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get component: %v", err)
	}

	// Extract values from map[string]interface{}
	componentID, _ := component["component_id"].(string)
	componentData, _ := component["component_data"].(string)
	context, _ := component["context"].(string)
	creator, _ := component["creator"].(string)
	status, _ := component["status"].(string)
	txhash, _ := component["txhash"].(string)

	return &pb.GetComponentResponse{
		ComponentId:   componentID,
		ComponentData: componentData,
		Context:       context,
		Creator:       creator,
		Status:        status,
		Txhash:        txhash,
	}, nil
}

func (s *Server) GetComponentIdentity(ctx context.Context, req *pb.GetComponentIdentityRequest) (*pb.GetComponentIdentityResponse, error) {
	identity, err := s.blockchainClient.GetComponentIdentity(ctx, req.ComponentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get component identity: %v", err)
	}

	// Extract values from map[string]interface{}
	identityValue, _ := identity["identity"].(string)
	status, _ := identity["status"].(string)

	return &pb.GetComponentIdentityResponse{
		ComponentId: req.ComponentId,
		Identity:    identityValue,
		Status:      status,
	}, nil
}

func (s *Server) VerifyComponent(ctx context.Context, req *pb.VerifyComponentRequest) (*pb.VerifyComponentResponse, error) {
	result, err := s.blockchainClient.VerifyComponent(ctx, req.Verifier, req.ComponentId, req.Context)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify component: %v", err)
	}

	// Extract values from map[string]interface{}
	componentID, _ := result["component_id"].(string)
	verifier, _ := result["verifier"].(string)
	status, _ := result["status"].(string)
	txhash, _ := result["txhash"].(string)

	return &pb.VerifyComponentResponse{
		ComponentId: componentID,
		Verifier:    verifier,
		Status:      status,
		Txhash:      txhash,
	}, nil
}

// LCT Management
func (s *Server) CreateLCT(ctx context.Context, req *pb.CreateLCTRequest) (*pb.CreateLCTResponse, error) {
	result, err := s.blockchainClient.CreateLCT(ctx, req.Creator, req.ComponentA, req.ComponentB, req.Context, req.ProxyId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create LCT: %v", err)
	}

	// Extract values from map[string]interface{}
	lctID, _ := result["lct_id"].(string)
	componentA, _ := result["component_a"].(string)
	componentB, _ := result["component_b"].(string)
	context, _ := result["context"].(string)
	proxyID, _ := result["proxy_id"].(string)
	status, _ := result["status"].(string)
	createdAt, _ := result["created_at"].(int64)
	creator, _ := result["creator"].(string)
	txhash, _ := result["txhash"].(string)
	lctKeyHalf, _ := result["lct_key_half"].(string)
	deviceKeyHalf, _ := result["device_key_half"].(string)

	return &pb.CreateLCTResponse{
		LctId:         lctID,
		ComponentA:    componentA,
		ComponentB:    componentB,
		Context:       context,
		ProxyId:       proxyID,
		Status:        status,
		CreatedAt:     createdAt,
		Creator:       creator,
		Txhash:        txhash,
		LctKeyHalf:    lctKeyHalf,
		DeviceKeyHalf: deviceKeyHalf,
	}, nil
}

func (s *Server) GetLCT(ctx context.Context, req *pb.GetLCTRequest) (*pb.GetLCTResponse, error) {
	lct, err := s.blockchainClient.GetLCT(ctx, req.LctId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get LCT: %v", err)
	}

	// Extract values from map[string]interface{}
	lctID, _ := lct["lct_id"].(string)
	componentAID, _ := lct["component_a_id"].(string)
	componentBID, _ := lct["component_b_id"].(string)
	lctKeyHalf, _ := lct["lct_key_half"].(string)
	pairingStatus, _ := lct["pairing_status"].(string)
	createdAt, _ := lct["created_at"].(int64)
	updatedAt, _ := lct["updated_at"].(int64)
	lastContactAt, _ := lct["last_contact_at"].(int64)
	trustAnchor, _ := lct["trust_anchor"].(string)
	operationalContext, _ := lct["operational_context"].(string)
	proxyComponentID, _ := lct["proxy_component_id"].(string)
	authorizationRules, _ := lct["authorization_rules"].(string)

	return &pb.GetLCTResponse{
		LctId:              lctID,
		ComponentAId:       componentAID,
		ComponentBId:       componentBID,
		LctKeyHalf:         lctKeyHalf,
		PairingStatus:      pairingStatus,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		LastContactAt:      lastContactAt,
		TrustAnchor:        trustAnchor,
		OperationalContext: operationalContext,
		ProxyComponentId:   proxyComponentID,
		AuthorizationRules: authorizationRules,
	}, nil
}

func (s *Server) UpdateLCTStatus(ctx context.Context, req *pb.UpdateLCTStatusRequest) (*pb.UpdateLCTStatusResponse, error) {
	result, err := s.blockchainClient.UpdateLCTStatus(ctx, req.Creator, req.LctId, req.Status, req.Context)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update LCT status: %v", err)
	}

	// Extract values from map[string]interface{}
	lctID, _ := result["lct_id"].(string)
	status, _ := result["status"].(string)
	updatedAt, _ := result["updated_at"].(int64)

	return &pb.UpdateLCTStatusResponse{
		LctId:     lctID,
		Status:    status,
		UpdatedAt: updatedAt,
	}, nil
}

// Pairing
func (s *Server) InitiatePairing(ctx context.Context, req *pb.InitiatePairingRequest) (*pb.InitiatePairingResponse, error) {
	result, err := s.blockchainClient.InitiatePairing(ctx, req.Creator, req.ComponentA, req.ComponentB, req.OperationalContext, req.ProxyId, req.ForceImmediate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to initiate pairing: %v", err)
	}

	// Extract values from map[string]interface{}
	challengeID, _ := result["challenge_id"].(string)
	componentA, _ := result["component_a"].(string)
	componentB, _ := result["component_b"].(string)
	operationalContext, _ := result["operational_context"].(string)
	proxyID, _ := result["proxy_id"].(string)
	forceImmediate, _ := result["force_immediate"].(bool)
	status, _ := result["status"].(string)
	createdAt, _ := result["created_at"].(int64)
	creator, _ := result["creator"].(string)
	txhash, _ := result["txhash"].(string)

	return &pb.InitiatePairingResponse{
		ChallengeId:        challengeID,
		ComponentA:         componentA,
		ComponentB:         componentB,
		OperationalContext: operationalContext,
		ProxyId:            proxyID,
		ForceImmediate:     forceImmediate,
		Status:             status,
		CreatedAt:          createdAt,
		Creator:            creator,
		Txhash:             txhash,
	}, nil
}

func (s *Server) CompletePairing(ctx context.Context, req *pb.CompletePairingRequest) (*pb.CompletePairingResponse, error) {
	result, err := s.blockchainClient.CompletePairing(ctx, req.Creator, req.ChallengeId, req.ComponentAAuth, req.ComponentBAuth, req.SessionContext)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to complete pairing: %v", err)
	}

	// Extract values from map[string]interface{}
	lctID, _ := result["lct_id"].(string)
	sessionKeys, _ := result["session_keys"].(string)
	trustSummary, _ := result["trust_summary"].(string)
	txhash, _ := result["txhash"].(string)
	splitKeyA, _ := result["split_key_a"].(string)
	splitKeyB, _ := result["split_key_b"].(string)

	return &pb.CompletePairingResponse{
		LctId:        lctID,
		SessionKeys:  sessionKeys,
		TrustSummary: trustSummary,
		Txhash:       txhash,
		SplitKeyA:    splitKeyA,
		SplitKeyB:    splitKeyB,
	}, nil
}

func (s *Server) RevokePairing(ctx context.Context, req *pb.RevokePairingRequest) (*pb.RevokePairingResponse, error) {
	result, err := s.blockchainClient.RevokePairing(ctx, req.Creator, req.LctId, req.Reason, req.NotifyOffline)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to revoke pairing: %v", err)
	}

	// Extract values from map[string]interface{}
	lctID, _ := result["lct_id"].(string)
	status, _ := result["status"].(string)
	reason, _ := result["reason"].(string)

	return &pb.RevokePairingResponse{
		LctId:  lctID,
		Status: status,
		Reason: reason,
	}, nil
}

func (s *Server) GetPairingStatus(ctx context.Context, req *pb.GetPairingStatusRequest) (*pb.GetPairingStatusResponse, error) {
	pairingStatus, err := s.blockchainClient.GetPairingStatus(ctx, req.ChallengeId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get pairing status: %v", err)
	}

	// Extract values from map[string]interface{}
	statusValue, _ := pairingStatus["status"].(string)
	createdAt, _ := pairingStatus["created_at"].(int64)

	return &pb.GetPairingStatusResponse{
		ChallengeId: req.ChallengeId,
		Status:      statusValue,
		CreatedAt:   createdAt,
	}, nil
}

// Trust Tensor
func (s *Server) CreateTrustTensor(ctx context.Context, req *pb.CreateTrustTensorRequest) (*pb.CreateTrustTensorResponse, error) {
	result, err := s.blockchainClient.CreateTrustTensor(ctx, req.Creator, req.ComponentA, req.ComponentB, req.Context, req.InitialScore)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create trust tensor: %v", err)
	}

	// Extract values from map[string]interface{}
	tensorID, _ := result["tensor_id"].(string)
	score, _ := result["score"].(float64)
	status, _ := result["status"].(string)
	txhash, _ := result["txhash"].(string)

	return &pb.CreateTrustTensorResponse{
		TensorId: tensorID,
		Score:    score,
		Status:   status,
		Txhash:   txhash,
	}, nil
}

func (s *Server) GetTrustTensor(ctx context.Context, req *pb.GetTrustTensorRequest) (*pb.GetTrustTensorResponse, error) {
	tensor, err := s.blockchainClient.GetTrustTensor(ctx, req.TensorId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get trust tensor: %v", err)
	}

	// Extract values from map[string]interface{}
	score, _ := tensor["score"].(float64)
	status, _ := tensor["status"].(string)

	return &pb.GetTrustTensorResponse{
		TensorId: req.TensorId,
		Score:    score,
		Status:   status,
	}, nil
}

func (s *Server) UpdateTrustScore(ctx context.Context, req *pb.UpdateTrustScoreRequest) (*pb.UpdateTrustScoreResponse, error) {
	result, err := s.blockchainClient.UpdateTrustScore(ctx, req.Creator, req.TensorId, req.Score, req.Context)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update trust score: %v", err)
	}

	// Extract values from map[string]interface{}
	tensorID, _ := result["tensor_id"].(string)
	score, _ := result["score"].(float64)
	updatedAt, _ := result["updated_at"].(int64)

	return &pb.UpdateTrustScoreResponse{
		TensorId:  tensorID,
		Score:     score,
		UpdatedAt: updatedAt,
	}, nil
}

// Energy Operations
func (s *Server) CreateEnergyOperation(ctx context.Context, req *pb.CreateEnergyOperationRequest) (*pb.CreateEnergyOperationResponse, error) {
	result, err := s.blockchainClient.CreateEnergyOperation(ctx, req.Creator, req.ComponentA, req.ComponentB, req.OperationType, req.Amount, req.Context)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create energy operation: %v", err)
	}

	// Extract values from map[string]interface{}
	operationID, _ := result["operation_id"].(string)
	operationType, _ := result["operation_type"].(string)
	amount, _ := result["amount"].(float64)
	status, _ := result["status"].(string)
	txhash, _ := result["txhash"].(string)

	return &pb.CreateEnergyOperationResponse{
		OperationId:   operationID,
		OperationType: operationType,
		Amount:        amount,
		Status:        status,
		Txhash:        txhash,
	}, nil
}

func (s *Server) ExecuteEnergyTransfer(ctx context.Context, req *pb.ExecuteEnergyTransferRequest) (*pb.ExecuteEnergyTransferResponse, error) {
	result, err := s.blockchainClient.ExecuteEnergyTransfer(ctx, req.Creator, req.OperationId, req.Amount, req.Context)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to execute energy transfer: %v", err)
	}

	// Extract values from map[string]interface{}
	operationID, _ := result["operation_id"].(string)
	amount, _ := result["amount"].(float64)
	status, _ := result["status"].(string)
	txhash, _ := result["txhash"].(string)

	return &pb.ExecuteEnergyTransferResponse{
		OperationId: operationID,
		Amount:      amount,
		Status:      status,
		Txhash:      txhash,
	}, nil
}

func (s *Server) GetEnergyBalance(ctx context.Context, req *pb.GetEnergyBalanceRequest) (*pb.GetEnergyBalanceResponse, error) {
	balance, err := s.blockchainClient.GetEnergyBalance(ctx, req.ComponentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get energy balance: %v", err)
	}

	// Extract values from map[string]interface{}
	balanceValue, _ := balance["balance"].(float64)
	unit, _ := balance["unit"].(string)
	lastUpdated, _ := balance["last_updated"].(int64)

	return &pb.GetEnergyBalanceResponse{
		ComponentId: req.ComponentId,
		Balance:     balanceValue,
		Unit:        unit,
		LastUpdated: lastUpdated,
	}, nil
}

// Real-time Battery Monitoring
func (s *Server) StreamBatteryStatus(req *pb.StreamBatteryStatusRequest, stream pb.APIBridgeService_StreamBatteryStatusServer) error {
	ctx := stream.Context()
	interval := time.Duration(req.IntervalSeconds) * time.Second
	if interval == 0 {
		interval = 5 * time.Second // Default to 5 seconds
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Simulate battery status updates
			// In a real implementation, this would read from actual battery sensors
			update := &pb.BatteryStatusUpdate{
				ComponentId:   req.ComponentId,
				Voltage:       12.5 + float64(time.Now().Unix()%10)*0.1, // Simulate voltage variation
				Current:       2.1 + float64(time.Now().Unix()%5)*0.2,   // Simulate current variation
				Temperature:   25.0 + float64(time.Now().Unix()%8)*0.5,  // Simulate temperature variation
				StateOfCharge: 85.0 + float64(time.Now().Unix()%15)*0.3, // Simulate SOC variation
				Timestamp:     time.Now().Unix(),
			}

			if err := stream.Send(update); err != nil {
				return status.Errorf(codes.Internal, "failed to send battery status: %v", err)
			}
		}
	}
}
