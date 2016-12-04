package task

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var DB *sqlx.DB

type Task struct {
	ID        int    `db:"id"`
	TaskID    int    `db:"task_id"`
	Content   string `json:"content" db:"content"`
	EnrollCnt int    `json:"enroll_cnt" db:"enroll_cnt"`
}

func TaskByID(db *sqlx.DB, taskID int) (t Task, err error) {
	sqlStr := "SELECT id, content, enroll_cnt  FROM tasks WHERE id = ?"
	err = db.QueryRowx(sqlStr, taskID).Scan(&t.ID, &t.Content, &t.EnrollCnt)
	if err != nil {
		err = errors.Wrap(err, "tasks by ID error")
		return
	}
	return
}

func TaskCountByChat(db *sqlx.DB, chatID int64) (cnt int, err error) {
	sqlStr := "SELECT COUNT(*) FROM tasks WHERE chat_id = ?"
	err = db.QueryRowx(sqlStr, chatID).Scan(&cnt)
	if err != nil {
		err = errors.Wrap(err, "task count by chat error")
		return
	}
	return
}

func TasksByChat(db *sqlx.DB, chatID int64) (tl []Task, err error) {
	sqlStr := "SELECT id, task_id, content, enroll_cnt  FROM tasks WHERE chat_id = ? ORDER BY task_id ASC LIMIT 50"
	rows, err := db.Queryx(sqlStr, chatID)
	if err != nil {
		err = errors.Wrap(err, "tasks by chat error")
		return
	}
	defer rows.Close()
	t := Task{}
	for rows.Next() {
		err = rows.Scan(&t.ID, &t.TaskID, &t.Content, &t.EnrollCnt)
		if err != nil {
			err = errors.Wrap(err, "task by chat error: get rows error")
			return
		}
		tl = append(tl, t)
	}
	return
}

func TaskRealID(db *sqlx.DB, taskID int, chatID int64) (tID int, err error) {
	sqlStr := "SELECT id FROM tasks WHERE task_id = ? AND chat_id = ? ORDER BY `id` DESC LIMIT 1"
	err = db.QueryRowx(sqlStr, taskID, chatID).Scan(&tID)
	if err != nil {
		err = errors.Wrap(err, "task realID error")
		return
	}
	return
}

func AddTask(db *sqlx.DB, task string, enrollCnt int, chatID int64) (err error) {
	sqlStr := "INSERT INTO tasks (task_id, content, enroll_cnt, chat_id)VALUES(?, ?, ?, ?)"
	tot, err := TaskCountByChat(db, chatID)
	if err != nil {
		err = errors.Wrap(err, "add task error")
		return
	}
	tID := tot + 1
	_, err = db.Queryx(sqlStr, tID, task, enrollCnt, chatID)
	if err != nil {
		err = errors.Wrap(err, "add task error")
		return
	}
	return
}

func FinishCountByTaskID(db *sqlx.DB, taskID int) (cnt int, err error) {
	sqlStr := "SELECT COUNT(*) FROM task_done WHERE task_id = ?"
	err = db.QueryRowx(sqlStr, taskID).Scan(&cnt)
	if err != nil {
		err = errors.Wrap(err, "get finish count error")
		return
	}
	return
}

func IsDone(db *sqlx.DB, taskID int, by string) (done bool, err error) {
	sqlStr := "SELECT * FROM task_done WHERE `task_id` = ? AND `by` = ?"
	rows, err := db.Queryx(sqlStr, taskID, by)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() == false {
		done = false
		return
	}
	done = true
	return
}

func AddDone(db *sqlx.DB, taskID int, by string) (err error) {
	sqlStr := "INSERT INTO task_done (`task_id`, `by`) VALUES(?, ?)"
	_, err = db.Queryx(sqlStr, taskID, by)
	if err != nil {
		err = errors.Wrap(err, "add done error")
		return
	}
	return
}
