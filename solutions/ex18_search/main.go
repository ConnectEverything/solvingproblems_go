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
	subject := "match"
	url := os.Getenv("SP_NATS_SERVER")

	if url == "" {
		url = "nats://localhost:4222"
	}

	log.Printf("Connecting to %s", url)

	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	minResponses := 2
	responses := []string{}
	input := os.Args[1]
	replyTo := "solution.search"

	sub, err := nc.SubscribeSync(replyTo)
	if err != nil {
		log.Fatal(err)
	}
	nc.Flush()

	// Send the request
	nc.PublishRequest(subject, replyTo, []byte(input))

	// Wait for a single response
	max := 500 * time.Millisecond
	start := time.Now()
	for time.Now().Sub(start) < max {
		msg, err := sub.NextMsg(1 * time.Second)
		if err != nil {
			break
		}

		responses = append(responses, string(msg.Data))

		if len(responses) >= minResponses {
			break
		}
	}
	sub.Unsubscribe()

	if len(responses) == 0 {
		fmt.Println("No matches found")
		os.Exit(-1)
	}

	matches := map[string]string{}

	for _, response := range responses {
		lines := strings.Split(response, "\n")

		// skip the ldap
		for i := 1; i < len(lines); i++ {
			matches[lines[i]] = lines[i]
		}
	}

	fmt.Printf("Received matches from %d servers\n", len(responses))
	fmt.Println()

	for _, s := range matches {
		fmt.Println(s)
	}
}
