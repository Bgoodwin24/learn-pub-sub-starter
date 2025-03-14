package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func main() {
	fmt.Println("Starting Peril client...")

	const conString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(conString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v\n", err)
	}

	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Failed to set username: %v\n", err)
	}

	queueName := routing.PauseKey + "." + username

	ch, q, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		queueName,
		routing.PauseKey,
		pubsub.QueueTypeTransient,
	)
	if err != nil {
		log.Fatalf("Failed to declare and bind queue: %v\n", err)
	}

	defer ch.Close()

	fmt.Printf("Created queue: %s\n", q.Name)
	fmt.Println("Client running. Press Enter to exit...")
	fmt.Scanln()
}
