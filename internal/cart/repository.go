package cart

import (
	"context"
	"github.com/online-store/internal/domain"
	"github.com/online-store/pkg/database"
)

type Repository interface {
	InsertCartItem(ctx context.Context, data []domain.Cart) error
	FetchWithFilterAndPaginationAndOrderBy(ctx context.Context, page int, pageSize int, query string, countQuery string, orderBy string, model interface{}, args ...interface{}) (*database.Paginator, error)
	DeleteCartItem(ctx context.Context, cartID, customerID int) error
}
