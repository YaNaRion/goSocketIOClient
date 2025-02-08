package main

import (
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
		conn: wsConn,
		url:  url,
	}
}

type Client struct {
	url  string
	conn *websocket.Conn
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) WriteMessage(message []byte) {
	err := c.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println(err)
	}
}

func (c *Client) HandlerServerResponse() {
	for {
		_, message, err := c.conn.ReadMessage()
		if websocket.IsUnexpectedCloseError(err) {
			log.Println("Websocket is disconnected for unknow reason, trying to reconnect")
			break
		}
		log.Printf("Received message from the server: %s\n", message)
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
	log.Println(string(connectionMessageByte))
	c.WriteMessage(connectionMessageByte)
}
