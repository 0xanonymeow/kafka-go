package utils

import (
	"errors"

	"github.com/0xanonymeow/kafka-go/config"
)

func GetTopicByKey(k string) (string, error) {
	c, err := config.LoadConfig()

	if err != nil {
		return "", errors.New("failed to load config")
	}

	topics := c.Kafka.Producer.Topics

	for _, t := range topics {
		if t.Key == k {
			return t.Name, nil
		}
	}

	return "", errors.New("topic not found")
}
