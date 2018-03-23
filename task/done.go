package task

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type RankingObj struct {
	DoneBy  string `db:"by"`
	Count   int    `db:"count"`
	Wanders int    `db:"wanders"`
}

func Ranking(db *sqlx.DB, count int) (rls []RankingObj, err error) {
	// sqlStr := "SELECT `by`, COUNT(*) count FROM task_done GROUP BY `by` ORDER BY `count` desc LIMIT ?"
	sqlStr := "SELECT IF(u.disp_name IS NULL, \"UNKNOWN\" ,u.disp_name) `by` ,COUNT(*) `count`, SUM(wander_times) wanders FROM active_task_instance ati  LEFT JOIN users u ON ati.user_id = u.id WHERE instance_state = 2 GROUP BY u.id ORDER BY count DESC LIMIT ?"
	rows, er := db.Queryx(sqlStr, count)
	if er != nil {
		err = errors.Wrap(er, "Ranking:")
		return
	}
	defer rows.Close()
	robj := RankingObj{}
	for rows.Next() {
		err = rows.Scan(&robj.DoneBy, &robj.Count, &robj.Wanders)
		if err != nil {
			err = errors.Wrap(err, "Ranking:")
			return
		}
		rls = append(rls, robj)
	}
	return
}
