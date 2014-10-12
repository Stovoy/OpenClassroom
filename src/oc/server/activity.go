package server

import (
	"github.com/gorilla/mux"
	"net/http"

	"oc/db"
)

func activityLoadGlobalHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	activities, err := db.GetLast50Activities()

	if err != nil {
		errorJSONResponse(w, err)
	}

	printJSON(w, struct{ Activities interface{} }{Activities: activities})
}

func activityLoadUserHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["user"]

	activities, err := db.GetLast50UserActivities(username)

	if err != nil {
		errorJSONResponse(w, err)
	}

	printJSON(w, struct{ Activities interface{} }{Activities: activities})
}
