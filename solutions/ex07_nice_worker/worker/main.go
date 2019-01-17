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
	workSubject := "ex07_work"
	workQueue := "ex07_queue"
	url := os.Getenv("SP_NATS_SERVER")

	if url == "" {
		url = "nats://localhost:4222"
	}

	nc, err := nats.Connect(url, nats.DrainTimeout(10*time.Second))

	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	width := 800.0
	height := 600.0
	max := 255

	sub, err := nc.QueueSubscribeSync(workSubject, workQueue)
	if err != nil {
		log.Fatal(err)
	}

	count := 0

	for {
		msg, err := sub.NextMsg(1000 * time.Second)
		if err != nil {
			log.Fatal(err)
		}
		if msg == nil {
			log.Fatal("No more work")
		}

		work := string(msg.Data)
		parts := strings.Split(work, ",")
		col, _ := strconv.ParseFloat(parts[0], 64)
		row, _ := strconv.ParseFloat(parts[1], 64)
		x0 := (col - width/2.0) * 4.0 / width
		y0 := (row - height/2.0) * 4.0 / height

		x := 0.0
		y := 0.0
		iterations := 0

		for (x*x+y*y) < 4 && (iterations < max) {
			temp := (x*x - y*y) + x0
			y = (2 * x * y) + y0
			x = temp
			iterations++
		}

		response := fmt.Sprintf("%d,%d,%d", int(col), int(row), iterations)
		nc.Publish(msg.Reply, []byte(response))

		count++

		if count == 50000 {
			log.Println("To tired to go own")
			if err := nc.Drain(); err != nil {
				log.Fatal(err)
			}
			break
		}
	}

	nc.Close()
}
