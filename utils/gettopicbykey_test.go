package utils

import (
	"errors"
	"os"
	"testing"

	"github.com/0xanonymeow/kafka-go/testings"
)

func TestGetTopicByKey(t *testing.T) {

	subtests := []testings.Subtest{
		{
			Name:         "get_topic",
			ExpectedData: "producer-topic",
			ExpectedErr:  nil,
			Test: func() (interface{}, error) {
				return GetTopicByKey("example")
			},
		},
		{
			Name:         "get_topic_not_found",
			ExpectedData: "",
			ExpectedErr:  errors.New("topic not found"),
			Test: func() (interface{}, error) {
				return GetTopicByKey("non-exist")
			},
		},
		{
			Name:         "load_config_error",
			ExpectedData: "",
			ExpectedErr:  errors.New("failed to load config"),
			Test: func() (interface{}, error) {
				os.Setenv("KAFKA_GO_CONFIG_PATH", "invalid")

				return GetTopicByKey("example")
			},
		},
	}

	testings.RunSubtests(t, subtests)
}
