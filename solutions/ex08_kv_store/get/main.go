package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

	subject := fmt.Sprintf("%s.%s", subjectPrefix, key)
	resp, err := nc.Request(subject, nil, time.Second)
	if err != nil {
		fmt.Printf("An error occurred %s\n", err.Error())
	}

	fmt.Printf("%s = %q\n", key, string(resp.Data))
}
