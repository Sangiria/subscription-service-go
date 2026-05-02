package models

import (
	"subscription-service-go/internal/utils"
	"time"
)

const (
	TagCreate = "create"
	TagUpdate = "update"
)

type SubscriptionCreateReq struct {
	ServiceName string 		`json:"service_name" example:"YandexMusic" create:"required" update:"omitzero"`
	Price		int			`json:"price" example:"1200" create:"required,gt=-1" update:"omitzero,gt=-1"`
	UserId		string		`json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" create:"required,uuid"`
	StartDate	string		`json:"start_date" example:"10-2023" create:"required,datetime=01-2006" update:"omitzero,datetime=01-2006"`
	EndDate		string		`json:"end_date,omitzero" example:"10-2024" create:"omitzero,datetime=01-2006" update:"omitzero,datetime=01-2006"`
}

type SubscriptionUpdateReq struct {
	ServiceName *string 	`json:"service_name" update:"omitzero"`
	Price		*int		`json:"price" update:"omitzero,gt=-1"`
	StartDate	*string		`json:"start_date" update:"omitzero,datetime=01-2006"`
	EndDate		*string		`json:"end_date,omitzero" update:"omitzero,datetime=01-2006"`
}

type SumSubscriptionPriceParams struct {
	UserID  	string			`query:"user_id" validate:"required,uuid"`
	ServiceName string     		`query:"service_name" validate:"omitzero"`
	StartDate   string 			`query:"start_date" validate:"omitzero,datetime=01-2006"`
	EndDate     string 			`query:"end_date" validate:"omitzero,datetime=01-2006"`
}

type Subscription struct {
	Id			string			`json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ServiceName string 			`json:"service_name"`
	Price		int				`json:"price"`
	UserId		string			`json:"user_id"`
	StartDate	time.Time		`json:"start_date" swaggertype:"string" example:"2026-05-02T16:21:19Z"`
	EndDate		*time.Time		`json:"end_date,omitzero" swaggertype:"string" example:"2027-05-02T16:21:19Z"`
}

type ListParams struct {
	UserID *string 	`query:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" validate:"omitzero,uuid"`
	Limit  *int 	`query:"limit" example:"10" validate:"omitzero,gte=-1,lte=100"`
	Offset *int 	`query:"offset" example:"5" validate:"omitzero,gte=-1"`
}

func (req *SubscriptionUpdateReq) ToMap() map[string]any {
    updateData := make(map[string]any)
    
    if req.ServiceName != nil {
        updateData["service_name"] = *req.ServiceName
    }
    if req.Price != nil {
        updateData["price"] = *req.Price
    }
    if req.StartDate != nil {
        updateData["start_date"] = *(utils.ParseToDate(*req.StartDate))
    }
	if req.EndDate != nil {
        updateData["end_date"] = *(utils.ParseToDate(*req.EndDate))
    }
    
    return updateData
}