package task

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var mu = sync.RWMutex{}

var DB *sqlx.DB

type Task struct {
	ID        int    `db:"id"`
	TaskID    int    `db:"task_id"`
	Content   string `json:"content" db:"content"`
	EnrollCnt int    `json:"enroll_cnt" db:"enroll_cnt"`
	CreateBy  int    `db:"create_by"`
}

func (t Task) String() string {
	return fmt.Sprintf("[%d] %s", t.TaskID, t.Content)
}

func TaskByID(db *sqlx.DB, taskID int) (t Task, err error) {
	sqlStr := "SELECT id, task_id, content, enroll_cnt, create_by  FROM tasks WHERE id = ?"
	err = db.QueryRowx(sqlStr, taskID).Scan(&t.ID, &t.TaskID, &t.Content, &t.EnrollCnt, &t.CreateBy)
	if err != nil {
		err = errors.Wrap(err, "tasks by ID error")
		return
	}
	return
}

func TaskExist(db *sqlx.DB, taskID int) (ok bool, err error) {
	sqlStr := "SELECT COUNT(*) > 0 FROM tasks WHERE id = ?"
	err = db.QueryRowx(sqlStr, taskID).Scan(&ok)
	if err != nil {
		err = errors.Wrap(err, "TaskExist")
		return
	}
	return
}

func TaskCountByChat(db *sqlx.DB, chatID int64) (cnt int, err error) {
	sqlStr := "SELECT COUNT(*) FROM tasks WHERE chat_id = ?"
	mu.RLock()
	err = db.QueryRowx(sqlStr, chatID).Scan(&cnt)
	mu.RUnlock()
	if err != nil {
		err = errors.Wrap(err, "task count by chat error")
		return
	}
	return
}

func TasksByChat(db *sqlx.DB, chatID int64) (tl []Task, err error) {
	sqlStr := "SELECT id, task_id, content, enroll_cnt  FROM tasks WHERE chat_id = ? ORDER BY task_id ASC"
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

func DelTask(db *sqlx.DB, taskID int) (err error) {
	sqlStr := "DELETE FROM tasks WHERE id = ?"
	_, err = db.Queryx(sqlStr, taskID)
	if err != nil {
		err = errors.Wrap(err, "del task error")
		return
	}
	return
}

func AddTask(db *sqlx.DB, task string, enrollCnt int, chatID int64, createBy int) (tID int, err error) {
	sqlStr := "INSERT INTO tasks (task_id, content, enroll_cnt, chat_id, create_by)VALUES(?, ?, ?, ?, ?)"
	tot, err := TaskCountByChat(db, chatID)
	if err != nil {
		err = errors.Wrap(err, "add task error")
		return
	}
	tID = tot + 1
	mu.Lock()
	_, err = db.Queryx(sqlStr, tID, task, enrollCnt, chatID, createBy)
	mu.Unlock()
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
