//go:build prod

package log

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

func Debug(args ...any) {
}

func Error(message string, err error, fatal ...bool) {
	now := time.Now().Format("03:04:05 PM")
	fmt.Print("[" + now)
	skip := 1
	if len(fatal) > 0 && fatal[0] {
		skip = 2
	}
	pc, _, line, ok := runtime.Caller(skip)
	if !ok {
		panic("No caller information")
	}
	fmt.Print(" " + runtime.FuncForPC(pc).Name() + ":" + strconv.Itoa(line) + "]")
	fmt.Print(" " + message)
	if err == nil {
		println()
		return
	} else {
		fmt.Print(" | " + err.Error())
	}
}

func Fatal(message string, err error) {
	Error(message, err, true)
	os.Exit(1)
}
