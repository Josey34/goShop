package sqlite

import (
	"context"
	"database/sql"

	"github.com/Josey34/goshop/repository"
)

type UnitOfWork struct {
	db           *sql.DB
	tx           *sql.Tx
	productRepo  *ProductRepo
	orderRepo    *OrderRepo
	customerRepo *CustomerRepo
}

func NewUnitOfWork(db *sql.DB) *UnitOfWork {
	return &UnitOfWork{
		db: db,
	}
}

func (u *UnitOfWork) Begin(ctx context.Context) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	u.tx = tx
	u.productRepo = NewProductRepo(tx)
	u.orderRepo = NewOrderRepo(tx)
	u.customerRepo = NewCustomerRepo(tx)

	return nil
}

func (u *UnitOfWork) Commit() error {
	return u.tx.Commit()
}

func (u *UnitOfWork) Rollback() error {
	return u.tx.Rollback()
}

func (u *UnitOfWork) ProductRepository() repository.ProductRepository {
	return u.productRepo
}

func (u *UnitOfWork) OrderRepository() repository.OrderRepository {
	return u.orderRepo
}

func (u *UnitOfWork) CustomerRepository() repository.CustomerRepository {
	return u.customerRepo
}
