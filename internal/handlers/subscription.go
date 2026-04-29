package handlers

import (
	"errors"
	"net/http"
	"subscription-service-go/internal/models"
	"subscription-service-go/internal/repository"
	"subscription-service-go/internal/utils"
	"subscription-service-go/internal/validation"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type ApiError struct {
    Message string 	`json:"message"`
	Details string	`json:"details,omitzero"`
}

func sendError(c *echo.Context, code int, msg string, details string) error {
    return c.JSON(code, ApiError{
        Message: msg,
		Details: details,
    })
}

type SubscriptionHandler struct {
    repo repository.Repository
}

func NewSubscriptionHandler(repo repository.Repository) *SubscriptionHandler {
    return &SubscriptionHandler{repo: repo}
}

func (h *SubscriptionHandler) CreateSubscription(c *echo.Context) error {
	var subReq models.SubscriptionReq

	if err := c.Bind(subReq); err != nil {
		return sendError(c, http.StatusBadRequest, "Validation failed", err.Error())
    }

	if err := validation.Validate(&subReq, utils.Ptr(models.TagCreate)); err != nil {
		return sendError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	sub := models.Subscription{
		ServiceName: subReq.ServiceName,
		Price: subReq.Price,
		UserId: subReq.UserId,
		StartDate: *utils.ParseToDate(subReq.StartDate),
		EndDate: utils.ParseToDate(subReq.EndDate),
	}

	err := h.repo.Create(&sub)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return sendError(c, http.StatusConflict, "This subscription already exist", err.Error())
		}
		return sendError(c, http.StatusInternalServerError, "Couldn't create subscription record", err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"subscription": sub,
	})
}

func (h *SubscriptionHandler) GetSubscription(c *echo.Context) error {
	subId := c.Param("id")
	sub, err := h.repo.Get(subId)

	if _, err := uuid.Parse(subId); err != nil {
        return sendError(c, http.StatusBadRequest, "Invalid UUID format", err.Error())
    }
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(c, http.StatusNotFound, "This subscription doesn't exist", err.Error())
		}
		return sendError(c, http.StatusInternalServerError, "Error getting subscription record", err.Error())
	}

	return c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHandler) ListSubscriptions(c *echo.Context) error {
	limit, offset := utils.ToInt(c.QueryParam("limit"), -1), utils.ToInt(c.QueryParam("offset"), -1)

	if err := validation.Validate(&models.ListParams{Limit: limit, Offset: offset}, nil); err != nil {
		return sendError(c, http.StatusBadRequest, "Invalid parameters", err.Error())
	}

	subs, err := h.repo.List(limit, offset)
	if err != nil {
		return sendError(c, http.StatusInternalServerError, "Error getting subscription records", err.Error())
	}

	return c.JSON(http.StatusOK, subs)
}