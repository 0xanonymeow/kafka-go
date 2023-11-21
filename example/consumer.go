package example

import (
	"fmt"

	pkg "github.com/0xanonymeow/kafka-go"
	client "github.com/0xanonymeow/kafka-go/client"
	"github.com/0xanonymeow/kafka-go/config"
	"github.com/0xanonymeow/kafka-go/kafka"
)

type Consumer struct {
	kafka                pkg.Client
	handlerSpecificProps map[string]interface{}
}

func ConsumerExample() error {
	_c, err := config.LoadConfig()

	if err != nil {
		return err
	}

	k, err := kafka.NewKafka(_c)

	if err != nil {
		return nil
	}

	exampleProp := "example"
	csmr := Consumer{
		handlerSpecificProps: map[string]interface{}{
			"example": exampleProp,
		},
	}
	_c.Kafka.Consumers[0].Type = &csmr
	_, err = client.NewClient(k, _c)

	if err != nil {
		return err
	}

	return nil
}

func (c *Consumer) ExampleHandler(b []byte) error {
	prop := c.handlerSpecificProps["example"]

	fmt.Printf("example prop: %v", prop)

	return nil
}
