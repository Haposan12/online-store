package order

import (
	"context"
	"github.com/online-store/internal/domain"
	"gorm.io/gorm"
)

type Repository interface {
	DB() *gorm.DB
	InsertOrder(ctx context.Context, tx *gorm.DB, data domain.Order) (*domain.Order, error)
	InsertOrderItem(ctx context.Context, tx *gorm.DB, data []domain.OrderItem) error
	InsertPayment(ctx context.Context, tx *gorm.DB, data domain.Payment) (*domain.Payment, error)
	UpdateOrder(ctx context.Context, tx *gorm.DB, paymentID, orderID int) error
}
