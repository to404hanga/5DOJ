package domain

type TestCase struct {
	TestCaseId string `json:"testCaseId"`
	Input      string `json:"input"`
	Expected   string `json:"expected"`
	Score      int    `json:"score"`
}

type Problem struct {
	ProblemId    uint64     `json:"problemId"`
	TestCases    []TestCase `json:"testCases"`
	TotalScore   int        `json:"totalScore"`
	TimeLimitNS  uint64     `json:"timeLimitNS"`
	MemoryLimitB uint64     `json:"memoryLimitB"`
}
