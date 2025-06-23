package domain

import (
	"time"
)

type QueryView struct {
	RecordId      uint64
	ContestId     uint64
	ProblemId     uint64
	UserId        uint64
	Language      string
	Score         int
	Result        string
	TimeUsageMS   uint64
	MemoryUsageKB uint64
	Code          string
	SubmitTime    time.Time
	UserName      string
	ProblemTitle  string
}
