package service

import (
	"context"

	"github.com/luis13005/cleanarchitecture/internal/infra/grpc/pb"
	"github.com/luis13005/cleanarchitecture/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	OrderUseCase *usecase.CreateOrderUseCase
}

func NewOrderService(orderUseCase *usecase.CreateOrderUseCase) *OrderService {
	return &OrderService{OrderUseCase: orderUseCase}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.OrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}
