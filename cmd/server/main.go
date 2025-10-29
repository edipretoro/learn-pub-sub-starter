package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		fmt.Sprintf("%s.%s", routing.GameLogSlug, "*"),
		pubsub.SimpleQueueDurable,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gamelogic.PrintServerHelp()
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			fmt.Println("Sending a pause message...")
			playingState := routing.PlayingState{
				IsPaused: true,
			}
			err := pubsub.PublishJSON(
				c,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				playingState,
			)
			if err != nil {
				fmt.Printf("Could not publish pause message: %v\n", err)
			} else {
				fmt.Println("Pause message published.")
			}
		case "resume":
			fmt.Println("Sending a resume message...")
			playingState := routing.PlayingState{
				IsPaused: false,
			}
			err := pubsub.PublishJSON(
				c,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				playingState,
			)
			if err != nil {
				fmt.Printf("Could not publish resume message: %v\n", err)
			} else {
				fmt.Println("Resume message published.")
			}
		case "quit":
			fmt.Println("Quitting Peril server...")
			return
		case "help":
			gamelogic.PrintServerHelp()
		default:
			fmt.Printf("Unknown command: %s\n", words[0])
		}
	}
}
