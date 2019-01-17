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

	nc, err := nats.Connect(url)

	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	msg, err := nc.Request(subject, nil, time.Second)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(string(msg.Data))
}
