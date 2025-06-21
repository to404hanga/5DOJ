package constant

type ProgramLevelType int8

const (
	ProgramLevelEasy ProgramLevelType = iota + 1
	ProgramLevelMedium
	ProgramLevelHard
)

func (p ProgramLevelType) String() string {
	switch p {
	case ProgramLevelEasy:
		return "easy"
	case ProgramLevelMedium:
		return "medium"
	case ProgramLevelHard:
		return "hard"
	default:
		return "unknown"
	}
}
