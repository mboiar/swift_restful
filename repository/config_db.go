package repository

import (
	"database/sql"
	"fmt"
	"os"

	repository "swift-restful/repository/sqlc"

	"github.com/go-sql-driver/mysql"
)

func SetupDB() (*repository.Queries, error) {
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Net:    "tcp",
		Addr:   os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),
	}

	conn, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	db := repository.New(conn)

	fmt.Println("Connected!")
	return db, nil
}
