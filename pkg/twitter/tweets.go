/**
URY Tweet Board

Author: Michael Grace <michael.grace@ury.org.uk>
*/

package twitter

import (
	"os"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func GetTweetStream(hashtags []string, tweets chan<- *twitter.Tweet) {

	/**
	I just want to let the world know, I spent absolutely
	ages trying to work out why it was fine when I gave the
	strings to it directly but didn't work when I passed them
	through environment variables. It took for absolutely ever,
	and it just turns out they had carriage returns at the end
	of them.
	*/

	consumerKey := strings.TrimSuffix(os.Getenv("TWITTER_CONSUMER_KEY"), "\r")
	consumerSecret := strings.TrimSuffix(os.Getenv("TWITTER_CONSUMER_SECRET"), "\r")
	oauthToken := strings.TrimSuffix(os.Getenv("TWITTER_OAUTH_TOKEN"), "\r")
	oauthSecret := strings.TrimSuffix(os.Getenv("TWITTER_OAUTH_SECRET"), "\r")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(oauthToken, oauthSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	params := &twitter.StreamFilterParams{
		Track: hashtags,
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
