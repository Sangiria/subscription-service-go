package handlers

import (
	"net/http"
	"subscription-service-go/internal/models"
	"subscription-service-go/internal/validation"

	"github.com/labstack/echo/v5"
)

type ApiError struct {
    Status  int    	`json:"status"`
    Message string 	`json:"message"`
	Details string	`json:"details,omitzero"`
}

func sendError(c *echo.Context, code int, msg string, details string) error {
    return c.JSON(code, ApiError{
        Status:  code,
        Message: msg,
		Details: details,
    })
}

func CreateSubscription(c *echo.Context) error {
	var subscriptionReq models.Subscription

	if err := validation.BindAndValidate(c, &subscriptionReq, models.TagCreate); err != nil {
		return sendError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	return c.NoContent(http.StatusOK)
}