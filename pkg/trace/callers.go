package trace

import (
	"fmt"
	"runtime"
	"strings"
)

type Callers struct {
	Callers []Caller
	Path    string
}

type Caller struct {
	file     string
	line     int
	funcName string
}

func (c Callers) String() string {
	var message string
	format := "In %s:%d in function %s, called function %s\n"
	for i := 1; i < len(c.Callers); i++ {
		msg := fmt.Sprintf(format, c.Callers[i].file, c.Callers[i].line, c.Callers[i].funcName, c.Callers[i-1].funcName)
		message = msg + message
	}

	formattedPath := fmt.Sprintf("trace path: %s\n", c.Path)
	return formattedPath + message
}

func StackTrace() Callers {
	i := 2
	callerFunc := caller(i)
	callers := Callers{
		Callers: make([]Caller, 0),
		Path:    callerFunc.funcName,
	}

	const firstFuncName = "goexit"
	for ; ; i++ {
		if callerFunc = caller(i); callerFunc == nil || callerFunc.funcName == firstFuncName {
			break
		}

		callers.Callers = append(callers.Callers, *callerFunc)
		callers.Path = callerFunc.funcName + "->" + callers.Path
	}
	return callers
}

func caller(callerSkip int) *Caller {
	counter, file, line, ok := runtime.Caller(callerSkip)
	if !ok {
		return nil
	}
	return &Caller{
		file:     file,
		line:     line,
		funcName: baseNameWithMethod(runtime.FuncForPC(counter).Name()),
	}
}

func baseNameWithMethod(name string) string {
	if i := strings.Index(name, "."); i != -1 {
		return name[i+1:]
	}
	return name
}
