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
	// Look at query parameters and create filter
	// 4 possible parameters: campus, code, credits, title
	filter := bson.M{}
	possible := [4]string{"campus", "code", "credits", "title"}
	for _, pos := range possible {
		newFilter, updateErr := updateFilter(request, filter, pos)
		if updateErr != nil {
			return apiError(updateErr)
		}
		filter = newFilter
	}

	// Get cursor, loop through all courses, append to results
	result, readErr := readCoursesFilter(connection.collection, filter)
	if readErr != nil {
		return apiError(readErr)
	}

	// Format body response
	tmap := make(map[string]interface{})
	tmap["length"] = len(result)
	tmap["courses"] = result

	body, jsonErr := json.Marshal(tmap)
	if jsonErr != nil {
		return apiError(jsonErr)
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

// Add parameter key and value to bson query
func updateFilter(request events.APIGatewayProxyRequest, filter bson.M, parameter string) (bson.M, error) {
	param, found := request.QueryStringParameters[parameter]
	if found {
		value, err := url.QueryUnescape(param)
		if err != nil {
			return nil, err
		}

		// Credits parameter must be turned back to float
		if parameter == "credits" {
			float, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			filter[parameter] = float
		} else {
			filter[parameter] = value
		}
	}
	return filter, nil
}

// Read all courses that match filter from db and append to results
func readCoursesFilter(courses *mongo.Collection, filter bson.M) ([]bson.M, error) {
	ctx := context.Background()
	cur, findErr := courses.Find(ctx, filter)
	if findErr != nil {
		return nil, findErr
	}

	// Decode all courses in cursor into result
	var result []bson.M
	if curErr := cur.All(ctx, &result); curErr != nil {
		return nil, curErr
	}

	cur.Close(ctx)
	return result, nil
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
