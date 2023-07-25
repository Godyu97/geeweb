package gee

import (
	"log"
	"time"
)

type GeeLogger interface {
	Info(a ...any)
	Error(a ...any)
}

var geeLogger GeeLogger = &defaultLog{}

func InitLogger(logger GeeLogger) {
	//todo logger init
	geeLogger = logger
}

type defaultLog struct{}

func (l *defaultLog) Info(a ...any) {
	log.Println(a...)
}

func (l *defaultLog) Error(a ...any) {
	log.Println(a...)
}

func Logger() HandlerFunc {
	return func(ctx *Context) {
		t := time.Now()
		ctx.Next()
		geeLogger.Info(ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}
