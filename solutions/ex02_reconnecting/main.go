package main

import (
	"log"
	"os"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
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
		nats.DisconnectHandler(func(_ *nats.Conn) {
			log.Printf("Connection disconnected")
		}),
		nats.ReconnectHandler(func(conn *nats.Conn) {
			log.Printf("Connection reconnected to %s", conn.ConnectedUrl())
		}),
		nats.ReconnectWait(time.Second*5),
		nats.MaxReconnects(5),
	)

	if err != nil {
		log.Fatal(err)
	}

	getStatusTxt := func(nc *nats.Conn) string {
		switch nc.Status() {
		case nats.CONNECTED:
			return "Connected"
		case nats.CLOSED:
			return "Closed"
		default:
			return "Other"
		}
	}
	log.Printf("The connection is %v\n", getStatusTxt(nc))

	time.Sleep(time.Second * 10 * 60)

	nc.Close()

	time.Sleep(time.Second * 1)
}
