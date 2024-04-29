package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"my-zhihu/config"
	"my-zhihu/router"
)

func main() {
	app := gin.Default()
	Init(app)
	_ = app.Run(config.Cfg.Server.Port)
}

func Init(engine *gin.Engine) {
	Session(engine)
	router.Route(engine)
}

func Session(engine *gin.Engine) {
	store := sessions.NewCookieStore([]byte(config.Cfg.Server.SessionSecret))
	store.Options = &sessions.Options{
		HttpOnly: true,
		MaxAge:   7 * 86400,
		Path:     "/",
	}
	sess := sessions.NewSession(store, "gin-session")
	engine.Use(sess)
}
