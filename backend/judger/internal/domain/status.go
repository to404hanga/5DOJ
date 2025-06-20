package domain

type Status string

const (
	StatusAccepted            Status = "Accepted"              // 正常情况
	StatusMemoryLimitExceeded Status = "Memory Limit Exceeded" // 内存超限
	StatusTimeLimitExceeded   Status = "Time Limit Exceeded"   // 时间超限
	StatusOutputLimitExceeded Status = "Output Limit Exceeded" // 输出超限
	StatusFileError           Status = "File Error"            // 文件错误
	StatusNonzeroExitStatus   Status = "Nonzero Exit Status"   // 非零返回值
	StatusSignalled           Status = "Signalled"             // 被信号杀死
	StatusInternalError       Status = "Internal Error"        // 内部错误
)
