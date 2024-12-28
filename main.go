package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	//"errors"
)

type question struct {
	QuestionID   string `json:"questionid"`
	QuestionType string `json:"questionType"`
	QuestionText string `json:"questionText"`
	Answer       bool   `json:"answer"`
	AnswerDetail string `json:"answerDetail"`
	Source       string `json:"source"`
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

func getQuestions(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, questions)
}

func main() {
	router := gin.Default()
	router.GET("/questions", getQuestions)
	router.Run("localhost:8080")
}
