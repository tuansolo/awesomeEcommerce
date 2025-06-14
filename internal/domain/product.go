package domain

import (
	"time"
)

// Product represents a product in the e-commerce system
type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Stock       int       `json:"stock" gorm:"not null"`
	SKU         string    `json:"sku" gorm:"size:50;uniqueIndex;not null"`
	ImageURL    string    `json:"image_url" gorm:"size:255"`
	CategoryID  uint      `json:"category_id"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// ProductCategory represents a category for products
type ProductCategory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	ParentID  *uint     `json:"parent_id" gorm:"default:null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for Product
func (Product) TableName() string {
	return "products"
}

// TableName specifies the table name for ProductCategory
func (ProductCategory) TableName() string {
	return "product_categories"
}
