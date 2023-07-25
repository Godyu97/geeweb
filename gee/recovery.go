package gee

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				geeLogger.Error(trace(fmt.Sprintf("%s", err)))
				ctx.Fail(http.StatusInternalServerError, errors.New("Internal Server Error"))
			}
		}()
		ctx.Next()
	}
}

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
