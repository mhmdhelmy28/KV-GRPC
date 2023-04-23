package main

import (
	"context"
	pb "kv-grpc/server/kv"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

)
var (
)
type KV struct {
	pb.UnimplementedKVServer
	data map[string]string
	mu   sync.Mutex
}

func (kv *KV) 	Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error){
	key := req.Key
	value := req.Value
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[key] = value
	return &pb.PutResponse{}, nil

}
type Err struct {
	error
	msg string
}
func (kv *KV) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error){
	key := req.Key
	kv.mu.Lock()
	defer kv.mu.Unlock()
	val, ok := kv.data[key]
	if ok {
		return &pb.GetResponse{Value: val}, nil
	}
	return nil, Err{msg: "not found"}
}

func (kv *KV) List (req *pb.ListRequest, stream pb.KV_ListServer) error{
	kv.mu.Lock()
	defer kv.mu.Unlock()
	for k,v := range kv.data {
		err := stream.Send(&pb.ListResponse{Key: k ,Value: v})
		if err != nil {
			return err
		}
	}
	return nil

}

func NewKVServer() *KV{
	return &KV{
		data: map[string]string{},
		mu: sync.Mutex{},
	}
}
func main(){
	l, err := net.Listen("tcp", "localhost:5051")
	if err != nil {
		log.Fatal("could not listen on port")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterKVServer(grpcServer, NewKVServer())
	log.Println("Server started")
	grpcServer.Serve(l)

}
