package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tealeg/xlsx/v3"

	"swift-restful/repository"
	sqlc "swift-restful/repository/sqlc"
	"swift-restful/utils"

	"github.com/go-sql-driver/mysql"
)

type processRecordParams struct {
	db             *sql.DB
	queries        *sqlc.Queries
	skipDuplicates bool
}

func parseInsertCsvData(filePath string, p processRecordParams) error {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
		return err
	}
	defer f.Close()
	reader := csv.NewReader(bufio.NewReader(f))
	_, err = reader.Read()
	if err != nil {
		log.Fatal("Failed to read header", err)
	}
	var i int
	for {
		log.Printf("Processing record %d", i)
		i = i + 1
		record, err := reader.Read()
		if err != nil {
			break // EOF
		}
		err = processRecord(record, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseInsertXlsxData(filePath string, p processRecordParams) error {
	f, err := xlsx.OpenFile(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
		return err
	}
	sheetLen := len(f.Sheets)
	if sheetLen == 0 {
		return errors.New("Empty spreadsheet " + filePath)
	}
	sheet := f.Sheets[0]
	defer sheet.Close()
	var record []string
	var i int

	err = sheet.ForEachRow(func(row *xlsx.Row) error {
		if row != nil {
			log.Printf("Processing record %d", i)
			i = i + 1
			record = record[:0]
			err := row.ForEachCell(func(cell *xlsx.Cell) error {
				str, err := cell.FormattedValue()
				if err != nil {
					return err
				}
				record = append(record, str)
				return nil
			})
			if err != nil {
				return err
			}
			err = processRecord(record, p)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func parseInsertSpreadsheetData(filePath string, p processRecordParams) error {
	if len(filePath) == 0 {
		return errors.New("empty filename")
	}
	filePathSplit := strings.Split(filePath, ".")
	fileExt := filePathSplit[len(filePathSplit)-1]
	var err error

	switch fileExt {
	case "csv":
		err = parseInsertCsvData(filePath, p)
	case "xlsx":
		err = parseInsertXlsxData(filePath, p)
	default:
		err = errors.New("Unknown spreadsheet extension " + fileExt)
	}
	return err

}

func processRecord(record []string, p processRecordParams) error {

	if len(record) != 8 {
		return fmt.Errorf("invalid record length %d", len(record))
	}

	CountryISO2 := strings.ToUpper(strings.TrimSpace(record[0]))
	SwiftCode := strings.TrimSpace(record[1])
	BankName := strings.TrimSpace(record[3])
	Address := strings.TrimSpace(record[4])
	CountryName := strings.ToUpper(strings.TrimSpace(record[6]))

	isHeadquarter_, err := utils.IsHeadquarter(SwiftCode)
	if err != nil {
		return err
	}
	params := repository.CreateBankUpdateRelationsParams{
		BankParams: sqlc.CreateBankParams{
			Address:       Address,
			BankName:      BankName,
			CountryISO2:   CountryISO2,
			IsHeadquarter: isHeadquarter_,
			SwiftCode:     SwiftCode},
		CountryParams: sqlc.CreateCountryParams{
			Iso2: CountryISO2,
			Name: CountryName},
	}
	ctx := context.Background()
	err = repository.CreateBankUpdateRelations(ctx, p.db, p.queries, params)
	if driverErr, ok := err.(*mysql.MySQLError); ok {
		if driverErr.Number == repository.MYSQL_ER_DUP_ENTRY && p.skipDuplicates {
			// log.Printf("skipping dublicate entry %s", params.BankParams)
			return nil
		}
	}

	return err
}

func main() {
	var filePath string
	var batchSize int
	var skipDuplicates bool
	var verbose bool
	flag.StringVar(&filePath, "f", "", "Spreadsheet file path (required)")
	flag.IntVar(&batchSize, "batch-size", 1000, "Batch size for db insertion")
	flag.BoolVar(&skipDuplicates, "skip-duplicates", false, "Skip duplicate bank entries")
	flag.BoolVar(&verbose, "verbose", true, "Output additional info")

	flag.Parse()

	if filePath == "" {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	db, queries, err := repository.SetupDB()
	if err != nil {
		log.Fatal("Cannot setup DB: ", err)
		os.Exit(1)
	}
	err = parseInsertSpreadsheetData(filePath, processRecordParams{db, queries, skipDuplicates})
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
