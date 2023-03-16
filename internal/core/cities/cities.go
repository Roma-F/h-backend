package cities

type Cities struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Translit string `json:"translit" db:"translit"`
}

type CitiesSearchParams struct {
	Page     int
	Limit    int
	Name     string
	Translit string
}

func (csp CitiesSearchParams) HasFilters() bool {
	return csp.Name != "" || csp.Translit != ""
}
