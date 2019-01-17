package main

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/png"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/nats-io/go-nats"
)

func main() {
	workSubject := "ex05_work"
	completeSubject := "ex05_complete"
	url := os.Getenv("SP_NATS_SERVER")

	if url == "" {
		url = "nats://localhost:4222"
	}

	nc, err := nats.Connect(url,
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Printf("Error: %v", err)
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			log.Printf("Connection closed")
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	width := 800
	height := 600
	remaining := width * height
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	palette := palette.Plan9

	workQueue := make(chan string, width*height)

	// Add all the work
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			work := fmt.Sprintf("%d,%d", col, row) // x,y
			workQueue <- work
		}
	}

	fmt.Print("Working ")

	if _, err := nc.Subscribe(workSubject, func(msg *nats.Msg) {
		select {
		case work := <-workQueue:
			err := nc.Publish(msg.Reply, []byte(work))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		default:
		}
	}); err != nil {
		log.Fatal(err)
	}

	if _, err := nc.Subscribe(completeSubject, func(msg *nats.Msg) {
		work := string(msg.Data)
		parts := strings.Split(work, ",")
		col, _ := strconv.Atoi(parts[0])
		row, _ := strconv.Atoi(parts[1])
		iterations, _ := strconv.Atoi(parts[2])
		color := palette[iterations]

		img.Set(col, row, color)

		if remaining%10000 == 0 {
			fmt.Printf(".")
		}

		remaining--

		if remaining == 0 {
			fmt.Println(" done")
			out, err := os.Create("mandelbrot.png")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = png.Encode(out, img)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Image saved")
			os.Exit(0)
		}
	}); err != nil {
		log.Fatal(err)
	}

	runtime.Goexit()
}
