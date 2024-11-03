package repository

import (
	"context"
	"fmt"
	"order_service/internal/models"
	"sync"
	"sync/atomic"
)

const (
	errNoOrder = "no order with this Id"
)

type OrderRepository struct {
	db        map[string]models.Order
	rw        *sync.RWMutex
	idCounter int64
}

func NewOrderRepository(mapa map[string]models.Order, rw *sync.RWMutex) *OrderRepository {
	return &OrderRepository{mapa, rw, 0}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order models.Order) (*models.Order, error) {
	r.rw.Lock()
	defer r.rw.Unlock()
	order.ID = fmt.Sprintf("%d", atomic.AddInt64(&r.idCounter, 1))

	if _, ok := r.db[order.ID]; !ok {
		r.db[order.ID] = order
	}

	return &order, nil
}

func (r *OrderRepository) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	order := models.Order{}
	r.rw.Lock()
	defer r.rw.Unlock()

	if v, ok := r.db[id]; ok {
		order = v
	} else {
		return nil, fmt.Errorf(errNoOrder)
	}

	return &order, nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, newOrder models.Order) (*models.Order, error) {
	r.rw.Lock()
	defer r.rw.Unlock()

	id := newOrder.ID

	if _, ok := r.db[id]; ok {
		r.db[id] = newOrder
	} else {
		return nil, fmt.Errorf(errNoOrder)
	}

	return &newOrder, nil
}

func (r *OrderRepository) DeleteOrder(ctx context.Context, id string) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	if _, ok := r.db[id]; ok {
		delete(r.db, id)
	} else {
		return fmt.Errorf(errNoOrder)
	}

	return nil
}

func (r *OrderRepository) ListOrders(ctx context.Context) (*[]models.Order, error) {
	var orders []models.Order
	r.rw.Lock()
	defer r.rw.Unlock()

	for _, v := range r.db {
		orders = append(orders, v)
	}

	return &orders, nil
}
