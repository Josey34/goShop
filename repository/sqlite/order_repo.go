package sqlite

import (
	"context"
	"time"

	"github.com/Josey34/goshop/domain/entity"
	"github.com/Josey34/goshop/domain/valueobject"
)

type OrderRepo struct {
	db DBTX
}

func NewOrderRepo(db DBTX) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(ctx context.Context, order *entity.Order) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO orders (id, customer_id, total, status, created_at, updated_at) VALUES (?,?,?,?,?,?)`,
		order.ID, order.CustomerID, order.Total.Value(), string(order.Status), order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = r.db.ExecContext(ctx,
			`INSERT INTO order_items (id, order_id, product_id, name, price, quantity) VALUES (?,?,?,?,?,?)`,
			item.ID, item.OrderID, item.ProductID, item.Name, item.Price.Value(), item.Quantity,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *OrderRepo) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	var total int64
	var status string

	order := new(entity.Order)

	err := r.db.QueryRowContext(ctx,
		`SELECT id, customer_id, total, status, created_at, updated_at FROM orders WHERE id = ?`,
		id,
	).Scan(&order.ID, &order.CustomerID, &total, &status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	order.Total = valueobject.NewMoney(total)
	order.Status = valueobject.OrderStatus(status)

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, order_id, product_id, name, price, quantity FROM order_items WHERE order_id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var item entity.OrderItem
		var price int64
		rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Name, &price, &item.Quantity)
		item.Price = valueobject.NewMoney(price)
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, id string, status valueobject.OrderStatus) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE orders SET status=?, updated_at=? WHERE id=?`,
		string(status), time.Now(), id,
	)
	return err
}

func (r *OrderRepo) FindByCustomer(ctx context.Context, customerID string, pagination valueobject.Pagination) ([]*entity.Order, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, customer_id, total, status, created_at, updated_at FROM orders WHERE customer_id = ? LIMIT ? OFFSET ?`,
		customerID, pagination.Limit, pagination.Offset,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []*entity.Order

	for rows.Next() {
		var total int64
		var status string
		order := new(entity.Order)

		err := rows.Scan(&order.ID, &order.CustomerID, &total, &status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}

		order.Total = valueobject.NewMoney(total)
		order.Status = valueobject.OrderStatus(status)

		orders = append(orders, order)
	}

	return orders, nil
}
