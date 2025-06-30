package domain

import "time"

type ProblemView struct {
	Id            uint64
	Title         string
	Level         int8
	CreatedBy     string
	UpdatedBy     string
	Enabled       bool
	TimeLimit     int
	MemoryLimit   int
	TotalScore    int
	TotalTestCase int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Markdown      string
}
