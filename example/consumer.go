package example

import (
	pkg "github.com/0xanonymeow/kafka-go"
	client "github.com/0xanonymeow/kafka-go/client"
	"github.com/0xanonymeow/kafka-go/config"
	_consumer "github.com/0xanonymeow/kafka-go/consumer"
	"github.com/0xanonymeow/kafka-go/kafka"
)

type Consumer struct {
	kafka pkg.Client
}

func consumer() error {
	_c, err := config.LoadConfig()

	if err != nil {
		return err
	}

	k := kafka.NewKafka(_c)
	c := client.NewClient(k, _c)
	csmr := Consumer{
		kafka: c,
	}

	c.RegisterConsumerHandler(_consumer.Topic{
		Topic:   "consumer",
		Name:    "Handler",
		Handler: &csmr,
	})

	return nil
}

func (c *Consumer) Handler(b []byte) error {

	return nil
}
