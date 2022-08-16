package players

import (
	"context"
	"encoding/json"

	"github.com/mpraski/clusters"
	"github.com/valyala/fasthttp"
	packagetypes "github.com/vidur2/fantasyAutomation/packageTypes"
	"github.com/vidur2/fantasyAutomation/prisma/db"
	"github.com/vidur2/fantasyAutomation/util"
)

func FetchPlayers() error {
	var parsedPlayers []packagetypes.PlayerRes

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	endpoint := util.LinkMap["players"]

	req.SetRequestURI(endpoint)

	util.Client.Do(req, res)

	err := json.Unmarshal(res.Body(), &parsedPlayers)

	if err != nil {
		return err
	}

	qb := make([][]float64, 0)
	rb := make([][]float64, 0)
	wr := make([][]float64, 0)
	te := make([][]float64, 0)

	playerMap := map[string][][]float64{
		"QB": qb,
		"RB": rb,
		"WR": wr,
		"TE": te,
	}

	playerNamesMap := map[string][]int{
		"QB": make([]int, 0),
		"RB": make([]int, 0),
		"WR": make([]int, 0),
		"TE": make([]int, 0),
	}

	for _, player := range parsedPlayers {
		if player.Pos != "K" && player.Pos != "DEF" {
			stats, err := scrapePlayers(player)
			stats["depth_chart_position_adj"] = 100 / float32(player.DepthChartPos)

			if err != nil {
				continue
			}

			playerArr := make([]float64, len(stats))

			if player.Pos == "QB" {
				idx := 0
				for stat, wgt := range util.PASS_STATS {
					playerArr[idx] = float64(stats[stat]) * float64(wgt)
					idx++
				}
			} else {
				idx := 0
				for stat, wgt := range util.RUSH_RECV_STATS {
					playerArr[idx] = float64(stats[stat]) * float64(wgt)
					idx++
				}
			}

			playerMap[string(player.Pos)] = append(playerMap[string(player.Pos)], playerArr)
			playerNamesMap[string(player.Pos)] = append(playerNamesMap[string(player.Pos)], int(player.Id))
		}
	}

	client := db.NewClient()

	err = client.Connect()

	if err != nil {
		return err
	}

	users, err := client.User.FindMany().Exec(context.Background())

	if err != nil {
		return err
	}

	exploredLeagues := make(map[int]bool)

	go func() {
		client2 := db.NewClient()

		err := client.Connect()

		if err != nil {
			return
		}

		for _, user := range users {
			exploredLeagues, _ = RelinkPlayers(client2, user.FantasyOwnerID, exploredLeagues)
		}
	}()

	for pos, data := range playerMap {
		c, err := clusters.KMeans(1000, 8, clusters.EuclideanDistance)
		if err != nil {
			return err
		}

		err = c.Learn(data)

		if err != nil {
			return err
		}

		for idx, guess := range c.Guesses() {
			playerId := playerNamesMap[pos][idx]
			_, err := client.Players.UpsertOne(db.Players.ID.Equals(playerId)).Update(db.Players.Category.Set(guess)).Exec(context.Background())

			if err != nil {
				_, err = client.Players.CreateOne(db.Players.ID.Set(playerId), db.Players.Position.Set(db.Position(pos)), db.Players.Category.Set(guess)).Exec(context.Background())

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
