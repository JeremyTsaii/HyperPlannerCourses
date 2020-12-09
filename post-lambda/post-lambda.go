package main

import (
	"context"
	"encoding/json"
	"net/url"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection Struct stores handle to collection used by Lambda
type Connection struct {
	collection *mongo.Collection
}

// Course Struct for course entry
type Course struct {
	Campus  string  `json:"campus"`
	Code    string  `json:"code"`
	Credits float64 `json:"credits"`
	Title   string  `json:"title"`
}

func (connection Connection) handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

}

func apiError(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 502}, err
}

func main() {
	// Open Connection
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("CONNECT_STR")))

	if err != nil {
		panic(err)
	}

	// close db connection after main returns
	defer client.Disconnect(ctx)

	// Store connection for connection pooling
	connection := Connection{
		collection: client.Database("HyperPlanner").Collection("Courses"),
	}

	// Invoke handleRequest with connection as receiver
	lambda.Start(connection.handleRequest)
}