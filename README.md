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

### 2️⃣ Clone the Repository
```sh
git clone https://github.com/mboiar/swift-restful.git
cd swift-restful
```

### 4️⃣ Run the Application with Docker-compose
```sh
docker-compose up --build
```
This will build and start the API along with a MySQL database container.

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
