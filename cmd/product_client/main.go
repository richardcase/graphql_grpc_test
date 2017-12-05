package main

import (
	"context"
	"log"

	pb "github.com/richardcase/graphql_grpc_test/product"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductsClient(conn)

	request := &pb.ProductsRequest{Id: 1234}
	r, err := c.GetProducts(context.Background(), request)
	if err != nil {
		log.Fatalf("could not get products: %v", err)
	}
	log.Printf("Products: %s", r.Products[0].Name)

}
