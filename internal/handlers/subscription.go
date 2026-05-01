package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"subscription-service-go/internal/models"
	"subscription-service-go/internal/repository"
	"subscription-service-go/internal/utils"
	"subscription-service-go/internal/validation"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type apiError struct {
	Message string `json:"message"`
	Details string `json:"details,omitzero"`
}

func sendError(c echo.Context, code int, msg string, details string) error {
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

func (h *SubscriptionHandler) CreateSubscription(c echo.Context) error {
	var subReq models.SubscriptionCreateReq

	if err := c.Bind(&subReq); err != nil {
		slog.Error("failed to bind subscription create data", "error", err)
		return sendError(c, http.StatusBadRequest, "Invalid parameters", err.Error())
	}

	if err := validation.Validate(&subReq, new(models.TagCreate)); err != nil {
		slog.Error("failed to validate subscription create parameters", "error", err)
		return sendError(c, http.StatusBadRequest, "Invalid parameters", err.Error())
	}

	sub := models.Subscription{
		ServiceName: subReq.ServiceName,
		Price:       subReq.Price,
		UserId:      subReq.UserId,
		StartDate:   *utils.ParseToDate(subReq.StartDate),
		EndDate:     utils.ParseToDate(subReq.EndDate),
	}

	err := h.repo.Create(&sub)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			slog.Warn("subscription already exists",
				"user_id", subReq.UserId,
				"service_name", subReq.ServiceName,
				"error", err,
			)
			return sendError(c, http.StatusConflict, "This subscription already exist", err.Error())
		}
		slog.Error("failed to create subscription record in database",
			"user_id", subReq.UserId,
			"service_name", subReq.ServiceName,
			"error", err,
		)
		return sendError(c, http.StatusInternalServerError, "Error creating subscription record", err.Error())
	}

	slog.Info("subscription record created",
		"sub_id", sub.Id,
		"user_id", sub.UserId,
		"service_name", sub.ServiceName,
		"start_date", sub.StartDate,
		"end_date", sub.EndDate,
	)

	return c.JSON(http.StatusOK, map[string]any{
		"subscription": sub,
	})
}

func (h *SubscriptionHandler) GetSubscription(c echo.Context) error {
	subId := c.Param("id")
	if _, err := uuid.Parse(subId); err != nil {
		slog.Error("failed to validate id parameter", "error", err)
		return sendError(c, http.StatusBadRequest, "Invalid parameter", err.Error())
	}

	sub, err := h.repo.Get(subId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("subscription not found", "id", subId, "error", err)
			return sendError(c, http.StatusNotFound, "This subscription doesn't exist", err.Error())
		}
		slog.Error("failed to get subscription record from database", "id", subId, "error", err)
		return sendError(c, http.StatusInternalServerError, "Error getting subscription record", err.Error())
	}

	slog.Info("fetched subscription record",
		"sub_id", sub.Id,
		"user_id", sub.UserId,
		"service_name", sub.ServiceName,
		"start_date", sub.StartDate,
		"end_date", sub.EndDate,
	)

	return c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHandler) ListSubscriptions(c echo.Context) error {
	var subReq models.ListParams
	if err := c.Bind(&subReq); err != nil {
		slog.Error("failed to bind list parameters", "error", err)
		return sendError(c, http.StatusBadRequest, "Invalid parameters", err.Error())
	}

	if err := validation.Validate(&subReq, nil); err != nil {
		slog.Error("failed to validate list parameters", "error", err)
		return sendError(c, http.StatusBadRequest, "Invalid parameters", err.Error())
	}

	subs, err := h.repo.List(subReq)
	if err != nil {
		slog.Error("failed to get subscription records from database", "error", err)
		return sendError(c, http.StatusInternalServerError, "Error getting subscription records", err.Error())
	}

	slog.Info("fetched subscription record",
		"records_count", len(subs),
		"limit", subReq.Limit,
		"offset", subReq.Offset,
		"user_id", subReq.UserID,
	)

	return c.JSON(http.StatusOK, subs)
}

func (h *SubscriptionHandler) DeleteSubscriptions(c echo.Context) error {
	subId := c.Param("id")
	if _, err := uuid.Parse(subId); err != nil {
		slog.Error("failed to validate id parameter", "error", err)
		return sendError(c, http.StatusBadRequest, "Invalid parameter", err.Error())
	}

	if err := h.repo.Delete(subId); err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Warn("subscription not found", "id", subId, "error", err)
			return sendError(c, http.StatusNotFound, "This subscription doesn't exist", err.Error())
		}
		slog.Error("failed to delete subscription record from database", "error", err)
		return sendError(c, http.StatusInternalServerError, "Error deleting subscription record", err.Error())
	}

	slog.Info("subscription record deleted", "sub_id", subId)

	return c.NoContent(http.StatusOK)
}

func (h *SubscriptionHandler) UpdateSubscriptions(c echo.Context) error {
	subId := c.Param("id")
	if _, err := uuid.Parse(subId); err != nil {
		slog.Error("failed to validate id parameter", "error", err)
		return sendError(c, http.StatusBadRequest, "Invalid parameters", err.Error())
	}

	var subReq models.SubscriptionUpdateReq
	if err := c.Bind(&subReq); err != nil {
		slog.Error("failed to bind subscription update data", "error", err)
		return sendError(c, http.StatusBadRequest, "Invalid parameters", err.Error())
	}

	if err := validation.Validate(&subReq, new(models.TagUpdate)); err != nil {
		slog.Error("failed to validate subscription update parameters", "error", err)
		return sendError(c, http.StatusBadRequest, "Invalid parameters", err.Error())
	}

	fields := subReq.ToMap()
	if len(fields) == 0 {
		slog.Warn("subscription update parameters is empty")
		return sendError(c, http.StatusBadRequest, "Nothing to update", "No valid fields provided for update")
	}

	sub, err := h.repo.Update(subId, fields)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			slog.Warn("subscription not found", "id", subId, "error", err)
			return sendError(c, http.StatusNotFound, "This subscription doesn't exist", err.Error())
		}
		slog.Error("failed to update subscription record from database", "error", err)
		return sendError(c, http.StatusInternalServerError, "Error updating subscription record", err.Error())
	}

	slog.Info("subscription record updated",
		"sub_id", sub.Id,
		"user_id", sub.UserId,
		"service_name", sub.ServiceName,
		"start_date", sub.StartDate,
		"end_date", sub.EndDate,
	)

	return c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHandler) SumSubscriptionsPrice(c echo.Context) error {
	var subReq models.SumSubscriptionPriceParams
	if err := c.Bind(&subReq); err != nil {
		slog.Error("failed to bind sum data", "error", err)
		return sendError(c, http.StatusBadRequest, "Invalid parameters", err.Error())
	}

	if err := validation.Validate(&subReq, nil); err != nil {
		slog.Error("failed to validate sum parameters", "error", err)
		return sendError(c, http.StatusBadRequest, "Invalid parameters", err.Error())
	}

	total, err := h.repo.Sum(subReq)
	if err != nil {
		slog.Error("failed to calculate sum of subscriptions price", "error", err)
		return sendError(c, http.StatusInternalServerError, "Error calculating subscription sum price", err.Error())
	}

	slog.Info("subscriptions sum price calculated",
		"total_sum", total,
		"user_id", subReq.UserID,
		"service_name", subReq.ServiceName,
		"start_date", subReq.StartDate,
		"end_date", subReq.EndDate,
	)

	return c.JSON(http.StatusOK, map[string]any{
		"total": total,
	})
}
