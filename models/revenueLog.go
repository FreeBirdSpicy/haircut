package models

type RevenueLog struct {
	Id       int    `json:"id" form:"id"`
	TicketId int    `json:"t_id" form:"t_id" gorm:"column:t_id"`
	Price    string `json:"price" form:"price"`
	EmpId    int    `json:"emp_id" form:"emp_id"`
	UserId   int    `json:"user_id" form:"user_id"`
	Plat     int    `json:"plat" form:"plat" gorm:"default:1"`
	State    int    `json:"state" form:"state" gorm:"default:1"`
	Dated    string `json:"dated" form:"dated"`
}

func (RevenueLog) TableName() string {
	return "revenue_log"
}

// BeforeCreate 在GORM创建记录前设置默认值。注意：BeforeCreate只在通过 GORM 的 Create 方法创建记录时触发
/* func (r *RevenueLog) BeforeCreate(tx *gorm.DB) error {
	if r.Plat == 0 {
		r.Plat = 1 // 设置 Plat 默认值为 1
	}
	if r.State == 0 {
		r.State = 1 // 设置 State 默认值为 1
	}
	return nil
} */
