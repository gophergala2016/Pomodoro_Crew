package routes

import (
	"fmt"

	"github.com/gophergala2016/Pomodoro_Crew/session"

	"github.com/martini-contrib/render"
	"github.com/gophergala2016/Pomodoro_Crew/models"
)

type Home struct {
	User *models.User
	IsAuthorized bool
}

func IndexHandler(rnd render.Render, s *session.Session) {
	fmt.Println(s.Username)
	fmt.Println(s.Id)
	user := models.NewUser(s.Username)
	rnd.HTML(200, "index", Home{User: user, IsAuthorized:s.IsAuthorized})
}
