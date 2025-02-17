package socketIOClient

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

func Connection(Ipaddr string, port string) *Client {
	factoClient := NewFactoryClient(Ipaddr, port)
	client := factoClient.NewClient()
	client.ConnectionWebSocket()
	return client
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
		OnEvent: make(map[string]func(EventHandlerInt)),
		handler: &EventHandler{
			conn: wsConn,
		},
	}
}

type Client struct {
	handler *EventHandler
	url     string
	conn    *websocket.Conn
	OnEvent map[string]func(EventHandlerInt)
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Emit(event string, message *string) error {
	err := c.handler.Emit(event, message)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) On(event string, handler func(EventHandlerInt)) {
	c.OnEvent[event] = handler
}

func (c *Client) writeMessage(message []byte) error {
	err := c.handler.writeMessage(message)
	if err != nil {
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
		c.handler.conn = c.conn
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	connectionMessage := newSocketIOMessage(CONNECT, nil, nil)
	connectionMessageByte, _ := connectionMessage.messageToMapOfByte()
	err = c.writeMessage(connectionMessageByte)
	if err != nil {
		log.Println(err)
	}
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
		value := c.OnEvent[trimMessage.EventName]
		if value != nil {
			c.handler.payload = trimMessage.Data
			value(*c.handler)
		}
	}
}

type TrimMessage struct {
	EventName string
	Data      string
}

func (c *Client) trimMessage(message string) *TrimMessage {
	var rawData []interface{}
	err := json.Unmarshal([]byte(message), &rawData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}

	var trimMessage TrimMessage
	if command, ok := rawData[0].(string); ok {
		trimMessage.EventName = command
	} else {
		fmt.Println("Error: Expected a string for the command")
		return nil
	}
	if rawData[1] != nil {
		trimMessage.Data = fmt.Sprint(rawData[1])
	}

	return &trimMessage
}
