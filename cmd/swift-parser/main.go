/*
Swift-parser transfers SWIFT data from a spreadsheet file to a MySQL database.
It uses bulk insertion for fast uploads.

Required flag -f specifies spreadsheet file. CSV and XLSX spreadsheet formats are supported.

Usage

	swift-parser [flags] [-f ...]

Optional flags are:

	-db-config
		Load database configuration (.env file).
		See [repository.SetupDB] for details.
	-batch-size
		Specify batch size for bulk inserts.
		Default is 1000
	-skip-duplicates
		Do not raise an error if duplicate SWIFT entries encountered
	-v
		Verbose mode
	-load-data-local
		Use LOAD DATA LOCAL MySQL instruction for faster bulk inserts.
		Requires additional configuration. See "https://dev.mysql.com/doc/refman/8.4/en/load-data-local-security.html"
*/
package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/tealeg/xlsx/v3"

	"swift-restful/repository"
	sqlc "swift-restful/repository/sqlc"
)

// SwiftParserParams represents SwiftParser parameters
type SwiftParserParams struct {
	db             *sql.DB
	queries        *sqlc.Queries
	skipDuplicates bool
	loadDataLocal  bool
	batchSize      uint
}

// A Record represents values of a single SWIFT spreadsheet row
type Record struct {
	CountryIso2 string
	SwiftCode   string
	CodeType    string
	Name        string
	Address     string
	TownName    string
	CountryName string
	TimeZone    string
}

// FromArr constructs a Record from an array of strings representing its fields
func (r *Record) FromArr(recordArr []string) error {
	if len(recordArr) != 8 {
		return fmt.Errorf("invalid record length %d", len(recordArr))
	}
	r.CountryIso2 = recordArr[0]
	r.SwiftCode = recordArr[1]
	r.CodeType = recordArr[2]
	r.Name = recordArr[3]
	r.Name = recordArr[4]
	r.Address = recordArr[5]
	r.TownName = recordArr[6]
	r.CountryName = recordArr[7]
	r.TimeZone = recordArr[8]
	return nil
}

type processBulkArgs struct {
	bankBulkArgs    []sqlc.CreateBankBulkParams
	countryBulkArgs []sqlc.CreateCountryBulkParams
}

// A SwiftParser transfers SWIFT data from a spreadsheet into a MySQL database.
type SwiftParser struct {
	params SwiftParserParams
}

func (sp SwiftParser) parseCsvData(filePath string) error {
	var err error
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	reader := csv.NewReader(bufio.NewReader(f))
	_, err = reader.Read()
	if err != nil {
		slog.Error("Failed to read csv header")
		return err
	}
	var i uint
	var argsArr processBulkArgs
	var bankArgs sqlc.CreateBankBulkParams
	var countryArgs sqlc.CreateCountryBulkParams
	var record *Record
	for {
		slog.Debug(fmt.Sprintf("Processing record %d", i))
		i = i + 1
		recordArr, err := reader.Read()
		if err == io.EOF {
			break // EOF
		}
		err = record.FromArr(recordArr)
		if err != nil {
			return err
		}
		bankArgs, countryArgs, err = sp.processRecord(record)
		if err != nil {
			return err
		}
		argsArr.bankBulkArgs = append(argsArr.bankBulkArgs, bankArgs)
		argsArr.countryBulkArgs = append(argsArr.countryBulkArgs, countryArgs)
		if i%sp.params.batchSize == 0 {
			err = sp.processBatch(argsArr)
			if err != nil {
				return err
			}
			argsArr = processBulkArgs{}
		}
	}
	err = sp.processBatch(argsArr)
	return err
}

