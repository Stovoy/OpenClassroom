package db

import (
	"strings"
	"time"
)

func SendMessage(username string, message string, page string) error {
	page = strings.ToLower(page)
	username = strings.ToLower(username)

	_, err := db.Exec(
		`INSERT INTO chat_messages (chat_id, user_id, time, message)
		 SELECT chats.id, users.id, $1, $2
		 FROM chats, wiki, users
		 WHERE wiki.page=$3 AND
		 users.lowername=$4 AND
		 chats.wiki_id=wiki.id`,
		time.Now(), message, page, username)
	if err != nil {
		return err
	}
	return nil
}

func GetChatMessagesAfter(chat string, lastMessage string) ([]ChatMessage, error) {
	chat = strings.ToLower(chat)

	var messages []ChatMessage = make([]ChatMessage, 0)

	rows, err := db.Query(
		`SELECT m.id, users.name, m.time, m.message
		 FROM chat_messages m, wiki w, users
		 WHERE m.chat_id=w.id AND
		 w.page=$1 AND
		 m.id>$2 AND
		 m.user_id=users.id
		 ORDER BY m.id`, chat, lastMessage)
	if err != nil {
		return messages, err
	}
	for rows.Next() {
		message := ChatMessage{}
		var t time.Time

		err = rows.Scan(&message.ID, &message.User, &t, &message.Message)
		if err != nil {
			return messages, err
		}
		message.Time = t.Format("Mon, 02 Jan 2006 15:04:05")
		messages = append(messages, message)
	}
	return messages, nil
}

type ChatMessage struct {
	ID      string
	User    string
	Time    string
	Message string
}
