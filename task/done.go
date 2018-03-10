package task

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type RankingObj struct {
	DoneBy string `db:"by"`
	Count  int    `db:"count"`
}

func Ranking(db *sqlx.DB, count int) (rls []RankingObj, err error) {
	sqlStr := "SELECT `by`, COUNT(*) count FROM task_done GROUP BY `by` ORDER BY `count` desc LIMIT ?"
	rows, er := db.Queryx(sqlStr, count)
	if er != nil {
		err = errors.Wrap(er, "Ranking:")
		return
	}
	defer rows.Close()
	robj := RankingObj{}
	for rows.Next() {
		err = rows.Scan(&robj.DoneBy, &robj.Count)
		if err != nil {
			err = errors.Wrap(err, "Ranking:")
			return
		}
		rls = append(rls, robj)
	}
	return
}
