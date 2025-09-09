package models

import (
    "time"
)

// StarGoal 代表一个星目标
type StarGoal struct {
    ID          int       `json:"id" db:"id"`                    // 目标ID
    Title       string    `json:"title" db:"title"`              // 目标标题
    Description string    `json:"description" db:"description"`  // 目标描述
    Category    string    `json:"category" db:"category"`        // 目标类别
    Stars       int       `json:"stars" db:"stars"`              // 星数
    CreatedAt   time.Time `json:"created_at" db:"created_at"`    // 创建时间
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`    // 更新时间
}