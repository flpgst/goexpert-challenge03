package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/flpgst/golang-studies/55-CleanArch/configs"
	"github.com/flpgst/golang-studies/55-CleanArch/internal/event/handler"
	"github.com/flpgst/golang-studies/55-CleanArch/internal/infra/graph"
	"github.com/flpgst/golang-studies/55-CleanArch/internal/infra/grpc/pb"
	"github.com/flpgst/golang-studies/55-CleanArch/internal/infra/grpc/service"
	"github.com/flpgst/golang-studies/55-CleanArch/pkg/events"
	"github.com/go-chi/chi/v5"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(configs)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrderUseCase(db)

	// webserver

	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	router := chi.NewRouter()
	router.Route("/order", func(r chi.Router) {
		r.Post("/", webOrderHandler.Create)
		r.Get("/", webOrderHandler.List)
	})
	// webserver := webserver.NewWebServer(configs.WebServerPort)
	// webserver.AddHandler("/order", webOrderHandler.Create)
	// webserver.AddHandler("/order", webOrderHandler.List)
	fmt.Println("Listening on port", configs.WebServerPort)
	go http.ListenAndServe(":"+configs.WebServerPort, router)

	// grpc
	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)

	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Listening GRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	// graphql
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   listOrderUseCase,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", configs.GraphQLServerPort)
	log.Fatal(http.ListenAndServe(":"+configs.GraphQLServerPort, nil))

}

func getRabbitMQChannel(configs *configs.Conf) *amqp.Channel {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", configs.RABBITMQ_USER, configs.RABBITMQ_PASSWORD, configs.RABBITMQ_HOST, configs.RABBITMQ_PORT))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
