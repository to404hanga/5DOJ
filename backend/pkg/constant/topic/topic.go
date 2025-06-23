package topic

const TopicSubmitEvent = "submit_event"

type SubmitEvent struct {
	RecordId  uint64 `json:"recordId"`
	ContestId uint64 `json:"contestId"`
	ProblemId uint64 `json:"problemId"`
	UserId    uint64 `json:"userId"`
	Code      string `json:"code"`
	Language  string `json:"language"`
	Mode      string `json:"mode"`
}
