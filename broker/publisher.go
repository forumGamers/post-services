package broker

import (
	"context"

	h "github.com/post-services/helper"
	m "github.com/post-services/models"
	"github.com/rabbitmq/amqp091-go"
)

const (
	POSTQUEUE = "Post-Queue"
)

type Publisher interface {
	SendPost(ctx context.Context, data m.Post) error
}

type publisherImpl struct {
	Channel *amqp091.Channel
}

func NewPublisher(ch *amqp091.Channel) Publisher {
	return &publisherImpl{
		Channel: ch,
	}
}

func (p *publisherImpl) QueueDeclare(name string) error {
	if _, err := p.Channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}
	return nil
}

func (p *publisherImpl) SendPost(ctx context.Context, data m.Post) error {
	if err := p.QueueDeclare(POSTQUEUE); err != nil {
		return err
	}

	if err := p.Channel.PublishWithContext(
		ctx,
		"",
		POSTQUEUE,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(h.ParseToJson(data)),
		},
	); err != nil {
		return err
	}

	return nil
}
