package example

import (
	client "github.com/0xanonymeow/kafka-go/client"
	"github.com/0xanonymeow/kafka-go/config"
	"github.com/0xanonymeow/kafka-go/kafka"
	"github.com/0xanonymeow/kafka-go/message"
	"github.com/0xanonymeow/kafka-go/utils"
)

func ProducerExample() error {
	_c, err := config.LoadConfig()

	if err != nil {
		return err
	}

	k, err := kafka.NewKafka(_c)

	if err != nil {
		return err
	}

	c, err := client.NewClient(k, _c)

	if err != nil {
		return err
	}

	t, err := utils.GetTopicByKey("example")

	if err != nil {
		return err
	}

	m := message.Message{
		Topic: t,
		Value: []byte("test"),
	}

	err = c.Produce(m)

	return err
}
