package service

import (
	"context"
	"order_service/internal/models"
)

type OrderRepo interface {
	CreateOrder(ctx context.Context, order models.Order) (*models.Order, error)
	GetOrder(ctx context.Context, id string) (*models.Order, error)
	DeleteOrder(ctx context.Context, id string) error
	ListOrders(ctx context.Context) (*[]models.Order, error)
	UpdateOrder(ctx context.Context, newOrder models.Order) (*models.Order, error)
}

type OrderService struct {
	Repo OrderRepo
}

func NewOrderService(repo OrderRepo) *OrderService {
	return &OrderService{
		Repo: repo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, position models.Order) (*models.Order, error) {
	// Здесь все дополнительные проверки, походы в другие сервисы, склеивание данных из разных баз и тд
	return s.Repo.CreateOrder(ctx, position)
}

func (s *OrderService) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	return s.Repo.GetOrder(ctx, id)
}

func (s *OrderService) DeleteOrder(ctx context.Context, id string) error {
	return s.Repo.DeleteOrder(ctx, id)
}

func (s *OrderService) ListOrders(ctx context.Context) (*[]models.Order, error) {
	return s.Repo.ListOrders(ctx)
}

func (s *OrderService) UpdateOrder(ctx context.Context, newOrder models.Order) (*models.Order, error) {
	return s.Repo.UpdateOrder(ctx, newOrder)
}
