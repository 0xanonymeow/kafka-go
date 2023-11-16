package config

import (
	"errors"
	"os"
	"testing"

	"github.com/0xanonymeow/kafka-go/testings"
	"github.com/go-playground/validator"
)

func TestConfig(t *testing.T) {
	subtests := []testings.Subtest{
		{
			Name:         "load_config",
			ExpectedData: nil,
			ExpectedErr:  nil,
			Test: func() (interface{}, error) {
				config, err := LoadConfig()

				validate := validator.New()
				err = validate.Struct(config)

				return nil, err
			},
		},
		{
			Name:         "load_config_not_found",
			ExpectedData: nil,
			ExpectedErr:  errors.New("open invalid: no such file or directory"),
			Test: func() (interface{}, error) {
				os.Setenv("KAFKA_GO_CONFIG_PATH", "invalid")

				_, err := LoadConfig()

				return nil, err
			},
		},
		{
			Name:         "load_config_error",
			ExpectedData: nil,
			ExpectedErr:  errors.New("open false: no such file or directory"),
			Test: func() (interface{}, error) {
				os.Setenv("KAFKA_GO_CONFIG_PATH", "false")

				_, err := LoadConfig()

				return nil, err
			},
		},
		{
			Name:         "validate_error",
			ExpectedData: nil,
			ExpectedErr:  errors.New("validator: (nil *config.Config)"),
			Test: func() (interface{}, error) {
				toml := `
					[invalid]
					config = false
				`

				os.Setenv("KAFKA_GO_CONFIG_PATH", toml)

				config, err := LoadConfig()

				validate := validator.New()
				err = validate.Struct(config)

				return nil, err
			},
		},
	}

	testings.RunSubtests(t, subtests)
}
