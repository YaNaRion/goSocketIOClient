package main

import (
	"log"
)

func testHandler(event EventHandlerInt) {
	log.Printf("%v", event.Payload())
	event.Emit("testOn", "BONJOURS FROM SOCKET ON CLIENT SIDE")
}

func main() {
	port := "3030"
	ipAdd := "localhost"

	client := Connection(ipAdd, port)

	defer client.Close()
	// Establish a WebSocket connection

	// Handle messages from the server
	client.On("test", testHandler)
	err := client.Emit("test", "WESH ALORS")

	if err != nil {
		log.Println(err)
	}
	client.HandlerServerResponse()
}
