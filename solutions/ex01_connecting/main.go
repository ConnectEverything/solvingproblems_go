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
		nats.DiscoveredServersHandler(func(nc *nats.Conn) {
			log.Printf("Known servers: %v\n", nc.Servers())
			log.Printf("Discovered servers: %v\n", nc.DiscoveredServers())
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			log.Printf("Connection closed")
		}),
		nats.DisconnectHandler(func(_ *nats.Conn) {
			log.Printf("Connection disconnected")
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Printf("Connection reconnected")
		}),
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
	log.Printf("Known servers: %v\n", nc.Servers())
	log.Printf("Discovered servers: %v\n", nc.DiscoveredServers())

	time.Sleep(time.Second * 5)

	nc.Close()

	time.Sleep(time.Second * 1)
}
