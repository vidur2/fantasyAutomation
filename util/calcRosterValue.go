package util

import (
	"context"

	packagetypes "github.com/vidur2/fantasyAutomation/packageTypes"
	"github.com/vidur2/fantasyAutomation/prisma/db"
)

func CalcRosterVal(leagueId int, client *db.PrismaClient) ([]packagetypes.RosterValue, error) {
	league, err := client.League.FindUnique(db.League.LeagueID.Equals(leagueId)).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	retVal := make([]packagetypes.RosterValue, len(league.Rosters()))

	for _, roster := range league.Rosters() {
		rosterNew := new(packagetypes.RosterValue)
		for _, player := range roster.Player() {
			switch player.Position {
			case "QB":
				rosterNew.QbVal.Value += player.Category / 4
				rosterNew.QbVal.PlayerIds = append(rosterNew.RbVal.PlayerIds, player)
			case "RB":
				rosterNew.RbVal.Value += player.Category
				rosterNew.RbVal.PlayerIds = append(rosterNew.RbVal.PlayerIds, player)
			case "WR":
				rosterNew.WrVal.Value += player.Category
				rosterNew.WrVal.PlayerIds = append(rosterNew.RbVal.PlayerIds, player)
			case "TE":
				rosterNew.TeVal.Value += player.Category
				rosterNew.TeVal.PlayerIds = append(rosterNew.RbVal.PlayerIds, player)
			}
		}
	}

	return retVal, nil
}
