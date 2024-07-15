package models

import "time"

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Passport string `json:"passport"`
	Surname  string `json:"surname"`
	Address  string `json:"address"`
	Patr     string `json:"patronomic"`
}

type ApiUser struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Address string `json:"address"`
	Patr    string `json:"patronomic"`
}

type Task struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user-id"`
	StartDate time.Time `json:"start-date"`
	EndDate   time.Time `json:"end-date"`
	LastStart time.Time `json:"last-start"`
	Duration  int       `json:"duration"`
}
