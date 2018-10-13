package models

import (
	fmt "fmt"
	utils "hermes-messaging-service/utils"

	amqp "github.com/streadway/amqp"
)

const (
	// WorkQueueName : Name of the Job Queue containing incoming messages
	WorkQueueName = "messages-incoming"

	// ConversationsIncomingExchangeName : Name of incoming messages exchange
	ConversationsIncomingExchangeName = "conversations.incoming.dx"

	// ConversationsPrivateExchangeName : Name of private messages exchange
	ConversationsPrivateExchangeName = "conversations.private.tx"
)

// MessagingInterface : Messaging interface component to be injected through dependency injection accross the packages
type MessagingInterface interface {
	DeclareExchange(name string, kind string, durable bool, autodelete bool, internal bool, noWait bool, args amqp.Table) error
	WorkQueue() (*amqp.Queue, error)
	BindQueueWithExchange(queue *amqp.Queue, routingKey string, exchange string) error
	CloseChannelAndConnection()
}

// RabbitMQ : RabbitMQ Struct encapsulating connection to broker and channel
type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

// NewRabbitMQ : Return a new RabbitMQ communication interface
func NewRabbitMQ(brokerURL string) (*RabbitMQ, error) {

	// TODO: Change this to use DialTLS

	// Get connection to broker
	// Connection abstracts the socket connection and
	// takes care of protocol version negotiation and authentication
	conn, err := amqp.Dial(brokerURL)

	if err != nil {
		utils.PanicOnError(err, "Failed to connect to RabbitMQ")
		return nil, err
	}

	// Open a RabbitMQ Channel
	// The channel is where most of the API for getting things done resides
	ch, err := conn.Channel()

	if err != nil {
		utils.PanicOnError(err, "Failed to open RabbitMQ Channel")
		return nil, err
	}

	return &RabbitMQ{conn, ch}, nil
}

// WorkQueue : Initialize/Get the worker queue responsible to queue incoming messages
// This queue is durable and non-auto-deleted. It is then designed to survive server restarts and remain
// when there are no remaining consumers or bindings.
// Persistent publishings will be restored in this queue on server restart.
// This queue must be bound the durable direct exchange << conversations.incoming.dx >>
func (rabbitMQ *RabbitMQ) WorkQueue() (*amqp.Queue, error) {

	q, err := rabbitMQ.Channel.QueueDeclare(
		WorkQueueName, // Queue Name
		true,          // Durable
		false,         // Auto-Delete
		false,         // Exclusive
		false,         // No-Wait
		nil,           // Arguments
	)

	if err != nil {
		utils.PanicOnError(err, "Failed to create worker queue")
		return nil, err
	}

	return &q, nil
}

// DeclareExchange : Init exchange with params
func (rabbitMQ *RabbitMQ) DeclareExchange(name string, kind string, durable bool, autodelete bool, internal bool, noWait bool, args amqp.Table) error {

	err := rabbitMQ.Channel.ExchangeDeclare(
		name,       // Exchange Name
		kind,       // Type
		durable,    // Durable
		autodelete, // Auto-Delete
		internal,   // Internal
		noWait,     // No-Wait
		args,       // Arguments
	)

	if err != nil {
		utils.PanicOnError(err, fmt.Sprintf("Failed to declare << %s >> exchange", name))
		return err
	}

	return nil
}

// BindQueueWithExchange : Bind Queue with chosen exchange
func (rabbitMQ *RabbitMQ) BindQueueWithExchange(queue *amqp.Queue, routingKey string, exchange string) error {

	err := rabbitMQ.Channel.QueueBind(
		queue.Name, // Queue name
		routingKey, // Routing key
		exchange,   // Exchange
		false,      // No-Wait
		nil,        // Args
	)

	if err != nil {
		utils.PanicOnError(err, fmt.Sprintf("Failed to bind queue : %s with exchange : << %s >>", queue.Name, exchange))
		return err
	}

	return nil
}

// CloseChannelAndConnection : Close RabbitMQ channel and connection
func (rabbitMQ *RabbitMQ) CloseChannelAndConnection() {

	err := rabbitMQ.Channel.Close()
	err = rabbitMQ.Connection.Close()

	if err != nil {
		utils.PanicOnError(err, "Failed to close RabbitMQ channel or connection")
	}
}
