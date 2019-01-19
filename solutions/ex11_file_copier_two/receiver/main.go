package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	subject := "ex11"
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

	var words []string
	var received []bool

	for {
		msg, err := sub.NextMsg(1000 * time.Second)
		if err != nil || msg == nil {
			fmt.Println("Message timeout, continuing to work")
			continue
		}

		msgString := string(msg.Data)
		split := strings.Split(msgString, "|")
		index, _ := strconv.Atoi(split[0])
		max, _ := strconv.Atoi(split[1])
		word := ""

		if len(split) == 3 {
			word = split[2]
		}

		if words == nil {
			words = make([]string, max)
			received = make([]bool, max)
		}

		if rand.Float64() < 0.05 {
			fmt.Printf("Skipping %d\n", index)
			continue
		}

		words[index] = word
		received[index] = true
		nc.Publish(msg.Reply, nil) // send the ack

		if index == max-1 {
			break
		}
	}
	nc.Flush() // get all the acks to the server

	fmt.Println("Received file...")
	fmt.Println()

	for i, w := range words {
		if i != 0 {
			fmt.Print(" ")
		}
		fmt.Print(w)
	}
	fmt.Println()
	nc.Close()
}
