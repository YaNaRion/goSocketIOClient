package socketIOClient

import (
	"github.com/gorilla/websocket"
	"log"
)

type EventHandlerInt interface {
	Payload() string
	Emit(string, *string) error
}

type EventHandler struct {
	payload string
	conn    *websocket.Conn
}

func (c EventHandler) Emit(event string, message *string) error {
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

func (c *EventHandler) writeMessage(message []byte) error {
	err := c.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c EventHandler) Payload() string {
	return c.payload
}
