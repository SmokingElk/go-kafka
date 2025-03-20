package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"go-kafka/internal/infrastructure"
)

var (
	groupsCnt = 3
	msgsCnt   = 100
)

func newGroups(count int) []string {
	res := make([]string, count)

	for i := 0; i < count; i++ {
		res[i] = uuid.NewString()
	}

	return res
}

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal(err)
	}

	serversStr := os.Getenv("SERVERS")
	servers := strings.Split(serversStr, ",")

	topic := os.Getenv("TOPIC")

	pr, err := infrastructure.NewProducer(servers)

	if err != nil {
		log.Fatal(err)
	}

	defer pr.Close()

	groups := newGroups(groupsCnt)

	for i := 0; i < msgsCnt; i++ {
		id := groups[i%len(groups)]
		err := pr.Produce(id, fmt.Sprintf("Message %d from producer", i), topic)
		if err != nil {
			fmt.Println(err)
		}
	}

	log.Printf("Sended %d messages\n", msgsCnt)
}
