package server
import (
	"github.com/googollee/go-socket.io"
	"github.com/gophergala2016/Pomodoro_Crew/session"
	"log"
	"github.com/gophergala2016/Pomodoro_Crew/models"
)

const (
	RoomName = "common"
)

type SocketServer struct {
	*socketio.Server
	session *session.Session
}

func (s *SocketServer) SetSession(session *session.Session) {
	s.session = session
}

func NewServer() (*SocketServer, error) {
	s, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}

	server := &SocketServer{s, nil}

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



	return server, nil
}
