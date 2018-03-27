package task

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Phrase struct {
	UUID      string         `db:"uuid" json:"uuid"`
	Phrase    string         `db:"phrase" json:"phrase"`
	CreateBy  int            `db:"create_by" json:"create_by"`
	CreateAt  mysql.NullTime `db:"create_at" json:"create_at"`
	UpdateAt  mysql.NullTime `db:"update_at" json:"update_at"`
	Show      string         `db:"show" json:"show"`
	GroupUUID string         `db:"group_uuid" json:"group_uuid"`
}

func InsertPhrase(db *sqlx.DB, p Phrase) (err error) {
	sqlStr := "INSERT INTO phrases(uuid, phrase, create_by, group_uuid)VALUES(?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, p.UUID, p.Phrase, p.CreateBy, p.GroupUUID)
	if err != nil {
		err = errors.Wrap(err, "InsertPhrase")
		return err
	}
	return
}

func SelectPhrasesByGroupUUID(db *sqlx.DB, UUID string) (pl []Phrase, err error) {
	sqlStr := "SELECT * FROM phrases WHERE group_uuid = ?"
	rows, er := db.Queryx(sqlStr, UUID)
	if er != nil {
		err = errors.Wrap(er, "SelectPhraseByGroupUUID")
		return
	}
	for rows.Next() {
		p := Phrase{}
		err = rows.StructScan(&p)
		if err != nil {
			err = errors.Wrap(err, "SelectPhraseByGroupUUID")
			return
		}
		pl = append(pl, p)
	}
	return
}

func SelectPhraseByUUID(db *sqlx.DB, UUID string) (p Phrase, err error) {
	sqlStr := "SELECT * FROM phrases WHERE uuid = ?"
	err = db.QueryRowx(sqlStr, UUID).StructScan(&p)
	if err != nil {
		err = errors.Wrap(err, "SelectPhraseByUUID")
		return
	}
	return
}

func DeletePhraseByUUID(db *sqlx.DB, UUID string) (err error) {
	sqlStr := "DELETE FROM phrases WHERE uuid = ?"
	_, err = db.Exec(sqlStr, UUID)
	if err != nil {
		err = errors.Wrap(err, "DeletePhraseByUUID")
		return
	}
	return
}
