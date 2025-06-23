package model

import "time"

type Record struct {
	Id            uint64    `gorm:"primaryKey;autoIncrement;type:bigint unsigned;comment:记录 id"`                        // 记录 id
	UserId        uint64    `gorm:"type:bigint unsigned;not null;index:idx_user_id;comment:用户 id"`                      // 用户 id
	ProblemId     uint64    `gorm:"type:bigint unsigned;not null;index:idx_problem_id;comment:题目 id"`                   // 题目 id
	ContestId     uint64    `gorm:"type:bigint unsigned;not null;index:idx_contest_id;comment:比赛 id"`                   // 比赛 id
	Language      string    `gorm:"type:varchar(20);not null;index:idx_language;comment:语言"`                            // 语言
	Score         int       `gorm:"type:int;not null;default:0;comment:分数"`                                             // 分数
	Result        string    `gorm:"type:varchar(20);not null;default:'In the Queue';comment:结果"`                        // 结果
	timeUsageMS   uint64    `gorm:"column:time_usage_ms;type:bigint unsigned;not null;default:0;comment:时间使用量"`         // 时间使用量
	memoryUsageKB uint64    `gorm:"column:memory_usage_kb;type:bigint unsigned;not null;default:0;comment:内存使用量"`       // 内存使用量
	CodeId        string    `gorm:"type:varchar(50);not null;default:'';comment:代码 id"`                                 // 代码 id
	CreatedAt     time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;index:idx_created_at;comment:创建时间"` // 创建时间
}
