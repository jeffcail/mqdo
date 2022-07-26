package main

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func failOnError2(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError2(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError2(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	failOnError2(err, "Failed to declare a queue")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError2(err, "Failed to register a consumer")

	var forever chan struct{}
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] waiting for message. To exit press CTRL+C")
	<-forever
}
