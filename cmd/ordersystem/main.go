package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/luis13005/cleanarchitecture/configs"
	"github.com/luis13005/cleanarchitecture/internal/event/handler"
	"github.com/luis13005/cleanarchitecture/pkg/events"
	"github.com/streadway/amqp"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", configs.DBName, configs.DBPassword, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch
}
