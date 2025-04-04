package schemas

type CreateSwiftEntry struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName" binding:"required"`
	CountryIso2   string `json:"countryISO2" binding:"required"`
	SwiftCode     string `json:"swiftCode" binding:"required"`
	CountryName   string `json:"countryName" binding:"required"`
	IsHeadquarter *bool  `json:"isHeadquarter" binding:"required"`
}

type GetBranchEntry struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryIso2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type GetHeadquarterEntry struct {
	Address       string           `json:"address"`
	BankName      string           `json:"bankName"`
	CountryIso2   string           `json:"countryISO2"`
	CountryName   string           `json:"countryName"`
	IsHeadquarter bool             `json:"isHeadquarter"`
	SwiftCode     string           `json:"swiftCode"`
	Branches      []GetBranchEntry `json:"branches"`
}
