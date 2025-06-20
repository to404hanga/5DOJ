package domain

type Result struct {
	Status     Status            `json:"status"`
	Error      string            `json:"error"`      // 详细错误信息
	ExitStatus int               `json:"exitStatus"` // 程序返回值
	Time       int               `json:"time"`       // 程序运行 CPU 时间，单位 ns
	Memory     int               `json:"memory"`     // 程序运行内存，单位 byte
	ProcPeak   int               `json:"procPeak"`   // 程序运行最大线程数量
	RunTime    int               `json:"runTime"`    // 程序运行现实时间，单位 ns
	Files      map[string]string `json:"files"`      // copyOut 和 pipeCollector 指定的文件内容
	FileIds    map[string]string `json:"fileIds"`    // copyFileCached 指定的文件 id
	FileError  []FileError       `json:"fileError"`  // 文件错误信息
}
