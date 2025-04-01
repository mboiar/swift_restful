# Swift-restful - RESTful API for SWIFT data handling

## 📌 Overview
Swift-restful is a RESTful API for retrieving and modifying SWIFT bank data records stored in a MySQL database, written in Go with Gin framework.

## 🚀 Features
- CRUD operations for accessing and modifying SWIFT data
- Fast look-up by SWIFT code and country ISO2 code
- Containerized deployment with Docker
- Utility for fast upload of spreadsheet-stored data

## 🏗️ Setup & Installation

### 1️⃣ Prerequisites
Ensure you have the following installed:
- [Go](https://go.dev/dl/)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

### 2️⃣ Clone the Repository
```sh
git clone https://github.com/mboiar/swift-restful.git
cd swift-restful
```

### 4️⃣ Run the Application with Docker-compose
```sh
docker-compose up --build
```
This will build and start the API along with a MySQL database container. Requests will be accepted at `localhost:8080`.

To populate the database with example data from a spreadsheet, run the CLI parser (See parser [documentation](https://github.com/mboiar/swift_restful/blob/main/cmd/swift-parser/README.md)):
```sh
    docker exec swift_restful-web-1 go run cmd/swift-parser/main.go -f '/go/src/github.com/mboiar/swift-restful/example/Interns_2025_SWIFT_CODES - Sheet1.csv' --skip-duplicates
```

## 📡 API Endpoints

| Method | Endpoint                                   | Description                          |
|--------|--------------------------------------------|--------------------------------------|
| GET    | `/v1/swift-codes/:swift-code`              | Get all banks by SWIFT code          |
| POST   | `/v1/swift-codes/`                         | Create a SWIFT data entry            |
| GET    | `/v1/swift-codes/country/:countryISO2code` | Get SWIFT data by country ISO2 code  |
| DELETE | `/v1/swift-codes/:swift-code`              | Delete bank by SWIFT code            |

## 🧪 Running Tests

To run unit and integration tests:
```sh
go test ./... -v
```
