package customer

import (
	"context"
	"github.com/online-store/internal/domain"
)

type Repository interface {
	InsertCustomer(ctx context.Context, request domain.Customer) error
	GetUserByEmail(ctx context.Context, email string) (domain.Customer, error)
}
