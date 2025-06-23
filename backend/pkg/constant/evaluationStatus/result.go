package evaluationStatus

type EvaluationStatusType string

const (
	AC  EvaluationStatusType = "Accept"
	CE  EvaluationStatusType = "Compile Error"
	WA  EvaluationStatusType = "Wrong Answer"
	TLE EvaluationStatusType = "Time Limit Exceed"
	MLE EvaluationStatusType = "Memory Limit Exceed"
	RE  EvaluationStatusType = "Runtime Error"
	UKE EvaluationStatusType = "Unknown Error"
	IQ  EvaluationStatusType = "In the Queue"
)
