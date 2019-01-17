package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	subject := "ex08.*"
	url := os.Getenv("SP_NATS_SERVER")

	if url == "" {
		url = "nats://localhost:4222"
	}

	nc, err := nats.Connect(url)

	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	database := map[string][]byte{}

	sub, err := nc.SubscribeSync(subject)
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg, err := sub.NextMsg(1000 * time.Second)
		if err != nil || msg == nil {
			fmt.Printf("Message timeout, continuing to work")
			continue
		}

		subject := msg.Subject
		split := strings.Split(subject, ".")

		if len(split) != 2 {
			fmt.Printf("Got subject with wrong number of tokens %s (%d)\n", subject, len(split))
			continue
		}

		key := split[1]
		replyTo := msg.Reply
		value := msg.Data

		// Check for set/delete
		if replyTo == "" {
			if value == nil || len(value) == 0 {
				fmt.Printf("Deleting %s\n", key)
				delete(database, key)
			} else {
				fmt.Printf("Setting %s\n", key)
				database[key] = value
			}
			continue
		}

		if value != nil && len(value) > 0 {
			fmt.Printf("Received get request with data, ignoring, %s\n", subject)
			continue
		}

		fmt.Printf("Returning %s\n", key)
		value = database[key]
		nc.Publish(replyTo, value)
	}
}
