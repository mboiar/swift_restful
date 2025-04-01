// Repository implements SWIFT database and routines for interacting with it.
package repository

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	repository "swift-restful/repository/sqlc"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// SetupDB returns database and sqlc Queries instance configured based on provided .env file.
// If provided, configuration file should define DB_USER, DB_PASSWORD, DB_NAME, DB_HOST and DB_PORT environment variables.
func SetupDB(dbcfg_path *string) (*sql.DB, *repository.Queries, error) {
	if dbcfg_path != nil && *dbcfg_path != "" {
		err := godotenv.Load(*dbcfg_path)
		if err != nil {
			slog.Error("error loading config " + *dbcfg_path)
			return nil, nil, err
		}
	}
	user := os.Getenv("MYSQL_USER")
	passwd := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	addr := host + ":" + port
	dbname := os.Getenv("MYSQL_DATABASE")
	if dbname == "" || host == "" || port == "" {
		return nil, nil, fmt.Errorf("MYSQL_DATABASE (%s), DB_HOST (%s) and DB_PORT (%s) can not be empty: make sure appropriate environment variables are set", dbname, host, port)
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
