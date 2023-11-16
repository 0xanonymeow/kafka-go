package message

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	Topic         string
	Partition     int
	Offset        int64
	HighWaterMark int64
	Key           []byte
	Value         []byte
	Headers       []kafka.Header
	WriterData    interface{}
	Time          time.Time
}
