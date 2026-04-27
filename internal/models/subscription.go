package models

type Subscription struct {
	ServiceName string 		`json:"service_name" create:"required"`
	Price		int			`json:"price" create:"required,gt=0"`
	UserId		string		`json:"user_id" create:"required,uuid"`
	StartDate	string		`json:"start_date" create:"required,datetime=01-2006"`
	EndDate		string		`json:"end_date,omitzero" create:"omitzero,datetime=01-2006"`
}

const (
    TagCreate = "create"
    TagUpdate = "update"
)