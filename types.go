package goevents

import "sync"

// A single event in the event bus
type Event struct {
	Name string
}

// The data associated with an event
// This data will be passed to the event handlers when the event is emitted
type EventData struct {
	Message string
	Payload any
}

// A function that will be called when an event is emitted
type EventHandler func(event *EventData, args ...string)

// EventFactory is responsible for creating and managing events
type EventFactory struct {
	mu             *sync.Mutex
	wg             *sync.WaitGroup
	eventGroup     []*Event
	registeredFunc map[*Event][]EventHandler
}

// Decodes the payload of an EventData into the specified type
func DecodeDataPayload[T any](data *EventData) (T, bool) {
	payload, ok := data.Payload.(T)
	return payload, ok
}