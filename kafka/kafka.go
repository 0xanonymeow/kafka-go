package kafka

import (
	"github.com/0xanonymeow/kafka-go"
	"github.com/0xanonymeow/kafka-go/config"

	"context"
	"crypto/tls"
	"time"

	logger "github.com/0xanonymeow/kafka-go/logger"
	"github.com/pkg/errors"
	_kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type Kafka struct {
	writer       *_kafka.Writer
	readerConfig *_kafka.ReaderConfig
	reader       *_kafka.Reader
}

func NewKafka(_c interface{}) (kafka.Kafka, error) {
	c, err := config.LoadConfig(&_c)

	if err != nil {
		return nil, err
	}

	logger := logger.GetLogger(c)

	var saslMechanism sasl.Mechanism
	var tlsConfig *tls.Config

	if c.Secure.Enable {
		if c.Secure.ApiKey != "" && c.Secure.ApiSecret != "" {
			mech, err := scram.Mechanism(
				scram.SHA512,
				c.Secure.ApiKey,
				c.Secure.ApiSecret,
			)

			if err != nil {
				logger.Fatal(errors.Wrap(err, "failed to create sasl mechanism"))
			}

			saslMechanism = mech
		}

		if c.Secure.CertFile != "" && c.Secure.KeyFile != "" {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
				Certificates: []tls.Certificate{
					{
						Certificate: [][]byte{
							[]byte(c.Secure.CertFile),
						},
						PrivateKey: []byte(c.Secure.KeyFile),
					},
				},
				MinVersion: tls.VersionTLS12,
			}
		}
	}

	transport := &_kafka.Transport{
		SASL: saslMechanism,
		TLS:  tlsConfig,
	}

	d := &_kafka.Dialer{
		SASLMechanism: saslMechanism,
		Timeout:       3 * time.Minute,
		TLS:           tlsConfig,
		DualStack:     true,
	}

	writer := &_kafka.Writer{
		Addr:                   _kafka.TCP(c.Kafka.Connection),
		Balancer:               &_kafka.LeastBytes{},
		BatchSize:              1,
		BatchBytes:             1024,
		AllowAutoTopicCreation: c.Kafka.AllowAutoTopicCreation || true,
		Transport:              transport,
		RequiredAcks:           1,
	}

	rc := &_kafka.ReaderConfig{
		Brokers:  []string{c.Kafka.Connection},
		GroupID:  c.Kafka.Group,
		MaxWait:  10 * time.Second,
		MinBytes: 1,
		MaxBytes: 1024,
		Dialer:   d,
	}

	return &Kafka{
		writer:       writer,
		readerConfig: rc,
		reader:       nil,
	}, nil
}

func (k *Kafka) DialContext(ctx context.Context, network, addr string) (*_kafka.Conn, error) {
	return _kafka.DialContext(ctx, network, addr)
}

func (k *Kafka) NewReader(t []string) *_kafka.Reader {
	k.readerConfig.GroupTopics = t
	k.reader = _kafka.NewReader(*k.readerConfig)

	return k.reader
}

func (k *Kafka) WriteMessages(ctx context.Context, m _kafka.Message) error {
	return k.writer.WriteMessages(ctx, m)
}

func (k *Kafka) ReadMessage(ctx context.Context) (_kafka.Message, error) {
	if k.reader == nil {
		return _kafka.Message{}, errors.New("reader not initialized")
	}

	return k.reader.ReadMessage(ctx)
}

func (k *Kafka) FetchMessage(ctx context.Context) (_kafka.Message, error) {
	if k.reader == nil {
		return _kafka.Message{}, errors.New("reader not initialized")
	}

	return k.reader.FetchMessage(ctx)
}

func (k *Kafka) CommitMessages(ctx context.Context, m _kafka.Message) error {
	if k.reader == nil {
		return errors.New("reader not initialized")
	}

	return k.reader.CommitMessages(ctx, m)
}

func (k *Kafka) Close() error {
	if k.reader == nil {
		return errors.New("reader not initialized")
	}

	return k.reader.Close()
}
