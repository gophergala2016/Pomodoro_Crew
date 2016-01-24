package server
import (
	"github.com/googollee/go-socket.io"
	"github.com/gophergala2016/Pomodoro_Crew/session"
	"log"
	"github.com/gophergala2016/Pomodoro_Crew/models"
	"time"
	"github.com/google/cayley"
	"strconv"
)

const (
	RoomName = "common"
)

type SocketServer struct {
	*socketio.Server
	session *session.Session
	storage *models.Storage
}

func (s *SocketServer) SetSession(session *session.Session) {
	s.session = session
}

func (s *SocketServer) NotifyStop(t int64) {
	freeAt := strconv.FormatInt(t, 10)
	p := cayley.StartPath(s.storage, freeAt).In("free at")
	it := p.BuildIterator()
	for cayley.RawNext(it) {
		s.BroadcastTo(RoomName, "stop", s.storage.NameOf(it.Result()))
	}
}

func (s *SocketServer) initTimer() {
	ticker := time.NewTicker(time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				go func() {
					s.NotifyStop(time.Now().Unix())
				}()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func NewServer(storage *models.Storage) (*SocketServer, error) {
	s, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}

	server := &SocketServer{s, nil, storage}

	server.On("connection", func(so socketio.Socket) {
		so.Join(RoomName)
		so.On("start", func() {
			if server.session == nil || !server.session.IsAuthorized {
				return
			}
			user := models.NewUser(server.session.Username)
			user.Start(models.Iterate30Minutes)
			so.BroadcastTo(RoomName, "start", user.Name)
		})
		so.On("stop", func() {
			if server.session == nil || !server.session.IsAuthorized {
				return
			}
			user := models.NewUser(server.session.Username)
			user.Start(models.Iterate30Minutes)
			so.BroadcastTo(RoomName, "stop", user.Name)
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	server.initTimer()

	return server, nil
}
