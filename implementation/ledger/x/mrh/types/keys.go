package types

const (
	// ModuleName defines the module name
	ModuleName = "mrh"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_mrh"
	
	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

// Key prefixes for store
var (
	// MRHGraphPrefix is the prefix for MRH graph storage
	MRHGraphPrefix = []byte{0x01}
	
	// LCTMRHMappingPrefix maps LCT IDs to their MRH graph hashes
	LCTMRHMappingPrefix = []byte{0x02}
	
	// WitnessRelationshipPrefix stores witness relationships
	WitnessRelationshipPrefix = []byte{0x03}
	
	// ContextCachePrefix caches computed contexts
	ContextCachePrefix = []byte{0x04}
	
	// TrustPathCachePrefix caches trust paths between LCTs
	TrustPathCachePrefix = []byte{0x05}
)

// GetMRHGraphKey returns the store key for an MRH graph
func GetMRHGraphKey(hash string) []byte {
	return append(MRHGraphPrefix, []byte(hash)...)
}

// GetLCTMRHMappingKey returns the store key for LCT to MRH mapping
func GetLCTMRHMappingKey(lctID string) []byte {
	return append(LCTMRHMappingPrefix, []byte(lctID)...)
}

// GetWitnessRelationshipKey returns the store key for a witness relationship
func GetWitnessRelationshipKey(witnessLCT, subjectLCT string) []byte {
	key := append(WitnessRelationshipPrefix, []byte(witnessLCT)...)
	key = append(key, []byte(":")...)
	return append(key, []byte(subjectLCT)...)
}

// GetContextCacheKey returns the store key for a cached context
func GetContextCacheKey(centerLCT string, radius uint32) []byte {
	key := append(ContextCachePrefix, []byte(centerLCT)...)
	key = append(key, []byte(":")...)
	return append(key, uint32ToBytes(radius)...)
}

// GetTrustPathCacheKey returns the store key for a cached trust path
func GetTrustPathCacheKey(fromLCT, toLCT string) []byte {
	key := append(TrustPathCachePrefix, []byte(fromLCT)...)
	key = append(key, []byte(":")...)
	return append(key, []byte(toLCT)...)
}

// Helper function to convert uint32 to bytes
func uint32ToBytes(n uint32) []byte {
	return []byte{
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	}
}