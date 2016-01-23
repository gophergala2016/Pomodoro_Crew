package routes

import (
	"net/http"

	//"github.com/gophergala2016/Pomodoro_Crew/models"
	"github.com/gophergala2016/Pomodoro_Crew/session"
	//"github.com/gophergala2016/Pomodoro_Crew/utils"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	//"github.com/boltdb/bolt"
)

func ViewHandler(s *session.Session, rnd render.Render, r *http.Request, params martini.Params) {

	rnd.HTML(200, "view", nil)
}

func GetHtmlHandler(rnd render.Render, r *http.Request) {

	rnd.JSON(200, nil)
}
