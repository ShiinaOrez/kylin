package logger

import "fmt"

type Logger interface {
	Info(string)
	Warning(string)
	Debug(string)
	Fatal(string)
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