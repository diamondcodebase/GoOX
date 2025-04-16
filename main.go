package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"errors"
)

// need to execute "go get github.com/gin-contrib/cors" to overcome the cross origin block from frontend call API

// for config
type Config struct {
	Backend struct {
		Port string `json:"port"`
	}
}

func LoadConfiguration(filename string) (Config, error) {
	var config Config
	configFile, err := os.Open(filename)
	defer configFile.Close()

	if err != nil {
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}

type Question struct {
	QuestionID   int32  `json:"questionId"`
	QuestionType string `json:"questionType"`
	QuestionText string `json:"questionText"`
	Answer       bool   `json:"answer"`
	AnswerDetail string `json:"answerDetail"`
	Source       string `json:"source"`
	Script       string `json:"script"`
	Chapter      string `json:"chapter"`
	Verse        string `json:"verse"`
}

type Comment struct {
	CommentID   string `json:"commentid"`
	CommentText string `json:"commentText"`
}

// data stored in memory and endpoint testing
var questions = []Question{
	{
		QuestionID:   1,
		QuestionText: "The genealogy of Jesus Christ in the Gospel of Matthew begins with Adam.",
		Answer:       false,
		AnswerDetail: "The genealogy of Jesus Christ in Matthew begins with Abraham (Matthew 1:1)",
	},
	{
		QuestionID:   2,
		QuestionText: "Five women's names appear in the genealogy of Jesus in the Gospel of Matthew.",
		Answer:       true,
		AnswerDetail: "(Matthew 1:3,5,6,16)",
	},
	{
		QuestionID:   3,
		QuestionText: "Matthew 1:23 'The virgin will conceive and give birth to a son, and they will call him Immanuel” comes from the Old Testament Psalms.",
		Answer:       false,
		AnswerDetail: "It comes from Isaiah 7 in the Old Testament (Isaiah 7:14)",
	},
	{
		QuestionID:   4,
		QuestionText: "Matthew 3:3 'A voice of one calling in the wilderness, 'Prepare the way for the Lord, make straight paths for him.'' comes from Isaiah chapter 40.",
		Answer:       true,
		AnswerDetail: "(Matthew 3:3, Isaiah 40:3)",
	},
	{
		QuestionID:   5,
		QuestionText: "Jesus was born in Bethlehem, Judah.",
		Answer:       true,
		AnswerDetail: "(Matthew 2:1)",
	},
}

var testResult = []Question{
	{
		QuestionID:   1,
		QuestionText: "The first question.",
		Answer:       false,
		AnswerDetail: "The first answer detail",
	},
	{
		QuestionID:   2,
		QuestionText: "The second question.",
		Answer:       true,
		AnswerDetail: "The second answer detail",
	},
	{
		QuestionID:   3,
		QuestionText: "The third question.",
		Answer:       false,
		AnswerDetail: "The third answer detail",
	},
}

var comments = []Comment{}

func getQuestions(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, questions)
}

func getTestResult(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, testResult)
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
	var url = "mongodb://localhost:27017/"
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
	client, collection, err := getMongoDBConnection("ox", "comments")
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
	commentid := c.Query("id")

	// Get the MongoDB client and collection
	client, collection, err := getMongoDBConnection("ox", "comments")
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer client.Disconnect(context.Background())

	// Retrieve the comment from MongoDB
	cursor, err := collection.Find(context.Background(), bson.M{"commentid": commentid})
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
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

// Function of getting a specific question according to question id
func getQuestionByQuestionID(c *gin.Context) {
	// Get the comment ID from the request URL
	questionIdstr := c.Query("id")
	questionId, err := strconv.Atoi(questionIdstr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Get the MongoDB client and collection
	client, collection, err := getMongoDBConnection("ox", "bible_questions")
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer client.Disconnect(context.Background())

	// Retrieve the comment from MongoDB
	cursor, err := collection.Find(context.Background(), bson.M{"questionId": questionId})
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer cursor.Close(context.Background())

	var questions []Question
	for cursor.Next(context.Background()) {
		var question Question
		if err := cursor.Decode(&question); err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			fmt.Printf(err.Error())
			return
		}
		questions = append(questions, question)
	}
	c.IndentedJSON(http.StatusOK, questions)
}

// Function of finding multiple questions according to question ids
func getBibleQuestionSet(c *gin.Context) {
	// Get length integer from query parameter len
	lengthStr := c.Query("len")
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Get the MongoDB client and collection
	client, collection, err := getMongoDBConnection("ox", "bible_questions")
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer client.Disconnect(context.Background())

	// Get question collection size
	collectionSize, err := collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Generate a ramdom questionIds array
	questionIds := generateQuestionNoArray(length, int(collectionSize))
	fmt.Println("questionIds are ", questionIds)

	// Create the filter using $in operator
	filter := bson.M{"questionId": bson.M{"$in": questionIds}}

	// Retrieve an array of questions from MongoDB according to the random questionIds array
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer cursor.Close(context.Background())

	var questions []Question
	for cursor.Next(context.Background()) {
		var question Question
		if err := cursor.Decode(&question); err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			fmt.Printf(err.Error())
			return
		}
		questions = append(questions, question)
	}
	c.IndentedJSON(http.StatusOK, questions)
}

