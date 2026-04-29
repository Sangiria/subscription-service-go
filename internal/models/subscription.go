package models

import "time"

type SubscriptionReq struct {
	ServiceName string 		`json:"service_name" create:"required"`
	Price		int			`json:"price" create:"required,gt=0"`
	UserId		string		`json:"user_id" create:"required,uuid"`
	StartDate	string		`json:"start_date" create:"required,datetime=01-2006"`
	EndDate		string		`json:"end_date,omitzero" create:"omitzero,datetime=01-2006"`
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
	Limit  int `query:"limit" validate:"gte=-2,lte=100"`
	Offset int `query:"offset" validate:"gte=-2"`
}

const (
    TagCreate = "create"
    TagUpdate = "update"
)
