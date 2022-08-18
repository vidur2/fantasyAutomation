package rosteranalysis

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

func ScoreTransactions(league_id int, client *db.PrismaClient) error {
	lim, err := getLeagueState()

	if err != nil {
		return err
	}

	leagueMap, err := getLeagueMap(league_id, client)

	if err != nil {
		return err
	}

	playerSums := make(map[int]int)

	for i := 0; i <= lim; i++ {
		uri := fmt.Sprintf(util.LinkMap["transactions"], league_id, i)
		req := fasthttp.AcquireRequest()
		res := fasthttp.AcquireResponse()

		req.SetRequestURI(uri)

		err := util.Client.Do(req, res)

		if err != nil {
			continue
		}

		var transactionRes []packagetypes.TransactionsRes

		req = fasthttp.AcquireRequest()
		res = fasthttp.AcquireResponse()

		err = util.Client.Do(req, res)

		if err != nil {
			continue
		}

		err = json.Unmarshal(res.Body(), &transactionRes)

		if err != nil {
			continue
		}

		for _, transaction := range transactionRes {
			for _, player := range transaction.Players {
				playerInt, err := strconv.ParseInt(player.PlayerId, 10, 32)

				if err != nil {
					break
				}

				playerFetched, err := client.Players.FindFirst(db.Players.ID.Equals(int(playerInt))).Exec(context.Background())

				if err != nil {
					break
				}

				rtg := playerFetched.Category
				playerSums[player.OwnerId] += rtg
				playerSums[player.PreviousOwnerId] -= rtg
			}
		}
	}

	leagueModel, err := client.League.FindFirst(db.League.LeagueID.Equals(league_id)).Exec(context.Background())

	if err != nil {
		return err
	}

	for _, roster := range leagueModel.Rosters() {
		rosterId := leagueMap[roster.UserID]
		score := playerSums[rosterId]
		client.Roster.UpsertOne(db.Roster.RosterID.Equals(roster.RosterID)).Update(db.Roster.RosterScore.Set(score)).Exec(context.Background())
	}

	return nil
}

func getLeagueState() (int, error) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	var state packagetypes.NflStateRes

	req.SetRequestURI("https://api.sleeper.app/v1/state/nfl")

	err := util.Client.Do(req, res)

	if err != nil {
		return -1, err
	}

	err = json.Unmarshal(res.Body(), &state)

	if err != nil {
		return -1, err
	}

	return state.Week, nil
}

func getLeagueMap(league_id int, client *db.PrismaClient) (map[int]int, error) {
	leagueInfo, err := client.League.FindFirst(db.League.LeagueID.Equals(league_id)).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	var parsedLi map[int]int

	mapString, valid := leagueInfo.RosterMap()

	if !valid {
		return nil, fmt.Errorf("prisma fetch error")
	}

	err = json.Unmarshal([]byte(mapString), &parsedLi)

	if err != nil {
		return nil, err
	}

	return parsedLi, nil
}
