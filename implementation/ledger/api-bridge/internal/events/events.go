package events

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Event represents a generic event to be emitted
// Extend Data as needed for your use case
// Type: event type string (e.g. "component_registered")
// Data: event payload (should be serializable)
type Event struct {
	Type      string      `json:"event_type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
	Attempts  int         `json:"-"` // for retry logic
}

// EventQueue manages event emission and retries
// Only enabled if sinks is non-nil and non-empty
// Sinks: map of event type to list of endpoint URLs
// MaxRetries: max attempts per event
// Backoff: initial backoff duration (doubles each retry)
type EventQueue struct {
	sinks      map[string][]string
	maxRetries int
	backoff    time.Duration
	queue      chan *Event
	logger     zerolog.Logger
	wg         sync.WaitGroup
	quit       chan struct{}
	enabled    bool
}

// NewEventQueue creates a new event queue (enabled only if sinks is non-empty)
func NewEventQueue(sinks map[string][]string, maxRetries int, backoff time.Duration, logger zerolog.Logger) *EventQueue {
	enabled := len(sinks) > 0
	eq := &EventQueue{
		sinks:      sinks,
		maxRetries: maxRetries,
		backoff:    backoff,
		queue:      make(chan *Event, 100),
		logger:     logger,
		quit:       make(chan struct{}),
		enabled:    enabled,
	}
	if enabled {
		eq.wg.Add(1)
		go eq.worker()
		logger.Info().Msg("Event queue enabled and worker started")
	} else {
		logger.Info().Msg("Event queue not enabled (no sinks configured)")
	}
	return eq
}

// Emit adds an event to the queue (no-op if not enabled)
func (eq *EventQueue) Emit(eventType string, data interface{}) {
	if !eq.enabled {
		return
	}
	eq.queue <- &Event{
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Data:      data,
		Attempts:  0,
	}
}

// worker processes the event queue
func (eq *EventQueue) worker() {
	defer eq.wg.Done()
	for {
		select {
		case event := <-eq.queue:
			eq.processEvent(event)
		case <-eq.quit:
			return
		}
	}
}

// processEvent POSTs the event to all sinks, with retries
func (eq *EventQueue) processEvent(event *Event) {
	endpoints := eq.sinks[event.Type]
	if len(endpoints) == 0 {
		eq.logger.Debug().Str("event", event.Type).Msg("No endpoints configured for event")
		return
	}
	payload, _ := json.Marshal(event)
	for _, url := range endpoints {
		success := false
		attempts := event.Attempts
		backoff := eq.backoff
		for attempts < eq.maxRetries {
			resp, err := http.Post(url, "application/json", bytes.NewReader(payload))
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				eq.logger.Info().Str("endpoint", url).Str("event", event.Type).Msg("Event POSTed successfully")
				resp.Body.Close()
				success = true
				break
			}
			if resp != nil {
				resp.Body.Close()
			}
			eq.logger.Warn().Str("endpoint", url).Str("event", event.Type).Int("attempt", attempts+1).Err(err).Msg("Failed to POST event, will retry")
			time.Sleep(backoff)
			backoff *= 2
			attempts++
		}
		if !success {
			eq.logger.Error().Str("endpoint", url).Str("event", event.Type).Msg("Event delivery failed after max retries")
		}
	}
}

// Shutdown gracefully stops the worker
func (eq *EventQueue) Shutdown() {
	if !eq.enabled {
		return
	}
	close(eq.quit)
	eq.wg.Wait()
}
