package controllers

import (
	"database/sql"
	"net/http"
	"starpool/config"
	"starpool/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CommentController 处理评论相关的HTTP请求
type CommentController struct{}

// CreateComment 创建新评论
// @Summary 创建新评论
// @Description 为指定目标创建新评论或回复已有评论
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "目标ID"
// @Param comment body models.Comment true "评论信息"
// @Success 201 {object} models.Comment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /goals/{id}/comments [post]
func (cc *CommentController) CreateComment(c *gin.Context) {
	// 获取路径参数
	goalId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目标ID"})
		return
	}

	// 解析请求体
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	// 如果提供了父评论ID，检查父评论是否存在
	if comment.ParentID != nil {
		var parentComment models.Comment
		query = `SELECT id FROM comments WHERE id = ? AND goal_id = ?`
		err = config.DB.QueryRow(query, *comment.ParentID, goalId).Scan(&parentComment.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "父评论未找到或不属于该目标"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
	}

	// 插入数据库
	query = `INSERT INTO comments(goal_id, parent_id, content, created_at) VALUES (?, ?, ?, NOW())`
	result, err := config.DB.Exec(query, goalId, comment.ParentID, comment.Content)
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

	comment.ID = int(id)
	comment.GoalID = goalId
	comment.CreatedAt = time.Now()

	// 返回创建的评论
	c.JSON(http.StatusCreated, comment)
}

// GetCommentsByGoalID 获取指定目标的所有评论（嵌套结构）
// @Summary 获取指定目标的所有评论
// @Description 获取指定目标的所有评论，包括回复，并以嵌套结构返回
// @Tags comments
// @Produce json
// @Param id path int true "目标ID"
// @Success 200 {array} models.Comment
// @Failure 404 {object} map[string]string
// @Router /goals/{id}/comments [get]
func (cc *CommentController) GetCommentsByGoalID(c *gin.Context) {
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

    // 查询评论
    query = `SELECT id, goal_id, parent_id, content, created_at FROM comments WHERE goal_id = ? ORDER BY created_at ASC`
    rows, err := config.DB.Query(query, goalId)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    // 遍历结果
    var comments []models.Comment
    for rows.Next() {
        var comment models.Comment
        err := rows.Scan(&comment.ID, &comment.GoalID, &comment.ParentID, &comment.Content, &comment.CreatedAt)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        comments = append(comments, comment)
    }

    // 检查遍历过程中是否有错误
    if err = rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 构建嵌套评论结构
    nestedComments := buildNestedComments(comments)
    
    // 计算总评论数
    totalComments := countTotalComments(comments)

    // 返回嵌套评论结构和总评论数
    c.JSON(http.StatusOK, gin.H{
        "comments": nestedComments,
        "total":    totalComments,
    })
}

// countTotalComments 计算评论总数（包括回复）
func countTotalComments(comments []models.Comment) int {
    return len(comments)
}

// buildNestedComments 构建嵌套评论结构
func buildNestedComments(comments []models.Comment) []map[string]interface{} {
    // 创建评论映射以便快速查找
    commentsMap := make(map[int]map[string]interface{})
    var rootComments []map[string]interface{}

    // 初始化评论映射
    for _, comment := range comments {
        commentsMap[comment.ID] = map[string]interface{}{
            "id":         comment.ID,
            "goal_id":    comment.GoalID,
            "parent_id":  comment.ParentID,
            "content":    comment.Content,
            "created_at": comment.CreatedAt,
            "children":   []map[string]interface{}{},
        }
    }

    // 构建评论树结构
    for _, comment := range comments {
        commentMap := commentsMap[comment.ID]
        if comment.ParentID == nil {
            rootComments = append(rootComments, commentMap)
        } else {
            // 确保父评论存在
            if parent, exists := commentsMap[*comment.ParentID]; exists {
                parent["children"] = append(parent["children"].([]map[string]interface{}), commentMap)
            }
        }
    }

    return rootComments
}
