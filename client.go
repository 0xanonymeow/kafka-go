package kafka

import (
	"github.com/0xanonymeow/kafka-go/message"
)

type Client interface {
	Produce(message.Message) error
}
