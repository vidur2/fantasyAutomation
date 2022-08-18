package packagetypes

import "github.com/vidur2/fantasyAutomation/prisma/db"

type RosterValue struct {
	RbVal PosVal
	WrVal PosVal
	QbVal PosVal
	TeVal PosVal
}

type PosVal struct {
	Value     int
	PlayerIds []db.PlayersModel
}
