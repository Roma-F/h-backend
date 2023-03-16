package main

import (
	"encoding/json"
	"fmt"
	"io"
	"lbr-backend/internal/database"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := database.SetUpDB("development")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	h := Handler{db}

	http.HandleFunc("/", h.Handle)
	http.ListenAndServe(":8080", nil)
}

// Брянск - City
// Россия, Брянск, Бежицкий район, Комсомольская улица, 46 - FormattedAddress
// 46 - HouseNumber
// 241035 - Postcode
// Брянская область - State
// городской округ Брянск - StateDistrict
// Комсомольская улица - Street
// Бежицкий район - District
// а так же longitude и latitude

type PostReq struct {
	Id                    int    `json:"id"`
	AdministrativeArea    string `json:"administrativeAreaName"`
	SubAdministrativeArea string `json:"subAdministrativeAreaName"`
	City                  string `json:"city"`
	DependentLocality     string `json:"dependentLocality"`
	Street                string `json:"street"`
	HouseNumber           string `json:"house"`
	PostalCode            string `json:"postalCode"`
	Latitude              string `json:"latitude"`
	Longitude             string `json:"longitude"`
}

type Entry struct {
	Id      int    `json:"id" db:"id"`
	Address string `json:"address" db:"address"`
}

type Handler struct {
	db *sqlx.DB
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Accept", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet {
		h.Get(w, r)
		return
	}

	if r.Method == http.MethodPost {
		h.Post(w, r)
		return
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	getAdressesQuery := `
			select ads.id id, address_opt.opt_value_str as address
			from ads
				left join ads_opts as address_opt on ads.id = address_opt.ad_id and address_opt.opt_type = 1
			where ads.address_status = 0
				and opt_value_str is not null
			limit 100;
		`
	rows, err := h.db.Query(getAdressesQuery)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var respBody []Entry

	var entry Entry
	for rows.Next() {
		err := rows.Scan(&entry.Id, &entry.Address)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBody = append(respBody, Entry{entry.Id, entry.Address})
	}

	w.WriteHeader(http.StatusOK)
	bodyStr, err := json.Marshal(respBody)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(bodyStr)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var reqBody PostReq
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	insertStrOptQuery := `
		insert into ads_opts (ad_id, opt_type, opt_value_str)
		values (?, ?, ?) on duplicate key update opt_value_str = ?;
	`

	// begin db tx
	tx, err := h.db.Begin()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hasData := false

	if len(reqBody.AdministrativeArea) > 0 {
		_, err := tx.Exec(insertStrOptQuery, reqBody.Id, 41, reqBody.AdministrativeArea, reqBody.AdministrativeArea)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			tx.Rollback()
			return
		}
		hasData = true
	}
	if len(reqBody.SubAdministrativeArea) > 0 {
		_, err := tx.Exec(insertStrOptQuery, reqBody.Id, 42, reqBody.SubAdministrativeArea, reqBody.SubAdministrativeArea)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			tx.Rollback()
			return
		}
		hasData = true
	}
	if len(reqBody.City) > 0 {
		_, err := tx.Exec(insertStrOptQuery, reqBody.Id, 43, reqBody.City, reqBody.City)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			tx.Rollback()
			return
		}
		hasData = true
	}
	if len(reqBody.DependentLocality) > 0 {
		_, err := tx.Exec(insertStrOptQuery, reqBody.Id, 44, reqBody.DependentLocality, reqBody.DependentLocality)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			tx.Rollback()
			return
		}
		hasData = true
	}
	if len(reqBody.Street) > 0 {
		_, err := tx.Exec(insertStrOptQuery, reqBody.Id, 45, reqBody.Street, reqBody.Street)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			tx.Rollback()
			return
		}
		hasData = true
	}
	if len(reqBody.HouseNumber) > 0 {
		_, err := tx.Exec(insertStrOptQuery, reqBody.Id, 46, reqBody.HouseNumber, reqBody.HouseNumber)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			tx.Rollback()
			return
		}
		hasData = true
	}
	if len(reqBody.PostalCode) > 0 {
		_, err := tx.Exec(insertStrOptQuery, reqBody.Id, 47, reqBody.PostalCode, reqBody.PostalCode)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			tx.Rollback()
			return
		}
		hasData = true
	}
	if len(reqBody.Latitude) > 0 {
		_, err := tx.Exec(insertStrOptQuery, reqBody.Id, 48, reqBody.Latitude, reqBody.Latitude)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			tx.Rollback()
			return
		}
		hasData = true
	}
	if len(reqBody.Longitude) > 0 {
		_, err := tx.Exec(insertStrOptQuery, reqBody.Id, 49, reqBody.Longitude, reqBody.Longitude)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			tx.Rollback()
			return
		}
		hasData = true
	}

	if hasData {
		_, err := tx.Exec(`update ads set address_status = 1 where id = ?`, reqBody.Id)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			tx.Rollback()
			return
		}

		fmt.Println(reqBody.Id, "updated successfully")
		w.WriteHeader(http.StatusCreated)
		return
	} else {
		tx.Rollback()
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
