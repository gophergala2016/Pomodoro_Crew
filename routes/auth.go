package routes

/*import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gophergala2016/Pomodoro_Crew/session"

	"github.com/dgrijalva/jwt-go"
	"github.com/martini-contrib/render"
)

func GetLoginHandler(rnd render.Render) {
	rnd.HTML(200, "login", nil)
}

func PostLoginHandler(w http.ResponseWriter, rnd render.Render, r *http.Request, s *session.Session) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Println(username)
	fmt.Println(password)

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["Name"] = "token"
	token.Claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	tokenString, err := token.SignedString([]byte(username + session.TOKEN_STR))
	fmt.Println(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Token Signing error: %v\n", err)
		fmt.Fprintln(w, "Sorry, error while Signing Token!")
	}

	http.SetCookie(w, &http.Cookie{
		Name:  session.TOKEN_NAME,
		Value: tokenString,
	})
	s.Username = username
	s.Id = tokenString
	s.IsAuthorized = true
	rnd.Redirect("/")
}

func LogoutHandler(rnd render.Render, r *http.Request, s *session.Session) {
	s.Username = ""
	s.IsAuthorized = false

	rnd.Redirect("/")
}*/
