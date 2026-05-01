package routes

import (
	"subscription-service-go/internal/handlers"

	"github.com/labstack/echo/v4"
)

func InitSubscriptionRoutes(e *echo.Echo, h *handlers.SubscriptionHandler) {
	g := e.Group("/subscriptions")

	g.POST("", h.CreateSubscription)
	g.GET("", h.ListSubscriptions)
	g.GET("/:id", h.GetSubscription)
	g.DELETE("/:id", h.DeleteSubscriptions)
	g.PATCH("/:id", h.UpdateSubscriptions)
	g.GET("/sum", h.SumSubscriptionsPrice)
}