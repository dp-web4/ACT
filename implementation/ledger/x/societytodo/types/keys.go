package types

const (
	// ModuleName defines the module name
	ModuleName = "societytodo"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_societytodo"
)

var (
	// Key prefixes
	SocietyTodoListPrefix = []byte{0x01}
	CitizenRequestPrefix  = []byte{0x02}
	DelegationPrefix      = []byte{0x03}
	CycleHistoryPrefix    = []byte{0x04}
	FederationSharePrefix = []byte{0x05}

	// Counters
	TodoSequenceKey      = []byte{0x10}
	RequestSequenceKey   = []byte{0x11}
	DelegationSequenceKey = []byte{0x12}
)

// GetSocietyTodoListKey returns the store key for a society's todo list
func GetSocietyTodoListKey(societyLCT string) []byte {
	return append(SocietyTodoListPrefix, []byte(societyLCT)...)
}

// GetCitizenRequestKey returns the store key for a citizen request
func GetCitizenRequestKey(requestID string) []byte {
	return append(CitizenRequestPrefix, []byte(requestID)...)
}

// GetDelegationKey returns the store key for an ATP delegation
func GetDelegationKey(citizenLCT, poolID string) []byte {
	return append(DelegationPrefix, append([]byte(citizenLCT), []byte(poolID)...)...)
}

// GetCycleHistoryKey returns the store key for cycle history
func GetCycleHistoryKey(societyLCT string, cycleNumber uint64) []byte {
	return append(CycleHistoryPrefix, append([]byte(societyLCT), UintToBytes(cycleNumber)...)...)
}

// UintToBytes converts uint64 to bytes
func UintToBytes(n uint64) []byte {
	bz := make([]byte, 8)
	for i := 0; i < 8; i++ {
		bz[7-i] = byte(n >> (8 * i))
	}
	return bz
}