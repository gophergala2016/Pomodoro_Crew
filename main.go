package main

import (
	"fmt"
	"html/template"
	"log"
	//"net/http"

	"github.com/go-martini/martini"
	"github.com/googollee/go-socket.io"
	"github.com/martini-contrib/render"

	"github.com/gophergala2016/Pomodoro_Crew/routes"
	"github.com/gophergala2016/Pomodoro_Crew/session"

	"github.com/google/cayley"
	_ "github.com/google/cayley/graph/bolt"

	"github.com/google/cayley/graph"
	"net/http"
	"time"
)

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {

	path := "/tmp/pc"

	graph.InitQuadStore("bolt", path, nil)

	store, err := cayley.NewGraph("bolt", path, nil)
	if err != nil {
		log.Fatalln(err)
	}

	initTimer(store)

	m := martini.Classic()

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

	m.Use(session.Middleware)

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                         // Specify what path to load the templates from.
		Layout:     "layout",                            // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"},          // Specify extensions to load for templates.
		Funcs:      []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
		Charset:    "UTF-8",                             // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                                // Output human readable JSON
	}))

	server, err := initSocket(store)
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
		server.ServeHTTP(w, r)
		fmt.Println("^^^^^^^^^^^")
		fmt.Println(s)
	})
	m.Run()
}

func initSocket(store *cayley.Handle) (*socketio.Server, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("chat")
		so.On("start", func() {
			log.Println("emit:", so.Emit("start"))
			so.BroadcastTo("chat", "start")
			//addUser(name, store)
			//makeBusy(name, store)
		})
		so.On("stop", func() {
			log.Println("emit:", so.Emit("stop"))
			so.BroadcastTo("chat", "stop")
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

func initTimer(store *cayley.Handle) {
	ticker := time.NewTicker(time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				//release(time.Now().Unix(), store)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
