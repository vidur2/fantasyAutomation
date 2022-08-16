package packagetypes

type RosterValue struct {
	RbVal PosVal
	WrVal PosVal
	QbVal PosVal
	TeVal PosVal
}

type PosVal struct {
	Value     int
	PlayerIds []int
}
