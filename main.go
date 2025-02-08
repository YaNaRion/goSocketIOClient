package main

import (
	"log"
)

func trash(data *string) {
	if data != nil {
		log.Println(*data)
	}
}

func main() {
	port := "3030"
	ipAdd := "localhost"

	factoClient := NewFactoryClient(ipAdd, port)
	client := factoClient.NewClient()
	client.ConnectionWebSocket()
	defer client.Close()
	// Establish a WebSocket connection

	// Handle messages from the server
	client.On("test", trash)
	err := client.Emit("test", nil)

	if err != nil {
		log.Println(err)
	}
	client.HandlerServerResponse()
}
