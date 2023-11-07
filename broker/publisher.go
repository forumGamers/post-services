package broker

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	h "github.com/post-services/helper"
	"github.com/rabbitmq/amqp091-go"
)

type PublisherImpl struct {
	Channel *amqp091.Channel
}

type Publisher interface {
	PublishMessage(ctx context.Context, exchangeName, queueName, ContentType string, data any) error
	DeclareExchangeAndQueue()
	BindQueue(exchange string, queues []string)
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
		fmt.Sprintf("%s.%s", exchangeName, queueName),
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

	conn, err := amqp091.DialConfig(rabbitMqServerUrl, amqp091.Config{
		Heartbeat: 10,
	})
	h.PanicIfError(err)

	ch, err := conn.Channel()
	h.PanicIfError(err)

	notifyClose := conn.NotifyClose(make(chan *amqp091.Error))
	go func() {
		retries := 0
		for {
			select {
			case err := <-notifyClose:
				if err != nil && retries < 10 {
					newConn, newErr := amqp091.DialConfig(rabbitMqServerUrl, amqp091.Config{
						Heartbeat: 10,
					})
					if newErr != nil {
						log.Printf("Gagal melakukan koneksi ulang: %s", newErr)
						continue
					}

					newCh, newErr := newConn.Channel()
					if newErr != nil {
						newConn.Close()
						log.Printf("Gagal membuat channel baru: %s", newErr)
						continue
					}

					Broker = &PublisherImpl{
						Channel: newCh,
					}
					notifyClose = conn.NotifyClose(make(chan *amqp091.Error))
				}
				break
			}
		}
	}()

	Broker = &PublisherImpl{
		Channel: ch,
	}

	Broker.DeclareExchangeAndQueue()

	log.Println("connection to broker success")
}

func (ch *PublisherImpl) DeclareExchangeAndQueue() {
	exchanges := []string{POSTEXCHANGE, LIKEEXCHANGE, COMMENTEXCHANGE, REPLYEXCHANGE, SHAREEXCHANGE}
	queues := []string{
		NEWPOSTQUEUE, DELETEPOSTQUEUE, BULKPOSTQUEUE,
		NEWLIKEQUEUE, DELETELIKEQUEUE, BULKLIKEQUEUE,
		NEWCOMMENTQUEUE, DELETECOMMENTQUEUE, BULKCOMMENTQUEUE,
		NEWREPLYQUEUE, DELETEREPLYQUEUE,
		NEWSHAREQUEUE, DELETESHAREQUEUE,
	}
	for _, exchangeName := range exchanges {
		h.PanicIfError(
			ch.Channel.ExchangeDeclare(
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

	for _, queueName := range queues {
		_, err := ch.Channel.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			nil,
		)
		h.PanicIfError(err)
	}

	for _, exchange := range exchanges {
		var queues []string
		switch exchange {
		case POSTEXCHANGE:
			queues = []string{NEWPOSTQUEUE, DELETEPOSTQUEUE, BULKPOSTQUEUE}
		case LIKEEXCHANGE:
			queues = []string{NEWLIKEQUEUE, DELETELIKEQUEUE, BULKLIKEQUEUE}
		case COMMENTEXCHANGE:
			queues = []string{NEWCOMMENTQUEUE, DELETECOMMENTQUEUE, BULKCOMMENTQUEUE}
		}
		ch.BindQueue(exchange, queues)
	}
}

func (ch *PublisherImpl) BindQueue(exchange string, queues []string) {
	for _, queue := range queues {
		h.PanicIfError(
			ch.Channel.QueueBind(
				queue,
				fmt.Sprintf("%s.%s", exchange, queue),
				exchange,
				false,
				nil,
			),
		)
	}
}
