package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Peril is shutting down.")
}
