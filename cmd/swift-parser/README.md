# SwiftParser - A command-line tool for uploading SWIFT data to a database

## 📌 Overview
SwiftParser is a command-line tool for transferring SWIFT data from a spreadsheet file to a MySQL database using bulk insertion for fast uploads. CSV and XLSX spreadsheet formats are supported.

## 🚀 Usage
Run the CLI tool with:

```sh
go run ./cmd/swift-parser/main.go [flags]
```

or if built:
- Linux
    ```sh
    ./swift-parser [flags]
    ```
- Windows
    ```cmd
    ./swift-parser.exe [flags]
    ```
## Installation

- Linux
    ```sh
    go build -o bin/swift-parser ./cmd/swift-parser
    ```
- Windows
    ```sh
    go build -o bin/swift-parser.exe ./cmd/swift-parser
    ```
## Examples

```sh
# load data with default options
./swift-parser -f "path/to/spreadsheet" --db-config "path/to/env/file"
# specify batch size for insertion, use verbose mode
./swift-parser -f "file.csv" --db-config "db.env" --batch-size 500 -v
# skip duplicate bank entries without throwing an error
./swift-parser -f "file.csv" --db-config "db.env" --skip-duplicates
# use LOAD DATA LOCAL instruction for fast inserts of large number of rows
./swift-parser -f "file.csv" --db-config "db.env" --load-data-local
```

Warning: LOAD DATA LOCAL requires additional MySQL configuration. See [official documentation](https://dev.mysql.com/doc/refman/8.4/en/load-data-local-security.html).