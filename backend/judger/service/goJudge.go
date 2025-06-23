package service

import (
	"5DOJ/judger/global"
	"5DOJ/pkg/constant/contestMode"
	"5DOJ/pkg/constant/evaluationStatus"
	"5DOJ/pkg/constant/language"
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type GoJudgeJudgerService struct {
	baseUrl string
	client  *http.Client
}

var _ IJudgerService = (*GoJudgeJudgerService)(nil)

//go:embed json/compiler.template.json
var goJudgeCompilerTemplate string

//go:embed json/runner.template.json
var goJudgeRunnerTemplate string

func NewGoJudgeJudgerService(baseUrl string) *GoJudgeJudgerService {
	uuid.SetRand(rand.Reader)
	return &GoJudgeJudgerService{
		baseUrl: baseUrl,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (j *GoJudgeJudgerService) Preheater(ctx context.Context, contestId uint64) (err error) {
	key := fmt.Sprintf("contest:preheater:%d", contestId)

	var val string
	val, err = global.Rds.Get(ctx, key).Result()
	if err != nil {
		return
	}

	return json.Unmarshal([]byte(val), &global.CP)
}

func (j *GoJudgeJudgerService) Judge(ctx context.Context, recordId, problemId uint64, lang language.LanguageType, userCode string, mode contestMode.ContestModeType) (evalutionStatus evaluationStatus.EvaluationStatusType, timeUsageMS, memoryUsageKB uint64, err error) {
	var fileId, filenameWithoutExt string
	fileId, filenameWithoutExt, err = j.compile(ctx, lang, userCode)
	if err != nil {
		evalutionStatus = evaluationStatus.CE
		return
	}

	problem := global.CP[problemId]
	var (
		op    string
		tu    uint64
		mu    uint64
		score int
	)
	defer func() {
		if errI := global.MySQL.WithContext(ctx).Where("id = ?", recordId).Updates(map[string]any{
			"status":          evalutionStatus,
			"time_usage_ms":   timeUsageMS,
			"memory_usage_kb": memoryUsageKB,
			"score":           score,
		}); errI != nil {
			// TODO 记录日志
		}
		if errI := j.remove(fileId); errI != nil {
			// TODO 记录日志
		}
	}()
	switch mode {
	case contestMode.ACM:
		for _, testCase := range problem.TestCases {
			op, tu, mu, err = j.run(ctx, fileId, filenameWithoutExt, testCase.Input, problem.TimeLimitNS, problem.MemoryLimitB)
			timeUsageMS = max(timeUsageMS, tu/1000000)
			memoryUsageKB = max(memoryUsageKB, mu>>10)
			if err != nil {
				evalutionStatus = evaluationStatus.RE
				return
			}
			if mu > problem.MemoryLimitB {
				evalutionStatus = evaluationStatus.MLE
				return
			}
			if tu > problem.TimeLimitNS {
				evalutionStatus = evaluationStatus.TLE
				return
			}
			if op != testCase.Expected {
				evalutionStatus = evaluationStatus.WA
				return
			}
		}
		evalutionStatus = evaluationStatus.AC
		score = problem.TotalScore
	case contestMode.IOI:
		for _, testCase := range problem.TestCases {
			op, tu, mu, err = j.run(ctx, fileId, filenameWithoutExt, testCase.Input, problem.TimeLimitNS, problem.MemoryLimitB)
			timeUsageMS = max(timeUsageMS, tu/1000000)
			memoryUsageKB = max(memoryUsageKB, mu>>10)
			if err != nil {
				evalutionStatus = evaluationStatus.RE
				continue
			}
			if mu > problem.MemoryLimitB {
				evalutionStatus = evaluationStatus.MLE
				continue
			}
			if tu > problem.TimeLimitNS {
				evalutionStatus = evaluationStatus.TLE
				continue
			}
			if op != testCase.Expected {
				evalutionStatus = evaluationStatus.WA
				continue
			}
			score += testCase.Score
		}
		if evalutionStatus != "" {
			return
		}
		evalutionStatus = evaluationStatus.AC
		score = problem.TotalScore
	default:
		evalutionStatus = evaluationStatus.UKE
	}
	return
}

func (j *GoJudgeJudgerService) remove(fileId string) (err error) {
	req, _ := http.NewRequest(http.MethodDelete, j.baseUrl+"/file/"+fileId, nil)
	_, err = j.client.Do(req)
	return
}

func (j *GoJudgeJudgerService) compile(ctx context.Context, lang language.LanguageType, userCode string) (fileId, filenameWithoutExt string, err error) {
	filenameWithoutExt = uuid.New().String()
	filename := filenameWithoutExt + "." + lang.String()
	compiler := lang.Compiler()

	compilerRequest := fmt.Sprintf(goJudgeCompilerTemplate, compiler, filename, filenameWithoutExt, filename, userCode, filenameWithoutExt)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, j.baseUrl+"/run", strings.NewReader(compilerRequest))
	req.Header.Add("Content-Type", "application/json")

	var resp *http.Response
	resp, err = j.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var result []GoJudgeResult
	if err = json.Unmarshal(body, &result); err != nil {
		return
	}

	if result[0].Status != Accepted {
		err = fmt.Errorf("compile error: %s", result[0].Error)
		return
	}

	fileId = result[0].FileIds[filenameWithoutExt]
	return
}

func (j *GoJudgeJudgerService) run(ctx context.Context, fileId, filenameWithoutExt, input string, timeLimitNS, memoryLimitB uint64) (output string, timeUsageNS, memoryUsageB uint64, err error) {
	runnerRequest := fmt.Sprintf(goJudgeRunnerTemplate, filenameWithoutExt, input, (memoryLimitB*3)>>1, timeLimitNS<<1, filenameWithoutExt, fileId)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, j.baseUrl+"/run", strings.NewReader(runnerRequest))
	req.Header.Add("Content-Type", "application/json")

	var resp *http.Response
	resp, err = j.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var result []GoJudgeResult
	if err = json.Unmarshal(body, &result); err != nil {
		return
	}

	if result[0].Status != Accepted {
		err = fmt.Errorf("compile error: %s", result[0].Error)
		return
	}

	return result[0].Files["stdout"], result[0].Time / 1000000, result[0].Memory >> 10, nil
}

type GoJudgeResult struct {
	Status     GoJudgeStatus     `json:"status"`
	Error      string            `json:"error"`      // 详细错误信息
	ExitStatus int               `json:"exitStatus"` // 程序返回值
	Time       uint64            `json:"time"`       // 程序运行 CPU 时间，单位 ns
	Memory     uint64            `json:"memory"`     // 程序运行内存，单位 byte
	ProcPeak   int               `json:"procPeak"`   // 程序运行最大线程数量
	RunTime    uint64            `json:"runTime"`    // 程序运行现实时间，单位 ns
	Files      map[string]string `json:"files"`      // copyOut 和 pipeCollector 指定的文件内容
	FileIds    map[string]string `json:"fileIds"`    // copyFileCached 指定的文件 id
	FileError  []FileErrorS      `json:"fileError"`  // 文件错误信息
}

type GoJudgeStatus string

const (
	Accepted            GoJudgeStatus = "Accepted"              // 正常情况
	MemoryLimitExceeded GoJudgeStatus = "Memory Limit Exceeded" // 内存超限
	TimeLimitExceeded   GoJudgeStatus = "Time Limit Exceeded"   // 时间超限
	OutputLimitExceeded GoJudgeStatus = "Output Limit Exceeded" // 输出超限
	FileError           GoJudgeStatus = "File Error"            // 文件错误
	NonzeroExitStatus   GoJudgeStatus = "Nonzero Exit Status"   // 非零返回值
	Signalled           GoJudgeStatus = "Signalled"             // 被信号杀死
	InternalError       GoJudgeStatus = "Internal Error"        // 内部错误
)

type FileErrorType string

const (
	CopyInOpenFile        FileErrorType = "CopyInOpenFile"
	CopyInCreateFile      FileErrorType = "CopyInCreateFile"
	CopyInCopyContent     FileErrorType = "CopyInCopyContent"
	CopyOutOpen           FileErrorType = "CopyOutOpen"
	CopyOutNotRegularFile FileErrorType = "CopyOutNotRegularFile"
	CopyOutSizeExceeded   FileErrorType = "CopyOutSizeExceeded"
	CopyOutCreateFile     FileErrorType = "CopyOutCreateFile"
	CopyOutCopyContent    FileErrorType = "CopyOutCopyContent"
	CollectSizeExceeded   FileErrorType = "CollectSizeExceeded"
)

type FileErrorS struct {
	Name    string        `json:"name"`
	Type    FileErrorType `json:"type"`
	Message string        `json:"message"`
}
