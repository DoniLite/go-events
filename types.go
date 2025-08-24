package goevents

import "sync"

type Event struct {
	Name string
}

type EventData struct {
	Message string
	Payload any
}

type EventHandler func(event *EventData, args ...string)

type EventFactory struct {
	mu             *sync.Mutex
	wg             *sync.WaitGroup
	eventGroup     []*Event
	registeredFunc map[*Event][]EventHandler
}


func DecodeDataPayload[T any](data *EventData) T {
	return  data.Payload.(T)
}