package models
import (
	"github.com/google/cayley/graph"
	"github.com/google/cayley"
	"strconv"
)

const BoltPath = "/tmp/pc"

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

func (s *Storage) GetUsersFreeAt(t *int64) {
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
	if s == nil {
		graph.InitQuadStore("bolt", BoltPath, nil)
		s, err = cayley.NewGraph("bolt", BoltPath, nil)
		storage = s
	} else {
		s = storage
	}

	return
}
