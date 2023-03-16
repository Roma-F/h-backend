package ads

import (
	search_ads "lbr-backend/internal/service/ads_search"
	storage_ads "lbr-backend/internal/storage/ads"
)

type AdsController struct {
	adsStore  *storage_ads.AdsStore
	adsSearch *search_ads.AdsSearchService
}

func NewAdsController(adsStore *storage_ads.AdsStore, adsSearch *search_ads.AdsSearchService) *AdsController {
	return &AdsController{
		adsStore:  adsStore,
		adsSearch: adsSearch,
	}
}
