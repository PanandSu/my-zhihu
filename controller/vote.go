package controller

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"my-zhihu/models"
	"net/http"
	"strconv"
)

type VoteType struct {
	Type string
}

func VoteAnswer(c *gin.Context) {
	var vote VoteType
	var flag bool
	aid := c.Param("aid")
	if aid == "" {
		c.JSON(200, gin.H{
			"success": flag,
		})
		return
	}
	user, _ := Visitor(c)
	if user == nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": flag,
		})
		return
	}
	err := json.NewDecoder(c.Request.Body).Decode(&vote)
	if err != nil {
		c.JSON(200, gin.H{
			"success": flag,
		})
	}
	switch vote.Type {
	case "up":
		flag = user.VoteUp(aid)
	case "down":
		flag = user.VoteDown(aid)
	case "neutral":
		flag = user.Neutral(aid)
	default:
	}
	c.JSON(200, gin.H{
		"success": flag,
	})
}

func AnswerVoters(c *gin.Context) {
	aid := c.Param("aid")
	if aid == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	uid := VisitorId(c)
	page := &models.Page{
		Session: sessions.Default(c),
	}
	offset, _ := strconv.Atoi(c.Request.FormValue("offset"))
	voters := page.AnswerVoters(aid, offset, uid)
	c.JSON(200, gin.H{
		"paging": page.Paging,
		"data":   voters,
	})
}
