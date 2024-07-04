package repository

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/online-store/internal/cart"
	"github.com/online-store/internal/domain"
	"github.com/online-store/pkg/database"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) cart.Repository {
	return &CartRepository{db: db}
}

func (r *CartRepository) InsertCartItem(ctx context.Context, data []domain.Cart) error {
	err := r.db.WithContext(newrelic.NewContext(ctx, newrelic.FromContext(ctx))).Omit("CartID", "UpdatedAt", "UpdatedBy", "DeletedAt", "DeletedBy").CreateInBatches(&data, 20).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CartRepository) FetchWithFilterAndPaginationAndOrderBy(ctx context.Context, page int, pageSize int, query string, countQuery string, orderBy string, model interface{}, args ...interface{}) (*database.Paginator, error) {
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

func (r *CartRepository) DeleteCartItem(ctx context.Context, cartID, customerID int) error {
	return r.db.WithContext(newrelic.NewContext(ctx, newrelic.FromContext(ctx))).
		Table("cart").Where("cart_id = ? AND customer_id = ?", cartID, customerID).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
			"deleted_by": "System",
		}).Error
}
