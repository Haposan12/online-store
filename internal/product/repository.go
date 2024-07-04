package product

import (
	"context"
	"github.com/online-store/pkg/database"
)

type Repository interface {
	FetchWithFilterAndPaginationAndOrderBy(ctx context.Context, page int, pageSize int, query string, countQuery string, orderBy string, model interface{}, args ...interface{}) (*database.Paginator, error)
}
