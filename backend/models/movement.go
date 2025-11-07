package models

import "time"

type Movement struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ProductID    uint      `gorm:"not null" json:"product_id"`
	Product      Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type         string    `gorm:"type:enum('entrada','salida');not null" json:"type"`
	Quantity     int       `gorm:"not null" json:"quantity"`
	Description  string    `json:"description"`
	MovementDate time.Time `gorm:"autoCreateTime" json:"movement_date"`
}
