package routes

import (
	"fmt"

	"github.com/gophergala2016/Pomodoro_Crew/session"

	"github.com/martini-contrib/render"
)

func IndexHandler(rnd render.Render, s *session.Session) {
	fmt.Println(s.Username)
	fmt.Println(s.Id)
	rnd.HTML(200, "index", nil)
}
