package ads

import (
	"errors"
	"fmt"
	core_ads "lbr-backend/internal/core/ads"
	store_ads "lbr-backend/internal/storage/ads"

	"github.com/rs/zerolog/log"
)

type AdsFether interface {
	Search(asp core_ads.AdsSearchParams) ([]core_ads.Ad, int, error)
	GetPageAds(limit, offset int) ([]core_ads.Ad, error)
	GetAdsPagesCount(limit int) (int, error)
}
type AdsSearchService struct {
	adsStore AdsFether
}

func NewAdsSearchService(adsStore AdsFether) *AdsSearchService {
	return &AdsSearchService{adsStore}
}

func (as *AdsSearchService) Search(asp core_ads.AdsSearchParams) ([]core_ads.Ad, int, error) {
	if !asp.HasFilters() {
		return as.getAdsWithoutFilters(asp.Page, asp.Limit)
	}
	fmt.Println(asp, "ads search service")
	if asp.City == "" {
		return nil, 0, errors.New("city is required")
	}

	ads, pageCnt, err := as.adsStore.Search(asp)
	if err != nil {
		log.Err(err).Msgf("AdsSearchService#Search - can't load ads")
		return nil, 0, err
	}

	return ads, pageCnt, nil
}

func (as *AdsSearchService) getAdsWithoutFilters(page, limit int) ([]core_ads.Ad, int, error) {
	ads, err := as.adsStore.GetPageAds(limit, (page-1)*limit)
	if err == store_ads.ErrNoSuchPageExists {
		return nil, 0, err
	} else if err != nil {
		log.Err(err).Msgf("AdsSearchService#Search - can't load ads")
		return nil, 0, err
	}

	pagesCnt, err := as.adsStore.GetAdsPagesCount(limit)
	if err != nil {
		log.Err(err).Msgf("AdsSearchService#Search - can't load ads pages cnt")
		return nil, 0, err
	}

	return ads, pagesCnt, nil
}
