package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"my-zhihu/config"
	"my-zhihu/models"
)

func Visitor(c *gin.Context) (*models.User, uint) {
	sess := sessions.Default(c)
	uid := sess.Get(config.Cfg.Server.SessionKey)
	if uid == nil {
		return nil, 0
	}
	user := models.GetUserByID(uid.(uint))
	if user == nil {
		return nil, 0
	}
	return user, uid.(uint)
}

func VisitorId(c *gin.Context) uint {
	_, uid := Visitor(c)
	return uid
}
