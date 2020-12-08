package main

import (
	"context"
	"os"
	"net/url"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Course is struct for a course entry
type Course struct {
	Campus  string  `json:"campus"`
	Code    string  `json:"code"`
	Credits float64 `json:"credits"`
	Title   string  `json:"title"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Open Connection
	uri := os.Getenv("CONNECT_STR")
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Connect to db
	db := client.Database("HyperPlanner")
	courses := db.Collection("Courses")

	// Look at query parameters and create filter
	// 4 possible parameters: campus, code, credits, title
	filter := bson.M{}
	possible := [4]string{"campus", "code", "credits", "title"}
	for _, pos := range possible {
		if err := updateFilter(request, filter, pos); err != nil {
			return apiError(err)
		}
	}
	
	// Get cursor, loop through all courses, append to results
	var result []Course
	length, err := readCoursesFilter(courses, result, filter)
	if err != nil {
		return apiError(err)
	}

	// Format body response
	tmap := make(map[string]interface{})
	tmap["length"] = length
	tmap["courses"] = result

	body, err := json.Marshal(tmap)
	if err != nil {
		return apiError(err)
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

// Add parameter key and value to bson query
func updateFilter(request events.APIGatewayProxyRequest, filter bson.M, parameter string) (error) {
	param, found := request.QueryStringParameters[parameter]
	if found {
		value, err := url.QueryUnescape(param)
		if nil != err {
			return err
		}
		filter[parameter] = value
	}
	return nil
}

// Read all courses that match filter from db and append to results
func readCoursesFilter(courses *mongo.Collection, result []Course, filter bson.M) (int, error) {
	cur, err := courses.Find(context.Background(), filter)
	if err != nil {
		return -1, err
	}

	// Loop through all courses in cursor and append to results
	length := 0
	for cur.Next(context.Background()) {
		var course Course
		if err := cur.Decode(&course); err != nil {
			return -1, err
		}
		length++
		result = append(result, course)
	}

	if err = cur.Err(); err != nil {
		return -1, err
	}

	cur.Close(context.Background())
	return length, nil
}

func apiError(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{}, err
}

func main() {
	lambda.Start(handleRequest)
}