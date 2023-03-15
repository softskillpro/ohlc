package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"ohcl/config"
	"ohcl/controllers/data"
	_ "ohcl/docs/swagger"
	"ohcl/middleware"
)

// @title          	Historical OHCL Price Data docs
// @version         1.0
// @termsOfService  http://swagger.io/terms/
// @description OHLC is large amount of historical OHLC price data in CSV files format, which now needs to be centralized and digitized.
func main() {

	router := gin.Default()

	if err := Initializer(); err != nil {
		log.Fatalf("Gotcha Error on initializer %s", err.Error())
	}

	data.RouterRegister(router.RouterGroup)

	// initialize the swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.NoRoute(middleware.NotFoundPath)
	router.NoMethod(middleware.NotFoundMethod)

	port, err := config.Get("PORT")
	if err != nil {
		log.Fatalf("Gotcha Error on get config %s form .env file", err.Error())
	}

	if err := router.Run(port); err != nil {
		log.Fatalf("Gotcha Error on running web service %s", err.Error())
	}
}
