package logger

type Logger interface {
	Info(string)
	Warning(string)
	Debug(string)
	Fatal(string)
}