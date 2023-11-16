package config

import (
	"os"

	"github.com/0xanonymeow/kafka-go/consumer"
	"github.com/0xanonymeow/kafka-go/producer"
	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

type Env struct {
	Development bool
}

type Kafka struct {
	Connection             string                   `toml:"conn"`
	ApiKey                 string                   `toml:"api_key"`
	ApiSecret              string                   `toml:"api_secret"`
	ProducerTopics         []producer.Topic `toml:"producer_topics"`
	ConsumerTopics         []consumer.Topic `toml:"consumer_topics"`
	Group                  string                   `toml:"grp"`
	AllowAutoTopicCreation bool                     `toml:"auto_create_topic"`
}

type Config struct {
	Env Env `toml:"env"`
	Log struct {
		Level string `toml:"level"`
	} `toml:"log"`
	Server struct {
		CertFile string `toml:"cert_file"`
		KeyFile  string `toml:"key_file"`
	} `toml:"server"`
	Kafka Kafka `toml:"kafka"`
}

func LoadConfig() (*Config, error) {
	path := os.Getenv("KAFKA_GO_CONFIG_PATH")

	if path == "" {
		godotenv.Load("../.env.example")

		path = "../" + os.Getenv("KAFKA_GO_CONFIG_PATH")
	}

	c := Config{}

	if _, err := toml.DecodeFile(path, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
