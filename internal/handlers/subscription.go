package handlers

import (
	"errors"
	"net/http"
	"subscription-service-go/internal/models"
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

type SubscriptionHandler struct {
    db *gorm.DB
}

func NewSubscriptionHandler(db *gorm.DB) *SubscriptionHandler {
    return &SubscriptionHandler{db: db}
}

func (h *SubscriptionHandler) CreateSubscription(c *echo.Context) error {
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

	result := h.db.Create(&sub)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return sendError(c, http.StatusConflict, "Record already exist", (result.Error).Error())
		}
		return sendError(c, http.StatusInternalServerError, "Couldn't create record", (result.Error).Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"subscription": sub,
	})
}