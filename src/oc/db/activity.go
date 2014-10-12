package db

import (
	"strings"
	"time"
)

func WriteActivity(activity string, username string) error {
	username = strings.ToLower(username)
	_, err := db.Exec(
		`INSERT INTO activity (user_id, action, time)
		 SELECT users.id, $1, $2
		 FROM users
		 WHERE users.lowername=$3`, activity, time.Now(), username)
	return err
}

func GetLast50Activities() ([]Activity, error) {
	var activities []Activity = make([]Activity, 0)

	rows, err := db.Query(
		`SELECT a.id, users.name, a.action, a.time
		 FROM activity a, users
		 WHERE a.user_id=users.id
		 ORDER BY a.id
		 LIMIT (50)`)
	if err != nil {
		return activities, err
	}
	for rows.Next() {
		activity := Activity{}
		var t time.Time

		err = rows.Scan(&activity.ID, &activity.User, &activity.Action, &t)
		if err != nil {
			return activities, err
		}
		activity.Time = t.Format("02 Jan 06 15:04")
		activities = append(activities, activity)
	}
	return activities, nil
}

func GetLast50UserActivities(username string) ([]Activity, error) {
	username = strings.ToLower(username)
	var activities []Activity = make([]Activity, 0)

	rows, err := db.Query(
		`SELECT a.id, users.name, a.action, a.time
		 FROM activity a, users
		 WHERE a.user_id=users.id AND
		 users.lowername=$1
		 ORDER BY a.id
		 LIMIT (50)`, username)
	if err != nil {
		return activities, err
	}
	for rows.Next() {
		activity := Activity{}
		var t time.Time

		err = rows.Scan(&activity.ID, &activity.User, &activity.Action, &t)
		if err != nil {
			return activities, err
		}
		activity.Time = t.Format("02 Jan 06 15:04")
		activities = append(activities, activity)
	}
	return activities, nil
}

type Activity struct {
	ID     string
	User   string
	Action string
	Time   string
}
