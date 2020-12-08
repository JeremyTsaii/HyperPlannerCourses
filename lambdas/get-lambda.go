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
		filter, err = updateFilter(request, filter, pos)
		if err != nil {
			return apiError(err)
		}
	}

	// Get cursor, loop through all courses, append to results
	result, err := readCoursesFilter(courses, filter)
	if err != nil {
		return apiError(err)
	}

	// Format body response
	tmap := make(map[string]interface{})
	tmap["length"] = len(result)
	tmap["courses"] = result

	body, err := json.Marshal(tmap)
	if err != nil {
		return apiError(err)
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

// Add parameter key and value to bson query
func updateFilter(request events.APIGatewayProxyRequest, filter bson.M, parameter string) (bson.M, error) {
	param, found := request.QueryStringParameters[parameter]
	if found {
		value, err := url.QueryUnescape(param)
		if nil != err {
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
func readCoursesFilter(courses *mongo.Collection, filter bson.M) ([]Course, error) {
	var result []Course

	cur, err := courses.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	// Loop through all courses in cursor and append to results
	for cur.Next(context.Background()) {
		var course Course
		if err := cur.Decode(&course); err != nil {
			return nil, err
		}
		result = append(result, course)
	}

	if err = cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(context.Background())
	return result, nil
}

func apiError(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{}, err
}

func main() {
	lambda.Start(handleRequest)
}
