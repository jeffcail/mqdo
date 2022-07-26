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

	err = ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	failOnError2(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	failOnError2(err, "Failed to declare a queue")

	err = ch.QueueBind(q.Name, "", "logs", false, nil)
	failOnError2(err, "Failed to bind a queue")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError2(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [x] Waiting for logs. To exit press CTRL+C")
	<-forever
}
