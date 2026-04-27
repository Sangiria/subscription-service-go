package main

import (
	"subscription-service-go/internal/environment"
	"subscription-service-go/internal/repository"
	"subscription-service-go/internal/routes"

	"github.com/labstack/echo/v5"
)

func main(){
	environment.LoadEnvVariables()
	repository.GetDBConnection()
	
	e := echo.New()

	routes.InitSubscriptionRoutes(e)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}