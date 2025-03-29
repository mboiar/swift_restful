package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	repository "swift-restful/repository/sqlc"

	"github.com/go-sql-driver/mysql"
)

const (
	MYSQL_ER_DUP_ENTRY = 1062
)

func SetupDB() (*sql.DB, *repository.Queries, error) {
	user := os.Getenv("DB_USER")
	passwd := os.Getenv("DB_PASSWORD")
	addr := os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	if dbname == "" || addr == "" {
		return nil, nil, fmt.Errorf("DB name (%s) and address (%s) can not be empty", dbname, addr)
	}
	cfg := mysql.Config{
		User:   user,
		Passwd: passwd,
		Net:    "tcp",
		Addr:   addr,
		DBName: dbname,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, nil, err
	}
	queries := repository.New(db)

	log.Println("Connected to database", cfg.DBName)
	return db, queries, nil
}
