package main

import (
	"context"
	"fmt"
	"encoding/json"
	"log"
	"os"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// Connect to db
	db := client.Database("HyperPlanner")
	courses := db.Collection("Courses")

	// Insert course
	// testCourse := Course{"hmc", "CSCI134", 3.0, "Operating Systems"}
	// insertCourse(courses, testCourse)

	// Read course
	// var course bson.M
	// readCourse(courses, course)
	var allCourses []Course
	allFilter := bson.M{"campus": "hmc"}
	readCoursesFilter(courses, allCourses, allFilter)

	var codeCourses []Course
	codeFilter := bson.M{"code": "CSCI134"}
	readCoursesFilter(courses, codeCourses, codeFilter)

	var titleCourses []Course
	titleFilter := bson.M{"title": "Software Development"}
	readCoursesFilter(courses, titleCourses, titleFilter)

	var creditCourses []Course
	creditFilter := bson.M{"credits": 3.0}
	readCoursesFilter(courses, creditCourses, creditFilter)

}

// Insert given course into db
func insertCourse(courses *mongo.Collection, course Course) {
	result, err := courses.InsertOne(
		context.Background(),
		course)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("API Result:", result.InsertedID)
}

// Read single course from db
func readCourse(courses *mongo.Collection, course bson.M, filter bson.M) {
	err := courses.FindOne(context.Background(), filter).Decode(&course)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Found single document:", course)
}

// Read all courses from db and output as json string
func readCoursesFilter(courses *mongo.Collection, results []Course, filter bson.M) {
	cursor, err := courses.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through all courses in cursor and append to results
	total := 0
	for cursor.Next(context.Background()) {
		var course Course
		if err := cursor.Decode(&course); err != nil {
			log.Fatal(err)
		}
		total++
		results = append(results, course)
	}

	if err = cursor.Err(); err != nil {
		log.Fatal(err)
	}

	cursor.Close(context.Background())

	ret, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Found documents: ", total)
	fmt.Println(string(ret))
}
