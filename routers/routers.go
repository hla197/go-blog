package routers

import (
	"net/http"

	"github.com/gavin/blog/handlers"
	"github.com/gavin/blog/middleware"
	"github.com/gin-gonic/gin"
)

func InitApi(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	authHandler := &handlers.AuthHandler{}
	commentHandle := &handlers.CommentHandle{}
	postHandler := &handlers.PostHandler{}

	// 公共接口（不需要 token）
	public := router.Group("/auth")
	{
		public.POST("/login", authHandler.Login)
		public.POST("/register", authHandler.Register)
	}

	auth := router.Group("")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		post := auth.Group("/post")
		post.POST("add", postHandler.AddPost)
		post.POST("update", postHandler.UpdatePost)
		post.GET(":id", postHandler.GetPost)
		post.DELETE(":id", postHandler.DeletePost)
		post.GET("user", postHandler.GetUserPost)
		post.POST("page", postHandler.GetPagePosts)

		comment := auth.Group("/comment")
		comment.POST("add", commentHandle.AddComment)
		comment.POST("update", commentHandle.UpdateComment)
		comment.GET(":id", commentHandle.GetComment)
		comment.DELETE(":id", commentHandle.DeleteComment)
		comment.GET("user", commentHandle.GetUserComment)
		comment.POST("page", commentHandle.GetPageComments)
	}
}
