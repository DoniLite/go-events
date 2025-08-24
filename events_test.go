package goevents

import (
	"sync"
	"testing"
	"time"
)

func TestCreateEvent_ShouldRegisterAndReturnSameEvent(t *testing.T) {
	bus := EventBus

	ev1 := bus.CreateEvent("test:event")
	ev2 := bus.CreateEvent("test:event")

	if ev1 != ev2 {
		t.Errorf("Expected same event instance, got different ones")
	}
}

func TestOn_ShouldRegisterHandler(t *testing.T) {
	bus := EventBus
	event := bus.CreateEvent("on:event")

	called := false
	handler := func(data *EventData, args ...string) {
		called = true
	}

	bus.On(event, handler)
	bus.Emit(event, &EventData{Message: "Hello"})
	bus.Wait()

	if !called {
		t.Errorf("Handler was not called after Emit")
	}
}

func TestEmit_ShouldCallAllHandlers(t *testing.T) {
	bus := EventBus
	event := bus.CreateEvent("multi:handler")

	var mu sync.Mutex
	calls := 0

	handler1 := func(data *EventData, args ...string) {
		mu.Lock()
		calls++
		mu.Unlock()
	}
	handler2 := func(data *EventData, args ...string) {
		mu.Lock()
		calls++
		mu.Unlock()
	}

	bus.On(event, handler1)
	bus.On(event, handler2)
	bus.Emit(event, &EventData{Message: "Event triggered"})
	bus.Wait()

	if calls != 2 {
		t.Errorf("Expected 2 handlers to be called, got %d", calls)
	}
}

func TestOff_ShouldRemoveHandler(t *testing.T) {
	bus := EventBus
	event := bus.CreateEvent("off:event")

	called := false
	handler := func(data *EventData, args ...string) {
		called = true
	}

	bus.On(event, handler)
	bus.Off(event, handler)
	bus.Emit(event, &EventData{Message: "Should not be called"})
	bus.Wait()

	if called {
		t.Errorf("Handler was called after being removed")
	}
}

func TestEmit_ShouldPassArguments(t *testing.T) {
	bus := EventBus
	event := bus.CreateEvent("args:event")

	var received string
	handler := func(data *EventData, args ...string) {
		if len(args) > 0 {
			received = args[0]
		}
	}

	bus.On(event, handler)
	bus.Emit(event, &EventData{}, "Doni")
	bus.Wait()

	if received != "Doni" {
		t.Errorf("Expected argument 'Doni', got '%s'", received)
	}
}

func TestEmit_ShouldBeAsynchronous(t *testing.T) {
	bus := EventBus
	event := bus.CreateEvent("async:event")

	start := time.Now()
	handler := func(data *EventData, args ...string) {
		time.Sleep(50 * time.Millisecond)
	}

	bus.On(event, handler)
	bus.Emit(event, &EventData{})
	bus.Wait()

	elapsed := time.Since(start)
	if elapsed > 100*time.Millisecond {
		t.Errorf("Emit took too long, expected async behavior")
	}
}

func TestSubscribe_ShouldRegisterEventForEach(t *testing.T) {
	eventNames := []string{"event:one", "event:two", "event:three"}
	var events []*Event
	bus := EventBus
	called := 0

	handler := func(data *EventData, args ...string) {
		called++
	}

	for _, name := range eventNames {
		event := bus.CreateEvent(name)
		events = append(events, event)
	}

	bus.Subscribe(handler)

	for _, event := range events {
		bus.Emit(event, &EventData{})
	}

	bus.Wait()

	if called != len(events) {
		t.Errorf("Expected handler to be called %d times, got %d", len(events), called)
	}
}

func TestSubscribe_ShouldRegisterEventForTargets(t *testing.T) {
	eventNames := []string{"event:one", "event:two", "event:three"}
	var events []*Event
	bus := EventBus
	called := 0

	handler := func(data *EventData, args ...string) {
		called++
	}

	for _, name := range eventNames {
		event := bus.CreateEvent(name)
		events = append(events, event)
	}

	bus.Subscribe(handler, events...)

	for _, event := range events {
		bus.Emit(event, &EventData{})
	}

	bus.Wait()

	if called != len(events) {
		t.Errorf("Expected handler to be called %d times, got %d", len(events), called)
	}
}

func TestDecodeDataPayload_ShouldReturnCorrectType(t *testing.T) {
	expected := "test"
	data := &EventData{Payload: expected}

	result, ok := DecodeDataPayload[string](data)

	if !ok {
		t.Errorf("Failed to decode payload")
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}
