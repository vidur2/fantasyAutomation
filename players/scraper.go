package players

import (
	"strconv"

	"github.com/anaskhan96/soup"
	"github.com/valyala/fasthttp"
	packagetypes "github.com/vidur2/fantasyAutomation/packageTypes"
	"github.com/vidur2/fantasyAutomation/util"
)

var stats = map[string]string{
	"QB": "https://www.pro-football-reference.com/years/2022/passing.htm",
	"RB": "https://www.pro-football-reference.com/years/2022/rushing.htm",
	"WR": "https://www.pro-football-reference.com/years/2022/recieving.htm",
	"TE": "https://www.pro-football-reference.com/years/2022/recieving.htm",
}

var players = make(map[string]string)

const BASE_URI = "https://www.pro-football-reference.com/"

func scrapePlayers(player packagetypes.PlayerRes) (map[string]float32, error) {
	playerName := player.Name

	href := players[player.Name]

	var playerFinal map[string]float32

	if href != "" {
		playerDoc, err := getRoot(href)

		if err != nil {
			return nil, err
		}

		playerFinal = getPlayerData(playerDoc)

	} else {
		posLink, err := getRoot(stats[players[string(player.Pos)]])

		if err != nil {
			return nil, err
		}

		for _, link := range posLink.FindAll("a") {
			path := link.Attrs()["href"]
			txt := link.Text()

			if txt == playerName {
				playerDoc, err := getRoot(BASE_URI + path + "/gamelog/2022/advanced")

				if err != nil {
					return nil, err
				}

				playerFinal = getPlayerData(playerDoc)

				break
			} else {
				players[player.Name] = BASE_URI + path + "/gamelog/2022/advanced"
			}
		}
	}

	return playerFinal, nil
}

func getRoot(link string) (soup.Root, error) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	req.SetRequestURI(link)

	err := util.Client.Do(req, res)

	if err != nil {
		return soup.Root{}, err
	}

	doc := soup.HTMLParse(res.String())

	return doc, nil
}

func getPlayerData(playerDoc soup.Root) map[string]float32 {
	table := playerDoc.Find("table")
	tot := table.FindAll("tr")
	retVal := make(map[string]float32)
	for wkNum, row := range tot {
		adjWkNm := wkNum + 1
		for _, data := range row.FindAll("td") {
			cat := data.Attrs()["data-stat"]

			stat, err := strconv.ParseFloat(data.Text(), 32)

			if err != nil {
				continue
			}

			avgStat := (retVal[cat]*float32(wkNum) + float32(stat)) / float32(adjWkNm)
			retVal[cat] = avgStat
		}
	}

	return retVal
}
