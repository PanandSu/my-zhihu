package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"my-zhihu/config"
	"net/http"
)

func SigninRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		uid := sess.Get(config.Cfg.Server.SessionKey)
		if uid == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "not authorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RefreshSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		uid := sess.Get(config.Cfg.Server.SessionKey)
		sess.Clear()
		if uid != nil {
			sess.Set(config.Cfg.Server.SessionKey, uid.(uint))
		}
		_ = sess.Save()
		c.Next()
	}
}
