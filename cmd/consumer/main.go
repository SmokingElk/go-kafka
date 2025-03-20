package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"

	"go-kafka/internal/infrastructure"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal(err)
	}

	serversStr := os.Getenv("SERVERS")
	servers := strings.Split(serversStr, ",")

	consumerGroup := os.Getenv("CONSUMER_GROUP")
	topic := os.Getenv("TOPIC")

	cg, err := infrastructure.NewConsumer(servers, consumerGroup, topic, func(key, value string) {
		log.Printf("Get message {\n\tKey: %s\n\tValue: %s\n}\n", key, value)
	})

	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cg.Close()
}
