package types

// LCT status constants
const (
	// StatusPending indicates the LCT is in a pending state, typically during creation or key exchange
	StatusPending = "pending"

	// StatusActive indicates the LCT is active and operational
	StatusActive = "active"

	// StatusTerminated indicates the LCT has been terminated
	StatusTerminated = "terminated"

	// StatusInactive indicates the LCT is inactive but not terminated
	StatusInactive = "inactive"

	// StatusMaintenance indicates the LCT is in maintenance mode
	StatusMaintenance = "maintenance"
)

// IsValidLCTStatus checks if the given status is a valid LCT status
func IsValidLCTStatus(status string) bool {
	switch status {
	case StatusPending, StatusActive, StatusTerminated, StatusInactive, StatusMaintenance:
		return true
	default:
		return false
	}
}
