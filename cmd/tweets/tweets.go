/**
URY Tweet Board
Copyright (C) 2022 Michael Grace <michael.grace@ury.org.uk>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
