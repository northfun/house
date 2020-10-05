package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Config struct {
	Path  string
	Env   string
	Level int
}

var (
	l  *zap.Logger
	sl *zap.SugaredLogger
)

func Init(conf *Config) error {
	var zConf zap.Config

	switch conf.Env {
	case "production":
		zConf = zap.NewProductionConfig()
	default:
		zConf = zap.NewDevelopmentConfig()
	}

	zConf.OutputPaths = []string{conf.Path}

	var err error
	l, err = zConf.Build()
	if err != nil {
		Println("logger build err, %v,%s,%s", err, conf.Path, ".")
		return err
	}

	sl = l.Sugar()

	return nil
}

func Sync() error {
	if l == nil {
		return nil
	}

	return l.Sync()
}

func Debug(msg string, fields ...zap.Field) {
	if l == nil {
		return
	}
	l.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	if l == nil {
		return
	}
	l.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	if l == nil {
		return
	}
	l.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	if l == nil {
		return
	}
	l.Error(msg, fields...)
}

func Debugf(template string, args ...interface{}) {
	if sl == nil {
		return
	}
	sl.Debugf(template, args)
}

func Infof(template string, args ...interface{}) {
	if sl == nil {
		return
	}
	sl.Info(template, args)
}

func Warnf(template string, args ...interface{}) {
	if sl == nil {
		return
	}
	sl.Warnf(template, args)
}

func Errorf(template string, args ...interface{}) {
	if sl == nil {
		return
	}
	sl.Errorf(template, args)
}

func Println(template string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(template, args...))
}
