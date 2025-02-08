package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const MessageWS = "4"

const (
	CONNECT       = "0"
	DISCONNECT    = "1"
	EVENT         = "2"
	ACK           = "3"
	CONNECT_ERROR = "4"
	BINARY_ACK    = "5"
)

type FactoryClient struct {
	Port   string
	IpAddr string
}

func NewFactoryClient(Ipaddr string, port string) *FactoryClient {
	return &FactoryClient{
		Port:   port,
		IpAddr: Ipaddr,
	}
}

func (f *FactoryClient) createURLConnection() string {
	return fmt.Sprintf("ws://%s:%s/socket.io?EIO=4&transport=websocket", f.IpAddr, f.Port)
}

func (f *FactoryClient) NewClient() *Client {
	var wsConn *websocket.Conn
	url := f.createURLConnection()
	return &Client{
		conn:    wsConn,
		url:     url,
		OnEvent: make(map[string]func(*string)),
	}
}

type Client struct {
	url     string
	conn    *websocket.Conn
	OnEvent map[string]func(*string)
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Emit(event string, message interface{}) error {
	connectionMessage := newSocketIOMessage(EVENT, &event, message)
	connectionMessageByte, err := connectionMessage.messageToMapOfByte()
	if err != nil {
		log.Println(err)
		return err
	}
	err = c.writeMessage(connectionMessageByte)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) On(event string, handler func(*string)) {

	c.OnEvent[event] = handler
}

func (c *Client) writeMessage(message []byte) error {
	err := c.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c *Client) HandlerServerResponse() {
	for {
		_, message, err := c.conn.ReadMessage()
		if websocket.IsUnexpectedCloseError(err) {
			log.Println("Websocket is disconnected for unknow reason, trying to reconnect")
			c.ConnectionWebSocket()
		}
		c.handlerServerMessage(string(message))
	}
}

func (c *Client) ConnectionWebSocket() {
	var err error
	for {
		c.conn, _, err = websocket.DefaultDialer.Dial(c.url, nil)
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	connectionMessage := newSocketIOMessage(CONNECT, nil, nil)
	connectionMessageByte, _ := connectionMessage.messageToMapOfByte()
	c.writeMessage(connectionMessageByte)
}

func (c *Client) handlerServerMessage(message string) {
	var socketIOMessageType string
	if len(message) > 1 {
		socketIOMessageType = message[1:2]
	}
	switch socketIOMessageType {
	case CONNECT:
		if len(message) > 1 {
			log.Printf("Connection with server, ID is: %s", message[2:])
		}
	case EVENT:
		trimMessage := c.trimMessage(message[2:])
		value, _ := c.OnEvent[trimMessage[0]]
		if value != nil {
			value(&trimMessage[1])
		}
	}
}

func (c *Client) trimMessage(message string) []string {
	var result []string
	err := json.Unmarshal([]byte(message), &result)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return result
}
