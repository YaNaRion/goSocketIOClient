package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Hello world")

	port := "3030"
	ipAdd := "localhost"

	factoClient := NewFactoryClient(ipAdd, port)
	client := factoClient.NewClient()
	client.ConnectionWebSocket()
	defer client.Close()
	// Establish a WebSocket connection

	// Handle messages from the server

	eventTest := "test"
	testMessage := newSocketIOMessage(EVENT, &eventTest, nil)
	testMessageByte, err := testMessage.messageToMapOfByte()
	client.WriteMessage(testMessageByte)
	if err != nil {
		log.Println(err)
	}
	client.HandlerServerResponse()
}
