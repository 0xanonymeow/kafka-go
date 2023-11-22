package utils

import (
	"errors"

	"github.com/0xanonymeow/kafka-go/producer"
)

func GetTopicByKey(topics []producer.Topic, k string) (string, error) {
	for _, t := range topics {
		if t.Key == k {
			return t.Name, nil
		}
	}

	return "", errors.New("topic not found")
}
