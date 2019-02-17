package model

// Order is the database representation of an order
type Order struct {
	ID          int64  `gorm:"type:int" json:"id"`
	Description string `gorm:"type:varchar" json:"description"`
	Ts          int64  `gorm:"timestamp" json:"ts"`
}
