package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	repository "swift-restful/repository/sqlc"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	tmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
)

type parserTestCase struct {
	parserName        string
	data              string
	expected_err      error
	expected_db_count uint
}

// createTestDB returns an instance of database for testing purposes
func createTestDB() (*sql.DB, *repository.Queries, func(), error) {
	ctx := context.Background()
	// create MySQL instance with docker
	container, err := tmysql.Run(ctx,
		"mysql:latest",
		tmysql.WithDatabase("test_db"),
		tmysql.WithUsername("root"),
		tmysql.WithPassword("password"),
		tmysql.WithScripts(filepath.Join("testdata", "schema.sql")),
	)
	cleanup := func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to start MySQL container: %w", err)
	}

	// connect to MySQL instance
	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		cleanup()
		return nil, nil, cleanup, err
	}
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, nil, cleanup, fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	queries := repository.New(db)
	return db, queries, cleanup, nil
}

func resetTestDB(db *sql.DB) error {
	_, err := db.Exec(`TRUNCATE TABLE bank`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`SET FOREIGN_KEY_CHECKS = 0`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`TRUNCATE TABLE country`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`SET FOREIGN_KEY_CHECKS = 1`)
	return err
}

// writeTempFile creates a temporary test file with sample data.
func writeTempFile(content string, extension string) (string, error) {
	tmpFile, err := os.CreateTemp("", "testfile-*."+extension)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

// TestParse verifies that the parser correctly reads and inserts data.
func TestParse(t *testing.T) {
	db, q, cleanup, err := createTestDB()
	if err != nil {
		t.Fatalf("Failed to create test DB: %v", err)
	}
	defer cleanup()

	parser1 := SwiftParser{params: SwiftParserParams{ // skips duplicates
		db:             db,
		queries:        q,
		skipDuplicates: true,
		loadDataLocal:  false,
		batchSize:      1000,
	}}
	parser2 := SwiftParser{params: SwiftParserParams{ // raises an error on duplicate entry
		db:             db,
		queries:        q,
		skipDuplicates: false,
		loadDataLocal:  false,
		batchSize:      1000,
	}}

	csvTestStrValid := "COUNTRY ISO2 CODE,SWIFT CODE,CODE TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIME ZONE\nAL,AAISALTRXXX,BIC11,UNITED BANK OF ALBANIA SH.A,\"HYRJA 3 RR. DRITAN HOXHA ND. 11 TIRANA, TIRANA, 1023\",TIRANA,ALBANIA,Europe/Tirane"
	csvTestStrDuplicates := "COUNTRY ISO2 CODE,SWIFT CODE,CODE TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIME ZONE\nAL,AAISALTRXXX,BIC11,UNITED BANK OF ALBANIA SH.A,\"HYRJA 3 RR. DRITAN HOXHA ND. 11 TIRANA, TIRANA, 1023\",TIRANA,ALBANIA,Europe/Tirane\nAL,AAISALTRXXX,BIC11,UNITED BANK OF ALBANIA SH.A,\"HYRJA 3 RR. DRITAN HOXHA ND. 11 TIRANA, TIRANA, 1023\",TIRANA,ALBANIA,Europe/Tirane"

	testCases := []parserTestCase{
		{parserName: "parser1", data: csvTestStrValid, expected_err: nil, expected_db_count: 1},
		{parserName: "parser1", data: csvTestStrDuplicates, expected_err: nil, expected_db_count: 1},
		{parserName: "parser2", data: csvTestStrDuplicates, expected_err: &mysql.MySQLError{Number: 1062}, expected_db_count: 0},
	}

	var parserTested *SwiftParser
	for _, testCase := range testCases {
		f, err := writeTempFile(testCase.data, "csv")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(f)
		if testCase.parserName == "parser1" {
			parserTested = &parser1
		} else if testCase.parserName == "parser2" {
			parserTested = &parser2
		} else {
			t.Fatalf("Invalid testcase %v", testCase)
		}

		err = resetTestDB(db)
		if err != nil {
			t.Fatalf("Failed to reset database: %v", err)
		}
		err = parserTested.Parse(f)
		if (err != nil && testCase.expected_err == nil) || !errors.Is(err, testCase.expected_err) {
			t.Fatalf("Parse() failed: for input %v expected output %v, got %v", testCase.data, testCase.expected_err, err)
		}

		var count uint
		err = db.QueryRow("SELECT COUNT(*) FROM bank").Scan(&count)
		if err != nil {
			t.Fatalf("Failed to query database: %v", err)
		}
		if count != testCase.expected_db_count {
			t.Errorf("Expected %d records, got %d", testCase.expected_db_count, count)
		}
	}
}
