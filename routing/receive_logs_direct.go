package main

import (
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

func failOnError2(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError2(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError2(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare("logs_direct", "direct",
		true, false, false, false, nil)
	failOnError2(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	failOnError2(err, "Failed to declare a queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(1)
	}

	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "logs_direct", s)
		err = ch.QueueBind(q.Name, s, "logs_direct", false, nil)
		failOnError2(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError2(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
