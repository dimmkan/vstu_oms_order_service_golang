package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"log"
	"vstu_oms_order_service/config"
	"vstu_oms_order_service/service"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ChangeOrderDescription(ctx context.Context, d amqp.Delivery, ch *amqp.Channel) {
	message, _ := service.Deserialize[service.ChangeOrderDescriptionType](d.Body)

	request_url := fmt.Sprintf("%s/items/orders/%s", config.New().Directus.DIRECTUS_HOST, message.Order_id)
	client := &http.Client{Timeout: time.Second * 10}

	req, err := http.NewRequest("PATCH", request_url, bytes.NewBuffer(d.Body))
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

	response, err := service.Serialize(map[string]bool{
		"success": true,
	})
	service.FailOnError(err, "Failed to response serialized")

	err = ch.PublishWithContext(ctx,
		"",        // exchange
		d.ReplyTo, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: d.CorrelationId,
			Body:          response,
		})
	service.FailOnError(err, "Failed to publish a message")
	d.Ack(false)
}
