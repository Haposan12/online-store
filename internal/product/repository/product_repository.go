package repository

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/online-store/internal/product"
	"github.com/online-store/pkg/database"
	"gorm.io/gorm"
	"reflect"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) product.Repository {
	return &ProductRepository{db: db}
}

func (r ProductRepository) FetchWithFilterAndPaginationAndOrderBy(ctx context.Context, page int, pageSize int, query string, countQuery string, orderBy string, model interface{}, args ...interface{}) (*database.Paginator, error) {
	linq.From(args).Where(func(item interface{}) bool {
		if reflect.TypeOf(item).Kind() == reflect.Slice {
			return len(item.([]string)) != 0
		}
		return item != ""
	}).ToSlice(&args)
	paginate := database.NewPaginator(r.db, page, pageSize, model).Raw(query, args, countQuery, args)

	if err := paginate.FindWithOrderBy(ctx, orderBy).Error; err != nil {
		return paginate, err
	}
	return paginate, nil
}
