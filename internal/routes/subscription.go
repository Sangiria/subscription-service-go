package routes

import (
	"subscription-service-go/internal/handlers"

	"github.com/labstack/echo/v5"
)

func InitSubscriptionRoutes(e *echo.Echo, h *handlers.SubscriptionHandler) {
	e.POST("/subscriptions", h.CreateSubscription)
	e.GET("/subscriptions/:id", h.GetSubscription)
}