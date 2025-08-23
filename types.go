package goevents

import "sync"

type Event struct {
	Name string
}

type EventData struct {
	Message string
}

type EventHandler func(event *EventData, args ...string)

type EventFactory struct {
	Mu             *sync.Mutex
	Wg             *sync.WaitGroup
	eventGroup     []*Event
	registeredFunc map[*Event][]EventHandler
}
