package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/benbjohnson/tatslack"
)

const CommunityChannel = "C029P2S9X"

const TATChannel = "C03HCH8RB"

func main() {
	log.SetFlags(0)

	// Parse the flags.
	token := flag.String("token", "", "slack API token")
	flag.Parse()

	// Ensure we have an API token.
	if *token == "" {
		log.Fatal("token required")
	}

	// Open the database.
	db, err := tatslack.Open("db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/*
		// Create client.
		c := tatslack.Client{
			Token: *token,
		}

		// Read twoalarmtuesday channel messages.
		resp, err := c.ChannelHistory(CommunityChannel)
		if err != nil {
			log.Fatal(err)
		}

		// Write the messages in.
		if err := db.SaveMessages(CommunityChannel, resp.Messages); err != nil {
			log.Fatal(err)
		}
	*/

	log.Println("listening on http://localhost:9000")

	// Run HTTP server.
	h := &tatslack.Handler{
		DB:      db,
		Channel: CommunityChannel,
	}
	http.ListenAndServe(":9000", h)
}
