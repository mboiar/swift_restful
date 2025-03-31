// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: bank.sql

package repository

import (
	"context"
	"database/sql"
)

const createBank = `-- name: CreateBank :execresult
INSERT INTO bank(
    ` + "`" + `address` + "`" + `,
    ` + "`" + `name` + "`" + `,
    ` + "`" + `country_ISO2` + "`" + `,
    ` + "`" + `is_headquarter` + "`" + `,
    ` + "`" + `swift_code` + "`" + `
) VALUES (?, ?, ?, ?, ?)
`

type CreateBankParams struct {
	Address       sql.NullString `json:"address"`
	Name          string         `json:"name"`
	CountryIso2   string         `json:"country_iso2"`
	IsHeadquarter bool           `json:"is_headquarter"`
	SwiftCode     string         `json:"swift_code"`
}

func (q *Queries) CreateBank(ctx context.Context, arg CreateBankParams) (sql.Result, error) {
	return q.exec(ctx, q.createBankStmt, createBank,
		arg.Address,
		arg.Name,
		arg.CountryIso2,
		arg.IsHeadquarter,
		arg.SwiftCode,
	)
}

const createBankBulk = `-- name: CreateBankBulk :copyfrom
INSERT INTO bank(
    ` + "`" + `address` + "`" + `,
    ` + "`" + `name` + "`" + `,
    ` + "`" + `country_ISO2` + "`" + `,
    ` + "`" + `is_headquarter` + "`" + `,
    ` + "`" + `swift_code` + "`" + `
) VALUES (
    ?, ?, ?, ?, ?
)
`

type CreateBankBulkParams struct {
	Address       sql.NullString `json:"address"`
	Name          string         `json:"name"`
	CountryIso2   string         `json:"country_iso2"`
	IsHeadquarter bool           `json:"is_headquarter"`
	SwiftCode     string         `json:"swift_code"`
}

const deleteBank = `-- name: DeleteBank :exec
DELETE FROM bank
WHERE swift_code = ?
`

func (q *Queries) DeleteBank(ctx context.Context, swiftCode string) error {
	_, err := q.exec(ctx, q.deleteBankStmt, deleteBank, swiftCode)
	return err
}

const getBankBySwiftCode = `-- name: GetBankBySwiftCode :one
SELECT bank.id, bank.name, bank.address, bank.is_headquarter, bank.swift_code, bank.country_iso2, country.name FROM bank
INNER JOIN country
ON bank.` + "`" + `country_ISO2` + "`" + ` = country.` + "`" + `ISO2` + "`" + `
WHERE swift_code = ? LIMIT 1
`

type GetBankBySwiftCodeRow struct {
	ID            int32          `json:"id"`
	Name          string         `json:"name"`
	Address       sql.NullString `json:"address"`
	IsHeadquarter bool           `json:"is_headquarter"`
	SwiftCode     string         `json:"swift_code"`
	CountryIso2   string         `json:"country_iso2"`
	Name_2        string         `json:"name_2"`
}

func (q *Queries) GetBankBySwiftCode(ctx context.Context, swiftCode string) (GetBankBySwiftCodeRow, error) {
	row := q.queryRow(ctx, q.getBankBySwiftCodeStmt, getBankBySwiftCode, swiftCode)
	var i GetBankBySwiftCodeRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.IsHeadquarter,
		&i.SwiftCode,
		&i.CountryIso2,
		&i.Name_2,
	)
	return i, err
}

const getBranchesByCountryISO2 = `-- name: GetBranchesByCountryISO2 :many
SELECT id, name, address, is_headquarter, swift_code, country_iso2 FROM bank
WHERE ` + "`" + `country_ISO2` + "`" + ` = ? LIMIT ?
`

type GetBranchesByCountryISO2Params struct {
	CountryIso2 string `json:"country_iso2"`
	Limit       int32  `json:"limit"`
}

func (q *Queries) GetBranchesByCountryISO2(ctx context.Context, arg GetBranchesByCountryISO2Params) ([]Bank, error) {
	rows, err := q.query(ctx, q.getBranchesByCountryISO2Stmt, getBranchesByCountryISO2, arg.CountryIso2, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Bank
	for rows.Next() {
		var i Bank
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.IsHeadquarter,
			&i.SwiftCode,
			&i.CountryIso2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBranchesBySwiftCode = `-- name: GetBranchesBySwiftCode :many
SELECT id, name, address, is_headquarter, swift_code, country_iso2 from bank
WHERE LEFT(bank.swift_code, 8) = LEFT(?, 8) LIMIT ?
`

type GetBranchesBySwiftCodeParams struct {
	LEFT  string `json:"LEFT"`
	Limit int32  `json:"limit"`
}

func (q *Queries) GetBranchesBySwiftCode(ctx context.Context, arg GetBranchesBySwiftCodeParams) ([]Bank, error) {
	rows, err := q.query(ctx, q.getBranchesBySwiftCodeStmt, getBranchesBySwiftCode, arg.LEFT, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Bank
	for rows.Next() {
		var i Bank
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Address,
			&i.IsHeadquarter,
			&i.SwiftCode,
			&i.CountryIso2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
