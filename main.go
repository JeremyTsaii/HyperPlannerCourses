package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Course is struct for a course entry
type Course struct {
	campus string
	code string
	credits float32
	title string
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
	result, err := courses.InsertOne(
		context.Background(),
		bson.D{
			{"campus", "hmc"},
			{"code", "CSCI121"},
			{"credits", 3.0},
			{"title", "Software Development"},
		})
	
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("API Result:", result);
	fmt.Println("Successfully inserted.");

	// Read course
	var course bson.M 
	if err = courses.FindOne(context.Background(), bson.M{}).Decode(&course); err != nil {
		log.Fatal(err);
	}
	fmt.Println(course);
	fmt.Println("Successfully read.")
}
