package kafka

import (
	pkg "github.com/0xanonymeow/kafka-go"
	"github.com/0xanonymeow/kafka-go/config"

	"context"
	"crypto/tls"
	"time"

	"github.com/pkg/errors"
	_kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/sirupsen/logrus"
)

type Kafka struct {
	writer       *_kafka.Writer
	readerConfig *_kafka.ReaderConfig
	reader       *_kafka.Reader
}

func NewKafka(c *config.Config) pkg.Kafka {
	var saslMechanism sasl.Mechanism
	var tlsConfig *tls.Config

	if c.Env.Development {
		saslMechanism = nil
		tlsConfig = nil
	} else {
		mech, err := scram.Mechanism(
			scram.SHA512,
			c.Kafka.ApiKey,
			c.Kafka.ApiSecret,
		)

		if err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to create sasl mechanism"))
		}

		saslMechanism = mech

		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
			Certificates: []tls.Certificate{
				{
					Certificate: [][]byte{
						[]byte(c.Server.CertFile),
					},
					PrivateKey: []byte(c.Server.KeyFile),
				},
			},
			MinVersion: tls.VersionTLS12,
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
		MaxWait:  3 * time.Minute,
		MinBytes: 1,
		MaxBytes: 1024,
		Dialer:   d,
	}

	return &Kafka{
		writer:       writer,
		readerConfig: rc,
		reader:       nil,
	}
}

func (k *Kafka) DialContext(ctx context.Context, network, addr string) (*_kafka.Conn, error) {
	return _kafka.DialContext(ctx, network, addr)
}

func (k *Kafka) NewReader(t string) *_kafka.Reader {
	k.readerConfig.Topic = t
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
