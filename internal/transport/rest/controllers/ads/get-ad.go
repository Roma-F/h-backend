package ads

import (
	"encoding/json"
	core_ads "homeho-backend/internal/core/ads"
	storage_ads "homeho-backend/internal/storage/ads"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type getAdResponse struct {
	Ad *core_ads.Ad `json:"ad"`
}

func (ac *AdsController) GetAd(rw http.ResponseWriter, r *http.Request) {
	adId, err := strconv.Atoi(chi.URLParam(r, "ad_id"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	ad, err := ac.adsStore.GetAdById(adId)
	if err == storage_ads.ErrAdNotFound {
		log.Warn().Msgf("AdsController#GetAd - ad[%d] not found: %v", adId, err.Error())
		rw.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Err(err).Msgf("AdsController#GetAd - error loading ad %d", ad.Id)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultJson, err := json.Marshal(getAdResponse{ad})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(resultJson))
}
