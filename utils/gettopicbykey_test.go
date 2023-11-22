package utils

import (
	"errors"
	"testing"

	"github.com/0xanonymeow/go-subtest"
	"github.com/0xanonymeow/kafka-go/producer"
)

func TestGetTopicByKey(t *testing.T) {
	topics := []producer.Topic{
		{
			Key:  "example",
			Name: "producer",
		},
	}

	subtests := []subtest.Subtest{
		{
			Name:         "get_topic",
			ExpectedData: "producer",
			ExpectedErr:  nil,
			Test: func() (interface{}, error) {
				return GetTopicByKey(topics, "example")
			},
		},
		{
			Name:         "get_topic_not_found",
			ExpectedData: "",
			ExpectedErr:  errors.New("topic not found"),
			Test: func() (interface{}, error) {
				return GetTopicByKey(topics, "non-exist")
			},
		},
	}

	subtest.RunSubtests(t, subtests)
}
