package cities

import (
	"encoding/json"
	cities "homeho-backend/internal/core/cities"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type getCitiesResponse struct {
	Cities []cities.Cities `json:"cities"`
}

func (cc *CitiesController) GetCities(rw http.ResponseWriter, r *http.Request) {
	cities, err := cc.citiesStore.GetCities()
	if err != nil {
		log.Err(err).Msg("can't load cities")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	resultJson, err := json.Marshal(getCitiesResponse{cities})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(resultJson))
}

func (cc *CitiesController) GetCitiesSearch(rw http.ResponseWriter, r *http.Request) {
	cityStr := chi.URLParam(r, "city_str")

	cities, err := cc.citiesStore.GetCitiesByName(cityStr)
	if err != nil {
		log.Err(err).Msg("can't load cities")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	resultJson, err := json.Marshal(getCitiesResponse{cities})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(resultJson))

}
