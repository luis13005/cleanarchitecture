package handler

import (
	"encoding/json"
	"sync"

	"github.com/luis13005/cleanarchitecture/internal/event"
	"github.com/streadway/amqp"
)

type OrderCreatedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewOrderCreatedHandler(rabbitMQChannel *amqp.Channel) *OrderCreatedHandler {
	return &OrderCreatedHandler{RabbitMQChannel: rabbitMQChannel}
}

func (h *OrderCreatedHandler) Handle(event event.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	jsonOutput, _ := json.Marshal(event.GetPayLoad())

	msgRabbitmq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	h.RabbitMQChannel.Publish(
		"amq.direct",
		"",
		false,
		false,
		msgRabbitmq,
	)
}
