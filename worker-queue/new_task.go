package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("worker-queue-test", false, false, false, false, nil)
	failOnError(err, "Failed to declare q queue")

	body := bodyFrom(os.Args)
	err = ch.PublishWithContext(context.Background(), "", q.Name, false, false, amqp091.Publishing{
		ContentType:  "test/plain",
		DeliveryMode: amqp091.Persistent,
		Body:         []byte(body),
	})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}
