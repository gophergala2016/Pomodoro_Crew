package main

import (
	"html/template"
	"log"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/gophergala2016/Pomodoro_Crew/routes"
	"github.com/gophergala2016/Pomodoro_Crew/server"
	"github.com/gophergala2016/Pomodoro_Crew/session"

	"github.com/google/cayley/graph"
	"net/http"
)

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {

	path := "/tmp/pc"

	graph.InitQuadStore("bolt", path, nil)

	m := martini.Classic()

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

	m.Use(session.Middleware)

	m.Use(render.Renderer(render.Options{
		Directory:  "templates", // Specify what path to load the templates from.
		Layout:     "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		Funcs:      []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
		Charset:    "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true, // Output human readable JSON
	}))

	server, err := server.NewServer()
	if err != nil {
		log.Fatalln(err)
	}

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))
	m.Get("/", routes.IndexHandler)
	m.Get("/login", routes.GetLoginHandler)
	m.Get("/logout", routes.LogoutHandler)
	m.Post("/login", routes.PostLoginHandler)
	m.Get("/view:id", routes.ViewHandler)
	m.Post("/gethtml", routes.GetHtmlHandler)
	m.Get("/socket.io/", func(w http.ResponseWriter, rnd render.Render, r *http.Request, s *session.Session) {
		server.SetSession(s)
		server.ServeHTTP(w, r)
	})
	m.Run()
}
