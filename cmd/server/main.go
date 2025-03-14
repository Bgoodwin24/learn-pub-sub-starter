package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const conString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(conString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v\n", err)
	}

	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	connCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create new channel: %v", err)
	}

	err = pubsub.PublishJSON(connCh,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{
			IsPaused: true,
		},
	)

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Peril is shutting down.")
}
