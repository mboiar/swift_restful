package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	repository "swift-restful/repository/sqlc"
)

func InsertCountryMultiRow(ctx context.Context, args []repository.CreateCountryBulkParams, db *sql.DB) error {
	placeholderArr := make([]string, len(args))
	argsArr := make([]interface{}, len(args)*2)
	for i, arg := range args {
		placeholderArr[i] = "(?, ?)"
		argsArr[i*2] = arg.Iso2
		argsArr[i*2+1] = arg.Name
	}
	rawQuery := fmt.Sprintf("INSERT IGNORE INTO country (iso2, name) VALUES %s", strings.Join(placeholderArr, ","))
	_, err := db.ExecContext(ctx, rawQuery, argsArr...)
	return err
}
