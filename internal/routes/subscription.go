package routes

import (
	"subscription-service-go/internal/handlers"

	"github.com/labstack/echo/v5"
)

func InitSubscriptionRoutes(e *echo.Echo) {
	e.POST("/subscriptions", handlers.CreateSubscription)
}