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

// type Comment struct {
// 	CommentID   string `json:"commentid"`
// 	CommentText string `json:"commentText"`
// }

// func main() {
// 	// Connect to MongoDB
// 	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
// 	if err != nil {
// 		panic(fmt.Errorf("Failed to connect to MongoDB: %w", err))
// 	}
// 	defer client.Disconnect(context.Background())

// 	// Get the database and collection
// 	db := client.Database("local")
// 	collection := db.Collection("comments")

// 	// Create a new Gin router
// 	r := gin.Default()

// 	// Define the route to retrieve a comment
// 	r.GET("/comments", func(c *gin.Context) {
// 		// Retrive all the comments from the request URL

// 		cursor, err := collection.Find(context.Background(), bson.M{})
// 		if err != nil {
// 			c.AbortWithStatus(http.StatusInternalServerError)
// 			return
// 		}

// 		defer cursor.Close(context.Background())

// 		// Convert the cursor to a slice of records
// 		var comments []Comment
// 		for cursor.Next(context.Background()) {
// 			var comment Comment
// 			if err := cursor.Decode(&comment); err != nil {
// 				c.AbortWithStatus(http.StatusInternalServerError)
// 				return
// 			}
// 			comments = append(comments, comment)
// 		}

// 		// Return the records as JSON
// 		c.JSON(http.StatusOK, comments)
// 	})

// 	// Start the server
// 	r.Run(":8080")
// }
