package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	var (
		keyFile   string
		usersFile string
		tweetID   string
	)
	flag.StringVar(&keyFile, "key", "keys.json", "the file where you store your consumer key and secret for the api")
	flag.StringVar(&usersFile, "users", "users.csv", "the file where users who have retweeted the tweet at stored. This will be created if it does not exists")
	flag.StringVar(&tweetID, "tweet", "", "The ID of the Tweet you wish to find retweeters of")
	flag.Parse()

	key, secret, err := keys(keyFile)
	if err != nil {
		panic(err)
	}
	client, err := twitterClient(key, secret)
	if err != nil {
		panic(err)
	}
	newUsernames, err := retweeters(client, tweetID)
	if err != nil {
		panic(err)
	}
	existUsernames := existing(usersFile)
	allUsernames := merge(newUsernames, existUsernames)

	err = writeUsers(usersFile, allUsernames)
	if err != nil {
		panic(err)
	}

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

func keys(keyFile string) (key, secret string, err error) {
	var keys struct {
		Key    string `json:"consumer_key"`
		Secret string `json:"consumer_secret"`
	}
	f, err := os.Open(keyFile)
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	dec.Decode(&keys)
	return keys.Key, keys.Secret, nil
}

func retweeters(client *http.Client, tweetID string) ([]string, error) {
	url := fmt.Sprintf("https://api.x.com/1.1/statuses/retweets/%s.json", tweetID)
	res, err := client.Get(url)
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

func existing(usersFile string) []string {
	f, err := os.Open(usersFile)
	if err != nil {
		return []string{}
	}
	defer f.Close()
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	fmt.Println(lines)
	users := make([]string, 0, len(lines))
	for _, line := range lines {
		users = append(users, line[0])
	}
	return users
}

func merge(a, b []string) []string {
	uniq := make(map[string]struct{}, 0)
	for _, user := range a {
		uniq[user] = struct{}{}
	}
	for _, user := range b {
		uniq[user] = struct{}{}
	}
	ret := make([]string, 0, len(uniq))
	for user := range uniq {
		ret = append(ret, user)
	}
	return ret
}

func writeUsers(usersFile string, users []string) error {
	f, err := os.OpenFile(usersFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	for _, username := range users {
		if err := w.Write([]string{username}); err != nil {
			return err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	return nil
}