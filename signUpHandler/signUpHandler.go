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

func HandleSignUp(r packagetypes.SignUpReq) error {
	prismaCtx := context.Background()
	client := db.NewClient()
	err := client.Connect()

	if err != nil {
		return err
	}

	fantasyOwnerId, err := getUserInformation(r.Username)

	if err != nil {
		client.Disconnect()
		return err
	}

	_, err = client.User.CreateOne(db.User.Email.Set(r.Email), db.User.Password.Set(r.Password), db.User.FantasyOwnerID.Set(int(fantasyOwnerId))).Exec(prismaCtx)

	if err != nil {
		client.Disconnect()
		return err
	}

	err = iterateLeagues(int(fantasyOwnerId), client, prismaCtx)

	client.Disconnect()
	if err != nil {
		return err
	}

	return nil
}

func getUserInformation(username string) (int64, error) {
	reqEndpoint := fmt.Sprintf(util.LinkMap["player_id"], username)
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	req.SetRequestURI(reqEndpoint)

	err := util.Client.Do(req, res)

	if err != nil {
		return 0, err
	}

	var parsed packagetypes.UserRes
	json.Unmarshal(res.Body(), &parsed)

	return strconv.ParseInt(parsed.UserId, 10, 32)
}

func iterateLeagues(userId int, client *db.PrismaClient, ctx context.Context) error {
	reqEndpoint := fmt.Sprintf(util.LinkMap["league"], userId)

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	req.SetRequestURI(reqEndpoint)

	util.Client.Do(req, res)
	var parsed []packagetypes.LeagueRes
	err := json.Unmarshal(res.Body(), &parsed)

	if err != nil {
		return err
	}

	for _, league := range parsed {
		parsedId, _ := strconv.ParseUint(league.LeagueId, 10, 32)
		err := handleLeague(uint(parsedId), uint64(userId), client, ctx)

		if err != nil {
			return err
		}
	}

	return nil
}
