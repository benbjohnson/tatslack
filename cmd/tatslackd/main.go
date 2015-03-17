package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/benbjohnson/tatslack"
)

const (
	CommunityChannel   = "C029P2S9X"
	TATChannel         = "C03HCH8RB"
	Channel1409        = "C02HNLCBU"
	Channel1410        = "C02Q114D3"
	Channel1412        = "C035BNUFY"
	Channel1502        = "C03J81F1L"
	WatercoolerChannel = "C02AZBZ81"
	TuringChannel      = "C03HS2LS0"
)

var Channels = []string{
	CommunityChannel, TuringChannel, TATChannel,
	Channel1409, Channel1410, Channel1412, Channel1502,
}

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

	// Create client.
	c := tatslack.Client{
		Token: *token,
	}

	// Read channel messages.
	log.Println("retrieving messages...")
	for i, channel := range Channels {
		log.Printf("%d. %s", i, channel)

		resp, err := c.ChannelHistory(channel)
		if err != nil {
			log.Fatal(err)
		}

		// Write the messages in.
		if err := db.SaveMessages(channel, resp.Messages); err != nil {
			log.Fatal(err)
		}
	}

	// Read users.
	log.Println("retrieving users...")
	resp, err := c.UsersList()
	if err != nil {
		log.Fatal(err)
	}
	db.SetUsers(resp.Members)

	log.Println("listening on http://localhost:9000")

	// Run HTTP server.
	h := &tatslack.Handler{
		DB:      db,
		Channel: CommunityChannel,
	}
	http.ListenAndServe(":9000", h)
}
