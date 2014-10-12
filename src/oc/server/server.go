package server

import (
	"net/http"
	"text/template"

	"code.google.com/p/go-html-transform/h5"
	"code.google.com/p/go-html-transform/html/transform"
	"code.google.com/p/go.net/html"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"

	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"oc/db"
	"strings"
)

var tokenCookieName string = "OpenClassroomToken"
var usernameCookieName string = "OpenClassroomUsername"

type baseHandler func(*Context, http.ResponseWriter, *http.Request)

func Start() {
	router := mux.NewRouter()

	handleFunc(router, "/", homeHandler)

	// User
	handleFunc(router, "/login/", loginHandler)
	handleFunc(router, "/logout/", logoutHandler)
	handleFunc(router, "/register/", registerHandler)
	handleFunc(router, "/user/{username:.*}", userHandler)

	// Pages
	handleFunc(router, "/search/{page:.+}", searchHandler)
	handleFunc(router, "/wiki/{page:.+}", wikiHandler)

	// Chat
	handleFunc(router, "/chat/loadNew/", chatLoadNewHandler)
	handleFunc(router, "/chat/message/", chatMessageHandler)

	// Static
	router.HandleFunc("/js/{file:.+}", jsHandler)
	router.HandleFunc("/css/{file:.+}", cssHandler)

	http.Handle("/", router)

	fmt.Println("Running...")
	fmt.Println(http.ListenAndServe(":80", nil))
}

func handleFunc(router *mux.Router, path string, base baseHandler) {
	wrappedHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		c := &Context{}
		tokenCookie, err := r.Cookie(tokenCookieName)
		if err == http.ErrNoCookie {
			// No token cookie
			c.Authenticated = false
		} else if err != nil {
			errorResponse(w, err)
			return
		} else {
			// Has token cookie
			usernameCookie, err := r.Cookie(usernameCookieName)
			if err == http.ErrNoCookie {
				// No username cookie
				c.Authenticated = false
			} else if err != nil {
				errorResponse(w, err)
				return
			} else {
				// Has username cookie
				token := tokenCookie.Value
				username := usernameCookie.Value
				valid, err := db.CheckToken(username, token)
				if err != nil {
					errorResponse(w, err)
					return
				}
				if valid {
					c.Authenticated = true
					c.Username = username
				}
			}
		}

		base(c, w, r)
	}

	router.HandleFunc(path, wrappedHandler)
}

func homeHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(
		"../static/html/home.html",
		"../static/html/head.html",
		"../static/html/top-bar.html"))
	err := t.Execute(w, createStructFromContext(c))
	if err != nil {
		errorResponse(w, err)
		return
	}
}

func searchHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	page := mux.Vars(r)["page"]
	page = strings.ToLower(page)
	hasContent, err := db.HasContent("/wiki/" + page)
	if err != nil {
		errorJSONResponse(w, err)
		return
	}
	if hasContent {
		printJSON(w, struct{ Result string }{"/wiki/" + page})
		return
	}
	_, err = loadFromWikipedia(page)
	if err != nil {
		errorJSONResponse(w, err)
		return
	}
	printJSON(w, struct{ Result string }{"/wiki/" + page})
}

func wikiHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	page := mux.Vars(r)["page"]
	hasContent, err := db.HasContent("/wiki/" + page)
	if err != nil {
		errorResponse(w, err)
	}
	var content string
	if hasContent || !hasContent {
		content, err = loadFromWikipedia(page)
		err = db.LoadContent("/wiki/"+page, content)
		if err != nil {
			errorJSONResponse(w, err)
		}
	} else {
		content, err = db.GetContent("/wiki/" + page)
	}
	if err != nil {
		errorResponse(w, err)
		return
	}
	t := template.Must(template.ParseFiles(
		"../static/html/wiki.html",
		"../static/html/head.html",
		"../static/html/top-bar.html"))

	err = t.Execute(w, struct {
		Authenticated bool
		Username      string
		Content       string
		Page          string
	}{c.Authenticated, c.Username, content, page})
	if err != nil {
		errorResponse(w, err)
	}
}

func loadFromWikipedia(page string) (string, error) {
	doc, err := goquery.NewDocument(fmt.Sprintf(
		"https://en.wikipedia.org/w/index.php?search=%s&title=Special%3ASearch&go=Go", page))
	if err != nil {
		return "", err
	}
	if strings.Contains(doc.Url.Path, "w/index.php") {
		return "", fmt.Errorf("Page %s does not exist.", page)
	} else {
		sel := doc.Find("html")
		htmlResult, err := sel.Html()
		if err != nil {
			return "", err
		}
		reader := bytes.NewBufferString(htmlResult)
		tree, err := h5.New(reader)
		if err != nil {
			return "", err
		}
		t := transform.New(tree)

		removeFromParent := func() transform.TransformFunc {
			return func(n *html.Node) {
				n.Parent.RemoveChild(n)
			}
		}
		err = t.Apply(removeFromParent(), "#mw-panel")
		err = t.Apply(removeFromParent(), "#ma-panel")
		err = t.Apply(removeFromParent(), "#ca-talk")
		err = t.Apply(removeFromParent(), "#left-navigation")
		err = t.Apply(removeFromParent(), "#right-navigation")
		err = t.Apply(removeFromParent(), "#mw-head")
		err = t.Apply(removeFromParent(), "#footer")
		err = t.Apply(removeFromParent(), "#mw-page-base")
		err = t.Apply(removeFromParent(), "#mw-page-head")
		err = t.Apply(removeFromParent(), "#mw-head-base")
		if err != nil {
			return "", err
		}
		result := t.String()
		err = db.LoadContent(doc.Url.Path, result)
		if err != nil {
			return "", err
		}
		return result, nil
	}
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
	file := mux.Vars(r)["file"]
	bytes, err := ioutil.ReadFile("../static/js/" + file)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/javascript")
	fmt.Fprint(w, string(bytes))
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	file := mux.Vars(r)["file"]
	bytes, err := ioutil.ReadFile("../static/css/" + file)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/css")
	fmt.Fprint(w, string(bytes))
}

func createStructFromContext(c *Context) interface{} {
	return struct {
		Authenticated bool
		Username      string
	}{c.Authenticated, c.Username}
}

func printJSON(w http.ResponseWriter, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		errorResponse(w, err)
		return
	}
	fmt.Fprint(w, string(bytes))
}

func errorResponse(w http.ResponseWriter, err error) {
	fmt.Fprint(w, "Error: "+err.Error())
}

func errorJSONResponse(w http.ResponseWriter, err error) {
	printJSON(w, struct{ Error string }{err.Error()})
}

func createTokenCookie() *http.Cookie {
	bytes := securecookie.GenerateRandomKey(64)
	token := base64.StdEncoding.EncodeToString(bytes)
	return &http.Cookie{Name: tokenCookieName, Value: token, Path: "/"}
}

func createUsernameCookie(username string) (*http.Cookie, error) {
	original, err := db.GetOriginalName(username)
	if err != nil {
		return nil, err
	}
	return &http.Cookie{Name: usernameCookieName, Value: original, Path: "/"}, nil
}

type Context struct {
	Authenticated bool
	Username      string
}
