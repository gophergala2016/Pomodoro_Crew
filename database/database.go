package models

import (
	"github.com/google/cayley"
	"github.com/google/cayley/graph"
	_ "github.com/google/cayley/graph/bolt"
	"strconv"
)

const BoltPath = "/tmp/pomodoro_crew"

type Storage struct {
	*cayley.Handle
}

func (s *Storage) SaveUser(u *User) {
	s.AddQuad(cayley.Quad(u.Name, "is", "user", ""))

	p := cayley.StartPath(u.getStorage(), u.Name).Out("free at")

	iterationTime := u.IterationTime()
	it := p.BuildIterator()
	for cayley.RawNext(it) {
		s.RemoveQuad(cayley.Quad(u.Name, "free at", s.NameOf(it.Result()), ""))
	}

	s.AddQuad(cayley.Quad(u.Name, "free at", strconv.FormatInt(iterationTime, 10), ""))
}

func (s *Storage) GetUsers(exclude ...*User) []*User {
	users := []*User{}
	p := cayley.StartPath(s, "user").In("is")
	it := p.BuildIterator()
	for cayley.RawNext(it) {
		name := s.NameOf(it.Result())
		excluded := false
		if len(exclude) > 0 {
			for _, u := range exclude {
				if name == u.Name {
					excluded = true
					break
				}
			}
		}
		if !excluded {
			users = append(users, NewUser(name))
		}
	}

	return users
}

func (s *Storage) GetUsersFreeAt(t int64) []*User {
	freeAt := strconv.FormatInt(t, 10)

	users := []*User{}
	p := cayley.StartPath(s, freeAt).In("free at")
	it := p.BuildIterator()
	for cayley.RawNext(it) {
		users = append(users, NewUser(s.NameOf(it.Result())))
	}

	return users
}

var storage *Storage

func GetStorage() (s *Storage, err error) {
	if storage == nil {
		graph.InitQuadStore("bolt", BoltPath, nil)
		var handle *cayley.Handle
		handle, err = cayley.NewGraph("bolt", BoltPath, nil)
		s = &Storage{handle}
		storage = s
	} else {
		s = storage
	}

	return s, err
}
