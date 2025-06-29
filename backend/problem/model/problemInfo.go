package model

import "time"

type ProblemInfo struct {
	Id            uint64 `gorm:"type:bigint unsigned;primaryKey;autoIncrement;comment:'题目 ID'"`                   // 题目 ID
	Title         string `gorm:"type:varchar(20);not null;comment:'题目标题'"`                                        // 题目标题
	Level         int8   `gorm:"type:tinyint;not null;index:idx_level;comment:'题目难度等级: 0=未定义, 1=简单, 2=中等, 3=困难'"` // 题目难度等级: 0=未定义, 1=简单, 2=中等, 3=困难
	CreatedBy     uint64 `gorm:"type:bigint unsigned;not null;comment:'创建者 ID'"`                                  // 创建者 ID
	UpdatedBy     uint64 `gorm:"type:bigint unsigned;not null;comment:'最后更新者 ID'"`                                // 最后更新者 ID
	Enabled       bool   `gorm:"type:tinyint;not null;default:1;comment:'是否启用: 0=禁用, 1=启用'"`                      // 是否启用: 0=禁用, 1=启用
	TimeLimit     int    `gorm:"type:int;not null;default:1000;comment:'时间限制: 单位 ms'"`                            // 时间限制: 单位 ms
	MemoryLimit   int    `gorm:"type:int;not null;default:256;comment:'内存限制: 单位 MB'"`                             // 内存限制: 单位 MB
	TotalScore    int    `gorm:"type:int;not null;default:0;comment:'总分数（仅 IOI 赛制有效）'"`                           // 总分数（仅 IOI 赛制有效）
	TotalTestCase int    `gorm:"type:int;not null;default:0;comment:'总测试点数'"`                                     // 总测试点数

	CreatedAt time.Time // 题目创建时间
	UpdatedAt time.Time // 题目最后更新时间
}
