package domain

import "time"

type (
	Product struct {
		ID           int     `gorm:"column:id" json:"id"`
		Name         string  `gorm:"column:name" json:"name"`
		Description  string  `gorm:"column:description" json:"description"`
		CategoryID   string  `gorm:"column:category_id" json:"category_id"`
		CategoryName string  `gorm:"column:category_name" json:"category_name"`
		Price        float64 `gorm:"column:price" json:"price"`
		Stock        int     `gorm:"column:stock" json:"stock"`

		CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
		CreatedBy string     `gorm:"column:created_by" json:"created_by"`
		UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
		UpdatedBy *string    `gorm:"column:updated_by" json:"updated_by"`
		DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
		DeletedBy *string    `gorm:"column:deleted_by" json:"deleted_by"`
	}

	GetProductListRequest struct {
		Page            int    `json:"-"`
		Limit           int    `json:"-"`
		ProductCategory string `json:"product_category"`
		Search          string `json:"search"`
	}
)
