package tatslack

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// Client connects to the Slack API.
type Client struct {
	Token string
}

// ChannelHistory returns a list of messages.
func (c *Client) ChannelHistory(channel string) (*Response, error) {
	// Generate URL.
	u := &url.URL{
		Scheme: "https",
		Host:   "slack.com",
		Path:   "/api/channels.history",
	}
	v := url.Values{}
	v.Set("token", c.Token)
	v.Set("channel", channel)
	u.RawQuery = v.Encode()

	// Request from Slack.
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Unmarshal response.
	r := &Response{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	} else if !r.OK {
		return nil, errors.New(r.Error)
	}

	return r, nil
}

// UsersList returns a list of users.
func (c *Client) UsersList() (*Response, error) {
	// Generate URL.
	u := &url.URL{
		Scheme: "https",
		Host:   "slack.com",
		Path:   "/api/users.list",
	}
	v := url.Values{}
	v.Set("token", c.Token)
	u.RawQuery = v.Encode()

	// Request from Slack.
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Unmarshal response.
	r := &Response{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	} else if !r.OK {
		return nil, errors.New(r.Error)
	}

	return r, nil
}

type Response struct {
	OK        bool       `json:"ok"`
	Error     string     `json:"error"`
	Messages  []*Message `json:"messages"`
	Members   []*User    `json:"members"`
	HasMore   bool       `json:"has_more"`
	IsLimited bool       `json:"is_limited"`
}

type Message struct {
	Type string `json:"type"`
	TS   string `json:"ts"`
	User string `json:"user"`
	Text string `json:"text"`
}

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Profile struct {
		Image24 string `json:"image_24"`
		Image32 string `json:"image_32"`
		Image48 string `json:"image_48"`
	} `json:"profile"`
}

func warn(v ...interface{})              { fmt.Fprintln(os.Stderr, v...) }
func warnf(msg string, v ...interface{}) { fmt.Fprintf(os.Stderr, msg+"\n", v...) }
