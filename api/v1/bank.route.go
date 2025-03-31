package routes

import (
	"swift-restful/controllers"

	"github.com/gin-gonic/gin"
)

type SwiftRoutes struct {
	swiftController controllers.SwiftController
}

func NewRouteSwift(swiftController controllers.SwiftController) SwiftRoutes {
	return SwiftRoutes{swiftController}
}

func (sr *SwiftRoutes) SwiftRoute(rg *gin.RouterGroup) {

	router := rg.Group("/v1/swift-codes")
	router.POST("/", sr.swiftController.CreateBank)
	router.GET("/country/:countryISO2code", sr.swiftController.GetCountryData)
	router.GET("/:swift-code", sr.swiftController.GetSwiftData)
	router.DELETE("/:swift-code", sr.swiftController.DeleteBankBySwiftCode)
}
