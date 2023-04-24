package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	pb "kv-grpc/server/kv"
)
func main(){
	conn, err := grpc.Dial("localhost:5051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewKVClient(conn)
	if _, err := client.Put(context.Background(), &pb.PutRequest{Key: "k", Value: "KTBFFH"}); err != nil {
		log.Fatal(err)
	}
	
}