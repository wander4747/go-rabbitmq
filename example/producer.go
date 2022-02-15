package main

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/wander4747/go-rabbitmq"
)

type Message struct {
	Text string
}

func producer() {

	config := rabbitmq.Config{
		User:     "guest",
		Password: "guest",
		Host:     "localhost",
		Port:     5672,
		Vhost:    "",
	}

	queue := rabbitmq.Queue{
		Name:       "new_code_queue",
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}

	exchange := rabbitmq.Exchange{
		Name:       "new_code_exchange",
		Type:       "direct",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}

	rabbit := rabbitmq.NewRabbitMQ(&config)

	producer, _ := rabbit.NewProducer()

	for i := 0; i < 5; i++ {
		u, err := json.Marshal(Message{Text: fmt.Sprintf("text %v", i)})
		if err != nil {
			log.Println(err)
		}

		options := rabbitmq.Options{
			RoutingKey:   "",
			ContentType:  "application/json",
			Message:      string(u),
			DeliveryMode: amqp.Persistent,
		}

		err = producer.Publish(queue, exchange, options)
		if err != nil {
			log.Println(err)
		}
	}
}
