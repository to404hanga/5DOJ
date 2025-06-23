package language

import "fmt"

type LanguageType string

const (
	C    LanguageType = "c"
	CPP  LanguageType = "cpp"
	PY   LanguageType = "py"
	GO   LanguageType = "go"
	JAVA LanguageType = "java"
)

func (l LanguageType) Compiler() string {
	switch l {
	case C:
		return "/usr/bin/gcc"
	case CPP:
		return "/usr/bin/g++"
	default:
		panic(fmt.Sprintf("unsupported language %s", l))
	}
}

func (l LanguageType) String() string {
	return string(l)
}
