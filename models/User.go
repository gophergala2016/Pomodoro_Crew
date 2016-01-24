package models

import (
	"github.com/google/cayley"
	"strconv"
	"time"
)

const (
	Iterate15Minutes = 15 * 60
	Iterate30Minutes = 30 * 60
	Iterate45Minutes = 45 * 60
)

type User struct {
	Name          string
	iterationTime int64
	storage       *Storage
}

func NewUser(name string) *User {
	user := &User{name, 0, nil}
	return user
}

func (u *User) Id() string {
	return u.Name
}

func (u *User) Iteration() bool {
	return u.IterationTime() > time.Now().Unix()
}

func (u *User) IterationTime() int64 {
	if u.iterationTime == 0 {
		p := cayley.StartPath(u.getStorage(), u.Name).Out("free at")

		it := p.BuildIterator()
		if cayley.RawNext(it) {
			u.iterationTime, _ = strconv.ParseInt(u.getStorage().NameOf(it.Result()), 10, 64)
		} else {
			u.iterationTime = time.Now().Unix()
		}
	}

	return u.iterationTime
}

func (u *User) Start(duration int64) {
	u.iterationTime = time.Now().Unix() + duration
	u.getStorage().SaveUser(u)
}

func (u *User) Stop() {
	u.iterationTime = time.Now().Unix()
	u.getStorage().SaveUser(u)
}

func (u *User) getStorage() *Storage {
	if u.storage == nil {
		u.storage, _ = GetStorage()
	}

	return u.storage
}
