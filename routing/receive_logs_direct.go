package main

import (
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

func failsOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", err, msg)
	}
}

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672")
	failsOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failsOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	failsOnError(err, "Failed to declare a queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failsOnError(err, "Failed to declare a queue")

	if len(os.Args) < 1 {
		log.Printf("Usage : %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	}
	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "log_direct", s)
		err = ch.QueueBind(
			q.Name,
			s,
			"logs_direct",
			false,
			nil,
		)
		failsOnError(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failsOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [*] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To Exit press CTRL + C\n")
	<-forever
}
