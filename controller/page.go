package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"my-zhihu/models"
	"net/http"
)

type DataState struct {
	Page     string           `json:"page"`
	Question *models.Question `json:"question"`
	Answers  []*models.Answer `json:"answers"`
	TopStory []*models.Action `json:"top_story"`
}

func AnswerGet(c *gin.Context) {
	aid := c.Param("aid")
	user, uid := Visitor(c)
	answer := models.GetAnswer(aid, uid)
	if answer == nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	v := DataState{
		Page:     "answer",
		Question: answer.Question,
	}
	v.Answers = append(v.Answers, answer)
	dataState, _ := json.Marshal(&v)
	c.HTML(http.StatusOK, "answer.html", gin.H{
		"answer":    answer,
		"question":  answer.Question,
		"user":      user,
		"dataState": string(dataState),
	})
}

func QuestionGet(c *gin.Context) {
	user, uid := Visitor(c)
	qid := c.Param("qid")
	if qid == "" {
		c.HTML(http.StatusNotFound, "404.html", nil)
		log.Println("controllers.QuestionGet():no question id")
		return
	}
	question := models.GetQuestionWithAnswers(qid, uid)
	if question == nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	v := DataState{
		Page:     "question",
		Question: question,
	}
	for _, answer := range question.Answers {
		v.Answers = append(v.Answers, answer)
	}
	dataState, _ := json.Marshal(&v)
	c.HTML(http.StatusOK, "question.html", gin.H{
		"question":  question,
		"user":      user,
		"dataState": string(dataState),
	})
}
