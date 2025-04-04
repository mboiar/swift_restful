package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	routes "swift-restful/api/v1"
	"swift-restful/controllers"
	repository "swift-restful/repository/sqlc"
	"swift-restful/schemas"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	tmysql "github.com/testcontainers/testcontainers-go/modules/mysql"
)

type APITestCase struct {
	requestType   string
	route         string
	payload       string
	expected_data string
	expected_code int
}

type Response map[string]interface{}

// createTestDB returns an instance of database for testing purposes
func createTestDB() (*sql.DB, *repository.Queries, func(), error) {
	ctx := context.Background()
	// create MySQL instance with docker
	container, err := tmysql.Run(ctx,
		"mysql:latest",
		tmysql.WithDatabase("test_db"),
		tmysql.WithUsername("root"),
		tmysql.WithPassword("password"),
		tmysql.WithScripts(
			filepath.Join("repository", "migration", "000001_init_schema.up.sql"),
			filepath.Join("testdata", "values.sql"),
		),
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

func TestGetSwiftDataBySwiftCode(t *testing.T) {
	ctx = context.TODO()
	_, q, cleanup, err := createTestDB()
	if err != nil {
		t.Fatalf("Failed to create test DB: %v", err)
	}
	defer cleanup()

	SwiftController = *controllers.NewSwiftController(q, ctx)
	SwiftRoutes = routes.NewRouteSwift(SwiftController)

	server = gin.Default()
	server.SetTrustedProxies(nil)
	router := server.Group("/")
	SwiftRoutes.SwiftRoute(router)
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("The specified route %s not found", ctx.Request.URL)})
	})

	testCases := []APITestCase{
		{requestType: "GET", route: "/v1/swift-codes/123", payload: "", expected_data: `{"message":"invalid SWIFT code format"}`, expected_code: 400},
		{requestType: "GET", route: "/v1/swift-codes/CRBAALTRXXX", payload: "", expected_data: `{"address":"TIRANA, TIRANA","bankName":"BANKA OTP ALBANIA SH.A","countryISO2":"AL","countryName":"ALBANIA","isHeadquarter":true,"swiftCode":"CRBAALTRXXX","branches":null}`, expected_code: 200},
		{requestType: "GET", route: "/v1/swift-codes/CrBaALtRxxx", payload: "", expected_data: `{"address":"TIRANA, TIRANA","bankName":"BANKA OTP ALBANIA SH.A","countryISO2":"AL","countryName":"ALBANIA","isHeadquarter":true,"swiftCode":"CRBAALTRXXX","branches":null}`, expected_code: 200},
		{requestType: "GET", route: "/v1/swift-codes/STANALT1SHY", payload: "", expected_data: `{"message":"failed to retrieve SWIFT data for SWIFT code"}`, expected_code: 404},
		{requestType: "GET", route: "/v1/swift-codes/STANALT1SHX", payload: "", expected_data: `{"address":"SHKODER, SHKODER, 4001","bankName":"BANK OF ALBANIA","countryISO2":"AL","countryName":"ALBANIA","isHeadquarter":false,"swiftCode":"STANALT1SHX"}`, expected_code: 200},
		{requestType: "GET", route: "/v1/swift-codes/BPKOPLPWXXX", payload: "", expected_data: `{"address":"UL. PULAWSKA 15  WARSZAWA, MAZOWIECKIE, 02-515","bankName":"PKO BANK POLSKI S.A.","countryISO2":"PL","countryName":"POLAND","isHeadquarter":true,"swiftCode":"BPKOPLPWXXX","branches":[{"address":"TYSIACLECIA PANSTWA POLSKIEGO 6  BIALYSTOK, PODLASKIE, 15-111","bankName":"PKO BANK POLSKI S.A.","countryISO2":"PL","isHeadquarter":false,"swiftCode":"BPKOPLPWBIA"}]}`, expected_code: 200},
	}
	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tc.requestType, tc.route, nil)
		server.ServeHTTP(w, req)

		assert.Equal(t, tc.expected_code, w.Code)
		assert.Equal(t, tc.expected_data, w.Body.String())
	}
}

