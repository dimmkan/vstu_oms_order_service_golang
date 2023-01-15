package main

import (
	"context"
	"fmt"
	"time"

	"log"
	"vstu_oms_order_service/config"
	"vstu_oms_order_service/handlers"
	"vstu_oms_order_service/service"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found!")
	}
}

func main() {
	conf := config.New()
	amqp_connection_string := fmt.Sprintf("amqp://%s:%s@%s:5672/", conf.Amqp.AMQP_USER, conf.Amqp.AMQP_PASSWORD, conf.Amqp.AMQP_HOST)

	// Open connect
	conn, err := amqp.Dial(amqp_connection_string)
	service.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//Create channel
	ch, err := conn.Channel()
	service.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	//Exchange declare
	err = ch.ExchangeDeclare(
		conf.Amqp.AMQP_EXCHANGE, // name
		"topic",                 // type
		true,                    // durable
		false,                   // auto-deleted
		false,                   // internal
		false,                   // no-wait
		nil,                     // arguments
	)
	service.FailOnError(err, "Failed to declare an exchange")

	//Queue declare
	q, err := ch.QueueDeclare(
		conf.Amqp.AMQP_QUEUE, // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	service.FailOnError(err, "Failed to declare a queue")

	//Bind queue
	err = ch.QueueBind(
		q.Name,                  // queue name
		"#",                     // routing key
		conf.Amqp.AMQP_EXCHANGE, // exchange
		false,
		nil)
	service.FailOnError(err, "Failed to bind a queue")

	//Set QoS
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	service.FailOnError(err, "Failed to set QoS")

	//Register consumer
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	service.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		for d := range msgs {
			switch d.RoutingKey {
			case "order.create.command":
				{
					handlers.CreateOrder(ctx, d, ch)
				}
			case "order.changestatus.command":
				{
					handlers.ChangeOrderStatus(ctx, d, ch)
				}
			case "order.changedescription.command":
				{
					handlers.ChangeOrderDescription(ctx, d, ch)
				}
			case "order.delete.command":
				{
					handlers.DeleteOrder(ctx, d, ch)
				}
			case "order.getbyuser.query":
				{
					handlers.GetUserOrders(ctx, d, ch)
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for requests. To exit press CTRL+C")
	<-forever
}
