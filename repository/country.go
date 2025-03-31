package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
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
	res, err := db.ExecContext(ctx, rawQuery, argsArr...)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	rowsIgnored := len(args) - int(rowsAffected)
	slog.Info(fmt.Sprintf("Executed: insert values into country table. %d rows affected", rowsAffected))
	if rowsIgnored > 0 {
		slog.Info(fmt.Sprintf("%d duplicate rows ignored", rowsIgnored))
	}
	return nil
}
