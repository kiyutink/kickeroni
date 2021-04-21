package server

import (
	"errors"
	"fmt"

	"github.com/kiyutink/kickeroni/db"
	"github.com/kiyutink/kickeroni/lib"
	"github.com/kiyutink/kickeroni/templates"
	"github.com/slack-go/slack"
)

type Game struct {
	Id        string     `json:"id"`
	MessageTS string     `json:"message_ts"`
	ChannelID string     `json:"channel_id"`
	Team0     *[2]string `json:"team0"`
	Team1     *[2]string `json:"team1"`
	Winner    string     `json:"winner"`
	Status    string     `json:"status"`
}

func FindActiveGame() (*Game, error) {
	game := new(Game)

	found, err := db.GetFirstRecordByFieldValue(db.MatchesTable, "status", "active", game)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("there are no active games")
	}
	return game, err
}

func (game *Game) Teams() (*[2]string, *[2]string) {
	return game.Team0, game.Team1
}

func (game *Game) Save() error {
	return db.ReplaceRecord(db.MatchesTable, game.Id, game)
}

func (game *Game) IsPlayerPresent(id string) bool {
	allPlayers := []string{}
	allPlayers = append(allPlayers, game.Team0[:]...)
	allPlayers = append(allPlayers, game.Team1[:]...)
	for _, p := range allPlayers {
		if p == id {
			return true
		}
	}

	return false
}

func (game *Game) AddPlayerToTeam(teamIndex uint8, userID string) error {
	if teamIndex > 1 {
		return errors.New("teamIndex should be either 0 or 1")
	}

	if game.IsPlayerPresent(userID) {
		return nil
	}

	team := game.Team0
	if teamIndex == 1 {
		team = game.Team1
	}

	for i, p := range team {
		if p == "" {
			team[i] = userID
			return nil
		}
	}

	return fmt.Errorf("team is %v full", teamIndex)
}

func (game *Game) Render() error {
	var template []slack.Block
	switch game.Status {
	case "active":
		if game.IsFull() {
			template = templates.InProgress(game)
		} else {
			template = templates.Starting(game)
		}
	case "finished":
		winners := game.Team0
		losers := game.Team1
		if game.Winner == "team1" {
			winners = game.Team1
			losers = game.Team0
		}
		template = templates.Finished(winners, losers)
	}

	_, _, _, err := lib.Api.UpdateMessage(
		game.ChannelID,
		game.MessageTS,
		slack.MsgOptionBlocks(template...),
	)
	return err
}

func (game *Game) Delete() error {
	_, _, err := lib.Api.DeleteMessage(game.ChannelID, game.MessageTS)
	if err != nil {
		return err
	}
	err = db.DeleteRecord(db.MatchesTable, game.Id)

	return err
}

func (game *Game) IsFull() bool {
	allPlayers := []string{}
	allPlayers = append(allPlayers, game.Team0[:]...)
	allPlayers = append(allPlayers, game.Team1[:]...)

	for _, p := range allPlayers {
		if p == "" {
			return false
		}
	}

	return true
}

func (game *Game) GetWinnerIds() *[2]string {
	if game.Winner == "team0" {
		return game.Team0
	}

	return game.Team1
}

func (game *Game) GetLosersIds() *[2]string {
	if game.Winner == "team0" {
		return game.Team1
	}

	return game.Team0
}
