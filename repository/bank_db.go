package repository

import (
	"context"
	"database/sql"

	repository "swift-restful/repository/sqlc"
)

func CreateBankUpdateRelations(ctx context.Context, db *sql.DB, queries *repository.Queries, params repository.CreateBankParams) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)
	var res sql.Result
	if res, err = qtx.CreateBank(ctx, params); err != nil {
		return err
	}
	bankId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	if err = qtx.CreateCountry(ctx, repository.CreateCountryParams{
		Iso2: params.CountryISO2,
		Name: params.CountryName}); err != nil {
		return err
	}
	if params.IsHeadquarter {
		if err := qtx.UpdateBranchesHeadquarter(ctx, repository.UpdateBranchesHeadquarterParams{
			HeadquarterId: sql.NullInt32{Int32: int32(bankId), Valid: true},
			SwiftCode:     params.SwiftCode}); err != nil {
			return err
		}
	} else {
		if _, err = qtx.SetBranchHeadquarter(ctx, int32(bankId)); err != nil {
			return err
		}
	}
	return tx.Commit()
}
