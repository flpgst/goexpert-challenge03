//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/flpgst/golang-studies/55-CleanArch/internal/entity"
	"github.com/flpgst/golang-studies/55-CleanArch/internal/event"
	"github.com/flpgst/golang-studies/55-CleanArch/internal/infra/database"
	"github.com/flpgst/golang-studies/55-CleanArch/internal/infra/web"
	"github.com/flpgst/golang-studies/55-CleanArch/internal/usecase"
	"github.com/flpgst/golang-studies/55-CleanArch/pkg/events"
	"github.com/google/wire"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

// var setEventDispatcherDependency = wire.NewSet(
// 	events.NewEventDispatcher,
// 	event.NewOrderCreated,
// 	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
// 	wire.Bind(new(events.EventDispacherInterface), new(*events.EventDispatcher)),
// )

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispacherInterface) *usecase.CreateOrderUseCase {
	wire.Build(setOrderRepositoryDependency, setOrderCreatedEvent, usecase.NewCreateOrderUseCase)
	return &usecase.CreateOrderUseCase{}
}

func NewListOrderUseCase(db *sql.DB) *usecase.ListOrderUseCase {
	wire.Build(setOrderRepositoryDependency, usecase.NewListOrderUseCase)
	return &usecase.ListOrderUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispacherInterface) *web.WebOrderHandler {
	wire.Build(setOrderRepositoryDependency, setOrderCreatedEvent, web.NewWebOrderHandler)
	return &web.WebOrderHandler{}
}
