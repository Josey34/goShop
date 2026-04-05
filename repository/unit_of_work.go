package repository

import "context"

type UnitOfWork interface {
	Begin(ctx context.Context) error
	Commit() error
	Rollback() error
	ProductRepository() ProductRepository
	OrderRepository() OrderRepository
	CustomerRepository() CustomerRepository
}
