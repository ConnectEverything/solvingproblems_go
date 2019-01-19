package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

	for i := 0; i < wordCount; i++ {
		wordToSend := fmt.Sprintf("%d|%d|%s", i, wordCount, words[i])
		ack, err := nc.Request(subject, []byte(wordToSend), time.Millisecond*50)

		if ack == nil || err != nil {
			fmt.Printf("Retrying %d\n", i)
			i-- // repeat this one
		}
	}

	fmt.Println("File Sent.")
	fmt.Println()

	nc.Close()
}
