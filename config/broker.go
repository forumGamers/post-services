package config

import (
	"fmt"
	"os"

	h "github.com/post-services/helper"
	"github.com/rabbitmq/amqp091-go"
)

var Broker *amqp091.Connection

func BrokerConnection() {
	rabbitMqServerUrl := os.Getenv("RABBITMQURL")

	if rabbitMqServerUrl == "" {
		rabbitMqServerUrl = "amqp://guest:guest@172.19.0.5:5672/"
	}

	conn, err := amqp091.Dial(rabbitMqServerUrl)
	h.PanicIfError(err)
	fmt.Println("connection to broker success")

	defer conn.Close()

	Broker = conn
}
