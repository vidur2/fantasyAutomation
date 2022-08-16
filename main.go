package main

import (
	"encoding/json"
	"time"

	"github.com/valyala/fasthttp"
	packagetypes "github.com/vidur2/fantasyAutomation/packageTypes"
	"github.com/vidur2/fantasyAutomation/players"
	signuphandler "github.com/vidur2/fantasyAutomation/signUpHandler"
	"github.com/vidur2/fantasyAutomation/util"
)

func handler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/sign_up":
		var reqBody packagetypes.SignUpReq
		err := json.Unmarshal(ctx.Request.Body(), &reqBody)

		if err != nil {
			handleError(err, ctx)
		}

		err = signuphandler.HandleSignUp(reqBody)

		if err != nil {
			handleError(err, ctx)
		}
	}

}

func handleError(err error, ctx *fasthttp.RequestCtx) {
	body := packagetypes.ErrorType{
		Err: err,
	}

	bodyBytes, _ := json.Marshal(body)
	ctx.Response.AppendBody(bodyBytes)
}

func main() {
	util.InitClient()
	go func() {
		for {
			players.FetchPlayers()
			time.Sleep(24 * time.Hour)
		}
	}()
	fasthttp.ListenAndServe(":8080", handler)
}
