package task

import (
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	ATI_STATE_WORKING   = 1  // user explicitly working on this task
	ATI_STATE_FINISHED  = 2  // user finished the task
	ATI_STATE_INACTIVE  = 3  // user worked on it before but now he asked to give up for now
	ATI_STATE_IMMUTABLE = 4  // user set a reminder, the task is immutable, will keep repeated every cron expression
	ATI_STATE_INVALID   = -1 // the task no longer exist
)

type ActiveTaskInstance struct {
	InstanceUUID       string         `db:"instance_uuid"`
	TaskID             int            `db:"task_id"`
	UserID             int            `db:"user_id"`
	InstanceState      int            `db:"instance_state"`
	ReminderState      int            `db:"reminder_state"`
	ReminderExpression string         `db:"reminder_expression"`
	WanderTimes        int            `db:"wander_times"`
	NotifyID           int64          `db:"notify_to_id"`
	Cooldown           int            `db:"cooldown"`
	StartAt            mysql.NullTime `db:"start_at"`
	EndAt              mysql.NullTime `db:"end_at"`
}

func SelectATIByUUID(db *sqlx.DB, UUID string) (atil []ActiveTaskInstance, err error) {
	return
}

func SelectATIByUserIDAndState(db *sqlx.DB, uid int, state int) (atil []ActiveTaskInstance, err error) {
	sqlx := "SELECT * FROM active_task_instance WHERE user_id = ? AND instance_state = ?"
	rows, er := db.Queryx(sqlx, uid, state)
	if er != nil {
		err = errors.Wrap(er, "SelectATIByUserIDAndStateForUpdate")
		return
	}
	for rows.Next() {
		ati := ActiveTaskInstance{}
		rows.StructScan(&ati)
		atil = append(atil, ati)
	}
	return
}

func UpdateATIStateByUUID(db *sqlx.DB, UUID string, state int) (err error) {
	sqlx := "UPDATE active_task_instance SET instance_state = ? WHERE instance_uuid = ?"
	_, err = db.Exec(sqlx, state, UUID)
	if err != nil {
		err = errors.Wrap(err, "UpdateATIStateByUUID")
		return err
	}
	return
}

func SelectATIByUserIDAndChatIDAndState(db *sqlx.DB, uid int, cid int64, state int) (atil []ActiveTaskInstance, err error) {
	sqlx := "SELECT * FROM active_task_instance WHERE user_id = ? AND notify_to_id = ? AND instance_state = ?"
	rows, er := db.Queryx(sqlx, uid, cid, state)
	if er != nil {
		err = errors.Wrap(er, "SelectATIByUserIDAndChatIDAndStateForUpdate")
		return
	}
	for rows.Next() {
		ati := ActiveTaskInstance{}
		rows.StructScan(&ati)
		atil = append(atil, ati)
	}
	return
}

func FinishATI(db *sqlx.DB, uuid string) (err error) {
	sqlx := "UPDATE active_task_instance SET end_at = ?, instance_state = ? WHERE instance_uuid = ?"
	tx, er := db.Begin()
	if er != nil {
		err = errors.Wrap(er, "FinishATI")
		return
	}
	_, err = tx.Exec(sqlx, mysql.NullTime{Time: time.Now(), Valid: true}, ATI_STATE_FINISHED, uuid)
	if err != nil {
		err = errors.Wrap(err, "FinishATI")
		return
	}
	if err = tx.Commit(); err != nil {
		err = errors.Wrap(err, "FinishATI")
		return
	}
	return
}

func SelectATIByUserIDAndStateForUpdate(db *sqlx.DB, uid int, state int) (atil []ActiveTaskInstance, err error) {
	return SelectATIByUserIDAndState(db, uid, state)
}

func SelectATIByUserIDAndChatIDAndStateForUpdate(db *sqlx.DB, uid int, cid int64, state int) (atil []ActiveTaskInstance, err error) {
	return SelectATIByUserIDAndChatIDAndState(db, uid, cid, state)
}

func InsertATI(db *sqlx.DB, ati ActiveTaskInstance) (err error) {
	sqlx := `INSERT INTO active_task_instance (instance_uuid, task_id,
	instance_state, reminder_state, 
	reminder_expression, user_id, 
	notify_to_id, start_at, end_at)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlx, ati.InstanceUUID,
		ati.TaskID, ati.InstanceState,
		ati.ReminderState, ati.ReminderExpression, ati.UserID, ati.NotifyID,
		ati.StartAt, ati.EndAt)

	if err != nil {
		return errors.Wrap(err, "InsertATI")
	}
	return
}

func IncWanderTimes(db *sqlx.DB, uuid string) (err error) {
	sqlx := "UPDATE active_task_instance SET wander_times = wander_times + 1 WHERE instance_uuid = ?"
	_, err = db.Exec(sqlx, uuid)
	if err != nil {
		err = errors.Wrap(err, "IncWanderTimes")
		return
	}
	return
}
