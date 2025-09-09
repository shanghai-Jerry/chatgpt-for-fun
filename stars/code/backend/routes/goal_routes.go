package routes

import (
	"starpool/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterGoalRoutes 注册星目标相关的路由
func RegisterGoalRoutes(router *gin.Engine) {
	goalController := &controllers.GoalController{}
	commentController := &controllers.CommentController{}

	// 目标管理路由
	router.POST("/goals", goalController.CreateGoal)
	router.GET("/goals", goalController.GetGoals)
	router.GET("/goals/:id", goalController.GetGoalByID)
	router.PUT("/goals/:id", goalController.UpdateGoal)
	router.DELETE("/goals/:id", goalController.DeleteGoal)
	router.GET("/goals/category/:category", goalController.GetGoalsByCategory)
	// 添加获取总星数的路由
	router.GET("/stars", goalController.GetTotalStars)
	// 添加每日评分路由
	router.POST("/goals/:id/daily-rating", goalController.AddDailyRating)
	router.GET("/goals/:id/daily-ratings", goalController.GetDailyRatings)
	
	// 添加评论路由
	router.POST("/goals/:id/comments", commentController.CreateComment)
	router.GET("/goals/:id/comments", commentController.GetCommentsByGoalID)
}
