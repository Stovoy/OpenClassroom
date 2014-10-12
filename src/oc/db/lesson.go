package db

import (
	"strings"

	_ "github.com/lib/pq"
	"time"
)

func WriteLesson(username string, lessonName string) (int, error) {
	username = strings.ToLower(username)
	var id int
	err := db.QueryRow(
		`INSERT INTO lessons (teacher_id, name, running, start_time)
		 SELECT users.id, $1, false, $2
		 FROM users
		 WHERE users.lowername=$3
		 RETURNING id`, lessonName, time.Now(), username).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func WriteLessonItem(username string, lessonName string, action string) error {
	username = strings.ToLower(username)
	_, err := db.Exec(
		`INSERT INTO lesson_items (lesson_id, number, action)
		 SELECT lessons.id, (
		 	SELECT COUNT(*)
		 	FROM lesson_items
		 	WHERE lesson_id=lessons.id
		 ), $2
		 FROM lessons
		 WHERE lessons.name=$1`, lessonName, action)
	return err
}

func GetLesson(lessonID string) ([]LessonItem, error) {
	var lessonItems []LessonItem = make([]LessonItem, 0)

	rows, err := db.Query(
		`SELECT li.id, li.number, li.action
		 FROM lesson_items li, lessons l, users u
		 WHERE l.id=$1 AND
		 u.id=l.teacher_id AND
		 li.lesson_id = l.id
		 ORDER BY li.number`, lessonID)
	if err != nil {
		return lessonItems, err
	}
	for rows.Next() {
		lessonItem := LessonItem{}

		err = rows.Scan(&lessonItem.ID, &lessonItem.Number, &lessonItem.Action)
		if err != nil {
			return lessonItems, err
		}
		lessonItems = append(lessonItems, lessonItem)
	}
	return lessonItems, nil
}

func GetLessonHeader(lessonID string) (Lesson, error) {
	var lesson Lesson = Lesson{}

	rows, err := db.Query(
		`SELECT l.id, l.name, u.name, l.running, l.start_time, l.end_time
		 FROM lessons l, users u
		 WHERE l.id=$1 AND
		 u.id=l.teacher_id`, lessonID)
	if err != nil {
		return lesson, err
	}
	if rows.Next() {
		var startT time.Time
		var endT *time.Time
		err = rows.Scan(&lesson.ID, &lesson.Name, &lesson.Teacher,
			&lesson.Running, &startT, &endT)
		if err != nil {
			return lesson, err
		}
		lesson.StartTime = startT.Format("Mon, 02 Jan 2006 15:04:05")
		if endT == nil {
			lesson.EndTime = ""
		} else {
			lesson.EndTime = endT.Format("Mon, 02 Jan 2006 15:04:05")
		}
	}
	return lesson, nil
}

func GetLessonsForUser(username string) ([]Lesson, error) {
	username = strings.ToLower(username)
	var lessons []Lesson = make([]Lesson, 0)

	rows, err := db.Query(
		`SELECT l.id, l.name, users.name, l.running, l.start_time, l.end_time
		 FROM lessons l, users
		 WHERE l.teacher_id=users.id
		 ORDER BY l.id DESC`)
	if err != nil {
		return lessons, err
	}
	for rows.Next() {
		lesson := Lesson{}
		var startT time.Time
		var endT *time.Time

		err = rows.Scan(&lesson.ID, &lesson.Name, &lesson.Teacher, &lesson.Running,
			&startT, &endT)
		if err != nil {
			return lessons, err
		}
		lesson.StartTime = startT.Format("Mon, 02 Jan 2006 15:04:05")
		if endT == nil {
			lesson.EndTime = ""
		} else {
			lesson.EndTime = endT.Format("Mon, 02 Jan 2006 15:04:05")
		}
		lessons = append(lessons, lesson)
	}
	return lessons, nil
}

func GetGlobalLessons() ([]Lesson, error) {
	var lessons []Lesson = make([]Lesson, 0)

	rows, err := db.Query(
		`SELECT l.id, users.name, l.running, l.start_time, l.end_time
		 FROM lessons l, users
		 ORDER BY l.id DESC`)
	if err != nil {
		return lessons, err
	}
	for rows.Next() {
		lesson := Lesson{}
		var startT time.Time
		var endT *time.Time

		err = rows.Scan(&lesson.ID, &lesson.Name, &lesson.Teacher, &lesson.Running,
			&startT, &endT)
		if err != nil {
			return lessons, err
		}
		lesson.StartTime = startT.Format("Mon, 02 Jan 2006 15:04:05")
		if endT == nil {
			lesson.EndTime = ""
		} else {
			lesson.EndTime = endT.Format("Mon, 02 Jan 2006 15:04:05")
		}
		lessons = append(lessons, lesson)
	}
	return lessons, nil
}

type Lesson struct {
	ID        string
	Name      string
	Teacher   string
	Running   bool
	StartTime string
	EndTime   string
}

type LessonItem struct {
	ID     string
	Number uint64
	Action string
}
