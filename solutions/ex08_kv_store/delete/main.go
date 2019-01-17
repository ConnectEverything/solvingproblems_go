package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nats-io/go-nats"
)

func main() {
	subjectPrefix := "ex08"
	url := os.Getenv("SP_NATS_SERVER")

	if url == "" {
		url = "nats://localhost:4222"
	}

	nc, err := nats.Connect(url)

	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	key := os.Args[1]

	fmt.Printf("Deleting %s\n", key)
	subject := fmt.Sprintf("%s.%s", subjectPrefix, key)
	err = nc.Publish(subject, nil)
	if err != nil {
		fmt.Printf("An error occurred %s\n", err.Error())
	}
}
