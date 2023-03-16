package ads

import (
	"encoding/json"
	core_ads "lbr-backend/internal/core/ads"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type getRecommendedAdsResponse struct {
	Ads []core_ads.Ad `json:"ads"`
}

func (ac *AdsController) GetRecommendedAds(rw http.ResponseWriter, r *http.Request) {
	cityTranslit := chi.URLParam(r, "city")
	ads, err := ac.adsStore.GetRecommended(10, cityTranslit)
	if err != nil {
		log.Err(err).Msg("can't load recommended ads")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultJson, err := json.Marshal(getRecommendedAdsResponse{ads})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(resultJson))
}
