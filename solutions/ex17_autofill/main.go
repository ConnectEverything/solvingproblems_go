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
	subject := "prefix"
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

	var response string
	input := os.Args[1]
	replyTo := "solution.autofill"

	sub, err := nc.SubscribeSync(replyTo)
	if err != nil {
		log.Fatal(err)
	}
	nc.Flush()

	// Send the request
	nc.PublishRequest(subject, replyTo, []byte(input))

	// Wait for a single response
	for {
		msg, err := sub.NextMsg(1 * time.Second)
		if err != nil {
			log.Fatal(err)
		}

		response = string(msg.Data)
		break
	}
	sub.Unsubscribe()

	if response == "" {
		fmt.Println("No prefixes found")
		os.Exit(-1)
	}

	lines := strings.Split(response, "\n")

	fmt.Printf("Prefixes from %s\n", lines[0])

	// skip the ldap
	for i := 1; i < len(lines); i++ {
		fmt.Println(lines[i])
	}
}
