package model

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type User struct {
	UUID           string         `db:"uuid" json:"uuid"`
	ID             int            `db:"id" json:"id"`
	UserName       string         `db:"user_name" json:"user_name"`
	DispName       string         `db:"disp_name" json:"disp_name"`
	CreateAt       mysql.NullTime `db:"create_at" json:"-"`
	UpdateAt       mysql.NullTime `db:"update_at" json:"-"`
	Exist          bool           `db:"exist" json:"exist"` // if Exist = false, the object is treated as nil
	DontTrack      string         `db:"dont_track" json:"dont_track"`
	MoyuPhraseUUID string         `db:"moyu_phrase_uuid" json:"-"`
	PhraseUUID     string         `db:"phrase_uuid" json:"-"`
}

func SelectUser(db *sqlx.DB, id int) (u User, err error) {
	sqlStr := "SELECT COUNT(*)>0 FROM users WHERE id = ?"
	ok := false
	err = db.QueryRowx(sqlStr, id).Scan(&ok)
	if err != nil {
		err = errors.Wrap(err, "SelectUser")
		return
	}
	if !ok {
		u.Exist = false
		return
	}
	sqlStr = "SELECT * FROM users WHERE id = ?"
	err = db.QueryRowx(sqlStr, id).StructScan(&u)
	if err != nil {
		err = errors.Wrap(err, "SelectUser")
		return
	}
	u.Exist = true
	return
}

func UpdateUser(db *sqlx.DB, u User) (err error) {
	sqlStr := "UPDATE users SET user_name = ?, disp_name = ?, dont_track = ? ,phrase_uuid = ? WHERE id = ?"
	_, err = db.Exec(sqlStr, u.UserName, u.DispName, u.DontTrack, u.PhraseUUID, u.ID)
	if err != nil {
		err = errors.Wrap(err, "UpdateUser")
	}
	return
}

func CreateUser(db *sqlx.DB, u User) (err error) {
	sqlStr := "INSERT INTO users (uuid, id, user_name, disp_name)VALUES(?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, u.UUID, u.ID, u.UserName, u.DispName)
	if err != nil {
		err = errors.Wrap(err, "CreateUser")
		return
	}
	return
}

func ListUser(db *sqlx.DB, page int) (ul []User, err error) {
	sqlStr := "SELECT * FROM users ORDER BY create_at DESC LIMIT 20 OFFSET ?"
	rows, er := db.Queryx(sqlStr, (page-1)*20)
	if er != nil {
		err = errors.Wrap(er, "ListUser")
		return
	}
	defer rows.Close()
	for rows.Next() {
		u := User{}
		err = rows.StructScan(&u)
		if err != nil {
			err = errors.Wrap(err, "ListUser")
			return
		}
		u.Exist = true
		ul = append(ul, u)
	}
	return
}

// func SelectUser(db *sqlx.DB, id int) (u User, err error) {
// 	sqlStr := "SELECT * FROM users WHERE id = ?"
// 	err = db.QueryRowx(sqlStr, id).StructScan(&u)
// 	if err != nil {
// 		err = errors.Wrap(err, "SelectUser")
// 		return
// 	}
// 	return
// }
