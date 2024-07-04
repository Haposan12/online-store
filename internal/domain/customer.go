package domain

import "time"

type (
	Customer struct {
		CustomerID  int    `gorm:"column:customer_id" json:"customer_id"`
		FirstName   string `gorm:"column:first_name" json:"first_name"`
		LastName    string `gorm:"column:last_name" json:"last_name"`
		Email       string `gorm:"column:email" json:"email"`
		Password    string `gorm:"column:password" json:"password"`
		Address     string `gorm:"column:address" json:"address"`
		PhoneNumber string `gorm:"column:phone_number" json:"phone_number"`

		CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
		CreatedBy string     `gorm:"column:created_by" json:"created_by"`
		UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
		UpdatedBy *string    `gorm:"column:updated_by" json:"updated_by"`
		DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
		DeletedBy *string    `gorm:"column:deleted_by" json:"deleted_by"`
	}

	InsertCustomerRequest struct {
		FirstName   string `json:"first_name" validate:"required,name"`
		LastName    string `json:"last_name" validate:"required,name"`
		Email       string `json:"email" validate:"required,email_address"`
		Password    string `json:"password" validate:"required"`
		Address     string `json:"address" validate:"required,address"`
		PhoneNumber string `json:"phone_number" validate:"required,number,len=12"`
	}

	LoginRequest struct {
		Email    string `json:"email" validate:"required,email_address"`
		Password string `json:"password" validate:"required"`
	}
)

func (Customer) TableName() string {
	return "customer"
}
