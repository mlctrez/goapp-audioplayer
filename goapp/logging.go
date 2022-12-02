package goapp

import (
	"fmt"
	"runtime"
	"strings"
)

var BuildDir string
var Module string

type Logging struct {
}

func (l *Logging) Log(message string) {
	caller, fileRaw, line, ok := runtime.Caller(1)
	if ok {
		funcName, file := trim(caller, fileRaw)
		fmt.Printf("%s:%d %s %s\n", file, line, funcName, message)
	}
}

func (l *Logging) Logf(format string, args ...any) {
	caller, fileRaw, line, ok := runtime.Caller(1)
	if ok {
		funcName, file := trim(caller, fileRaw)
		message := fmt.Sprintf(format, args...)
		fmt.Printf("%s:%d %s %s\n", file, line, funcName, message)
	}
}

func trim(caller uintptr, file string) (string, string) {
	funcName := runtime.FuncForPC(caller).Name()
	funcName = strings.TrimPrefix(funcName, Module+"/")
	file = strings.TrimPrefix(file, BuildDir+"/")
	return funcName, file
}
