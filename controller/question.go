package controller

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"my-zhihu/models"
	"net/http"
	"strconv"
)

func PostQuestion(c *gin.Context) {
	var question *models.Question
	err := json.NewDecoder(c.Request.Body).Decode(&question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		log.Println("controller.PostQuestion:", err)
		return
	}
	if question.Title == "" || len(question.TopicURLTokens) == 0 {
		c.JSON(200, gin.H{
			"success": false,
		})
	}
	uid := VisitorId(c)
	if uid == 0 {
		c.JSON(200, gin.H{
			"success": false,
		})
	}
	err = models.InsertQuestion(question, uid)
	if err != nil {
		c.JSON(200, gin.H{
			"success": false,
		})
	}
	c.JSON(200, gin.H{
		"success":    true,
		"questionId": question.ID,
	})
}

func FollowQuestion(c *gin.Context) {
	qid := c.Param("qid")
	if qid == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	uid := VisitorId(c)
	err := models.FollowQuestion(qid, uid)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		log.Println("controller.FollowQuestion:", err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}

func UnfollowQuestion(c *gin.Context) {
	qid := c.Param("qid")
	if qid == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	uid := VisitorId(c)
	err := models.UnfollowQuestion(qid, uid)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		log.Println("controller.UnfollowQuestion:", err)
		return
	}
	c.JSON(200, gin.H{
		"success": true,
	})
}

func QuestionFollowers(c *gin.Context) {
	qid := c.Param("qid")
	if qid == "" {
		c.JSON(http.StatusNotFound, nil)
		log.Println("controllers.QuestionFollowers")
		return
	}
	uid := VisitorId(c)
	page := &models.Page{
		Session: sessions.Default(c),
	}
	offset, _ := strconv.Atoi(c.Request.FormValue("offset"))
	followers := page.QuestionFollowers(qid, offset, uid)
	if followers == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	c.JSON(200, gin.H{
		"paging": page.Paging,
		"data":   followers,
	})
}
