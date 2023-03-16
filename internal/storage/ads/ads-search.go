package ads

import (
	"fmt"
	core_ads "homeho-backend/internal/core/ads"
	"math"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/rs/zerolog/log"
)

const adsOptsTable = "ads_opts"

func (as *AdsStore) Search(asp core_ads.AdsSearchParams) ([]core_ads.Ad, int, error) {
	cntSql, adsSql, err := as.buildSql(asp)
	if err != nil {
		return nil, 0, err
	}

	var cnt int
	err = as.sqlxDb.QueryRow(cntSql).Scan(&cnt)
	if err != nil {
		log.Err(err).Msgf("AdsSearchService#Search - can't fetch ads count: %v", err)
		return nil, 0, err
	} else if cnt == 0 {
		return nil, 0, ErrNoSuchPageExists
	}
	log.Debug().Msgf("AdsSearchService#Search - cnt: %d", cnt)

	pageCnt := int(math.Ceil(float64(cnt) / float64(asp.Limit)))
	if asp.Page > pageCnt {
		return nil, 0, ErrNoSuchPageExists
	}

	rows, err := as.sqlxDb.Query(adsSql)
	if err != nil {
		log.Err(err).Msgf("AdsSearchService#Search - can't fetch ads ids: %v", err)
		return nil, 0, err
	}

	adsIds := make([]int, 0, asp.Limit)
	for rows.Next() {
		var adId int
		err = rows.Scan(&adId)
		if err != nil {
			log.Err(err).Msgf("AdsSearchService#Search - can't scan ad id: %v", err)
			return nil, 0, err
		}
		adsIds = append(adsIds, adId)
	}
	err = rows.Err()
	if err != nil {
		log.Err(err).Msgf("AdsSearchService#Search - can't fetch ads ids: %v", err)
		return nil, 0, err
	}

	ads, err := as.GetAdsByIds(adsIds)
	if err != nil {
		log.Err(err).Msgf("AdsSearchService#Search - can't fetch ads: %v", err)
		return nil, 0, err
	}

	return ads, pageCnt, nil

}

