package usecase

import (
	"github.com/flpgst/golang-studies/55-CleanArch/internal/dto"
	"github.com/flpgst/golang-studies/55-CleanArch/internal/entity"
)

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrderUseCase(OrderRepository entity.OrderRepositoryInterface) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrderUseCase) Execute() ([]dto.OrderOutputDTO, error) {
	orders, err := c.OrderRepository.List()
	if err != nil {
		return []dto.OrderOutputDTO{}, err
	}
	dtos := []dto.OrderOutputDTO{}
	for _, order := range orders {
		dto := dto.OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}
