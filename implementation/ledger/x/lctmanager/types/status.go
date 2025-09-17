package types

const (
    StatusPending    = "pending"
    StatusActive     = "active"
    StatusTerminated = "terminated"
)

func IsValidLCTStatus(status string) bool {
    switch status {
    case StatusPending, StatusActive, StatusTerminated:
        return true
    default:
        return false
    }
}