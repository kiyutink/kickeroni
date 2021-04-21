package lib

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/slack-go/slack"
)

var Api = slack.New(Getenv("SLACK_BOT_TOKEN"))
