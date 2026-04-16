package sqlite

import (
	"context"
	"database/sql"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/errors"
	"github.com/Josey34/goshop/domain/valueobject"
)

type CustomerRepo struct {
	db DBTX
}

func NewCustomerRepo(db DBTX) *CustomerRepo {
	return &CustomerRepo{
		db: db,
	}
}

func (r *CustomerRepo) Create(ctx context.Context, customer *entity.Customer) error {
	query := `INSERT INTO customers (id, name, email, phone, street, city, province, postal_code, password_hash, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		customer.ID,
		customer.Name,
		customer.Email.Value(),
		customer.Phone.Value(),
		customer.Address.Street,
		customer.Address.City,
		customer.Address.Province,
		customer.Address.PostalCode,
		customer.PasswordHash,
		customer.CreatedAt,
		customer.UpdatedAt,
	)
	return err
}

func (r *CustomerRepo) FindByID(ctx context.Context, id string) (*entity.Customer, error) {
	query := `SELECT id, name, email, phone, street, city, province, postal_code, password_hash, created_at, updated_at
              FROM customers WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)

	var c entity.Customer

	var email, phone, street, city, province, postalCode string

	err := row.Scan(
		&c.ID,
		&c.Name,
		&email,
		&phone,
		&street,
		&city,
		&province,
		&postalCode,
		&c.PasswordHash,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFound("customer", id)
	}

	if err != nil {
		return nil, err
	}

	c.Email, _ = valueobject.NewEmail(email)
	c.Phone, _ = valueobject.NewPhone(phone)
	c.Address, _ = valueobject.NewAddress(street, city, province, postalCode)

	return &c, nil
}

func (r *CustomerRepo) FindByEmail(ctx context.Context, email string) (*entity.Customer, error) {
	query := `SELECT id, name, email, phone, street, city, province, postal_code, password_hash, created_at, updated_at
              FROM customers WHERE email = ?`

	row := r.db.QueryRowContext(ctx, query, email)

	var c entity.Customer
	var emailVal, phone, street, city, province, postalCode string

	err := row.Scan(
		&c.ID, &c.Name, &emailVal, &phone,
		&street, &city, &province, &postalCode,
		&c.PasswordHash, &c.CreatedAt, &c.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFound("customer", email)
	}
	if err != nil {
		return nil, err
	}

	c.Email, _ = valueobject.NewEmail(emailVal)
	c.Phone, _ = valueobject.NewPhone(phone)
	c.Address, _ = valueobject.NewAddress(street, city, province, postalCode)

	return &c, nil
}

func (r *CustomerRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM customers WHERE email = ?`, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
