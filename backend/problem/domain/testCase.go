package domain

type TestCaseView struct {
	Id        string
	Input     string
	Output    string
	Score     int
	CreatedBy string
	UpdatedBy string
	Enabled   bool
}
