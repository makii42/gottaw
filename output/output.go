package output

import (
	n "github.com/0xAX/notificator"
	"github.com/fatih/color"
	c "github.com/makii42/gottaw/config"
	"fmt"
)

type Level int

const (
	L_QUIET  Level = iota
	L_NOTICE
	L_TRACE
)

type Log struct {
	cfg                                        *c.Config
	errors, notices, triggers, success, normal *color.Color
	n                                          *n.Notificator
	level                                      Level
}

type Logger interface {
	Errorf(format string, a ...interface{})
	Noticef(format string, a ...interface{})
	Triggerf(format string, a ...interface{})
	Successf(format string, a ...interface{})
	Tracef(format string, a ...interface{})
}

func (o *Log) growl(title, msg, icon, urgency string) {
	if o.n != nil {
		o.n.Push(title, msg, icon, urgency)
	}
}

func NewLog(trace, quiet bool, cfg *c.Config) (Logger, error) {

	var notificator *n.Notificator
	if cfg.Growl {
		if notificator == nil {
			notificator = makeNotificator()
		}
	} else {
		notificator = nil
	}

	l := Log{
		cfg:      cfg,
		errors:   color.New(color.FgHiRed),
		notices:  color.New(color.FgBlue),
		triggers: color.New(color.FgYellow),
		success:  color.New(color.FgGreen),
		normal:   color.New(color.FgHiWhite),
		n:        notificator,
		level:    lvl,
	}
	return &l, nil
}

func (l *Log) GetLog() Logger {
	return l
}

func (l *Log) Errorf(format string, a ...interface{}) {
	l.errors.Printf(format, a)
	l.growl("Error", fmt.Sprintf(format, a...), "", n.UR_CRITICAL)
}

func (l *Log) Noticef(format string, a ...interface{}) {
	l.notices.Printf(format, a)
}

func (l *Log) Triggerf(format string, a ...interface{}) {
	l.triggers.Printf(format, a)
}

func (l *Log) Successf(format string, a ...interface{}) {
	l.success.Printf(format, a)
	l.growl("Pipeline Success", fmt.Sprintf(format, a...), "", n.UR_NORMAL)
}

func (l *Log) Tracef(format string, a ...interface{}) {
	if l.level >= L_TRACE {
		l.normal.Printf(format, a)
	}
}

func makeNotificator() *n.Notificator {
	return n.New(n.Options{
		AppName: "Gotta Watch!",
	})
}
