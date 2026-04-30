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

type apiError struct {
    Message string 	`json:"message"`
	Details string	`json:"details,omitzero"`
}

func sendError(c *echo.Context, code int, msg string, details string) error {
    return c.JSON(code, apiError{
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
	var subReq models.SubscriptionCreateReq

	if err := c.Bind(&subReq); err != nil {
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
	if _, err := uuid.Parse(subId); err != nil {
		return sendError(c, http.StatusBadRequest, "Invalid UUID format", err.Error())
    }

	sub, err := h.repo.Get(subId)
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return sendError(c, http.StatusNotFound, "This subscription doesn't exist", err.Error())
		}
		return sendError(c, http.StatusInternalServerError, "Error getting subscription record", err.Error())
	}

	return c.JSON(http.StatusOK, sub)
}

//TODO: refactor
//TODO: add user_id query_param

func (h *SubscriptionHandler) ListSubscriptions(c *echo.Context) error {
	var subReq models.ListParams
	if err := c.Bind(&subReq); err != nil {
        return sendError(c, http.StatusBadRequest, "Invalid parameters", err.Error())
    }

	if err := validation.Validate(&subReq, nil); err != nil {
		return sendError(c, http.StatusBadRequest, "Invalid parameters", err.Error())
	}

	subs, err := h.repo.List(subReq)
	if err != nil {
		return sendError(c, http.StatusInternalServerError, "Error getting subscription records", err.Error())
	}

	return c.JSON(http.StatusOK, subs)
}

func (h *SubscriptionHandler) DeleteSubscriptions(c *echo.Context) error {
	subId := c.Param("id")
	if _, err := uuid.Parse(subId); err != nil {
        return sendError(c, http.StatusBadRequest, "Invalid UUID format", err.Error())
    }

	if err := h.repo.Delete(subId); err != nil {
		if err == gorm.ErrRecordNotFound {
			return sendError(c, http.StatusNotFound, "This subscription doesn't exist", err.Error())
		}

		return sendError(c, http.StatusInternalServerError, "Error deleting subscription record", err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *SubscriptionHandler) UpdateSubscriptions(c *echo.Context) error {
	subId := c.Param("id")
	if _, err := uuid.Parse(subId); err != nil {
        return sendError(c, http.StatusBadRequest, "Invalid UUID format", err.Error())
    }

	var subReq models.SubscriptionUpdateReq
	if err := c.Bind(&subReq); err != nil {
		return sendError(c, http.StatusBadRequest, "Validation failed", err.Error())
    }

	fields := subReq.ToMap()
	if len(fields) == 0 {
		return sendError(c, http.StatusBadRequest, "Update failed", "No valid fields provided for update")
	}
	
	if err := validation.Validate(&subReq, utils.Ptr(models.TagUpdate)); err != nil {
		return sendError(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	sub, err := h.repo.Update(subId, fields)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return sendError(c, http.StatusNotFound, "This subscription doesn't exist", err.Error())
		}

		return sendError(c, http.StatusInternalServerError, "Error updating subscription record", err.Error())
	}

	return c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHandler) SumSubscriptionsPrice(c *echo.Context) error {
	var subReq models.SumSubscriptionPriceParams
	if err := c.Bind(&subReq); err != nil {
		return sendError(c, http.StatusBadRequest, "Invalid parameters", err.Error())
    }

	if err := validation.Validate(&subReq, nil); err != nil {
		return sendError(c, http.StatusBadRequest, "Invalid parameters", err.Error())
	}

	total, err := h.repo.Sum(subReq)
	if err != nil {
		return sendError(c, http.StatusInternalServerError, "Error calculating subscription sum price", err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"total": total,
	})
}