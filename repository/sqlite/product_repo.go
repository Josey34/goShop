package sqlite

import (
	"context"
	"database/sql"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/errors"
	"github.com/Josey34/goshop/domain/valueobject"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

func (r *ProductRepo) Create(ctx context.Context, product *entity.Product) error {
	query := `INSERT INTO products (id, name, description, price, stock, image_url, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, product.ID, product.Name, product.Description, product.Price.Value(), product.Stock, product.ImageURL, product.CreatedAt, product.UpdatedAt)
	return err
}

func (r *ProductRepo) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	query := `SELECT id, name, description, price, stock, image_url, created_at, updated_at
              FROM products WHERE id = ?`

	row := r.db.QueryRowContext(ctx, query, id)

	var p entity.Product
	var price int64

	err := row.Scan(
		&p.ID, &p.Name, &p.Description,
		&price, &p.Stock, &p.ImageURL,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFound("product", id)
	}
	if err != nil {
		return nil, err
	}

	p.Price = valueobject.NewMoney(price)
	return &p, nil
}

func (r *ProductRepo) FindAll(ctx context.Context, pagination valueobject.Pagination) ([]*entity.Product, error) {
	query := `SELECT id, name, description, price, stock, image_url, created_at, updated_at
              FROM products LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, pagination.Limit, pagination.Offset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var p entity.Product
		var price int64
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &price, &p.Stock, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		p.Price = valueobject.NewMoney(price)
		products = append(products, &p)
	}
	return products, nil
}

func (r *ProductRepo) Update(ctx context.Context, product *entity.Product) error {
	query := `UPDATE products SET name=?, description=?, price=?, stock=?, image_url=?, updated_at=?
              WHERE id=?`

	_, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price.Value(),
		product.Stock,
		product.ImageURL,
		product.UpdatedAt,
		product.ID,
	)
	return err
}

func (r *ProductRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
