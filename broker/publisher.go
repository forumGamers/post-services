package broker

import "github.com/rabbitmq/amqp091-go"

type Publisher interface {
	Send(data any)
}

type PublisherImpl struct {
	Channel *amqp091.Channel
}

func NewPublisher(ch *amqp091.Channel) Publisher {
	return &PublisherImpl{
		Channel: ch,
	}
}

func (p *PublisherImpl) Send(data any) {
	
}
