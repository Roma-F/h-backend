package ads

import (
	"encoding/json"
	storage_ads "lbr-backend/internal/storage/ads"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type getAdPhoneResponse struct {
	Phone string `json:"phone"`
}

func (ac *AdsController) GetAdPhone(rw http.ResponseWriter, r *http.Request) {
	adId, err := strconv.Atoi(chi.URLParam(r, "ad_id"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	phone, err := ac.adsStore.GetAdPhone(adId)
	if err == storage_ads.ErrNoPhone {
		rw.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Err(err).Msgf("can't load ad[%d] phone", adId)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultJson, err := json.Marshal(getAdPhoneResponse{phone})
	if err != nil {

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(resultJson))
}
