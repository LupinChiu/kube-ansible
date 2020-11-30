package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var dbConnErr error

type Player struct {
	gorm.Model
	Name     string
	Account  string
	Password string
}

const indexPage = `
<h1>Login</h1>
<form method="post" action="/v1/login">
    <label for="account">Account</label>
    <input type="text" id="account" name="account">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

const internalPage = `
<h1>Internal</h1>
<hr>
<small>User: %s</small>
<form method="post" action="/v1/logout">
    <button type="submit">Logout</button>
</form>
`

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		fmt.Fprintf(response, internalPage, userName)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	account := request.FormValue("account")
	pass := request.FormValue("password")

	var player Player
	db.First(&player, "account = ?", account, pass)

	fmt.Println("login test.")
	redirectTarget := "/index"
	if player.Name != "" {
		// .. check credentials ..
		setSession(player.Name, response)
		redirectTarget = "/internal"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/index", 302)
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

var router = mux.NewRouter()

func main() {
	dsn := fmt.Sprintf("root:1234@tcp(%s:%s)/arpg?charset=utf8&parseTime=true", os.Args[1], os.Args[2])
	db, dbConnErr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if dbConnErr != nil {
		panic("failed to connect database")
	}

	router.HandleFunc("/index", indexPageHandler)
	router.HandleFunc("/internal", internalPageHandler)

	router.HandleFunc("/v1/login", loginHandler).Methods("POST")
	router.HandleFunc("/v1/logout", logoutHandler).Methods("POST")

	fmt.Println("server start.")
	http.ListenAndServe(":8033", router)

}
