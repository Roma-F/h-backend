package rest

import (
	search_ads "homeho-backend/internal/service/ads_search"
	store_ads "homeho-backend/internal/storage/ads"
	store_cities "homeho-backend/internal/storage/cities"
	controller_ads "homeho-backend/internal/transport/rest/controllers/ads"
	controller_cities "homeho-backend/internal/transport/rest/controllers/cities"
	"homeho-backend/internal/transport/rest/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Deps struct {
	AdsStore    *store_ads.AdsStore
	AdsSearch   *search_ads.AdsSearchService
	CitiesStore *store_cities.CitiesStore
}

func NewRestHandler(deps Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.CORS)
	r.Use(middleware.CommonHeaders)

	adsController := controller_ads.NewAdsController(deps.AdsStore, deps.AdsSearch)
	citiesController := controller_cities.NewCitiesController(deps.CitiesStore)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/room", func(r chi.Router) {
			r.Get("/index", adsController.GetAds)
			r.Get("/{ad_id}", adsController.GetAd)
			r.Get("/{ad_id}/phone", adsController.GetAdPhone)
			r.Get("/recommended/{city}", adsController.GetRecommendedAds)
			r.Get("/cities", citiesController.GetCities)
			r.Get("/cities/{city_str}", citiesController.GetCitiesSearch)
		})
	})

	return r
}
