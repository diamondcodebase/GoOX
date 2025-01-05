// package main

// import (
// 	"context"
// 	"fmt"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// type Record struct {
// 	ID    string `bson:"_id"`
// 	Name  string `bson:"name"`
// 	Age   int    `bson:"age"`
// 	Email string `bson:"email"`
// }

// func getRecords(c *gin.Context) {
// 	// Get the MongoDB client and collection
// 	client, collection := getMongoClient()
// 	defer client.Disconnect(context.Background())

// 	// Retrieve all records from MongoDB
// 	cursor, err := collection.Find(context.Background(), bson.M{})
// 	if err != nil {
// 		c.AbortWithStatus(http.StatusInternalServerError)
// 		return
// 	}
// 	defer cursor.Close(context.Background())

// 	// Convert the cursor to a slice of records
// 	var records []Record
// 	for cursor.Next(context.Background()) {
// 		var record Record
// 		if err := cursor.Decode(&record); err != nil {
// 			c.AbortWithStatus(http.StatusInternalServerError)
// 			return
// 		}
// 		records = append(records, record)
// 	}

// 	// Return the records as JSON
// 	c.JSON(http.StatusOK, records)
// }

// func getMongoClient(dbName string, collectionName string) (*mongo.Client, *mongo.Collection) {
// 	// Connect to MongoDB
// 	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
// 	if err != nil {
// 		panic(fmt.Errorf("Failed to connect to MongoDB: %w", err))
// 	}

// 	// Get the database and collection
// 	db := client.Database(dbName)
// 	collection := db.Collection(collectionName)

// 	return client, collection
// }

// func main() {
// 	// Create a new Gin router
// 	r := gin.Default()

// 	// Define the route to retrieve all records
// 	r.GET("/records", getRecords)

// 	// Start the server
// 	r.Run(":8080")
// }
