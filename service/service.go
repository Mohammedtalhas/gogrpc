package main

import (
	"context"
	"fmt"
	pb "grpccrud/proto"
	"grpccrud/service/Config"
	"net"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func main() {
	var errs error
	Config.DB, errs = gorm.Open("mysql", Config.DbURL(Config.BuildDBConfig()))
	if errs != nil {
		fmt.Println("Status:", errs)
	}
	defer Config.DB.Close()
	Config.DB.AutoMigrate(&pb.ToDo{})
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
	var err error
	if a == apiVersion {

		if err = Config.DB.Create(b).Error; err != nil {
			return &pb.CreateResponse{Api: a, Id: 0}, err
		}
	}

	return &pb.CreateResponse{Api: "v2", Id: b.Id}, nil
}

func (s *server) Read(ctx context.Context, request *pb.ReadRequest) (*pb.ReadResponse, error) {
	a := request.GetApi()
	b := request.GetId()
	var todo = &pb.ToDo{}
	var err error
	if a == apiVersion {

		if err = GetDataByID(todo, b); err != nil {
			return nil, err
		}
	}
	return &pb.ReadResponse{Api: a, ToDo: todo}, nil
}
func GetDataByID(data *pb.ToDo, id int64) (err error) {
	if err = Config.DB.Where("id = ?", id).First(data).Error; err != nil {
		return err
	}
	return nil
}

func (s *server) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	a, b := request.GetApi(), request.GetToDo()
	var todo = &pb.ToDo{}
	if a == apiVersion {
		err := GetDataByID(todo, b.Id)
		if err != nil {
			return nil, err
		}
		Config.DB.Save(b)
	}

	return &pb.UpdateResponse{Api: a, Updated: b.Id}, nil
}

func (s *server) Delete(ctx context.Context, request *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	a := request.GetApi()
	b := request.GetId()
	var todo = &pb.ToDo{}
	if a == apiVersion {
		Config.DB.Where("ID = ?", b).Delete(todo)

	}
	return &pb.DeleteResponse{Api: a, Deleted: b}, nil
}

func (s *server) ReadAll(ctx context.Context, request *pb.ReadAllRequest) (*pb.ReadAllResponse, error) {
	a := request.GetApi()
	todo := []pb.ToDo{}
	ToDos := make([]*pb.ToDo, 0)
	var err error
	if a == apiVersion {

		if err = Config.DB.Find(&todo).Error; err != nil {
			return nil, err
		}
	}
	for _, u := range todo {
		ToDos = append(ToDos,
			&pb.ToDo{
				Id: u.Id, Title: u.Title, Description: u.Description, Reminder: u.Reminder,
			})
	}
	return &pb.ReadAllResponse{Api: a, ToDos: ToDos}, nil
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
