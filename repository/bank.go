package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	repository "swift-restful/repository/sqlc"
)

func InsertBankMultiRow(ctx context.Context, args []repository.CreateBankBulkParams, db *sql.DB, skipDuplicates bool) error {
	placeholderArr := make([]string, len(args))
	argsArr := make([]interface{}, len(args)*4)
	for i, arg := range args {
		placeholderArr[i] = "(?, ?, ?, ?)"
		argsArr[i*4] = arg.Address
		argsArr[i*4+1] = arg.Name
		argsArr[i*4+2] = arg.CountryIso2
		argsArr[i*4+3] = arg.SwiftCode
	}
	var queryStr string
	if skipDuplicates {
		queryStr = "INSERT IGNORE INTO bank (address, name, country_iso2, swift_code) VALUES %s"
	} else {
		queryStr = "INSERT INTO bank (address, name, country_iso2, swift_code) VALUES %s"
	}
	rawQuery := fmt.Sprintf(queryStr, strings.Join(placeholderArr, ","))
	_, err := db.ExecContext(ctx, rawQuery, argsArr...)
	return err
}
