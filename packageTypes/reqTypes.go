package packagetypes

type SignUpReq struct {
	Email    string
	Password []byte
	UserId   []uint
	RosterId uint
}
