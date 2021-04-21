package server

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kiyutink/kickeroni/lib"
	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"
)

func Ranks(c echo.Context) error {
	command, err := slack.SlashCommandParse(c.Request())
	if err != nil {
		return err
	}

	players, err := ListPlayers()
	if err != nil {
		return err
	}
	msg := ""
	for _, p := range players {
		msg += fmt.Sprintf("<@%v> - %v - %v wins / %v losses\n", p.Id, p.Rank, p.Wins, p.Losses)
	}

	_, _, err = lib.Api.PostMessage(command.ChannelID, slack.MsgOptionText(msg, false))
	if err != nil {
		return err
	}
	return nil
}

// Slash command /play
func Play(c echo.Context) error {

	command, err := slack.SlashCommandParse(c.Request())
	if err != nil {
		return err
	}

	player, _ := GetPlayerById(command.UserID)

	if player == nil {
		player = NewPlayer(command.UserID)
		player.Save()
	}

	_, err = FindActiveGame()

	if err == nil {
		return errors.New("there's an active game already in progress")
	}

	messageChannelID, messageTimestamp, _ := lib.Api.PostMessage(
		command.ChannelID,
		slack.MsgOptionText("Starting a new game...", false),
	)

	newGame := Game{
		Id:        uuid.NewString(),
		ChannelID: messageChannelID,
		MessageTS: messageTimestamp,
		Team0:     &[2]string{command.UserID},
		Team1:     &[2]string{},
		Winner:    "",
		Status:    "active",
	}

	err = newGame.Render()

	if err != nil {
		return err
	}

	err = newGame.Save()

	return err
}
