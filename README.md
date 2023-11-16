# kafka-go

This project is a simple kafka go client that wraps the [github.com/segmentio/kafka-go](https://github.com/segmentio/kafka-go) library. It aims to provide a simplified interface for producing and consuming messages using kafka go.

## Mapping Producer and Consumer Topics

To effectively produce and consume messages using kafka, it's important to correctly map producer and consumer topics.

Follow these steps to map topics:

1. Create a new kafka topic for your producer using the kafka command-line tool or an admin library like [kafka.apache.org/quickstart](https://kafka.apache.org/quickstart) or set `auto_create_topic` to `true` in `.confg.toml`.

> Check [github.com/0xanonymeow/kafka-go/.config.toml.example](https://github.com/0xanonymeow/kafka-go/blob/main/example/.config.toml.example)

2. Update your `KAFKA_GO_CONFIG_PATH` in your env to point to the `config.toml`
3. In your producer code, specify the key to get the topic name then send a message to kafka along with the value.

```go
k := kafka.NewKafka(config)
c := kafkaclient.NewClient(k, config)

t, err := utils.GetTopicByKey("example")

if err != nil {
	return err
}

m := message.Message{
	Topic: t,
	Value: []byte("test"),
}

err = c.Produce(m)
```

> Check [github.com/0xanonymeow/kafka-go/example/producer.go](https://github.com/0xanonymeow/kafka-go/blob/main/example/producer.go)

4. In your consumer code, make sure to create a struct method of your handler, beginning with an uppercase letter, and embed `client` into your struct. Always make certain to designate `[]byte` as your first parameter, as the reader loop will transmit the message through this parameter.

```go
type Consumer struct {
	kafka Client
}

func (h *Handler) Handler([]byte) error {}
```

- 4.1 After that, register the handler to `handlerMapper`

```go
c := client.NewClient(k, config)
csmr := Consumer{
	kafka: c,
}

c.RegisterConsumerHandler(consumer.Topic{
	Topic:   "consumer",
	Name:    "Handler",
	Handler: &csmr,
})
```

> Check [github.com/0xanonymeow/kafka-go/example/consumer.go](https://github.com/0xanonymeow/kafka-go/blob/main/example/consumer.go)

5. Ensure that the topic name is consistent between the producer and consumer code. This ensures that messages are properly sent and received.

<br/>
By following these steps, you can effectively map producer and consumer topics in your kafka application.

## Contributing

Contributions are welcome! If you have any suggestions, bug reports, or feature requests, please open an issue on the GitHub repository.

## License

This project is released under the [The Unlicense](https://choosealicense.com/licenses/unlicense/).

## Acknowledgements

This project is built on top of the [github.com/segmentio/kafka-go](https://github.com/segmentio/kafka-go) library. Special thanks to the authors and contributors of that library for their excellent work.
