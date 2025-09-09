package models

import (
    "time"
)

// Comment 代表一个评论
type Comment struct {
    ID        int       `json:"id" db:"id"`                 // 评论ID
    GoalID    int       `json:"goal_id" db:"goal_id"`       // 关联的目标ID
    ParentID  *int      `json:"parent_id" db:"parent_id"`   // 父评论ID（用于回复评论，可以为空）
    Content   string    `json:"content" db:"content"`       // 评论内容
    CreatedAt time.Time `json:"created_at" db:"created_at"` // 创建时间
}