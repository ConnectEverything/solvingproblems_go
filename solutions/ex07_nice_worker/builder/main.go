package main

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nats-io/go-nats"
)

func main() {
	workSubject := "ex07_work"
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
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	palette := palette.Plan9

	rand.Shuffle(len(palette), func(i, j int) {
		palette[i], palette[j] = palette[j], palette[i]
	})

	fmt.Print("Working ")

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			work := fmt.Sprintf("%d,%d", col, row) // x,y
			response, err := nc.Request(workSubject, []byte(work), time.Second*5)

			if err != nil {
				fmt.Print("x")
				col--
				continue
			}

			str := string(response.Data)
			parts := strings.Split(str, ",")
			col, _ := strconv.Atoi(parts[0])
			row, _ := strconv.Atoi(parts[1])
			iterations, _ := strconv.Atoi(parts[2])
			color := palette[iterations]

			img.Set(col, row, color)
		}
		fmt.Print(".")
	}

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
	nc.Close()
}
