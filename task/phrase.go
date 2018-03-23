package task

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Phrase struct {
	UUID      string         `db:"uuid"`
	Phrase    string         `db:"phrase"`
	CreateBy  int            `db:"create_by"`
	CreateAt  mysql.NullTime `db:"create_at"`
	UpdateAt  mysql.NullTime `db:"update_at"`
	Show      string         `db:"show"`
	GroupUUID string         `db:"group_uuid"`
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
