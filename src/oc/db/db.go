package db

import (
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"fmt"

	_ "github.com/lib/pq"
	"strings"
)

var db *sql.DB

func Start() error {
	connection, err := sql.Open("postgres", "user=root dbname=openclassroom sslmode=disable")
	if err != nil {
		return err
	}

	db = connection
	return nil
}

func WriteToken(username string, token string) error {
	username = strings.ToLower(username)
	_, err := db.Exec(
		`UPDATE users
		 SET token=$1
		 WHERE lowername=$2`, token, username)
	return err
}

func Logout(username string) error {
	username = strings.ToLower(username)
	_, err := db.Exec(
		`UPDATE users
		 SET token=''
		 WHERE lowername=$1`, username)
	return err
}

func GetOriginalName(username string) (string, error) {
	username = strings.ToLower(username)
	var originalUsername string
	err := db.QueryRow(
		`SELECT name FROM users
	 	 WHERE lowername=$1`, username).Scan(&originalUsername)
	if err != nil {
		return "", err
	}
	return originalUsername, nil
}

func CheckToken(username string, token string) (bool, error) {
	username = strings.ToLower(username)
	var readToken string
	var count int
	err := db.QueryRow(
		`SELECT COUNT(*) FROM users
		 WHERE lowername=$1`, username).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		err := db.QueryRow(
			`SELECT token FROM users
	 	 	 WHERE lowername=$1`, username).Scan(&readToken)
		if err != nil {
			return false, err
		}
		return token == readToken, nil
	}
	return false, nil
}

func Register(username string, password string) error {
	lowerUsername := strings.ToLower(username)
	var count int
	encryptedPassword := shaEncryption(password)
	err := db.QueryRow(
		`SELECT COUNT(*) FROM users
		 WHERE lowername=$1`, username).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("User already exists")
	}
	_, err = db.Exec(
		`INSERT INTO users (name, lowername, encrypted_password, token)
		 VALUES ($1, $2, $3, '')`, username, lowerUsername, encryptedPassword)
	if err != nil {
		return err
	}
	return nil
}

func CheckLogin(username string, password string) (bool, error) {
	username = strings.ToLower(username)
	var count int
	encryptedPassword := shaEncryption(password)
	err := db.QueryRow(
		`SELECT COUNT(*) FROM users
		 WHERE lowername=$1 AND encrypted_password=$2`,
		username, encryptedPassword).Scan(&count)
	if err != nil {
		return false, err
	}
	if count == 1 {
		return true, nil
	}
	return false, nil
}

func HasContent(page string) (bool, error) {
	page = strings.ToLower(page)
	var count int
	err := db.QueryRow(
		`SELECT COUNT(*) FROM wiki
		 WHERE page=$1`,
		page).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func LoadContent(page string, content string) error {
	page = strings.ToLower(page)
	var count int
	err := db.QueryRow(
		`SELECT COUNT(*) FROM wiki
		 WHERE page=$1`,
		page).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = db.Exec(
			`INSERT INTO wiki (page, content)
			 VALUES ($1, $2)`, page, content)
		if err != nil {
			return err
		}
		_, err = db.Exec(
			`INSERT INTO chats (wiki_id)
			 SELECT wiki.id
			 FROM wiki
			 WHERE wiki.page=$1`, page)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetContent(page string) (string, error) {
	page = strings.ToLower(page)
	var content string
	err := db.QueryRow(
		`SELECT content FROM wiki
		 WHERE page=$1`,
		page).Scan(&content)
	if err != nil {
		return "", fmt.Errorf("Page %s does not exist.", page)
	}
	return content, nil
}

func shaEncryption(text string) string {
	sha := sha512.New()
	sha.Write([]byte(text))
	return base64.StdEncoding.EncodeToString(sha.Sum(nil))
}
