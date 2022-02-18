package common

import (
	"runtime"
	"strings"
)

func GetFuncName() string {
	counter, _, _, success := runtime.Caller(1)
	if !success {
		panic("[common.GetFuncName()] Can't get function name.")
	}

	fullModulePath := runtime.FuncForPC(counter).Name()
	pathSplitted := strings.Split(fullModulePath, "/")
	return pathSplitted[len(pathSplitted)-1]
}
