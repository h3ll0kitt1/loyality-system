package domain

type CtxLoginKey struct{}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
