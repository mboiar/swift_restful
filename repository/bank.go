package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	repository "swift-restful/repository/sqlc"
)

func InsertBankMultiRow(ctx context.Context, args []repository.CreateBankBulkParams, db *sql.DB, skipDuplicates bool) error {
	nRows := len(args)
	placeholderArr := make([]string, nRows)
	argsArr := make([]interface{}, nRows*5)
	for i, arg := range args {
		placeholderArr[i] = "(?, ?, ?, ?, ?)"
		argsArr[i*5] = arg.Address
		argsArr[i*5+1] = arg.Name
		argsArr[i*5+2] = arg.CountryIso2
		argsArr[i*5+3] = arg.IsHeadquarter
		argsArr[i*5+4] = arg.SwiftCode
	}
	var queryStr string
	if skipDuplicates {
		queryStr = "INSERT IGNORE INTO bank (address, name, country_iso2, is_headquarter, swift_code) VALUES %s"
	} else {
		queryStr = "INSERT INTO bank (address, name, country_iso2, is_headquarter, swift_code) VALUES %s"
	}
	rawQuery := fmt.Sprintf(queryStr, strings.Join(placeholderArr, ","))
	res, err := db.ExecContext(ctx, rawQuery, argsArr...)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	rowsIgnored := nRows - int(rowsAffected)
	slog.Info(fmt.Sprintf("Executed: insert values into bank table. %d rows affected", rowsAffected))
	if rowsIgnored > 0 {
		slog.Info(fmt.Sprintf("%d duplicate rows ignored", rowsIgnored))
	}

	return nil
}
