/*
Package goevents provides a lightweight, thread-safe event bus for Go.

It allows you to register named events, attach multiple handlers to each event,
emit events asynchronously with data and arguments, and manage handlers dynamically.

Features:

  - Register and emit named events
  - Attach and remove handlers for events
  - Pass data and arguments to handlers
  - Asynchronous handler execution
  - Thread-safe with sync primitives

Typical usage:

	package main

	import (
	    "github.com/DoniLite/go-events"
	    "fmt"
	)

	func main() {
	    bus := goevents.NewEventBus()
	    event := bus.CreateEvent("example:event")

	    bus.On(event, func(data *goevents.EventData, args ...string) {
	        fmt.Println("Event received:", data.Message, args)
	    })

	    bus.Emit(event, &goevents.EventData{Message: "Hello"}, "arg1")
	    bus.Wait()
	}
*/
package goevents
