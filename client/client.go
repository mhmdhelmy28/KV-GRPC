package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "kv-grpc/server/kv"
)
func Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Printf("unary: %s", method)

		return invoker(metadata.AppendToOutgoingContext(ctx, "key", "val"), method, req, reply, cc, opts...)
		
	}
}
func main(){
	conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure(), grpc.WithUnaryInterceptor(Unary()),)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewKVClient(conn)
	if _, err := client.Put(context.Background(), &pb.PutRequest{Key: "k", Value: "KTBFFH"}); err != nil {
		log.Fatal(err)
	}
	
}