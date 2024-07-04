package repository

import (
	"context"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/online-store/internal/customer"
	"github.com/online-store/internal/domain"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) customer.Repository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) InsertCustomer(ctx context.Context, request domain.Customer) error {
	err := r.db.WithContext(newrelic.NewContext(ctx, newrelic.FromContext(ctx))).Omit("CustomerID", "UpdatedAt", "UpdatedBy", "DeletedAt", "DeletedBy").Create(&request).Error
	return err
}

func (r *CustomerRepository) GetUserByEmail(ctx context.Context, email string) (domain.Customer, error) {
	var data domain.Customer

	result := r.db.WithContext(newrelic.NewContext(ctx, newrelic.FromContext(ctx))).Where("email = ?", email).First(&data)
	return data, result.Error
}
