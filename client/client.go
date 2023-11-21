package client

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/0xanonymeow/kafka-go"
	"github.com/0xanonymeow/kafka-go/config"
	"github.com/0xanonymeow/kafka-go/consumer"
	log "github.com/0xanonymeow/kafka-go/logger"
	"github.com/0xanonymeow/kafka-go/message"
	"github.com/0xanonymeow/kafka-go/producer"
	"github.com/0xanonymeow/kafka-go/utils"
	"github.com/pkg/errors"
	_kafka "github.com/segmentio/kafka-go"

	"github.com/sirupsen/logrus"
)

type Client struct {
	kafka   kafka.Kafka
	mapping map[string]Handler
}

type Handler func(interface{}, []byte) error

var handlerMapper map[string]Handler
var Producer producer.Producer
var Consumers []consumer.Consumer

var produceMu sync.Mutex
var (
	waitingLoop   bool
	waitingLoopMu sync.Mutex
)
var (
	retryLoop   bool
	retryLoopMu sync.Mutex
)
var retryInterval time.Duration
var errch chan error
var logger *logrus.Logger

func registerHandler(_type interface{}, t string) (string, error) {
	m := utils.TopicToMethodName(t)

	method, exists := reflect.TypeOf(_type).MethodByName(m)

	if !exists {
		return "", errors.Errorf("method '%s' not found", m)
	}

	handlerMapper[t] = func(_ interface{}, b []byte) error {
		reflectArgs := []reflect.Value{reflect.ValueOf(_type), reflect.ValueOf(b)}
		result := method.Func.Call([]reflect.Value(reflectArgs))

		if err, ok := result[0].Interface().(error); ok && err != nil {
			return err
		}

		return nil
	}

	return m, nil
}

func _init(c *config.Config) error {
	logger = log.GetLogger(c)
	errch = make(chan error, 10)

	waitingLoopMu.Lock()
	defer waitingLoopMu.Unlock()
	waitingLoop = true

	retryLoopMu.Lock()
	defer retryLoopMu.Unlock()
	retryLoop = true

	retryInterval = 3 * time.Second

	Producer = c.Kafka.Producer
	Consumers = c.Kafka.Consumers

	handlerMapper = make(map[string]Handler)

	for _, csmr := range Consumers {
		for _, ct := range csmr.Topics {
			if csmr.Type == nil {
				return errors.Errorf("handler of topic %v cannot be nil", ct.Name)
			}

			name, err := registerHandler(csmr.Type, ct.Name)

			if err != nil {
				return err
			}

			logger.Debugf("consumer handler '%s' added", name)
		}
	}

	return nil
}

func NewClient(k kafka.Kafka, _c interface{}) (kafka.Client, error) {
	c, err := config.LoadConfig(&_c)

	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	err = _init(c)

	if err != nil {
		return nil, err
	}

	go func() {
		for retryLoop {
			_, err := k.DialContext(ctx, "tcp", c.Kafka.Connection)

			if err != nil {
				err = errors.Wrap(err, "failed to connect")

				select {
				case errch <- err:
				default:
				}

				logger.Error(err)
				logger.Infof("retrying in %v ...", retryInterval)

				time.Sleep(retryInterval)
				continue
			}

			retryLoopMu.Lock()
			defer retryLoopMu.Unlock()
			retryLoop = false
		}

		for _, csmr := range Consumers {
			for _, ct := range csmr.Topics {
				go func(_type interface{}, t string) {
					logger.Debugf("waiting for messages on %v", t)

					k.NewReader(t)
					defer k.Close()

					for waitingLoop {
						m, err := k.FetchMessage(ctx)

						if err != nil {
							err = errors.Wrap(err, "failed to read topic")

							select {
							case errch <- err:
							default:
							}

							logger.Error(err)
							continue
						}

						logger.Debugf("consumed event from topic %s", t)

						handler := handlerMapper[m.Topic]

						err = handler(_type, m.Value)

						if err != nil {
							err = errors.Wrap(err, "failed to handle topic")

							select {
							case errch <- err:
							default:
							}

							logger.Error(err)
							continue
						}

						err = k.CommitMessages(ctx, m)

						if err != nil {
							err = errors.Wrap(err, "failed to commit message")
							select {
							case errch <- err:
							default:
							}

							logger.Error(err)
						}
					}
				}(csmr.Type, ct.Name)
			}
		}
	}()

	return &Client{
		kafka: k,
	}, nil
}

func (c *Client) Produce(m message.Message) error {
	produceMu.Lock()
	defer produceMu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.kafka.WriteMessages(ctx, _kafka.Message(m))

	if err != nil {
		return errors.Wrap(err, "failed to produce message")
	}

	logger.Debugf("successfully produced record to topic %s", m.Topic)

	return nil
}
