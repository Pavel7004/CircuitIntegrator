package runtime

import (
	"path"
	"reflect"
	"runtime"
	"strings"
)

func GetFuncModule(fn interface{}) string {
	fullPath := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	name := path.Base(fullPath)
	return strings.Split(name, ".")[0]
}

func GetParentFuncName() string {
	counter, _, _, success := runtime.Caller(2)
	if !success {
		panic("[misc.GetParentFuncName()] Can't get function name.")
	}

	return path.Base(runtime.FuncForPC(counter).Name())
}
