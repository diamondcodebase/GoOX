package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"errors"
)

type Question struct {
	QuestionID   string `json:"questionId"`
	QuestionType string `json:"questionType"`
	QuestionText string `json:"questionText"`
	Answer       bool   `json:"answer"`
	AnswerDetail string `json:"answerDetail"`
	Source       string `json:"source"`
}

type Comment struct {
	CommentID   string `json:"commentid"`
	CommentText string `json:"commentText"`
}

// data stored in memory and endpoint testing
var questions = []Question{
	{QuestionID: "1",
		QuestionText: "The genealogy of Jesus Christ in the Gospel of Matthew begins with Adam.",
		Answer:       false,
		AnswerDetail: "The genealogy of Jesus Christ in Matthew begins with Abraham (Matthew 1:1)",
	},
	{
		QuestionID:   "2",
		QuestionText: "Five women's names appear in the genealogy of Jesus in the Gospel of Matthew.",
		Answer:       true,
		AnswerDetail: "(Matthew 1:3,5,6,16)",
	},
	{
		QuestionID:   "3",
		QuestionText: "Matthew 1:23 'The virgin will conceive and give birth to a son, and they will call him Immanuel” comes from the Old Testament Psalms.",
		Answer:       false,
		AnswerDetail: "It comes from Isaiah 7 in the Old Testament (Isaiah 7:14)",
	},
	{
		QuestionID:   "4",
		QuestionText: "Matthew 3:3 'A voice of one calling in the wilderness, 'Prepare the way for the Lord, make straight paths for him.'' comes from Isaiah chapter 40.",
		Answer:       true,
		AnswerDetail: "(Matthew 3:3, Isaiah 40:3)",
	},
	{
		QuestionID:   "5",
		QuestionText: "Jesus was born in Bethlehem, Judah.",
		Answer:       true,
		AnswerDetail: "(Matthew 2:1)",
	},
}

var comments = []Comment{}

func getQuestions(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, questions)
}

func createComment(c *gin.Context) {
	var newComment Comment
	if err := c.BindJSON(&newComment); err != nil {
		return
	}

	comments = append(comments, newComment)
	c.IndentedJSON(http.StatusCreated, newComment)
}

func getComments(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, comments)
}

// data from MongoDB
// This function closes mongoDB connection and cancel context
func close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

// This is basic function connect mongoDB
// func connectDB(url string) *mongo.Client {
// 	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
// 	if err != nil {
// 		panic(fmt.Errorf("Failed to connect to MongoDB: %w", err))
// 	}
// 	return client
// }

// Function of getting client and collection according to url, dbname and collection name
func getMongoDBConnection(dbname string, collectname string) (*mongo.Client, *mongo.Collection, error) {
	// connect mongoDB
	var url = "mongodb://localhost:27017"
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to connect to MongoDB: %w", err)
	}

	// Get the database and collection
	db := client.Database(dbname)
	collection := db.Collection(collectname)
	return client, collection, nil
}

// Function of getting all comments
func getAllComments(c *gin.Context) {
	// Get the MongoDB client and collection
	client, collection, err := getMongoDBConnection("local", "comments")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Retrieve all records from MongoDB
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	// Convert the cursor to a slice of records
	var allComments []Comment
	for cursor.Next(context.Background()) {
		var comment Comment
		if err := cursor.Decode(&comment); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		allComments = append(allComments, comment)
	}

	c.IndentedJSON(http.StatusOK, allComments)
}

// Function of getting a specific comment according to comment id
func getCommentByCommentID(c *gin.Context) {
	// Get the comment ID from the request URL
	commentid := c.Query("commentid")

	// Get the MongoDB client and collection
	client, collection, err := getMongoDBConnection("local", "comments")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Retrieve the comment from MongoDB
	cursor, err := collection.Find(context.Background(), bson.M{"commentid": commentid})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var comments []Comment
	for cursor.Next(context.Background()) {
		var comment Comment
		if err := cursor.Decode(&comment); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}
	c.IndentedJSON(http.StatusOK, comments)
}

// query method returns a cursor and error.
func query(client *mongo.Client, ctx context.Context,
	dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {

	// select database and collection.
	collection := client.Database(dataBase).Collection(col)

	// collection has an method Find,
	// that returns a mongo.cursor
	// based on query and field.
	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return
}

// main function
func main() {
	// Create a new Gin router
	router := gin.Default()

	// Define the route to retrieve all records
	router.GET("/comments", getAllComments)
	router.GET("/comments/search", getCommentByCommentID)
	// Run the router
	router.Run("localhost:8080")
}