func (as *AdsStore) buildSql(asp core_ads.AdsSearchParams) (string, string, error) {
	builder := as.goquBuilder
	fmt.Println(asp, "ads search params 2")
	query := builder.From(adsOptsTable).Select("ad_id").Where(goqu.Ex{
		"opt_type":      AdOptType_TranslitAddressCity,
		"opt_value_str": asp.City,
	})

	if asp.ConstructionYear > 0 {
		ex := goqu.Ex{"ad_id": query, "opt_type": AdOptType_ConstructionYear, "opt_value_uint": asp.ConstructionYear}
		query = builder.From(adsOptsTable).Select("ad_id").Where(ex)
	}
	if asp.WindowView != "" {
		ex := goqu.Ex{"ad_id": query, "opt_type": AdOptType_WindowView, "opt_value_str": asp.WindowView}
		query = builder.From(adsOptsTable).Select("ad_id").Where(ex)
	}
	if asp.Renovation != "" {
		ex := goqu.Ex{"ad_id": query, "opt_type": AdOptType_Renovation, "opt_value_str": asp.Renovation}
		query = builder.From(adsOptsTable).Select("ad_id").Where(ex)
	}
	if asp.BuildingType != "" {
		ex := goqu.Ex{"ad_id": query, "opt_type": AdOptType_BuildingType, "opt_value_str": asp.BuildingType}
		query = builder.From(adsOptsTable).Select("ad_id").Where(ex)
	}
	if asp.JoistType != "" {
		ex := goqu.Ex{"ad_id": query, "opt_type": AdOptType_JoistType, "opt_value_str": asp.JoistType}
		query = builder.From(adsOptsTable).Select("ad_id").Where(ex)
	}

	if asp.TranslitAddressDependentLocality != "" {
		ex := goqu.Ex{"ad_id": query, "opt_type": AdOptType_TranslitAddressDependentLocality, "opt_value_str": asp.TranslitAddressDependentLocality}
		query = builder.From(adsOptsTable).Select("ad_id").Where(ex)
	}

	if asp.TranslitAddressAddressStreet != "" {
		ex := goqu.Ex{"ad_id": query, "opt_type": AdOptType_TranslitAddressAddressStreet, "opt_value_str": asp.TranslitAddressAddressStreet}
		query = builder.From(adsOptsTable).Select("ad_id").Where(ex)
	}

	if asp.AddressHouseNumber != "" {
		ex := goqu.Ex{"ad_id": query, "opt_type": AdOptType_AddressHouseNumber, "opt_value_str": asp.AddressHouseNumber}
		query = builder.From(adsOptsTable).Select("ad_id").Where(ex)
	}

	if asp.TranslitAddressHouseNumber != "" {
		ex := goqu.Ex{"ad_id": query, "opt_type": AdOptType_TranslitAddressHouseNumber, "opt_value_str": asp.TranslitAddressHouseNumber}
		query = builder.From(adsOptsTable).Select("ad_id").Where(ex)
	}

	if asp.IsStudio {
		ex := goqu.Ex{"ad_id": query, "opt_type": AdOptType_IsStudio, "opt_value_bool": true}
		query = builder.From(adsOptsTable).Select("ad_id").Where(ex)
	} else if asp.RoomsCount > 0 {
		ex := goqu.Ex{"ad_id": query, "opt_type": AdOptType_RoomsCount, "opt_value_uint": asp.RoomsCount}
		query = builder.From(adsOptsTable).Select("ad_id").Where(ex)
	}

	if asp.HasFloor {
		if asp.Floor < 0 {
			ex := goqu.Ex{"ad_id": query, "opt_type": AddOptType_NegativeFloor, "opt_value_bool": true}
			query = builder.From(adsOptsTable).Select("ad_id").Where(ex)
		}
		ex := goqu.Ex{"ad_id": query, "opt_type": AddOptType_Floor, "opt_value_uint": asp.Floor}
		query = builder.From(adsOptsTable).Select("ad_id").Where(ex)
	}
	if asp.Elevator {
		ex := goqu.Ex{"ad_id": query, "opt_type": AdOptType_Elevator}
		query = builder.From(adsOptsTable).Select("ad_id").Where(goqu.And(ex, goqu.C("opt_value_str").IsNotNull()))
	}

	if asp.PriceMin > 0 || asp.PriceMax > 0 {
		query = as.buildUintMinMaxSubQuery(query, AdOptType_Price, asp.PriceMin, asp.PriceMax)
	}
	if asp.SpaceMin > 0 || asp.SpaceMax > 0 {
		query = as.buildUintMinMaxSubQuery(query, AdOptType_TotalSpace, asp.SpaceMin, asp.SpaceMax)
	}
	if asp.FloorMin > 0 || asp.FloorMax > 0 {
		query = as.buildUintMinMaxSubQuery(query, AddOptType_Floor, asp.FloorMin, asp.FloorMax)
	}
	if asp.LivingSpaceMin > 0 || asp.LivingSpaceMax > 0 {
		query = as.buildUintMinMaxSubQuery(query, AdOptType_LivingSpace, asp.LivingSpaceMin, asp.LivingSpaceMax)
	}
	if asp.KitchenSpaceMin > 0 || asp.KitchenSpaceMax > 0 {
		query = as.buildUintMinMaxSubQuery(query, AdOptType_KitchenSpace, asp.KitchenSpaceMin, asp.KitchenSpaceMax)
	}
	if asp.PriceForSqmMin > 0 || asp.PriceForSqmMax > 0 {
		query = as.buildUintMinMaxSubQuery(query, AdOptType_PriceSqM, asp.PriceForSqmMin, asp.PriceForSqmMax)
	}
	if asp.CeilingHeightMin > 0 || asp.CeilingHeightMax > 0 {
		query = as.buildUintMinMaxSubQuery(query, AdOptType_CeilingHeight, asp.CeilingHeightMin, asp.CeilingHeightMax)
	}
	if asp.TotalFloorsMin > 0 || asp.TotalFloorsMax > 0 {
		query = as.buildUintMinMaxSubQuery(query, AddOptType_TotalFloors, asp.TotalFloorsMin, asp.TotalFloorsMax)
	}
	if asp.BalconyCnt > 0 {
		query = as.buildUintMinMaxSubQuery(query, AdOptType_BalconyCnt, asp.BalconyCnt, 0)
	}
	if asp.LoggiaCnt > 0 {
		query = as.buildUintMinMaxSubQuery(query, AdOptType_LoggiaCnt, asp.LoggiaCnt, 0)
	}
	if asp.CombinedBathroomCnt > 0 {
		query = as.buildUintMinMaxSubQuery(query, AdOptType_CombinedBathroomCnt, asp.CombinedBathroomCnt, 0)
	}
	if asp.SeparateBathroomCnt > 0 {
		query = as.buildUintMinMaxSubQuery(query, AdOptType_SeparateBathroomCnt, asp.SeparateBathroomCnt, 0)
	}

	cntSql, _, err := query.Select(goqu.COUNT("*")).ToSQL()
	if err != nil {
		log.Err(err).Msgf("AdsSearchService#Search - can't build count sql query: %v", err)
		return "", "", err
	}

	log.Info().Msgf("AdsStore#buildSql - CNT_QUERY: %s", cntSql)

	adsSql, _, err := query.Select("ad_id").
		Limit(uint(asp.Limit)).
		Offset(uint((asp.Page - 1) * asp.Limit)).
		ToSQL()
	if err != nil {
		log.Err(err).Msgf("AdsSearchService#Search - can't build ads sql query: %v", err)
		return "", "", err
	}

	return cntSql, adsSql, nil
}

func (as *AdsStore) buildUintMinMaxSubQuery(query *goqu.SelectDataset, optType AdOptType, min int, max int) *goqu.SelectDataset {
	ex := goqu.Ex{
		"ad_id":    query,
		"opt_type": optType,
	}
	if min > 0 && max > 0 {
		ex["opt_value_uint"] = exp.Op{"between": goqu.Range(min, max)}
	} else if max > 0 {
		ex["opt_value_uint"] = exp.Op{"lte": max}
	} else if min > 0 {
		ex["opt_value_uint"] = exp.Op{"gte": min}
	}

	return as.goquBuilder.From(adsOptsTable).Select("ad_id").Where(ex)
}
