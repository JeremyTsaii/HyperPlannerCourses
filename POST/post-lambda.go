package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

// CoursesArr Struct for list of courses in POST body
type CoursesArr struct {
	Courses []Course `json:"courses"`
}

func (connection Connection) handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := request.Body

	// Insert courses specified in body of request
	length, insertErr := insertCourses(connection.collection, body)
	if insertErr != nil {
		return apiError(insertErr)
	}

	// Format body response
	tmap := make(map[string]interface{})
	tmap["length"] = length

	resBody, jsonErr := json.Marshal(tmap)
	if jsonErr != nil {
		return apiError(jsonErr)
	}

	return events.APIGatewayProxyResponse{Body: string(resBody), StatusCode: 200}, nil
}

func insertCourses(collection *mongo.Collection, rawStr string) (int, error) {
	ctx := context.Background()

	// Empty POST body
	if rawStr == "" {
		return 0, nil
	}

	// Parse json string into usable json
	var courses CoursesArr
	decodeErr := json.Unmarshal([]byte(rawStr), &courses)
	if decodeErr != nil {
		return -1, decodeErr
	}

	// Convert to format necessary for insertion
	var iCourses []interface{}
	for _, course := range courses.Courses {
		iCourses = append(iCourses, course)
	}

	// Insert into db
	_, insertErr := collection.InsertMany(ctx, iCourses)
	if insertErr != nil {
		return -1, insertErr
	}

	return len(courses.Courses), nil
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
