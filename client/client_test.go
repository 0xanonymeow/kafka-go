package client

import (
	"io"
	"sync"
	"testing"

	pkg "github.com/0xanonymeow/kafka-go"
	"github.com/0xanonymeow/kafka-go/config"
	"github.com/0xanonymeow/kafka-go/consumer"
	"github.com/0xanonymeow/kafka-go/message"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"

	"github.com/0xanonymeow/go-subtest"
	mock_kafka "github.com/0xanonymeow/kafka-go/testings/mock_kafka"

	"go.uber.org/mock/gomock"
)

const topicTest = "test-topic"
const topicTestMethodNotExist = "test-topic-method-not-exist"
const topicTestHandlerNotExist = "test-topic-handler-not-exist"
const topicTestMethodError = "test-topic-method-error"

type TestHandler struct {
	kafka pkg.Client
}

func (h *TestHandler) TestTopicHandler([]byte) error {
	return nil
}

func (h *TestHandler) TestTopicMethodErrorHandler([]byte) error {
	return errors.New("error expected")
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
	h := TestHandler{}

	subtests := []subtest.Subtest{
		{
			Name: "new_client",
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
				csmrs := []consumer.Consumer{
					{
						Topics: []consumer.Topic{
							{
								Name: t,
							},
						},
						Type: &h,
					},
				}
				_c := config.Config{
					Kafka: config.Kafka{
						Connection: "localhost:9092",
						Consumers:  csmrs,
					},
				}
				c, err := NewClient(k, &_c)

				wg.Wait()

				return c, err
			},
		},
		{
			Name: "new_client_empty_handler",
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
				c, err := NewClient(k, _c)

				wg.Wait()

				return c, err
			},
		},
		{
			Name:         "new_client_method_not_found",
			ExpectedData: nil,
			ExpectedErr:  errors.Errorf("method '%s' not found", "TestTopicMethodNotExistHandler"),
			Test: func() (interface{}, error) {
				csmrs := []consumer.Consumer{
					{
						Topics: []consumer.Topic{
							{
								Name: topicTestMethodNotExist,
							},
						},
						Type: &h,
					},
				}
				_c := config.Config{
					Kafka: config.Kafka{
						Connection: "localhost:9092",
						Consumers:  csmrs,
					},
				}
				c, err := NewClient(k, &_c)

				return c, err
			},
		},
		{
			Name:         "new_client_handler_not_found",
			ExpectedData: nil,
			ExpectedErr:  errors.Errorf("handler of topic %v cannot be nil", topicTestHandlerNotExist),
			Test: func() (interface{}, error) {
				csmrs := []consumer.Consumer{
					{
						Topics: []consumer.Topic{
							{
								Name: topicTestHandlerNotExist,
							},
						},
					},
				}
				_c := config.Config{
					Kafka: config.Kafka{
						Connection: "localhost:9092",
						Consumers:  csmrs,
					},
				}
				c, err := NewClient(k, &_c)

				return c, err
			},
		},
		{
			Name:         "new_client_connection_error",
			ExpectedData: nil,
			ExpectedErr:  errors.Errorf("failed to connect: %v", "failed to open connection"),
			Test: func() (interface{}, error) {
				wg := sync.WaitGroup{}
				wg.Add(1)

				err := errors.New("failed to open connection")
				k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, err).
					Do(func(arg0, arg1, arg2 interface{}) {
						defer wg.Done()

						retryLoopMu.Lock()
						defer retryLoopMu.Unlock()
						retryLoop = false
					})
				c, err := NewClient(k, _c)

				if err != nil {
					return nil, err
				}

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
				csmrs := []consumer.Consumer{
					{
						Topics: []consumer.Topic{
							{
								Name: topicTest,
							},
						},
						Type: &h,
					},
				}
				_c := config.Config{
					Kafka: config.Kafka{
						Connection: "localhost:9092",
						Consumers:  csmrs,
					},
				}
				c, err := NewClient(k, &_c)

				if err != nil {
					return nil, err
				}

				wg.Wait()

				return c, nil
			},
		},
		{
			Name:         "consume_message_fetch_error",
			ExpectedData: nil,
			ExpectedErr:  errors.Errorf("failed to read topic: %v", io.ErrUnexpectedEOF),
			Test: func() (interface{}, error) {
				wg := sync.WaitGroup{}
				wg.Add(2)

				k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, nil).
					Do(func(arg0, arg1, arg2 interface{}) {
						defer wg.Done()
					})
				k.EXPECT().NewReader(gomock.Any()).Return(&kafka.Reader{})
				k.EXPECT().Close().Return(nil)
				err := io.ErrUnexpectedEOF
				k.EXPECT().FetchMessage(gomock.Any()).Return(kafka.Message{}, err).
					Do(func(arg0 interface{}) {
						defer wg.Done()

						waitingLoopMu.Lock()
						defer waitingLoopMu.Unlock()
						waitingLoop = false

					})
				csmrs := []consumer.Consumer{
					{
						Topics: []consumer.Topic{
							{
								Name: topicTest,
							},
						},
						Type: &h,
					},
				}
				_c := config.Config{
					Kafka: config.Kafka{
						Connection: "localhost:9092",
						Consumers:  csmrs,
					},
				}
				c, err := NewClient(k, &_c)

				if err != nil {
					return nil, err
				}

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
				t := topicTestMethodError
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
				csmrs := []consumer.Consumer{
					{
						Topics: []consumer.Topic{
							{
								Name: t,
							},
						},
						Type: &h,
					},
				}
				_c := config.Config{
					Kafka: config.Kafka{
						Connection: "localhost:9092",
						Consumers:  csmrs,
					},
				}
				c, err := NewClient(k, &_c)

				if err != nil {
					return nil, err
				}

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
				err := io.ErrClosedPipe
				k.EXPECT().CommitMessages(gomock.Any(), m).Return(err).
					Do(func(arg0, arg1 interface{}) {
						defer wg.Done()
						waitingLoopMu.Lock()
						defer waitingLoopMu.Unlock()
						waitingLoop = false
					})
				csmrs := []consumer.Consumer{
					{
						Topics: []consumer.Topic{
							{
								Name: t,
							},
						},
						Type: &h,
					},
				}
				_c := config.Config{
					Kafka: config.Kafka{
						Connection: "localhost:9092",
						Consumers:  csmrs,
					},
				}
				c, err := NewClient(k, &_c)

				if err != nil {
					return nil, err
				}

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
				c, err := NewClient(k, _c)

				if err != nil {
					return nil, err
				}

				wg.Wait()

				t := topicTest
				b := []byte("test")
				m := kafka.Message{
					Topic: t,
					Value: b,
				}
				k.EXPECT().WriteMessages(gomock.Any(), m).Return(nil)
				err = c.Produce(message.Message(m))

				return nil, err
			},
		},
		{
			Name:         "produce_message_error",
			ExpectedData: nil,
			ExpectedErr:  errors.Errorf("failed to produce message: %v", "unknown server error"),
			Test: func() (interface{}, error) {
				wg := sync.WaitGroup{}
				wg.Add(1)

				k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, nil).
					Do(func(arg0, arg1, arg2 interface{}) {
						defer wg.Done()
					})
				c, err := NewClient(k, _c)

				if err != nil {
					return nil, err
				}

				wg.Wait()

				b := []byte("test")
				m := message.Message{
					Topic: topicTest,
					Value: b,
				}
				err = errors.New("unknown server error")
				k.EXPECT().WriteMessages(gomock.Any(), gomock.Any()).Return(err)
				err = c.Produce(m)

				return nil, err
			},
		},
		// {
		// 	Name:         "register_consumer_handler",
		// 	ExpectedData: nil,
		// 	ExpectedErr:  nil,
		// 	Test: func() (interface{}, error) {
		// 		wg := sync.WaitGroup{}
		// 		wg.Add(1)

		// 		k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, nil).
		// 			Do(func(arg0, arg1, arg2 interface{}) {
		// 				defer wg.Done()
		// 			})
		// 		c := NewClient(k, _c)
		// 		err := c.RegisterConsumerHandlers([]consumer.Topic{
		// 			{
		// 				Name:   topicTest,
		// 				Type: c,
		// 			},
		// 		})

		// 		wg.Wait()

		// 		return nil, err
		// 	},
		// },
		// {
		// 	Name:         "register_consumer_handler_not_found",
		// 	ExpectedData: nil,
		// 	ExpectedErr:  errors.Errorf("handler '%s' not found", TestMethodNameNotExist),
		// 	Test: func() (interface{}, error) {
		// 		wg := sync.WaitGroup{}
		// 		wg.Add(1)

		// 		k.EXPECT().DialContext(gomock.Any(), gomock.Any(), gomock.Any()).Return(&kafka.Conn{}, nil).
		// 			Do(func(arg0, arg1, arg2 interface{}) {
		// 				defer wg.Done()
		// 			})
		// 		c := NewClient(k, _c)
		// 		err := c.RegisterConsumerHandlers([]consumer.Topic{
		// 			{
		// 				Name: topicTestMethodNotExist,
		// 			},
		// 		})

		// 		wg.Wait()

		// 		return nil, err
		// 	},
		// },
	}

	subtest.RunSubtests(t, subtests)
}
