package main

import (
	"os"

	"github.com/gavin/blog/config"
	"github.com/gavin/blog/logger"
	"github.com/gavin/blog/middleware"
	"github.com/gavin/blog/routers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	defer deferClose()
	// 初始化日志
	logger.InitLogger()

	// 加载.env配置文件
	err := godotenv.Load()
	if err != nil {
		logger.Log.Error("Error loading .env file")
	}

	write := logger.Log.GetIoWriter()
	// 将 Gin 的日志输出指向 Zap
	// 重定向必须在 gin.New() 前
	gin.DefaultWriter = write
	gin.DefaultErrorWriter = write

	logger.Log.Infof("starting handlers")

	router := gin.New()
	router.Use(middleware.GinLogMiddleware())
	router.Use(middleware.GinRecoveryWithLogger())

	// 初始化数据库
	config.InitDB()

	config.Migrate()

	routers.InitApi(router)

	port := os.Getenv("PORT")
	logger.Log.Infof("handlers started  addr %s", port)
	router.Run(port) // 监听并在 0.0.0.0:8080 上启动服务
}

func deferClose() {
	// 延迟关闭日志，确保所有日志都写入
	logger.Log.Close()
}
