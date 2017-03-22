package main

import (
	"fmt"

	n "github.com/0xAX/notificator"
	"github.com/fatih/color"
)

type Logger struct {
	cfg                                *Config
	errors, notices, triggers, success *color.Color
	n                                  *n.Notificator
}

func NewLogger(cfg *Config) *Logger {
	l := Logger{
		cfg:      cfg,
		errors:   color.New(color.FgHiRed),
		notices:  color.New(color.FgBlue),
		triggers: color.New(color.FgYellow),
		success:  color.New(color.FgGreen),
	}
	return &l
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	l.errors.Printf(format, a...)
	l.growl("Pipeline Error", fmt.Sprintf(format, a...), "", n.UR_CRITICAL)
}

func (l *Logger) Noticef(format string, a ...interface{}) {
	l.notices.Printf(format, a...)
}

func (l *Logger) Triggerf(format string, a ...interface{}) {
	l.triggers.Printf(format, a...)
}

func (l *Logger) Successf(format string, a ...interface{}) {
	l.success.Printf(format, a...)
	l.growl("Pipeline Success", fmt.Sprintf(format, a...), "", n.UR_NORMAL)
}

func (l *Logger) growl(title, msg, icon, urgency string) {
	if cfg.Growl {
		if l.n == nil {
			l.n = makeNotificator()
		}
		l.n.Push(title, msg, icon, urgency)
	} else if l.n != nil {
		l.n = nil
	}
}

func makeNotificator() *n.Notificator {
	return n.New(n.Options{
		AppName: "Gotta Watch!",
	})
}
