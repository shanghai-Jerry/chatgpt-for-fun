package main

import (
	"log"
	"starpool/config"
	"starpool/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	config.ConnectDB()

	// 创建gin路由器
	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	// 配置CORS
	config := cors.Config{
		AllowOrigins:     []string{"*"}, // 前端服务地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))

	// 注册路由
	routes.RegisterGoalRoutes(router)

	// 启动服务器
	log.Println("服务器启动在端口 8080 ，模式为 DebugMode")
	router.Run(":8080")
}
