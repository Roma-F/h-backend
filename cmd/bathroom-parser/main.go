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

const updateBathroomStatus = `update ads set bathroom_status = ? where id = ?`
const insertCountOpt = `insert into ads_opts (ad_id, opt_type, opt_value_uint) values (?, ?, ?)`
const getBathroomQuery = `
	select ads.id, bathroom_opt.opt_value_str
	from ads
			left join ads_opts as bathroom_opt on ads.id = bathroom_opt.ad_id and bathroom_opt.opt_type = 34
	where ads.bathroom_status = 0
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
	rows, err := db.Query(getBathroomQuery, batchSize)
	if err != nil {
		return fmt.Errorf("can't query ads: %v", err)
	}

	hasRows := false
	for rows.Next() {
		if !hasRows {
			hasRows = true
		}

		var id int
		var bathroom string
		err = rows.Scan(&id, &bathroom)
		if err != nil {
			return fmt.Errorf("can't scan ad: %v", err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("can't begin tx: %v", err)
		}

		combinedBathroomCnt, separateBathroomCnt, err := parseBathroom(bathroom)
		if err != nil {
			_, err = tx.Exec(updateBathroomStatus, 2, id)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("can't insert status 2: %v", err)
			}
			errCommit := tx.Commit()
			if errCommit != nil {
				tx.Rollback()
				return fmt.Errorf("can't commit status 2: %v", errCommit)
			}
			continue
		}
		if combinedBathroomCnt > 0 {
			_, err = tx.Exec(insertCountOpt, id, 54, combinedBathroomCnt)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("can't commit status 2: %v", err)
			}
		}
		if separateBathroomCnt > 0 {
			_, err = tx.Exec(insertCountOpt, id, 55, separateBathroomCnt)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("can't commit status 2: %v", err)
			}
		}
		_, err = tx.Exec(updateBathroomStatus, 1, id)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("can't exec bathroom insert: %v", err)
		}
		err = tx.Commit()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("can't commit bathroom count: %v", err)
		}
		fmt.Println("Parsed ad", id)
	}

	if !hasRows {
		return sql.ErrNoRows
	}

	return nil
}

func parseBathroom(bathroomStr string) (int, int, error) {
	bathroomParams := strings.Split(bathroomStr, ", ")
	if len(bathroomParams) == 2 {
		combinedBathroomCnt, err := strconv.Atoi(string(bathroomParams[0][0]))
		if err != nil {
			return 0, 0, err
		}
		separateBathroomCnt, err := strconv.Atoi(string(bathroomParams[1][0]))
		if err != nil {
			return 0, 0, err
		}
		return combinedBathroomCnt, separateBathroomCnt, nil
	} else if len(bathroomParams) == 1 {
		if strings.Contains(bathroomParams[0], "совмещен") {
			combinedBathroomCnt, err := strconv.Atoi(string(bathroomParams[0][0]))
			if err != nil {
				return 0, 0, err
			}
			return combinedBathroomCnt, 0, nil
		} else if strings.Contains(bathroomParams[0], "раздель") {
			separateBathroomCnt, err := strconv.Atoi(string(bathroomParams[0][0]))
			if err != nil {
				return 0, 0, err
			}
			return 0, separateBathroomCnt, nil
		}
		return 0, 0, errors.New("unknown bathroom type")
	}
	return 0, 0, errors.New("unknown bathroom type")
}
