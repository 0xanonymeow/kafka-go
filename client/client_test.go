package client

import (
	"io"
	"sync"
	"testing"

	pkg "github.com/0xanonymeow/kafka-go"
	"github.com/0xanonymeow/kafka-go/config"
	"github.com/0xanonymeow/kafka-go/consumer"
	"github.com/0xanonymeow/kafka-go/message"
	"github.com/segmentio/kafka-go"

	"github.com/0xanonymeow/kafka-go/testings"
	mock_kafka "github.com/0xanonymeow/kafka-go/testings/mock_kafka"

	"github.com/pkg/errors"
	"go.uber.org/mock/gomock"
)

const topicTest = "test-topic"
const topicTestHandlerNotExist = "test-topic-handler-not-exist"
const topicTestHandlerError = "test-topic-handler-error"
const TestHandlerName string = "TopicHandler"
const TestHandlerNameNotExist string = "TopicHandlerNotExist"
const TestHandlerNameError string = "TopicHandlerError"

func (*Client) TopicHandler([]byte) error {
	return nil
}

func (*Client) TopicHandlerError([]byte) error {
	return errors.New("error expected")
}

type TestHandler struct {
	kafka pkg.Client
}

func (h *TestHandler) Handler([]byte) error {
	return nil
}

func TestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	k := mock_kafka.NewMockKafka(ctrl)
	_c := &config.Config{
		Kafka: config.Kafka{
			Connection: "localhost:9092",
		},
	}

	subtests := []testings.Subtest{
		{
			Name: "new_client",
			ExpectedData: &Client{
				kafka: k,
			},
			ExpectedErr: nil,
			Test: func() (interface{}, error) {

				wg := sync.WaitGroup{}
				wg.Add(1)

				k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, nil).
					Do(func(arg0, arg1, arg2 interface{}) {
						defer wg.Done()
					})
				c := NewClient(k, _c)

				wg.Wait()

				return c, nil
			},
		},
		{
			Name:         "new_client_error",
			ExpectedData: nil,
			ExpectedErr:  errors.Errorf("failed to connect: %v", "failed to open connection"),
			Test: func() (interface{}, error) {
				err := errors.New("failed to open connection")

				wg := sync.WaitGroup{}
				wg.Add(1)

				k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, err).
					Do(func(arg0, arg1, arg2 interface{}) {
						defer wg.Done()

						retryLoopMu.Lock()
						defer retryLoopMu.Unlock()
						retryLoop = false
					})
				c := NewClient(k, _c)

				wg.Wait()

				receivedErr := <-errch

				return c, receivedErr
			},
		},
		{
			Name: "consume_message",
			ExpectedData: &Client{
				kafka: k,
			},
			ExpectedErr: nil,
			Test: func() (interface{}, error) {
				wg := sync.WaitGroup{}
				wg.Add(3)

				k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, nil).
					Do(func(arg0, arg1, arg2 interface{}) {
						defer wg.Done()
					})
				k.EXPECT().NewReader(gomock.Any()).Return(&kafka.Reader{})
				k.EXPECT().Close().Return(nil)
				t := topicTest
				b := []byte("test")
				m := kafka.Message{
					Topic: t,
					Value: b,
				}
				k.EXPECT().FetchMessage(gomock.Any()).Return(m, nil).
					Do(func(arg0 interface{}) {
						defer wg.Done()
					})
				k.EXPECT().CommitMessages(gomock.Any(), m).Return(nil).
					Do(func(arg0, arg1 interface{}) {
						defer wg.Done()

						waitingLoopMu.Lock()
						defer waitingLoopMu.Unlock()
						waitingLoop = false
					})
				_c := &config.Config{
					Kafka: config.Kafka{
						Connection: "localhost:9092",
						ConsumerTopics: []consumer.Topic{
							{
								Topic:   t,
								Name:    TestHandlerName,
								Handler: &Client{},
							},
						},
					},
				}

				c := NewClient(k, _c)

				wg.Wait()

				return c, nil
			},
		},
		{
			Name:         "consume_message_fetch_error",
			ExpectedData: nil,
			ExpectedErr:  errors.Errorf("failed to read topic: %v", io.ErrUnexpectedEOF),
			Test: func() (interface{}, error) {
				err := io.ErrUnexpectedEOF

				wg := sync.WaitGroup{}
				wg.Add(2)

				k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, nil).
					Do(func(arg0, arg1, arg2 interface{}) {
						defer wg.Done()
					})
				k.EXPECT().NewReader(gomock.Any()).Return(&kafka.Reader{})
				k.EXPECT().Close().Return(nil)
				k.EXPECT().FetchMessage(gomock.Any()).Return(kafka.Message{}, err).
					Do(func(arg0 interface{}) {
						defer wg.Done()

						waitingLoopMu.Lock()
						defer waitingLoopMu.Unlock()
						waitingLoop = false

					})
				_c := &config.Config{
					Kafka: config.Kafka{
						ConsumerTopics: []consumer.Topic{
							{
								Topic:   topicTest,
								Name:    TestHandlerName,
								Handler: &Client{},
							},
						},
					},
				}
				c := NewClient(k, _c)

				wg.Wait()

				receivedErr := <-errch

				return c, receivedErr
			},
		},
		{
			Name:         "message_handler_not_found",
			ExpectedData: nil,
			ExpectedErr:  errors.Errorf("method '%v' not found", TestHandlerNameNotExist),
			Test: func() (interface{}, error) {
				wg := sync.WaitGroup{}
				wg.Add(2)

				k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, nil).
					Do(func(arg0, arg1, arg2 interface{}) {
						defer wg.Done()
					})
				k.EXPECT().NewReader(gomock.Any()).Return(&kafka.Reader{})
				k.EXPECT().Close().Return(nil)
				t := topicTestHandlerNotExist
				b := []byte("test")
				m := kafka.Message{
					Topic: t,
					Value: b,
				}
				k.EXPECT().FetchMessage(gomock.Any()).Return(m, nil).
					Do(func(arg0 interface{}) {
						defer wg.Done()

						waitingLoopMu.Lock()
						defer waitingLoopMu.Unlock()
						waitingLoop = false
					})
				_c := &config.Config{
					Kafka: config.Kafka{
						Connection: "localhost:9092",
						ConsumerTopics: []consumer.Topic{
							{
								Topic:   topicTest,
								Name:    TestHandlerNameNotExist,
								Handler: &Client{},
							},
						},
					},
				}
				c := NewClient(k, _c)

				wg.Wait()

				receivedErr := <-errch

				return c, receivedErr
			},
		},
		{
			Name:         "message_handler_error",
			ExpectedData: nil,
			ExpectedErr:  errors.Errorf("failed to handle topic: %v", "error expected"),
			Test: func() (interface{}, error) {
				wg := sync.WaitGroup{}
				wg.Add(2)

				k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, nil).
					Do(func(arg0, arg1, arg2 interface{}) {
						defer wg.Done()
					})
				k.EXPECT().NewReader(gomock.Any()).Return(&kafka.Reader{})
				k.EXPECT().Close().Return(nil)
				t := topicTestHandlerError
				b := []byte("test")
				m := kafka.Message{
					Topic: t,
					Value: b,
				}
				k.EXPECT().FetchMessage(gomock.Any()).Return(m, nil).
					Do(func(arg0 interface{}) {
						defer wg.Done()

						waitingLoopMu.Lock()
						defer waitingLoopMu.Unlock()
						waitingLoop = false
					})
				_c := &config.Config{
					Kafka: config.Kafka{
						ConsumerTopics: []consumer.Topic{
							{
								Topic:   t,
								Name:    TestHandlerNameError,
								Handler: &Client{},
							},
						},
					},
				}
				c := NewClient(k, _c)

				wg.Wait()

				receivedErr := <-errch

				return c, receivedErr
			},
		},
		{
			Name:         "consume_message_commit_error",
			ExpectedData: nil,
			ExpectedErr:  errors.Errorf("failed to commit message: %v", io.ErrClosedPipe),
			Test: func() (interface{}, error) {
				err := io.ErrClosedPipe

				wg := sync.WaitGroup{}
				wg.Add(3)

				k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, nil).
					Do(func(arg0, arg1, arg2 interface{}) {
						defer wg.Done()
					})
				k.EXPECT().NewReader(gomock.Any()).Return(&kafka.Reader{})
				k.EXPECT().Close().Return(nil)
				t := topicTest
				b := []byte("test")
				m := kafka.Message{
					Topic: t,
					Value: b,
				}
				k.EXPECT().FetchMessage(gomock.Any()).Return(m, nil).
					Do(func(arg0 interface{}) {
						defer wg.Done()
					})
				k.EXPECT().CommitMessages(gomock.Any(), m).Return(err).
					Do(func(arg0, arg1 interface{}) {
						defer wg.Done()
						waitingLoopMu.Lock()
						defer waitingLoopMu.Unlock()
						waitingLoop = false
					})
				_c := &config.Config{
					Kafka: config.Kafka{
						ConsumerTopics: []consumer.Topic{
							{
								Topic:   t,
								Name:    TestHandlerName,
								Handler: &Client{},
							},
						},
					},
				}
				c := NewClient(k, _c)

				wg.Wait()

				receivedErr := <-errch

				return c, receivedErr
			},
		},
		{
			Name:         "produce_message",
			ExpectedData: nil,
			ExpectedErr:  nil,
			Test: func() (interface{}, error) {
				wg := sync.WaitGroup{}
				wg.Add(1)

				k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, nil).
					Do(func(arg0, arg1, arg2 interface{}) {
						defer wg.Done()
					})
				c := NewClient(k, _c)

				wg.Wait()

				t := topicTest
				b := []byte("test")
				m := kafka.Message{
					Topic: t,
					Value: b,
				}
				k.EXPECT().WriteMessages(gomock.Any(), m).Return(nil)
				err := c.Produce(message.Message(m))

				return nil, err
			},
		},
		{
			Name:         "produce_message_error",
			ExpectedData: nil,
			ExpectedErr:  errors.Errorf("failed to produce message: %v", "unknown server error"),
			Test: func() (interface{}, error) {
				err := errors.New("unknown server error")

				wg := sync.WaitGroup{}
				wg.Add(1)

				k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, nil).
					Do(func(arg0, arg1, arg2 interface{}) {
						defer wg.Done()
					})
				c := NewClient(k, _c)

				wg.Wait()

				b := []byte("test")
				m := message.Message{
					Topic: topicTest,
					Value: b,
				}
				k.EXPECT().WriteMessages(gomock.Any(), gomock.Any()).Return(err)
				err = c.Produce(m)

				return nil, err
			},
		},
		{
			Name:         "register_consumer_handler",
			ExpectedData: nil,
			ExpectedErr:  nil,
			Test: func() (interface{}, error) {
				wg := sync.WaitGroup{}
				wg.Add(1)

				k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, nil).
					Do(func(arg0, arg1, arg2 interface{}) {
						defer wg.Done()
					})
				c := NewClient(k, _c)
				err := c.RegisterConsumerHandler(consumer.Topic{
					Topic:   topicTest,
					Name:    TestHandlerName,
					Handler: c,
				})

				wg.Wait()

				return nil, err
			},
		},
		{
			Name:         "register_consumer_handler_not_found",
			ExpectedData: nil,
			ExpectedErr:  errors.Errorf("handler '%s' not found", TestHandlerNameNotExist),
			Test: func() (interface{}, error) {
				wg := sync.WaitGroup{}
				wg.Add(1)

				k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, nil).
					Do(func(arg0, arg1, arg2 interface{}) {
						defer wg.Done()
					})
				c := NewClient(k, _c)
				err := c.RegisterConsumerHandler(consumer.Topic{
					Topic: topicTestHandlerNotExist,
					Name:  TestHandlerNameNotExist,
				})

				wg.Wait()

				return nil, err
			},
		},
	}

	testings.RunSubtests(t, subtests)
}
