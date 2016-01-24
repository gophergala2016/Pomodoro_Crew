package models

import (
	"time"
	"github.com/google/cayley"
	"strconv"
)

type User struct {
	Name	string
	storage *Storage
}

func NewUser(name string) *User {
	return &User{name}
}

func (u *User) Id() string {
	return u.Name
}

func (u *User) Iteration() bool {
	return u.IterationTime() > time.Now().Unix()
}

func (u *User) IterationTime() int64 {
	p := cayley.StartPath(u.getStorage(), u.Name).Out("free at")

	it := p.BuildIterator()
	if cayley.RawNext(it) {
		return strconv.ParseInt(u.getStorage().NameOf(it.Result()), 10, 64)
	} else {
		return time.Now().Unix()
	}
}

func (u *User) getStorage() *Storage {
	if u.storage == nil {
		u.storage = GetStorage()
	}

	return u.storage
}
