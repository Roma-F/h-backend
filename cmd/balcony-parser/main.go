package main

import (
	"database/sql"
	"errors"
	"fmt"
	"homeho-backend/internal/database"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

const updateBalconyStatus = `update ads set balcony_status = ? where id = ?`
const insertCountOpt = `insert into ads_opts (ad_id, opt_type, opt_value_uint) values (?, ?, ?)`
const getBalconyQuery = `
	select ads.id, balcony_opt.opt_value_str
	from ads
			left join ads_opts as balcony_opt on ads.id = balcony_opt.ad_id and balcony_opt.opt_type = 29
	where ads.balcony_status = 0
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
	rows, err := db.Query(getBalconyQuery, batchSize)
	if err != nil {
		return fmt.Errorf("can't query ads: %v", err)
	}

	hasRows := false
	for rows.Next() {
		if !hasRows {
			hasRows = true
		}

		var id int
		var balcony string
		err = rows.Scan(&id, &balcony)
		if err != nil {
			return fmt.Errorf("can't scan ad: %v", err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("can't begin tx: %v", err)
		}

		balconyCnt, loggiaCnt, err := extractBalconyAndLoggia(balcony)
		if err != nil {
			_, errExec := tx.Exec(updateBalconyStatus, 2, id)
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

		if balconyCnt > 0 {
			_, err = tx.Exec(insertCountOpt, id, 52, balconyCnt)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("can't commit status 2: %v", err)
			}
		}
		if loggiaCnt > 0 {
			_, err = tx.Exec(insertCountOpt, id, 53, loggiaCnt)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("can't commit status 2: %v", err)
			}
		}

		_, err = tx.Exec(updateBalconyStatus, 1, id)
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

func extractBalconyAndLoggia(balconyStr string) (int, int, error) {
	balconyParams := strings.Split(balconyStr, ", ")
	if len(balconyParams) == 2 {
		balconyCnt, err := strconv.Atoi(string(balconyParams[0][0]))
		if err != nil {
			return 0, 0, err
		}
		loggiaCnt, err := strconv.Atoi(string(balconyParams[1][0]))
		if err != nil {
			return 0, 0, err
		}
		return balconyCnt, loggiaCnt, nil
	} else if len(balconyParams) == 1 {
		if strings.Contains(balconyParams[0], "балкон") {
			balconyCnt, err := strconv.Atoi(string(balconyParams[0][0]))
			if err != nil {
				return 0, 0, err
			}
			return balconyCnt, 0, nil
		} else if strings.Contains(balconyParams[0], "лоджи") {
			loggiaCnt, err := strconv.Atoi(string(balconyParams[0][0]))
			if err != nil {
				return 0, 0, err
			}
			return 0, loggiaCnt, nil
		}
		return 0, 0, errors.New("invalid balcony")
	}
	return 0, 0, errors.New("invalid balcony")
}
