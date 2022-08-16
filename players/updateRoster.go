package players

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"
	packagetypes "github.com/vidur2/fantasyAutomation/packageTypes"
	"github.com/vidur2/fantasyAutomation/prisma/db"
	"github.com/vidur2/fantasyAutomation/util"
)

func RelinkPlayers(client *db.PrismaClient, userId int, exploredLeagues map[int]bool) (map[int]bool, error) {
	userInfo, err := client.User.FindFirst(db.User.FantasyOwnerID.Equals(userId)).Exec(context.Background())

	if err != nil {
		return exploredLeagues, err
	}

	for _, league := range userInfo.League() {
		id := league.LeagueID

		if !exploredLeagues[id] {
			link := fmt.Sprintf(util.LinkMap["rosters"], id)

			var resBody []packagetypes.RosterRes

			req := fasthttp.AcquireRequest()
			res := fasthttp.AcquireResponse()

			req.SetRequestURI(link)

			err := util.Client.Do(req, res)

			if err != nil {
				return exploredLeagues, err
			}

			err = json.Unmarshal(res.Body(), &resBody)

			if err != nil {
				return exploredLeagues, err
			}

			for _, roster := range resBody {
				userId, err := strconv.ParseInt(roster.Owner, 10, 32)
				if err != nil {
					continue
				}

				prevRoster, err := client.Roster.FindFirst(db.Roster.UserID.Equals(int(userId)), db.Roster.LeagueLeagueID.Equals(id)).Exec(context.Background())

				if err != nil {
					continue
				}

				unlinkClauses := make([]db.PlayersWhereParam, 0)
				for _, player := range prevRoster.Player() {
					unlinkClauses = append(unlinkClauses, db.Players.ID.Equals(player.ID))
				}

				_, err = client.Roster.UpsertOne(db.Roster.RosterID.Equals(prevRoster.RosterID)).Update(db.Roster.Player.Unlink(unlinkClauses...)).Exec(context.Background())

				if err != nil {
					return exploredLeagues, err
				}

				linkClauses := make([]db.PlayersWhereParam, 0)
				for _, player := range roster.Players {
					intPid, err := strconv.ParseInt(player, 10, 32)

					if err != nil {
						continue
					}

					linkClauses = append(linkClauses, db.Players.ID.Equals(int(intPid)))
				}

				_, err = client.Roster.UpsertOne(db.Roster.RosterID.Equals(prevRoster.RosterID)).Update(db.Roster.Player.Link(linkClauses...)).Exec(context.Background())

				if err != nil {
					return exploredLeagues, err
				}
			}

			exploredLeagues[id] = true
		}
	}

	return exploredLeagues, nil
}
