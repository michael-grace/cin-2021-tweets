/**
URY Tweet Board

Author: Michael Grace <michael.grace@ury.org.uk>
*/

package main

import (
	t "github.com/dghubble/go-twitter/twitter"
	"github.com/michael-grace/cin-2021-tweets/pkg/twitter"
	"github.com/michael-grace/cin-2021-tweets/pkg/web"
)

func main() {
	tweets := make(chan *t.Tweet)

	go twitter.GetTweetStream(tweets)
	web.StartWebServer(tweets)
}
