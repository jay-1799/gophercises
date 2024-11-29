package x

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func New(key, secret string) (*Client, error) {
	client, err := twitterClient(key, secret)
	if err != nil {
		return nil, err
	}
	return &Client{
		client: client,
	}, nil
}

type Client struct {

	//unexported
	client *http.Client
}

func twitterClient(key, secret string) (*http.Client, error) {
	req, err := http.NewRequest("POST", "https://api.x.com/oauth2/token",
		strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(key, secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var token oauth2.Token
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&token)
	if err != nil {
		return nil, err
	}
	// fmt.Println(token)

	var conf oauth2.Config
	return conf.Client(context.Background(), &token), nil
}

func (c *Client) Retweeters(tweetID string) ([]string, error) {

	url := fmt.Sprintf("https://api.x.com/1.1/statuses/retweets/%s.json", tweetID)
	res, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var retweets []struct {
		User struct {
			ScreenName string `json:"screen_name"`
		} `json: "user"`
	}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&retweets)
	if err != nil {
		return nil, err
	}
	// var usernames []string
	usernames := make([]string, 0, len(retweets))
	for _, retweet := range retweets {
		usernames = append(usernames, retweet.User.ScreenName)
	}
	return usernames, nil
}
