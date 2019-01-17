package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	subject := "ex03"
	prefix := "solution"
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

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			r := rand.Int31n(10000)
			notification := fmt.Sprintf("%s.%d", prefix, r)
			nc.Publish(subject, []byte(notification))
		}
	}()

	sub, err := nc.SubscribeSync(subject)
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg, err := sub.NextMsg(100 * time.Second)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s", msg.Data)
	}
}
