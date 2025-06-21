package domain

type ProgramView struct {
	ProgramId     uint64   `json:"programId"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	CreatedBy     uint64   `json:"createdBy"`
	Creator       string   `json:"creator"`
	UpdatedBy     uint64   `json:"updatedBy"`
	Updator       string   `json:"updator"`
	Level         string   `json:"level"`
	TimeLimitMS   uint64   `json:"timeLimitMS"`
	MemoryLimitMB uint64   `json:"memoryLimitMB"`
	Tags          []string `json:"tags"`
	CreatedAt     string   `json:"createdAt"`
	UpdatedAt     string   `json:"updatedAt"`
	Enabled       bool     `json:"enabled"`
}

type TestCaseView struct {
	TestCaseId string `json:"testCaseId"`
	Inputs     string `json:"input"`
	Expecteds  string `json:"expected"`
}
