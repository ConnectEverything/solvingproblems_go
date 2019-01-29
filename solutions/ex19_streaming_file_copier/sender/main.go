package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	stan "github.com/nats-io/go-nats-streaming"
)

func main() {
	channel := "ex19"
	url := os.Getenv("SP_NATS_SERVER")

	if url == "" {
		url = "nats://localhost:4223"
	}

	sc, err := stan.Connect("nats_course", "ex19", stan.NatsURL(url))

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

	for i, w := range words {
		sc.Publish(channel, []byte(fmt.Sprintf("%d|%d|%s", i, wordCount, w)))
	}

	sc.Close()
}
