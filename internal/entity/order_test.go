package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGivenAnEmptyID_WhenCreateNewOrder_ShouldReceiveAnError(t *testing.T) {
	order := Order{}
	assert.Error(t, order.IsValid(), "invalid id")
}

func TestGivenAnEmptyPrice_WhenCreateNewOrder_ShouldReceiveAnError(t *testing.T) {
	order := Order{
		ID:    "123",
		Price: 0,
	}
	assert.Error(t, order.IsValid(), "invalid price")
}

func TestGivenAnEmptyTax_WhenCreateNewOrder_ShouldReceiveAnError(t *testing.T) {
	order := Order{
		ID:    "123",
		Price: 10,
		Tax:   0,
	}
	assert.Error(t, order.IsValid(), "invalid tax")
}

func TestGivenValidParams_WhenNewOrderIsCreated_ShouldReceiveCreatedOrderWithAllParams(t *testing.T) {
	order := Order{
		ID:    "123",
		Price: 10,
		Tax:   0.5,
	}

	assert.Equal(t, "123", order.ID)
	assert.Equal(t, 10.0, order.Price)
	assert.Equal(t, 0.5, order.Tax)
	assert.Nil(t, order.IsValid())
}

func TestGivenValidParams_WhenNewOrderIsCalled_ShouldReceiveCreatedOrderWithAllParams(t *testing.T) {
	order, err := NewOrder("123", 10.0, 0.5)

	assert.Equal(t, "123", order.ID)
	assert.Equal(t, 10.0, order.Price)
	assert.Equal(t, 0.5, order.Tax)
	assert.NoError(t, err, order)
}

func TestGivenPriceAndTax_WhenCalculatePriceIsCalled_ShouldSetFinalPrice(t *testing.T) {
	order, err := NewOrder("123", 10.0, 0.5)
	assert.NoError(t, err, order)
	assert.Nil(t, order.CalculateFinalPrice())
	assert.Equal(t, 10.5, order.FinalPrice)
}
