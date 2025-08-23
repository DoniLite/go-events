package goevents

import "sync"

var (
	EventBus *EventFactory
)

func init() {
	EventBus = &EventFactory{
		Mu:             &sync.Mutex{},
		Wg:             &sync.WaitGroup{},
		eventGroup:     []*Event{},
		registeredFunc: make(map[*Event][]EventHandler),
	}
}

func NewEventBus() *EventFactory {
	return &EventFactory{
		Mu:             &sync.Mutex{},
		Wg:             &sync.WaitGroup{},
		eventGroup:     []*Event{},
		registeredFunc: make(map[*Event][]EventHandler),
	}
}