package main

import (
	"html/template"
	"log"

	"github.com/google/cayley/graph"
	"net/http"
)

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {

	path := "/tmp/pc"

	graph.InitQuadStore("bolt", path, nil)

	storage, err := models.GetStorage()
	if err != nil {
		log.Fatalln(err)
	}

	server, err := server.NewServer(storage)
	if err != nil {
		log.Fatalln(err)
	}

	user := models.NewUser("admin")
	user.Iteration()

	m.Get("/socket.io/", func(w http.ResponseWriter, rnd render.Render, r *http.Request, s *session.Session) {
		server.SetSession(s)
		server.ServeHTTP(w, r)
	})
	m.Run()
}
