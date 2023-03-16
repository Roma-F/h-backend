package ads

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	core_ads "homeho-backend/internal/core/ads"
	"homeho-backend/internal/storage"
	"homeho-backend/pkg/dotdb"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

//go:embed queries.sql
var queries string

var ErrAdNotFound = errors.New("ad is not found")
var ErrNoSuchPageExists = errors.New("no such ads page exists")
var ErrNoPhone = errors.New("ad no phone")

var PrivateOptsTypes = []AdOptType{AdOptType_Phone}

type AdsStore struct {
	sqlxDb      *sqlx.DB
	db          *dotdb.DotDB
	goquBuilder goqu.DialectWrapper
	optTypes    []storage.DbOptType
}

func NewAdsStore(db *sqlx.DB, goquBuilder goqu.DialectWrapper) (*AdsStore, error) {
	dotDb, err := dotdb.NewDotDB(db, queries)
	if err != nil {
		return nil, fmt.Errorf("can't create new AdsStore, initing dotdb: %v", err)
	}

	var AdOptTypes []storage.DbOptType // todo: extract outside and pass here as params
	err = dotDb.GetStructsSlice("get-ad-opt-types", &AdOptTypes)
	if err != nil {
		return nil, fmt.Errorf("can't create new AdsStore, reading adOptTypes: %v", err)
	}

	return &AdsStore{
		sqlxDb:      db,
		db:          dotDb,
		goquBuilder: goquBuilder,
		optTypes:    AdOptTypes,
	}, nil
}

func (as *AdsStore) GetAdById(id int) (*core_ads.Ad, error) {
	adRow, err := as.db.GetRow("fetch-ad", id)
	if err != nil {
		return nil, fmt.Errorf("can't fetch ad: %v", err)
	}

	var userId int
	err = adRow.Scan(&userId)
	if err == sql.ErrNoRows {
		return nil, ErrAdNotFound
	} else if err != nil {
		return nil, fmt.Errorf("error scanning userId for ad: %v", err)
	}

	ad := core_ads.Ad{
		Id:     id,
		UserId: userId,
	}

	rows, err := as.db.GetRows("fetch-ads-opts", []int{id}, PrivateOptsTypes)
	if err != nil {
		return nil, fmt.Errorf("error loading opts: %v", err)
	}
	defer rows.Close()

	opts, err := as.gatherOpts(rows, 1)
	if err != nil {
		return nil, fmt.Errorf("error gathering opts: %v", err)
	}
	ad.Opts = *opts[ad.Id]
	return &ad, nil
}

func (as *AdsStore) GetAdsByIds(ids []int) ([]core_ads.Ad, error) {
	rows, err := as.db.GetRows("fetch-ads-opts", ids, PrivateOptsTypes)
	if err != nil {
		return nil, fmt.Errorf("error loading opts: %v", err)
	}
	defer rows.Close()

	adsOpts, err := as.gatherOpts(rows, len(ids))
	if err != nil {
		return nil, fmt.Errorf("error gathering opts: %v", err)
	}

	ads := make([]core_ads.Ad, len(ids))
	for idx, id := range ids {
		opts, ok := adsOpts[id]
		if !ok {
			continue
		}
		ads[idx] = core_ads.Ad{
			Id:   id,
			Opts: *opts,
		}
	}

	return ads, nil
}

func (as *AdsStore) GetPageAds(limit, offset int) ([]core_ads.Ad, error) {
	var ads []core_ads.Ad
	fmt.Println(limit, offset, "IIIIIIIII")
	err := as.db.GetStructsSlice("fetch-page-ads", &ads, limit, offset)
	if err != nil {
		return nil, err
	} else if len(ads) == 0 {
		return nil, ErrNoSuchPageExists
	}

	adsIds := make([]int, len(ads))
	for i, ad := range ads {
		adsIds[i] = ad.Id
	}

	rows, err := as.db.GetRows("fetch-ads-opts", adsIds, PrivateOptsTypes)
	if err != nil {
		return nil, fmt.Errorf("error loading opts: %v", err)
	}
	defer rows.Close()

	adsOpts, err := as.gatherOpts(rows, len(ads))
	if err != nil {
		return nil, fmt.Errorf("error gathering opts: %v", err)
	}

	for idx := range ads {
		opts, ok := adsOpts[ads[idx].Id]
		if !ok {
			continue
		}
		ads[idx].Opts = *opts
	}

	return ads, nil
}

