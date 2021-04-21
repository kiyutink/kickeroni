package templates

import (
	"fmt"

	"github.com/slack-go/slack"
)

// Starting template
func InProgress(t Teamer) []slack.Block {
	team0, team1 := t.Teams()
	mainText := slack.NewTextBlockObject("mrkdwn",
		"Game in progress",
		false,
		false)
	mainSection := slack.NewSectionBlock(mainText, nil, nil)

	team0Text := slack.NewTextBlockObject(
		"mrkdwn",
		fmt.Sprintf("Team0: <@%s>, <@%s>", team0[0], team0[1]),
		false,
		false,
	)
	team0Section := slack.NewSectionBlock(team0Text, nil, nil)

	team1Text := slack.NewTextBlockObject(
		"mrkdwn",
		fmt.Sprintf("Team1: <@%s>, <@%s>", team1[0], team1[1]),
		false,
		false,
	)
	team1Section := slack.NewSectionBlock(team1Text, nil, nil)

	team0WinButton := slack.NewButtonBlockElement("win_team0", "", slack.NewTextBlockObject("plain_text", "Team 0 won", false, false))
	team0WinButton.Style = slack.StylePrimary

	team1WinButton := slack.NewButtonBlockElement("win_team1", "", slack.NewTextBlockObject("plain_text", "Team 1 won", false, false))
	team1WinButton.Style = slack.StylePrimary

	actions := slack.NewActionBlock("", team0WinButton, team1WinButton)

	cancelButton := slack.NewButtonBlockElement("cancel", "cancel", slack.NewTextBlockObject("plain_text", "Cancel game", false, false))
	cancelButton.Style = slack.StyleDanger
	cancelButtonSection := slack.NewActionBlock("", cancelButton)

	return []slack.Block{mainSection, team0Section, team1Section, actions, cancelButtonSection}
}
