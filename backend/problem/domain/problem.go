package domain

import "time"

type ProblemView struct {
	Id          uint64
	Title       string
	Level       int8
	CreatedBy   uint64
	UpdatedBy   uint64
	Enabled     bool
	TimeLimit   int
	MemoryLimit int
	TotalScore  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Markdown    string
}
