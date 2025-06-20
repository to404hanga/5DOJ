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
			}
			t.Logf("fileId: %v", fileId)
		})
	}
}
