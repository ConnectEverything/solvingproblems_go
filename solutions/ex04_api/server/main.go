package main

import (
	"log"
	"os"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	subject := "ex04"
	url := os.Getenv("SP_NATS_SERVER")

	if url == "" {
		url = "nats://localhost:4222"
	}

	log.Printf("Connecting to %s", url)

	nc, err := nats.Connect(url,
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Printf("Error: %v", err)
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			log.Printf("Connection closed")
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	sub, err := nc.SubscribeSync(subject)
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg, err := sub.NextMsg(1000 * time.Second)
		if err != nil {
			log.Fatal(err)
		}
		now := time.Now()
		response := now.Format(time.RFC822)
		nc.Publish(msg.Reply, []byte(response))
	}
}
