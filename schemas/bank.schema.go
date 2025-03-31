package schemas

import "database/sql"

type CreateSwiftEntry struct {
	Address       sql.NullString `json:"address" binding:"required"`
	BankName      string         `json:"bankName" binding:"required"`
	CountryIso2   string         `json:"countryISO2" binding:"required"`
	SwiftCode     string         `json:"swiftCode" binding:"required"`
	CountryName   string         `json:"countryName" binding:"required"`
	IsHeadquarter bool           `json:"isHeadaquarter" binding:"required"`
}

type GetBranchEntry struct {
	Address       sql.NullString `json:"address"`
	BankName      string         `json:"bankName"`
	CountryIso2   string         `json:"countryISO2"`
	IsHeadquarter bool           `json:"isHeadaquarter"`
	SwiftCode     string         `json:"swiftCode"`
}

type GetHeadquarterEntry struct {
	Address       sql.NullString   `json:"address"`
	BankName      string           `json:"bankName"`
	CountryIso2   string           `json:"countryISO2"`
	CountryName   string           `json:"countryName"`
	IsHeadquarter bool             `json:"isHeadaquarter"`
	SwiftCode     string           `json:"swiftCode"`
	Branches      []GetBranchEntry `json:"branches"`
}
