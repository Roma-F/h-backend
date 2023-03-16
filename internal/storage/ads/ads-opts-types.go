package ads

type AdOptType int

const (
	AddOptType_Floor         AdOptType = 2
	AddOptType_TotalFloors   AdOptType = 3
	AddOptType_NegativeFloor AdOptType = 39

	AdOptType_TotalSpace   AdOptType = 4
	AdOptType_LivingSpace  AdOptType = 5
	AdOptType_KitchenSpace AdOptType = 12

	AdOptType_Price    AdOptType = 13
	AdOptType_PriceSqM AdOptType = 36

	AdOptType_RoomsCount AdOptType = 50
	AdOptType_IsStudio   AdOptType = 51

	AdOptType_ConstructionYear                 AdOptType = 6
	AdOptType_Phone                            AdOptType = 9
	AdOptType_BuildingType                     AdOptType = 15
	AdOptType_CeilingHeight                    AdOptType = 18
	AdOptType_Renovation                       AdOptType = 20
	AdOptType_WindowView                       AdOptType = 21
	AdOptType_JoistType                        AdOptType = 23
	AdOptType_Elevator                         AdOptType = 38
	AdOptType_AddressCity                      AdOptType = 43
	AdOptType_BalconyCnt                       AdOptType = 52
	AdOptType_LoggiaCnt                        AdOptType = 53
	AdOptType_CombinedBathroomCnt              AdOptType = 54
	AdOptType_SeparateBathroomCnt              AdOptType = 55
	AdOptType_AddressHouseNumber               AdOptType = 46
	AdOptType_TranslitAddressCity              AdOptType = 58
	AdOptType_TranslitAddressDependentLocality AdOptType = 59
	AdOptType_TranslitAddressAddressStreet     AdOptType = 60
	AdOptType_TranslitAddressHouseNumber       AdOptType = 61
)
