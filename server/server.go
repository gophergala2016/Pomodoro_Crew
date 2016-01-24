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
	stopAt := strconv.FormatInt(t, 10)
	p := cayley.StartPath(s.storage, stopAt).In("free at")
	it := p.BuildIterator()
	for cayley.RawNext(it) {
		s.BroadcastTo(RoomName, "stop", s.storage.NameOf(it.Result()))
	}
}

func (s *SocketServer) NotifyEnable(t int64) {
	stopedAt := strconv.FormatInt(t - models.Wait5Minutes, 10)
	p := cayley.StartPath(s.storage, stopedAt).In("free at")
	it := p.BuildIterator()
	for cayley.RawNext(it) {
		//TODO send to only one user
		s.BroadcastTo(RoomName, "enable", s.storage.NameOf(it.Result()))
	}
}

func (s *SocketServer) initTimer() {
	ticker := time.NewTicker(time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				t := time.Now().Unix()
				go func() {
					s.NotifyStop(t)
				}()
				go func() {
					s.NotifyEnable(t)
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
			user.Stop()
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
