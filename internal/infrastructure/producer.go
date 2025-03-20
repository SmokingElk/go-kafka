package infrastructure

import (
	"errors"
	"time"

	"github.com/IBM/sarama"
)

var ErrTimeout = errors.New("interrupt by timeout")

type Producer struct {
	producer sarama.AsyncProducer
}

func NewProducer(servers []string) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForLocal
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true

	p, err := sarama.NewAsyncProducer(servers, cfg)

	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: p,
	}, nil
}

func (p *Producer) Produce(key, value, topic string) error {
	m := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	p.producer.Input() <- m

	select {
	case <-p.producer.Successes():
		return nil
	case err := <-p.producer.Errors():
		return err
	case <-time.After(time.Duration(100) * time.Second):
		return ErrTimeout
	}
}

func (p *Producer) Close() {
	p.producer.Close()
}
