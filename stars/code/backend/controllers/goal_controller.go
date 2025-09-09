package controllers

import (
	"database/sql"
	"net/http"
	"starpool/config"
	"starpool/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GoalController 处理星目标相关的HTTP请求
type GoalController struct{}

// CreateGoal 创建新目标
// @Summary 创建新目标
// @Description 创建一个新的星目标
// @Tags goals
// @Accept json
// @Produce json
// @Param goal body models.StarGoal true "目标信息"
// @Success 201 {object} models.StarGoal
// @Failure 400 {object} map[string]string
// @Router /goals [post]
func (gc *GoalController) CreateGoal(c *gin.Context) {
	var goal models.StarGoal

	// 解析请求体
	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 插入数据库
	query := `INSERT INTO star_goals(title, description, category, stars, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())`
	result, err := config.DB.Exec(query, goal.Title, goal.Description, goal.Category, goal.Stars)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取插入记录的ID
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	goal.ID = int(id)

	// 返回创建的目标
	c.JSON(http.StatusCreated, goal)
}

// GetGoals 获取所有目标
// @Summary 获取所有目标
// @Description 获取所有已创建的星目标
// @Tags goals
// @Produce json
// @Success 200 {array} models.StarGoal
// @Router /goals [get]
func (gc *GoalController) GetGoals(c *gin.Context) {
	// 查询数据库
	query := `SELECT id, title, description, category, stars, created_at, updated_at FROM star_goals`
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// 遍历结果
	var goals []models.StarGoal
	for rows.Next() {
		var goal models.StarGoal
		err := rows.Scan(&goal.ID, &goal.Title, &goal.Description, &goal.Category, &goal.Stars, &goal.CreatedAt, &goal.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		goals = append(goals, goal)
	}

	// 检查遍历过程中是否有错误
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回所有目标
	c.JSON(http.StatusOK, goals)
}

// GetGoalByID 根据ID获取单个目标
// @Summary 根据ID获取单个目标
// @Description 根据ID获取特定的星目标
// @Tags goals
// @Produce json
// @Param id path int true "目标ID"
// @Success 200 {object} models.StarGoal
// @Failure 404 {object} map[string]string
// @Router /goals/{id} [get]
func (gc *GoalController) GetGoalByID(c *gin.Context) {
	// 获取路径参数
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目标ID"})
		return
	}

	// 查询数据库
	var goal models.StarGoal
	query := `SELECT id, title, description, category, stars, created_at, updated_at FROM star_goals WHERE id = ?`
	err = config.DB.QueryRow(query, id).Scan(&goal.ID, &goal.Title, &goal.Description, &goal.Category, &goal.Stars, &goal.CreatedAt, &goal.UpdatedAt)

	// 处理查询结果
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "目标未找到"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 返回目标
	c.JSON(http.StatusOK, goal)
}

// UpdateGoal 更新目标
// @Summary 更新目标
// @Description 更新特定星目标的信息
// @Tags goals
// @Accept json
// @Produce json
// @Param id path int true "目标ID"
// @Param goal body models.StarGoal true "更新的目标信息"
// @Success 200 {object} models.StarGoal
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /goals/{id} [put]
func (gc *GoalController) UpdateGoal(c *gin.Context) {
	// 获取路径参数
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目标ID"})
		return
	}

	// 解析请求体
	var goal models.StarGoal
	if err := c.ShouldBindJSON(&goal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新数据库
	query := `UPDATE star_goals SET title = ?, description = ?, category = ?, stars = ?, updated_at = NOW() WHERE id = ?`
	result, err := config.DB.Exec(query, goal.Title, goal.Description, goal.Category, goal.Stars, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 检查是否有记录被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "目标未找到"})
		return
	}

	// 返回更新后的目标
	goal.ID = id
	c.JSON(http.StatusOK, goal)
}

// DeleteGoal 删除目标
// @Summary 删除目标
// @Description 删除特定的星目标
// @Tags goals
// @Produce json
// @Param id path int true "目标ID"
// @Success 204 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /goals/{id} [delete]
func (gc *GoalController) DeleteGoal(c *gin.Context) {
	// 获取路径参数
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目标ID"})
		return
	}

	// 删除数据库记录
	query := `DELETE FROM star_goals WHERE id = ?`
	result, err := config.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 检查是否有记录被删除
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "目标未找到"})
		return
	}

	// 返回成功响应
	c.Status(http.StatusNoContent)
}

