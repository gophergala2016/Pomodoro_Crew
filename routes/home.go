package routes

import (
	"github.com/gophergala2016/Pomodoro_Crew/session"

	"github.com/martini-contrib/render"
	"github.com/gophergala2016/Pomodoro_Crew/models"
	"log"
	"time"
)

type Home struct {
	Current     *models.User
	Users       []*models.User
	SecondsLeft int64
}

func IndexHandler(rnd render.Render, s *session.Session) {
	user := models.NewUser(s.Username)
	storage, err := models.GetStorage()
	if err != nil {
		log.Fatalln(err)
	}

	data := Home{
		Current: user,
		Users:storage.GetUsers(user),
	}

	if user.Iteration() {
		data.SecondsLeft = user.IterationTime() - time.Now().Unix()
	} else if !user.CanStart() {
		data.SecondsLeft = models.Wait5Minutes - (time.Now().Unix() - user.IterationTime())
	}

	rnd.HTML(200, "index", data)
}
