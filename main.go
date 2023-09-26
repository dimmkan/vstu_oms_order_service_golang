package main

import (
	"fmt"
	"log"
	"sync"
	"vstu_oms_order_service/config"
	"vstu_oms_order_service/handlers"
	"vstu_oms_order_service/service"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found!")
	}
}

func isValidRoutingKey(routingKey string) bool {
	validKeys := map[string]struct{}{
		"order.create.command":            struct{}{},
		"order.changestatus.command":      struct{}{},
		"order.changedescription.command": struct{}{},
		"order.delete.command":            struct{}{},
		"order.getbyuser.query":           struct{}{},
	}

	_, ok := validKeys[routingKey]
	return ok
}

func main() {
	conf := config.New()
	amqpConnectionString := fmt.Sprintf("amqp://%s:%s@%s:5672/", conf.Amqp.AMQP_USER, conf.Amqp.AMQP_PASSWORD, conf.Amqp.AMQP_HOST)

	// Open connect
	conn, err := amqp.Dial(amqpConnectionString)
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

	log.Printf(" [*] Waiting for requests. To exit press CTRL+C")

	// Ожидание завершения всех горутин
	var wg sync.WaitGroup

	// Обработка полученных сообщений
	for msg := range msgs {
		if !isValidRoutingKey(msg.RoutingKey) {
			// Если RoutingKey не в списке допустимых, игнорируем его
			log.Printf("Received message with invalid routing key: %s", msg.RoutingKey)
			msg.Nack(false, false) // Отправляем сообщение обратно в очередь
			continue
		}

		wg.Add(1)
		go func(msg amqp.Delivery, channel *amqp.Channel) {
			defer wg.Done()
			switch msg.RoutingKey {
			case "order.create.command":
				{
					handlers.CreateOrder(msg, channel)
				}
			case "order.changestatus.command":
				{
					handlers.ChangeOrderStatus(msg, channel)
				}
			case "order.changedescription.command":
				{
					handlers.ChangeOrderDescription(msg, channel)
				}
			case "order.delete.command":
				{
					handlers.DeleteOrder(msg, channel)
				}
			case "order.getbyuser.query":
				{
					handlers.GetUserOrders(msg, channel)
				}
			}
		}(msg, ch)
	}

	// Ожидание завершения всех горутин
	wg.Wait()
}
