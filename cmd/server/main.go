package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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

	gamelogic.PrintServerHelp()

	ch, _, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.QueueTypeDurable,
	)
	if err != nil {
		log.Fatalf("Failed to declare and bind queue: %v\n", err)
	}

	defer ch.Close()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("Peril is shutting down.")
		os.Exit(0)
	}()

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		} else if words[0] == "pause" {
			log.Println("Sending pause message.")
			err := pubsub.PublishJSON(connCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: true,
				},
			)
			if err != nil {
				log.Fatalf("Failed to publish pause message: %v\n", err)
			}
		} else if words[0] == "resume" {
			log.Println("Sending resume message.")
			err := pubsub.PublishJSON(connCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				},
			)
			if err != nil {
				log.Fatalf("Failed to publish resume message: %v\n", err)
			}
		} else if words[0] == "quit" {
			log.Println("Exiting.")
			break
		} else {
			log.Println("Unknown command.")
		}
	}
}
