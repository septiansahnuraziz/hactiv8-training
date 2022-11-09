package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"todo-grpc/common/model"
	"todo-grpc/config"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TodoServer struct {
	model.UnimplementedTodosServer
}

var Todos = make(map[string]*model.Todo)

func main() {
	server := grpc.NewServer()
	userServer := new(TodoServer)
	model.RegisterTodosServer(server, userServer)

	log.Println("Starting Todo Server At ", config.SERVICE_TODO_PORT)

	listener, err := net.Listen("tcp", "localhost"+config.SERVICE_TODO_PORT)
	if err != nil {
		log.Fatalf("could not listen. Err: %+v\n", err)
	}

	go func() {
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		grpcServerEndpoint := flag.String("grpc-server-endpoint", "localhost"+config.SERVICE_TODO_PORT, "gRPC server endpoint")
		_ = model.RegisterTodosHandlerFromEndpoint(context.Background(), mux, *grpcServerEndpoint, opts)
		log.Println("Starting Todo Server HTTP at 9001")
		http.ListenAndServe(":9001", mux)
	}()

	log.Fatalln(server.Serve(listener))
}

func (t *TodoServer) CreateTodo(ctx context.Context, req *model.Todo) (*model.MutationResponse, error) {
	Todos[req.GetId()] = &model.Todo{
		Id:   req.GetId(),
		Name: req.GetName(),
		Todo: req.GetTodo(),
	}

	message := fmt.Sprintf("%v successfully appended", req.GetId())

	return &model.MutationResponse{Success: message}, nil
}

func (t *TodoServer) GetAll(ctx context.Context, req *emptypb.Empty) (*model.GetAllResponse, error) {
	var todos []*model.Todo
	for _, v := range Todos {
		todos = append(todos, &model.Todo{
			Id:   v.GetId(),
			Name: v.GetName(),
			Todo: v.GetTodo(),
		})
	}

	return &model.GetAllResponse{Data: todos}, nil
}

func (t *TodoServer) Update(ctx context.Context, req *model.UpdateRequest) (*model.MutationResponse, error) {
	Todos[req.GetId()] = &model.Todo{
		Id:   req.GetId(),
		Name: req.GetName(),
		Todo: req.GetTodo(),
	}
	message := req.GetId() + "successfully appended"

	return &model.MutationResponse{Success: message}, nil
}

func (t *TodoServer) Delete(ctx context.Context, req *model.DeleteRequest) (*model.MutationResponse, error) {
	delete(Todos, req.GetId())
	message := req.GetId() + "successfully deleted"

	return &model.MutationResponse{Success: message}, nil
}
