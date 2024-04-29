package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"my-zhihu/models"
	"net/http"
)

func MemberInfo(c *gin.Context) {
	uid := VisitorId(c)
	urlToken := c.Param("url_token")
	if urlToken == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	member := models.GetUserByURLToken(urlToken, uid)
	if member != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	c.JSON(http.StatusOK, member)
}

func FollowMember(c *gin.Context) {
	urlToken := c.Param("url_token")
	if urlToken == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	uid := VisitorId(c)
	err := models.FollowMember(urlToken, uid)
	flag := true
	if err != nil {
		flag = false
		log.Println("controllers.FollowMember(): ", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"succeed": flag,
	})
}

func UnfollowMember(c *gin.Context) {
	urlToken := c.Param("url_token")
	if urlToken == "" {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	uid := VisitorId(c)
	err := models.UnfollowMember(urlToken, uid)
	flag := true
	if err != nil {
		flag = false
		log.Println("controllers.UnfollowMember(): ", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"succeed": flag,
	})
}
