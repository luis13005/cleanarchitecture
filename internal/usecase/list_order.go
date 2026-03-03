package usecase

import (
	"github.com/luis13005/cleanarchitecture/internal/entity"
	"github.com/luis13005/cleanarchitecture/pkg/events"
)

type ListOrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	List            events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewListOrderUseCase(orderRepository entity.OrderRepositoryInterface,
	list events.EventInterface,
	eventDispatcher events.EventDispatcherInterface) *ListOrderUseCase {
	return &ListOrderUseCase{OrderRepository: orderRepository,
		List:            list,
		EventDispatcher: eventDispatcher}
}

func (useCase *ListOrderUseCase) ListOrder() ([]OrderOutputDTO, error) {
	orders, err := useCase.OrderRepository.List()
	if err != nil {
		return nil, err
	}

	var ordensDTO []OrderOutputDTO

	for _, v := range orders {
		var ordemDTO OrderOutputDTO

		ordemDTO.ID = v.ID
		ordemDTO.Price = v.Price
		ordemDTO.Tax = v.Tax
		ordemDTO.FinalPrice = v.FinalPrice

		ordensDTO = append(ordensDTO, ordemDTO)
	}

	return ordensDTO, nil
}
