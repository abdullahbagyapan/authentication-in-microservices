package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"microservices/auth-service/internal/ports"
)

type Adapter struct {
	ch *amqp.Channel
}

func NewMsgBrokerAdapter() *Adapter {
	conn, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("Error connecting to rabbitmq %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating to channel %v", err)
	}

	_, err = declareQueue(ch)
	if err != nil {
		log.Fatalf("Error declaring to queue %v", err)
	}

	return &Adapter{ch: ch}
}

func connectRabbitMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	return conn, err
}

func (adapter Adapter) PublishMessage(messageStruct *ports.MsgBrokerUserInfo) error {

	byteMsg, err := json.Marshal(messageStruct)

	if err != nil {
		log.Printf("Error converting to byte  %v", err)
	}

	err = adapter.ch.Publish(
		"",            // exchange
		"authservice", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        byteMsg,
		})

	return err
}

func declareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		"authservice", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)

	return q, err
}
