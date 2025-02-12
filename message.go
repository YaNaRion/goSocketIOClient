package socketIOClient

import (
	"encoding/json"
	"fmt"
)

type SocketIOMessage struct {
	WSMessage string
	EventType string
	EventName string
	data      *string
}

func formatEventName(s *string) string {
	if s != nil {
		return fmt.Sprintf("\"%s\"", *s)
	}
	return "/"
}

func newSocketIOMessage(EventType string, EventName *string, data *string) *SocketIOMessage {
	return &SocketIOMessage{
		WSMessage: MessageWS,
		EventType: EventType,
		EventName: formatEventName(EventName),
		data:      data,
	}
}

func (m *SocketIOMessage) messageToMapOfByte() ([]byte, error) {
	var bytesMessage []byte
	bytesMessage = append(bytesMessage, []byte(m.WSMessage)...)
	bytesMessage = append(bytesMessage, []byte(m.EventType)...)
	jsonConnect, err := json.Marshal(m.data)
	if err != nil {
		return nil, err
	}
	payloadString := fmt.Sprintf("[%s, %s]", m.EventName, string(jsonConnect))
	bytesMessage = append(bytesMessage, []byte(payloadString)...)
	return bytesMessage, nil
}
