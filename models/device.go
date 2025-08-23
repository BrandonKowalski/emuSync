package models

type Device struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Model    string `json:"model"`
}
