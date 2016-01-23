package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"blog-example/db/documents"
	"blog-example/models"
	"blog-example/session"
	"blog-example/utils"

	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"

	"gopkg.in/mgo.v2"
)

const COOKIE_NAME = "sessionId"

var postsCollection *mgo.Collection
var inMemorySession *session.Session

func main() {
	fmt.Println("Listening on port :3000")

	inMemorySession = session.NewSession()
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	postsCollection = session.DB("blog").C("posts")
	m := martini.Classic()

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

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
	m.Get("/login", getLoginHandler)
	m.Post("/", postLoginHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit/:id", editHandler)
	m.Get("/delete/:id", deleteHandler)
	m.Post("/SavePost", savePostHandler)

	m.Post("/gethtml", getHtmlHandler)

}

func indexHandler(rnd render.Render, w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie(COOKIE_NAME)
	if cookie != nil {
		fmt.Println(inMemorySession.Get(cookie.Value))
	}

	postDocuments := []documents.PostDocument{}
	postsCollection.Find(nil).All(&postDocuments)

	posts := []models.Post{}
	for _, doc := range postDocuments {
		post := models.Post{doc.Id, doc.Title, doc.ContentHtml, doc.ContentMarkdown}
		posts = append(posts, post)
	}

	rnd.HTML(200, "index", posts)
}

func writeHandler(rnd render.Render) {
	post := models.Post{}
	rnd.HTML(200, "write", post)
}

func unescape(x string) interface{} {
	return template.HTML(x)
}

func savePostHandler(rnd render.Render, r *http.Request) {
	id := utils.GenerateId()
	title := r.FormValue("title")
	contentMarkdown := r.FormValue("content")
	contentHtml := utils.ConvertMarkdownToHtml(contentMarkdown)

	postDocument := documents.PostDocument{id, title, contentHtml, contentMarkdown}
	if id != "" {
		postsCollection.UpdateId(id, postDocument)
	} else {
		id = GenerateId()
		postDocument.Id = id
		postsCollection.Insert(postDocument)
	}
	rnd.Redirect("/")
}

func getHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	html := utils.ConvertMarkdownToHtml(md)

	rnd.JSON(200, map[string]interface{}{"html": string(html)})
}

func editHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]

	postDocument := documents.PostDocument{}
	err := postsCollection.FindId(id).One(&postDocument)

	if err != nil {
		rnd.Redirect("/")
	}
	post := models.Post{postDocument.Id, postDocument.Title, postDocument.ContentHtml, postDocument.ContentMarkdown}

	rnd.HTML(200, "write", post)
}

func deleteHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]
	if id != "" {
		rnd.Redirect("/")
	}

	postsCollection.RemoveId(id)
	rnd.Redirect("/")
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
