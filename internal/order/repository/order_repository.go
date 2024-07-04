package repository

import (
	"context"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/online-store/internal/domain"
	"github.com/online-store/internal/order"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) order.Repository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) DB() *gorm.DB {
	return r.db
}

func (r *OrderRepository) InsertOrder(ctx context.Context, tx *gorm.DB, data domain.Order) (*domain.Order, error) {
	err := tx.WithContext(newrelic.NewContext(ctx, newrelic.FromContext(ctx))).Omit("UpdatedAt", "UpdatedBy", "DeletedAt", "DeletedBy").Create(&data).Error

	return &data, err
}

func (r *OrderRepository) InsertOrderItem(ctx context.Context, tx *gorm.DB, data []domain.OrderItem) error {
	err := tx.WithContext(newrelic.NewContext(ctx, newrelic.FromContext(ctx))).Omit("OrderItemID", "UpdatedAt", "UpdatedBy", "DeletedAt", "DeletedBy").CreateInBatches(&data, 20).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) InsertPayment(ctx context.Context, tx *gorm.DB, data domain.Payment) (*domain.Payment, error) {
	err := tx.WithContext(newrelic.NewContext(ctx, newrelic.FromContext(ctx))).Omit("UpdatedAt", "UpdatedBy", "DeletedAt", "DeletedBy").Create(&data).Error

	return &data, err
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, tx *gorm.DB, paymentID, orderID int) error {
	return tx.WithContext(newrelic.NewContext(ctx, newrelic.FromContext(ctx))).
		Table("order").Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"payment_id": paymentID,
		}).Error
}
