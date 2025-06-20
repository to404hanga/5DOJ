package service

import (
	"context"
	"testing"

	"github.com/to404hanga/5DOJ/judger/internal/domain"
)

func InitJudger() *JudgerService {
	return NewJudgerService("http://localhost:5050")
}

func TestCompile(t *testing.T) {
	judger := InitJudger()

	testcases := []struct {
		Name     string
		Compiler domain.Compiler
		Filename string
		Content  string
	}{
		{
			Name:     "编译 C++ 代码",
			Compiler: domain.CompilerGPP,
			Filename: "a.cc",
			Content:  "#include <iostream>\\nusing namespace std;\\nint main() {\\nint a, b;\\ncin >> a >> b;\\ncout << a + b << endl;\\n}",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			fileId, err := judger.Compile(context.Background(), testcase.Compiler, testcase.Filename, testcase.Content)
			if err != nil {
				t.Errorf("编译失败: %v", err)
				return
			}
			t.Logf("fileId: %v", fileId)
		})
	}
}

func TestRun(t *testing.T) {
	judger := InitJudger()

	testcases := []struct {
		Name                     string
		FilenameWithoutExtension string
		FileId                   string
		Input                    string
		TimeLimit                int
		MemoryLimit              int
	}{
		{
			Name:                     "运行 C++ 代码",
			FilenameWithoutExtension: "a",
			FileId:                   "NGM5LGZU",
			Input:                    "1 2",
			TimeLimit:                1000000000, // 1s
			MemoryLimit:              536870912,  // 512MB
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.Name, func(t *testing.T) {
			output, timeUsage, memoryUsage, err := judger.Run(context.Background(), testcase.FilenameWithoutExtension, testcase.FileId, testcase.Input, testcase.TimeLimit, testcase.MemoryLimit)
			if err != nil {
				t.Errorf("运行失败: %v", err)
				return
			}
			t.Logf("output: %v", output)
			t.Logf("timeUsage: %v", timeUsage)
			t.Logf("memoryUsage: %v", memoryUsage)
		})
	}
}
