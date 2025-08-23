
# go-events

A lightweight, thread-safe event bus for Go, supporting asynchronous event handling and handler management.

## Features

- Register and emit named events
- Attach multiple handlers to events
- Remove handlers dynamically
- Pass data and arguments to handlers
- Asynchronous handler execution
- Thread-safe with `sync.Mutex` and `sync.WaitGroup`

## Installation

```sh
go get github.com/DoniLite/go-events
```

## Usage

```go
package main

import (
    "github.com/DoniLite/go-events"
    "fmt"
)

func main() {
    // Create or use the default event bus
    bus := goevents.NewEventBus()

    // Create an event
    event := bus.CreateEvent("user:created")

    // Register a handler
    bus.On(event, func(data *goevents.EventData, args ...string) {
        fmt.Println("User created:", data.Message, args)
    })

    // Emit the event
    bus.Emit(event, &goevents.EventData{Message: "Alice"}, "extraArg")
    bus.Wait() // Wait for all handlers to finish
}
```

## API

### Types

- `Event`: Represents an event with a name.
- `EventData`: Data passed to handlers.
- `EventHandler`: Handler function signature.
- `EventFactory`: The event bus.

### Functions

- `CreateEvent(name string) *Event`: Register or get an event.
- `On(event *Event, handler EventHandler)`: Register a handler.
- `Off(event *Event, handler EventHandler)`: Remove a handler.
- `Emit(event *Event, data *EventData, args ...string)`: Emit an event asynchronously.
- `Wait()`: Wait for all handlers to complete.
- `NewEventBus() *EventFactory`: Create a new event bus instance.

## Testing

Run all tests:

```sh
go test ./...
```

## License

[Apache License](./LICENSE)
