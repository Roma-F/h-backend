package cities

import (
	_ "embed"
	"errors"
	"fmt"
	core_cities "homeho-backend/internal/core/cities"
	"homeho-backend/internal/storage"
	"homeho-backend/pkg/dotdb"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

//go:embed queries.sql
var queries string

var ErrAdNotFound = errors.New("ad is not found")

type CitiesStore struct {
	sqlxDb      *sqlx.DB
	db          *dotdb.DotDB
	goquBuilder goqu.DialectWrapper
	optTypes    []storage.DbOptType
}

func NewCitiesStore(db *sqlx.DB, goquBuilder goqu.DialectWrapper) (*CitiesStore, error) {
	dotDb, err := dotdb.NewDotDB(db, queries)
	if err != nil {
		return nil, fmt.Errorf("can't create new CitiesStore, initing dotdb: %v", err)
	}

	var AdOptTypes []storage.DbOptType // todo: extract outside and pass here as params
	// err = dotDb.GetStructsSlice("get-ad-opt-types", &AdOptTypes)
	// if err != nil {
	// 	return nil, fmt.Errorf("can't create new CitiesStore, reading adOptTypes: %v", err)
	// }
	return &CitiesStore{
		sqlxDb:      db,
		db:          dotDb,
		goquBuilder: goquBuilder,
		optTypes:    AdOptTypes,
	}, nil
}

func (as *CitiesStore) GetCities() ([]core_cities.Cities, error) {
	var cities []core_cities.Cities
	err := as.db.GetStructsSlice("get-cities", &cities)
	if err != nil {
		return nil, err
	}
	ids := make([]int, len(cities))
	for i, city := range cities {
		ids[i] = city.Id
	}
	return cities, nil
}

// func (as *CitiesStore) GetCities(c string) ([]core_cities.Cities, error) {
// 	var cities []core_cities.Cities
// 	// c to lover case
// 	// cityPref = c // http://localhost:3030/v1/room/cities?cityPref=аба
// 	// доставать из гет параметров
// 	err := as.db.GetStructsSlice("get-ciries-by-name-or-translit", &cities, c+"%")

// 	if err != nil {
// 		return nil, err
// 	}

// 	ids := make([]int, len(cities))
// 	for i, city := range cities {
// 		ids[i] = city.Id
// 	}

// 	return cities, nil
// }

func (as *CitiesStore) GetCitiesByName(name string) ([]core_cities.Cities, error) {
	var cities []core_cities.Cities
	var cityName = strings.ToLower(name)
	err := as.db.GetStructsSlice("get-ciries-by-name-like", &cities, cityName+"%")

	if err != nil {
		return nil, fmt.Errorf("can't fetch cities: %v", err)
	}
	ids := make([]int, len(cities))
	for i, city := range cities {
		ids[i] = city.Id
	}

	fmt.Println(cities, err)
	return cities, nil
}
