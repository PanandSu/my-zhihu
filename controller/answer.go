package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"my-zhihu/models"
	"net/http"
)

func PostAnswer(c *gin.Context) {
	qid := c.Param("qid")
	if qid == "" {
		c.JSON(http.StatusNotFound, nil)
		log.Println("controllers.PostAnswer(): no question id")
		return
	}
	uid := VisitorId(c)
	content := struct {
		Content string `json:"content"`
	}{}
	err := json.NewDecoder(c.Request.Body).Decode(&content)
	if err != nil {
		log.Println("controller.PostAnswer(): ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	aid, err := models.InsertAnswer(qid, content.Content, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": aid})
}

func DeleteAnswer(c *gin.Context) {
	aid := c.Param("aid")
	if aid == "" {
		c.JSON(http.StatusNotFound, nil)
		log.Println("controllers.DeleteAnswer(): no answer id")
		return
	}
	uid := VisitorId(c)
	err := models.DeleteAnswer(aid, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(200, gin.H{"success": true})
}

func RestoreAnswer(c *gin.Context) {
	aid := c.Param("aid")
	if aid == "" {
		c.JSON(http.StatusNotFound, gin.H{})
		log.Println("controllers.RestoreAnswer(): no answer id")
		return
	}
	uid := VisitorId(c)
	err := models.RestoreAnswer(aid, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(200, gin.H{"success": true})
}
