package main

import (
	"bufio"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/nats-io/go-nats"
)

func main() {
	ldap := "solution"
	matchSubject := "match"
	prefixSubject := "prefix"
	url := os.Getenv("SP_NATS_SERVER")

	if url == "" {
		url = "nats://localhost:4222"
	}

	if len(os.Args) > 1 {
		ldap = os.Args[1]
	}

	log.Printf("Connecting to %s", url)

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

	wordList := []string{}

	file, err := os.Open("resources/words_alpha.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wordList = append(wordList, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	file.Close()

	if _, err := nc.Subscribe(matchSubject, func(msg *nats.Msg) {
		toMatch := string(msg.Data)
		response := []string{ldap}

		for _, word := range wordList {

			if strings.Contains(word, toMatch) {
				response = append(response, word)
			}

			if len(response) > 10 {
				break
			}
		}

		log.Printf("Returning %d matches for %s\n", len(response)-1, toMatch)
		nc.Publish(msg.Reply, []byte(strings.Join(response, "\n")))

	}); err != nil {
		log.Fatal(err)
	}

	if _, err := nc.Subscribe(prefixSubject, func(msg *nats.Msg) {
		toMatch := string(msg.Data)
		response := []string{ldap}

		for _, word := range wordList {

			if strings.HasPrefix(word, toMatch) {
				response = append(response, word)
			}

			if len(response) > 10 {
				break
			}
		}

		log.Printf("Returning %d prefixes for %s\n", len(response)-1, toMatch)
		nc.Publish(msg.Reply, []byte(strings.Join(response, "\n")))
	}); err != nil {
		log.Fatal(err)
	}

	runtime.Goexit()
}
