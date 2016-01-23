package main

import (
	"fmt"
	"html/template"

	"github.com/gophergala2016/Pomodoro_Crew/routes"
	"github.com/gophergala2016/Pomodoro_Crew/session"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/boltdb/bolt"
)

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {
	fmt.Println("Listening on port :3000")

	db, err := bolt.Open("pomororo_crew.db", 0600, nil)
	if err != nil {
		panic(err)
	}

	m := martini.Classic()

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

	m.Map(db)

	m.Use(session.Middleware)

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                         // Specify what path to load the templates from.
		Layout:     "layout",                            // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"},          // Specify extensions to load for templates.
		Funcs:      []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
		Charset:    "UTF-8",                             // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                                // Output human readable JSON
	}))

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))
	m.Get("/", routes.IndexHandler)
	m.Get("/login", routes.GetLoginHandler)
	m.Get("/logout", routes.LogoutHandler)
	m.Post("/login", routes.PostLoginHandler)
	m.Get("/view:id", routes.ViewHandler)
	m.Post("/gethtml", routes.GetHtmlHandler)

	m.Run()
}
