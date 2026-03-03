package web

import (
	"encoding/json"
	"net/http"

	"github.com/luis13005/cleanarchitecture/internal/entity"
	"github.com/luis13005/cleanarchitecture/internal/usecase"
	"github.com/luis13005/cleanarchitecture/pkg/events"
)

type WebOrderHandler struct {
	EventDispatcher   events.EventDispatcherInterface
	OrderRepository   entity.OrderRepositoryInterface
	OrderCreatedEvent events.EventInterface
}

func NewWebOrderHandler(
	eventDispatcher events.EventDispatcherInterface,
	orderRepository entity.OrderRepositoryInterface,
	orderCreatedEvent events.EventInterface) *WebOrderHandler {
	return &WebOrderHandler{
		EventDispatcher:   eventDispatcher,
		OrderRepository:   orderRepository,
		OrderCreatedEvent: orderCreatedEvent,
	}
}

func (web *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createOrder := usecase.NewCreateOrderUseCase(web.OrderRepository, web.OrderCreatedEvent, web.EventDispatcher)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (web *WebOrderHandler) ListarOrdens(w http.ResponseWriter, r *http.Request) {
	usecase := usecase.NewListOrderUseCase(web.OrderRepository, web.OrderCreatedEvent, web.EventDispatcher)
	ordens, err := usecase.ListOrder()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(ordens)
}
