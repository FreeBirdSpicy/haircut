package models

type Ticket struct {
	Id    int     `json:"id" form:"id"`
	Name  string  `json:"name" form:"name"`
	Price float64 `json:"price" form:"price"`
	State int     `json:"state" form:"state"`
}

func (Ticket) TableName() string {
	return "ticket"
}
