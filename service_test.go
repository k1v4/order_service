package _83984_chvik1_course_1195

import (
	"context"
	"google.golang.org/grpc"
	pb "order_service/pkg/api/3_1" // Путь к сгенерированным protobuf-файлам
	"testing"
)

func TestOrderService(t *testing.T) {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure()) // замените порт и адрес на свои
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)

	ctx := context.Background()

	// Test CreateOrder
	t.Run("CreateOrder", func(t *testing.T) {
		req := &pb.CreateOrderRequest{
			Item:     "Test Item",
			Quantity: 5,
		}
		resp, err := client.CreateOrder(ctx, req)
		if err != nil {
			t.Fatalf("CreateOrder failed: %v", err)
		}
		if resp.Id == "" {
			t.Errorf("Expected non-empty order ID, got empty")
		}
	})

	t.Run("GetOrder", func(t *testing.T) {
		createResp, err := client.CreateOrder(ctx, &pb.CreateOrderRequest{
			Item:     "Test Item",
			Quantity: 5,
		})
		if err != nil {
			t.Fatalf("Failed to create order: %v", err)
		}

		req := &pb.GetOrderRequest{Id: createResp.Id}
		resp, err := client.GetOrder(ctx, req)
		if err != nil {
			t.Fatalf("GetOrder failed: %v", err)
		}
		if resp.Order.Id != createResp.Id {
			t.Errorf("Expected order ID %s, got %s", createResp.Id, resp.Order.Id)
		}
	})

	t.Run("UpdateOrder", func(t *testing.T) {
		createResp, err := client.CreateOrder(ctx, &pb.CreateOrderRequest{
			Item:     "Test Item",
			Quantity: 5,
		})
		if err != nil {
			t.Fatalf("Failed to create order: %v", err)
		}

		req := &pb.UpdateOrderRequest{
			Id:       createResp.Id,
			Item:     "Updated Item",
			Quantity: 10,
		}
		resp, err := client.UpdateOrder(ctx, req)
		if err != nil {
			t.Fatalf("UpdateOrder failed: %v", err)
		}
		if resp.Order.Item != req.Item || resp.Order.Quantity != req.Quantity {
			t.Errorf("UpdateOrder returned incorrect data: got %+v", resp.Order)
		}
	})

	t.Run("DeleteOrder", func(t *testing.T) {
		createResp, err := client.CreateOrder(ctx, &pb.CreateOrderRequest{
			Item:     "Test Item",
			Quantity: 5,
		})
		if err != nil {
			t.Fatalf("Failed to create order: %v", err)
		}

		req := &pb.DeleteOrderRequest{Id: createResp.Id}
		resp, err := client.DeleteOrder(ctx, req)
		if err != nil {
			t.Fatalf("DeleteOrder failed: %v", err)
		}
		if !resp.Success {
			t.Errorf("Expected DeleteOrder to succeed, but it failed")
		}
	})

	t.Run("ListOrders", func(t *testing.T) {
		_, _ = client.CreateOrder(ctx, &pb.CreateOrderRequest{
			Item:     "Item 1",
			Quantity: 3,
		})
		_, _ = client.CreateOrder(ctx, &pb.CreateOrderRequest{
			Item:     "Item 2",
			Quantity: 4,
		})

		resp, err := client.ListOrders(ctx, &pb.ListOrdersRequest{})
		if err != nil {
			t.Fatalf("ListOrders failed: %v", err)
		}
		if len(resp.Orders) < 2 {
			t.Errorf("Expected at least 2 orders, got %d", len(resp.Orders))
		}
	})
}
