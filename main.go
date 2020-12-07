package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

// Course is struct for a course entry
type Course struct {
	Campus  string  `json:"campus"`
	Code    string  `json:"code"`
	Credits float32 `json:"credits"`
	Title   string  `json:"title"`
}

func main() {
	// Open Connection
	uri := os.Getenv("CONNECT_STR")
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		panic(err)
	}

	// Ping the primary
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected and pinged.")

	// Connect to db
	db := client.Database("HyperPlanner")
	courses := db.Collection("Courses")

	// Insert course
	testCourse := Course{"hmc", "CSCI134", 3.0, "Operating Systems"}
	insertCourse(courses, testCourse)

	// Read course
	var course bson.M
	readCourse(courses, course)
}

func insertCourse(courses *mongo.Collection, course Course) {
	result, err := courses.InsertOne(
		context.Background(),
		course)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("API Result:", result.InsertedID)
	fmt.Println("Successfully inserted.")
}

func readCourse(courses *mongo.Collection, course bson.M) {
	err := courses.FindOne(context.Background(), bson.M{}).Decode(&course)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Found single document:", course)
	fmt.Println("Successfully read.")
}
