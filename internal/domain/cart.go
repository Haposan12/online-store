package domain

import "time"

type (
	Cart struct {
		ProductID  int `gorm:"column:product_id" json:"product_id"`
		Quantity   int `gorm:"column:quantity" json:"quantity"`
		CustomerID int `gorm:"column:customer_id" json:"customer_id"`

		CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
		CreatedBy string     `gorm:"column:created_by" json:"created_by"`
		UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
		UpdatedBy *string    `gorm:"column:updated_by" json:"updated_by"`
		DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
		DeletedBy *string    `gorm:"column:deleted_by" json:"deleted_by"`
	}

	CreateCartRequest struct {
		CartItem   []CartItem `json:"cart_item" validate:"required,dive"`
		CustomerID int        `json:"customer_id"`
	}

	CartItem struct {
		ProductID int `json:"product_id" validate:"required,number"`
		Quantity  int `json:"quantity" validate:"required,number"`
	}

	GetListCartRequest struct {
		Page       int `json:"-"`
		Limit      int `json:"-"`
		CustomerID int `json:"customer_id"`
	}

	CartProduct struct {
		CartID             int     `gorm:"column:cart_id" json:"cart_id"`
		ProductName        string  `gorm:"column:product_name" json:"product_name"`
		ProductDescription string  `gorm:"column:product_description" json:"product_description"`
		CategoryName       string  `gorm:"column:category_name" json:"category_name"`
		ProductPrice       float64 `gorm:"column:product_price" json:"product_price"`
		Quantity           int     `gorm:"column:quantity" json:"quantity"`
	}
)

func (Cart) TableName() string {
	return "cart"
}
