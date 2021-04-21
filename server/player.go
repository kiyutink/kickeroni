package server

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/kiyutink/kickeroni/db"
)

const DEFAULT_RANK = 1600

type Player struct {
	Id     string `json:"id"`
	Rank   int    `json:"rank"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
}

func NewPlayer(id string) *Player {
	return &Player{
		Id:     id,
		Rank:   DEFAULT_RANK,
		Wins:   0,
		Losses: 0,
	}
}

func (p *Player) Save() error {
	return db.ReplaceRecord(db.PlayersTable, p.Id, p)
}

func (p *Player) AdjustRank(delta int) error {
	p.Rank += delta
	return p.Save()
}

func GetPlayerById(id string) (*Player, error) {
	p := NewPlayer(id)
	found, err := db.GetFirstRecordByFieldValue(db.PlayersTable, "id", id, p)

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("there's no player with id %v", id)
	}
	return p, err
}

func ListPlayers() ([]*Player, error) {
	players := []*Player{}
	res, err := db.GetAllRecords(db.PlayersTable)
	if err != nil {
		return players, err
	}
	for _, p := range res {
		player := &Player{}
		if err := dynamodbattribute.UnmarshalMap(p, player); err != nil {
			return players, err
		}
		players = append(players, player)
	}
	return players, nil
}
