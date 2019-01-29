package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	stan "github.com/nats-io/go-nats-streaming"
)

func main() {
	channel := "ex19"
	startAt := 1
	url := os.Getenv("SP_NATS_SERVER")

	if url == "" {
		url = "nats://localhost:4223"
	}

	if len(os.Args) > 1 {
		startAt, _ = strconv.Atoi(os.Args[1])
	}

	sc, err := stan.Connect("nats_course", "ex19", stan.NatsURL(url))
	if err != nil {
		log.Fatal(err)
	}

	defer sc.Close()

	var words []string

	wg := sync.WaitGroup{}
	wg.Add(1)

	fmt.Println("Listening for file...")
	sub, err := sc.Subscribe(channel, func(msg *stan.Msg) {
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
		}

		words[index] = word

		if index == max-1 {
			wg.Done()
		}
	}, stan.StartAtSequence(uint64(startAt)))

	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	fmt.Println("Received file...")
	fmt.Println()
	for i, w := range words {
		if i != 0 {
			fmt.Print(" ")
		}
		fmt.Print(w)
	}
	fmt.Println()

	sub.Close()
	sc.Close()
}
