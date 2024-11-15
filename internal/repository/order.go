package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"order_service/internal/models"
	"order_service/pkg/db/postgres"
)

const (
	errNoOrder = "no order with this Id"
)

type OrderRepository struct {
	db *postgres.DB
}

func NewOrderRepository(db *postgres.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order models.Order) (*models.Order, error) {
	var result models.Order
	err := sq.Insert("orders").
		Columns("item", "quantity").
		Values(order.Item, order.Quantity).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db.Db).
		QueryRow().
		Scan(&result.ID, &result.Item, &result.Quantity)
	if err != nil {
		return nil, fmt.Errorf("repository.CreateOrder: %w", err)
	}

	return &result, nil
}

func (r *OrderRepository) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	var result models.Order

	err := sq.Select("*").
		From("orders").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db.Db).
		QueryRow().
		Scan(&result.ID, &result.Item, &result.Quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(errNoOrder)
		}
		return nil, fmt.Errorf("repository.GetOrder: %w", err)
	}

	return &result, nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, newOrder models.Order) (*models.Order, error) {
	_, err := sq.Update("orders").
		Set("quantity", newOrder.Quantity).
		Set("item", newOrder.Item).
		Where(sq.Eq{"id": newOrder.ID}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db.Db).
		Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(errNoOrder)
		}
		return nil, fmt.Errorf("repository.UpdateOrder: %w", err)
	}

	return &newOrder, nil
}

func (r *OrderRepository) DeleteOrder(ctx context.Context, id string) error {

	_, err := sq.Delete("orders").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db.Db).
		Query()
	fmt.Println(err)
	if err != nil {
		return fmt.Errorf("repository.DeleteOrder: %w", err)
	}

	return nil
}

func (r *OrderRepository) ListOrders(ctx context.Context) (*[]models.Order, error) {
	var orders []models.Order
	rows, err := sq.Select("*").
		From("orders").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db.Db).
		Query()
	if err != nil {
		return nil, fmt.Errorf("repository.ListOrders: %w", err)
	}

	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.ID, &order.Item, &order.Quantity)
		if err != nil {
			return nil, fmt.Errorf("repository.ListOrders: %w", err)
		}

		orders = append(orders, order)
	}

	return &orders, nil
}
