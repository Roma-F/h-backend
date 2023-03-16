package cities

import (
	storage_cities "lbr-backend/internal/storage/cities"
)

type CitiesController struct {
	citiesStore *storage_cities.CitiesStore
}

func NewCitiesController(citiesStore *storage_cities.CitiesStore) *CitiesController {
	return &CitiesController{
		citiesStore: citiesStore,
	}
}
