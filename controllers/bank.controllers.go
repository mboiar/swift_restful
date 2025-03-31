package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	repository "swift-restful/repository/sqlc"
	"swift-restful/schemas"
	"swift-restful/utils"

	"github.com/gin-gonic/gin"
)

type SwiftController struct {
	q   *repository.Queries
	ctx context.Context
}

func NewSwiftController(q *repository.Queries, ctx context.Context) *SwiftController {
	return &SwiftController{q, ctx}
}

func (sc *SwiftController) CreateBank(ctx *gin.Context) {
	var payload *schemas.CreateSwiftEntry

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload", "error": err.Error()})
		return
	}
	countryArgs := &repository.CreateCountryParams{
		Iso2: payload.CountryIso2,
		Name: payload.CountryName,
	}
	bankArgs := &repository.CreateBankParams{
		Address:     payload.Address,
		Name:        payload.BankName,
		CountryIso2: payload.CountryIso2,
		SwiftCode:   payload.SwiftCode,
	}
	_, err := sc.q.CreateCountry(ctx, *countryArgs)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "failed retrieving SWIFT data", "error": err.Error()})
		return
	}
	_, err = sc.q.CreateBank(ctx, *bankArgs)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "failed retrieving SWIFT data", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully created SWIFT entry"})
}

func (sc *SwiftController) GetSwiftData(ctx *gin.Context) {
	swiftCode := ctx.Param("swift-code")
	isBankHeadquarter, err := utils.IsHeadquarter(swiftCode)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid SWIFT code format"})
		return
	}
	if isBankHeadquarter {
		branches, err := sc.q.GetBranchesBySwiftCode(ctx, swiftCode)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"message": "failed to retrieve SWIFT data for SWIFT code"})
				return
			}
			ctx.JSON(http.StatusBadGateway, gin.H{"message": "failed retireving SWIFT data", "error": err.Error()})
			return
		}
		var headquarter *schemas.GetHeadquarterEntry
		var branchEntries []schemas.GetBranchEntry
		for _, bank := range branches {
			if bank.IsHeadquarter {
				headquarter = &schemas.GetHeadquarterEntry{
					Address:       bank.Address.String,
					BankName:      bank.Name,
					CountryIso2:   bank.CountryIso2,
					CountryName:   bank.Name_2,
					IsHeadquarter: bank.IsHeadquarter,
					SwiftCode:     bank.SwiftCode,
				}
			} else {
				branchEntries = append(branchEntries, schemas.GetBranchEntry{
					Address:       bank.Address.String,
					BankName:      bank.Name,
					CountryIso2:   bank.CountryIso2,
					IsHeadquarter: bank.IsHeadquarter,
					SwiftCode:     bank.SwiftCode,
				})
			}
		}
		fmt.Printf("%v", branches)
		if headquarter == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "failed to retrieve SWIFT data for SWIFT code"})
			return
		}
		headquarter.Branches = branchEntries
		ctx.JSON(http.StatusOK, headquarter)
		return
	} else {
		branch, err := sc.q.GetBankBySwiftCode(ctx, swiftCode)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, gin.H{"message": "failed to retrieve SWIFT data for SWIFT code"})
				return
			}
			ctx.JSON(http.StatusBadGateway, gin.H{"message": "failed retireving SWIFT data", "error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"address":       branch.Address.String,
			"bankName":      branch.Name,
			"countryISO2":   branch.CountryIso2,
			"countryName":   branch.Name_2,
			"ssHeadquarter": branch.IsHeadquarter,
			"swiftCode":     branch.SwiftCode,
		})
		return
	}
}

func (sc *SwiftController) GetCountryData(ctx *gin.Context) {
	iso2Code := ctx.Param("countryISO2code")
	if !utils.IsValidISO2Code(iso2Code) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid country ISO2 code format"})
		return
	}

	branches, err := sc.q.GetBranchesByCountryISO2(ctx, iso2Code)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "failed to retrieve SWIFT data for ISO2 code"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "failed retireving SWIFT data", "error": err.Error()})
		return
	}
	var branchEntries []schemas.GetBranchEntry
	if len(branches) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("no entries for code %s", iso2Code)})
		return
	}
	countryName := branches[0].CountryName
	for _, bank := range branches {
		branchEntries = append(branchEntries, schemas.GetBranchEntry{
			Address:       bank.Address.String,
			BankName:      bank.Name,
			CountryIso2:   bank.CountryIso2,
			IsHeadquarter: bank.IsHeadquarter,
			SwiftCode:     bank.SwiftCode,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"countryISO2": iso2Code,
		"countryName": countryName,
		"swiftCodes":  branchEntries,
	})
}

func (sc *SwiftController) DeleteBankBySwiftCode(ctx *gin.Context) {
	swiftCode := ctx.Param("swift-code")
	isSwiftCode := utils.IsValidSwiftCode(swiftCode)
	if !isSwiftCode {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid SWIFT code format"})
		return
	}
	_, err := sc.q.GetBankBySwiftCode(ctx, swiftCode)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "failed to retrieve SWIFT data for given SWIFT code"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "failed retireving SWIFT data", "error": err.Error()})
		return
	}
	err = sc.q.DeleteBank(ctx, swiftCode)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "failed deleting SWIFT entry", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully deleted SWIFT entry",
	})
}
