package domain

type TestCaseView struct {
	Id        string
	Input     string
	Output    string
	Score     int
	CreatedBy uint64
	UpdatedBy uint64
	Enabled   bool
}
