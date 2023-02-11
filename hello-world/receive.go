package main

import (
	"log"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(errors.Wrap(err, "Failed To Connect to Rabbitmq: "))
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(errors.Wrap(err, "Failed to get Channel"))
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("Hello", false, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "Failed to get the queue : "))
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Recieve a message : %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting fo the Message. To Exit Press Ctrl + C\n")
	<-forever
}
