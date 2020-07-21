package main

import (
	"log"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/streadway/amqp"
)

func main() {
	url := os.Getenv("AMQP_URL")

	if url == "" {
		url = "amqp://guest:guest@localhost:5672"
	}

	connection, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}

	channel, err := connection.Channel()
	if err != nil {
		panic("could not open channel " + err.Error())
	}

	err = channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	message := amqp.Publishing{
		Body:      []byte("YOOOOOO"),
		Timestamp: time.Now(),
	}

	id, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("uuid %s", err)
	}
	err = channel.Publish("events", id.String(), false, false, message)
	if err != nil {
		panic("cannot publish message" + err.Error())
	}

	_, err = channel.QueueDeclare("test", true, false, false, false, nil)
	if err != nil {
		panic("error declaring queue " + err.Error())
	}

	err = channel.QueueBind("test", "#", "events", false, nil)
	if err != nil {
		panic("error bind queue " + err.Error())
	}

}
