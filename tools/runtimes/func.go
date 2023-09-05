package runtimes

import "runtime"

func CurrentFuncName(skip ...int) string {
	pc := make([]uintptr, 1)
	s := 2
	if len(skip) > 0 {
		s += skip[0]
	}
	runtime.Callers(s, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}
