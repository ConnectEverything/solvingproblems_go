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
			if err == nats.ErrSlowConsumer {
				dropped, _ := sub.Dropped()
				fmt.Println()
				fmt.Printf("Slow consumer on subject %s dropped %d messages\n",
					sub.Subject, dropped)
				return
			}
			fmt.Printf("Error: %v", err)
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			fmt.Println()
			fmt.Println("Connection closed")
		}),
		nats.DisconnectHandler(func(_ *nats.Conn) {
			fmt.Println()
			fmt.Println("Connection disconnected")
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			fmt.Println()
			fmt.Println("Connection reconnected")
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

	sub.SetPendingLimits(100, 1024*8)

	count := 0

	fmt.Println()
	fmt.Print("Listening ")
	for {
		_, err := sub.NextMsg(100 * time.Second)
		if err != nil {
			log.Fatal(err)
		}
		count++

		if count == 1000 {
			fmt.Print(".")
			count = 0
		}
	}
}
