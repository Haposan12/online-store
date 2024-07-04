package usecase

import (
	"context"
	"fmt"
	beegoContext "github.com/beego/beego/v2/server/web/context"
	jsoniter "github.com/json-iterator/go"
	"github.com/online-store/internal/domain"
	"github.com/online-store/internal/product"
	"github.com/online-store/pkg/cache"
	"github.com/online-store/pkg/database"
	"github.com/online-store/pkg/zaplogger"
)

type ProductUseCase struct {
	productRepo product.Repository
	cacheRepo   cache.RedisRepository
	zapLogger   zaplogger.Logger
}

func NewProductUseCase(productRepo product.Repository, cacheRepo cache.RedisRepository, zapLogger zaplogger.Logger) product.UseCase {
	return &ProductUseCase{
		productRepo: productRepo,
		cacheRepo:   cacheRepo,
		zapLogger:   zapLogger,
	}
}

func (u ProductUseCase) GetListProduct(beegoCtx *beegoContext.Context, req domain.GetProductListRequest) (*database.Paginator, error) {
	var entities []domain.Product
	cacheKey := fmt.Sprintf("%s:%s", domain.ProductKeyCache, fmt.Sprintf("%s|%d|%d|%s|%s", "ALL", req.Page, req.Limit, req.ProductCategory, req.Search))

	//check cache
	redisResult, err := u.cacheRepo.Fetch(beegoCtx.Request.Context(), cacheKey)
	if err != nil {
		query := `SELECT 
					p.id, p."name", p.description, p.category_id, c."name" AS category_name, p.price, p.stock, p.created_at, p.created_by, p.updated_at, p.updated_by, p.deleted_at, p.deleted_by 
				FROM product p
				JOIN category c ON p.category_id = c.id
				WHERE p.deleted_at IS NULL AND c.deleted_at IS NULL`
		countQuery := `SELECT COUNT(*) from product WHERE 1=1`

		if req.Search != "" {
			query += ` AND p."name" LIKE = ?`
			countQuery += ` AND p."name" LIKE = ?`
		}

		if req.ProductCategory != "" {
			query += ` AND c."name" LIKE = ?`
			countQuery += ` AND c."name" LIKE = ?`
		}

		data, err := u.productRepo.FetchWithFilterAndPaginationAndOrderBy(
			context.Background(),
			req.Page,
			req.Limit,
			query,
			countQuery,
			"",
			&entities,
		)

		if err != nil {
			beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(err))
			return nil, err
		}

		err = u.cacheRepo.Save(beegoCtx.Request.Context(), cacheKey, *data, domain.HalfCacheExpiration)
		if err != nil {
			beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(err))
		}

		return data, nil
	}

	var paginator = new(database.Paginator)
	if err := jsoniter.UnmarshalFromString(*redisResult, paginator); err != nil {
		return nil, err
	}

	return paginator, nil
}
