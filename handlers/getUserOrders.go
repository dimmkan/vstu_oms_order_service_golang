package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"log"
	"vstu_oms_order_service/config"
	"vstu_oms_order_service/service"

	amqp "github.com/rabbitmq/amqp091-go"
)

func GetUserOrders(ctx context.Context, d amqp.Delivery, ch *amqp.Channel) {
	message, _ := service.Deserialize[service.GetUserOrdersType](d.Body)

	request_url := fmt.Sprintf("%s/items/orders?filter[user_id][_eq]=%v", config.New().Directus.DIRECTUS_HOST, message.User_id)
	client := &http.Client{Timeout: time.Second * 10}

	req, err := http.NewRequest("GET", request_url, nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("staticToken", config.New().Directus.ADMIN_API_KEY)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	service.FailOnError(err, "Failed to read get response")

	err = ch.PublishWithContext(ctx,
		"",        // exchange
		d.ReplyTo, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: d.CorrelationId,
			Body:          []byte(body),
		})
	service.FailOnError(err, "Failed to publish a message")
	d.Ack(false)
}
