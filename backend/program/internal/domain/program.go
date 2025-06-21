package domain

type ProgramView struct {
	ProgramId   uint64   `json:"programId"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	CreatedBy   uint64   `json:"createdBy"`
	Creator     string   `json:"creator"`
	UpdatedBy   uint64   `json:"updatedBy"`
	Updator     string   `json:"updator"`
	Level       string   `json:"level"`
	TimeLimit   uint64   `json:"timeLimit"`
	MemoryLimit uint64   `json:"memoryLimit"`
	Tags        []string `json:"tags"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}
