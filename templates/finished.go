package templates

import (
	"fmt"

	"github.com/slack-go/slack"
)

func Finished(winners, losers *[2]string) []slack.Block {
	text := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf(
		"Game over: <@%v> and <@%v> defeat <@%v> and <@%v>",
		winners[0],
		winners[1],
		losers[0],
		losers[1],
	),
		false,
		false,
	)
	section := slack.NewSectionBlock(text, nil, nil)

	return []slack.Block{section}
}
