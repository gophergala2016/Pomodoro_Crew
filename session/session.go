package session

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-martini/martini"
)

const (
	TOKEN_NAME = "token"
	TOKEN_STR  = "Pomodoro crew"
)

type Session struct {
	Id           string
	Username     string
	IsAuthorized bool
}

type SessionStore struct {
	data map[string]*Session
}

func NewSessionStore() *SessionStore {
	s := new(SessionStore)
	s.data = make(map[string]*Session)
	return s
}

func (store *SessionStore) Get(token string) *Session {
	session := store.data[token]
	if session == nil {
		return &Session{Id: token}
	}
	return session
}

func (store *SessionStore) Set(session *Session) {
	store.data[session.Id] = session
}

func ensureCookie(r *http.Request, w http.ResponseWriter) string {

	tokenCookie, err := r.Cookie(TOKEN_NAME)
	switch {
	case err == http.ErrNoCookie:
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, r.Cookie)
		return "Dont have valid token!"
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Cookie parse error: %v\n", err)
		return "Error while Parsing cookie!"
	}
	if tokenCookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Dont have valid token!")
	}
	token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(TOKEN_STR), nil
	})

	switch err.(type) {

	// TODO check whether it's needed
	case nil:

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return "Invalid Token!"
		}

	case *jwt.ValidationError:
		vErr := err.(*jwt.ValidationError)

		switch vErr.Errors {
		case jwt.ValidationErrorExpired:
			w.WriteHeader(http.StatusUnauthorized)
			return "Token Expired, get a new one."
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("ValidationError error: %+v\n", vErr.Error())
			return "Error while Parsing Token!"
		}

	default:
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Token parse error: %v\n", err)
		return "Error while Parsing Token!"
	}

	tokenString, err := token.SignedString([]byte(TOKEN_STR))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Token Signing error: %v\n", err)
		return "Sorry, error while Signing Token!"
	}

	return tokenString

}

var SessionStorage = NewSessionStore()

func Middleware(ctx martini.Context, r *http.Request, w http.ResponseWriter) {
	assets := regexp.MustCompile(`.*assets.*`)
	if r.URL.Path != "/login" && !assets.MatchString(r.URL.Path) && r.URL.Path != "/" {
		sessionId := ensureCookie(r, w)
		session := SessionStorage.Get(sessionId)

		ctx.Map(session)

		ctx.Next()

		SessionStorage.Set(session)
	}
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["Name"] = "token"
	token.Claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	tokenString, err := token.SignedString([]byte(TOKEN_STR))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Token Signing error: %v\n", err)
		fmt.Fprintln(w, "Sorry, error while Signing Token!")
	}

	session := SessionStorage.Get(tokenString)

	ctx.Map(session)

	ctx.Next()

	SessionStorage.Set(session)
}
