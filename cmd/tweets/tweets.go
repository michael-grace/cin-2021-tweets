/**
URY Tweet Board

Author: Michael Grace <michael.grace@ury.org.uk>
*/

package main

import (
	"fmt"
	"os"
	"strings"

	t "github.com/dghubble/go-twitter/twitter"
	"github.com/michael-grace/cin-2021-tweets/pkg/twitter"
	"github.com/michael-grace/cin-2021-tweets/pkg/web"
)

func main() {

	requiredEnvs := []string{
		"HASHTAG",
		"AUTH_USER",
		"AUTH_PASS",
		"TWITTER_CONSUMER_KEY",
		"TWITTER_CONSUMER_SECRET",
		"TWITTER_OAUTH_TOKEN",
		"TWITTER_OAUTH_SECRET",
	}

	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			panic(fmt.Sprintf("%s not set", env))
		}
	}

	var hashtags []string

	hashes := os.Getenv("HASHTAG")
	hashtags = strings.Split(strings.ReplaceAll(hashes, " ", ""), ",")

	tweets := make(chan *t.Tweet)

	go twitter.GetTweetStream(hashtags, tweets)
	web.StartWebServer(hashtags, tweets)
}