func TestGetSwiftDataByCountryCode(t *testing.T) {
	ctx = context.TODO()
	_, q, cleanup, err := createTestDB()
	if err != nil {
		t.Fatalf("Failed to create test DB: %v", err)
	}
	defer cleanup()

	SwiftController = *controllers.NewSwiftController(q, ctx)
	SwiftRoutes = routes.NewRouteSwift(SwiftController)

	server = gin.Default()
	server.SetTrustedProxies(nil)
	router := server.Group("/")
	SwiftRoutes.SwiftRoute(router)
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("The specified route %s not found", ctx.Request.URL)})
	})

	testCases := []APITestCase{
		{requestType: "GET", route: "/v1/swift-codes/country/pll", payload: "", expected_data: `{"message":"invalid country ISO2 code format"}`, expected_code: 400},
		{requestType: "GET", route: "/v1/swift-codes/country/PL", payload: "", expected_data: `{"countryISO2":"PL","countryName":"POLAND","swiftCodes":[{"address":"UL. PULAWSKA 15  WARSZAWA, MAZOWIECKIE, 02-515","bankName":"PKO BANK POLSKI S.A.","countryISO2":"PL","isHeadquarter":true,"swiftCode":"BPKOPLPWXXX"},{"address":"TYSIACLECIA PANSTWA POLSKIEGO 6  BIALYSTOK, PODLASKIE, 15-111","bankName":"PKO BANK POLSKI S.A.","countryISO2":"PL","isHeadquarter":false,"swiftCode":"BPKOPLPWBIA"}]}`, expected_code: 200},
		{requestType: "GET", route: "/v1/swift-codes/country/pl", payload: "", expected_data: `{"countryISO2":"PL","countryName":"POLAND","swiftCodes":[{"address":"UL. PULAWSKA 15  WARSZAWA, MAZOWIECKIE, 02-515","bankName":"PKO BANK POLSKI S.A.","countryISO2":"PL","isHeadquarter":true,"swiftCode":"BPKOPLPWXXX"},{"address":"TYSIACLECIA PANSTWA POLSKIEGO 6  BIALYSTOK, PODLASKIE, 15-111","bankName":"PKO BANK POLSKI S.A.","countryISO2":"PL","isHeadquarter":false,"swiftCode":"BPKOPLPWBIA"}]}`, expected_code: 200},
		{requestType: "GET", route: "/v1/swift-codes/country/BP", payload: "", expected_data: `{"message":"no entries for code BP"}`, expected_code: 404},
	}
	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tc.requestType, tc.route, nil)
		server.ServeHTTP(w, req)

		assert.Equal(t, tc.expected_code, w.Code)
		assert.Equal(t, tc.expected_data, w.Body.String())
	}
}

func TestDeleteSwiftData(t *testing.T) {
	ctx = context.TODO()
	_, q, cleanup, err := createTestDB()
	if err != nil {
		t.Fatalf("Failed to create test DB: %v", err)
	}
	defer cleanup()

	SwiftController = *controllers.NewSwiftController(q, ctx)
	SwiftRoutes = routes.NewRouteSwift(SwiftController)

	server = gin.Default()
	server.SetTrustedProxies(nil)
	router := server.Group("/")
	SwiftRoutes.SwiftRoute(router)
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("The specified route %s not found", ctx.Request.URL)})
	})

	testCases := []APITestCase{
		{requestType: "DELETE", route: "/v1/swift-codes/123", payload: "", expected_data: `{"message":"invalid SWIFT code format"}`, expected_code: 400},
		{requestType: "DELETE", route: "/v1/swift-codes/CRBAALTRXXX", payload: "", expected_data: `{"message":"successfully deleted SWIFT entry"}`, expected_code: 200},
		{requestType: "DELETE", route: "/v1/swift-codes/STANALT1SHX", payload: "", expected_data: `{"message":"successfully deleted SWIFT entry"}`, expected_code: 200},
		{requestType: "DELETE", route: "/v1/swift-codes/CRBAALTRXXY", payload: "", expected_data: `{"message":"failed to retrieve SWIFT data for given SWIFT code"}`, expected_code: 404},
	}
	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tc.requestType, tc.route, nil)
		server.ServeHTTP(w, req)

		assert.Equal(t, tc.expected_code, w.Code)
		assert.Equal(t, tc.expected_data, w.Body.String())
	}
}