func (as *AdsStore) GetRecommended(limit int, cityTranslit string) ([]core_ads.Ad, error) {
	var ads []core_ads.Ad
	fmt.Println(cityTranslit, "IIIIIIIII")
	err := as.db.GetStructsSlice("get-recommended-city-ads", &ads, cityTranslit, limit)
	if err != nil {
		return nil, err
	}
	adsIds := make([]int, len(ads))
	for i, ad := range ads {
		adsIds[i] = ad.Id
	}
	rows, err := as.db.GetRows("fetch-ads-opts", adsIds, PrivateOptsTypes)
	if err != nil {
		return nil, fmt.Errorf("error loading opts: %v", err)
	}
	defer rows.Close()

	adsOpts, err := as.gatherOpts(rows, len(ads))
	if err != nil {
		return nil, fmt.Errorf("error gathering opts: %v", err)
	}

	for idx := range ads {
		opts, ok := adsOpts[ads[idx].Id]
		if !ok {
			continue
		}
		ads[idx].Opts = *opts
	}
	return ads, nil
}

func (as *AdsStore) GetAdsPagesCount(limit int) (int, error) {
	row, err := as.db.GetRow("fetch-ads-total-pages-count", limit)
	if err != nil {
		return 0, fmt.Errorf("error loading ads total pages: %v", err)
	}

	var cnt int
	err = row.Scan(&cnt)
	if err != nil {
		return 0, fmt.Errorf("error scanning ads total pages: %v", err)
	}

	return cnt, nil
}

func (as *AdsStore) GetAdPhone(adId int) (string, error) {
	row, err := as.db.GetRow("fetch-ad-string-opt", adId, AdOptType_Phone)
	if err == sql.ErrNoRows {
		return "", ErrNoPhone
	} else if err != nil {
		return "", err
	}

	var phone sql.NullString
	row.Scan(&phone)
	if !phone.Valid {
		return "", ErrNoPhone
	}

	return phone.String, nil
}

func (as *AdsStore) gatherOpts(rows *sqlx.Rows, optsSetCount int) (map[int]*[]core_ads.AdOpt, error) {
	adsOpts := make(map[int]*[]core_ads.AdOpt, optsSetCount)
	for rows.Next() {
		var adId int
		var optType int
		var valStr sql.NullString
		var valInt sql.NullInt64
		var valBool sql.NullBool
		var valBlob sql.RawBytes

		if err := rows.Scan(&adId, &optType, &valStr, &valInt, &valBool, &valBlob); err != nil {
			return nil, fmt.Errorf("error scanning opt: %v", err)
		}
		dbOptType, found := as.findDbOptType(optType)
		if !found {
			log.Warn().Msgf("unknown ad opt id %d\n", optType)
			continue
		}

		opt := core_ads.AdOpt{
			TypeName:  dbOptType.Name,
			ValueType: storage.OptValueTypeNames[dbOptType.ValueType],
		}

		switch dbOptType.ValueType {
		case storage.OVT_UINT:
			if !valInt.Valid {
				log.Warn().Msgf("got opt with null value, adId - %d, optType -  %d", adId, optType)
				continue
			}
			opt.ValueInt64 = valInt.Int64
		case storage.OVT_STRING:
			if !valStr.Valid {
				log.Warn().Msgf("got opt with null value, adId - %d, optType -  %d", adId, optType)
				continue
			}
			opt.ValueStr = valStr.String
		case storage.OVT_BOOLEAN:
			if !valBool.Valid {
				log.Warn().Msgf("got opt with null value, adId - %d, optType -  %d", adId, optType)
				continue
			}
			opt.ValueBool = valBool.Bool
		case storage.OVT_BLOB:
			if len(valBlob) == 0 {
				log.Warn().Msgf("got opt with null value, adId - %d, optType -  %d", adId, optType)
				continue
			}
			opt.ValueBlob = string(valBlob)
		default:
			log.Warn().Msgf("error - unknown opt value type %d, adId - %d", optType, adId)
			continue
		}

		valOptsSet, ok := adsOpts[adId]
		if !ok {
			optSet := make([]core_ads.AdOpt, 0, len(as.optTypes))
			optSet = append(optSet, opt)
			adsOpts[adId] = &optSet
		} else {
			*valOptsSet = append(*valOptsSet, opt)
		}
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating rows: %v", rows.Err())
	}
	return adsOpts, nil
}

func (as *AdsStore) findDbOptType(optType int) (storage.DbOptType, bool) {
	for _, ot := range as.optTypes {
		if ot.OptType == optType {
			return ot, true
		}
	}
	return storage.DbOptType{}, false
}
