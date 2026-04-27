package handlers

import (
	"context"
	"net/http"
	"subscription-service-go/internal/models"
	"subscription-service-go/internal/repository"
	"subscription-service-go/internal/utils"
	"subscription-service-go/internal/validation"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
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
	var subReq models.SubscriptionReq

	if err := validation.BindAndValidate(c, &subReq, models.TagCreate); err != nil {
		return sendError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	sub := models.Subscription{
		ServiceName: subReq.ServiceName,
		Price: subReq.Price,
		UserId: subReq.UserId,
		StartDate: *utils.ParseToDate(subReq.StartDate),
		EndDate: utils.ParseToDate(subReq.EndDate),
	}

	if err := gorm.G[models.Subscription](repository.Db).Create(context.Background(), &sub); err != nil {
		return sendError(c, http.StatusInternalServerError, "Failed to create record", err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"subscription": sub,
	})
}