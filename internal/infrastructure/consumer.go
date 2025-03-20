package infrastructure

import (
	"context"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type handler struct {
	handleFunc func(string, string)
}

func (h *handler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *handler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			h.handleFunc(string(message.Key), string(message.Value))
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

type Consumer struct {
	group sarama.ConsumerGroup
	stop  context.CancelFunc
}

func NewConsumer(servers []string, group, topic string, handleFunc func(key, value string)) (*Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.AutoCommit.Enable = true
	cfg.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	g, err := sarama.NewConsumerGroup(servers, group, cfg)

	if err != nil {
		return nil, err
	}

	h := &handler{
		handleFunc: handleFunc,
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := g.Consume(ctx, []string{topic}, h); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	return &Consumer{
		group: g,
		stop:  cancel,
	}, nil
}

func (c *Consumer) Close() {
	c.stop()
	c.group.Close()
}
