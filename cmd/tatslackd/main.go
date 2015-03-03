package main

import (
	"flag"
	"log"

	"github.com/benbjohnson/tatslack"
)

func main() {
	log.SetFlags(0)

	// Parse the flags.
	token := flag.String("token", "", "slack API token")
	flag.Parse()

	// Ensure we have an API token.
	if *token == "" {
		log.Fatal("token required")
	}

	// Create client.
	c := tatslack.Client{
		Token: *token,
	}

	// Read twoalarmtuesday channel messages.
	resp, err := c.ChannelHistory("C03HCH8RB")
	if err != nil {
		log.Fatal(err)
	}

	// Print messages.
	log.Println("MESSAGES")
	log.Println("========")
	for _, m := range resp.Messages {
		log.Println(">", m.Text)
	}
}
