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
	subject := ">"
	url := os.Getenv("SP_NATS_SERVER")

	if url == "" {
		url = "nats://localhost:4222"
	}

	nc, err := nats.Connect(url)

	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	sub, err := nc.SubscribeSync(subject)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Monitoring...")
	for {
		msg, err := sub.NextMsg(1000 * time.Second)
		if err != nil || msg == nil {
			fmt.Println("<empty message>")
			continue
		}

		msgString := string(msg.Data)

		msgString = strings.Replace(msgString, "|", " ", -1) // fix shell print issue
		fmt.Printf("%s: %s\n", msg.Subject, msgString)
	}
}
