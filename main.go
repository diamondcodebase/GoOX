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

type question struct {
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

// func getQueryComment(c *gin.Context) {
// 	// convert bson.M to Struct
// 	var queryComment comment
// 	var m = c.IndentedJSON(http.StatusOK, comments)
// }

// This function closes mongoDB connection and cancel context
func close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
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

// get method for endpoint
func getData() (comments []bson.D) {
	// Get Client, Context, CancelFunc and err from connect method.
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	// Free the resource when main function is returned
	defer close(client, ctx, cancel)

	// create a filter an option of type interface,
	// that stores bjson objects.
	var filter, option interface{}

	// filter  gets all document,
	// with maths field greater that 70
	filter = bson.D{{"commentid", "22749003"}}

	//  option remove id field from all documents
	option = bson.D{{"_id", 0}}

	// call the query method with client, context, database name, collection  name, filter and option
	cursor, err := query(client, ctx, "local",
		"comments", filter, option)
	// handle the errors.
	if err != nil {
		panic(err)
	}

	var results []bson.D

	if err := cursor.All(ctx, &results); err != nil {
		// handle the error
		panic(err)
	}

	// printing the result of query.
	fmt.Println("Query Result")
	for _, doc := range results {
		fmt.Println(doc)
	}
	return results
}

// main function
func main() {
	getData()
	// router := gin.Default()
	// router.GET("/questions", getQuestions)
	// router.GET("/comments", getComments)
	// router.POST("/comment", createComment)
	// router.Run("localhost:8080")
}
