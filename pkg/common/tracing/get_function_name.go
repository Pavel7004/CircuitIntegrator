package tracing

import (
	"context"
	"runtime"
	"strings"

	"github.com/opentracing/opentracing-go"
)

func StartSpanFromContext(ctx context.Context) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, getFuncName())
}

func getFuncName() string {
	counter, _, _, success := runtime.Caller(2)
	if !success {
		panic("[common.GetFuncName()] Can't get function name.")
	}

	fullModulePath := runtime.FuncForPC(counter).Name()
	pathSplitted := strings.Split(fullModulePath, "/")
	return pathSplitted[len(pathSplitted)-1]
}
