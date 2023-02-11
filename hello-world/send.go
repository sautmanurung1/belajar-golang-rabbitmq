package main

import (
	"fmt"
	"os"

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
	
	err = ch.Publish("", q.Name, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte(os.Args[1]),
	})

	if err != nil {
		panic(errors.Wrap(err, "Failed to publish message: "))
	}

	fmt.Println("Send Message", os.Args[1])
}
