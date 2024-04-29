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

func QuestionComments(c *gin.Context) {
	qid := c.Param("qid")
	if qid == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	uid := VisitorId(c)
	page := &models.Page{
		Session: sessions.Default(c),
	}
	offset, _ := strconv.Atoi(c.Request.FormValue("offset"))
	comments := page.QuestionComments(qid, offset, uid)
	if comments == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	c.JSON(200, gin.H{
		"paging": page.Paging,
		"data":   comments,
	})
}

func PostQuestionComment(c *gin.Context) {
	qid := c.Param("qid")
	if qid == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	uid := VisitorId(c)
	content := struct {
		Content string `json:"content"`
	}{}
	err := json.NewDecoder(c.Request.Body).Decode(&content)
	if err != nil {
		log.Println("controllers.PostQuestionComment(): ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	comment, err := models.InsertQuestionComment(qid, content.Content, uid)
	if err != nil {
		log.Println("controllers.PostQuestionComment(): ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, comment)
}

func DeleteQuestionComment(c *gin.Context) {
	qid := c.Param("id")
	cid := c.Param("cid")
	if qid == "" || cid == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	uid := VisitorId(c)
	err := models.DeleteQuestionComment(qid, cid, uid)
	if err != nil {
		log.Println("controllers.DeleteQuestionComment(): ", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func LikeQuestionComment(c *gin.Context) {
	cid := c.Param("cid")
	if cid == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	uid := VisitorId(c)
	err := models.LikeQuestionComment(cid, uid)
	if err != nil {
		log.Println("controllers.LikeQuestionComment(): ", err)
		if err.Error() == "reply is zero" {
			c.JSON(http.StatusBadRequest, nil)
		} else {
			c.JSON(http.StatusInternalServerError, nil)
		}
		return
	}
	c.JSON(http.StatusOK, nil)
}

func UnlikeQuestionComment(c *gin.Context) {
	cid := c.Param("cid")
	if cid == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	uid := VisitorId(c)
	err := models.UnlikeQuestionComment(cid, uid)
	if err != nil {
		log.Println("controllers.UndoLikeQuestionComment(): ", err)
		if err.Error() == "reply is zero" {
			c.JSON(http.StatusBadRequest, nil)
		} else {
			c.JSON(http.StatusInternalServerError, nil)
		}
		return
	}
	c.JSON(http.StatusOK, nil)
}
