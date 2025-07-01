package models

type Camera struct {
	ID       int    `json:"id"`
	IP       string `json:"ip"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
