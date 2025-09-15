package keeper

import (
	"context"
	"fmt"
	"time"
	
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"racecarweb/x/mrh/types"
)

// Keeper maintains the MRH module state
type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	logger       log.Logger
	
	// MRH storage backend (local or IPFS)
	mrhStorage   types.MRHStorage
	mrhTraversal *types.LocalMRHTraversal
	
	// Cache for frequently accessed data
	graphCache   map[string]*types.MRHGraph
	pathCache    map[string][]string // Cached trust paths
}

// NewKeeper creates a new MRH keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	storageConfig types.StorageConfig,
) (*Keeper, error) {
	// Initialize storage backend
	storage, err := types.NewMRHStorage(storageConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MRH storage: %w", err)
	}
	
	// Initialize traversal (assumes local storage for now)
	var traversal *types.LocalMRHTraversal
	if localStorage, ok := storage.(*types.LocalMRHStorage); ok {
		traversal = types.NewLocalMRHTraversal(localStorage)
	}
	
	return &Keeper{
		cdc:          cdc,
		storeService: storeService,
		logger:       logger,
		mrhStorage:   storage,
		mrhTraversal: traversal,
		graphCache:   make(map[string]*types.MRHGraph),
		pathCache:    make(map[string][]string),
	}, nil
}

// Logger returns the module logger
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", "x/"+types.ModuleName)
}

// CreateMRHGraph creates a new MRH graph for an LCT
func (k Keeper) CreateMRHGraph(ctx context.Context, lctID string) (*types.MRHGraph, error) {
	// Create new graph
	graph := &types.MRHGraph{
		LctID:        lctID,
		Version:      1,
		CreatedAt:    time.Now().Unix(),
		ContextDepth: 1, // Start with depth 1
		Triples:      []types.Triple{},
		Metadata:     make(map[string]string),
	}
	
	// Add self-reference triple
	graph.AddTriple(lctID, "rdf:type", "web4:LCT", 1.0)
	graph.AddTriple(lctID, "web4:hasContext", lctID, 1.0)
	
	// Store the graph
	hash, err := k.mrhStorage.Store(ctx, graph)
	if err != nil {
		return nil, fmt.Errorf("failed to store MRH graph: %w", err)
	}
	
	// Store LCT to graph mapping
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.storeService.OpenKVStore(sdkCtx)
	mappingKey := types.GetLCTMRHMappingKey(lctID)
	if err := store.Set(mappingKey, []byte(hash)); err != nil {
		return nil, fmt.Errorf("failed to store LCT-MRH mapping: %w", err)
	}
	
	// Cache the graph
	k.graphCache[hash] = graph
	
	return graph, nil
}

// GetMRHGraph retrieves the MRH graph for an LCT
func (k Keeper) GetMRHGraph(ctx context.Context, lctID string) (*types.MRHGraph, error) {
	// Get graph hash from mapping
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.storeService.OpenKVStore(sdkCtx)
	mappingKey := types.GetLCTMRHMappingKey(lctID)
	
	hashBytes, err := store.Get(mappingKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get LCT-MRH mapping: %w", err)
	}
	if hashBytes == nil {
		return nil, types.ErrGraphNotFound
	}
	
	hash := string(hashBytes)
	
	// Check cache first
	if cached, ok := k.graphCache[hash]; ok {
		return cached, nil
	}
	
	// Retrieve from storage
	graph, err := k.mrhStorage.Retrieve(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve MRH graph: %w", err)
	}
	
	// Cache it
	k.graphCache[hash] = graph
	
	return graph, nil
}

// AddRelationship adds a relationship triple to an LCT's MRH graph
func (k Keeper) AddRelationship(
	ctx context.Context,
	subjectLCT string,
	predicate string,
	objectLCT string,
	weight float64,
) error {
	// Get the subject's graph
	graph, err := k.GetMRHGraph(ctx, subjectLCT)
	if err != nil {
		// Create new graph if it doesn't exist
		graph, err = k.CreateMRHGraph(ctx, subjectLCT)
		if err != nil {
			return fmt.Errorf("failed to create MRH graph: %w", err)
		}
	}
	
	// Add the triple
	graph.AddTriple(subjectLCT, predicate, objectLCT, weight)
	graph.Version++
	
	// Store updated graph
	hash, err := k.mrhStorage.Store(ctx, graph)
	if err != nil {
		return fmt.Errorf("failed to store updated MRH graph: %w", err)
	}
	
	// Update mapping
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.storeService.OpenKVStore(sdkCtx)
	mappingKey := types.GetLCTMRHMappingKey(subjectLCT)
	if err := store.Set(mappingKey, []byte(hash)); err != nil {
		return fmt.Errorf("failed to update LCT-MRH mapping: %w", err)
	}
	
	// Clear cache entries that might be affected
	delete(k.graphCache, hash)
	k.clearPathCache(subjectLCT)
	
	return nil
}

