package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"my-zhihu/models"
	"net/http"
)

func IndexGet(c *gin.Context) {
	user, uid := Visitor(c)
	topStory := models.HomeTimeLine(uid)
	v := DataState{
		Page:     "index",
		TopStory: topStory,
	}
	dataState, _ := json.Marshal(&v)
	c.HTML(http.StatusOK, "index.tmpl",
		gin.H{
			"user":      user,
			"topStory":  topStory,
			"dataState": string(dataState),
		},
	)
}
