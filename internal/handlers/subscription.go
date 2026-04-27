package handlers

import (
	"time"

	"github.com/labstack/echo/v5"
)

type Subscription struct {
	ServiceName string 		`json:"service_name"`
	Price		int			`json:"price"`
	UserId		string		`json:"user_id"`
	StartDate	time.Time	`json:"start_date"`
	EndDate		time.Time	`json:"end_date,omitzero"`
}

func CreateSubscription(c *echo.Context) error {
	return nil
}