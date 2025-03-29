package repository

import (
	"context"
	"database/sql"

	repository "swift-restful/repository/sqlc"
)

type CreateBankUpdateRelationsParams struct {
	BankParams    repository.CreateBankParams
	CountryParams repository.CreateCountryParams
}

func CreateBankUpdateRelations(ctx context.Context, db *sql.DB, queries *repository.Queries, params CreateBankUpdateRelationsParams) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)
	var res sql.Result
	if _, err = qtx.CreateCountry(ctx, params.CountryParams); err != nil {
		// log.Printf("unable to create Country entry %s", params.CountryParams)
		return err
	}
	if res, err = qtx.CreateBank(ctx, params.BankParams); err != nil {
		// log.Printf("unable to create Bank entry %s", params.BankParams)
		return err
	}
	bankId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	if params.BankParams.IsHeadquarter {
		if err := qtx.UpdateBranchesHeadquarter(ctx, repository.UpdateBranchesHeadquarterParams{
			HeadquarterId: sql.NullInt32{Int32: int32(bankId), Valid: true},
			SwiftCode:     params.BankParams.SwiftCode}); err != nil {
			return err
		}
	} else {
		if _, err = qtx.SetBranchHeadquarter(ctx, int32(bankId)); err != nil {
			return err
		}
	}
	return tx.Commit()
}
