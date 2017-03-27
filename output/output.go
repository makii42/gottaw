package output

import (
	"fmt"

	n "github.com/0xAX/notificator"
	"github.com/fatih/color"
	c "github.com/makii42/gottaw/config"
)

type Level int

const (
	QUIET Level = iota
	NOTICE
	TRACE
)

type Logger struct {
	cfg                                *c.Config
	errors, notices, triggers, success *color.Color
	n                                  *n.Notificator
	level                              Level
}

func NewLogger(lvl Level, cfg *c.Config) *Logger {
	var notificator *n.Notificator
	if cfg.Growl {
		if notificator == nil {
			notificator = makeNotificator()
		}
	} else {
		notificator = nil
	}

	l := Logger{
		cfg:      cfg,
		errors:   color.New(color.FgHiRed),
		notices:  color.New(color.FgBlue),
		triggers: color.New(color.FgYellow),
		success:  color.New(color.FgGreen),
		n:        notificator,
		level:    lvl,
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

func (l *Logger) Tracef(format string, a ...interface{}) {
	if l.level >= TRACE {
		fmt.Printf(format, a...)
	}
}

func (l *Logger) Traceln(a ...interface{}) (n int, err error) {
	if l.level >= TRACE {
		return fmt.Println(a...)
	}
	return 0, nil
}

func (l *Logger) growl(title, msg, icon, urgency string) {
	if l.n != nil {
		l.n.Push(title, msg, icon, urgency)
	}
}

func makeNotificator() *n.Notificator {
	return n.New(n.Options{
		AppName: "Gotta Watch!",
	})
}
