// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createBankStmt, err = db.PrepareContext(ctx, createBank); err != nil {
		return nil, fmt.Errorf("error preparing query CreateBank: %w", err)
	}
	if q.createCountryStmt, err = db.PrepareContext(ctx, createCountry); err != nil {
		return nil, fmt.Errorf("error preparing query CreateCountry: %w", err)
	}
	if q.deleteBankStmt, err = db.PrepareContext(ctx, deleteBank); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteBank: %w", err)
	}
	if q.getBankBySwiftCodeStmt, err = db.PrepareContext(ctx, getBankBySwiftCode); err != nil {
		return nil, fmt.Errorf("error preparing query GetBankBySwiftCode: %w", err)
	}
	if q.getBranchesByCountryISO2Stmt, err = db.PrepareContext(ctx, getBranchesByCountryISO2); err != nil {
		return nil, fmt.Errorf("error preparing query GetBranchesByCountryISO2: %w", err)
	}
	if q.getBranchesByHeadquarterIdStmt, err = db.PrepareContext(ctx, getBranchesByHeadquarterId); err != nil {
		return nil, fmt.Errorf("error preparing query GetBranchesByHeadquarterId: %w", err)
	}
	if q.getCountryByCountryISO2Stmt, err = db.PrepareContext(ctx, getCountryByCountryISO2); err != nil {
		return nil, fmt.Errorf("error preparing query GetCountryByCountryISO2: %w", err)
	}
	if q.setBranchHeadquarterStmt, err = db.PrepareContext(ctx, setBranchHeadquarter); err != nil {
		return nil, fmt.Errorf("error preparing query SetBranchHeadquarter: %w", err)
	}
	if q.updateBranchesHeadquarterStmt, err = db.PrepareContext(ctx, updateBranchesHeadquarter); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateBranchesHeadquarter: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createBankStmt != nil {
		if cerr := q.createBankStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createBankStmt: %w", cerr)
		}
	}
	if q.createCountryStmt != nil {
		if cerr := q.createCountryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createCountryStmt: %w", cerr)
		}
	}
	if q.deleteBankStmt != nil {
		if cerr := q.deleteBankStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteBankStmt: %w", cerr)
		}
	}
	if q.getBankBySwiftCodeStmt != nil {
		if cerr := q.getBankBySwiftCodeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getBankBySwiftCodeStmt: %w", cerr)
		}
	}
	if q.getBranchesByCountryISO2Stmt != nil {
		if cerr := q.getBranchesByCountryISO2Stmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getBranchesByCountryISO2Stmt: %w", cerr)
		}
	}
	if q.getBranchesByHeadquarterIdStmt != nil {
		if cerr := q.getBranchesByHeadquarterIdStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getBranchesByHeadquarterIdStmt: %w", cerr)
		}
	}
	if q.getCountryByCountryISO2Stmt != nil {
		if cerr := q.getCountryByCountryISO2Stmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getCountryByCountryISO2Stmt: %w", cerr)
		}
	}
	if q.setBranchHeadquarterStmt != nil {
		if cerr := q.setBranchHeadquarterStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing setBranchHeadquarterStmt: %w", cerr)
		}
	}
	if q.updateBranchesHeadquarterStmt != nil {
		if cerr := q.updateBranchesHeadquarterStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateBranchesHeadquarterStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                             DBTX
	tx                             *sql.Tx
	createBankStmt                 *sql.Stmt
	createCountryStmt              *sql.Stmt
	deleteBankStmt                 *sql.Stmt
	getBankBySwiftCodeStmt         *sql.Stmt
	getBranchesByCountryISO2Stmt   *sql.Stmt
	getBranchesByHeadquarterIdStmt *sql.Stmt
	getCountryByCountryISO2Stmt    *sql.Stmt
	setBranchHeadquarterStmt       *sql.Stmt
	updateBranchesHeadquarterStmt  *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                             tx,
		tx:                             tx,
		createBankStmt:                 q.createBankStmt,
		createCountryStmt:              q.createCountryStmt,
		deleteBankStmt:                 q.deleteBankStmt,
		getBankBySwiftCodeStmt:         q.getBankBySwiftCodeStmt,
		getBranchesByCountryISO2Stmt:   q.getBranchesByCountryISO2Stmt,
		getBranchesByHeadquarterIdStmt: q.getBranchesByHeadquarterIdStmt,
		getCountryByCountryISO2Stmt:    q.getCountryByCountryISO2Stmt,
		setBranchHeadquarterStmt:       q.setBranchHeadquarterStmt,
		updateBranchesHeadquarterStmt:  q.updateBranchesHeadquarterStmt,
	}
}
