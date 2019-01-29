package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
)

func main() {
	channel := "ex20"
	maxInFlight := 1
	url := os.Getenv("SP_NATS_SERVER")

	if url == "" {
		url = "nats://localhost:4223"
	}

	if len(os.Args) > 1 {
		maxInFlight, _ = strconv.Atoi(os.Args[1])
	}

	sc, err := stan.Connect("nats_course", "ex20", stan.NatsURL(url))
	if err != nil {
		log.Fatal(err)
	}

	defer sc.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	count := 0

	sub, err := sc.Subscribe(channel, func(msg *stan.Msg) {
		msgString := string(msg.Data)
		fmt.Println(msgString)

		if count%2 == 0 {
			msg.Ack()
		}

		count++
	}, stan.DeliverAllAvailable(), stan.MaxInflight(maxInFlight),
		stan.SetManualAckMode(), stan.AckWait(time.Second))

	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	sub.Close()
	sc.Close()
}
