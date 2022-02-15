package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	*RabbitMQ
	channel *amqp.Channel
	log     *log.Logger
}

type ConsumerOptions struct {
	RoutingKey string
	Consumer   string
	AutoAck    bool
	Exclusive  bool
	NoLocal    bool
	NoWait     bool
	Args       amqp.Table
}

func (r *RabbitMQ) NewConsumer() (*Consumer, error) {
	rabbitmq, err := r.Connect()
	if err != nil {
		return nil, err
	}

	channel, err := r.conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		RabbitMQ: rabbitmq,
		channel:  channel,
		log:      log.Default(),
	}, nil
}

func (c *Consumer) Consumer(messageChannel chan amqp.Delivery, queue Queue, exchange Exchange, options ConsumerOptions) error {
	return c.consumer(messageChannel, queue, exchange, options)
}

func (r *RabbitMQ) consumer(messageChannel chan amqp.Delivery, queue Queue, exchange Exchange, options ConsumerOptions) error {

	ch, err := r.Channel()
	if err != nil {
		return err
	}

	_, err = r.queueDeclare(queue, ch)
	if err != nil {
		return err
	}

	channel, err := r.exchangeDeclare(exchange, ch)
	if err != nil {
		return err
	}

	err = channel.QueueBind(
		queue.Name,
		options.RoutingKey,
		exchange.Name,
		queue.NoWait,
		queue.Args,
	)

	if err != nil {
		return err
	}

	msgs, err := channel.Consume(
		queue.Name,
		options.Consumer,
		options.AutoAck,
		options.Exclusive,
		options.NoLocal,
		options.NoWait,
		options.Args,
	)

	if err != nil {
		return err
	}

	go func() {
		for m := range msgs {
			log.Println("Incoming new message")
			messageChannel <- m
		}
		log.Println("RabbitMQ channel closed")
		close(messageChannel)
	}()

	return nil
}
