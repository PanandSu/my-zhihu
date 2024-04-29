package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"my-zhihu/config"
	"my-zhihu/models"
	u "my-zhihu/utils"
	"net/http"
)

func SigninGet(c *gin.Context) {
	uid := VisitorId(c)
	if uid == 0 {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	c.HTML(http.StatusOK, "signin.html", nil)
}

// SigninPost 登录
func SigninPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	err := u.ValidateUsername(username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Msg,
			"code":    err.Code,
		})
		return
	}
	if len(password) < 6 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "请输入6位以上的密码",
			"code":    u.ErrBadPasswordFormat,
		})
	}
	pwdEncrypted := u.EncryptPassword(username, password)
	user := models.GetUserByUsername(username)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "该邮箱账号未注册知乎",
			"code":    u.ErrAccountNotFound,
		})
	}
	if user.Password != pwdEncrypted {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "账号或密码错误",
			"code":    u.ErrIncorrectPassword,
		})
	}
	sess := sessions.Default(c)
	sess.Clear()
	sess.Set(config.Cfg.Server.SessionKey, user.ID)
	_ = sess.Save()
	c.JSON(http.StatusCreated, nil)
}

func SignupGet(c *gin.Context) {
	uid := VisitorId(c)
	if uid == 0 {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	c.HTML(http.StatusOK, "signup.html", nil)
}

func SignupPost(c *gin.Context) {
	var err *u.Err
	fullName := c.PostForm("fullName")
	username := c.PostForm("username")
	password := c.PostForm("password")
	err = u.ValidateFullName(c, fullName)
	if err != nil {
		log.Println("fullName err:", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Msg,
			"code":    err.Code,
		})
		return
	}
	err = u.ValidateUsername(username)
	if err != nil {
		log.Println("username err:", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Msg,
			"code":    err.Code,
		})
		return
	}
	err = u.ValidatePassword(password)
	if err != nil {
		log.Println("password err:", err)
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Msg,
			"code":    err.Code,
		})
		return
	}
	encryptedPassword := u.EncryptPassword(username, password)
	user := &models.User{
		Name:     fullName,
		Email:    username,
		Password: encryptedPassword,
	}
	uid, err := models.InsertUser(user)
	println(uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "该邮箱已经注册,请直接登录",
			"code":    u.ErrDuplicatedEmail,
		})
		return
	}
	sess := sessions.Default(c)
	sess.Clear()
	sess.Set(config.Cfg.Server.SessionKey, uid)
	_ = sess.Save()
	c.JSON(http.StatusCreated, nil)
}

func LogoutGet(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Clear()
	_ = sess.Save()
	nextPath := c.GetHeader("referer")
	c.Redirect(http.StatusSeeOther, nextPath)
}
