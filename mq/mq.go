package mq

import (
	"log"
	"errors"
	// "fmt"
	// "math/rand"

	"github.com/streadway/amqp"
)

// Channel returns amqp channel
type Channel struct {
	Ch	*amqp.Channel
	Conn	*amqp.Connection
	Queue	amqp.Queue
	ExcName	string
}

func failOnError(err error, msg string) error {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		return errors.New(msg)
	}
	return nil
}

// TokenGenerator generates token for routing key
// func TokenGenerator() string {
// 	b := make([]byte, 18)
// 	rand.Read(b)
// 	return fmt.Sprintf("%x", b)
// }

// InitMQ initialize RabbitMQ connection
func InitMQ(url string, vhost string) (*Channel, error) {
	conn, err := amqp.Dial(url + vhost)
	if err1 := failOnError(err, "Failed to connect to RabbitMQ"); err1 != nil {
		return nil, err1
	}
	// defer conn.Close()

	ch, err := conn.Channel()
	if err1 := failOnError(err, "Failed to open a channel"); err1 != nil {
		return nil, err1
	}
	// defer ch.Close()

	out := new(Channel)
	out.Ch = ch
	out.Conn = conn

	return out, nil
}

// ExcDeclare declares exchange
func (ch *Channel) ExcDeclare(excName string, excType string) error {
	err =: ch.Ch.ExchangeDeclare(excName, excType, false, false, false, false, nil)
	if err1 := failOnError(err, "Failed to declare exchange"); err1 != nil {
		return err1
	}
	ch.ExcName = excName
	return nil
}

// QueueDeclare declares queue
func (ch *Channel) QueueDeclare(queueName string) error {
	q, err := ch.Ch.QueueDeclare(
		queueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err1 := failOnError(err, "Failed to declare queue"); err1 != nil {
		return err1
	}
	ch.Queue = q
	return nil
}

// PostMessage posts message to RabbitMQ
func (ch *Channel) PostMessage(msg string, rKey string) error {
	err := ch.Ch.Publish(ch.ExcName, rKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
	if err1 := failOnError(err, "Publish Error"); err1 != nil {
		return err1
	}
	return nil
}
