package main

import (
	"net/http"

	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"errors"
)

type question struct {
	QuestionID   string `json:"questionId"`
	QuestionType string `json:"questionType"`
	QuestionText string `json:"questionText"`
	Answer       bool   `json:"answer"`
	AnswerDetail string `json:"answerDetail"`
	Source       string `json:"source"`
}

type comment struct {
	CommentID   string `json:"commentid"`
	CommentText string `json:"commentText"`
}

var questions = []question{
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

var comments = []comment{}

func getQuestions(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, questions)
}

func createComment(c *gin.Context) {
	var newComment comment
	if err := c.BindJSON(&newComment); err != nil {
		return
	}

	comments = append(comments, newComment)
	c.IndentedJSON(http.StatusCreated, newComment)
}

func getComments(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, comments)
}

func Connect() *mongo.Collection {
	// Find .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Get value from .env
	MONGO_URI := os.Getenv("MONGO_URI")

	// Connect to the database.
	clientOption := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create collection
	collection := client.Database("testdb").Collection("test")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to db")

	return collection
}

const uri = "mongodb://localhost:27017/test"

func dbconnect_test() {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func main() {
	dbconnect_test()
	router := gin.Default()
	router.GET("/questions", getQuestions)
	router.GET("/comments", getComments)
	router.POST("/comment", createComment)
	router.Run("localhost:8080")
}
