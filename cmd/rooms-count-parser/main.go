package main

import (
	"database/sql"
	"errors"
	"fmt"
	"lbr-backend/internal/database"
	"strings"

	"github.com/jmoiron/sqlx"
)

const updateRoomsStatus = `update ads set rooms_status = ? where id = ?`
const insertRoomsCountOpt = `insert into ads_opts (ad_id, opt_type, opt_value_uint) values (?, ?, ?)`
const insertRoomsBoolOpt = `insert into ads_opts (ad_id, opt_type, opt_value_bool) values (?, ?, ?)`
const getTitleQuery = `
	select ads.id, title_opt.opt_value_str
	from ads
			left join ads_opts as title_opt on ads.id = title_opt.ad_id and title_opt.opt_type = 11
	where ads.rooms_status = 0
	and opt_value_str is not null
	limit ?
`

func main() {
	db, err := database.SetUpDB("development")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for {
		err = parseBatch(db, 100)
		if err == sql.ErrNoRows {
			fmt.Println("Work done")
			return
		} else if err != nil {
			panic(err)
		}
	}
}

func parseBatch(db *sqlx.DB, batchSize int) error {
	rows, err := db.Query(getTitleQuery, batchSize)
	if err != nil {
		return fmt.Errorf("can't query ads: %v", err)
	}

	hasRows := false
	for rows.Next() {
		if !hasRows {
			hasRows = true
		}

		var id int
		var title string
		err = rows.Scan(&id, &title)
		if err != nil {
			return fmt.Errorf("can't scan ad: %v", err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("can't begin tx: %v", err)
		}

		roomsCnt, isStudio, err := extractRoomCountFromTitle(title)
		if err != nil {
			_, errExec := tx.Exec(updateRoomsStatus, 2, id)
			if errExec != nil {
				tx.Rollback()
				return fmt.Errorf("can't insert status 2: %v", errExec)
			}
			errCommit := tx.Commit()
			if errCommit != nil {
				tx.Rollback()
				return fmt.Errorf("can't commit status 2: %v", errCommit)
			}
			continue
		}
		_, err = tx.Exec(insertRoomsCountOpt, id, 50, roomsCnt)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("can't commit status 2: %v", err)
		}
		if isStudio {
			_, err = tx.Exec(insertRoomsBoolOpt, id, 51, true)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("can't commit status 2: %v", err)
			}
		}
		_, err = tx.Exec(updateRoomsStatus, 1, id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("can't exec rooms insert: %v", err)
		}
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("can't commit rooms count: %v", err)
		}
		fmt.Println("Parsed ad", id)
	}

	if !hasRows {
		return sql.ErrNoRows
	}
	return nil
}

func extractRoomCountFromTitle(title string) (int, bool, error) {
	if strings.HasPrefix(title, "Студия") {
		return 1, true, nil
	}
	if strings.HasPrefix(title, "Апартаменты-студия") {
		return 1, true, nil
	}
	if strings.HasPrefix(title, "Комната") {
		return 1, false, nil
	}
	for i := 0; i < 10; i++ {
		if strings.HasPrefix(title, fmt.Sprintf("%d-комн.", i)) {
			return i, false, nil
		}
		if strings.Contains(title, fmt.Sprintf("%d комнаты", i)) {
			return i, false, nil
		}
	}
	return 0, false, errors.New("no rooms count detected")
}
