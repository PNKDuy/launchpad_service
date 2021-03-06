package main

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"launchpad_service/controller"
	_ "launchpad_service/docs"
	_ "net/http/pprof"
	"time"
)

// @title Swagger Launchpad Service API
// @version 1.0
// @description This is Launchpad service server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
func main() {
	go controller.DoEvery(2 *time.Second, controller.GetPriceAndUpdateList)

		server := echo.New()

		launchpad := server.Group("/launchpad")
		{
			launchpad.POST("/create", controller.Create)
			launchpad.GET("/get", controller.Get)
			launchpad.GET("/get-by-id/:id", controller.GetById)
			launchpad.PUT("/update/:id", controller.Update)
			launchpad.PUT("/deactivate-token/:id", controller.DeactivateToken)
		}

		token := server.Group("/token")
		{
			token.GET("/price/:token", controller.GetPrice)
			token.GET("/price", controller.GetAllPrice)
			token.GET("/klines/:token/:interval", controller.GetKlines)
			token.GET("/transaction/:hash", controller.GetTransaction)
			token.GET("/price-by-currency/:token/:currency", controller.GetPriceByCurrency)
		}

		server.GET("/swagger/*", echoSwagger.WrapHandler)
		server.Logger.Fatal(server.Start(":8081"))
}
