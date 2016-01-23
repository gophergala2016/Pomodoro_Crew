package models

import (
	"time"
)

type User struct {
	Id              string
	Name            string
	Password        string
	Teams           []string
	CurrentIterTime *time.Time
	Iteration       bool
}

func NewUser(id, name, password string, teams []string, iterationTime *time.Time, iteration bool) *User {
	return &User{id, name, password, teams, iterationTime, iteration}
}
