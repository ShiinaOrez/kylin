package logger

import (
	"context"
	"fmt"
)

type Logger interface {
	Info(string)
	Warning(string)
	Debug(string)
	Fatal(string)
}

func GetLogger(ctx context.Context) Logger {
	if ctx == nil {
		return DefaultLogger{}
	}
	i := ctx.Value("kylin-logger")
	if i == nil {
		return DefaultLogger{}
	}
	if l, ok := ctx.Value("kylin-logger").(Logger); !ok {
		return DefaultLogger{}
	} else {
		return l
	}
}

type DefaultLogger struct {}

func (l DefaultLogger) Info(s string) {
	fmt.Println("[ INFO ]:", s)
}

func (l DefaultLogger) Warning(s string) {
	fmt.Println("[ WARN ]:", s)
}

func (l DefaultLogger) Debug(s string) {
	fmt.Println("[ DEBUG]:", s)
}

func (l DefaultLogger) Fatal(s string) {
	fmt.Println("[ Fatal]:", s)
}