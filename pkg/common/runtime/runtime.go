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
