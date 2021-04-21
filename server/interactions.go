package server

import (
	"math/rand"

	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"

	elogo "github.com/kortemy/elo-go"
)

func GetPlayerRankById(id string) (int, error) {
	player, err := GetPlayerById(id)
	if err != nil {
		return 0, err
	}

	return player.Rank, nil
}

func GetTeamAvgRank(players *[2]string) (int, error) {
	sum := 0
	for _, id := range players {
		playerRank, err := GetPlayerRankById(id)
		if err != nil {
			return 0, err
		}
		sum += playerRank
	}
	return sum / 2, nil
}

func AdjustTeamRanksAndCounts(players *[2]string, delta int) error {
	for _, id := range players {
		player, err := GetPlayerById(id)
		if err != nil {
			return err
		}
		err = player.AdjustRank(delta)
		if err != nil {
			return err
		}
		if delta > 0 {
			player.Wins += 1
		} else {
			player.Losses += 1
		}
		player.Save()
	}
	return nil
}

// Interactions is the entry point for all user interactions
func Interactions(c echo.Context) error {
	ic := new(slack.InteractionCallback)
	ic.UnmarshalJSON([]byte(c.FormValue("payload")))
	activeGame, err := FindActiveGame()
	if err != nil {
		return err
	}

	player, _ := GetPlayerById(ic.User.ID)

	if player == nil {
		player = NewPlayer(ic.User.ID)
		player.Save()
	}

	for _, action := range ic.ActionCallback.BlockActions {

		switch action.ActionID {

		case "join_team0", "join_team1":
			teamIndex := uint8(0)
			if action.ActionID == "join_team1" {
				teamIndex = 1
			}
			err = activeGame.AddPlayerToTeam(teamIndex, ic.User.ID)
			if err != nil {
				return err
			}
			activeGame.Save()
			err = activeGame.Render()
			if err != nil {
				return err
			}

		case "join_random":
			teamIndex, otherTeamIndex := uint8(0), uint8(1)

			if rand.Float32() > 0.5 {
				teamIndex, otherTeamIndex = otherTeamIndex, teamIndex
			}

			err = activeGame.AddPlayerToTeam(teamIndex, ic.User.ID)

			if err != nil {
				err = activeGame.AddPlayerToTeam(otherTeamIndex, ic.User.ID)
			}
			if err != nil {
				return err
			}
			activeGame.Save()
			activeGame.Render()

		case "win_team0", "win_team1":
			outcome := 1.0
			winner := "team0"
			if action.ActionID == "win_team1" {
				winner = "team1"
				outcome = 0
			}
			activeGame.Status = "finished"
			activeGame.Winner = winner
			activeGame.Save()
			activeGame.Render()
			team0AvgRank, err := GetTeamAvgRank(activeGame.Team0)
			if err != nil {
				return err
			}
			team1AvgRank, err := GetTeamAvgRank(activeGame.Team1)
			if err != nil {
				return err
			}

			elo := elogo.NewElo()
			delta := elo.RatingDelta(team0AvgRank, team1AvgRank, outcome)
			AdjustTeamRanksAndCounts(activeGame.GetWinnerIds(), delta)
			AdjustTeamRanksAndCounts(activeGame.GetLosersIds(), -delta)

		case "cancel":
			activeGame.Delete()
		}
	}

	return nil
}