// Function of finding multiple questions about Canada according to question ids
func getCanadaQuestionSet(c *gin.Context) {
	// Get length integer from query parameter len
	lengthStr := c.Query("len")
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Get the MongoDB client and collection
	client, collection, err := getMongoDBConnection("ox", "canada_questions")
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer client.Disconnect(context.Background())

	// Get question collection size
	collectionSize, err := collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Generate a ramdom questionIds array
	questionIds := generateQuestionNoArray(length, int(collectionSize))
	fmt.Println("questionIds are ", questionIds)

	// Create the filter using $in operator
	filter := bson.M{"questionId": bson.M{"$in": questionIds}}

	// Retrieve an array of questions from MongoDB according to the random questionIds array
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer cursor.Close(context.Background())

	var questions []Question
	for cursor.Next(context.Background()) {
		var question Question
		if err := cursor.Decode(&question); err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			fmt.Printf(err.Error())
			return
		}
		questions = append(questions, question)
	}
	c.IndentedJSON(http.StatusOK, questions)
}

// Function of finding multiple questions about Hong Kong according to question ids
func getHongKongQuestionSet(c *gin.Context) {
	// Get length integer from query parameter len
	lengthStr := c.Query("len")
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Get the MongoDB client and collection
	client, collection, err := getMongoDBConnection("ox", "hongkong_questions")
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer client.Disconnect(context.Background())

	// Get question collection size
	collectionSize, err := collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Generate a ramdom questionIds array
	questionIds := generateQuestionNoArray(length, int(collectionSize))
	fmt.Println("questionIds are ", questionIds)

	// Create the filter using $in operator
	filter := bson.M{"questionId": bson.M{"$in": questionIds}}

	// Retrieve an array of questions from MongoDB according to the random questionIds array
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer cursor.Close(context.Background())

	var questions []Question
	for cursor.Next(context.Background()) {
		var question Question
		if err := cursor.Decode(&question); err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			fmt.Printf(err.Error())
			return
		}
		questions = append(questions, question)
	}
	c.IndentedJSON(http.StatusOK, questions)
}

// Function of generating a series of random numbers
func generateRandomNos(round int) {
	for i := 0; i < round; i++ {
		randomNumber := rand.Intn(5)
		fmt.Println(randomNumber)
	}
}

// Functio of generating an array of random numbers with 5 distinct integers
func generateQuestionNoArray(round int, maxNo int) []int {
	arr := rand.Perm(maxNo)

	// every integers in the array add 1 to avoid zero value id
	for i := range arr {
		arr[i]++
	}

	sl := arr[:round]
	fmt.Println("slice equals to ", sl)
	return sl
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
	// Test to load config file
	// config, _ := LoadConfiguration("config.json")
	// var port = ":" + config.Backend.Port
	// fmt.Println(port)

	// Create a new Gin router
	router := gin.Default()
	// Apply middleware to overcome CORS policy during API call from frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://main.d2q5e0gpmjipai.amplifyapp.com", "https://main.d2q5e0gpmjipai.amplifyapp.com:3000", "https://diamondbackend.click"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	// Test random number function
	// generateRandomNos(20)
	// Test generate random array function
	// generateQuestionNoArray(5, 200)

	// Define the route to retrieve all records
	router.GET("/comments", getAllComments)
	router.GET("/comment", getCommentByCommentID)               // test by cmd: curl localhost:8080/comment?id=22749003
	router.GET("/test", getTestResult)                          // test by cmd: curl localhost:8080/test
	router.GET("/question", getQuestionByQuestionID)            // test by cmd: curl localhost:8080/question?id=12
	router.GET("/questionset/bible", getBibleQuestionSet)       // test by cmd: curl localhost:8080/questionset/bible?len=5
	router.GET("/questionset/canada", getCanadaQuestionSet)     // test by cmd: curl localhost:8080/questionset/canada?len=5
	router.GET("/questionset/hongkong", getHongKongQuestionSet) // test by cmd: curl localhost:8080/questionset/hongkong?len=5
	// Run the server on localhost
	// router.Run("localhost:8080")

	// Azure App Service sets the port as an Environmental
	// This can be random, so needs to be loaed at start
	// port := os.Getenv("HTTP_PLATFORM_PORT")

	// // default back to 8080 for local dev
	// if port == "" {
	// 	port = "8080"
	// }

	port := "8080"

	// Run the server on a domain
	router.Run("localhost:" + port)
}
