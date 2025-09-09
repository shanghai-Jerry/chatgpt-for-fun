package models

import (
	"time"
)

// DailyRating 代表目标的每日评分记录
type DailyRating struct {
	ID        int       `json:"id" db:"id"`                 // 记录ID
	GoalID    int       `json:"goal_id" db:"goal_id"`       // 目标ID
	Rating    int       `json:"rating" db:"rating"`         // 评分 (1-5星)
	Date      time.Time `json:"date" db:"date"`             // 评分日期
	CreatedAt time.Time `json:"created_at" db:"created_at"` // 创建时间
}
