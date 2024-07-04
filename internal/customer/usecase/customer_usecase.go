package usecase

import (
	"errors"
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/jackc/pgconn"
	"github.com/online-store/internal/customer"
	"github.com/online-store/internal/domain"
	"github.com/online-store/pkg/zaplogger"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type CustomerUseCase struct {
	customerRepo customer.Repository
	zapLogger    zaplogger.Logger
}

func NewCustomerUseCase(customerRepo customer.Repository, zapLogger zaplogger.Logger) customer.UseCase {
	return &CustomerUseCase{
		customerRepo: customerRepo,
		zapLogger:    zapLogger,
	}
}

func (u *CustomerUseCase) InsertCustomer(beegoCtx *beegoContext.Context, req domain.InsertCustomerRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(err))
		return err
	}
	err = u.customerRepo.InsertCustomer(beegoCtx.Request.Context(), domain.Customer{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Password:    string(hashedPassword),
		Address:     req.Address,
		PhoneNumber: req.PhoneNumber,
		CreatedAt:   time.Now(),
		CreatedBy:   "System",
	})

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

	return nil
}

func (u *CustomerUseCase) LoginCustomer(beegoCtx *beegoContext.Context, req domain.LoginRequest) (*domain.Customer, error) {
	user, err := u.customerRepo.GetUserByEmail(beegoCtx.Request.Context(), req.Email)
	if err != nil {
		beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(err))
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(errors.New("invalid credentials")))
		return nil, errors.New("invalid credentials")
	}
	return &user, nil
}
