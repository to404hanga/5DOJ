package service

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/to404hanga/5DOJ/judger/internal/domain"
)

var (
	//go:embed json/compileRequestTemplate.json
	compileRequestTemplate string
	//go:embed json/runRequestTemplate.json
	runRequestTemplate string
)

type JudgerService struct {
	baseUrl string
	client  *http.Client
}

func NewJudgerService(baseUrl string) *JudgerService {
	return &JudgerService{
		baseUrl: baseUrl,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (j *JudgerService) Compile(ctx context.Context, compiler domain.Compiler, filename string, content string) (fileId string, err error) {
	filenameWithoutExtension := strings.Split(filename, ".")[0]
	compilerRequest := fmt.Sprintf(compileRequestTemplate, compiler, filename, filenameWithoutExtension, filename, content, filenameWithoutExtension)

	req, _ := http.NewRequest("POST", j.baseUrl+"/run", strings.NewReader(compilerRequest))
	req.Header.Add("Content-Type", "application/json")

	resp, err := j.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result []domain.Result
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result[0].Status != domain.StatusAccepted {
		return "", fmt.Errorf("compile failed: %s, filename: %s", result[0].Status, filename)
	}

	return result[0].FileIds[filenameWithoutExtension], nil
}

func (j *JudgerService) Run(ctx context.Context, filenameWithoutExtension, fileId, input string, timeLimit, memoryLimit int) (output string, timeUsage, memoryUsage int, err error) {
	compilerRequest := fmt.Sprintf(runRequestTemplate, filenameWithoutExtension, input, memoryLimit, timeLimit*2, filenameWithoutExtension, fileId)

	req, _ := http.NewRequest("POST", j.baseUrl+"/run", strings.NewReader(compilerRequest))
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

	var result []domain.Result
	if err = json.Unmarshal(body, &result); err != nil {
		return
	}

	if result[0].Status != domain.StatusAccepted {
		return "", result[0].Time, result[0].Memory, fmt.Errorf("run failed: %s, fileId: %s", result[0].Status, fileId)
	}

	return result[0].Files["stdout"], result[0].Time, result[0].Memory, nil
}
