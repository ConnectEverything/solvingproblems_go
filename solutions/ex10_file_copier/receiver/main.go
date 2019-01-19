package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	subject := "ex10"
	replayPrefix := "ex10.replay"
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

		words[index] = word
		received[index] = true

		if index == max-1 {
			break
		}
	}

	fmt.Println("Checking for missing sequence numbers...")
	for i := range words {
		if !received[i] {
			fmt.Printf("\tRequesting %d\n", i)
			replay, err := nc.Request(fmt.Sprintf("%s.%d", replayPrefix, i), nil, time.Second)

			if err != nil || replay == nil {
				fmt.Println("No response on replay, exiting...")
				os.Exit(-1)
			}

			replayString := string(replay.Data)
			replaySplit := strings.Split(replayString, "|")
			index, _ := strconv.Atoi(replaySplit[0])
			word := ""

			if len(replaySplit) == 3 {
				word = replaySplit[2]
			}

			words[index] = word
		}
	}

	// Tell the sender we are done
	nc.Publish(fmt.Sprintf("%s.%d", replayPrefix, -1), nil)

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
