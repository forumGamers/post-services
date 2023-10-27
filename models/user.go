package models

type User struct {
	UUID     string `json:"UUID"`
	LoggedAs string `json:"loggedAs"`
}
