package types

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// LocalMRHStorage implements MRHStorage using local filesystem
// This is designed to be easily replaceable with IPFS
type LocalMRHStorage struct {
	basePath string
	mu       sync.RWMutex
	cache    map[string]*MRHGraph // In-memory cache
}

// NewLocalMRHStorage creates a new local storage instance
func NewLocalMRHStorage(basePath string) (*LocalMRHStorage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create MRH storage directory: %w", err)
	}
	
	return &LocalMRHStorage{
		basePath: basePath,
		cache:    make(map[string]*MRHGraph),
	}, nil
}

// Store saves an MRH graph and returns its content hash (mimics IPFS behavior)
func (s *LocalMRHStorage) Store(ctx context.Context, graph *MRHGraph) (string, error) {
	// Serialize graph to JSON
	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to serialize MRH graph: %w", err)
	}
	
	// Calculate content hash (same as IPFS would)
	hash := s.calculateHash(data)
	
	// Store to filesystem
	filePath := s.getFilePath(hash)
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write MRH graph to disk: %w", err)
	}
	
	// Update cache
	s.mu.Lock()
	s.cache[hash] = graph
	s.mu.Unlock()
	
	return hash, nil
}

// Retrieve gets an MRH graph by its content hash
func (s *LocalMRHStorage) Retrieve(ctx context.Context, hash string) (*MRHGraph, error) {
	// Check cache first
	s.mu.RLock()
	if cached, ok := s.cache[hash]; ok {
		s.mu.RUnlock()
		return cached, nil
	}
	s.mu.RUnlock()
	
	// Load from filesystem
	filePath := s.getFilePath(hash)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("MRH graph not found: %s", hash)
		}
		return nil, fmt.Errorf("failed to read MRH graph: %w", err)
	}
	
	// Deserialize
	var graph MRHGraph
	if err := json.Unmarshal(data, &graph); err != nil {
		return nil, fmt.Errorf("failed to deserialize MRH graph: %w", err)
	}
	
	// Update cache
	s.mu.Lock()
	s.cache[hash] = &graph
	s.mu.Unlock()
	
	return &graph, nil
}

// Delete removes an MRH graph
func (s *LocalMRHStorage) Delete(ctx context.Context, hash string) error {
	// Remove from cache
	s.mu.Lock()
	delete(s.cache, hash)
	s.mu.Unlock()
	
	// Remove from filesystem
	filePath := s.getFilePath(hash)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete MRH graph: %w", err)
	}
	
	return nil
}

// Exists checks if a graph exists
func (s *LocalMRHStorage) Exists(ctx context.Context, hash string) (bool, error) {
	// Check cache first
	s.mu.RLock()
	if _, ok := s.cache[hash]; ok {
		s.mu.RUnlock()
		return true, nil
	}
	s.mu.RUnlock()
	
	// Check filesystem
	filePath := s.getFilePath(hash)
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Helper methods

func (s *LocalMRHStorage) calculateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (s *LocalMRHStorage) getFilePath(hash string) string {
	// Use first 2 chars as subdirectory for better filesystem performance
	subdir := hash[:2]
	dirPath := filepath.Join(s.basePath, subdir)
	
	// Create subdirectory if needed
	os.MkdirAll(dirPath, 0755)
	
	return filepath.Join(dirPath, hash+".json")
}

// LocalMRHTraversal implements MRHTraversal for local storage
type LocalMRHTraversal struct {
	storage *LocalMRHStorage
	graphs  map[string]*MRHGraph // Cache of loaded graphs
	mu      sync.RWMutex
}

// NewLocalMRHTraversal creates a new traversal instance
func NewLocalMRHTraversal(storage *LocalMRHStorage) *LocalMRHTraversal {
	return &LocalMRHTraversal{
		storage: storage,
		graphs:  make(map[string]*MRHGraph),
	}
}

// FindPath finds the shortest trust path between two LCTs using BFS
func (t *LocalMRHTraversal) FindPath(from, to string, maxDepth uint32) ([]string, error) {
	if from == to {
		return []string{from}, nil
	}
	
	// BFS to find shortest path
	type node struct {
		lct   string
		path  []string
		depth uint32
	}
	
	queue := []node{{lct: from, path: []string{from}, depth: 0}}
	visited := make(map[string]bool)
	visited[from] = true
	
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		
		if current.depth >= maxDepth {
			continue
		}
		
		// Get neighbors from MRH graph
		neighbors, err := t.getNeighbors(current.lct)
		if err != nil {
			continue // Skip if can't get neighbors
		}
		
		for _, neighbor := range neighbors {
			if neighbor == to {
				// Found the target
				return append(current.path, to), nil
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
	
	return nil, fmt.Errorf("no path found between %s and %s within depth %d", from, to, maxDepth)
}

// CalculateTrust calculates trust between two LCTs based on MRH distance
func (t *LocalMRHTraversal) CalculateTrust(from, to string) (float64, error) {
	// Find shortest path
	path, err := t.FindPath(from, to, 6) // Max depth of 6
	if err != nil {
		return 0, err
	}
	
	// Calculate trust degradation based on distance
	distance := len(path) - 1
	if distance == 0 {
		return 1.0, nil // Self-trust is always 1.0
	}
	
	// Trust degrades exponentially with distance
	// Using formula: trust = baseТrust * (decayFactor ^ distance)
	baseTrust := 1.0
	decayFactor := 0.8 // 20% degradation per hop
	
	trust := baseTrust * pow(decayFactor, float64(distance))
	
	// Apply minimum trust threshold
	if trust < 0.01 {
		trust = 0.01
	}
	
	return trust, nil
}

// GetContext returns all LCTs within a context boundary
func (t *LocalMRHTraversal) GetContext(center string, radius uint32) (*MRHContext, error) {
	context := &MRHContext{
		CenterLCT:    center,
		Radius:       radius,
		TrustDecay:   0.2, // 20% per hop
		IncludedLCTs: []string{center},
	}
	
	// BFS to find all LCTs within radius
	visited := make(map[string]bool)
	visited[center] = true
	
	type node struct {
		lct   string
		depth uint32
	}
	
	queue := []node{{lct: center, depth: 0}}
	
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		
		if current.depth >= radius {
			continue
		}
		
		neighbors, err := t.getNeighbors(current.lct)
		if err != nil {
			continue
		}
		
		for _, neighbor := range neighbors {
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
	
	return context, nil
}

// GetWitnesses returns all witnesses for an LCT
func (t *LocalMRHTraversal) GetWitnesses(lctID string) ([]*WitnessRelationship, error) {
	// In a real implementation, this would query the blockchain
	// For now, return empty list
	return []*WitnessRelationship{}, nil
}

// Helper methods

func (t *LocalMRHTraversal) getNeighbors(lctID string) ([]string, error) {
	// This would normally load the graph for the LCT and extract neighbors
	// For now, return empty list
	return []string{}, nil
}

// Simple power function for trust calculation
func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}