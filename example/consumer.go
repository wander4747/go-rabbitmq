package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/wander4747/go-rabbitmq"
)

func consumer() {

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

	options := rabbitmq.ConsumerOptions{
		RoutingKey: "",
		Consumer:   "",
		AutoAck:    false,
		Exclusive:  false,
		NoLocal:    false,
		NoWait:     false,
		Args:       nil,
	}
	rabbit := rabbitmq.NewRabbitMQ(&config)

	consumer, _ := rabbit.NewConsumer()

	messageChannel := make(chan amqp.Delivery)

	err := consumer.Consumer(messageChannel, queue, exchange, options)
	if err != nil {
		log.Println(err)
	}

	for msg := range messageChannel {
		log.Printf(" [x] %s", msg.Body)
		if err := msg.Ack(false); err != nil {
			log.Printf("Error acknowledging message : %s", err)
		} else {
			log.Printf("Acknowledged message")
		}
	}
}
