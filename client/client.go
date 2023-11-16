package client

import (
	"context"
	"reflect"
	"sync"
	"time"

	pkg "github.com/0xanonymeow/kafka-go"
	"github.com/0xanonymeow/kafka-go/config"
	"github.com/0xanonymeow/kafka-go/consumer"
	"github.com/0xanonymeow/kafka-go/message"
	"github.com/0xanonymeow/kafka-go/producer"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"

	"github.com/sirupsen/logrus"
)

type Client struct {
	kafka pkg.Kafka
}

type Handler func(interface{}, []byte) error

var handlerMapper map[string]Handler
var ProducerTopics []producer.Topic
var ConsumerTopics []consumer.Topic

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

func registerHandler(ct consumer.Topic) error {
	v := reflect.ValueOf(ct)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if v.Type().Field(i).Name == "Handler" && field.IsNil() {
			return errors.Errorf("handler '%s' not found", ct.Name)
		}
	}

	method, exists := reflect.TypeOf(ct.Handler).MethodByName(ct.Name)

	if !exists {
		return errors.Errorf("method '%s' not found", ct.Name)
	}

	handlerMapper[ct.Topic] = func(_ interface{}, b []byte) error {
		reflectArgs := []reflect.Value{reflect.ValueOf(ct.Handler), reflect.ValueOf(b)}
		result := method.Func.Call([]reflect.Value(reflectArgs))

		if err, ok := result[0].Interface().(error); ok && err != nil {
			return err
		}

		return nil
	}

	return nil
}

func _init(c *config.Config) error {
	logger = logrus.StandardLogger()
	errch = make(chan error, 10)

	waitingLoopMu.Lock()
	defer waitingLoopMu.Unlock()
	waitingLoop = true

	retryLoopMu.Lock()
	defer retryLoopMu.Unlock()
	retryLoop = true

	retryInterval = 3 * time.Second

	ProducerTopics = c.Kafka.ProducerTopics
	ConsumerTopics = c.Kafka.ConsumerTopics

	handlerMapper = make(map[string]Handler)

	for _, ct := range ConsumerTopics {
		err := registerHandler(ct)

		if err != nil {
			select {
			case errch <- err:
			default:
			}

			logger.Error(err)
			continue
		}

		logger.Debugf("consumer handler '%s' added", ct.Name)
	}

	return nil
}

func NewClient(k pkg.Kafka, _c *config.Config) pkg.Client {
	ctx := context.Background()

	_init(_c)

	go func() {
		for retryLoop {
			_, err := k.DialContext(ctx, "tcp", _c.Kafka.Connection)

			if err != nil {
				err = errors.Wrap(err, "failed to connect")

				select {
				case errch <- err:
				default:
				}

				logger.Error(err)
				logger.Infof("retrying in %v ...", retryInterval)

				time.After(retryInterval)
				continue
			}

			retryLoopMu.Lock()
			defer retryLoopMu.Unlock()
			retryLoop = false
		}
	}()

	c := &Client{
		kafka: k,
	}

	for _, ct := range ConsumerTopics {
		go func(ct consumer.Topic) {
			logger.Debugf("waiting for messages on %v", ct.Topic)

			k.NewReader(ct.Topic)
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

				logger.Debugf("consumed event from topic %s", ct.Topic)

				handler, found := handlerMapper[m.Topic]

				if !found {
					err = errors.New("handler function is not found")

					select {
					case errch <- err:
					default:
					}

					logger.Error(err)
					continue
				}

				err = handler(ct.Handler, m.Value)

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
		}(ct)
	}

	return c
}

func (c *Client) Produce(m message.Message) error {
	logger := logrus.StandardLogger()

	produceMu.Lock()
	defer produceMu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.kafka.WriteMessages(ctx, kafka.Message(m))
	err = errors.Wrap(err, "failed to produce message")

	if err != nil {
		return err
	}

	logger.Debugf("successfully produced record to topic %s", m.Topic)

	return nil
}

func (c *Client) RegisterConsumerHandler(ct consumer.Topic) error {
	return registerHandler(ct)
}
