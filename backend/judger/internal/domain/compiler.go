package domain

type Compiler string

const (
	CompilerGCC Compiler = "/usr/bin/gcc"
	CompilerGPP Compiler = "/usr/bin/g++"
)
