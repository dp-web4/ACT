package types

// QueueStatus represents comprehensive queue status for a component
type QueueStatus struct {
	ComponentId       string `json:"component_id"`
	Timestamp         int64  `json:"timestamp"`
	QueuedRequests    int32  `json:"queued_requests"`
	OfflineOperations int32  `json:"offline_operations"`
	TotalRetries      int32  `json:"total_retries"`
	FailedOperations  int32  `json:"failed_operations"`
	AverageRetryDelay int64  `json:"average_retry_delay"`
	LastProcessedAt   int64  `json:"last_processed_at"`
	QueueHealth       string `json:"queue_health"` // "healthy", "warning", "critical"
}
