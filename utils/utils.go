package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"my-zhihu/config"
	"regexp"
	"time"
)

type Err struct {
	Msg  string
	Code int
}

func (e *Err) Error() string {
	return e.Msg
}

const (
	ErrAccountNotFound = 100000 + iota
	ErrIncorrectPassword
	ErrDuplicatedEmail
	ErrBadFullNameFormat
	ErrBadEmailFormat
	ErrBadPasswordFormat
)

func ValidateFullName(c *gin.Context, fullName string) *Err {
	re := regexp.MustCompile(`^[\p{Han}\w]+([\p{Han}\w\s.-]*)$`)
	if !re.MatchString(fullName) {
		return &Err{
			Msg:  "用户名中含有特殊字符",
			Code: ErrBadFullNameFormat,
		}
	}
	return nil
}

func ValidateUsername(email string) *Err {
	// panjinaho0320@qq.com
	re := regexp.MustCompile(`^[a-zA-Z0-9]+@(\w+).(\w{2,5})$`)
	if !re.MatchString(email) {
		return &Err{
			Msg:  "邮箱格式错误",
			Code: ErrBadEmailFormat,
		}
	}
	return nil
}

func ValidatePassword(password string) *Err {
	re := regexp.MustCompile(`^(\w+[\w[:graph:]]*){6,}$`)
	if !re.MatchString(password) {
		return &Err{
			Msg:  "密码格式错误",
			Code: ErrBadPasswordFormat,
		}
	}
	return nil
}

const (
	Second = 1
	Minute = 60 * Second
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Month  = 30 * Day
	Year   = 12 * Month
)

func FormatUnixTime(t int64) string {
	var res string
	paramTime := time.Unix(t, 0)
	paramYear, paramMonth, paramDay := paramTime.Date()
	now := time.Now()
	year, month, day := now.Date()
	if year == paramYear && month == paramMonth && day == paramDay {
		res = paramTime.Format("15:04")
	} else if now.Unix()-t < 2*Day {
		res = "昨天" + paramTime.Format("15:04")
	} else {
		res = paramTime.Format("2006-01-02 15:04")
	}
	return res
}

func FormatBeforeUnixTime(t int64) string {
	var res string
	diff := time.Now().Unix() - t
	switch {
	case diff < Minute:
		res = "刚刚"
	case Minute < diff && diff <= Hour:
		res = fmt.Sprintf("%d分钟前", diff/Minute)
	case Hour < diff && diff <= Day:
		res = fmt.Sprintf("%d小时前", diff/Hour)
	case Day < diff && diff <= Month:
		res = fmt.Sprintf("%d天前", diff/Day)
	case Month < diff && diff <= Year:
		res = fmt.Sprintf("%d月前", diff/Month)
	case diff >= Year:
		res = fmt.Sprintf("%d年前", diff/Year)
	}
	return res
}

func EncryptPassword(username, password string) string {
	h := md5.New()
	_, _ = io.WriteString(h, password)
	md5Password := fmt.Sprintf("%x", h.Sum(nil))
	_, _ = io.WriteString(h, config.Cfg.Server.Salt)
	_, _ = io.WriteString(h, username)
	_, _ = io.WriteString(h, md5Password)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func URLToken(urlToken *string, urlTokenCode int) {
	if urlTokenCode != 0 {
		*urlToken = fmt.Sprintf("%s-%d", *urlToken, urlTokenCode)
	}
}
