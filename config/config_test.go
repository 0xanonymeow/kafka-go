package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/0xanonymeow/go-subtest"
	"github.com/go-playground/validator"
)

func TestConfig(t *testing.T) {
	subtests := []subtest.Subtest{
		{
			Name:         "load_config_default",
			ExpectedData: nil,
			ExpectedErr:  nil,
			Test: func() (interface{}, error) {
				config, err := LoadConfig()

				if err != nil {
					return nil, err
				}

				validate := validator.New()
				err = validate.Struct(config)

				return nil, err
			},
			Setup: func() {
				currentDir, _ := os.Getwd()
				currentDir = filepath.Dir(currentDir)
				path := filepath.Join(currentDir, ".config.toml.example")

				os.Setenv("KAFKA_GO_CONFIG_PATH", path)
			},
		},
		{
			Name:         "load_config_with_single_param",
			ExpectedData: nil,
			ExpectedErr:  nil,
			Test: func() (interface{}, error) {
				c := &Config{
					Kafka: Kafka{
						Connection: "localhost:9092",
					},
				}
				var param interface{}
				param = c

				config, err := LoadConfig(&param)

				if err != nil {
					return nil, err
				}

				validate := validator.New()
				err = validate.Struct(config)

				return nil, err
			},
		},
		{
			Name:         "load_config_with_single_param_invalid",
			ExpectedData: nil,
			ExpectedErr:  errors.New("invalid config conversion"),
			Test: func() (interface{}, error) {
				type wrongStruct struct{}
				w := &wrongStruct{}
				var param interface{}
				param = w

				config, err := LoadConfig(&param)

				if err != nil {
					return nil, err
				}

				validate := validator.New()
				err = validate.Struct(config)

				return nil, err
			},
		},
		{
			Name:         "load_config_with_single_param_error",
			ExpectedData: nil,
			ExpectedErr:  errors.New("require struct"),
			Test: func() (interface{}, error) {
				type wrongKind int
				var w wrongKind
				var param interface{}
				param = w

				config, err := LoadConfig(&param)

				if err != nil {
					return nil, err
				}

				validate := validator.New()
				err = validate.Struct(config)

				return nil, err
			},
		},
		{
			Name:         "load_config_too_many_params_error",
			ExpectedData: nil,
			ExpectedErr:  errors.New("too many arguments"),
			Test: func() (interface{}, error) {
				type someStruct struct{}
				s := &someStruct{}
				var param interface{}
				param = s

				config, err := LoadConfig(&param, &param)

				if err != nil {
					return nil, err
				}

				validate := validator.New()
				err = validate.Struct(config)

				return nil, err
			},
		},
		{
			Name:         "load_config_error",
			ExpectedData: nil,
			ExpectedErr:  errors.New("open invalid: no such file or directory"),
			Test: func() (interface{}, error) {
				return LoadConfig()
			},
			Setup: func() {
				os.Setenv("KAFKA_GO_CONFIG_PATH", "invalid")
			},
		},
		{
			Name:         "load_config_empty",
			ExpectedData: nil,
			ExpectedErr:  errors.New("kafka config path cannot be empty"),
			Test: func() (interface{}, error) {
				return LoadConfig()
			},
		},
		{
			Name:         "validate_error",
			ExpectedData: nil,
			ExpectedErr:  errors.New("kafka conn cannot be empty"),
			Test: func() (interface{}, error) {
				config, err := LoadConfig()

				if err != nil {
					return nil, err
				}

				validate := validator.New()
				err = validate.Struct(config)

				return nil, err
			},
			Setup: func() {
				toml := `
					[invalid]
					config = false
				`

				path := ".config.toml.test"

				file, _ := os.Create(path)

				defer file.Close()

				content := []byte(toml)

				file.Write(content)

				os.Setenv("KAFKA_GO_CONFIG_PATH", path)

			},
			Teardown: func() {
				path := ".config.toml.test"

				os.Remove(path)
			},
		},
	}

	subtest.RunSubtests(t, subtests)
}