func (sp SwiftParser) parseXlsxData(filePath string) error {
	f, err := xlsx.OpenFile(filePath)
	if err != nil {
		return err
	}
	sheetLen := len(f.Sheets)
	if sheetLen == 0 {
		return errors.New("Empty spreadsheet " + filePath)
	}
	sheet := f.Sheets[0]
	defer sheet.Close()
	var recordArr []string
	var i uint
	var argsArr processBulkArgs
	var bankArgs sqlc.CreateBankBulkParams
	var countryArgs sqlc.CreateCountryBulkParams
	var record *Record

	err = sheet.ForEachRow(func(row *xlsx.Row) error {
		if row != nil {
			if i == 0 {
				i = i + 1
				return nil
			}
			slog.Debug(fmt.Sprintf("Processing record %d", i))
			i = i + 1
			recordArr = recordArr[:0]
			err := row.ForEachCell(func(cell *xlsx.Cell) error {
				str, err := cell.FormattedValue()
				if err != nil {
					return err
				}
				if str != "" {
					recordArr = append(recordArr, str)
				}
				return nil
			})
			if err != nil {
				return err
			}
			err = record.FromArr(recordArr)
			if err != nil {
				return err
			}
			bankArgs, countryArgs, err = sp.processRecord(record)
			if err != nil {
				return err
			}
			argsArr.bankBulkArgs = append(argsArr.bankBulkArgs, bankArgs)
			argsArr.countryBulkArgs = append(argsArr.countryBulkArgs, countryArgs)
			if i%sp.params.batchSize == 0 {
				err = sp.processBatch(argsArr)
				if err != nil {
					return err
				}
				argsArr = processBulkArgs{}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = sp.processBatch(argsArr)
	return err
}

// processBatch inserts bulk SWIFT data into a MySQL database.
func (sp SwiftParser) processBatch(argsArr processBulkArgs) error {
	slog.Info(fmt.Sprintf("Processing batch (n=%d)", len(argsArr.bankBulkArgs)))
	if sp.params.skipDuplicates {
		slog.Debug("Ignoring duplicate bank entries")
	}
	ctx := context.Background()
	tx, err := sp.params.db.Begin()
	slog.Info("Beginning transaction")
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if sp.params.loadDataLocal {
		slog.Debug("Using LOAD DATA LOCAL")
		qtx := sp.params.queries.WithTx(tx)
		_, err = qtx.CreateCountryBulk(ctx, argsArr.countryBulkArgs)
		if err != nil {
			return err
		}
		_, err = qtx.CreateBankBulk(ctx, argsArr.bankBulkArgs)
		if err != nil {
			return err
		}
	} else {
		slog.Debug("Using multi-row insert")
		err = repository.InsertCountryMultiRow(ctx, argsArr.countryBulkArgs, sp.params.db)
		if err != nil {
			slog.Info("Aborting transaction")
			return err
		}
		err = repository.InsertBankMultiRow(ctx, argsArr.bankBulkArgs, sp.params.db, sp.params.skipDuplicates)
		if err != nil {
			slog.Info("Aborting transaction")
			return err
		}
	}
	return tx.Commit()
}

// Parse loads SWIFT data into a database by calling a method matching given file extension.
// Spreadsheet file should have 8 columns representing Record fields and a header.
func (sp SwiftParser) Parse(filePath string) error {
	if len(filePath) == 0 {
		return errors.New("empty filename")
	}
	filePathSplit := strings.Split(filePath, ".")
	fileExt := filePathSplit[len(filePathSplit)-1]
	var err error

	switch fileExt {
	case "csv":
		err = sp.parseCsvData(filePath)
	case "xlsx":
		err = sp.parseXlsxData(filePath)
	default:
		err = fmt.Errorf("spreadsheet extension %s is not recognized (allowed extensions: %v)", fileExt, []string{"csv", "xlsx"})
	}
	return err

}

// processRecord converts a Record to args for insertion into a database.
func (sp SwiftParser) processRecord(record *Record) (sqlc.CreateBankBulkParams, sqlc.CreateCountryBulkParams, error) {
	CountryISO2 := strings.ToUpper(strings.TrimSpace(record.CountryIso2))
	SwiftCode := strings.TrimSpace(record.SwiftCode)
	BankName := strings.TrimSpace(record.Name)
	addressStr := strings.TrimSpace(record.Address)
	Address := sql.NullString{
		String: addressStr,
		Valid:  len(addressStr) > 0}
	CountryName := strings.ToUpper(strings.TrimSpace(record.CountryName))

	bankArgs := sqlc.CreateBankBulkParams{
		Address:     Address,
		Name:        BankName,
		CountryIso2: CountryISO2,
		SwiftCode:   SwiftCode}
	countryArgs := sqlc.CreateCountryBulkParams{
		Iso2: CountryISO2,
		Name: CountryName}
	return bankArgs, countryArgs, nil
}

func main() {
	var filePath string
	var batchSize uint
	var skipDuplicates bool
	var verbose bool
	var loadDataLocal bool
	var dbConfigPath string
	flag.StringVar(&filePath, "f", "", "Spreadsheet file path (required)")
	flag.UintVar(&batchSize, "batch-size", 1000, "Batch size for db insertion")
	flag.BoolVar(&skipDuplicates, "skip-duplicates", false, "Skip duplicate bank entries")
	flag.BoolVar(&loadDataLocal, "load-data-local", false, "Use LOAD DATA LOCAL statement for faster bulk insert (requires additional client and server configuration)")
	flag.BoolVar(&verbose, "v", false, "Output additional info")
	flag.StringVar(&dbConfigPath, "db-config", "", "Database config .env file")

	flag.Parse()

	if filePath == "" {
		fmt.Println("Usage: swift-parser [flags]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var logLevel slog.Level
	if verbose {
		logLevel = slog.LevelDebug
		fmt.Println("Using verbose mode")
	} else {
		logLevel = slog.LevelInfo
	}
	slog.SetLogLoggerLevel(logLevel)

	db, queries, err := repository.SetupDB(dbConfigPath)
	if err != nil {
		log.Fatal("Cannot setup DB: ", err)
	}
	p := SwiftParserParams{
		db:             db,
		queries:        queries,
		skipDuplicates: skipDuplicates,
		loadDataLocal:  loadDataLocal,
		batchSize:      batchSize,
	}
	sp := SwiftParser{params: p}
	err = sp.Parse(filePath)
	if err != nil {
		log.Fatal(err)
	}
}
