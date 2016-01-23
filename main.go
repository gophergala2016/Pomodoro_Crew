package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
)

func main() {
	fmt.Println("Listening on port :3000")

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
		Charset:    "UTF-8",
		Funcs:      []template.FuncMap{unescapeFuncMap},
		IndentJSON: true,
	}))

	m.Run()
	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))

	m.Get("/", indexHandler)
	m.Get("/login", routes.getLoginHandler)
	m.Post("/", routes.postLoginHandler)

	m.Post("/gethtml", getHtmlHandler)

}

func indexHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {

	rnd.HTML(200, "index", posts)
}

func unescape(x string) interface{} {
	return template.HTML(x)
}

func getHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	html := utils.ConvertMarkdownToHtml(md)

	rnd.JSON(200, map[string]interface{}{"html": string(html)})
}

func getLoginHandler(rnd render.Render) {
	rnd.HTML(200, "login", nil)
}

func postLoginHandler(rnd render.Render, w http.ResponceWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	sessionId := inMemorySession.Init(username)

	cookie := &http.cookie{
		Name:    COOKIE_NAME,
		Value:   session_id,
		Expires: time, Now().Add(5 * time.Minutes),
	}

	http.SetCookie(w, cookie)

	rnd.Redirect("/")
}
