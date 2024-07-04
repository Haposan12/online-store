package domain

import "time"

type (
	CreateOrderCheckoutRequest struct {
		Order      []OrderRequest `json:"order" validate:"required,dive"`
		CustomerID int            `json:"customer_id"`
	}

	OrderRequest struct {
		ProductID int     `json:"product_id" validate:"required,number"`
		Quantity  int     `json:"quantity" validate:"required,number"`
		Price     float64 `json:"price" validate:"required,number"`
	}

	OrderItem struct {
		ProductID int     `json:"product_id"`
		Price     float64 `json:"price"`
		Quantity  int     `json:"quantity"`
		OrderID   int     `json:"order_id"`

		CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
		CreatedBy string     `gorm:"column:created_by" json:"created_by"`
		UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
		UpdatedBy *string    `gorm:"column:updated_by" json:"updated_by"`
		DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
		DeletedBy *string    `gorm:"column:deleted_by" json:"deleted_by"`
	}

	Order struct {
		ID         int     `gorm:"column:id" json:"id"`
		TotalPrice float64 `json:"total_price"`
		CustomerID int     `json:"customer_id"`

		CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
		CreatedBy string     `gorm:"column:created_by" json:"created_by"`
		UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
		UpdatedBy *string    `gorm:"column:updated_by" json:"updated_by"`
		DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
		DeletedBy *string    `gorm:"column:deleted_by" json:"deleted_by"`
	}

	PaymentRequest struct {
		OrderID string  `json:"order_id"`
		Amount  float64 `json:"amount" validate:"required,number"`
		Method  string  `json:"method" validate:"required"`
	}

	Payment struct {
		ID     int     `json:"id"`
		Method string  `json:"method"`
		Amount float64 `json:"amount"`
		Status string  `json:"status"`

		CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
		CreatedBy string     `gorm:"column:created_by" json:"created_by"`
		UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
		UpdatedBy *string    `gorm:"column:updated_by" json:"updated_by"`
		DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
		DeletedBy *string    `gorm:"column:deleted_by" json:"deleted_by"`
	}
)

func (OrderItem) TableName() string {
	return "order_item"
}
func (Order) TableName() string {
	return "order"
}
func (Payment) TableName() string {
	return "payment"
}
