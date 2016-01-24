package routes

import (
	"github.com/gophergala2016/Pomodoro_Crew/session"

	"github.com/martini-contrib/render"
	"github.com/gophergala2016/Pomodoro_Crew/models"
)

type Home struct {
	User *models.User
}

func IndexHandler(rnd render.Render, s *session.Session) {
	user := models.NewUser(s.Username)
	rnd.HTML(200, "index", Home{User: user})
}
