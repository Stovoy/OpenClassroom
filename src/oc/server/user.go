package server

import (
	"net/http"

	"fmt"
	"github.com/gorilla/mux"
	"oc/db"
	"text/template"
)

func loginHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) == 0 || len(username) > 20 || len(password) < 8 {
		errorJSONResponse(w, fmt.Errorf("Username and Password badly formatted"))
		return
	}

	valid, err := db.CheckLogin(username, password)
	if err != nil {
		errorJSONResponse(w, err)
		return
	}
	if valid {
		tokenCookie := createTokenCookie()
		err = db.WriteToken(username, tokenCookie.Value)
		if err != nil {
			errorJSONResponse(w, err)
			return
		}
		http.SetCookie(w, tokenCookie)
		usernameCookie, err := createUsernameCookie(username)
		if err != nil {
			errorJSONResponse(w, err)
			return
		}
		http.SetCookie(w, usernameCookie)
		printJSON(w, struct{ Result string }{"Successful"})
	} else {
		errorJSONResponse(w, fmt.Errorf("Invalid username or password"))
	}
}

func logoutHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	if c.Authenticated {
		err := db.Logout(c.Username)
		if err != nil {
			errorJSONResponse(w, err)
			return
		}
		printJSON(w, struct{ Result string }{"Successful"})
	} else {
		errorJSONResponse(w, fmt.Errorf("Not logged in"))
	}
}

func registerHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) == 0 || len(password) < 8 {
		errorJSONResponse(w, fmt.Errorf("Username and Password badly formatted"))
		return
	}

	err := db.Register(username, password)
	if err != nil {
		errorJSONResponse(w, err)
		return
	}
	tokenCookie := createTokenCookie()
	err = db.WriteToken(username, tokenCookie.Value)
	if err != nil {
		errorJSONResponse(w, err)
		return
	}
	http.SetCookie(w, tokenCookie)
	usernameCookie, err := createUsernameCookie(username)
	if err != nil {
		errorJSONResponse(w, err)
		return
	}
	http.SetCookie(w, usernameCookie)
	printJSON(w, struct{ Result string }{"Successful"})
}

func userHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	user := mux.Vars(r)["user"]

	t := template.Must(template.ParseFiles(
		"../static/html/user.html",
		"../static/html/head.html",
		"../static/html/top-bar.html"))

	err := t.Execute(w, struct {
		Authenticated bool
		Username      string
		PageUsername  string
	}{c.Authenticated, c.Username, user})
	if err != nil {
		errorResponse(w, err)
	}
}
