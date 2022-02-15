package rabbitmq

import (
	"errors"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Exchange struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	Vhost    string
}

type RabbitMQ struct {
	conn   *amqp.Connection
	config *Config
	log    *log.Logger
}

func NewRabbitMQ(config *Config) *RabbitMQ {
	return &RabbitMQ{
		conn:   nil,
		config: config,
		log:    log.Default(),
	}
}

func (r *RabbitMQ) Connect() (*RabbitMQ, error) {
	if r.config == nil {
		return nil, errors.New("config is nil")
	}

	conf := amqp.URI{
		Scheme:   "amqp",
		Host:     r.config.Host,
		Port:     r.config.Port,
		Username: r.config.User,
		Password: r.config.Password,
		Vhost:    r.config.Vhost,
	}.String()

	var err error

	r.conn, err = amqp.Dial(conf)

	if err != nil {
		return nil, err
	}

	log.Println("connected")
	return r, err
}

func (r *RabbitMQ) Channel() (*amqp.Channel, error) {
	if r.conn == nil {
		return nil, errors.New("not connected")
	}

	return r.conn.Channel()
}

func (r *RabbitMQ) queueDeclare(queue Queue, channel *amqp.Channel) (amqp.Queue, error) {
	q, err := channel.QueueDeclare(
		queue.Name,
		queue.Durable,
		queue.AutoDelete,
		queue.Exclusive,
		queue.NoWait,
		queue.Args,
	)

	return q, err
}

func (r *RabbitMQ) exchangeDeclare(exchange Exchange, channel *amqp.Channel) (*amqp.Channel, error) {
	return channel, channel.ExchangeDeclare(
		exchange.Name,
		exchange.Type,
		exchange.Durable,
		exchange.AutoDelete,
		exchange.Internal,
		exchange.NoWait,
		exchange.Args,
	)
}