// GetGoalsByCategory 根据类别获取目标
// @Summary 根据类别获取目标
// @Description 根据类别获取星目标
// @Tags goals
// @Produce json
// @Param category path string true "目标类别"
// @Success 200 {array} models.StarGoal
// @Router /goals/category/{category} [get]
func (gc *GoalController) GetGoalsByCategory(c *gin.Context) {
	// 获取路径参数
	category := c.Param("category")

	// 查询数据库
	query := `SELECT id, title, description, category, stars, created_at, updated_at FROM star_goals WHERE category = ?`
	rows, err := config.DB.Query(query, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// 遍历结果
	var goals []models.StarGoal
	for rows.Next() {
		var goal models.StarGoal
		err := rows.Scan(&goal.ID, &goal.Title, &goal.Description, &goal.Category, &goal.Stars, &goal.CreatedAt, &goal.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		goals = append(goals, goal)
	}

	// 检查遍历过程中是否有错误
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回符合条件的目标
	c.JSON(http.StatusOK, goals)
}

// GetTotalStars 获取所有目标的总星数
// @Summary 获取所有目标的总星数
// @Description 获取所有星目标的星数总和
// @Tags goals
// @Produce json
// @Success 200 {object} map[string]int
// @Router /stars/total [get]
func (gc *GoalController) GetTotalStars(c *gin.Context) {
	// 查询数据库获取总星数
	var totalStars int
	query := `SELECT COALESCE(SUM(stars), 0) FROM star_goals`
	err := config.DB.QueryRow(query).Scan(&totalStars)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回总星数
	result := map[string]int{"total_stars": totalStars}
	c.JSON(http.StatusOK, result)
}

// AddDailyRating 为指定目标添加每日评分
// @Summary 为指定目标添加每日评分
// @Description 为指定目标添加每日评分记录
// @Tags goals
// @Accept json
// @Produce json
// @Param id path int true "目标ID"
// @Param rating body models.DailyRating true "评分信息"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /goals/{id}/daily-rating [post]
func (gc *GoalController) AddDailyRating(c *gin.Context) {
	// 获取路径参数
	goalId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目标ID"})
		return
	}

	// 解析请求体
	var rating models.DailyRating
	if err := c.ShouldBindJSON(&rating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证评分值
	if rating.Rating < 1 || rating.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评分必须在1-5之间"})
		return
	}

	// 检查目标是否存在
	var goal models.StarGoal
	query := `SELECT id FROM star_goals WHERE id = ?`
	err = config.DB.QueryRow(query, goalId).Scan(&goal.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "目标未找到"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 插入或更新每日评分记录
	query = `INSERT INTO daily_ratings (goal_id, rating, date, created_at) VALUES (?, ?, ?, NOW()) 
             ON DUPLICATE KEY UPDATE rating = ?, created_at = NOW()`
	_, err = config.DB.Exec(query, goalId, rating.Rating, rating.Date, rating.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 重新计算并更新目标的总星数
	var totalStars int
	query = `SELECT COALESCE(SUM(rating), 0) FROM daily_ratings WHERE goal_id = ?`
	err = config.DB.QueryRow(query, goalId).Scan(&totalStars)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新目标的总星数
	query = `UPDATE star_goals SET stars = ?, updated_at = NOW() WHERE id = ?`
	_, err = config.DB.Exec(query, totalStars, goalId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回成功响应
	response := map[string]string{"message": "评分记录成功"}
	c.JSON(http.StatusOK, response)
}

// GetDailyRatings 获取指定目标的所有每日评分记录
// @Summary 获取指定目标的所有每日评分记录
// @Description 获取指定目标的所有每日评分记录
// @Tags goals
// @Produce json
// @Param id path int true "目标ID"
// @Success 200 {array} models.DailyRating
// @Failure 404 {object} map[string]string
// @Router /goals/{id}/daily-ratings [get]
func (gc *GoalController) GetDailyRatings(c *gin.Context) {
	// 获取路径参数
	goalId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目标ID"})
		return
	}

	// 检查目标是否存在
	var goal models.StarGoal
	query := `SELECT id FROM star_goals WHERE id = ?`
	err = config.DB.QueryRow(query, goalId).Scan(&goal.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "目标未找到"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 查询每日评分记录
	query = `SELECT id, goal_id, rating, date, created_at FROM daily_ratings WHERE goal_id = ? ORDER BY date DESC limit 7 `
	rows, err := config.DB.Query(query, goalId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// 遍历结果
	var ratings []models.DailyRating
	for rows.Next() {
		var rating models.DailyRating
		err := rows.Scan(&rating.ID, &rating.GoalID, &rating.Rating, &rating.Date, &rating.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ratings = append(ratings, rating)
	}

	// 检查遍历过程中是否有错误
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回评分记录
	c.JSON(http.StatusOK, ratings)
}
