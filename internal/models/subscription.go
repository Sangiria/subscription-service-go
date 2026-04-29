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
	ServiceName string 		`json:"service_name" create:"required" update:"omitzero"`
	Price		int			`json:"price" create:"required,gt=-1" update:"omitzero,gt=-1"`
	UserId		string		`json:"user_id" create:"required,uuid"`
	StartDate	string		`json:"start_date" create:"required,datetime=01-2006" update:"omitzero,datetime=01-2006"`
	EndDate		string		`json:"end_date,omitzero" create:"omitzero,datetime=01-2006" update:"omitzero,datetime=01-2006"`
}

type SubscriptionUpdateReq struct {
	ServiceName *string 	`json:"service_name" update:"omitzero"`
	Price		*int		`json:"price" update:"omitzero,gt=-1"`
	StartDate	*string		`json:"start_date" update:"omitzero,datetime=01-2006"`
	EndDate		*string		`json:"end_date,omitzero" update:"omitzero,datetime=01-2006"`
}

type SumSubscriptionPrice struct {
	UserID  	string		`validate:"required,uuid"`
	Name    	string     	`validate:"omitzero"`
	StartDate   string 		`validate:"omitzero,datetime=01-2006"`
	EndDate     string 		`validate:"omitzero,datetime=01-2006"`
}

type Subscription struct {
	Id			string			`json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ServiceName string 			`json:"service_name"`
	Price		int				`json:"price"`
	UserId		string			`json:"user_id"`
	StartDate	time.Time		`json:"start_date"`
	EndDate		*time.Time		`json:"end_date,omitzero"`
}

type ListParams struct {
	Limit  int `validate:"gte=0,lte=100"`
	Offset int `validate:"gte=0"`
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