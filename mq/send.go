package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go_log/mq/defs"
	"log"
)

type SliceMock struct {
	addr uintptr
	len int
	cap int
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	failOnError(err, "Failed to declare a queue")

	//body := "Hello World!"
	testStruct := defs.Data{"127.0.0.1","xiaolin"}
	body, err := json.Marshal(testStruct)
	fmt.Println(body)

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)

	failOnError(err, "Failed to publish a message")
}
