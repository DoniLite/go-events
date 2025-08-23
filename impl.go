package goevents

import "reflect"

func (bus *EventFactory) CreateEvent(eventName string) *Event {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	newEvent := &Event{Name: eventName}
	for _, ev := range bus.eventGroup {
		if ev.Name == eventName {
			return ev
		}
	}
	bus.eventGroup = append(bus.eventGroup, newEvent)
	return newEvent
}

func (bus *EventFactory) On(event *Event, handler EventHandler) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	handlers := bus.registeredFunc[event]
	for _, fn := range handlers {
		if &fn == &handler {
			return
		}
	}
	bus.registeredFunc[event] = append(handlers, handler)
}

func (bus *EventFactory) Off(event *Event, handler EventHandler) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	handlers := bus.registeredFunc[event]
	filtered := make([]EventHandler, 0)

	for _, h := range handlers {
		if reflect.ValueOf(h).Pointer() != reflect.ValueOf(handler).Pointer() {
			filtered = append(filtered, h)
		}
	}

	if len(filtered) == 0 {
		delete(bus.registeredFunc, event)
	} else {
		bus.registeredFunc[event] = filtered
	}
}

func (bus *EventFactory) Emit(event *Event, data *EventData, args ...string) {
	bus.mu.Lock()
	handlers := bus.registeredFunc[event]
	bus.mu.Unlock()

	for _, handler := range handlers {
		bus.wg.Add(1)
		go func(fn EventHandler) {
			defer bus.wg.Done()
			fn(data, args...)
		}(handler)
	}
}

func (bus *EventFactory) Wait() {
	bus.wg.Wait()
}
