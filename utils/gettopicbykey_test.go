package utils

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/0xanonymeow/go-subtest"
)

func TestGetTopicByKey(t *testing.T) {

	subtests := []subtest.Subtest{
		{
			Name:         "get_topic",
			ExpectedData: "producer",
			ExpectedErr:  nil,
			Test: func() (interface{}, error) {
				return GetTopicByKey("example")
			},
			Setup: func() {
				currentDir, _ := os.Getwd()
				currentDir = filepath.Dir(currentDir)
				path := filepath.Join(currentDir, ".config.toml.example")

				os.Setenv("KAFKA_GO_CONFIG_PATH", path)
			},
		},
		{
			Name:         "get_topic_not_found",
			ExpectedData: "",
			ExpectedErr:  errors.New("topic not found"),
			Test: func() (interface{}, error) {
				return GetTopicByKey("non-exist")
			},
			Setup: func() {
				currentDir, _ := os.Getwd()
				currentDir = filepath.Dir(currentDir)
				path := filepath.Join(currentDir, ".config.toml.example")

				os.Setenv("KAFKA_GO_CONFIG_PATH", path)
			},
		},
		{
			Name:         "load_config_error",
			ExpectedData: "",
			ExpectedErr:  errors.New("failed to load config"),
			Test: func() (interface{}, error) {
				return GetTopicByKey("example")
			},
		},
	}

	subtest.RunSubtests(t, subtests)
}
