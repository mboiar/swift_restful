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
) VALUES (
    ?, ?, ?, ?, ?
)
`

type CreateBankParams struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

func (q *Queries) CreateBank(ctx context.Context, arg CreateBankParams) (sql.Result, error) {
	return q.exec(ctx, q.createBankStmt, createBank,
		arg.Address,
		arg.BankName,
		arg.CountryISO2,
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
	Address       string `json:"address"`
	Name          string `json:"name"`
	CountryIso2   string `json:"country_iso2"`
	IsHeadquarter bool   `json:"is_headquarter"`
	SwiftCode     string `json:"swift_code"`
}

const deleteBank = `-- name: DeleteBank :exec
DELETE FROM bank
WHERE swift_code = ?
`

func (q *Queries) DeleteBank(ctx context.Context, swiftcode string) error {
	_, err := q.exec(ctx, q.deleteBankStmt, deleteBank, swiftcode)
	return err
}

const getBankBySwiftCode = `-- name: GetBankBySwiftCode :one
SELECT bank.id, bank.name, bank.address, bank.swift_code, bank.country_iso2, bank.is_headquarter, bank.headquarter_id, country.name FROM bank
INNER JOIN country
ON bank.` + "`" + `country_ISO2` + "`" + ` = country.` + "`" + `ISO2` + "`" + `
WHERE swift_code = ? LIMIT 1
`

type GetBankBySwiftCodeRow struct {
	ID            int32         `json:"id"`
	Name          string        `json:"name"`
	Address       string        `json:"address"`
	SwiftCode     string        `json:"swift_code"`
	CountryIso2   string        `json:"country_iso2"`
	IsHeadquarter bool          `json:"is_headquarter"`
	HeadquarterID sql.NullInt32 `json:"headquarter_id"`
	Name_2        string        `json:"name_2"`
}

func (q *Queries) GetBankBySwiftCode(ctx context.Context, swiftcode string) (GetBankBySwiftCodeRow, error) {
	row := q.queryRow(ctx, q.getBankBySwiftCodeStmt, getBankBySwiftCode, swiftcode)
	var i GetBankBySwiftCodeRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Address,
		&i.SwiftCode,
		&i.CountryIso2,
		&i.IsHeadquarter,
		&i.HeadquarterID,
		&i.Name_2,
	)
	return i, err
}

const getBranchesByCountryISO2 = `-- name: GetBranchesByCountryISO2 :many
SELECT id, name, address, swift_code, country_iso2, is_headquarter, headquarter_id FROM bank
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
			&i.SwiftCode,
			&i.CountryIso2,
			&i.IsHeadquarter,
			&i.HeadquarterID,
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

const getBranchesByHeadquarterId = `-- name: GetBranchesByHeadquarterId :many
SELECT id, name, address, swift_code, country_iso2, is_headquarter, headquarter_id from bank
WHERE ` + "`" + `headquarter_id` + "`" + ` = ? LIMIT ?
`

type GetBranchesByHeadquarterIdParams struct {
	HeadquarterID sql.NullInt32 `json:"headquarter_id"`
	Limit         int32         `json:"limit"`
}

func (q *Queries) GetBranchesByHeadquarterId(ctx context.Context, arg GetBranchesByHeadquarterIdParams) ([]Bank, error) {
	rows, err := q.query(ctx, q.getBranchesByHeadquarterIdStmt, getBranchesByHeadquarterId, arg.HeadquarterID, arg.Limit)
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
			&i.SwiftCode,
			&i.CountryIso2,
			&i.IsHeadquarter,
			&i.HeadquarterID,
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

const setBranchHeadquarter = `-- name: SetBranchHeadquarter :execresult
UPDATE bank AS branch
INNER JOIN bank AS headquarter
ON LEFT(branch.swift_code, 8) = LEFT(headquarter.swift_code, 8)
SET
branch.headquarter_id = headquarter.id
WHERE headquarter.is_headquarter AND branch.id = ?
`

func (q *Queries) SetBranchHeadquarter(ctx context.Context, id int32) (sql.Result, error) {
	return q.exec(ctx, q.setBranchHeadquarterStmt, setBranchHeadquarter, id)
}

const updateBranchesHeadquarter = `-- name: UpdateBranchesHeadquarter :exec
UPDATE bank
SET
` + "`" + `headquarter_id` + "`" + ` = ?
WHERE LEFT(` + "`" + `swift_code` + "`" + `, 8) = LEFT(?, 8) AND NOT ` + "`" + `is_headquarter` + "`" + `
`

type UpdateBranchesHeadquarterParams struct {
	HeadquarterId sql.NullInt32 `json:"headquarterId"`
	SwiftCode     string        `json:"swiftCode"`
}

func (q *Queries) UpdateBranchesHeadquarter(ctx context.Context, arg UpdateBranchesHeadquarterParams) error {
	_, err := q.exec(ctx, q.updateBranchesHeadquarterStmt, updateBranchesHeadquarter, arg.HeadquarterId, arg.SwiftCode)
	return err
}
