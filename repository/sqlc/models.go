// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repository

import (
	"database/sql"
)

type Bank struct {
	ID            int32          `json:"id"`
	Name          string         `json:"name"`
	Address       sql.NullString `json:"address"`
	IsHeadquarter bool           `json:"is_headquarter"`
	SwiftCode     string         `json:"swift_code"`
	CountryIso2   string         `json:"country_iso2"`
}

type Country struct {
	Iso2 string `json:"iso2"`
	Name string `json:"name"`
}
