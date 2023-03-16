package ads

type Ad struct {
	Id     int     `json:"id" db:"id"`
	UserId int     `json:"userId" db:"user_id"`
	Opts   []AdOpt `json:"opts"`
}

type AdOpt struct {
	TypeName   string `json:"typeName"`
	ValueType  string `json:"valueType"`
	ValueStr   string `json:"valStr"`
	ValueInt64 int64  `json:"valInt"`
	ValueBool  bool   `json:"valBool"`
	ValueBlob  string `json:"valBlob"`
}

type AdsSearchParams struct {
	Page                             int
	Limit                            int
	City                             string
	PriceMin                         int
	PriceMax                         int
	RoomsCount                       int
	IsStudio                         bool
	SpaceMin                         int
	SpaceMax                         int
	HasFloor                         bool
	Floor                            int
	FloorMin                         int
	FloorMax                         int
	LivingSpaceMin                   int
	LivingSpaceMax                   int
	KitchenSpaceMin                  int
	KitchenSpaceMax                  int
	PriceForSqmMin                   int
	PriceForSqmMax                   int
	CeilingHeightMin                 int
	CeilingHeightMax                 int
	TotalFloorsMin                   int
	TotalFloorsMax                   int
	ConstructionYear                 int
	WindowView                       string
	Renovation                       string
	BuildingType                     string
	JoistType                        string
	Elevator                         bool
	BalconyCnt                       int
	LoggiaCnt                        int
	CombinedBathroomCnt              int
	SeparateBathroomCnt              int
	TranslitAddressDependentLocality string
	TranslitAddressAddressStreet     string
	AddressHouseNumber               string
	TranslitAddressHouseNumber       string
	TranslitAddressCity              string
}

func (asp AdsSearchParams) HasFilters() bool {
	return asp.City != "" ||
		asp.PriceMin > 0 ||
		asp.PriceMax > 0 ||
		asp.RoomsCount > 0 ||
		asp.IsStudio ||
		asp.SpaceMin > 0 ||
		asp.SpaceMax > 0 ||
		asp.HasFloor ||
		asp.FloorMin > 0 ||
		asp.FloorMax > 0 ||
		asp.LivingSpaceMin > 0 ||
		asp.LivingSpaceMax > 0 ||
		asp.KitchenSpaceMin > 0 ||
		asp.KitchenSpaceMax > 0 ||
		asp.PriceForSqmMin > 0 ||
		asp.PriceForSqmMax > 0 ||
		asp.CeilingHeightMin > 0 ||
		asp.CeilingHeightMax > 0 ||
		asp.TotalFloorsMin > 0 ||
		asp.TotalFloorsMax > 0 ||
		asp.ConstructionYear > 0 ||
		asp.WindowView != "" ||
		asp.Renovation != "" ||
		asp.BuildingType != "" ||
		asp.JoistType != "" ||
		asp.Elevator ||
		asp.LoggiaCnt > 0 ||
		asp.BalconyCnt > 0 ||
		asp.CombinedBathroomCnt > 0 ||
		asp.SeparateBathroomCnt > 0 ||
		asp.TranslitAddressDependentLocality != "" ||
		asp.TranslitAddressAddressStreet != "" ||
		asp.AddressHouseNumber != "" ||
		asp.TranslitAddressHouseNumber != "" ||
		asp.TranslitAddressCity != ""
}
