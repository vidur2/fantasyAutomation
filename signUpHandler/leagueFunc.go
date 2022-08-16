package signuphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"
	packagetypes "github.com/vidur2/fantasyAutomation/packageTypes"
	db "github.com/vidur2/fantasyAutomation/prisma/db"
	"github.com/vidur2/fantasyAutomation/util"
)

func handleLeague(leagueId uint, userId uint64, client *db.PrismaClient, ctx context.Context) error {
	reqEndpoint := fmt.Sprintf(util.LinkMap["leagues"], leagueId)
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	req.SetRequestURI(reqEndpoint)

	err := util.Client.Do(req, res)

	if err != nil {
		return err
	}
	var parsedRes []packagetypes.RosterRes
	err = json.Unmarshal(res.Body(), &parsedRes)

	if err != nil {
		return err
	}

	client.League.CreateOne(db.League.LeagueID.Set(int(leagueId)))

	linkParams := make([]db.RosterWhereParam, len(parsedRes))

	for _, roster := range parsedRes {
		intOwner, err := strconv.ParseUint(roster.Owner, 10, 32)

		if err != nil {
			return err
		}

		linkClauses := make([]db.PlayersWhereParam, len(roster.Players))

		for idx, player := range roster.Players {
			intPlayer, err := strconv.ParseInt(player, 10, 32)

			if err != nil {
				return err
			}

			linkClauses[idx] = db.Players.ID.Equals(int(intPlayer))
		}

		roster, err := client.Roster.CreateOne(db.Roster.Player.Link(db.Players.Or(linkClauses...)), db.Roster.LeagueLeagueID.Set(int(leagueId))).Exec(ctx)

		if err != nil {
			return err
		}

		if intOwner == userId {
			linkParams = append(linkParams, db.Roster.RosterID.Equals(roster.RosterID))
		}
	}

	client.User.UpsertOne(db.User.FantasyOwnerID.Equals(int(userId))).Update(db.User.Rosters.Link(linkParams...), db.User.League.Link(db.League.LeagueID.Equals(int(leagueId)))).Exec(ctx)

	return nil
}
