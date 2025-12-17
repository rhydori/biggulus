package protocol

import (
	"strings"

	"github.com/rhydori/logs"
)

type Message struct {
	Entity string
	Action string
	Params []string
}

const paramIndex = 2

func ParseMessage(msg string) *Message {
	parts := strings.Split(msg, "|")
	if len(parts) < paramIndex {
		logs.Warnf("parseMessage: malformed message - %v", msg)
		return nil
	}
	message := &Message{
		Entity: parts[0],
		Action: parts[1],
	}
	if len(parts) > paramIndex {
		message.Params = parts[paramIndex:]
	}
	return message
}
func CreateMessageBytes(entity, action string, parts ...string) []byte {
	allParts := append([]string{entity, action}, parts...)
	msgString := strings.Join(allParts, "|")
	return []byte(msgString)
}
