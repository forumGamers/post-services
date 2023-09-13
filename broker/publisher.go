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
	log.Println("connection to broker success")
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
