package example

import (
	client "github.com/0xanonymeow/kafka-go/client"
	"github.com/0xanonymeow/kafka-go/config"
	"github.com/0xanonymeow/kafka-go/kafka"
	"github.com/0xanonymeow/kafka-go/message"
	"github.com/0xanonymeow/kafka-go/utils"
)

func producer() error {
	_c, err := config.LoadConfig()

	if err != nil {
		return err
	}

	k := kafka.NewKafka(_c)
	c := client.NewClient(k, _c)

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
