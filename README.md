# Tweet Board

###### Michael Grace, 2020, 2022

## Starting

Required environment variables:

-   `TWITTER_CONSUMER_KEY`
-   `TWITTER_CONSUMER_SECRET`
-   `TWITTER_OAUTH_TOKEN`
-   `TWITTER_OAUTH_SECRET`
-   `HASHTAG` (this must include the `#` character, and can be multiple hashtags, separated by `,`)
-   `AUTH_USER` - username for the the controller login
-   `AUTH_PASS` - password for the controller login

Build and run the Dockerfile, exposing port 3000.

## Accessing

-   The tweet wall can be accessed at `/`.
-   The controller can be accessed at `/control`. This requires the login defined by the above environment variables.
