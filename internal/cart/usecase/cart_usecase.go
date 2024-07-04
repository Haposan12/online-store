package usecase

import (
	"context"
	"fmt"
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/jackc/pgconn"
	jsoniter "github.com/json-iterator/go"
	"github.com/online-store/internal/cart"
	"github.com/online-store/internal/domain"
	"github.com/online-store/pkg/cache"
	"github.com/online-store/pkg/database"
	"github.com/online-store/pkg/zaplogger"
	"strconv"
	"time"
)

type CartUseCase struct {
	cartRepo  cart.Repository
	zapLogger zaplogger.Logger
	cacheRepo cache.RedisRepository
}

func NewCustomerUseCase(cartRepo cart.Repository, zapLogger zaplogger.Logger, cacheRepo cache.RedisRepository) cart.UseCase {
	return &CartUseCase{
		cartRepo:  cartRepo,
		zapLogger: zapLogger,
		cacheRepo: cacheRepo,
	}
}

func (u *CartUseCase) InsertCartItem(beegoCtx *beegoContext.Context, request domain.CreateCartRequest) error {
	var data []domain.Cart
	for _, v := range request.CartItem {
		data = append(data, domain.Cart{
			ProductID:  v.ProductID,
			Quantity:   v.Quantity,
			CustomerID: request.CustomerID,
			CreatedAt:  time.Now(),
			CreatedBy:  "System",
		})
	}
	err := u.cartRepo.InsertCartItem(beegoCtx.Request.Context(), data)
	if err != nil {
		beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(err))
		pgerr, ok := err.(*pgconn.PgError)
		if !ok {
			return err
		}
		switch pgerr.Code {
		case domain.PgCodeForeignKeyConstraint:
			return domain.ErrForeignKeyConstraint
		case domain.PgCodeUniqueConstraint:
			return domain.ErrUniqueConstraint
		default:
			return err
		}
	}

	//delete existing cache
	err = u.cacheRepo.Deletes(beegoCtx.Request.Context(), []string{
		domain.CartKeyCache,
	})

	if err != nil {
		beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(err))
	}

	return nil
}

func (u *CartUseCase) GetListCartItem(beegoCtx *beegoContext.Context, request domain.GetListCartRequest) (*database.Paginator, error) {
	var entities []domain.CartProduct

	cacheKey := fmt.Sprintf("%s:%s", domain.CartKeyCache, fmt.Sprintf("%s|%d|%d", "ALL", request.Page, request.Limit))

	//check cache
	redisResult, err := u.cacheRepo.Fetch(beegoCtx.Request.Context(), cacheKey)
	if err != nil {
		query := `SELECT 
					c.cart_id ,
					p.name as product_name,
					p.description as product_description,
					ca."name" as category_name,
					p.price as product_price,
					c.quantity as quantity
					from cart c 
					join product p ON p.id = c.product_id 
					join category ca on ca.id = p.category_id 
				WHERE p.deleted_at IS NULL AND c.deleted_at IS NULL AND ca.deleted_at IS NULL AND c.customer_id = ?`
		countQuery := `SELECT COUNT(*) from cart c 
					join product p ON p.id = c.product_id 
					join category ca on ca.id = p.category_id 
				WHERE p.deleted_at IS NULL AND c.deleted_at IS NULL AND ca.deleted_at IS NULL AND c.customer_id = ?`

		data, err := u.cartRepo.FetchWithFilterAndPaginationAndOrderBy(
			context.Background(),
			request.Page,
			request.Limit,
			query,
			countQuery,
			"ORDER BY c.created_at DESC",
			&entities,
			request.CustomerID,
		)

		if err != nil {
			beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(err))
			return nil, err
		}

		err = u.cacheRepo.Save(beegoCtx.Request.Context(), cacheKey, data, domain.HalfCacheExpiration)
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

func (u *CartUseCase) DeleteCartItem(beegoCtx *beegoContext.Context, cartIDReq string, customerIDReq int) error {
	cartID, err := strconv.Atoi(cartIDReq)
	if err != nil {
		beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(err))
		return err
	}

	err = u.cartRepo.DeleteCartItem(beegoCtx.Request.Context(), cartID, customerIDReq)
	if err != nil {
		beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(err))
		return err
	}

	//delete existing cache
	err = u.cacheRepo.Deletes(beegoCtx.Request.Context(), []string{
		domain.CartKeyCache,
	})

	if err != nil {
		beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(err))
	}

	return nil
}
