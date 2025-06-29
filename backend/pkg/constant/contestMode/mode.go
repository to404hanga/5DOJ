package contestMode

type ContestModeType int8

const (
	Unknown ContestModeType = iota
	ACM
	IOI
)

func (c ContestModeType) String() string {
	switch c {
	case 1:
		return "ACM"
	case 2:
		return "IOI"
	default:
		return "Unknown"
	}
}
