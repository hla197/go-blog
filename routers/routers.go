package routers

import (
	"net/http"

	"github.com/gavin/blog/middleware"
	"github.com/gavin/blog/server"
	"github.com/gin-gonic/gin"
)

func InitApi(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	// 公共接口（不需要 token）
	public := router.Group("/auth")
	{
		public.POST("/login", server.Login)
		public.POST("/register", server.Register)
	}

	auth := router.Group("")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		post := auth.Group("/post")
		post.POST("add", server.AddPost)
		post.POST("update", server.UpdatePost)
		post.GET(":id", server.GetPost)
		post.DELETE(":id", server.DeletePost)
		post.GET("user", server.GetUserPost)
		post.POST("page", server.GetPagePosts)

		comment := auth.Group("/comment")
		comment.POST("add", server.AddComment)
		comment.POST("update", server.UpdateComment)
		comment.GET(":id", server.GetComment)
		comment.DELETE(":id", server.DeleteComment)
		comment.GET("user", server.GetUserComment)
		comment.POST("page", server.GetPageComments)
	}
}
