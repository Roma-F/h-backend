package cities

import (
	cities_core "lbr-admin-backend/internal/core/cities"
)

type CitiesFether interface {
	Search(csp cities_core.CitiesSearchParams) ([]cities_core.Cities, int, error)
	GetPageCities(limit, offset int) ([]cities_core.Cities, error)
	GetCitiesPagesCount(limit int) (int, error)
}
