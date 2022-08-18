package packagetypes

import "github.com/vidur2/fantasyAutomation/prisma/db"

type PlayerRes struct {
	Name          string      `json:"name"`
	Pos           db.Position `json:"position"`
	Team          string      `json:"team"`
	Id            uint        `json:"player_id"`
	DepthChartPos uint        `json:"depth_chart_position"`
}

type RosterRes struct {
	Players  []string `json:"players"`
	Owner    string   `json:"owner_id"`
	RosterId int      `json:"roster_id"`
}

type UserRes struct {
	Username    string `json:"username"`
	UserId      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Avatar      string `json:"avatar"`
}

type LeagueRes struct {
	LeagueId string `json:"league_id"`
}

type NflStateRes struct {
	Week int `json:"week"`
}

type TransactionsRes struct {
	Players []Players `json:"players"`
}

type Players struct {
	PlayerId        string `json:"player_id"`
	PreviousOwnerId int    `json:"previous_owner_id"`
	OwnerId         int    `json:"owner_id"`
}
