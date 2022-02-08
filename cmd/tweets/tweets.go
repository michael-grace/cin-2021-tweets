/**
URY Tweet Board

Author: Michael Grace <michael.grace@ury.org.uk>
*/

package main

import (
	"os"
	"strings"

	t "github.com/dghubble/go-twitter/twitter"
	"github.com/michael-grace/cin-2021-tweets/pkg/twitter"
	"github.com/michael-grace/cin-2021-tweets/pkg/web"
)

func main() {

	var hashtags []string

	hashes := os.Getenv("HASHTAG")
	hashtags = strings.Split(strings.ReplaceAll(hashes, " ", ""), ",")

	tweets := make(chan *t.Tweet)

	go twitter.GetTweetStream(hashtags, tweets)
	web.StartWebServer(hashtags, tweets)
}
