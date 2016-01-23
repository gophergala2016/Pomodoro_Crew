package main

import (
	"fmt"
	"log"
	"net/http"
	"html/template"

	"github.com/googollee/go-socket.io"

	"github.com/google/cayley"
	_ "github.com/google/cayley/graph/bolt"

	"github.com/google/cayley/graph"
	"time"
	"strconv"
)

func main() {
	fmt.Println("Listening on port :3000")

	serve()
}

type Home struct {
	Users []User
}

type User struct {
	Name string
	FreeAt int64
}

func (u *User) IsBusy() bool {
	return u.FreeAt > time.Now().Unix()
}

func home(w http.ResponseWriter, r *http.Request, store *cayley.Handle) {
	p := cayley.StartPath(store, "user").In()
	it := p.BuildIterator()
	users := []User{}
	for cayley.RawNext(it) {
		userName := store.NameOf(it.Result())
		users = append(users, User{userName, time.Now().Unix()})
		log.Println("User:", store.NameOf(it.Result()))
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	t.Execute(w, Home{users})
}

func serve() {

	path := "/tmp/pc"

	graph.InitQuadStore("bolt", path, nil)

	store, err := cayley.NewGraph("bolt", path, nil)
	if err != nil {
		log.Fatalln(err)
	}

	server, err := initSocket(store)
	if err != nil {
		log.Fatalln(err)
	}

	initTimer(store)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		home(w, r, store)
	})
	http.Handle("/socket.io/", server)
	http.ListenAndServe(":3000", nil)
}

func initSocket(store *cayley.Handle) (*socketio.Server, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("chat")
		so.On("start", func(name string) {
			log.Println("emit:", so.Emit("start", name))
			so.BroadcastTo("chat", "start", name)
			addUser(name, store)
			makeBusy(name, store)
		})
		so.On("stop", func(name string) {
			log.Println("emit:", so.Emit("stop", name))
			so.BroadcastTo("chat", "stop", name)
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
				release(time.Now().Unix(), store)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func userExists(name string, store *cayley.Handle) bool {
	p := cayley.StartPath(store).Has(name)
	it := p.BuildIterator()
	if cayley.RawNext(it) {
		return true
	} else {
		return false
	}
}

func isBusy(name string, store *cayley.Handle) bool {
	p := cayley.StartPath(store, name).Has("free at")
	it := p.BuildIterator()
	if cayley.RawNext(it) {
		return true
	} else {
		return false
	}
}

func release(t int64, store *cayley.Handle) {
	freeAt := strconv.FormatInt(t, 10)
	p := cayley.StartPath(store, freeAt).In("free at")
	it := p.BuildIterator()
	for cayley.RawNext(it) {
		quad := cayley.Quad(store.NameOf(it.Result()), "free at", freeAt, "")
		store.RemoveQuad(quad)
	}
}

func following() {

}

func followers() {

}

func makeBusy(name string, store *cayley.Handle) {
	freeAt := strconv.FormatInt(time.Now().Unix() + 10, 10)
	quad := cayley.Quad(name, "free at", freeAt, "")
	log.Println(quad)
	store.AddQuad(quad)
}

func addUser(name string, store *cayley.Handle) {
	store.AddQuad(cayley.Quad(name, "is", "user", ""))
}