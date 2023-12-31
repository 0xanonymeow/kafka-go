package config

import (
	"errors"
	"os"
	"reflect"

	"github.com/0xanonymeow/kafka-go/consumer"
	"github.com/0xanonymeow/kafka-go/producer"
	"github.com/BurntSushi/toml"
	"github.com/mitchellh/mapstructure"
)

type Secure struct {
	Enable    bool   `toml:"enable"`
	ApiKey    string `toml:"api_key"`
	ApiSecret string `toml:"api_secret"`
	CertFile  string `toml:"cert_file"`
	KeyFile   string `toml:"key_file"`
}

type Kafka struct {
	Connection             string              `toml:"conn"`
	Producer               producer.Producer   `toml:"producer"`
	Consumers              []consumer.Consumer `toml:"consumers"`
	Group                  string              `toml:"grp"`
	AllowAutoTopicCreation bool                `toml:"auto_create_topic"`
}

type Config struct {
	Log struct {
		Level string `toml:"level"`
	} `toml:"log"`
	Secure Secure `toml:"secure"`
	Kafka  Kafka  `toml:"kafka"`
}

func LoadConfig(args ...*interface{}) (*Config, error) {
	c := Config{}

	if len(args) == 1 {
		v := reflect.ValueOf(*args[0])

		if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
			return nil, errors.New("require struct")
		}

		config := &mapstructure.DecoderConfig{
			Result: &c,
		}

		decoder, _ := mapstructure.NewDecoder(config)
		decoder.Decode(*args[0])
	} else if len(args) > 1 {
		return nil, errors.New("too many arguments")
	} else {
		path := os.Getenv("KAFKA_GO_CONFIG_PATH")

		if path == "" {
			return nil, errors.New("kafka config path cannot be empty")
		}

		if _, err := toml.DecodeFile(path, &c); err != nil {
			return nil, err
		}
	}

	if c.Kafka.Connection == "" {
		return nil, errors.New("kafka conn cannot be empty")
	}

	return &c, nil
}
