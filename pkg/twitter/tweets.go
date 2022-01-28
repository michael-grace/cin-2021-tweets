/**
URY Tweet Board

Author: Michael Grace <michael.grace@ury.org.uk>
*/

package twitter

import (
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func GetTweetStream(tweets chan<- *twitter.Tweet) {
	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_OAUTH_TOKEN"), os.Getenv("TWITTER_OAUTH_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	params := &twitter.StreamFilterParams{
		Track: []string{os.Getenv("HASHTAG")},
	}

	stream, err := client.Streams.Filter(params)

	if err != nil {
		panic(err)
	}

	defer stream.Stop()

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		tweets <- tweet
	}

	demux.HandleChan(stream.Messages)

}
