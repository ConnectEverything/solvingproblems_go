package main

import (
	"fmt"
	"log"
	"os"

	stan "github.com/nats-io/go-nats-streaming"
)

func main() {
	channel := "ex20"
	url := os.Getenv("SP_NATS_SERVER")

	if url == "" {
		url = "nats://localhost:4223"
	}

	sc, err := stan.Connect("nats_course", "ex20", stan.NatsURL(url))

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		sc.Publish(channel, []byte(fmt.Sprintf("%d", i)))
	}

	sc.Close()
}
