package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Kafka interface {
	DialContext(context.Context, string, string) (*kafka.Conn, error)
	NewReader([]string) *kafka.Reader
	WriteMessages(context.Context, kafka.Message) error
	ReadMessage(context.Context) (kafka.Message, error)
	FetchMessage(context.Context) (kafka.Message, error)
	CommitMessages(context.Context, kafka.Message) error
	Close() error
}