func TestAddSwiftData(t *testing.T) {
	ctx = context.TODO()
	_, q, cleanup, err := createTestDB()
	if err != nil {
		t.Fatalf("Failed to create test DB: %v", err)
	}
	defer cleanup()

	SwiftController = *controllers.NewSwiftController(q, ctx)
	SwiftRoutes = routes.NewRouteSwift(SwiftController)

	server = gin.Default()
	server.SetTrustedProxies(nil)
	router := server.Group("/")
	SwiftRoutes.SwiftRoute(router)
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("The specified route %s not found", ctx.Request.URL)})
	})

	tval := true
	fval := false
	testBankValid, _ := json.Marshal(schemas.CreateSwiftEntry{
		Address:       "Address",
		BankName:      "Bank",
		CountryIso2:   "UA",
		CountryName:   "Ukraine",
		SwiftCode:     "ODEKREUAXXX",
		IsHeadquarter: &tval,
	})
	testBankInconsistent, _ := json.Marshal(schemas.CreateSwiftEntry{
		Address:       "Address",
		BankName:      "Bank",
		CountryIso2:   "UA",
		CountryName:   "Ukraine",
		SwiftCode:     "ODEKREUAXXX",
		IsHeadquarter: &fval,
	})
	testBankInvalidIso2, _ := json.Marshal(schemas.CreateSwiftEntry{
		Address:       "Address",
		BankName:      "Bank",
		CountryIso2:   "UAA",
		CountryName:   "Ukraine",
		SwiftCode:     "ODEKREUAXXX",
		IsHeadquarter: &tval,
	})
	testBankInvalidSwiftCode, _ := json.Marshal(schemas.CreateSwiftEntry{
		Address:       "Address",
		BankName:      "Bank",
		CountryIso2:   "UA",
		CountryName:   "Ukraine",
		SwiftCode:     "ODEKREUAXXXt",
		IsHeadquarter: &tval,
	})
	testBankEmptyAddress, _ := json.Marshal(schemas.CreateSwiftEntry{
		Address:       "Address",
		BankName:      "Bank",
		CountryIso2:   "UA",
		CountryName:   "Ukraine",
		SwiftCode:     "ODEKREUAXXO",
		IsHeadquarter: &fval,
	})
	testCases := []APITestCase{
		{
			requestType:   "POST",
			route:         "/v1/swift-codes/",
			payload:       string(testBankValid),
			expected_data: `{"message":"successfully created SWIFT entry"}`,
			expected_code: 200,
		},
		{
			requestType:   "POST",
			route:         "/v1/swift-codes/",
			payload:       string(testBankInconsistent),
			expected_data: `{"message":"invalid payload: inconsistent headquarter info"}`,
			expected_code: 400,
		},
		{
			requestType:   "POST",
			route:         "/v1/swift-codes/",
			payload:       string(testBankInvalidIso2),
			expected_data: `{"message":"invalid payload: invalid code format"}`,
			expected_code: 400,
		},
		{
			requestType:   "POST",
			route:         "/v1/swift-codes/",
			payload:       string(testBankInvalidSwiftCode),
			expected_data: `{"message":"invalid payload: invalid code format"}`,
			expected_code: 400,
		},
		{
			requestType:   "POST",
			route:         "/v1/swift-codes/",
			payload:       string(testBankEmptyAddress),
			expected_data: `{"message":"successfully created SWIFT entry"}`,
			expected_code: 200,
		},
	}
	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(tc.requestType, tc.route, strings.NewReader(tc.payload))
		server.ServeHTTP(w, req)

		assert.Equal(t, tc.expected_code, w.Code)
		assert.Equal(t, tc.expected_data, w.Body.String())
	}
}
