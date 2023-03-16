package main

import (
	"homeho-backend/internal/database"
	search_ads "homeho-backend/internal/service/ads_search"
	store_ads "homeho-backend/internal/storage/ads"
	store_cities "homeho-backend/internal/storage/cities"
	"homeho-backend/internal/transport/rest"
	"os"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "02/01 15:04:05"})

	env, exists := os.LookupEnv("SIMPLE_AUTH_ENV")
	if !exists {
		env = "development"
	}

	db, err := database.SetUpDB(env)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	goquBuilder := goqu.Dialect("mysql")

	adsStore, err := store_ads.NewAdsStore(db, goquBuilder)
	if err != nil {
		panic(err)
	}

	citiesStore, err := store_cities.NewCitiesStore(db, goquBuilder)
	if err != nil {
		panic(err)
	}

	adsSearch := search_ads.NewAdsSearchService(adsStore)

	handlerDeps := rest.Deps{AdsStore: adsStore, AdsSearch: adsSearch, CitiesStore: citiesStore}
	handler := rest.NewRestHandler(handlerDeps)

	s := rest.NewServer(handler, "3030")
	log.Info().Msg("Starting server on port :3030")
	err = s.ListenAndServe()

	if err != nil {
		panic(err)
	}
	defer s.Close()
}
