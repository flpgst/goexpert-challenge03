package usecase

import (
	"github.com/flpgst/golang-studies/55-CleanArch/internal/dto"
	"github.com/flpgst/golang-studies/55-CleanArch/internal/entity"
	"github.com/flpgst/golang-studies/55-CleanArch/pkg/events"
)

type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispacherInterface
}

func NewCreateOrderUseCase(
	orderRepository entity.OrderRepositoryInterface,
	orderCreated events.EventInterface,
	eventDispatcher events.EventDispacherInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: orderRepository,
		OrderCreated:    orderCreated,
		EventDispatcher: eventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input dto.OrderInputDTO) (dto.OrderOutputDTO, error) {
	order, err := entity.NewOrder(input.ID, input.Price, input.Tax)
	if err != nil {
		return dto.OrderOutputDTO{}, err
	}
	order.CalculateFinalPrice()
	if err := c.OrderRepository.Save(order); err != nil {
		return dto.OrderOutputDTO{}, err
	}
	dto := dto.OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}
	c.OrderCreated.SetPayload(dto)
	c.EventDispatcher.Dispatch(c.OrderCreated)

	return dto, nil
}
