package main

import (
	"context"
	"fmt"
	pb "grpccrud/proto"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func main() {
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	pb.RegisterDataserviceServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

func (s *server) Create(ctx context.Context, request *pb.CreateRequest) (*pb.CreateResponse, error) {
	a := request.GetApi()
	b := request.GetToDo()
	fmt.Println(b)
	fmt.Println(a)

	return &pb.CreateResponse{Api: "v2", Id: 1}, nil
}

func (s *server) Read(ctx context.Context, request *pb.ReadRequest) (*pb.ReadResponse, error) {
	a := request.GetApi()
	b := request.GetId()

	fmt.Println(b)
	return &pb.ReadResponse{Api: a, ToDo: &pb.ToDo{}}, nil
}

func (s *server) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	a, b := request.GetApi(), request.GetToDo()
	fmt.Println(b)
	return &pb.UpdateResponse{Api: a, Updated: 1}, nil
}

func (s *server) Delete(ctx context.Context, request *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	a := request.GetApi()
	b := request.GetId()

	result := fmt.Sprint(a, b)

	return &pb.DeleteResponse{Api: result}, nil
}

func (s *server) ReadAll(ctx context.Context, request *pb.ReadAllRequest) (*pb.ReadAllResponse, error) {
	a := request.GetApi()

	result := a

	return &pb.ReadAllResponse{Api: result, ToDos: []*pb.ToDo{}}, nil
}

type DBConn struct {
	Username string
	Password string
	Addr     string
	Name     string
}

func (d *DBConn) Format() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s",
		d.Username,
		d.Password,
		d.Addr,
		d.Name)
}
