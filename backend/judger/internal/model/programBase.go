package model

import "time"

type ProgramBase struct {
	Id          uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement;comment:题目 id"`    // 题目 id
	Title       string    `gorm:"column:title;index:idx_title;type:varchar(255);not null;comment:题目名称"`     // 题目名称
	Level       int8      `gorm:"column:level;type:tinyint;index:idx_level;type:int;not null;comment:题目等级"` // 题目等级
	Tags        string    `gorm:"column:tags;type:varchar(255);default '';comment:题目标签，按英文逗号分隔"`            // 题目标签，按英文逗号分隔
	TestCaseNum int8      `gorm:"column:test_case_num;type:tinyint;not null;comment:测试用例数量"`                // 测试用例数量
	CreatedBy   uint64    `gorm:"column:created_by;type:bigint unsigned;not null;comment:创建者 id"`           // 创建者 id
	UpdatedBy   uint64    `gorm:"column:updated_by;type:bigint unsigned;not null;comment:更新者 id"`           // 更新者 id
	CreatedAt   time.Time `gorm:"column:created_at;not null;comment:创建时间"`                                  // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;comment:更新时间"`                                  // 更新时间
}

func (ProgramBase) TableName() string {
	return "program"
}
