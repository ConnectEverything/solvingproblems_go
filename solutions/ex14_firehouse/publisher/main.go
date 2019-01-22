package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	subject := "ex14"
	url := os.Getenv("SP_NATS_SERVER")

	if url == "" {
		url = "nats://localhost:4222"
	}

	fmt.Printf("Connecting to %s\n", url)

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

	counter := 0
	fmt.Print("Sending ")

	for {
		counter++
		if counter%1000 == 0 {
			fmt.Print(".")
		}

		err := nc.Publish(subject, []byte(fmt.Sprintf("%d", counter)))

		if err != nil {
			fmt.Println()
			fmt.Printf("Exception - %s\n", err.Error())
		}

		time.Sleep(1 * time.Millisecond)
	}
}
