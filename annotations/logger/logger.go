package logger

type Logger struct {
	Level string `annotation:"name=level,default=info,oneOf=info;debug"`
}
