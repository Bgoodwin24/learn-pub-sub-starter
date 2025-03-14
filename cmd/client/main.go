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

	go func() {
		fmt.Printf("Created queue: %s\n", q.Name)
		fmt.Println("Client running. Press Enter to exit...")
		fmt.Scanln()
	}()

	gamestate := gamelogic.NewGameState(username)

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		if words[0] == "spawn" {
			err := gamestate.CommandSpawn(words)
			if err != nil {
				fmt.Printf("Failed to spawn unit: %v\n", err)
			}
		} else if words[0] == "move" {
			_, err := gamestate.CommandMove(words)
			if err != nil {
				fmt.Printf("Failed to move unit: %v\n", err)
			} else {
				fmt.Println("Unit moved successfully")
			}
		} else if words[0] == "status" {
			gamestate.CommandStatus()
		} else if words[0] == "help" {
			gamelogic.PrintClientHelp()
		} else if words[0] == "spam" {
			fmt.Println("Spamming not allowed yet!")
		} else if words[0] == "quit" {
			gamelogic.PrintQuit()
			break
		} else {
			fmt.Println("Unknown command")
			continue
		}

	}

}
