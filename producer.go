package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	*RabbitMQ
	channel *amqp.Channel
	log     *log.Logger
}

type Options struct {
	RoutingKey   string
	ContentType  string
	DeliveryMode uint8
	Message      string
	Mandatory    bool
	Immediate    bool
}

func (r *RabbitMQ) NewProducer() (*Producer, error) {
	rabbitmq, err := r.Connect()
	if err != nil {
		return nil, err
	}

	channel, err := r.conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Producer{
		RabbitMQ: rabbitmq,
		channel:  channel,
		log:      log.Default(),
	}, nil
}

func (p *Producer) Publish(queue Queue, exchange Exchange, options Options) error {
	return p.publish(queue, exchange, options)
}

func (r *RabbitMQ) publish(queue Queue, exchange Exchange, options Options) error {
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

	err = channel.Publish(
		exchange.Name,
		options.RoutingKey,
		options.Mandatory,
		options.Immediate,
		amqp.Publishing{
			DeliveryMode: options.DeliveryMode,
			ContentType:  options.ContentType,
			Body:         []byte(options.Message),
		})

	if err != nil {
		return err
	}

	log.Printf("Send new message: %s", options.Message)
	return nil
}
