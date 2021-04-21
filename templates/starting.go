package templates

import (
	"fmt"

	"github.com/slack-go/slack"
)

type Teamer interface {
	Teams() (*[2]string, *[2]string)
}

// Starting template
func Starting(t Teamer) []slack.Block {
	team0, team1 := t.Teams()
	mainText := slack.NewTextBlockObject("mrkdwn",
		"Starting a game, please join!",
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

	joinRandomButton := slack.NewButtonBlockElement("join_random", "", slack.NewTextBlockObject("plain_text", "Join random team", false, false))
	joinRandomButton.Style = slack.StylePrimary

	join0Button := slack.NewButtonBlockElement("join_team0", "", slack.NewTextBlockObject("plain_text", "Join team 0", false, false))
	join0Button.Style = slack.StylePrimary

	join1Button := slack.NewButtonBlockElement("join_team1", "", slack.NewTextBlockObject("plain_text", "Join team 1", false, false))
	join1Button.Style = slack.StylePrimary

	isTeam0Full := true
	isTeam1Full := true
	for _, p := range team0 {
		if p == "" {
			isTeam0Full = false
		}
	}
	for _, p := range team1 {
		if p == "" {
			isTeam1Full = false
		}
	}
	actionSlice := []slack.BlockElement{}
	if !isTeam0Full && !isTeam1Full {
		actionSlice = append(actionSlice, joinRandomButton)
	}

	if !isTeam0Full {
		actionSlice = append(actionSlice, join0Button)
	}

	if !isTeam1Full {
		actionSlice = append(actionSlice, join1Button)
	}

	actions := slack.NewActionBlock("", actionSlice...)

	cancelButton := slack.NewButtonBlockElement("cancel", "cancel", slack.NewTextBlockObject("plain_text", "Cancel game", false, false))
	cancelButton.Style = slack.StyleDanger
	cancelButtonSection := slack.NewActionBlock("", cancelButton)

	return []slack.Block{mainSection, team0Section, team1Section, actions, cancelButtonSection}
}
