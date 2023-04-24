package main

import (
	"context"
	pb "kv-grpc/server/kv"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (kv *KV) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error){
	key := req.Key
	kv.mu.Lock()
	defer kv.mu.Unlock()
	val, ok := kv.data[key]
	if ok {
		return &pb.GetResponse{Value: val}, nil
	}
	return nil, status.Errorf(codes.NotFound, "key not found")
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
func Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error){
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Fatalf("error fetching metadata: %s", err)
		}
		log.Println(md["key"][0])
		log.Printf("%s: %s", info.FullMethod, time.Now())
		return handler(ctx, req)
	}
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

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(Unary()))
	pb.RegisterKVServer(grpcServer, NewKVServer())
	log.Println("Server started")
	grpcServer.Serve(l)

}
