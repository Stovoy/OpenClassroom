package server

import (
	"fmt"
	"net/http"
	"oc/db"
	"sync/atomic"
	"time"
)

var signedInChatConnections map[string]map[string]*SignedInConnection = make(map[string]map[string]*SignedInConnection)
var guestChatConnections map[string]map[string]*GuestConnection = make(map[string]map[string]*GuestConnection)
var guestNumber uint64 = 0
var timeout int = 10000

type GuestConnection struct {
	Name          string
	LastRefreshed time.Time
}

type SignedInConnection struct {
	Name          string
	LastRefreshed time.Time
}

func chatLoadNewHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	page := "/wiki/" + r.FormValue("page")
	lastMessage := r.FormValue("lastMessage")
	identifier := r.FormValue("identifier")

	refreshChat(c.Authenticated, page, c.Username, identifier)

	messages, err := db.GetChatMessagesAfter(page, lastMessage)
	if err != nil {
		errorJSONResponse(w, err)
		return
	}
	users, err := getChatUsers(page)
	if err != nil {
		errorJSONResponse(w, err)
		return
	}
	chatInfo := ChatInformation{Users: users, NewMessages: messages}
	printJSON(w, chatInfo)
}

func refreshChat(authenticated bool, page string, username string, identifier string) {
	if authenticated {
		if signedInChatConnections[page] == nil {
			signedInChatConnections[page] = make(map[string]*SignedInConnection)
		}
		if signedInChatConnections[page][username] == nil {
			signedInChatConnections[page][username] = &SignedInConnection{
				Name:          username,
				LastRefreshed: time.Now(),
			}
		} else {
			signedInChatConnections[page][username].LastRefreshed = time.Now()
		}
	} else {
		if guestChatConnections[page] == nil {
			guestChatConnections[page] = make(map[string]*GuestConnection)
		}
		if guestChatConnections[page][identifier] == nil {
			newGuestNumber := atomic.AddUint64(&guestNumber, 1)
			guestChatConnections[page][identifier] = &GuestConnection{
				Name:          "Guest" + string(newGuestNumber),
				LastRefreshed: time.Now(),
			}
		} else {
			guestChatConnections[page][identifier].LastRefreshed = time.Now()
		}
	}
}

func getChatUsers(page string) ([]User, error) {
	var users []User = make([]User, 0)

	for k, v := range signedInChatConnections[page] {
		user := User{}
		if v.LastRefreshed.Sub(time.Now()).Seconds() > 20 {
			delete(signedInChatConnections[page], k)
			continue
		}
		user.IsGuest = false
		user.Name = v.Name
		users = append(users, user)
	}
	for k, v := range guestChatConnections[page] {
		user := User{}
		if v.LastRefreshed.Sub(time.Now()).Seconds() > 20 {
			delete(guestChatConnections[page], k)
		}
		user.IsGuest = true
		user.Name = v.Name
		users = append(users, user)
	}
	return users, nil
}

type ChatInformation struct {
	Users       []User
	NewMessages []db.ChatMessage
}

type User struct {
	IsGuest bool
	Name    string
}

func chatWriteHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	page := r.FormValue("page")
	message := r.FormValue("message")

	if c.Authenticated {
		db.SendMessage(c.Username, message, page)
	}
	errorJSONResponse(w, fmt.Errorf("Not logged in"))
}
