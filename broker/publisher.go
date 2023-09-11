package broker

import (
	"context"
	"fmt"
	"os"
	"time"

	h "github.com/post-services/helper"
	"github.com/rabbitmq/amqp091-go"
)

const (
	POSTEXCHANGE    = "Post-Exchange"
	NEWPOSTQUEUE    = "New-Post-Queue"
	DELETEPOSTQUEUE = "Delete-Post-Queue"
)

type PublisherImpl struct {
	Channel *amqp091.Channel
}

type Publisher interface {
	PublishMessage(ctx context.Context, exchangeName, queueName, ContentType string, data any) error
}

var Broker Publisher

func (p *PublisherImpl) PublishMessage(
	ctx context.Context,
	exchangeName,
	queueName,
	ContentType string,
	data any,
) error {
	return p.Channel.PublishWithContext(
		ctx,
		exchangeName,
		fmt.Sprintf("%s.%s", exchangeName, NEWPOSTQUEUE),
		false,
		false,
		amqp091.Publishing{
			ContentType:  ContentType,
			Body:         []byte(h.ParseToJson(data)),
			DeliveryMode: 2,
			Timestamp:    time.Now(),
		},
	)
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

	for _, exchangeName := range []string{POSTEXCHANGE} {
		h.PanicIfError(
			ch.ExchangeDeclare(
				exchangeName,
				"direct",
				true,
				false,
				false,
				false,
				nil,
			),
		)
	}

	for _, queueName := range []string{NEWPOSTQUEUE, DELETEPOSTQUEUE} {
		_, err := ch.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			nil,
		)
		h.PanicIfError(err)
	}

	for _, queueName := range []string{NEWPOSTQUEUE, DELETEPOSTQUEUE} {
		for _, exchangeName := range []string{POSTEXCHANGE} {
			h.PanicIfError(
				ch.QueueBind(
					queueName,
					fmt.Sprintf("%s.%s", exchangeName, queueName),
					exchangeName,
					false,
					nil,
				),
			)
		}
	}

	notifyClose := conn.NotifyClose(make(chan *amqp091.Error))
	go func() { <-notifyClose }()

	Broker = &PublisherImpl{
		Channel: ch,
	}
	fmt.Println("connection to broker success")
}

type Media struct {
	Url  string `json:"url"`
	Type string `json:"type"`
	Id   string `json:"id"`
}

type PostDocument struct {
	Id           string `json:"id"`
	UserId       string `json:"userId"`
	Text         string `json:"text" bson:"text"`
	Media        Media
	AllowComment bool `json:"allowComment" default:"true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Tags         []string `json:"tags"`
	Privacy      string   `json:"privacy" default:"Public"`
}