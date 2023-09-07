package broker

import (
	"context"
	"fmt"
	"os"

	h "github.com/post-services/helper"
	"github.com/rabbitmq/amqp091-go"
)

const (
	NEWPOSTQUEUE    = "New-Post-Queue"
	DELETEPOSTQUEUE = "Delete-Post-Queue"
)

type PublisherImpl struct {
	Channel *amqp091.Channel
}

type Publisher interface {
	PublishMessage(ctx context.Context, queueName, ContentType string, data any) error
}

var Broker Publisher

func (p *PublisherImpl) PublishMessage(
	ctx context.Context,
	queueName,
	ContentType string,
	data any,
) error {

	q, err := p.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	if err := p.Channel.PublishWithContext(
		ctx,
		q.Name,
		NEWPOSTQUEUE,
		false,
		false,
		amqp091.Publishing{
			ContentType: ContentType,
			Body:        []byte(h.ParseToJson(data)),
		},
	); err != nil {
		return err
	}
	return nil
}

func BrokerConnection() {
	rabbitMqServerUrl := os.Getenv("RABBITMQURL")

	if rabbitMqServerUrl == "" {
		rabbitMqServerUrl = "amqp://user:password@localhost:5672"
	}

	conn, err := amqp091.Dial(rabbitMqServerUrl)
	h.PanicIfError(err)

	ch, err := conn.Channel()
	h.PanicIfError(err)

	Broker = &PublisherImpl{
		Channel: ch,
	}
	fmt.Println("connection to broker success")
}
