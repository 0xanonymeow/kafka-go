# kafka-go

This project is a simple kafka go client that wraps the [github.com/segmentio/kafka-go](https://github.com/segmentio/kafka-go) library. It aims to provide a simplified interface for producing and consuming messages using kafka go.

## Mapping Producer and Consumer Topics

To effectively produce and consume messages using kafka, it's important to correctly map producer and consumer topics.

Follow these steps to map topics:

1. Create a new kafka topic for your producer using the kafka command-line tool or an admin library like [kafka.apache.org/quickstart](https://kafka.apache.org/quickstart) or set `auto_create_topic` to `true` in `.confg.toml`.

> Check [github.com/0xanonymeow/kafka-go/.config.toml.example](https://github.com/0xanonymeow/kafka-go/blob/main/.config.toml.example)

1. Update your `KAFKA_GO_CONFIG_PATH` in your env to point to the `config.toml`
2. In your producer code, specify the key to get the topic name then send a message to kafka along with the value.

```go
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
```

> Check [github.com/0xanonymeow/kafka-go/example/producer.go](https://github.com/0xanonymeow/kafka-go/blob/main/example/producer.go)

1. Ensure that, in your consumer code, you create a struct method for your handler, commencing with an uppercase letter, and embed client into your struct. Additional attributes can be added as needed. Always adhere to the function signature, as the reader loop will pass the message through this parameter. Your handler method name should also conform to a specific naming pattern; it must begin with an uppercase letter and end with `Handler`. Ensure that your method name matches the topic name. For example, if your topic is `example`, your handler name should be `ExampleHandler`. If the topic contains a hyphen `-`, remove the hyphen and concatenate the words. If it contains a dot `.`, only consider the latter part.
   
2. 
```go
type Consumer struct {
	kafka Client
	handlerSpecificProps map[string]interface{}
}

func (h *Handler) ExampleHandler([]byte) error {}
```

> Check [github.com/0xanonymeow/kafka-go/example/consumer.go](https://github.com/0xanonymeow/kafka-go/blob/main/example/consumer.go)

1. Ensure that the topic name is consistent between the producer and consumer code. This ensures that messages are properly sent and received.

<br/>
By following these steps, you can effectively map producer and consumer topics in your kafka application.

## Contributing

Contributions are welcome! If you have any suggestions, bug reports, or feature requests, please open an issue on the GitHub repository.

## License

This project is released under the [The Unlicense](https://choosealicense.com/licenses/unlicense/).

## Acknowledgements

This project is built on top of the [github.com/segmentio/kafka-go](https://github.com/segmentio/kafka-go) library. Special thanks to the authors and contributors of that library for their excellent work.
