package server

import (
	"github.com/gorilla/mux"
	"net/http"

	"fmt"
	"oc/db"
	"text/template"
)

func lessonLoadGlobalHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	lessons, err := db.GetGlobalLessons()
	if err != nil {
		errorJSONResponse(w, err)
	}

	printJSON(w, struct{ lessons interface{} }{lessons})
}

func lessonLoadUserHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["user"]

	lessons, err := db.GetLessonsForUser(username)
	if err != nil {
		errorJSONResponse(w, err)
	}

	printJSON(w, struct{ lessons interface{} }{lessons})
}

func lessonCreateHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(
		"../static/html/lessonCreation.html",
		"../static/html/head.html",
		"../static/html/top-bar.html"))
	err := t.Execute(w, struct {
		Authenticated bool
		Username      string
	}{c.Authenticated, c.Username})
	if err != nil {
		errorResponse(w, err)
		return
	}
}

func lessonWriteHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	username := c.Username
	lessonName := r.PostFormValue("lesson")

	lesson, err := db.WriteLesson(username, lessonName)
	if err != nil {
		errorJSONResponse(w, err)
		return
	}

	printJSON(w, struct{ Result string }{fmt.Sprintf("/lesson/%d", lesson)})
}

func lessonWriteLessonItem(c *Context, w http.ResponseWriter, r *http.Request) {
	r.PostFormValue("action")
}

func lessonLoadHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	lesson := mux.Vars(r)["lesson"]
	t := template.Must(template.ParseFiles(
		"../static/html/lesson.html",
		"../static/html/head.html",
		"../static/html/top-bar.html"))
	lessonObj, err := db.GetLessonHeader(lesson)
	if err != nil {
		errorResponse(w, err)
		return
	}
	controller := false
	if c.Username == lessonObj.Teacher {
		controller = true
	}
	err = t.Execute(w, struct {
		Authenticated bool
		Username      string
		Controller    bool
	}{c.Authenticated, c.Username, controller})
	if err != nil {
		errorResponse(w, err)
		return
	}
}
