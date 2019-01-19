package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/nats-io/go-nats"
)

func main() {
	subject := "ex10"
	replaySubject := "ex10.replay.*"
	url := os.Getenv("SP_NATS_SERVER")

	if url == "" {
		url = "nats://localhost:4222"
	}

	nc, err := nats.Connect(url,
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			fmt.Printf("Error: %v", err)
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			fmt.Printf("Connection closed")
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	filePath := os.Args[1]
	encoded, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	data := string(encoded)

	words := strings.Split(data, " ")
	wordCount := len(words)

	if _, err := nc.Subscribe(replaySubject, func(msg *nats.Msg) {
		split := strings.Split(msg.Subject, ".")
		if len(split) != 3 {
			fmt.Printf("Got subject with wrong number of tokens %s (%d)\n", msg.Subject, len(split))
			return
		}

		index, err := strconv.Atoi(split[2])
		if err != nil {
			log.Fatal(err)
		}

		if index < 0 || index >= wordCount {
			fmt.Println("Copy complete, exiting...")
			nc.Close()
			os.Exit(0)
		}

		resp := fmt.Sprintf("%d|%d|%s", index, wordCount, words[index])
		fmt.Printf("Replaying %d %d %s\n", index, wordCount, words[index])
		nc.Publish(msg.Reply, []byte(resp))
	}); err != nil {
		log.Fatal(err)
	}

	for i, w := range words {
		if rand.Float64() < 0.05 { // skip some words
			fmt.Printf("Skipping %d\n", i)
			continue
		}
		nc.Publish(subject, []byte(fmt.Sprintf("%d|%d|%s", i, wordCount, w)))
	}

	runtime.Goexit()
}
