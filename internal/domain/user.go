package domain

import (
	"time"
)

// User represents a user in the e-commerce system
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"size:255;not null"` // Password is not exposed in JSON
	FirstName string    `json:"first_name" gorm:"size:100"`
	LastName  string    `json:"last_name" gorm:"size:100"`
	Phone     string    `json:"phone" gorm:"size:20"`
	Address   string    `json:"address" gorm:"type:text"`
	Role      string    `json:"role" gorm:"size:20;default:'customer'"`
	Cart      Cart      `json:"cart" gorm:"foreignKey:UserID"`
	Orders    []Order   `json:"orders" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}
