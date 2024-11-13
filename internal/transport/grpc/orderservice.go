package grpc

import (
	"context"
	"fmt"
	"github.com/AlekSi/pointer"
	"order_service/internal/models"
	order "order_service/pkg/api/3_1"
)

type Service interface {
	CreateOrder(ctx context.Context, order models.Order) (*models.Order, error)
	GetOrder(ctx context.Context, id string) (*models.Order, error)
	DeleteOrder(ctx context.Context, id string) error
	ListOrders(ctx context.Context) (*[]models.Order, error)
	UpdateOrder(ctx context.Context, newOrder models.Order) (*models.Order, error)
}

type OrderService struct {
	order.UnimplementedOrderServiceServer
	service Service
}

func NewOrderService(s Service) *OrderService {
	return &OrderService{service: s}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	resp, err := s.service.CreateOrder(ctx, models.Order{
		Item:     req.GetItem(),
		Quantity: req.GetQuantity(),
	})
	if err != nil {
		return nil, fmt.Errorf("CreateOrder: %w", err)
	}

	r := pointer.Get(resp)

	return &order.CreateOrderResponse{
		Id: r.ID,
	}, nil
}

func (s *OrderService) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	resp, err := s.service.GetOrder(ctx, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("GetOrder: %w", err)
	}

	r := pointer.Get(resp)

	return &order.GetOrderResponse{
		Order: &order.Order{
			Id:       r.ID,
			Item:     r.Item,
			Quantity: r.Quantity,
		},
	}, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, req *order.DeleteOrderRequest) (*order.DeleteOrderResponse, error) {
	err := s.service.DeleteOrder(ctx, req.GetId())
	if err != nil {
		return &order.DeleteOrderResponse{Success: false}, fmt.Errorf("DeleteOrder: %w", err)
	}

	return &order.DeleteOrderResponse{Success: true}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, req *order.ListOrdersRequest) (*order.ListOrdersResponse, error) {
	resp, err := s.service.ListOrders(ctx)
	if err != nil {
		fmt.Printf("ListOrders: %s", err)
		return nil, fmt.Errorf("ListOrders: %w", err)
	}

	r := pointer.Get(resp)
	var res []*order.Order

	for _, o := range r {
		res = append(res, &order.Order{
			Id:       o.ID,
			Item:     o.Item,
			Quantity: o.Quantity,
		})
	}

	return &order.ListOrdersResponse{
		Orders: res,
	}, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.UpdateOrderResponse, error) {
	resp, err := s.service.UpdateOrder(ctx, models.Order{
		ID:       req.GetId(),
		Item:     req.GetItem(),
		Quantity: req.GetQuantity(),
	})
	if err != nil {
		return nil, fmt.Errorf("UpdateOrder: %w", err)
	}

	r := pointer.Get(resp)

	return &order.UpdateOrderResponse{Order: &order.Order{
		Id:       r.ID,
		Item:     r.Item,
		Quantity: r.Quantity,
	}}, nil
}
