package graphql

import (
	"context"
	"strconv"

	gql "github.com/neelance/graphql-go"
	pb "github.com/richardcase/graphql_grpc_test/product"
	"google.golang.org/grpc"
)

var Schema = `
	schema {
		query: Query
	}
	type Query {
		products(id: ID): [Product]!
	}
	type Product {
		id: ID!
		name: String!
	}
`

type product struct {
	ID   gql.ID
	Name string
}

type Resolver struct{}

func (r *Resolver) Products(args struct{ ID *gql.ID }) []*productResolver {
	// Create fake products
	/*fakeProduct := &product{
		ID:   "1234",
		Name: "Product1",
	}
	products := []*product{fakeProduct}*/
	products, err := queryProductService(1)
	if err != nil {
		panic(err)
	}

	var l []*productResolver
	for _, prod := range products {
		l = append(l, &productResolver{prod})
	}
	return l
}

type productResolver struct {
	p *product
}

func (r *productResolver) ID() gql.ID {
	return r.p.ID
}

func (r *productResolver) Name() string {
	return r.p.Name
}

func queryProductService(id int64) ([]*product, error) {
	address := "localhost:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := pb.NewProductsClient(conn)

	request := &pb.ProductsRequest{Id: id}
	response, err := c.GetProducts(context.Background(), request)
	if err != nil {
		return nil, err
	}

	products := []*product{}
	for _, prod := range response.Products {
		strId := strconv.FormatInt(prod.Id, 10)
		products = append(products, &product{
			ID:   gql.ID(strId),
			Name: prod.Name,
		})
	}

	return products, nil
}
