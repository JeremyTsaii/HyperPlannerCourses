package main

import (
	"context"
	"os"

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

func (connection Connection) handleRequest(ctx context.Context, request map[string]interface{}) (interface{}, error) {
	// Query/mutation field and arguments
	field := request["field"].(string)

	switch field {
	case "getCourses":
		params := request["arguments"].(map[string]interface{})
		filter := bson.M{}

		possible := [4]string{"campus", "code", "credits", "title"}
		for _, pos := range possible {
			newFilter, updateErr := updateFilter(params, filter, pos)
			if updateErr != nil {
				return []Course{}, updateErr
			}
			filter = newFilter
		}

		// Get cursor, loop through all courses, append to results
		courses, readErr := readCoursesFilter(connection.collection, filter)
		if readErr != nil {
			return []Course{}, readErr
		}

		return courses, nil
	case "addCourses":
		args := request["arguments"].(map[string]interface{})
		courses := args["courses"].([]interface{})

		// Insert into db
		_, insertErr := connection.collection.InsertMany(ctx, courses)
		if insertErr != nil {
			return []Course{}, insertErr
		}

		return courses, nil
	}

	return []Course{}, nil
}

// Add parameter key and value to bson query
func updateFilter(request map[string]interface{}, filter bson.M, parameter string) (bson.M, error) {
	value, found := request[parameter]
	if found {
		filter[parameter] = value
	}
	return filter, nil
}

// Read all courses that match filter from db and append to results
func readCoursesFilter(collection *mongo.Collection, filter bson.M) ([]bson.M, error) {
	ctx := context.Background()
	cur, findErr := collection.Find(ctx, filter)
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
