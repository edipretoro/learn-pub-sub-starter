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
	const dsnRabbitMQ = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(dsnRabbitMQ)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s\n", err)
	}
	defer conn.Close()
	fmt.Println("Successfully connected to RabbitMQ")
	fmt.Println("Starting Peril server...")

	c, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s\n", err)
	}
	fmt.Println("Channel opened successfully")

	playingState := routing.PlayingState{
		IsPaused: true,
	}
	pubsub.PublishJSON(
		c,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		playingState,
	)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	s := <-signalChan
	fmt.Println("Received signal:", s)
	fmt.Println("Shutting down Peril server gracefully...")
}
