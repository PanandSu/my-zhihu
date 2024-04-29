package router

import (
	"github.com/gin-gonic/gin"
	c "my-zhihu/controller"
	mw "my-zhihu/middleware"
)

func Route(r *gin.Engine) {
	r.LoadHTMLGlob("views/*")

	r.Static("/static", "./static")

	r.GET("/foo", c.Handle404) // 传HTML
	r.GET("/user/:id", c.Handle404)
	r.GET("/", c.IndexGet) // 传HTML

	r.GET("/signin", c.SigninGet) // 传HTML
	r.POST("/signin", c.SigninPost)
	r.GET("/signup", c.SignupGet) // 传HTML
	r.POST("/signup", c.SignupPost)
	r.GET("/logout", c.LogoutGet) // 传HTML

	r.GET("question/:qid", c.QuestionGet) //传HTML
	r.GET("question/:qid/answers/:aid",
		mw.RefreshSession(),
		c.AnswerGet) // 传HTML
	r.GET("topics/autocomplete", c.SearchTopics)

	api := r.Group("/api")
	{
		api.POST("answers/:id/voters", c.AnswerVoters)
		api.POST("questions/:id/followers", c.QuestionFollowers)
		api.POST("questions/:id/comments", c.QuestionComments)
		api.POST("members/:url_token", c.MemberInfo)
	}
	{
		api.Use(mw.SigninRequired())

		api.POST("questions", c.PostQuestion)
		api.POST("answers/:aid/voters", c.VoteAnswer)

		api.POST("questions/:qid/answers", c.PostAnswer)
		api.DELETE("answers/:aid", c.DeleteAnswer)
		api.POST("answers/:aid/actions/restore", c.RestoreAnswer)

		api.POST("questions/:qid/comments", c.PostQuestionComment)
		api.DELETE("questions/:qid/comments/:id", c.DeleteQuestionComment)

		api.POST("questions/:qid/followers", c.FollowQuestion)
		api.DELETE("questions/:qid/followers", c.UnfollowQuestion)

		api.POST("questions/:qid/comments/:cid/actions/like",
			c.LikeQuestionComment)
		api.DELETE("questions/:qid/comments/:cid/actions/like",
			c.UnlikeQuestionComment)

		api.POST("members/:url_token/followers", c.FollowMember)
		api.DELETE("members/:url_token/followers", c.UnfollowMember)

		api.POST("topics", c.PostTopic)
	}
}