// CalculateContextBoundary calculates the context boundary for an LCT
func (k Keeper) CalculateContextBoundary(
	ctx context.Context,
	centerLCT string,
	radius uint32,
) (*types.MRHContext, error) {
	if radius > 10 {
		return nil, types.ErrContextTooLarge
	}
	
	// Check cache first
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := k.storeService.OpenKVStore(sdkCtx)
	cacheKey := types.GetContextCacheKey(centerLCT, radius)
	
	// For now, calculate fresh each time (implement caching later)
	context := &types.MRHContext{
		CenterLCT:    centerLCT,
		Radius:       radius,
		TrustDecay:   0.2, // 20% decay per hop
		IncludedLCTs: []string{centerLCT},
	}
	
	// BFS to find all LCTs within radius
	visited := make(map[string]bool)
	visited[centerLCT] = true
	
	type node struct {
		lct   string
		depth uint32
	}
	
	queue := []node{{lct: centerLCT, depth: 0}}
	
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		
		if current.depth >= radius {
			continue
		}
		
		// Get the graph for current LCT
		graph, err := k.GetMRHGraph(ctx, current.lct)
		if err != nil {
			continue // Skip if can't get graph
		}
		
		// Find all connected LCTs
		for _, triple := range graph.Triples {
			var neighbor string
			if triple.Subject == current.lct && triple.Object != current.lct {
				neighbor = triple.Object
			} else if triple.Object == current.lct && triple.Subject != current.lct {
				neighbor = triple.Subject
			} else {
				continue
			}
			
			if !visited[neighbor] {
				visited[neighbor] = true
				context.IncludedLCTs = append(context.IncludedLCTs, neighbor)
				queue = append(queue, node{
					lct:   neighbor,
					depth: current.depth + 1,
				})
			}
		}
	}
	
	// Store in cache
	contextBytes, _ := context.ToJSON()
	store.Set(cacheKey, contextBytes)
	
	return context, nil
}

// CalculateTrustPath finds the trust path between two LCTs
func (k Keeper) CalculateTrustPath(
	ctx context.Context,
	fromLCT string,
	toLCT string,
	maxDepth uint32,
) ([]string, float64, error) {
	if fromLCT == toLCT {
		return []string{fromLCT}, 1.0, nil
	}
	
	// Check cache
	cacheKey := fmt.Sprintf("%s->%s", fromLCT, toLCT)
	if cached, ok := k.pathCache[cacheKey]; ok {
		trust := k.calculateTrustFromPath(cached)
		return cached, trust, nil
	}
	
	// BFS to find shortest path
	type node struct {
		lct   string
		path  []string
		depth uint32
	}
	
	queue := []node{{lct: fromLCT, path: []string{fromLCT}, depth: 0}}
	visited := make(map[string]bool)
	visited[fromLCT] = true
	
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		
		if current.depth >= maxDepth {
			continue
		}
		
		// Get neighbors
		neighbors, err := k.getNeighbors(ctx, current.lct)
		if err != nil {
			continue
		}
		
		for _, neighbor := range neighbors {
			if neighbor == toLCT {
				// Found the target
				path := append(current.path, toLCT)
				trust := k.calculateTrustFromPath(path)
				
				// Cache the result
				k.pathCache[cacheKey] = path
				
				return path, trust, nil
			}
			
			if !visited[neighbor] {
				visited[neighbor] = true
				newPath := append([]string{}, current.path...)
				newPath = append(newPath, neighbor)
				queue = append(queue, node{
					lct:   neighbor,
					path:  newPath,
					depth: current.depth + 1,
				})
			}
		}
	}
	
	return nil, 0, types.ErrPathNotFound
}

// Helper methods

func (k Keeper) getNeighbors(ctx context.Context, lctID string) ([]string, error) {
	graph, err := k.GetMRHGraph(ctx, lctID)
	if err != nil {
		return nil, err
	}
	
	neighbors := make(map[string]bool)
	for _, triple := range graph.Triples {
		if triple.Subject == lctID && triple.Object != lctID {
			neighbors[triple.Object] = true
		} else if triple.Object == lctID && triple.Subject != lctID {
			neighbors[triple.Subject] = true
		}
	}
	
	result := make([]string, 0, len(neighbors))
	for neighbor := range neighbors {
		result = append(result, neighbor)
	}
	
	return result, nil
}

func (k Keeper) calculateTrustFromPath(path []string) float64 {
	if len(path) == 0 {
		return 0
	}
	if len(path) == 1 {
		return 1.0 // Self-trust
	}
	
	// Trust degrades exponentially with distance
	distance := len(path) - 1
	baseTrust := 1.0
	decayFactor := 0.8 // 20% degradation per hop
	
	trust := baseTrust
	for i := 0; i < distance; i++ {
		trust *= decayFactor
	}
	
	// Minimum trust threshold
	if trust < 0.01 {
		trust = 0.01
	}
	
	return trust
}

func (k Keeper) clearPathCache(lctID string) {
	// Clear all cached paths involving this LCT
	for key := range k.pathCache {
		if containsLCT(key, lctID) {
			delete(k.pathCache, key)
		}
	}
}

func containsLCT(pathKey, lctID string) bool {
	// Simple check - improve this
	return len(pathKey) > 0
}

// ToJSON helper for MRHContext
func (c *types.MRHContext) ToJSON() ([]byte, error) {
	// Implement JSON serialization
	return nil, nil
}