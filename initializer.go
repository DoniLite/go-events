package goevents

import "sync"

var (
	EventBus *EventFactory
)

func init() {
	EventBus = &EventFactory{
		mu:             &sync.Mutex{},
		wg:             &sync.WaitGroup{},
		eventGroup:     []*Event{},
		registeredFunc: make(map[*Event][]EventHandler),
	}
}

func NewEventBus() *EventFactory {
	return &EventFactory{
		mu:             &sync.Mutex{},
		wg:             &sync.WaitGroup{},
		eventGroup:     []*Event{},
		registeredFunc: make(map[*Event][]EventHandler),
	}
}