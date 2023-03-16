package ads

import (
	"encoding/json"
	"fmt"
	core_ads "lbr-backend/internal/core/ads"
	store_ads "lbr-backend/internal/storage/ads"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rs/zerolog/log"
)

type getAdsResponse struct {
	Ads        []core_ads.Ad `json:"ads"`
	TotalPages int           `json:"totalPages"`
	Limit      int           `json:"limit"`
	Page       int           `json:"page"`
}

func (ac *AdsController) GetAds(rw http.ResponseWriter, r *http.Request) {
	log.Debug().Msg(r.URL.RawQuery)
	sp, err := buildSearchParamsFromQuery(r.URL.Query())
	if err != nil {
		log.Err(err).Msgf("AdsController#GetAds - invalid input params")
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	ads, pagesCnt, err := ac.adsSearch.Search(sp)
	if err == store_ads.ErrNoSuchPageExists {
		rw.WriteHeader(http.StatusOK)
		return
	} else if err != nil {
		log.Err(err).Msgf("AdsController#GetAds - can't load ads")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultJson, err := json.Marshal(getAdsResponse{ads, pagesCnt, sp.Limit, sp.Page})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(resultJson))
}

func buildSearchParamsFromQuery(qv url.Values) (core_ads.AdsSearchParams, error) {
	fmt.Println(qv, "build search params from query")
	sp := core_ads.AdsSearchParams{
		IsStudio:                         qv.Get("is_studio") == "true",
		Elevator:                         qv.Get("elevator") == "true",
		City:                             qv.Get("city"),
		WindowView:                       qv.Get("window_view"),
		Renovation:                       qv.Get("renovation"),
		BuildingType:                     qv.Get("building_type"),
		JoistType:                        qv.Get("joist_type"),
		TranslitAddressDependentLocality: qv.Get("translit_address_dependent_locality"),
		TranslitAddressAddressStreet:     qv.Get("translit_address_address_street"),
		AddressHouseNumber:               qv.Get("address_house_number"),
		TranslitAddressHouseNumber:       qv.Get("translit_address_house_number"),
	}

	page, err := parsePositiveInt(qv.Get("page"))
	if err != nil {
		return sp, fmt.Errorf("invalid page: %v", err)
	}
	sp.Page = page

	limitParam := qv.Get("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		return sp, fmt.Errorf("invalid int limit[%s]", limitParam)
	}
	sp.Limit = limit

	if priceMinParam, ok := qv["price_min"]; ok {
		priceMin, err := parsePositiveInt(priceMinParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid price_min: %v", err)
		}
		sp.PriceMin = priceMin
	}
	if priceMaxParam, ok := qv["price_max"]; ok {
		priceMax, err := parsePositiveInt(priceMaxParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid price_max: %v", err)
		}
		sp.PriceMax = priceMax
	}
	if spaceMinParam, ok := qv["space_min"]; ok {
		spaceMin, err := parsePositiveInt(spaceMinParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid space_min: %v", err)
		}
		sp.SpaceMin = spaceMin
	}
	if spaceMaxParam, ok := qv["space_max"]; ok {
		spaceMax, err := parsePositiveInt(spaceMaxParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid space_max: %v", err)
		}
		sp.SpaceMax = spaceMax
	}
	if floorParam, ok := qv["floor"]; ok {
		floor, err := strconv.Atoi(floorParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid floor: %v", err)
		}
		sp.HasFloor = true
		sp.Floor = floor
	}
	if roomsCountParam, ok := qv["rooms_count"]; ok {
		roomsCount, err := parsePositiveInt(roomsCountParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid rooms_count: %v", err)
		}
		sp.RoomsCount = roomsCount
	}
	fmt.Println(qv)
	if floorMinParam, ok := qv["floor_min"]; ok {
		floorMin, err := parsePositiveInt(floorMinParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid floor_min: %v", err)
		}
		sp.FloorMin = floorMin
	}
	if floorMaxParam, ok := qv["floor_max"]; ok {
		floorMax, err := parsePositiveInt(floorMaxParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid floor_ax: %v", err)
		}
		sp.FloorMax = floorMax
	}
	if livingSpaceMinParam, ok := qv["living_space_min"]; ok {
		livingSpaceMin, err := parsePositiveInt(livingSpaceMinParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid living_space_min: %v", err)
		}
		sp.LivingSpaceMin = livingSpaceMin
	}
	if livingSpaceMaxParam, ok := qv["living_space_max"]; ok {
		livingSpaceMax, err := parsePositiveInt(livingSpaceMaxParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid living_space_max: %v", err)
		}
		sp.LivingSpaceMax = livingSpaceMax
	}
	if kitchenSpaceMinParam, ok := qv["kitchen_space_min"]; ok {
		kitchenSpaceMin, err := parsePositiveInt(kitchenSpaceMinParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid kitchen_space_min: %v", err)
		}
		sp.KitchenSpaceMin = kitchenSpaceMin
	}
	if kitchenSpaceMaxParam, ok := qv["kitchen_space_max"]; ok {
		kitchenSpaceMax, err := parsePositiveInt(kitchenSpaceMaxParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid kitchen_space_max: %v", err)
		}
		sp.KitchenSpaceMax = kitchenSpaceMax
	}
	if priceForSqmMinParam, ok := qv["price_for_sqm_min"]; ok {
		priceForSqmMin, err := parsePositiveInt(priceForSqmMinParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid price_for_sqm_min: %v", err)
		}
		sp.PriceForSqmMin = priceForSqmMin
	}
	if priceForSqmMaxParam, ok := qv["price_for_sqm_max"]; ok {
		priceForSqmMax, err := parsePositiveInt(priceForSqmMaxParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid price_for_sqm_max: %v", err)
		}
		sp.PriceForSqmMax = priceForSqmMax
	}
	if ceilingHeightMinParam, ok := qv["ceiling_height_min"]; ok {
		ceilingHeightMin, err := parsePositiveInt(ceilingHeightMinParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid ceiling_height_min: %v", err)
		}
		sp.CeilingHeightMin = ceilingHeightMin
	}
	if ceilingHeightMaxParam, ok := qv["ceiling_height_max"]; ok {
		ceilingHeightMax, err := parsePositiveInt(ceilingHeightMaxParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid ceiling_height_max: %v", err)
		}
		sp.CeilingHeightMax = ceilingHeightMax
	}
	if totalFloorsMinParam, ok := qv["total_floors_min"]; ok {
		totalFloorsMin, err := parsePositiveInt(totalFloorsMinParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid total_floors_min: %v", err)
		}
		sp.TotalFloorsMin = totalFloorsMin
	}
	if totalFloorsMaxParam, ok := qv["total_floors_max"]; ok {
		totalFloorsMax, err := parsePositiveInt(totalFloorsMaxParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid total_floors_max: %v", err)
		}
		sp.TotalFloorsMax = totalFloorsMax
	}
	if constructionYearParam, ok := qv["construction_year"]; ok {
		constructionYear, err := parsePositiveInt(constructionYearParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid construction_year: %v", err)
		}
		sp.ConstructionYear = constructionYear
	}
	if balconyCntParam, ok := qv["balcony_cnt"]; ok {
		balconyCnt, err := parsePositiveInt(balconyCntParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid balcony_cnt: %v", err)
		}
		sp.BalconyCnt = balconyCnt
	}
	if loggiaCntParam, ok := qv["loggia_cnt"]; ok {
		loggiaCnt, err := parsePositiveInt(loggiaCntParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid loggia_cnt: %v", err)
		}
		sp.LoggiaCnt = loggiaCnt
	}
	if combinedBathroomCntParam, ok := qv["combined_bathroom_cnt"]; ok {
		combinedBathroomCnt, err := parsePositiveInt(combinedBathroomCntParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid combined_bathroom_cnt: %v", err)
		}
		sp.CombinedBathroomCnt = combinedBathroomCnt
	}
	if separateBathroomCntParam, ok := qv["separate_bathroom_cnt"]; ok {
		separateBathroomCnt, err := parsePositiveInt(separateBathroomCntParam[0])
		if err != nil {
			return sp, fmt.Errorf("invalid separate_bathroom_cnt: %v", err)
		}
		sp.SeparateBathroomCnt = separateBathroomCnt
	}

	return sp, nil
}

func parsePositiveInt(value string) (int, error) {
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid int [%s]", value)
	}
	if parsedValue < 1 {
		return 0, fmt.Errorf("negative int [%d]", parsedValue)
	}
	return parsedValue, nil
}
