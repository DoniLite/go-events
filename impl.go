package goevents

import "reflect"

// Create an event that can be fire with the event bus
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

// register an event to the event bus based on the provided event notice that the provided handler will not be replaced or register again 
// if it already exist if you want to replace any function use the replace method or unregistered the handler before
func (bus *EventFactory) On(event *Event, handler EventHandler) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	handlers := bus.registeredFunc[event]
	for _, fn := range handlers {
		if reflect.ValueOf(fn).Pointer() == reflect.ValueOf(handler).Pointer() {
			return
		}
	}
	bus.registeredFunc[event] = append(handlers, handler)
}

// Unregister an event handler
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

// Emit the provided event and call all the registered handler
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

// Wait until all the registered handler func to be called an executed 
// Notice that this method is not async and calling the internal
// 
//	WaitGroup.Wait()
//
// Avoid calling this function as possible due to the fact that all handler are called in goroutines
// But this method can be helpful for testing purposes or if you want to get the result of the handlers
func (bus *EventFactory) Wait() {
	bus.wg.Wait()
}

// Subscribe a function to the event bus if you provide some events your handler will be called for all the events
// But if no event list is provided the handler will be called for all the events
func (bus *EventFactory) Subscribe(fn EventHandler, targetEvents ...*Event) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if len(targetEvents) == 0 {
		for _, ev := range bus.eventGroup {
			bus.registeredFunc[ev] = append(bus.registeredFunc[ev], fn)
		}
		return
	}

	for _, target := range targetEvents {
		if _, ok := bus.registeredFunc[target]; ok {
			bus.registeredFunc[target] = append(bus.registeredFunc[target], fn)
		}
	}
}
