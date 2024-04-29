package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"my-zhihu/models"
	"net/http"
)

func SearchTopics(c *gin.Context) {
	var topics []models.Topic
	token := c.Request.FormValue("token")
	if token == "" {
		c.JSON(200, topics)
	}
	topics = models.SearchTopics(token)
	c.JSON(200, topics)
}

func PostTopic(c *gin.Context) {
	var topic models.Topic
	var err error
	err = json.NewDecoder(c.Request.Body).Decode(&topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	err = models.UpdateTopic(&topic)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"topic":   topic,
	})
}
