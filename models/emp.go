package models

type Emp struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	State    int    `json:"state"`
}

func (Emp) TableName() string {
	return "emp"
}
