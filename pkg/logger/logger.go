package logger

type M map[string]interface{}

type Logger interface {
	Named(name string) Logger
	Debug(message string, args M)
	Info(message string, args M)
	Warn(message string, args M)
	Error(message string, args M)
	Fatal(message string, args M)
}
